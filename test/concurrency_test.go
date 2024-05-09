package test

import (
	"context"
	"sync"
	"testing"

	"github.com/foomo/gocontentful/test/testapi"
)

var (
	testProductID = "6dbjWqNd9SqccegcqYq224"
	testBrandID   = "651CQ8rLoIYCeY6G0QG22q"
	concurrency   = 10000
)

func readWorker(ctx context.Context, contentfulClient *testapi.ContentfulClient, i int) error {
	product, err := contentfulClient.GetProductByID(ctx, testProductID)
	if err != nil {
		return err
	}
	_, err = contentfulClient.GetAllProduct(ctx)
	if err != nil {
		return err
	}
	price := product.Price()
	testLogger.Infof("Read worker %d read price: %f", i, price)
	_ = product.Brand(ctx)
	_ = product.Categories(ctx)
	_ = product.Image(ctx)
	_ = product.Nodes()
	_ = product.ProductDescription()
	_ = product.ProductName()
	_ = product.Quantity()
	_ = product.SeoText()
	_ = product.Sizetypecolor()
	_ = product.Sku()
	_ = product.Slug()
	_ = product.Tags()
	_ = product.Website()
	_ = product.GetPublishingStatus()
	return nil
}

func parentWorker(ctx context.Context, contentfulClient *testapi.ContentfulClient, i int) error {
	brand, err := contentfulClient.GetBrandByID(ctx, testBrandID)
	if err != nil {
		return err
	}
	parents, err := brand.GetParents(ctx)
	if err != nil {
		return err
	}
	testLogger.Infof("Parent worker %d found %d brand parents", i, len(parents))
	return nil
}

func writeWorker(ctx context.Context, contentfulClient *testapi.ContentfulClient, i int) error {
	product, err := contentfulClient.GetProductByID(ctx, testProductID)
	if err != nil {
		return err
	}
	err = product.SetPrice(float64(i))
	if err != nil {
		return err
	}
	contentfulClient.SetProductInCache(product)
	testLogger.Infof("Write worker %d set price: %d", i, i)
	_ = product.SetBrand(testapi.ContentTypeSys{
		Sys: testapi.ContentTypeSysAttributes{
			ID:       "651CQ8rLoIYCeY6G0QG22q",
			Type:     "Link",
			LinkType: "Entry",
		},
	})
	_ = product.SetCategories([]testapi.ContentTypeSys{
		{Sys: testapi.ContentTypeSysAttributes{
			ID:       "7LAnCobuuWYSqks6wAwY2a",
			Type:     "Link",
			LinkType: "Entry",
		}},
	})
	_ = product.SetImage([]testapi.ContentTypeSys{
		{Sys: testapi.ContentTypeSysAttributes{
			ID:       "10TkaLheGeQG6qQGqWYqUI",
			Type:     "Link",
			LinkType: "Asset",
		}},
	})
	_ = product.SetNodes(nil)
	_ = product.SetProductDescription("")
	_ = product.SetProductName("")
	_ = product.SetQuantity(1)
	_ = product.SetSeoText("")
	_ = product.SetSizetypecolor("")
	_ = product.SetSku("")
	_ = product.SetSlug("")
	_ = product.SetTags([]string{""})
	_ = product.SetWebsite("")
	return nil
}

func TestConcurrentReadWrites(t *testing.T) {
	contentfulClient, err := getTestClient()
	if err != nil {
		testLogger.Errorf("testConcurrentReadWrites: %v", err)
	}
	var wg sync.WaitGroup
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			testLogger.Infof("testConcurrentReadWrites: caching run %d", i)
			_, _, err := contentfulClient.UpdateCache(context.TODO(), nil, false)
			if err != nil {
				testLogger.Errorf("testConcurrentReadWrites: %v", err)
			}
		}()
	}
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			err := writeWorker(context.TODO(), contentfulClient, i)
			if err != nil {
				testLogger.Errorf("testConcurrentReadWrites: %v", err)
			}
		}()
	}
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			err := readWorker(context.TODO(), contentfulClient, i)
			if err != nil {
				testLogger.Errorf("testConcurrentReadWrites: %v", err)
			}
		}()
	}
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			err := parentWorker(context.TODO(), contentfulClient, i)
			if err != nil {
				testLogger.Errorf("testConcurrentReadWrites: %v", err)
			}
		}()
	}
	wg.Wait()
}
