package erm

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func sliceIncludes(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

func firstCap(inputString string) (outputString string) {
	outputString = cases.Title(language.English, cases.Option(cases.NoLower)).String(inputString)
	outputString = strings.ReplaceAll(outputString, "_", "")
	return
}

func onlyLetters(inputString string) (outputString string) {
	re := regexp.MustCompile("[^A-Za-z]")
	return re.ReplaceAllString(inputString, "")
}
