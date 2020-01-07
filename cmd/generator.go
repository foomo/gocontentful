package main

import (
	"flag"
	"log"
	"strings"

	"github.com/foomo/contentful"
	erm "github.com/foomo/contentful-erm"
)

func main() {
	// Get Space ID and CMA Key from cmd line flags
	flagSpaceID := flag.String("spaceid", "", "Contentful space ID")
	flagCMAKey := flag.String("cmakey", "", "Contentful CMA key")
	flagContentTypes := flag.String("contenttypes", "", "[Optional] Content type IDs to parse, comma separated")
	flag.Parse()
	var flagContentTypesSlice []string
	if *flagContentTypes != "" {
		flagContentTypesSlice = strings.Split(*flagContentTypes, ",")
		log.Println("flagConteTypesSlice:", flagContentTypesSlice)
	}
	if *flagSpaceID == "" || *flagCMAKey == "" {
		flag.Usage()
		log.Fatal("You have to specify the cmd parameters correctly")
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

	err = erm.ProcessSpace(locales, filteredContentTypes)
	if err != nil {
		log.Fatal("Something went horribly wrong...", err)
	}
}
