// sliceOperations.go

package genLib

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// SplitAtEOL: split data to slice
func SplitAtEOL(data []byte) (outSlice []string) {
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	data = bytes.ReplaceAll(data, bCRLF, bLF)
	return strings.Split(string(data), string(bLF))
}

// SearchSl: Search in 2d string slice. cs=case sensitive, ww=whole word, rx=regex
func SearchSl(find string, table [][]string, caseSensitive, POSIXcharClass, POSIXstrictMode, regex, wholeWord bool) (out [][]string, err error) {
	if len(table) != 0 {
		if len(find) != 0 {
			search := find
			var outTable [][]string
			if POSIXcharClass {
				search = StringToCharacterClasses(find, caseSensitive, POSIXstrictMode)
			}
			search = `(` + search + `)`
			if wholeWord {
				search = `\b` + search + `\b`
			}
			if !caseSensitive && !POSIXcharClass {
				search = `(?i)` + search
			}
			regX, err := regexp.Compile(search)
			if err != nil {
				return out, err
			}
			fmt.Println("regex: " + search)
			if regex {
				regX, err = regexp.Compile(find)
				if err != nil {
					return out, err
				}
			}

			for idxRow := 0; idxRow < len(table); idxRow++ {
				for _, col := range table[idxRow] {
					if regX.MatchString(col) {
						outTable = append(outTable, table[idxRow])
						break // Avoid duplicate when element found twice in same row
					}
				}
			}
			if len(outTable) == 0 {
				return out, errors.New(find + "\n\nNot found ...")
			}
			return outTable, nil // Result found then return it
		}
	}
	return [][]string{}, errors.New("Nothing to search ...")
}

// GetStrIndex: Get index of a string in a slice, Return -1 if no entry found ...
func GetStrIndex(slice []string, item string) int {
	for i, _ := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// IsExistSl: if exist then  ...
func IsExistSl(slice []string, item string) bool {
	for _, mainRow := range slice {
		if mainRow == item {
			return true
		}
	}
	return false
}

// GetStrIndex2dCol: Search in 2d string slice if a column's value exist and return row number.
func GetStrIndex2dCol(slice [][]string, value string, col int) int {
	for idx, mainRow := range slice {
		if mainRow[col] == value {
			return idx
		}
	}
	return -1
}

// IsExist2dCol: Search in 2d string slice if a column's value exist.
func IsExist2dCol(slice [][]string, value string, col int) bool {
	for _, mainRow := range slice {
		if mainRow[col] == value {
			return true
		}
	}
	return false
}

// CmpRemSl2d: Compare slice1 and slice2 if row exist on both,
// the raw is removed from slice2 and result returned.
func CmpRemSl2d(sl1, sl2 [][]string) (outSlice [][]string) {
	var skip bool
	for _, row2 := range sl2 {
		for _, row1 := range sl1 {
			if reflect.DeepEqual(row1, row2) {
				skip = true
				break
			}
		}
		if !skip {
			outSlice = append(outSlice, row2)
		}
		skip = false
	}
	return outSlice
}

// IsExist2d Search in 2d string slice if a row exist.
func IsExist2d(slice [][]string, cmpRow []string) bool {
	for _, mainRow := range slice {
		if reflect.DeepEqual(mainRow, cmpRow) {
			return true
		}
	}
	return false
}

// CheckForDupSl2d: Return true if duplicated row is found
func CheckForDupSl2d(sl [][]string) bool {
	for idx := 0; idx < len(sl)-1; idx++ {
		for secIdx := idx + 1; secIdx < len(sl); secIdx++ {
			if reflect.DeepEqual(sl[idx], sl[secIdx]) {
				return true
			}
		}
	}
	return false
}

// RemoveDupSl: Remove duplicate entry in a string slice
// TODO rewrite this function using "reflect.DeepEqual"... on next use
func RemoveDupSl(slice []string) []string {
	var isExist bool
	tmpSlice := make([]string, 0)
	for _, inValue := range slice {
		isExist = false
		inValue = RemoveNonAlNum(inValue)
		for _, outValue := range tmpSlice {
			if outValue == inValue {
				isExist = true
				break
			}
		}
		if !isExist {
			tmpSlice = append(tmpSlice, inValue)
		}
	}
	fmt.Println(len(tmpSlice))
	return tmpSlice
}

// RemoveDupSl2d: Remove duplicate entry in a 2d string slice based on column number content.
// TODO rewrite this function using "reflect.DeepEqual"... on next use
func RemoveDupSl2d(slice [][]string, col int) (outSlice [][]string) {
	if len(slice) == 0 {
		return slice
	}
	var dupFlag bool
	outSlice = append(outSlice, slice[0])
	for primIdx := len(slice) - 1; primIdx >= 0; primIdx-- {
		dupFlag = false
		for secIdx := len(outSlice) - 1; secIdx >= 0; secIdx-- {
			if outSlice[secIdx][col] == slice[primIdx][col] {
				dupFlag = true
			}
		}
		if !dupFlag {
			outSlice = append(outSlice, slice[primIdx])
		}
	}
	return outSlice
}

// Preppend: Add data at the begining of a string slice
func Preppend(slice []string, prepend ...string) []string {
	return append(prepend, slice...)
}

// AppendAt: Add data at a specified position in slice of a string
func AppendAt(slice []string, pos int, insert ...string) []string {
	if pos > len(slice) {
		pos = len(slice)
	} else if pos < 0 {
		pos = 0
	}
	return append(slice[:pos], append(insert, slice[pos:]...)...)
}

// DeleteSl: Delete value at specified position in string slice
func DeleteSl(slice []string, pos int) []string {
	return append(slice[:pos], slice[pos+1:]...)
}

// DeleteSl1: Delete value at specified position in string slice
func DeleteSl1(slice []string, pos int) []string {
	copy(slice[pos:], slice[pos+1:])
	slice[len(slice)-1] = "" // or the zero value of T
	return slice[:len(slice)-1]
}

// SliceSortDate: Sort 2d string slice with date inside
func SliceSortDate(slice [][]string, fmtDate string, dateCol, secDateCol int, ascendant bool) [][]string {
	fieldsCount := len(slice[0]) // Get nb of columns
	var firstLine int
	var previous, after string
	var positiveidx, negativeidx int
	// compute unix date using given column numbers
	for idx := firstLine; idx < len(slice); idx++ {
		dateStr := FindDate(slice[idx][dateCol], fmtDate)
		if dateStr != nil { // search for 1st column
			slice[idx] = append(slice[idx], fmt.Sprintf("%d", FormatDate(fmtDate, dateStr[0]).Unix()))
		} else if secDateCol != -1 { // Check for second column if it was given
			dateStr = FindDate(slice[idx][secDateCol], fmtDate)
			if dateStr != nil { // If date was not found in 1st column, search for 2nd column
				slice[idx] = append(slice[idx], fmt.Sprintf("%d", FormatDate(fmtDate, slice[idx][secDateCol]).Unix()))
			} else { //  in case where none of the columns given contain date field, put null string if there is no way to find a date
				slice[idx] = append(slice[idx], ``)
			}
		} else { // put null string if there is no way to find a date
			slice[idx] = append(slice[idx], ``)
		}
	}
	// Ensure we always have a value in sorting field (get previous or next closer)
	for idx := firstLine; idx < len(slice); idx++ {
		if slice[idx][fieldsCount] == `` {
			for idxFind := firstLine + 1; idxFind < len(slice); idxFind++ {
				positiveidx = idx + idxFind
				negativeidx = idx - idxFind
				if positiveidx >= len(slice) { // Check index to avoiding 'out of range'
					positiveidx = len(slice) - 1
				}
				if negativeidx <= 0 {
					negativeidx = 0
				}
				after = slice[positiveidx][fieldsCount] // Get previous or next value
				previous = slice[negativeidx][fieldsCount]
				if previous != `` { // Set value, prioritise the previous one.
					slice[idx][fieldsCount] = previous
					break
				}
				if after != `` {
					slice[idx][fieldsCount] = after
					break
				}
			}
		}
	}
	tmpLines := make([][]string, 0)
	if ascendant != true {
		// Sort by unix date preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool { return slice[i][len(slice[i])-1] > slice[j][len(slice[i])-1] })
		for idx := firstLine; idx < len(slice); idx++ { // Store row count elements - 1
			tmpLines = append(tmpLines, slice[idx][:len(slice[idx])-1])
		}
	} else {
		// Sort by unix date preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool { return slice[i][len(slice[i])-1] < slice[j][len(slice[i])-1] })
		for idx := firstLine; idx < len(slice); idx++ { // Store row count elements - 1
			tmpLines = append(tmpLines, slice[idx][:len(slice[idx])-1])
		}
	}
	return tmpLines
}

// SliceSortString: Sort 2d string slice
func SliceSortString(slice [][]string, col int, ascendant, caseSensitive, numbered bool) {
	if numbered {
		var tmpWordList []string
		for _, wrd := range slice {
			tmpWordList = append(tmpWordList, wrd[col])
		}
		numberedWords := new(WordWithDigit)
		numberedWords.Init(tmpWordList)

		if ascendant != true {
			// Sort string preserving order descendant
			sort.SliceStable(slice, func(i, j int) bool {
				return numberedWords.FillWordToMatchMaxLength(slice[i][col]) > numberedWords.FillWordToMatchMaxLength(slice[j][col])
			})
		} else {
			// Sort string preserving order ascendant
			sort.SliceStable(slice, func(i, j int) bool {
				return numberedWords.FillWordToMatchMaxLength(slice[i][col]) < numberedWords.FillWordToMatchMaxLength(slice[j][col])
			})
		}
		return
	}

	toLowerCase := func(inString string) string {
		return inString
	}
	if !caseSensitive {
		toLowerCase = func(inString string) string { return strings.ToLower(inString) }
	}

	if ascendant != true {
		// Sort string preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool { return toLowerCase(slice[i][col]) > toLowerCase(slice[j][col]) })
	} else {
		// Sort string preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool { return toLowerCase(slice[i][col]) < toLowerCase(slice[j][col]) })
	}
}

// SliceSortFloat: Sort 2d string with float value
func SliceSortFloat(slice [][]string, col int, ascendant bool, decimalChar string) {
	if ascendant != true {
		// Sort string (float) preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool {
			return StringDecimalSwitchFloat(decimalChar, slice[i][col]) > StringDecimalSwitchFloat(decimalChar, slice[j][col])
		})
	} else {
		// Sort string (float) preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool {
			return StringDecimalSwitchFloat(decimalChar, slice[i][col]) < StringDecimalSwitchFloat(decimalChar, slice[j][col])
		})
	}
}

// Convert comma to dot if needed and return 0 if input string is empty.
func StringDecimalSwitchFloat(decimalChar, inString string) float64 {
	if inString == "" {
		inString = "0"
	}
	switch decimalChar {
	case ",":
		f, _ := strconv.ParseFloat(strings.Replace(inString, ",", ".", 1), 64)
		return f
	case ".":
		f, _ := strconv.ParseFloat(inString, 64)
		return f
	}
	return -1
}
