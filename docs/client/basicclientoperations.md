---
sidebar_label: Basic client operations
sidebar_position: 1
---

# Basic client operations

Let's consider a very simple use case. You have a Contentful space where you store information
about people and their pets.

To generate a go package to manipulate those entries, run the following in your terminal (refer to the [Setup section](setup) for details):

```shell
$ gocontentful -spaceid YOUR_SPACE_ID -cmakey YOUR_CMA_API_TOKEN -contenttypes person,pet path/to/your/go/project/folder/people
```

The **-contenttypes** parameter is optional. If not specified, gocontentful will generate an API that supports all the content types of the space.

Gocontentful will scan the space, download locales and content types and generate the Go API files in the target path:

```shell
path/to/your/go/project/folder/people
|-gocontentfulvobase.go
|-gocontentfulvolib_person.go // One file for each content type
|-gocontentfulvolib_pet.go    // One file for each content type
|-gocontentfulvolib.go
|-gocontentfulvo.go
```

[We recommend](setup) not passing the _-cmakey_ parameter but rather log in first using the Contentful CLI.
This will be remembered in all subsequent runs. See the [Setup section]

Note: Never modify the generated files. If you change the content model in Contentful, run gocontentful
again. This will update the files for you.

### Get a client

The generated files will be in the "people" subdirectory of your project. Your go program can get a Contentful
client from them:

```go
cc, err := people.NewContentfulClient(YOUR_SPACE_ID, people.ClientModeCDA, YOUR_API_KEY, 1000, contentfulLogger, people.LogDebug,false)
```

The parameters to pass to NewContentfulClient are:

- _spaceID_ (string)
- _clientMode_ (string) supports the constants ClientModeCDA, ClientModeCPA and ClientModeCMA. If you need to operate
  on multiple APIs (e.g. one for reading and CMA for writing) you need to get two clients
- _clientKey_ (string) is your API key (generate one for your API at Contentful)
- _optimisticPageSize_ (uint16) is the page size the client will use to download entries from the space for caching.
  Contentful's default is 100 but you can specify up to 1000: this might get you into an error because Contentful
  limits the payload response size to 70 KB but the client will handle the error and reduce the page size automatically
  until it finds a proper value. Hint: using a big page size that always fails is a waste of time and resources because
  a lot of initial calls will fail, whereas a too small one will not leverage the full download bandwidth. It's a
  trial-and-error and you need to find the best value for your case. For simple content types you can start with 1000,
  for very complex ones that include fat fields you might want to get down to 100 or even less.
- _logFn_ is a func(fields map[string]interface{}, level int, args ...interface{}) that the client will call whenever
  it needs to log something. It can be nil if you don't need logging and that will be handled gracefully but it's not
  recommended. A simple function you can pass that uses the https://github.com/Sirupsen/logrus package might look
  something like this:

```go
contentfulLogger := func(fields map[string]interface{}, level int, args ...interface{}) {
    switch level {
    case people.LogDebug:
        log.WithFields(fields).Debug(args)
    case people.LogInfo:
        log.WithFields(fields).Info(args)
    case people.LogWarn:
        log.WithFields(fields).Warn(args)
    case people.LogError:
        log.WithFields(fields).Error(args)
    default:
        return
    }
}
```

- _logLevel_ (int) is the debug level (see function above). Please note that LogDebug is very verbose and even logs
  when you request a field value but that is not set for the entry.
- _debug_ (bool) is the Contentful API client debug switch. If set to _true_ it will log on stdout all the CURL calls
  to Contentful. This is extremely verbose and extremely valuable when something fails in a call to the API because
  it's the only way to see the REST API response.

_NOTE:_ Gocontentful provides an offline version of the client that can load data from a JSON space export file
(as exported by the _contentful_ CLI tool). This is the way you can write unit tests against your generated API that
don't require to be online and the management of a safe API key storage. See the [API Reference](./api-reference)

### Environments support

Gocontentful supports Contentful environments in two ways:

- Code can be generated loading the content model from an environment other than master.
  This is done passing the -environment flag on the command line to specify the environment you want to load the model from.
- The gocontentful client in your application can be switched to any environment with the SetEnvironment method.
  For example, if your space has an extra environment named "devplayground" you can switch the API to use it with:

```go
cc.SetEnvironment("devplayground")
```

To reset the environment to master pass an empty string.

### Working with RichText

Contentful supports Rich Text fields. Behind the scenes, these are JSON objects that represent
the content through a Contentful-specific data model. Sooner or later you might want to convert such values to and from HTML.
Gocontentful supports the conversion both ways. For instance, you want a person's resume to be converted to HTML:

```go
htmlText := people.RichTextToHtml(person.Resume(), linkResolver, entryLinkResolver, imageResolver, embeddedEntryResolver locale)
```

The parameters linkResolver, entryLinkResolver, embeddedEntryResolver and imageResolver are all functions that you can pass
to convert various elements inserted by the user into the RichText field:

- linkResolver will allow you to create custom HTML tags for hyperlinks. If left blank, RichTextToHtml will just output an A tag.
- entryLinkResolver is used to create hyperlinks with custom URLs when the destination in Contentful is another entry.
  If you allow such links to be created in the editor then you must pass this function.
- imageResolver and embebbedEntryResolver are needed when the field accepts assets and entries embedded into the content, to turn
  these into actual HTML snippets

The conversion works the other way around too, when you need to source data from outside and create Contentful entries:

```go
myRichText := HtmlToRichText(htmlSrc)
```

See the [API Reference](./api-reference) for more details about these functions.

### More on references

When working with references it's often useful to know if there are any broken ones in the space.
This happens when a published entry references another that has been deleted after the parent
was published. This might create issues if your application code doesn't degrade content gracefully.
To get a report of all broken references you can use the following function:

```go
(cc *ContentfulClient) BrokenReferences() (brokenReferences []BrokenReference)
```

Note that this only works with cached clients. See [the next chapter on caching](./caching).

Also on references: when you want to reference entry B from entry A, you cannot assign
the value object of entry B to the reference field in A. First you need to convert the
object to a `ContentTypeSys` object because that's what Contentful expects in reference fields:

```go
(vo *CfPerson) ToReference() (refSys ContentTypeSys)
```

Finally, you can get the parents (AKA referring) entries of either an entry or
an EntryReference with the _GetParents()_ method. This returns a slice of `[]EntryReference`:

```go
(vo *CfPerson) GetParents() (parents []EntryReference, err error)
(ref *EntryReference) GetParents(cc *ContentfulClient) (parents []EntryReference, err error)
```

### Other useful functions

Another thing you might want to know is the content type of an entry with a given ID:

```go
(cc *ContentfulClient) GetContentTypeOfID(ID string) (contentType string)
```

### Caveats and limitations

- Avoid creating content types that have field IDs equal to reserved Go words (e.g. "type").
  Gocontentful won't scan for them and the generated code will break.
