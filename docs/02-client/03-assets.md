# Assets

Contentful allows upload and reference of binary assets and gocontentful fully supports them.
Assuming the dog entry references a picture in a field you can get it with:

```go
picture := dog.Picture(ctx) // you can pass a locale to this function as usual
```

This returns a \*contenful.AssetNoLocale object handling localization for you in two ways.
First, the field itself could be localized in the model, referencing two different assets altogether.
Secondly, the asset itself can have different files uploaded for different locales.
No matter what, the gocontentful API will return the right file:

```go
// Get the asset's URL at Contentful's CDN
if picture != nil && picture.Fields != nil && picture.Fields.File != nil {
    theURL := picture.Fields.File.URL
    // ...then do something with it
}
```

There are various functions and methods to work with assets, for example to create an asset
starting from an URL or to convert an asset to a reference to store it in a parent entry
field. See the [API Reference](./api-reference) chapter for details.

Note: there is no function to create a new asset in the generated code because the type `AssetNoLocale`
is from the _github.com/foomo/contentful_ package, just instantiate one if you need a blank asset.
