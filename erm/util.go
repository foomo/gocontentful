package erm

import (
	"regexp"
)

var onlyLettersRegex = regexp.MustCompile("[^A-Za-z]")

func sliceIncludes(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

func onlyLetters(inputString string) (outputString string) {
	return onlyLettersRegex.ReplaceAllString(inputString, "")
}
