package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/foomo/gocontentful/test"
	"github.com/foomo/gocontentful/test/testapi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTheAPI(t *testing.T) {
	ctx := context.Background()
	testLogger := logrus.StandardLogger()
	testFile, err := test.GetTestFile("./test/test-space-export.json")
	require.NoError(t, err)
	cc, errClient := testapi.NewOfflineContentfulClient(testFile,
		test.GetContenfulLogger(testLogger),
		test.LogDebug,
		true,
		true)
	require.NoError(t, errClient)
	prods, errProds := cc.GetAllProduct(ctx)
	require.NoError(t, errProds)
	testLogger.WithField("prods", len(prods)).Info("Loaded products")
	// Grab the first product
	prod, errProd := cc.GetProductByID(ctx, "6dbjWqNd9SqccegcqYq224")
	require.NoError(t, errProd)
	prodName := prod.ProductName()
	testLogger.WithField("name", prodName).Info("Product name")
	// Get the brand
	brandReference := prod.Brand(ctx)
	brand := brandReference.VO.(*testapi.CfBrand)
	testLogger.WithField("name", brand.CompanyName()).Info("Brand")
	parentRefs, errParents := brand.GetParents(ctx)
	require.NoError(t, errParents)
	testLogger.WithField("parent count", len(parentRefs)).Info("Parents")
	for _, parentRef := range parentRefs {
		switch parentRef.ContentType {
		case testapi.ContentTypeProduct:
			parentProduct := parentRef.VO.(*testapi.CfProduct)
			testLogger.WithField("name", parentProduct.ProductName()).Info("Parent product")
		}
	}
	testCleanUpUnicode, err := cc.GetProductByID(ctx, "6dbjWqNd9SqccegcqYq224")
	require.NoError(t, err)
	html, err := testapi.RichTextToHtml(testCleanUpUnicode.SeoText(testapi.SpaceLocaleGerman), nil, nil, nil, nil, testapi.SpaceLocaleGerman)
	require.NoError(t, err)
	assert.Equal(t, 2109, len(html))
	// testLogger.Info(html)
	// testLogger.Info(len(html))
	assert.Equal(t, 13, len(testCleanUpUnicode.ProductName()))
	// productName := testCleanUpUnicode.ProductName()
	// testLogger.Info(productName)
	// testLogger.Info(fmt.Printf("%q", productName))
	// testLogger.Info(len(testCleanUpUnicode.ProductName()))
	tags := []int{}
	for _, tag := range testCleanUpUnicode.Tags() {
		testLogger.Info(tag)
		tags = append(tags, len(tag))
	}
	assert.Equal(t, tags, []int{7, 11, 5, 11, 7})
	// testLogger.Info(tags)
}

//func TestReadExport(t *testing.T) {
//	fileBytes, err := ioutil.ReadFile("./test/test-space-export.json")
//	require.NoError(t, err)
//	var export erm.ExportFile
//	err = json.Unmarshal(fileBytes, &export)
//	require.NoError(t, err)
//	fmt.Println(export)
//}
//
//func TestLiveSpaceIsStillWorking(t *testing.T) {
//	testLogger := logrus.StandardLogger()
//	cc, err := testapi.NewContentfulClient("v0ro6q5qmdu8", testapi.ClientModeCDA, "PHvA-6ExNXvYvIIP-p255KWfZhu2eM2OrKHa3nhHSB4", 100, test.GetContenfulLogger(testLogger), test.LogDebug, true)
//	require.NoError(t, err)
//	prods, err := cc.GetAllProduct()
//	require.NoError(t, err)
//	testLogger.WithField("prods", len(prods)).Info("Loaded live products")
//}

func TestCleanInvisibles(t *testing.T) {
	invisibleChars := "\u0009D\u200bougl\u2028as Finn"
	fmt.Printf("\"%s\" is really %q\n", invisibleChars, invisibleChars)
	fmt.Println(len(invisibleChars))

	invisibleChars = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) || unicode.IsControl(r) {
			return r
		}
		return -1
	}, invisibleChars)

	fmt.Printf("IsGraphic: %q\n", invisibleChars)
	fmt.Println(len(invisibleChars))
}
