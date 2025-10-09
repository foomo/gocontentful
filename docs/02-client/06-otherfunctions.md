# Other useful functions

Another thing you might want to know is the content type of an entry with a given ID:

```go
(cc *ContentfulClient) GetContentTypeOfID(ctx, ID string) (contentType string)
```

## Field inheritance

When working with hierarchical content structures, you may need to retrieve field values that can be inherited from parent entries. The `GetOrInheritFieldValue` function provides this functionality:

```go
func GetOrInheritFieldValue(ctx context.Context, contentfulShop *ContentfulClient, entryID string, field string, parentContentTypes []string, locale Locale) (any, error)
```

This function retrieves a field value from a GenericEntry, traversing up the full parent hierarchy if the field is not found in the current entry. It's particularly useful for content types that have parent-child relationships where certain fields should be inherited from parent entries.

### Parameters

- `ctx`: Context for the operation
- `contentfulShop`: The ContentfulClient instance
- `entryID`: The ID of the entry to start the search from
- `field`: The field name to retrieve
- `parentContentTypes`: List of content type IDs that should be considered as valid parents
- `locale`: The locale to retrieve the field value in

### Return value

Returns the field value as `any` type, or an error if:
- The entry cannot be retrieved
- A circular reference is detected in the parent hierarchy
- The field is not found in any parent entry
- Any parent entry cannot be retrieved

### Example usage

```go
ctx := context.Background()
parentContentTypes := []string{"category", "section"}

// Try to get the "theme" field from the current entry or inherit it from parents
theme, err := GetOrInheritFieldValue(ctx, cc, "my-entry-id", "theme", parentContentTypes, contentful.DefaultLocale)
if err != nil {
    log.Printf("Theme not found: %v", err)
} else {
    log.Printf("Theme value: %v", theme)
}
```

### Important notes

- The function first attempts to get the field value from the current entry
- If not found, it recursively traverses up the parent hierarchy
- Circular references are detected and will return an error to prevent infinite loops
- The function respects the locale parameter for localized field values
- Only entries with content types specified in `parentContentTypes` are considered as valid parents

## Caveats and limitations

- Avoid creating content types that have field IDs equal to reserved Go words (e.g. "type").
  Gocontentful won't scan for them and the generated code will break.
