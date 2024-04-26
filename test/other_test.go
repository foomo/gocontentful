package test

import (
	"context"
	"testing"
	"time"

	"github.com/foomo/gocontentful/test/testapi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishingStatus(t *testing.T) {
	contentfulClient, err := getTestClient()
	require.NoError(t, err)
	time.Sleep(time.Second)
	draft, err := contentfulClient.GetProductByID(context.TODO(), "6dbjWqNd9SqccegcqYq224")
	require.NoError(t, err)
	require.Equal(t, testapi.StatusDraft, draft.GetPublishingStatus())
	published, err := contentfulClient.GetCategoryByID(context.TODO(), "7LAnCobuuWYSqks6wAwY2a")
	require.NoError(t, err)
	require.Equal(t, testapi.StatusPublished, published.GetPublishingStatus())
	changed, err := contentfulClient.GetProductByID(context.TODO(), "3DVqIYj4dOwwcKu6sgqOgg")
	require.NoError(t, err)
	require.Equal(t, testapi.StatusChanged, changed.GetPublishingStatus())
}

func TestCleanUpUnicode(t *testing.T) {
	testLogger := logrus.StandardLogger()
	cc, errClient := testapi.NewOfflineContentfulClient("./test-space-export.json",
		GetContenfulLogger(testLogger),
		LogDebug,
		true,
		true)
	require.NoError(t, errClient)
	testCleanUpUnicode, err := cc.GetProductByID(context.TODO(), "6dbjWqNd9SqccegcqYq224")
	require.NoError(t, err)
	html, err := testapi.RichTextToHtml(testCleanUpUnicode.SeoText(testapi.SpaceLocaleGerman), nil, nil, nil, nil, testapi.SpaceLocaleGerman)
	require.NoError(t, err)
	assert.Equal(t, 2109, len(html))
	assert.Equal(t, 13, len(testCleanUpUnicode.ProductName()))
	tags := []int{}
	for _, tag := range testCleanUpUnicode.Tags() {
		testLogger.Info(tag)
		tags = append(tags, len(tag))
	}
	assert.Equal(t, tags, []int{7, 11, 5, 11, 7})
}
