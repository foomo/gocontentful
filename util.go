package erm

import "strings"

import "regexp"

// SliceIncludes returns true if slice includes string
func SliceIncludes(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

func firstCap(inputString string) (outputString string) {
	outputString = strings.Title(inputString)
	return
}

func onlyLetters(inputString string) (outputString string) {
	re := regexp.MustCompile("[^A-Za-z]")
	return re.ReplaceAllString(inputString, "")
}
