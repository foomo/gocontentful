---
sidebar_label: Entries
sidebar_position: 2
---

# Working with entries

Refer to the [Getting started section](../gettingstarted) for an introduction on entry operations.
With your newly created client you can do things like:

```go
// Load all persons
persons, err := cc.GetAllPerson()
// Load a specific person
person, err := cc.GetPersonByID(THE_PERSON_ID)
// or pass a query
person, err := GetFilteredPerson(&contentful.Query{
	"contentType":"person",
    "exists": []string{"fields.resume"}
})
// The person's name
name := person.Name()
// The work title in a different localization. Available locales are generated as constants.
// If a space is configured to have a fallback from one locale to the default one,
// the getter functions will return that if the value is not set for locale passed to the function.
name := person.Title(people.SpaceLocaleItalian)
// Get references to the person's pets
petRefs := person.Pets()
// Deal with pets
for _, pet := range petRefs {
switch pet.ContentType {
case people.ContentTypeDog: // you have these constants in the generated code
dog := pet.VO.(*people.Dog)
// do something with dog
case people.ContentTypeCat:
// ...
}
```

Once you have loaded an entry, you can use any of the setter methods to alter the fields. For example:

```go
dog.SetAge(7)
```

This will only affect the Go object and doesn't automatically propagate to the space.
To save the entry to Contentful you need to explicitly call one of these methods:

```go
// Upsert (save) an entry
err := dog.UpsertEntry()
// Publish it (after it's been upserted)
err := dog.PublishEntry() // change your mind with err := dog.UnpublishEntry()
// Or do it in one step
err := dog.UpdateEntry() // upserts and publishes
// And delete it
err := dog.DeleteEntry()
```

If you want to know the publication status of an entry as represented in Contentful's UI you
can use the `GetPublishingStatus()` method on the entry itself. Possible return values are the
predefined constants `StatusDraft`, `StatusChanged` and `StatusPublished`.

When saving, publishing or deleting entries:

- You need a client that uses mode _ClientModeCMA_. Entries retrieved with ClientModeCDA
  or ClientModeCPA can be saved in memory (for example if you need to enrich the built-in cache) but not persisted to
  Contentful.
- Make sure you Get a fresh copy of the entry right before you manipulate it and upsert it / publish it to Contentful. In case it's
  saved by someone else in the meantime, the upsert will fail with a version mismatch error.

In case you need a completely new entry just create it, Contentful will fill in the technical details for you:

```go
NewCfPerson(contentfulClient ...*ContentfulClient) (cfPerson *CfPerson)
```
