package erm

import (
	"regexp"
)

var onlyLettersRegex = regexp.MustCompile("[^A-Za-z]")

func onlyLetters(inputString string) (outputString string) {
	return onlyLettersRegex.ReplaceAllString(inputString, "")
}
