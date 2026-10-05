package test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/foomo/gocontentful/test/testapi"
	"github.com/stretchr/testify/require"
)

const (
	regressedProductID = "6dbjWqNd9SqccegcqYq224"
	productImageID     = "10TkaLheGeQG6qQGqWYqUI"
	originalBrandID    = "651CQ8rLoIYCeY6G0QG22q"
	otherBrandID       = "4LgMotpNF6W20YKmuemW0a"
)

// exportVariant returns the test space export with edit applied to its entries and assets,
// indexed by ID. Entries and assets deleted from the maps are dropped from the export.
func exportVariant(t *testing.T, edit func(entries, assets map[string]map[string]any)) []byte {
	t.Helper()
	raw, err := GetTestFile("./test-space-export.json")
	require.NoError(t, err)
	export := map[string]any{}
	require.NoError(t, json.Unmarshal(raw, &export))
	index := func(key string) (map[string]map[string]any, []string) {
		items := map[string]map[string]any{}
		var order []string
		for _, item := range export[key].([]any) {
			typed := item.(map[string]any)
			id := typed["sys"].(map[string]any)["id"].(string)
			items[id] = typed
			order = append(order, id)
		}
		return items, order
	}
	entries, entryOrder := index("entries")
	assets, assetOrder := index("assets")
	edit(entries, assets)
	keep := func(items map[string]map[string]any, order []string) []any {
		var kept []any
		for _, id := range order {
			if item, ok := items[id]; ok {
				kept = append(kept, item)
			}
		}
		return kept
	}
	export["entries"] = keep(entries, entryOrder)
	export["assets"] = keep(assets, assetOrder)
	out, err := json.Marshal(export)
	require.NoError(t, err)
	return out
}

func entrySys(entry map[string]any) map[string]any {
	return entry["sys"].(map[string]any)
}

// cachedExport is the space as first cached: the product is at published version 5.
func cachedExport(t *testing.T) []byte {
	t.Helper()
	return exportVariant(t, func(entries, _ map[string]map[string]any) {
		entrySys(entries[regressedProductID])["publishedVersion"] = 5.0
	})
}

// regressedExport is the space after an environment reset: the product comes back at the older
// published version 3, with another name and brand, while another brand moves forward normally.
func regressedExport(t *testing.T) []byte {
	t.Helper()
	return exportVariant(t, func(entries, _ map[string]map[string]any) {
		product := entries[regressedProductID]
		entrySys(product)["publishedVersion"] = 3.0
		entrySys(product)["updatedAt"] = "2019-01-01T00:00:00.000Z"
		fields := product["fields"].(map[string]any)
		fields["productName"] = map[string]any{"de": "Older whisk"}
		fields["brand"] = map[string]any{"de": map[string]any{"sys": map[string]any{"type": "Link", "linkType": "Entry", "id": otherBrandID}}}
		brand := entries[otherBrandID]
		entrySys(brand)["publishedVersion"] = 2.0
		brand["fields"].(map[string]any)["companyName"] = map[string]any{"de": "Lemnos updated"}
	})
}

func newClientFromExport(t *testing.T, export []byte) *testapi.ContentfulClient {
	t.Helper()
	cc, err := testapi.NewOfflineContentfulClient(export, GetContenfulLogger(testLogger), LogDebug, true, true)
	require.NoError(t, err)
	return cc
}

func parentIDs(t *testing.T, cc *testapi.ContentfulClient, brandID string) []string {
	t.Helper()
	brand, err := cc.GetBrandByID(context.Background(), brandID)
	require.NoError(t, err)
	parents, err := brand.GetParents(context.Background())
	require.NoError(t, err)
	var ids []string
	for _, parent := range parents {
		ids = append(ids, parent.ID)
	}
	return ids
}

func TestUpdateCacheRetainsEntriesWithOlderVersions(t *testing.T) {
	ctx := context.Background()
	cc := newClientFromExport(t, cachedExport(t))
	product, err := cc.GetProductByID(ctx, regressedProductID)
	require.NoError(t, err)
	cachedName := product.ProductName()

	require.NoError(t, cc.SetOfflineFallback(regressedExport(t), 0))
	result, err := cc.UpdateCacheWithResult(ctx, nil, true)
	require.NoError(t, err)
	require.Equal(t, testapi.CacheUpdateModeRefresh, result.Mode)
	require.True(t, result.Degraded())
	require.Empty(t, result.Replaced)
	require.Len(t, result.Retained, 1)
	retained := result.Retained[0]
	require.Equal(t, testapi.ContentTypeProduct, retained.ContentType)
	require.Equal(t, regressedProductID, retained.EntryID)
	require.Equal(t, 5.0, retained.CachedVersion)
	require.Equal(t, 3.0, retained.IncomingVersion)
	require.Equal(t, "2020-01-15T10:09:36.578Z", retained.CachedUpdatedAt)
	require.Equal(t, "2019-01-01T00:00:00.000Z", retained.IncomingUpdatedAt)
	require.False(t, retained.RetainedSince.IsZero())

	// The retained entry, its generic entry and its references all come from the cached copy.
	product, err = cc.GetProductByID(ctx, regressedProductID)
	require.NoError(t, err)
	require.Equal(t, 5.0, product.Sys.PublishedVersion)
	require.Equal(t, cachedName, product.ProductName())
	genericProduct, err := cc.GetGenericEntry(regressedProductID)
	require.NoError(t, err)
	require.Equal(t, 5.0, genericProduct.Sys.PublishedVersion)
	genericName, err := genericProduct.FieldAsString("productName")
	require.NoError(t, err)
	require.Equal(t, cachedName, genericName)
	require.Contains(t, parentIDs(t, cc, originalBrandID), regressedProductID)
	require.NotContains(t, parentIDs(t, cc, otherBrandID), regressedProductID)

	// Everything else is refreshed.
	brand, err := cc.GetBrandByID(ctx, otherBrandID)
	require.NoError(t, err)
	require.Equal(t, "Lemnos updated", brand.CompanyName())

	// The entry stays retained, and keeps the time it was first retained, until a reset.
	again, err := cc.UpdateCacheWithResult(ctx, nil, true)
	require.NoError(t, err)
	require.Len(t, again.Retained, 1)
	require.Equal(t, retained.RetainedSince, again.Retained[0].RetainedSince)
	_, _, err = cc.UpdateCache(ctx, nil, true)
	require.NoError(t, err)
}

func TestForceUpdateCacheReplacesRetainedEntries(t *testing.T) {
	ctx := context.Background()
	cc := newClientFromExport(t, cachedExport(t))
	require.NoError(t, cc.SetOfflineFallback(regressedExport(t), 0))
	result, err := cc.UpdateCacheWithResult(ctx, nil, true)
	require.NoError(t, err)
	require.True(t, result.Degraded())

	result, err = cc.ForceUpdateCache(ctx, nil, true)
	require.NoError(t, err)
	require.Equal(t, testapi.CacheUpdateModeReset, result.Mode)
	require.False(t, result.Degraded())
	require.Len(t, result.Replaced, 1)
	replaced := result.Replaced[0]
	require.Equal(t, regressedProductID, replaced.EntryID)
	require.Equal(t, 5.0, replaced.CachedVersion)
	require.Equal(t, 3.0, replaced.IncomingVersion)
	require.True(t, replaced.RetainedSince.IsZero())

	product, err := cc.GetProductByID(ctx, regressedProductID)
	require.NoError(t, err)
	require.Equal(t, 3.0, product.Sys.PublishedVersion)
	require.Equal(t, "Older whisk", product.ProductName())
	require.Contains(t, parentIDs(t, cc, otherBrandID), regressedProductID)
	require.NotContains(t, parentIDs(t, cc, originalBrandID), regressedProductID)

	// After the reset the cache and Contentful agree again.
	result, err = cc.UpdateCacheWithResult(ctx, nil, true)
	require.NoError(t, err)
	require.False(t, result.Degraded())
	require.Empty(t, result.Retained)
}

func TestForceUpdateCacheValidatesSnapshot(t *testing.T) {
	ctx := context.Background()
	cc, err := getTestClient()
	require.NoError(t, err)

	empty := exportVariant(t, func(entries, _ map[string]map[string]any) {
		for id := range entries {
			delete(entries, id)
		}
	})
	require.NoError(t, cc.SetOfflineFallback(empty, 0))
	_, err = cc.ForceUpdateCache(ctx, nil, true)
	require.Error(t, err)
	stats, err := cc.GetCacheStats()
	require.NoError(t, err)
	require.Equal(t, 9, stats.EntryCount)

	// Legitimate shrinkage passes.
	shrunk := exportVariant(t, func(entries, _ map[string]map[string]any) {
		for id := range entries {
			if id != originalBrandID && id != otherBrandID {
				delete(entries, id)
			}
		}
	})
	require.NoError(t, cc.SetOfflineFallback(shrunk, 0))
	_, err = cc.ForceUpdateCache(ctx, nil, true)
	require.NoError(t, err)
	stats, err = cc.GetCacheStats()
	require.NoError(t, err)
	require.Equal(t, 2, stats.EntryCount)
}

func TestUpdateCacheReturnsRebuildErrors(t *testing.T) {
	ctx := context.Background()
	cc := newClientFromExport(t, cachedExport(t))
	broken := exportVariant(t, func(entries, _ map[string]map[string]any) {
		entries[regressedProductID]["fields"].(map[string]any)["price"] = map[string]any{"de": "not a number"}
	})
	require.NoError(t, cc.SetOfflineFallback(broken, 0))
	_, _, err := cc.UpdateCache(ctx, nil, true)
	require.Error(t, err)
	_, err = cc.ForceUpdateCache(ctx, nil, true)
	require.Error(t, err)

	// The previous cache is untouched.
	product, err := cc.GetProductByID(ctx, regressedProductID)
	require.NoError(t, err)
	require.Equal(t, 5.0, product.Sys.PublishedVersion)
}

func TestUpdateCacheReturnsContextErrors(t *testing.T) {
	cc, err := getTestClient()
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err = cc.UpdateCache(ctx, nil, true)
	require.True(t, errors.Is(err, context.Canceled))
	_, err = cc.ForceUpdateCache(ctx, nil, true)
	require.True(t, errors.Is(err, context.Canceled))
}

func TestForceUpdateCacheIsNotServedByAnOrdinaryRefresh(t *testing.T) {
	ctx := context.Background()
	for range 10 {
		cc := newClientFromExport(t, cachedExport(t))
		require.NoError(t, cc.SetOfflineFallback(regressedExport(t), 0))
		var wg sync.WaitGroup
		for range 20 {
			wg.Go(func() {
				result, err := cc.UpdateCacheWithResult(ctx, nil, true)
				require.NoError(t, err)
				require.Contains(t, []testapi.CacheUpdateMode{testapi.CacheUpdateModeRefresh, testapi.CacheUpdateModeReset}, result.Mode)
			})
		}
		var forced *testapi.CacheUpdateResult
		wg.Go(func() {
			var err error
			forced, err = cc.ForceUpdateCache(ctx, nil, true)
			require.NoError(t, err)
		})
		wg.Wait()
		require.Equal(t, testapi.CacheUpdateModeReset, forced.Mode)
		product, err := cc.GetProductByID(ctx, regressedProductID)
		require.NoError(t, err)
		require.Equal(t, 3.0, product.Sys.PublishedVersion)
	}
}

func TestRetainedEntriesKeepTheirDependencies(t *testing.T) {
	ctx := context.Background()
	cc := newClientFromExport(t, cachedExport(t))
	// The incoming snapshot regresses the product and lacks its brand and image.
	require.NoError(t, cc.SetOfflineFallback(exportVariant(t, func(entries, assets map[string]map[string]any) {
		entrySys(entries[regressedProductID])["publishedVersion"] = 3.0
		delete(entries, originalBrandID)
		delete(assets, productImageID)
	}), 0))
	result, err := cc.UpdateCacheWithResult(ctx, nil, true)
	require.NoError(t, err)
	require.Len(t, result.Retained, 1)

	product, err := cc.GetProductByID(ctx, regressedProductID)
	require.NoError(t, err)
	brand := product.Brand(ctx)
	require.NotNil(t, brand)
	require.Equal(t, originalBrandID, brand.ID)
	require.Len(t, product.Image(ctx), 1)
	require.Contains(t, parentIDs(t, cc, originalBrandID), regressedProductID)
	require.Contains(t, result.RetainedDependencies, testapi.RetainedDependency{SysType: "Entry", ContentType: testapi.ContentTypeBrand, ID: originalBrandID, RequiredBy: regressedProductID})
	require.Contains(t, result.RetainedDependencies, testapi.RetainedDependency{SysType: "Asset", ID: productImageID, RequiredBy: regressedProductID})

	// A reset drops the dependencies together with the retained entry.
	result, err = cc.ForceUpdateCache(ctx, nil, true)
	require.NoError(t, err)
	require.Empty(t, result.RetainedDependencies)
	_, err = cc.GetBrandByID(ctx, originalBrandID)
	require.Error(t, err)
}

func TestRetainingEntriesDoesNotRaceWithSetters(t *testing.T) {
	ctx := context.Background()
	cc := newClientFromExport(t, cachedExport(t))
	require.NoError(t, cc.SetOfflineFallback(regressedExport(t), 0))
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		for {
			select {
			case <-stop:
				return
			default:
			}
			product, err := cc.GetProductByID(ctx, regressedProductID)
			if err != nil {
				continue
			}
			_ = product.SetBrand(testapi.ContentTypeSys{Sys: testapi.ContentTypeSysAttributes{ID: otherBrandID, Type: "Link", LinkType: "Entry"}})
			_ = product.SetProductName("edited in memory")
		}
	})
	for range 20 {
		_, err := cc.UpdateCacheWithResult(ctx, nil, true)
		require.NoError(t, err)
	}
	close(stop)
	wg.Wait()
}
