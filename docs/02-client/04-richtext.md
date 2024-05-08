# RichText

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
