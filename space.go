package erm

import (
	"bytes"
	"encoding/json"
	"log"
	"runtime"
	"strings"
	"text/template"

	"github.com/foomo/contentful"
)

// SpaceConf is the space config object passed to the template
type SpaceConf struct {
	Filename     string
	FuncMap      map[string]interface{}
	PackageName  string
	Locales      []Locale
	ContentTypes []ContentType
}

// GetLocales retrieves locale definition from Contentful
func GetLocales(CMA *contentful.Contentful, spaceID *string) (locales []Locale, err error) {

	col, err := CMA.Locales.List(*spaceID).GetAll()
	if err != nil {
		log.Fatal("Couldn't get locales")
	}
	for _, item := range col.Items {
		var locale Locale
		byteArray, _ := json.Marshal(item)
		err = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&locale)
		if err != nil {
			break
		}
		locales = append(locales, locale)
	}
	return
}

// GetContentTypes retrieves content type definition from Contentful
func GetContentTypes(CMA *contentful.Contentful, spaceID *string) (contentTypes []ContentType, err error) {

	col := CMA.ContentTypes.List(*spaceID)
	//col.Query.Equal("name", "Flyout")
	_, err = col.GetAll()
	if err != nil {
		log.Fatal("Couldn't get locales")
	}
	for _, item := range col.Items {
		var contentType ContentType
		byteArray, _ := json.Marshal(item)
		err = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&contentType)
		if err != nil {
			break
		}
		for _, field := range contentType.Fields {
			if (field.Type == FieldTypeArray) && field.Items != nil && field.Items.Validations != nil {

			}
		}
		contentTypes = append(contentTypes, contentType)
	}
	return
}

// ProcessSpace calls the generators
func ProcessSpace(packageName string, locales []Locale, contentTypes []ContentType) (err error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	funcMap := template.FuncMap{"fieldIsAsset": fieldIsAsset, "fieldIsReference": fieldIsReference, "firstCap": strings.Title, "mapFieldType": mapFieldType}
	conf := SpaceConf{Filename: filename, FuncMap: funcMap, PackageName: packageName, Locales: locales, ContentTypes: contentTypes}

	err = GenerateVo(conf)
	if err != nil {
		panic(err)
	}

	err = GenerateLib(conf)
	if err != nil {
		panic(err)
	}

	return
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
	case FieldTypeObject: // JSON field
		return "Cf" + firstCap(contentTypeName) + firstCap(field.ID)
	case FieldTypeRichText:
		return "interface{}"
	case FieldTypeSymbol: // It's a text field
		return "string"
	default:
		return ""
	}
}

func fieldIsReference(field ContentTypeField) bool {
	if (field.Type == FieldTypeArray && field.Items.Type == FieldItemsTypeLink && field.Items.LinkType == FieldLinkTypeEntry) || (field.Type == FieldTypeLink && field.LinkType == FieldLinkTypeEntry) {
		return true
	}
	return false
}

func fieldIsAsset(field ContentTypeField) bool {
	if (field.Type == FieldTypeArray && field.Items.Type == FieldItemsTypeLink && field.Items.LinkType == FieldLinkTypeAsset) || (field.Type == FieldTypeLink && field.LinkType == FieldLinkTypeAsset) {
		return true
	}
	return false
}
