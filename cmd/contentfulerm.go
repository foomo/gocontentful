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
	log.Println("Contentful Entry-Reference Mapping Generator starting...")
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
		for _, contentType := range strings.Split(*flagContentTypes, ",") {
			flagContentTypesSlice = append(flagContentTypesSlice, strings.TrimSpace(contentType))
		}
	}

	// Get client
	CMA := contentful.NewCMA(*flagCMAKey)
	CMA.Debug = false

	// Get space locales
	locales, err := erm.GetLocales(CMA, flagSpaceID)
	if err != nil {
		log.Fatal("Could not get locales")
	}
	log.Println("Locales found:", locales)

	// Get content types
	contentTypes, err := erm.GetContentTypes(CMA, flagSpaceID)
	if err != nil {
		log.Fatal("Could not get content types")
	}
	log.Println("Content types found:", len(contentTypes))

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
	log.Println("Filtered Content types:", len(filteredContentTypes))

	err = erm.ProcessSpace(*flagPackage, locales, filteredContentTypes)
	if err != nil {
		log.Fatal("Something went horribly wrong...", err)
	}
	log.Println("ALL DONE!")

}
