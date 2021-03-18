package erm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/foomo/contentful"
)

// spaceConf is the space config object passed to the template
type spaceConf struct {
	FuncMap      map[string]interface{}
	PackageName  string
	PackageDir   string
	Locales      []Locale
	ContentTypes []ContentType
	ContentType  ContentType
}

// GetLocales retrieves locale definition from Contentful
func getLocales(CMA *contentful.Contentful, spaceID string) (locales []Locale, err error) {

	col, err := CMA.Locales.List(spaceID).GetAll()
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
func getContentTypes(CMA *contentful.Contentful, spaceID string) (contentTypes []ContentType, err error) {

	col := CMA.ContentTypes.List(spaceID)
	_, err = col.GetAll()
	if err != nil {
		log.Fatal("Couldn't get content types")
	}
	for _, item := range col.Items {
		var contentType ContentType
		byteArray, _ := json.Marshal(item)
		err = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&contentType)
		if err != nil {
			break
		}
		contentTypes = append(contentTypes, contentType)
	}
	return
}

func getData(spaceID, cmaKey string, flagContentTypes []string) (finalContentTypes []ContentType, locales []Locale, err error) {
	// Get client
	CMA := contentful.NewCMA(cmaKey)
	CMA.Debug = false

	// Get space locales
	locales, errGetLocales := getLocales(CMA, spaceID)
	if errGetLocales != nil {
		return nil, nil, errors.New("Could not get locales: " + errGetLocales.Error())
	}
	fmt.Println("Locales found:", locales)

	// Get content types
	contentTypes, err := getContentTypes(CMA, spaceID)
	if err != nil {
		return nil, nil, errors.New("Could not get content types")
	}
	fmt.Println("Content types found:", len(contentTypes))

	finalContentTypes = []ContentType{}
	if len(flagContentTypes) == 0 {
		finalContentTypes = contentTypes
	} else {
		for _, ct := range contentTypes {
			if sliceIncludes(flagContentTypes, ct.Sys.ID) {
				finalContentTypes = append(finalContentTypes, ct)
			}
		}
	}
	var finalContentTypesString []string
	for _, finalContentType := range finalContentTypes {
		finalContentTypesString = append(finalContentTypesString, finalContentType.Name)
	}
	fmt.Println("Filtered Content types:", len(finalContentTypes), strings.Join(finalContentTypesString, ", "))
	return finalContentTypes, locales, nil
}

// GenerateAPI calls the generators
func GenerateAPI(dir, packageName, spaceID, cmaKey string, flagContentTypes []string) (err error) {
	contentTypes, locales, errGetData := getData(spaceID, cmaKey, flagContentTypes)
	if errGetData != nil {
		return errGetData
	}

	packageDir := filepath.Join(dir, packageName)
	errMkdir := os.MkdirAll(packageDir, 0766)
	if errMkdir != nil {
		return errMkdir
	}
	funcMap := getFuncMap()
	conf := spaceConf{
		PackageDir:   packageDir,
		FuncMap:      funcMap,
		PackageName:  packageName,
		Locales:      locales,
		ContentTypes: contentTypes,
	}
	return generateCode(conf)
}
