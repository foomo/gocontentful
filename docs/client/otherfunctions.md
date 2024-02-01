---
sidebar_label: Other functions
sidebar_position: 6
---

# Other useful functions

Another thing you might want to know is the content type of an entry with a given ID:

```go
(cc *ContentfulClient) GetContentTypeOfID(ID string) (contentType string)
```

## Caveats and limitations

- Avoid creating content types that have field IDs equal to reserved Go words (e.g. "type").
  Gocontentful won't scan for them and the generated code will break.
