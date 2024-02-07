package test

import (
	"context"
	"testing"

	"github.com/foomo/gocontentful/test/testapi"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	contentfulClient, err := getTestClient()
	contentfulClient.ClientStats()
	require.NoError(t, err)
	stats, err := contentfulClient.GetCacheStats()
	require.NoError(t, err)
	require.Equal(t, 3, len(stats.ContentTypes))
	require.Equal(t, 12, stats.AssetCount)
	require.Equal(t, 9, stats.EntryCount)
	require.Equal(t, 6, stats.ParentCount)
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
	_, err = contentfulClient.GetAssetByID("Xc0ny7GWsMEMCeASWO2um")
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
	contentType, err := contentfulClient.GetContentTypeOfID("651CQ8rLoIYCeY6G0QG22q")
	require.NoError(t, err)
	require.Equal(t, "brand", contentType)
}

func TestGetParents(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	product, err := contentfulClient.GetProductByID("6dbjWqNd9SqccegcqYq224")
	require.NoError(t, err)
	brandRef := product.Brand()
	brandParents, err := brandRef.GetParents(contentfulClient)
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
	err = contentfulClient.UpdateCache(context.Background(), nil, false)
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
	err = contentfulClient.UpdateCache(context.TODO(), nil, false)
	require.NoError(t, err)
	brand, err := contentfulClient.GetBrandByID("JrePkDVYomE8AwcuCUyMi")
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
	require.Equal(t, "http://www.normann-copenhagen.com/", genericBrand.FieldAsString("website"))
	genericProduct, err := contentfulClient.GetGenericEntry("6dbjWqNd9SqccegcqYq224")
	require.NoError(t, err)
	require.Equal(t, 89.0, genericProduct.FieldAsFloat64("quantity"))
}
