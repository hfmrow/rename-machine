// sliceOperations.go

package slices

import (
	"fmt"
	"reflect"

	gst "github.com/hfmrow/genLib/strings"
)

// GetStrIndex: Get index of a string in a slice, Return -1 if no entry found ...
func GetStrIndex(slice []string, item string) int {
	for idx, row := range slice {
		if row == item {
			return idx
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
	for idx, row := range slice {
		if row[col] == value {
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

// IsExist2d Search in 2d string slice if a row exist (deepequal row).
func IsExist2d(slice [][]string, cmpRow []string) bool {
	for _, mainRow := range slice {
		if reflect.DeepEqual(mainRow, cmpRow) {
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
		inValue = gst.RemoveNonAlNum(inValue)
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
