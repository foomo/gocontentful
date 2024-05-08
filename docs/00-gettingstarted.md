# Getting Started

Before you install Gocontentful as a command-line tool to use it in your projects, we suggest you get a taste of how it works by playing with the test API from the Gocontentful repository. This doesn't yet require you to have access to Contentful.

## How to play with the test API

Clone the gocontentful repository from [https://github.com/foomo/gocontentful](https://github.com/foomo/gocontentful) and open it
in your IDE.

The repository includes an offline representation of a Contentful space that can is used for testing gocontentful
without depending on an online connection and an existing Contentful space.

First, open a terminal and install

```bash
go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
```

Then `cd` to the repository folder and make sure tests run fine on your machine

```bash
make test
```

Create a test file in the repository home directory (`api_test.go` might be a good choice).
Paste the following into the file:

```go
package main

import (
	"context"
    "testing"

	"github.com/foomo/gocontentful/test"
	"github.com/foomo/gocontentful/test/testapi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestTheAPI(t *testing.T) {
  testLogger := logrus.StandardLogger()
  testFile, err := test.GetTestFile("./test/test-space-export.json")
  require.NoError(t, err)
  cc, errClient := testapi.NewOfflineContentfulClient(testFile,
		test.GetContenfulLogger(testLogger),
		test.LogDebug,
		true, false)
  require.NoError(t, errClient)
  prods, errProds := cc.GetAllProduct(context.TODO())
  require.NoError(t, errProds)
  testLogger.WithField("prods", len(prods)).Info("Loaded products")
}
```

The first two lines in the unit test create a logger and an offline gocontentful client. This also
caches the content of the space in memory and lets you play with the API. The space includes three
content types (`brand`, `product` and `category`) and their relative entries.
A product has a reference to a brand and to one or more categories. If you open the `./test/test-space-export.json` file
you can see how the JSON representation of those entries is.

Getting all the products using the Contentful
Content Delivery API would normally require dealing with the connection, query and JSON payload, having
value object defined for all content types and functions to convert from/to those structs. With the Go API generated
by gocontentful all you need to do to load all the products is one single line:

```go
// First get a context, this is needed for all operations
// that potentially require a network connection to Contentful
ctx := context.TODO()
prods, errProds := cc.GetAllProduct(ctx)
```

Open a terminal and from the repository home directory run the test. Your output should looks similar to this:

```shell
$ go test -run TestTheAPI
INFO[0000] loading space from local file                 assets=12 entries=9
INFO[0000] contentful cache update queued                task=UpdateCache
INFO[0000] contentful cache worker starting              task=UpdateCache
INFO[0000] gonna use a local file                        task=UpdateCache
INFO[0000] cached all entries of content type            contentType=product method=updateCacheForContentType size=4
INFO[0000] cached all assets                             contentType=asset method=updateCacheForContentType size=12
INFO[0000] cached all entries of content type            contentType=brand method=updateCacheForContentType size=3
INFO[0000] cached all entries of content type            contentType=category method=updateCacheForContentType size=2
INFO[0000] space caching done, time recorded             task=UpdateCache time elapsed=179.357792ms
INFO[0000] contentful cache update returning             task=UpdateCache
INFO[0000] contentful cache update returning             task=UpdateCache
INFO[0000] Loaded products                               prods=4
PASS
ok      github.com/foomo/gocontentful   0.484s
```

The last line shows that we loaded 4 products. Let's go ahead and play with the API.
We'll load a specific product and log its name. Add this at the end of the unit test:

```go
prod, errProd := cc.GetProductByID(ctx, "6dbjWqNd9SqccegcqYq224")
require.NoError(t, errProd)
prodName := prod.ProductName("de")
testLogger.WithField("name", prodName).Info("Product loaded")
```

This will be the output at the end of the log when you run the test:

```shell
INFO[0000] Product loaded                                name="Whisk Beater"
```

The first line loads the product from the space. This is a `*testapi.CfProduct` pointer. The type is generated
and carries all the getter and setter methods to access all the fields and more, e.g. ProductName().
Note that when calling ProductName() we passed `"de"` as a parameter. This is the locale and it's
entirely optional and useful when your space supports multiple locales for translation.
If you omit it, the default space locale will be used.

Let's load the product's brand:

```go
// Get the brand
brandReference := prod.Brand(ctx)
brand := brandReference.VO.(*testapi.CfBrand)
testLogger.WithField("name", brand.CompanyName()).Info("Brand")
```

Note that:

- The product has a Brand() method that represents and retrieves the reference from the product entry to the brand entry
- The returned object is not a `*testapi.CfBrand` pointer as you might expect. This is because a reference field in Contentful
  can point to entries of multiple content types and that doesn't play nice with Go's static typing.
  The object returned is a generic `*testapi.EntryReference` that, among other, includes an `interface{}` attribute (VO) that
  is the actual `*testapi.CfBrand`. That's why in the second line we had to cast it.

The test now logs the brand company name:

```shell
INFO[0000] Brand                                         name="Normann Copenhagen"
```

What if we want to follow the reference the other way around and find out which entries link to this brand?

```go
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
```

Again, the `GetParents()` method returns references and not objects. It's a good idea to use the reference `ContentType` attribute
to switch before casting the VO to the type, because as we just said referenced objects can come in different types and casting
to the wrong one would make the runtime panic. Running the test we find out the two products that belong to this brand:

```shell
INFO[0000] Parents                                       parent count=2
INFO[0000] Parent product                                name="Whisk Beater"
INFO[0000] Parent product                                name="Hudson Wall Cup"
```
