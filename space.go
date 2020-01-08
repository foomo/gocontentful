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
	PackageName  string
	Locales      []Locale
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

func ProcessSpace(packageName string, locales []Locale, contentTypes []ContentType) (err error) {
	conf := Conf{PackageName: packageName, Locales: locales, ContentTypes: contentTypes}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	// VO base generation
	tmpl, err := template.New(VoBase + TplExt).ParseFiles(path.Dir(filename) + TplDir + VoBase + TplExt)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path.Dir(filename) + OutDir + VoBase + GoExt)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, conf)
	if err != nil {
		panic(err)
	}

	// Lib generation
	tmpl, err = template.New(VoLib + TplExt).ParseFiles(path.Dir(filename) + TplDir + VoLib + TplExt)
	if err != nil {
		panic(err)
	}

	f, err = os.Create(path.Dir(filename) + OutDir + VoLib + GoExt)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, conf)
	if err != nil {
		panic(err)
	}

	return
}
