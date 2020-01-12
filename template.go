package erm

// mapFieldType takes a ContentTypeField from the space model definition
// and returns a string that matches the type of the map[string] for the VO
func mapFieldType(contentTypeName string, field ContentTypeField) string {
	switch field.Type {
	case FieldTypeArray: // It's either a text list or a multiple reference
		switch field.Items.Type {
		case FieldItemsTypeSymbol:
			return "[]string"
		case FieldItemsTypeLink:
			return "[]ContentTypeSys"
		default:
			return ""
		}
	case FieldTypeBoolean:
		return "bool"
	case FieldTypeDate:
		return "string"
	case FieldTypeInteger:
		return "float64"
	case FieldTypeLink: // A single reference
		return "ContentTypeSys"
	case FieldTypeLocation:
		return "ContentTypeFieldLocation"
	case FieldTypeNumber: // Floating point
		return "float64"
	case FieldTypeJSON: // JSON field
		return "Cf" + firstCap(contentTypeName) + firstCap(field.ID)
	case FieldTypeRichText:
		return "interface{}"
	case FieldTypeText: // It's a text field
		return "string"
	default:
		return ""
	}
}

func fieldIsAsset(field ContentTypeField) bool {
	return (field.Type == FieldTypeArray && field.Items.Type == FieldItemsTypeLink && field.Items.LinkType == FieldLinkTypeAsset) || (field.Type == FieldTypeLink && field.LinkType == FieldLinkTypeAsset)
}

func fieldIsBoolean(field ContentTypeField) bool {
	return field.Type == FieldTypeBoolean
}

func fieldIsDate(field ContentTypeField) bool {
	return field.Type == FieldTypeDate
}

func fieldIsInteger(field ContentTypeField) bool {
	return field.Type == FieldTypeInteger
}

func fieldIsJSON(field ContentTypeField) bool {
	return field.Type == FieldTypeJSON
}

func fieldIsLink(field ContentTypeField) bool {
	return field.Type == FieldTypeLink
}

func fieldIsLocation(field ContentTypeField) bool {
	return field.Type == FieldTypeLocation
}

func fieldIsNumber(field ContentTypeField) bool {
	return field.Type == FieldTypeNumber
}

func fieldIsReference(field ContentTypeField) bool {
	return (field.Type == FieldTypeArray && field.Items.Type == FieldItemsTypeLink && field.Items.LinkType == FieldLinkTypeEntry) || (field.Type == FieldTypeLink && field.LinkType == FieldLinkTypeEntry)
}

func fieldIsRichText(field ContentTypeField) bool {
	return field.Type == FieldTypeRichText
}

func fieldIsText(field ContentTypeField) bool {
	return field.Type == FieldTypeText
}

func fieldIsTextList(field ContentTypeField) bool {
	return field.Type == FieldTypeArray && field.Items.Type == FieldItemsTypeSymbol
}




