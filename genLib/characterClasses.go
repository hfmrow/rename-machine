// characterClasses.go

package genLib

import (
	"regexp"
)

type cClass struct {
	regex string
}

// Convert char to char classes. strictMode mean mostly classes is applied, if set to false, matching will be made with less precision.
func charToCharClasses(inpString string, caseSensitive, strictMode bool) string {
	var charClassesCS []cClass
	var charClassesCI []cClass
	if strictMode {
		charClassesCS = []cClass{
			cClass{regex: `[[:upper:]]`},
			cClass{regex: `[[:lower:]]`},
			cClass{regex: `[[:digit:]]`},
			cClass{regex: `[[:blank:]]`},
			cClass{regex: `[[:space:]]`},
			cClass{regex: `[[:cntrl:]]`},
			cClass{regex: `[[:punct:]]`},
			cClass{regex: `[[:graph:]]`}}
		charClassesCI = []cClass{
			cClass{regex: `[[:alpha:]]`},
			cClass{regex: `[[:digit:]]`},
			cClass{regex: `[[:blank:]]`},
			cClass{regex: `[[:space:]]`},
			cClass{regex: `[[:cntrl:]]`},
			cClass{regex: `[[:punct:]]`},
			cClass{regex: `[[:graph:]]`}}
	} else {
		charClassesCS = []cClass{
			cClass{regex: `[[:upper:]]`},
			cClass{regex: `[[:lower:]]`},
			cClass{regex: `[[:digit:]]`},
			cClass{regex: `[[:punct:]]`}}
		charClassesCI = []cClass{
			cClass{regex: `[[:alpha:]]`},
			cClass{regex: `[[:digit:]]`},
			cClass{regex: `[[:punct:]]`}}
	}
	var tmpString string
	switch caseSensitive {
	case true: // Case sensitive matching ...
		for _, regex := range charClassesCS {
			regExpCompiled := regexp.MustCompile(regex.regex)
			if regExpCompiled.MatchString(inpString) {
				tmpString = regex.regex
				break
			} else {
				tmpString = charClassesCS[len(charClassesCS)-1].regex
			}
		}
	case false: // Case insensitive matching ...
		for _, regex := range charClassesCI {
			regExpCompiled := regexp.MustCompile(regex.regex)
			if regExpCompiled.MatchString(inpString) {
				tmpString = regex.regex
				break
			} else {
				tmpString = charClassesCI[len(charClassesCI)-1].regex
			}
		}
	}
	return tmpString
}

// Convert string to character classes equivalent's string
func StringToCharacterClasses(inpString string, caseSensitive, strictMode bool) (outString string) {
	bytesString := []byte(inpString)
	for _, charByte := range bytesString {
		outString += charToCharClasses(string(charByte), caseSensitive, strictMode)
	}
	return outString
}
