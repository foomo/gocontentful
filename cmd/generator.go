package main

import (
	"flag"
	"log"
	"regexp"
	"strings"

	"github.com/foomo/contentful"
	erm "github.com/foomo/contentful-erm"
)

func main() {
	// Get parameters from cmd line flags
	flagSpaceID := flag.String("spaceid", "", "Contentful space ID")
	flagCMAKey := flag.String("cmakey", "", "Contentful CMA key")
	flagPackage := flag.String("package", "", "Generated package name")
	flagContentTypes := flag.String("contenttypes", "", "[Optional] Content type IDs to parse, comma separated")

	flag.Parse()

	if *flagSpaceID == "" || *flagCMAKey == "" {
		flag.Usage()
		log.Fatal("Please specify the Contentful space ID and access Key")
	}

	matched, err := regexp.MatchString(`[a-z].{2,}`, *flagPackage)
	if !matched {
		flag.Usage()
		log.Fatal("Please specify the package name correctly (only small caps letters)")
	}

	var flagContentTypesSlice []string
	if *flagContentTypes != "" {
		flagContentTypesSlice = strings.Split(*flagContentTypes, ",")
	}

	// Get client
	CMA := contentful.NewCMA(*flagCMAKey)
	CMA.Debug = true

	// Get space locales
	locales, err := erm.GetLocales(CMA, flagSpaceID)
	if err != nil {
		log.Fatal("Could not get locales")
	}

	// Get content types
	contentTypes, err := erm.GetContentTypes(CMA, flagSpaceID)
	if err != nil {
		log.Fatal("Could not get locales")
	}

	filteredContentTypes := []erm.ContentType{}
	if len(flagContentTypesSlice) == 0 {
		filteredContentTypes = contentTypes
	} else {
		for _, ct := range contentTypes {
			if erm.SliceIncludes(flagContentTypesSlice, ct.Sys.ID) {
				filteredContentTypes = append(filteredContentTypes, ct)
			}
		}
	}

	err = erm.ProcessSpace(*flagPackage, locales, filteredContentTypes)
	if err != nil {
		log.Fatal("Something went horribly wrong...", err)
	}
}
