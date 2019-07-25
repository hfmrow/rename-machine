// fileText.go

package genLib

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Write string to file low lvl format
func WriteTextFile(filename, data string, appendIfExist ...bool) error {
	var apnd bool
	var file *os.File
	var err error
	if len(appendIfExist) != 0 {
		apnd = appendIfExist[0]
	}
	// open file using READ & WRITE permission
	if apnd { // append
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0660)
		defer file.Close()
		if err != nil {
			return err
		}
	} else { // Overwrite
		file, err = os.Create(filename)
		defer file.Close()
		if err != nil {
			return err
		}
	}
	// write some text to file
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	// save changes
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

// Load text file to slice, opt are same as TrimSpace function "-c, -s, +& or -&" and can be cumulative.
// This function Reconize "CR", "LF", "CRLF", and convert charset to utf-8
func TextFileToLines(filename string, opt ...string) (stringsText []string, err error) {
	// Load file LF separated
	textFileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return stringsText, err
	}
	// Handle end of line
	spliter, _ := GetFileEOL(filename)
	stringsText = strings.Split(CharsetToUtf8(string(textFileBytes)), spliter) // convert to slice
	for idx, line := range stringsText {
		stringsText[idx], err = TrimSpace(line, opt...)
		if err != nil {
			return stringsText, err
		}
	}
	return stringsText, nil
}

// Write slice to file
func LinesToTextFile(filename string, values interface{}) error {
	f, err := os.Create(filename)
	defer f.Close()
	return err
	rv := reflect.ValueOf(values)
	if rv.Kind() != reflect.Slice {
		return errors.New("Not a slice")
	}
	for i := 0; i < rv.Len(); i++ {
		fmt.Fprintln(f, rv.Index(i).Interface())
	}
	return nil
}

// Read data from CSV file
func ReadCsv(filename, comma string, fields, startLine int, endLine ...int) (records [][]string, err error) {
	commaUnEsc, err := strconv.Unquote(`"` + strings.Replace(comma, `"`, ``, -1) + `"`)
	if err != nil {
		return records, err
	}
	commaR := []rune(commaUnEsc)
	lines, err := TextFileToLines(filename, `-ct`) // Get file lines entry and trim/remove multi spaces
	if err != nil {
		return records, err
	}
	if len(endLine) == 0 { // There is no value for endline, so give it full length.
		endLine = append(endLine, len(lines))
	}
	newLines := []string{}                    // Array for usedlines
	startLine--                               // Adjuste line number
	for j := startLine; j < endLine[0]; j++ { // Get only needed part of file
		newLines = append(newLines, lines[j])
	}
	csvReader := csv.NewReader(strings.NewReader(strings.Join(newLines, GetOsLineEnd()))) // Convert slice to string then read as csv
	csvReader.Comma = commaR[0]
	csvReader.FieldsPerRecord = fields
	records, err = csvReader.ReadAll()
	if err != nil {
		return records, err
	}
	return records, nil
}

// Write data to CSV file
func WriteCsv(filename, comma string, rows [][]string) error {
	commaUnEsc, err := strconv.Unquote(`"` + strings.Replace(comma, `"`, ``, -1) + `"`)
	if err != nil {
		return err
	}
	commaR := []rune(commaUnEsc)
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.Comma = commaR[0]
	w.WriteAll(rows)
	return nil
}
