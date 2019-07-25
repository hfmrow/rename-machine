// csvfunc.go
package genLib

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type RowStore struct {
	Idx int
	Cnt int
	Tot int
	Str string
}

// Append to slice if not already exist (RowStore version)
func AppendIfMissing(inputSlice []RowStore, input RowStore) []RowStore {
	for _, element := range inputSlice {
		if element == input {
			return inputSlice
		}
	}
	return append(inputSlice, input)
}

// Append to slice if RowStore.cnt value does not already exist (Quote version)
func AppendIfMissingC(inputSlice []RowStore, input RowStore) []RowStore {
	for i, element := range inputSlice {
		if element.Cnt == input.Cnt {
			inputSlice[i].Tot = element.Tot + 1
			return inputSlice
		}
	}
	return append(inputSlice, input)
}

// Get count of non alphanum chars in a string in struct(RowStore) format
func FindCountStr(str, regExpression string) []RowStore {
	storeElements := []RowStore{}
	re := regexp.MustCompile(regExpression)
	submatchall := re.FindAllString(str, -1)
	for _, element := range submatchall {
		counted := strings.Count(str, element)
		storeElements = AppendIfMissing(storeElements, RowStore{counted, 0, 0, fmt.Sprintf(`%q`, element)})
	} // Sort slice higher to lower
	sort.Slice(storeElements, func(i, j int) bool {
		return storeElements[i].Idx > storeElements[j].Idx
	})
	return storeElements
}

// Get Quote char
func FindQuote(str string, sep string) []RowStore {
	mixedQuotes := []RowStore{}
	storeElements := []RowStore{}
	possiblesQuotes := []string{`"`, `'`}
	sepCln, _ := strconv.Unquote(sep)           // UnQuote char if needed
	for idx, quoteIn := range possiblesQuotes { // Create slice containing possibility quotes and comma
		mixedQuotes = append(mixedQuotes, RowStore{idx, 0, 0, sepCln + quoteIn})
		mixedQuotes = append(mixedQuotes, RowStore{idx, 0, 0, quoteIn + sepCln})
	}
	for _, quoteIn := range mixedQuotes {
		counted := strings.Count(str, quoteIn.Str)
		storeElements = append(storeElements, RowStore{quoteIn.Idx, counted, 0, fmt.Sprintf(`%q`, possiblesQuotes[quoteIn.Idx])})
	}
	return storeElements
}

// Get separator char
func GetSep(str string) []string {
	outputStr := []string{}
	storeQuotes := []RowStore{}                                      //		Exclude alpha, [:space:], ", ', % and numeric values:	`[^%^ ^"^'^[:alnum:]+]`
	storeElements := FindCountStr(str, `[^-: "'._/\\\*#[:alnum:]+]`) //		Exclude alpha, [:space:], ", ', / and numeric values:	`[^ ^"^'^/^[:alnum:]+]`

	for _, quote := range storeElements {
		storeQuotes = append(storeQuotes, FindQuote(str, quote.Str)...)
	} // Sort slice higher to lower
	sort.Slice(storeQuotes, func(i, j int) bool {
		return storeQuotes[i].Cnt > storeQuotes[j].Cnt
	})
	outputStr = append(outputStr, storeElements[0].Str) // Print out comma
	outputStr = append(outputStr, storeQuotes[0].Str)   // Print out Quote
	return outputStr
}
