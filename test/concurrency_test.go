package test

import (
	"context"
	"github.com/foomo/gocontentful/test/testapi"
	"sync"
	"testing"
)

var testProductID = "6dbjWqNd9SqccegcqYq224"

func readWorker(contentfulClient *testapi.ContentfulClient, i int) error {
	product, err := contentfulClient.GetProductByID(testProductID)
	if err != nil {
		return err
	}
	price := product.Price()
	testLogger.Infof("Read worker %d read price: %f", i, price)
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
	for i := 1; i <= 1000; i++ {
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
	for i := 1; i <= 1000; i++ {
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
	for i := 1; i <= 1000; i++ {
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
	wg.Wait()
}
