# More on references

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
(vo *CfPerson) GetParents(ctx) (parents []EntryReference, err error)
(ref *EntryReference) GetParents(cc *ContentfulClient) (parents []EntryReference, err error)
```
