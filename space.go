package erm

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
	"text/template"

	"github.com/foomo/contentful"
)

type Conf struct { 
    Locales []Locale
    ContentTypes []ContentType
} 

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

func ProcessSpace(locales []Locale, contentTypes []ContentType) (err error) {
	conf := Conf{Locales: locales, ContentTypes: contentTypes}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	tmpl, err := template.New(VoLib).ParseFiles(path.Dir(filename) + TplDir + VoLib)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, conf)
	if err != nil {
		panic(err)
	}
	return
}
