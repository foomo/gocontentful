package erm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	Version      string
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
	sort.Slice(
		contentTypes, func(i, j int) bool {
			return contentTypes[i].Name < contentTypes[j].Name
		},
	)
	return
}

func getData(spaceID, cmaKey, environment, exportFile string, flagContentTypes []string) (
	finalContentTypes []ContentType, locales []Locale, err error,
) {
	var contentTypes []ContentType
	if exportFile != "" {
		fileBytes, err := ioutil.ReadFile(exportFile)
		if err != nil {
			return nil, nil, fmt.Errorf("error reading export file: %v", err)
		}
		var export ExportFile
		err = json.Unmarshal(fileBytes, &export)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing export file: %v", err)
		}
		contentTypes = export.ContentTypes
		locales = export.Locales
	} else {
		// Get client
		CMA := contentful.NewCMA(cmaKey)
		CMA.Debug = false
		if environment != "" {
			CMA.Environment = environment
		}

		// Get space locales
		locales, err = getLocales(CMA, spaceID)
		if err != nil {
			return nil, nil, fmt.Errorf("could not get locales: %v", err)
		}
		fmt.Println("Locales found:", locales)

		// Get content types
		contentTypes, err = getContentTypes(CMA, spaceID)
		if err != nil {
			return nil, nil, fmt.Errorf("could not get content types: %v", err)
		}
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
func GenerateAPI(dir, packageName, spaceID, cmaKey, environment, exportFile string, flagContentTypes []string, version string) (err error) {
	contentTypes, locales, errGetData := getData(spaceID, cmaKey, environment, exportFile, flagContentTypes)
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
		Version:      version,
	}
	return generateCode(conf)
}
