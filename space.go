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
	funcMap := template.FuncMap{
		"firstCap": strings.Title, 
		"fieldIsAsset": fieldIsAsset, 
		"fieldIsBoolean": fieldIsBoolean, 
		"fieldIsDate": fieldIsDate, 
		"fieldIsInteger": fieldIsInteger, 
		"fieldIsJSON": fieldIsJSON, 
		"fieldIsLink": fieldIsLink, 
		"fieldIsLocation": fieldIsLocation, 
		"fieldIsNumber": fieldIsNumber, 
		"fieldIsReference": fieldIsReference, 
		"fieldIsRichText": fieldIsRichText, 
		"fieldIsText": fieldIsText, 
		"fieldIsTextList": fieldIsTextList, 
		"mapFieldType": mapFieldType,
	}
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
