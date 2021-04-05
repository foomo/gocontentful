package erm

import "regexp"

func sliceIncludes(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

var reOnlyLetters = regexp.MustCompile("[^A-Za-z]")

func onlyLetters(s string) string {
	return reOnlyLetters.ReplaceAllString(s, "")
}
