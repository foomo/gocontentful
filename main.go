package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/foomo/gocontentful/erm"
)

func fatal(comment string) {
	fmt.Println("ERROR:", comment)
	fmt.Printf("\nSYNOPSIS\n")
	fmt.Printf("     gocontentful -spaceid SpaceID -cmakey CMAKey [-contenttypes firsttype,secondtype...lasttype] path/to/target/package\n\n")
	flag.Usage()
	fmt.Printf("\nNote: The last segment of the path/to/target/package will be used as package name\n\n")
	os.Exit(1)
}

func main() {
	fmt.Printf("Contentful API Generator starting...\n\n")
	// Get parameters from cmd line flags
	flagSpaceID := flag.String("spaceid", "", "Contentful space ID")
	flagCMAKey := flag.String("cmakey", "", "Contentful CMA key")
	flagContentTypes := flag.String("contenttypes", "", "[Optional] Content type IDs to parse, comma separated")

	flag.Parse()

	if *flagSpaceID == "" || *flagCMAKey == "" {
		fatal("Please specify the Contentful space ID and access Key")
	}

	if len(flag.Args()) != 1 {
		fatal("Missing arg path/to/target/package")
	}

	path := flag.Arg(0)
	packageName := filepath.Base(path)

	if matched, _ := regexp.MatchString(`[a-z].{2,}`, packageName); !matched {
		fatal("Please specify the package name correctly (only small caps letters)")
	}

	var flagContentTypesSlice []string
	if *flagContentTypes != "" {
		for _, contentType := range strings.Split(*flagContentTypes, ",") {
			flagContentTypesSlice = append(flagContentTypesSlice, strings.TrimSpace(contentType))
		}
	}

	if err := erm.GenerateAPI(
		filepath.Dir(path),
		packageName,
		*flagSpaceID,
		*flagCMAKey,
		flagContentTypesSlice,
	); err != nil {
		log.Fatal("generating API:", err)
	}
	fmt.Println("ALL DONE!")

}
