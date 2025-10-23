package erm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/foomo/contentful"
	"github.com/pkg/errors"
)

// spaceConf is the space config object passed to the template
type spaceConf struct {
	FuncMap      map[string]interface{} `yaml:"funcMap"`
	PackageName  string                 `yaml:"packageName"`
	PackageDir   string                 `yaml:"packageDir"`
	Locales      []Locale               `yaml:"locales"`
	ContentTypes []ContentType          `yaml:"contentTypes"`
	ContentType  ContentType            `yaml:"contentType"`
	Version      string                 `yaml:"version"`
}

// GetLocales retrieves locale definition from Contentful
func getLocales(ctx context.Context, cma *contentful.Contentful, spaceID string) (locales []Locale, err error) {
	col, err := cma.Locales.List(ctx, spaceID).GetAll()
	if err != nil {
		log.Fatalf("Couldn't get locales: %v", err)
	}
	for _, item := range col.Items {
		var locale Locale
		byteArray, err := json.Marshal(item)
		if err != nil {
			break
		}
		err = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&locale)
		if err != nil {
			break
		}
		locales = append(locales, locale)
	}
	return
}

// GetContentTypes retrieves content type definition from Contentful
func getContentTypes(ctx context.Context, cma *contentful.Contentful, spaceID string) (contentTypes []ContentType, err error) {
	col := cma.ContentTypes.List(ctx, spaceID)
	col, err = col.GetAll()
	if err != nil {
		log.Fatal("Couldn't get content types")
	}
	for _, item := range col.Items {
		var contentType ContentType
		byteArray, err := json.Marshal(item)
		if err != nil {
			break
		}
		err = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&contentType)
		if err != nil {
			break
		}
		var filteredFields []ContentTypeField
		for _, field := range contentType.Fields {
			if !field.Omitted {
				filteredFields = append(filteredFields, field)
			}
		}
		contentType.Fields = filteredFields
		contentTypes = append(contentTypes, contentType)
	}
	sort.Slice(
		contentTypes, func(i, j int) bool {
			return contentTypes[i].Name < contentTypes[j].Name
		},
	)
	return
}

func getData(ctx context.Context, spaceID, cmaKey, environment, exportFile string, flagContentTypes []string) (
	finalContentTypes []ContentType, locales []Locale, err error,
) {
	var contentTypes []ContentType
	if exportFile != "" {
		fileBytes, err := os.ReadFile(exportFile)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error reading export file")
		}
		var export ExportFile
		err = json.Unmarshal(fileBytes, &export)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error parsing export file")
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
		locales, err = getLocales(ctx, CMA, spaceID)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not get locales")
		}
		fmt.Println("Locales found:", locales)

		// Get content types
		contentTypes, err = getContentTypes(ctx, CMA, spaceID)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not get content types")
		}
	}

	fmt.Println("Content types found:", len(contentTypes))

	finalContentTypes = []ContentType{}
	if len(flagContentTypes) == 0 {
		finalContentTypes = contentTypes
	} else {
		for _, ct := range contentTypes {
			if slices.Contains(flagContentTypes, ct.Sys.ID) {
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
func GenerateAPI(ctx context.Context, dir, packageName, spaceID, cmaKey, environment, exportFile string, flagContentTypes []string, version string) error {
	contentTypes, locales, errGetData := getData(ctx, spaceID, cmaKey, environment, exportFile, flagContentTypes)
	if errGetData != nil {
		return errors.Wrap(errGetData, "could not get data")
	}

	packageDir := filepath.Join(dir, packageName)
	errMkdir := os.MkdirAll(packageDir, 0o766)
	if errMkdir != nil {
		return errors.Wrap(errMkdir, "could not create target folder")
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
