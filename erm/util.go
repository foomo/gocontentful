package erm

import (
	"regexp"
)

func sliceIncludes(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

func onlyLetters(inputString string) (outputString string) {
	re := regexp.MustCompile("[^A-Za-z]")
	return re.ReplaceAllString(inputString, "")
}
