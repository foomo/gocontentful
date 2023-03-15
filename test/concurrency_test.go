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

func readWorker(contentfulClient *testapi.ContentfulClient, i int) error {
	product, err := contentfulClient.GetProductByID(testProductID)
	if err != nil {
		return err
	}
	price := product.Price()
	testLogger.Infof("Read worker %d read price: %f", i, price)
	return nil
}

func parentWorker(contentfulClient *testapi.ContentfulClient, i int) error {
	brand, err := contentfulClient.GetBrandByID(testBrandID)
	if err != nil {
		return err
	}
	parents, err := brand.GetParents()
	if err != nil {
		return err
	}
	testLogger.Infof("Parent worker %d found %d brand parents", i, len(parents))
	return nil
}

func writeWorker(contentfulClient *testapi.ContentfulClient, i int) error {
	product, err := contentfulClient.GetProductByID(testProductID)
	if err != nil {
		return err
	}
	err = product.SetPrice(float64(i))
	if err != nil {
		return err
	}
	contentfulClient.SetProductInCache(product)
	testLogger.Infof("Write worker %d set price: %d", i, i)
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
			err := contentfulClient.UpdateCache(context.TODO(), nil, false)
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
			err := writeWorker(contentfulClient, i)
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
			err := readWorker(contentfulClient, i)
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
			err := parentWorker(contentfulClient, i)
			if err != nil {
				testLogger.Errorf("testConcurrentReadWrites: %v", err)
			}
		}()
	}
	wg.Wait()
}
