package test

import (
	"context"
	"testing"

	"github.com/foomo/gocontentful/test/testapi"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NotNil(t, contentfulClient)
	require.NoError(t, err)
	contentfulClient.ClientStats()
	stats, err := contentfulClient.GetCacheStats()
	require.NoError(t, err)
	require.Equal(t, 3, len(stats.ContentTypes))
	require.Equal(t, 12, stats.AssetCount)
	require.Equal(t, 9, stats.EntryCount)
	require.Equal(t, 7, stats.ParentCount)
	err = contentfulClient.SetSyncMode(true)
	require.Error(t, err)
}

func TestBrokenReferences(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	brokenReferences := contentfulClient.BrokenReferences()
	require.Equal(t, 1, len(brokenReferences))
}

func TestCacheHasContentType(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	require.True(t, contentfulClient.CacheHasContentType("brand"))
}

func TestGetAsset(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	_, err = contentfulClient.GetAssetByID(context.TODO(), "Xc0ny7GWsMEMCeASWO2um")
	require.NoError(t, err)
	newAsset := testapi.NewAssetFromURL("12345", "https://example.com", "PNG", "New Asset")
	require.NotNil(t, newAsset)
}

func TestDeleteAssetFromCache(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	stats, err := contentfulClient.GetCacheStats()
	require.NoError(t, err)
	require.Equal(t, 12, stats.AssetCount)
	err = contentfulClient.DeleteAssetFromCache("Xc0ny7GWsMEMCeASWO2um")
	require.NoError(t, err)
	stats, err = contentfulClient.GetCacheStats()
	require.NoError(t, err)
	require.Equal(t, 11, stats.AssetCount)
}

func TestGetContentTypeOfID(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	contentType, err := contentfulClient.GetContentTypeOfID(context.TODO(), "651CQ8rLoIYCeY6G0QG22q")
	require.NoError(t, err)
	require.Equal(t, "brand", contentType)
}

func TestGetParents(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	product, err := contentfulClient.GetProductByID(context.TODO(), "6dbjWqNd9SqccegcqYq224")
	require.NoError(t, err)
	brandRef := product.Brand(context.TODO())
	brandParents, err := brandRef.GetParents(context.TODO())
	require.NoError(t, err)
	require.Equal(t, 2, len(brandParents))
	brandParents, err = brandRef.GetParents(context.TODO(), testapi.ContentTypeProduct)
	require.NoError(t, err)
	require.Equal(t, 2, len(brandParents))
	brandParents, err = brandRef.GetParents(context.TODO(), testapi.ContentTypeCategory)
	require.NoError(t, err)
	require.Equal(t, 0, len(brandParents))
	brandRef.CC = nil
	brandParents, err = brandRef.GetParents(context.TODO())
	require.NoError(t, err)
	require.Equal(t, 2, len(brandParents))
}

func TestCacheIfNewEntry(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	stats, err := contentfulClient.GetCacheStats()
	require.NoError(t, err)
	require.Equal(t, 9, stats.EntryCount)
	testFile, err := GetTestFile("./test-space-export-newer.json")
	require.NoError(t, err)
	err = contentfulClient.SetOfflineFallback(testFile)
	require.NoError(t, err)
	_, _, err = contentfulClient.UpdateCache(context.Background(), nil, false)
	require.NoError(t, err)
	stats, err = contentfulClient.GetCacheStats()
	require.NoError(t, err)
	require.Equal(t, 10, stats.EntryCount)
}

func TestPreserveCacheIfNewer(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	testFile, err := GetTestFile("./test-space-export-older.json")
	require.NoError(t, err)
	err = contentfulClient.SetOfflineFallback(testFile)
	require.NoError(t, err)
	_, _, err = contentfulClient.UpdateCache(context.TODO(), nil, false)
	require.NoError(t, err)
	brand, err := contentfulClient.GetBrandByID(context.TODO(), "JrePkDVYomE8AwcuCUyMi")
	require.NoError(t, err)
	require.Equal(t, 2.0, brand.Sys.Version)
}

func TestEntry(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	cfProduct := testapi.NewCfProduct(contentfulClient)
	err = cfProduct.SetProductName("dummy")
	require.NoError(t, err)
	require.NotNil(t, cfProduct.CC)
}

func TestGenericEntries(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	genericEntries, err := contentfulClient.GetAllGenericEntries()
	require.NoError(t, err)
	require.Equal(t, 9, len(genericEntries))
	genericBrand, err := contentfulClient.GetGenericEntry("651CQ8rLoIYCeY6G0QG22q")
	require.NoError(t, err)
	website, err := genericBrand.FieldAsString("website")
	require.NoError(t, err)
	require.Equal(t, "http://www.normann-copenhagen.com/", website)
	genericProduct, err := contentfulClient.GetGenericEntry("6dbjWqNd9SqccegcqYq224")
	require.NoError(t, err)
	quantity, err := genericProduct.FieldAsFloat64("quantity")
	require.NoError(t, err)
	require.Equal(t, 89.0, quantity)
	err = genericProduct.SetField("quantity", 90.0)
	require.NoError(t, err)
	quantityAny, err := genericProduct.FieldAsAny("quantity")
	require.NoError(t, err)
	require.Equal(t, 90.0, quantityAny.(float64))
	err = genericProduct.SetField("quantity2", 90.0)
	require.NoError(t, err)
	quantity2Any, err := genericProduct.FieldAsAny("quantity2")
	require.NoError(t, err)
	require.Equal(t, 90.0, quantity2Any.(float64))
	productBrand, err := genericProduct.FieldAsReference("brand")
	require.NoError(t, err)
	require.NotNil(t, productBrand)
	require.Equal(t, "651CQ8rLoIYCeY6G0QG22q", productBrand.ID)
	// inherit
	sku, err := genericProduct.FieldAsString("sku")
	require.Error(t, err)
	require.Equal(t, "", sku)
	ctx := context.Background()
	inheritedSKU, err := genericProduct.InheritAsString(ctx, "sku", nil)
	require.NoError(t, err)
	require.Equal(t, "B00MG4ULK2", inheritedSKU)
}
