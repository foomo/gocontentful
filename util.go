package erm

import "strings"

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
