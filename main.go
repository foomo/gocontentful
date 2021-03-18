package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/foomo/contentful-erm/erm"
)

func usageError(comment string) {
	fmt.Println(comment)
	flag.Usage()
	os.Exit(1)
}

func fatal(infos ...interface{}) {
	fmt.Println(infos...)
	os.Exit(1)
}

func main() {
	fmt.Println("Contentful Entry-Reference Mapping Generator starting...")
	// Get parameters from cmd line flags
	flagSpaceID := flag.String("spaceid", "", "Contentful space ID")
	flagCMAKey := flag.String("cmakey", "", "Contentful CMA key")
	flagContentTypes := flag.String("contenttypes", "", "[Optional] Content type IDs to parse, comma separated")

	flag.Parse()

	if *flagSpaceID == "" || *flagCMAKey == "" {
		usageError("Please specify the Contentful space ID and access Key")
	}

	if len(flag.Args()) != 1 {
		usageError("missing arg path/to/target/package")
	}

	path := flag.Arg(0)
	packageName := filepath.Base(path)

	matched, err := regexp.MatchString(`[a-z].{2,}`, packageName)
	if !matched {
		usageError("Please specify the package name correctly (only small caps letters)")
	}

	var flagContentTypesSlice []string
	if *flagContentTypes != "" {
		for _, contentType := range strings.Split(*flagContentTypes, ",") {
			flagContentTypesSlice = append(flagContentTypesSlice, strings.TrimSpace(contentType))
		}
	}

	err = erm.GenerateAPI(filepath.Dir(path), packageName, *flagSpaceID, *flagCMAKey, flagContentTypesSlice)
	if err != nil {
		fatal("Something went horribly wrong...", err)
	}
	fmt.Println("ALL DONE!")

}
