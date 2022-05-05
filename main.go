package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/foomo/gocontentful/erm"
)

var VERSION = "v1.0.10"

type contentfulRc struct {
	ManagementToken string `json:"managementToken"`
}

var Usage = func() {
	fmt.Printf("\nSYNOPSIS\n")
	fmt.Printf("     gocontentful -spaceid SpaceID -cmakey CMAKey [-contenttypes firsttype,secondtype...lasttype] path/to/target/package\n\n")
	flag.PrintDefaults()
	fmt.Printf("\nNotes:\n")
	fmt.Println("- The last segment of the path/to/target/package will be used as package name")
	fmt.Println("- The -cmakey parameter can be omitted if you logged in with the Contentful CLI")
	fmt.Println()
}

func usageError(comment string) {
	fmt.Println("ERROR:", comment)
	Usage()
	os.Exit(1)
}

func fatal(infos ...interface{}) {
	fmt.Println(infos...)
	os.Exit(1)
}

func getCmaKeyFromRcFile() string {
	currentUser, errGetUser := user.Current()
	if errGetUser != nil {
		return ""
	}
	contentfulRcBytes, errReadFile := ioutil.ReadFile(currentUser.HomeDir + "/.contentfulrc.json")
	if errReadFile != nil {
		return ""
	}
	var contentfulConfig contentfulRc
	errUnmarshal := json.Unmarshal(contentfulRcBytes, &contentfulConfig)
	if errUnmarshal != nil {
		return ""
	}
	return contentfulConfig.ManagementToken
}

func main() {
	// Get parameters from cmd line flags
	flagSpaceID := flag.String("spaceid", "", "Contentful space ID")
	flagCMAKey := flag.String("cmakey", "", "[Optional] Contentful CMA key")
	flagEnvironment := flag.String("environment", "", "[Optional] Contentful space environment")
	flagGenerateFromExport := flag.String("exportfile", "", "Space export file to generate the API from")
	flagContentTypes := flag.String("contenttypes", "", "[Optional] Content type IDs to parse, comma separated")
	flagVersion := flag.Bool("version", false, "Print version and exit")
	flagHelp := flag.Bool("help", false, "Print version and exit")
	flag.Parse()

	if *flagVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *flagHelp {
		Usage()
		os.Exit(0)
	}
	cmaKey := *flagCMAKey
	if cmaKey == "" && *flagGenerateFromExport == "" {
		cmaKey = getCmaKeyFromRcFile()
	}
	if *flagGenerateFromExport == "" && (*flagSpaceID == "" || cmaKey == "") ||
		*flagGenerateFromExport != "" && (*flagSpaceID != "" || cmaKey != "") {
		usageError("Please specify either a Contentful Space ID and CMA access token or an export file name")
	}

	if len(flag.Args()) != 1 {
		usageError("Missing arg path/to/target/package")
	}

	path := flag.Arg(0)
	packageName := filepath.Base(path)

	matched, err := regexp.MatchString(`[a-z].{2,}`, packageName)
	if !matched {
		usageError("Please specify the package name correctly (only small caps letters)")
	}

	fmt.Printf("Contentful API Generator %s starting...\n\n", VERSION)

	var flagContentTypesSlice []string
	if *flagContentTypes != "" {
		for _, contentType := range strings.Split(*flagContentTypes, ",") {
			flagContentTypesSlice = append(flagContentTypesSlice, strings.TrimSpace(contentType))
		}
	}

	err = erm.GenerateAPI(filepath.Dir(path), packageName, *flagSpaceID, cmaKey, *flagEnvironment, *flagGenerateFromExport, flagContentTypesSlice)
	if err != nil {
		fatal("Something went horribly wrong...", err)
	}
	fmt.Println("ALL DONE!")

}
