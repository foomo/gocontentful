package erm

import (
	"strings"
	"text/template"
)

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"firstCap":                 strings.Title,
		"fieldIsBasic":             fieldIsBasic,
		"fieldIsComplex":           fieldIsComplex,
		"fieldIsAsset":             fieldIsAsset,
		"fieldIsBoolean":           fieldIsBoolean,
		"fieldIsDate":              fieldIsDate,
		"fieldIsInteger":           fieldIsInteger,
		"fieldIsJSON":              fieldIsJSON,
		"fieldIsLink":              fieldIsLink,
		"fieldIsLocation":          fieldIsLocation,
		"fieldIsMultipleReference": fieldIsMultipleReference,
		"fieldIsMultipleAsset":     fieldIsMultipleAsset,
		"fieldIsNumber":            fieldIsNumber,
		"fieldIsReference":         fieldIsReference,
		"fieldIsRichText":          fieldIsRichText,
		"fieldIsSymbol":            fieldIsSymbol,
		"fieldIsSymbolList":        fieldIsSymbolList,
		"fieldIsText":              fieldIsText,
		"mapFieldType":             mapFieldType,
		"mapFieldTypeLiteral":      mapFieldTypeLiteral,
		"onlyLetters":              onlyLetters,
	}
}

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
		return "interface{}"
	case FieldTypeRichText:
		return "interface{}"
	case FieldTypeSymbol: // It's a text field
		return "string"
	case FieldTypeText: // It's a text field
		return "string"
	default:
		return ""
	}
}

// mapFieldTypeLiteral takes a ContentTypeField from the space model definition
// and returns an empty literal that matches the type of the map[string] for the VO
func mapFieldTypeLiteral(contentTypeName string, field ContentTypeField) string {
	switch field.Type {
	case FieldTypeBoolean:
		return "false"
	case FieldTypeDate, FieldTypeSymbol, FieldTypeText:
		return `""`
	case FieldTypeInteger, FieldTypeNumber:
		return "0"
	case FieldTypeArray, FieldTypeLink, FieldTypeLocation, FieldTypeJSON, FieldTypeRichText:
		return "nil"
	default:
		return ""
	}
}

func fieldIsAsset(field ContentTypeField) bool {
	return field.Type == FieldTypeLink && field.LinkType == FieldLinkTypeAsset
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

func fieldIsMultipleAsset(field ContentTypeField) bool {
	return field.Type == FieldTypeArray && field.Items.Type == FieldItemsTypeLink && field.Items.LinkType == FieldLinkTypeAsset
}
func fieldIsMultipleReference(field ContentTypeField) bool {
	return field.Type == FieldTypeArray && field.Items.Type == FieldItemsTypeLink && field.Items.LinkType == FieldLinkTypeEntry
}

func fieldIsNumber(field ContentTypeField) bool {
	return field.Type == FieldTypeNumber
}

func fieldIsReference(field ContentTypeField) bool {
	return field.Type == FieldTypeLink && field.LinkType == FieldLinkTypeEntry
}

func fieldIsRichText(field ContentTypeField) bool {
	return field.Type == FieldTypeRichText
}

func fieldIsSymbol(field ContentTypeField) bool {
	return field.Type == FieldTypeSymbol
}

func fieldIsSymbolList(field ContentTypeField) bool {
	return field.Type == FieldTypeArray && field.Items.Type == FieldItemsTypeSymbol
}

func fieldIsText(field ContentTypeField) bool {
	return field.Type == FieldTypeText
}

func fieldIsBasic(field ContentTypeField) bool {
	return fieldIsSymbolList(field) || fieldIsBoolean(field) || fieldIsInteger(field) || fieldIsNumber(field) || fieldIsSymbol(field) || fieldIsText(field) || fieldIsDate(field)
}
func fieldIsComplex(field ContentTypeField) bool {
	return field.Type == FieldTypeJSON || field.Type == FieldTypeLocation || field.Type == FieldTypeRichText
}
