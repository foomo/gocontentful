# API Reference

## Client and cache

```go
NewContentfulClient(
	spaceID string,
	clientMode string,
	clientKey string,
	optimisticPageSize uint16,
	logFn func(fields map[string]interface{}, level int, args ...interface{}),
	logLevel int,
	debug bool,
	) (*ContentfulClient, error)
```

Creates a Contentful client, [read this](02-client/01-basicclientoperations) for an explanation of all parameters.

```go
SetOfflineFallback(filename string) error
```

Sets a path to a space export JSON file to be used as a fallback in case
Contentful is not reachable when you call UpdateCache() on the client. This ensures availability
but can make your content look outdated if the export file is older (and typically it is).

```go
NewOfflineContentfulClient(
	file []byte,
	logFn func(fields map[string]interface{}, level int, args ...interface{}),
	logLevel int,
	cacheAssets bool,
	textJanitor bool,
	) (*ContentfulClient, error)
```

Creates an offline Contentful client that loads space data from a JSON file containing a space export.

```go
(cc *ContentfulClient) SetEnvironment(environment string)
```

Sets the Contentful client's environment. All subsequent API calls will be directed to that environment in the selected
space. Pass an empty string to reset to the _master_ environment.

```go
(cc *ContentfulClient) CacheHasContentType(contentTypeID string) bool
```

Returns true if the specified contentTypeID is cached by the client, false otherwise.

```go
(cc *ContentfulClient) SetCacheUpdateTimeout(seconds int64)
```

Sets the cache update timeout to the specified length. A new client by default times out
caching in 120 seconds. A timeout is used to prevent deadlocks when a service panics and recovers
while the gocontentful goroutines are running and the main caching job is waiting for all
them to finish.

```go
(cc *ContentfulClient) SetSyncMode(mode bool) error
```

Switches on/off the cache sync mode. This method will return an error if called on an offline client.

```go
(cc *ContentfulClient) ResetSync()
```

Resets the sync token: the next call to UpdateCache() will rebuild the cache from scratch.

```go
(cc *ContentfulClient) UpdateCache(ctx context.Context, contentTypes []string, cacheAssets bool) error
```

Builds or re-builds the entire client cache.

```go
(cc *ContentfulClient) UpdateCacheForEntity(ctx context.Context, sysType string, contentType string, entityID string) error
```

Updates a single entry or asset (the sysType can take const sysTypeEntry or sysTypeAsset values) in the cache.

```go
(cc *ContentfulClient) ClientMode() ClientMode
```

Returns the internal client mode. There are three constants defined in the generated API:

```go
const (
	ClientModeCDA ClientMode = "CDA"
	ClientModeCPA ClientMode = "CPA"
	ClientModeCMA ClientMode = "CMA"
)
```

## Content functions and methods

_For these we're assuming a content type named "Person"._

```go
NewCfPerson(contentfulClient ...*ContentfulClient) (cfPerson *CfPerson)
```

Creates a new Person entry. You can manipulate and upsert this later. The contentfulClient parameter is optional but you
might want to pass it most of the times or you won't be able to save the entry.

```go
(cc *ContentfulClient) GetAllPerson() (voMap map[string]*CfPerson, err error)
```

Retrieves all Person entries from the client and returnes a map where the key is the ID of the entry and the value is
the Go value object for that entry.

```go
(cc *ContentfulClient) GetFilteredPerson(query *contentful.Query) (voMap map[string]*CfPerson, err error)
```

Retrieves Person entries matching the specified query.

```go
(cc *ContentfulClient) GetPersonByID(id string, forceNoCache ...bool) (vo *CfPerson, err error)
```

Retrieves the Person entry with the specified ID. The optional _forceNoCache_ parameter, if true,
makes the function bypass the existing cache and load a fresh copy of the entry from Contentful.

```go
(ref ContentfulReferencedEntry) ContentType() (contentType string)

```

Returns the Sys.ID of the content type of the referenced entry

```go
(cc *ContentfulClient) GetContentTypeOfID(ID string) (contentType string)
```

Returns the Contentful content type of an entry ID.

```go
(vo *CfPerson) ToReference() (refSys ContentTypeSys)
```

Converts a value object into a reference that can be added to a reference field of an entry. Note that functions that
retrieve referenced entries return a more flexible and useful _[]\*EntryReference_ (see Quickstart above) but to store
a reference you need a ContentTypeSys.

```go
(vo *CfPerson) GetParents(ctx context.Context, contentType ...string) (parents []EntryReference, err error)
(ref *EntryReference) GetParents(ctx context.Context, contentType ...string) (parents []EntryReference, err error)
```

Return a slice of EntryReference objects that represent entries that reference the value object or the entry reference.

```go
HasAncestor(ctx context.Context, contentType string, entry entryOrRef, visited map[string]bool) (*EntryReference, error)
```

Returns the ancestor, if any, of the contentType specified for the entryOrRef. The function travels back in the references graph until it finds an ancestor, or no parent, or hits an infinite loop condition.

```go
(vo *CfPerson) GetPublishingStatus() string
```

Returns the publishing status of the entry as per the Contentful editor UI.
Value returned is one of the following:

```go
const (
  StatusDraft     = "draft"
  StatusChanged   = "changed"
  StatusPublished = "published"
)
```

## Entry field getters and setters

Field getters are named after the field ID in Contentful and return the proper type. For example, if the Person content
type has a Symbol (short text) field named 'Name', this will be the getter:

```go
(vo *CfPerson) Name(locale ...string) (string)
```

The locale parameter is optional and if not passed, the function will return the value for the default locale of the
space. If the locale is specified and it's not available for the space, an error is returned. If the locale is valid
but a value doesn't exist for the field and locale, the function will return the value for the default locale if that's
specified as a fallback locale in the space definition in Contentful, otherwise will return an error.

Possible return types are:

- _string_ for fields of types Symbol, Text, Date
- _[]string_ for fields of type List
- _float64_ for fields of type Integer or Number
- _bool_ for fields of type Boolean
- _\*ContentTypeSys_ for single reference fields
- _[]\*ContentTypeSys_ for multiple reference fields
- _\*ContentTypeFieldLocation_ for fields of type Location
- \*interface{} for fields of type Object or RichText

If logLevel is set to LogDebug retrieving the value of a field that is not set and so not available in the API response
even as a fallback to the default locale will log the event. This can become incredibly verbose, use with care.

Field setters are named after the field ID in Contentful and require to pass in the proper type. See FIELD GETTERS above
for a reference. Example:

```go
(vo *CfPerson) SetName(title string, locale ...string) (err error)
```

## Entry write ops (only available for _ClientModeCMA_)

```go
(vo *CfPerson) UpsertEntry(cc *ContentfulClient) (err error)
```

Upserts the entry. This will appear as "Draft" (if it's a new entry) or "Changed" if it's already existing. In the
latter case, you will need to retrieve the entry with one of the Manage\* functions above to acquire the Sys object
that contains the version information. Otherwise the API call will fail with a "Version mismatch" error.

```go
(vo *CfPerson) PublishEntry(cc *ContentfulClient) (err error)
```

Publishes the entry. Note that before publishing you will need to retrieve the entry with one of the Manage\* functions
above to acquire the Sys object that contains the version information. Otherwise the API call will fail with a "Version
mismatch" error. This is needed even if you have just upserted the entry with the function above!

```go
(vo *CfPerson) UnpublishEntry(cc *ContentfulClient) (err error)
```

Unpublishes the entry. Note that before unpublishing you will need to retrieve the entry with one of the Manage\*
functions above to acquire the Sys object that contains the version information. Otherwise the API call will fail with
a "Version mismatch" error. This is needed even if you have just upserted the entry with the function above!

```go
(vo *CfPerson) UpdateEntry(cc *ContentfulClient) (err error)
```

First upserts the entry and then it publishes it only if it was already published before upserting.  
The rationale is to respect the publishing status of entries and prevent unexpected go-live of content.
Note that before calling this you will need to retrieve theentry with one of the Manage\* functions above to acquire the Sys object that contains the version information. 
Otherwise the API call will fail with a "Version mismatch" error. Using this shortcut function avoids retrieving the entry twice.

```go
(vo *CfPerson) DeleteEntry(cc *ContentfulClient) (err error)
```

Unpublishes and deletes the entry

### Generic entries

```go
(cc *ContentfulClient) GetGenericEntry(entryID string) (*GenericEntry, error)
```

Retrieves a generic entry by ID

```go
(cc *ContentfulClient) GetAllGenericEntries() (map[string]*GenericEntry, error)
```

Retrieves all generic entries and returns a map where the key is the entry ID.

```go
(genericEntry *GenericEntry) FieldAsString(fieldName string, locale ...Locale) (string, error)
```

Returns the specified raw field as a string for the given locale.

```go
(genericEntry *GenericEntry) InheritAsString(fieldName string, locale ...Locale) (string, error)
```

Returns specified raw field from the entry's first parent's, if any, as a string for the given locale.

```go
(genericEntry *GenericEntry) FieldAsFloat64(fieldName string, locale ...Locale) (float64, error)
```

Returns the specified raw field as a float64 for the given locale.

```go
(genericEntry *GenericEntry) InheritAsFloat64(fieldName string, locale ...Locale) (float64, error)
```

Returns the specified raw field from the entry's first parent's, if any, as a float64 for the given locale.

```go
(genericEntry *GenericEntry) FieldAsReference(fieldName string, locale ...Locale) (*EntryReference, error)
```

Returns the specified raw field as a \*EntryReference for the given locale.

```go
(genericEntry *GenericEntry) InheritAsReference(fieldName string, locale ...Locale) (*EntryReference, error)
```

Returns the specified raw field from the entry's first parent's, if any, as a \*EntryReference for the given locale.

```go
(genericEntry *GenericEntry) SetField(fieldName string, fieldValue interface{}, locale ...Locale) error
```

Sets a generic entry's field value.

```go
(genericEntry *GenericEntry) Upsert() error
```

Upserts the generic entry to the space it came from.

```go
(genericEntry *GenericEntry) Update(ctx context.Context) (err error)
```

Upserts the generic entry and publishes it only if it was already published before upserting. Only available for
ClientModeCMA. Before calling this you should retrieve the entry to acquire the Sys version; otherwise the API may fail
with a "Version mismatch" error.

```go
(genericEntry *GenericEntry) GetPublishingStatus() string
```

Returns the publishing status of the entry as per the Contentful editor UI. The value is one of `StatusDraft`,
`StatusChanged`, or `StatusPublished`.

### Asset functions

```go
(cc *ContentfulClient) DeleteAsset(asset *contentful.Asset) error
```

Deletes an asset from the space (only available in CMA)

```go
(cc *ContentfulClient) DeleteAssetFromCache(key string) error {
```

Deletes an asset from the client's cache

```go
(cc *ContentfulClient) GetAllAssets() (map[string]*contentful.Asset, error)
```

Retrieve all assets from a space

```go
(cc *ContentfulClient) GetAssetByID(id string, forceNoCache ...bool) (*contentful.Asset, error)
```

Retrieve an asset from a space by its ID. The optional _forceNoCache_ parameter, if true,
makes the function bypass the existing cache and load a fresh copy of the asset from Contentful.

```go
NewAssetFromURL(id string, uploadUrl string, imageFileType string, title string, locale ...string) *contentful.Asset
```

Creates an Asset from an URL of an existing file online (you still need to upsert it later).

```go
ToAssetReference(asset *contentful.Asset) (refSys ContentTypeSys)
```

Converts the asset to a reference. You need to do this before you add the asset to a reference field of an entry.

```go
(cc *ContentfulClient) UpsertAsset(ctx context.Context, asset *contentful.Asset) error
```

Upserts an asset into the space. Only available for ClientModeCMA. Normalizes file URLs by removing the `https:` prefix
for each locale file, and tolerates idempotency errors returned by the SDK.

```go
(cc *ContentfulClient) PublishAsset(ctx context.Context, asset *contentful.Asset) error
```

Publishes an asset. Only available for ClientModeCMA. The call is idempotent and tolerates "Not published" responses
from the SDK.

```go
(cc *ContentfulClient) UpdateAsset(ctx context.Context, asset *contentful.Asset) (err error)
```

First upserts the asset and then publishes it only if it was already published before upserting. This respects the
current publishing status and avoids unintended go-live of assets. Only available for ClientModeCMA.

```go
(cc *ContentfulClient) DeleteAsset(asset *contentful.Asset) error
```

Deletes an asset from a space by its ID (only available for _ClientModeCMA_)

### Other helper functions and methods

```go
(cc *ContentfulClient) BrokenReferences() (brokenReferences []BrokenReference)
```

Returns a slice of BrokenReference objects with details of where entries have been
referenced but they are not found in the cache. This might naturally return false
positives for content types that are in the space but not included in the cache.

```go
FieldToObject(jsonField interface{}, targetObject interface{}) error
```

Converts a JSON field into an object. Make sure you pass a pointer to an object which type has JSON definition for all
fields you want to retrieve.

```go
HtmlToRichText(htmlSrc string) *RichTextNode
```

Converts an HTML fragment to a RichTextNode. This is useful to migrate data from third-party systems to Contentful or
support HTML paste operations in Web applications. It currently supports headings, paragraphs, hyperlinks, italic and
bold tags, horizontal rules, blockquote, ordered and unordered lists, code. Unknown tags are stripped. This function
doesn't return any error as it converts the input text into something as good as possible, without validation.

```go
RichTextToHtml(rt interface{}, linkResolver LinkResolverFunc, entryLinkResolver EntryLinkResolverFunc, imageResolver ImageResolverFunc, locale Locale) (string, error)
```

Converts an interface representing a Contentful RichText value (usually from a field getter) into HTML.
The function takes in three (optional) functions as parameters to resolve
hyperlink URLs, permalinks to entries and to derive IMG tag attributes for embedded image assets. The three functions
return a map of attributes for the HTML tag the RichTextToHtml function will emit (either an A or an IMG) and have the
following signature. Note that the ImageResolverFunc function must return a customHTML value that can be empty but if
set it will substitute the IMG tag with the returned HTML snippet. This allows you to emit custom mark-up for your
images, e.g. a PICTURE tag.

```go
type LinkResolverFunc func(url string) (resolvedAttrs map[string]string, resolveError error)

type EntryLinkResolverFunc func(entryID string, locale Locale) (resolvedAttrs map[string]string, resolveError error)

type ImageResolverFunc func(assetID string, locale Locale) (attrs map[string]string, customHTML string, resolveError error)

type EmbeddedEntryResolverFunc func(entryID string, locale Locale) (html string, resolveError error)
```

All the three functions above can be passed as nil with different levels of graceful degrading.

```go
RichTextToPlainText(rt interface{}, locale Locale) (string, error)

Converts the RichText to plain text.
```

### Constants and global variables

Each generated content type library file exports a constant with the Contentful ID of the content type itself, for
example in _contentful_vo_lib_person.go_:

```go
const ContentTypePerson = "person"
```

Constants are available for each locale supported by the space at the time of code generation, e.g.:

```go
const SpaceLocaleGerman Locale = "de"
const SpaceLocaleFrench Locale = "fr"
const defaultLocale Locale = SpaceLocaleGerman
```

Four levels of logging are supported (even if only partially used at this time):

```go
const (
    LogDebug = 0
    LogInfo  = 1
    LogWarn  = 2
    LogError = 3
)
```

A global variable named _SpaceContentTypeInfoMap_ contains an ID-indexed map of all content types
with their names and descriptions
