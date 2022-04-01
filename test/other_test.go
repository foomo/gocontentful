package test

import (
	"github.com/foomo/gocontentful/test/testapi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublishingStatus(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	draft, err := contentfulClient.GetProductByID("6dbjWqNd9SqccegcqYq224")
	require.NoError(t, err)
	require.Equal(t, testapi.StatusDraft, draft.GetPublishingStatus())
	published, err := contentfulClient.GetCategoryByID("7LAnCobuuWYSqks6wAwY2a")
	require.NoError(t, err)
	require.Equal(t, testapi.StatusPublished, published.GetPublishingStatus())
	changed, err := contentfulClient.GetProductByID("3DVqIYj4dOwwcKu6sgqOgg")
	require.NoError(t, err)
	require.Equal(t, testapi.StatusChanged, changed.GetPublishingStatus())
}
