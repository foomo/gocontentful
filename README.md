# Contentful ERM

A Contentful Entry-Reference Mapping code generator for Go. Initial features:

- Creation of Value Objects from Contentful content model
- CDA, CPA, CMA support for CRUD operations
- Automatic management/resolution of references

How to run and generate the library files
-----------------------------------------

Let's assume you want to generate a package named "people" and manipulate entries of content type ID "person". From the main directory run this:

`go run cmd/contentfulerm.go -spaceid=YOUR_SPACE_ID -cmakey=YOUR_MANAGEMENT_KEY -package=people -contenttypes=person`

Note: You can pass multiple values to contenttypes as a comma-separated list

The script will scan the space, download locales and content types and generate some files in a directory named after the package inside the "generated" directory:

<pre><code>generated/people
|-contentful_vo_base.go
|-contentful_vo_lib_person.go // One file for each content type
|-contentful_vo_lib.go
|-contentful_vo.go
</code></pre>

Copy these files into a subdirectory of your project and import the "people" package. 

_Note: Do NOT modify these files! If you change the content model in Contentful you will need to run the generator again and overwrite the files._

Public function set
---------------------

**BASE FUNCTIONS COMMON TO ALL CONTENT TYPES**

>**NewContentfulClient**(spaceID string, ck *ContentfulKeys, logFn func(contentType string, entryID string, method string, err error), debug bool) (*ContentfulClient, error)

Creates a Contentful client, this is the first function you need to call:

* _spaceID_ is the ID of the Contentful space the client attach to
* _ck_ is a struct you can fill in with all the keys for the APIs you need to work with. For instance, you will need a CDAKey if you want to use the Get methods for the Content Delivery API or a CMAKey if you nee the Manage methods. 
* _logFn_ is logging function you can pass to the client to be used by the shortcut value getter or conversion methods. Normal getters are named after the field name or content type and return both a value and an error, e.g. _Name()_ for a person. Often in the application code it's safe and much less verbose not to handle errors for each getter but accept a default  return value, for example an empty string for a string field that is not filled in. In these cases you can prefix the method with "ValueOf", e.g. _ValueOfName()_ to get only the value. If you pass a logging function to the client's constructor, it will be called transparently any time it's not possible to read or convert a value and the default is returned. This way the application code remains lean but you still get full logging of the underlying operations. Note that it's your responsibility to check in your application that some return values are safe before calling further methods, especially in case of _nil_ values like a en empty reference field that might panic if you call a conversion function on them!
* _debug_ enables or disables the debug mode in the Contentful client

---

**FUNCTIONS NAMED AFTER THE CONTENT TYPE**

_For these we're assuming a content type named "Person"._

>**NewCfPerson**() (cfPerson *CfPerson)

Creates a new Person entry. You can manipulate and upsert this later

>(cc *ContentfulClient) **GetAllPerson**() (vos []*CfPerson, err error)

Retrieves all Person entries from the Contentful Delivery API (CDA). You need to have the CDA Key setup in NewContentfulClient for this to work.

>(cc *ContentfulClient) **GetFilteredPerson**(query *contentful.Query) (vos []*CfPerson, err error) 

Retrieves Person entries matching the specified query from the CDA.

>(cc *ContentfulClient) **GetPersonByID**(id string) (vo *CfPerson, err error)

Retrieves the Person entry with the specified ID from the CDA.

>(cc *ContentfulClient) **ManageAllPerson**() (vos []*CfPerson, err error)

Retrieves all the draft versions of the Person entries from the Contentful Management API (CMA). You need to have the CMA Key setup in NewContentfulClient for this to work. Note that you will need this or one or the other Manage* functions if you want to update one or multiple existing entries in Contentful. The Preview API doesn't return the entry version and upserting the modified entry is not possible.

>(cc *ContentfulClient) **ManageFilteredPerson**(query *contentful.Query) (vos []*CfPerson, err error) 

Retrieves draft Person entries matching the specified query from the CMA.

>(cc *ContentfulClient) **ManagePersonByID**(id string) (vo *CfPerson, err error)

Retrieves the draft Person entry with the specified ID from the CMA.

>(cc *ContentfulClient) **PreviewAllPerson**() (vos []*CfPerson, err error)

Retrieves all the draft versions of the Person entries from the Contentful Preview API (CPA). You need to have the CPA Key setup in NewContentfulClient for this to work.

>(cc *ContentfulClient) **PreviewFilteredPerson**(query *contentful.Query) (vos []*CfPerson, err error) 

Retrieves draft Person entries matching the specified query from the CPA.

>(cc *ContentfulClient) **PreviewPersonByID**(id string) (vo *CfPerson, err error)

Retrieves the draft Person entry with the specified ID from the CPA.

>(ref ContentfulReferencedEntry) **ToCfPerson**() (vo *CfPerson, err error)

Converts a referenced entry to the specified value object. See the ContentType() function above.

>(ref ContentfulReferencedEntry) **ValueOfCfPerson**() (vo *CfPerson)

Shortcut (value only) version of the previous method, with automatic logging to the client if defined (see NewContentfulClient above)

---

**REFERENCE CONVERSION FUNCTIONS**

>(vo *CfPerson) **ToReference**() (refSys ContentTypeSys) 

Converts a value object into a reference that can be added to a reference field of an entry.

>(vo *CfPerson) **ToReferenceArray**() (refSysArray []ContentTypeSys) 

Converts a value object into a reference array that can be added to a multiple reference field of an entry.

>(ref ContentfulReferencedEntry) **ContentType**() (contentType string)

Returns the Contentful content type of a referenced entry. This is used for reference fields that validate multiple content types: you will want to switch/case this and manage the referenced type with the right value object.

---

**FIELD GETTERS**

Field getters are named after the field ID in Contentful and return the proper type. For example, if the Person content type has a Symbol (short text) field named 'Name', this will be the getter:

>(vo *CfPerson) **Name**(locale ...string) (string, error) 

The locale parameter is optional and if not passed, the function will return the value for the default locale of the space. If the locale is specified and it's not available for the space, an error is returned. If the locale is valid but a value doesn't exist for the field and locale, the function will return the value for the default locale if that's specified as a fallback locale in the space definition in Contentful, otherwise will return an error.

Possible return types are:

- _string_ for fields of types Symbol, Text, Date
- _[]string_ for fields of type List
- _float64_ for fields of type Integer or Number
- _bool_ for fields of type Boolean
- _*ContentTypeSys_ for single reference fields
- _[]*ContentTypeSys_ for multiple reference fields
- _*ContentTypeFieldLocation_ for fields of type Location
- *interface{} for fields of type Object or RichText

>(vo *CfPerson) **ValueOfName**(locale ...string) (string) 

Shortcut (value only) version of the previous method, with automatic logging to the client if defined (see NewContentfulClient above)

---

**FIELD SETTERS**

Field setters are named after the field ID in Contentful and require to pass in the proper type. See FIELD GETTERS above for a reference. Example:

>(vo *CfPerson) **SetName**(title string, locale ...string) (err error) 

---

**WRITE OPERATIONS**

>(vo *CfPerson) **UpsertEntry**(cc *ContentfulClient) (err error) 

Upserts the entry. This will appear as "Draft" (if it's a new entry) or "Changed" if it's already existing. In the latter case, you will need to retrieve the entry with one of the Manage* functions above to acquire the Sys object that contains the version information. Otherwise the API call will fail with a "Version mismatch" error.

>(vo *CfPerson) **PublishEntry**(cc *ContentfulClient) (err error) 

Publishes the entry. Note that before publshing you will need to retrieve the entry with one of the Manage* functions above to acquire the Sys object that contains the version information. Otherwise the API call will fail with a "Version mismatch" error. This is needed even if you have just upserted the entry with the function above!

>(vo *CfPerson) **UpdateEntry**(cc *ContentfulClient) (err error) 

Shortcut function that upserts and publishes the entry. Note that before calling this you will need to retrieve the entry with one of the Manage* functions above to acquire the Sys object that contains the version information. Otherwise the API call will fail with a "Version mismatch" error. Using this shortcut function avoids retrieving the entry twice.

>(vo *CfPerson) **DeleteEntry**(cc *ContentfulClient) (err error) 

Unpublishes and deletes the entry

---

**UTILITY FUNCTIONS**

>**NewAssetFromURL**(id string, uploadUrl string, imageFileType string, title string, locale ...string) *contentful.Asset

Creates an Asset from an URL of an existing file online (you still need to upsert it later).

>**ToAssetReference**(asset *contentful.Asset) (refSys ContentTypeSys) 

Converts the asset to a reference. You need to do this before you add the asset to a reference field of an entry.

>**HtmlToRichText**(htmlSrc string) *RichTextNode

Converts an HTML fragment to a RichTextNode. This is far from complete but useful to migrate data from third-party systems to Contentful. It currently supports headings, paragraphs, hyperlinks, italic and bold tags, horizontal rules, blockquote, ordered and unordered lists, code. Unknown tags are stripped. This function doesn't return any error as it converts the input text into something as good as possible, without  validation.

>**RichTextToHtml**(rt interface{}) string

Converts an interface representing a Contentful RichText value (usually from a field getter) into HTML. It currently supports all tags except for embedded and inline entries and assets.