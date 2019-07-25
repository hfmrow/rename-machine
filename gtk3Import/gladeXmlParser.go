// gladeXmlParser.go

/// +build ignore

/*
	GladeXmlParser v2.0 ©2019 H.F.M. MIT license

	parseGladeXmlFile: Create a parsed glade structure containing
	all objects with property, signals and packing informations.
*/

package gtk3Import

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var Name = "GladeXmlParser"
var Vers = "v2.0"
var Descr = "Golang parser for glade xml files"
var Creat = "H.F.M"
var YearCreat = "2019"
var LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
var LicenseAbrv = "License (MIT)"
var LicenseUrl = "https://opensource.org/licenses/mit-license.php"
var Repository = "github.com/..."

type GtkInterface struct {
	GXPVersion    string
	ObjectsCount  int
	UpdatedOn     string
	GladeFilename string
	Requires      requires
	Objects       []GtkObject
	Comments      []string
	// private
	objectsLoaded   bool
	skipNoId        bool
	skipLoweAtFirst bool
	idxLine         int
	prevIdxLine     int
	xmlSource       []string
	getValuesReg    *regexp.Regexp
}

type GtkObject struct {
	Class    string
	Id       string
	Property []GtkProps
	Signal   []GtkProps
	Packing  []GtkProps
}

type GtkProps struct {
	Name         string
	Value        string
	Swapped      string
	Translatable string
}
type requires struct {
	Lib     string
	Version string
}

// parseGladeXmlFile: Create new parsed glade structure containing
// all objects with property, signals, packing informations.
func GladeXmlParserNew(filename string, skipNoId, skipLoweAtFirst bool) (out *GtkInterface, err error) {
	iFace := new(GtkInterface)
	iFace.GladeFilename = filename
	iFace.GXPVersion = fmt.Sprintf("%s %s", Name, Vers)
	iFace.getValuesReg = regexp.MustCompile(`"(.*?)"|>(.*?)<`)
	iFace.skipNoId = skipNoId
	iFace.skipLoweAtFirst = skipLoweAtFirst

	err = iFace.readGladeXmlFile()
	if err == nil {
		iFace.parseGladeXmlFile()
		iFace.ObjectsCount = len(iFace.Objects)
		return iFace, err
	}
	return nil, err
}

// readGladeXmlFile:
func (iFace *GtkInterface) readGladeXmlFile() (err error) {
	var ok bool
	data, err := ioutil.ReadFile(iFace.GladeFilename)
	if err == nil {
		for _, line := range strings.Split(string(data), GetTextEOL(data)) {
			if strings.Contains(line, `<requires lib="gtk+"`) {
				ok = true
				break
			}
		}
		if ok {
			iFace.xmlSource, err = iFace.sanitizeXml(data)
			if err == nil {
				iFace.objectsLoaded = true
			}
		} else {
			return errors.New("Bad file format: " + filepath.Base(iFace.GladeFilename))
		}
	}
	return err
}

// parseGladeXmlFile:
func (iFace *GtkInterface) parseGladeXmlFile() {
	var line, tmpID string
	var values []string
	for iFace.idxLine = 0; iFace.idxLine < len(iFace.xmlSource); iFace.idxLine++ {
		values, line = iFace.getValues()
		if len(values) > 1 { // Check for ID
			tmpID = values[1]
		} else {
			tmpID = ""
		}
		if !(iFace.skipNoId && len(tmpID) == 0) { // Exclude no name objects if requires
			if !(iFace.skipLoweAtFirst && LowercaseAtFirst(tmpID)) { // Exclude lower at first char objects if requires
				switch {
				case strings.Contains(line, `<requires lib="`): // requires
					iFace.Requires.Lib = values[0]
					iFace.Requires.Version = tmpID
				case strings.Contains(line, `<object class="`): // object
					iFace.idxLine++
					var tmpGtkObj GtkObject
					tmpGtkObj = iFace.readObject(GtkObject{
						Class:    values[0],
						Id:       tmpID,
						Property: []GtkProps{},
						Signal:   []GtkProps{},
						Packing:  []GtkProps{}})
					iFace.Objects = append(iFace.Objects, tmpGtkObj)
				}
			}

			/* DEBUG purpose: In a case where DLV debugger can not recover the value of some slice structures,here:
			"Property", "Signal", "Packing" look like empty, it is really boring. But values still available in code ...*/
			// var tmpGtkObj GtkObject
			// fmt.Printf("%#v\n", iFace.Objects[len(iFace.Objects)-1].Property)
			// fmt.Printf("%#v\n", iFace.Objects[len(iFace.Objects)-1].Packing)
			// theObject := iFace.readObject(GtkObject{Class: values[0], Id: values[1], Property: []GtkProps{}, Signal: []GtkProps{}, Packing: []GtkProps{}})
			// tmpGtkObj.Class = theObject.Class
			// tmpGtkObj.Id = theObject.Id
			// tmpGtkObj.Property = make([]GtkProps, len(theObject.Property))
			// tmpGtkObj.Signal = make([]GtkProps, len(theObject.Signal))
			// tmpGtkObj.Packing = make([]GtkProps, len(theObject.Packing))
			// copy(tmpGtkObj.Property, theObject.Property)
			// copy(tmpGtkObj.Signal, theObject.Signal)
			// copy(tmpGtkObj.Packing, theObject.Packing)
			// iFace.Objects = append(iFace.Objects, tmpGtkObj)

		}
	}
}

// readObject: read object, property and packing.
func (iFace *GtkInterface) readObject(inObj GtkObject) (outObj GtkObject) {
	var newObjectCount int
	ok := true
	var line string
	var values []string
	iFace.prevIdxLine = iFace.idxLine
	for iFace.idxLine = iFace.idxLine; iFace.idxLine < len(iFace.xmlSource); iFace.idxLine++ {
		values, line = iFace.getValues()
		switch ok {
		case true:
			inObj, ok = iFace.readProps(line, values, inObj)
		case false:
			switch {
			case strings.Contains(line, `<object class="`): // object
				newObjectCount++
			case strings.Contains(line, `</object>`): // object
				newObjectCount--
			case newObjectCount == -1: // object end
				iFace.idxLine = iFace.prevIdxLine
				return inObj
			case strings.Contains(line, `<packing>`): // <packing>
				iFace.idxLine++
				for iFace.idxLine = iFace.idxLine; iFace.idxLine < len(iFace.xmlSource); iFace.idxLine++ {
					values, line = iFace.getValues()
					inObj, ok = readPacking(line, values, inObj)
					if !ok {
						iFace.idxLine = iFace.prevIdxLine
						return inObj
					}
				}
			}

		}
	}
	return outObj
}

// readProps:
func (iFace *GtkInterface) readProps(line string, values []string, inObj GtkObject) (outObj GtkObject, ok bool) {
	switch {
	case strings.Contains(line, `<property name="`): // property
		ok = true
		switch {
		case strings.Contains(line, `translatable="`): // property and translatable
			if len(values) > 2 {
				inObj.Property = append(inObj.Property, GtkProps{Name: values[0], Translatable: values[1], Value: values[2]})
			} else {
				fmt.Printf("%s: %d was skipped.\n", "error reading xml file at line", iFace.idxLine+1)
			}
		default: // property only
			if len(values) > 1 {
				inObj.Property = append(inObj.Property, GtkProps{Name: values[0], Value: values[1]})
			} else {
				fmt.Printf("%s: %d was skipped.\n", "error reading xml file at line", iFace.idxLine+1)
			}
		}
	case strings.Contains(line, `<signal name="`): // signal
		ok = true
		if len(values) > 2 {
			inObj.Signal = append(inObj.Signal, GtkProps{Name: values[0], Value: values[1], Swapped: values[2]})
		} else {
			fmt.Printf("%s: %d was skipped.\n", "error reading xml file at line", iFace.idxLine+1)
		}
	}
	if !ok {
		iFace.prevIdxLine = iFace.idxLine
	}
	return inObj, ok
}

// readPacking:
func readPacking(line string, values []string, inObj GtkObject) (outObj GtkObject, ok bool) {
	switch {
	case strings.Contains(line, `<property name="`): // property
		ok = true
		inObj.Packing = append(inObj.Packing, GtkProps{Name: values[0], Value: values[1]})
	}
	return inObj, ok
}

// getValues: get all values from a line and clean them
func (iFace *GtkInterface) getValues() (values []string, line string) {
	line = iFace.xmlSource[iFace.idxLine]
	tmpValues := iFace.getValuesReg.FindAllStringSubmatch(line, -1)
	for _, v := range tmpValues {
		values = append(values, strings.Trim(strings.Trim(strings.Trim(v[0], `"`), `>`), `<`))
	}
	return values, line
}

// SanitizeXml: Escape eol when multilines text are found,
// give an error if file is not a valid glade xml format.
func (iFace *GtkInterface) sanitizeXml(inBytes []byte) (newString []string, err error) {
	var row string
	regComments := regexp.MustCompile(`<!--(.*|\n+)+-->`)
	regCommentLine := regexp.MustCompile(`<!--(.*?)-->`)
	commentsBytes := regComments.FindAll(inBytes, -1)
	commentLinesBytes := regCommentLine.FindAll(inBytes, -1)
	for _, com := range commentsBytes {
		iFace.Comments = append(iFace.Comments, string(com))
	}
	for _, com := range commentLinesBytes {
		iFace.Comments = append(iFace.Comments, string(com))
	}

	// remove comments from XML file
	inBytes = regCommentLine.ReplaceAll(inBytes, []byte(""))
	inBytes = regComments.ReplaceAll(inBytes, []byte(""))

	regStart := regexp.MustCompile(`^(<)`)
	regPropertyAtStart := regexp.MustCompile(`^(</property>)`)
	regSpaceTab := regexp.MustCompile(`\s`)
	newString = strings.Split(string(inBytes), GetTextEOL(inBytes))
	// Replace LF inside labels or hints with "\n"
	for idx := len(newString) - 1; idx >= 0; idx-- {
		row = newString[idx]
		if len(row) != 0 {
			if !regStart.MatchString(regSpaceTab.ReplaceAllString(row, "")) || regPropertyAtStart.MatchString(regSpaceTab.ReplaceAllString(row, "")) {
				newString = append(newString[:idx], newString[idx+1:]...)
				newString[idx-1] += `\n` + row
			}
		}
	}
	// // Search for errors in glade xml file
	// var detectValidGladeXmlFile = func() error {
	// 	var newObjectCount int
	// 	enclosures := [][]string{
	// 		{"<child>", "</child>"},
	// 		{"<object class=", "</object>"},
	// 		{"<packing>", "</packing>"},
	// 		{"<interface>", "</interface>"}}
	// 	for _, enc := range enclosures {
	// 		for _, line := range newString {
	// 			switch {
	// 			case strings.Contains(line, enc[0]): // object
	// 				newObjectCount++
	// 			case strings.Contains(line, enc[1]): // object
	// 				newObjectCount--
	// 			}
	// 		}
	// 		if newObjectCount != 0 {
	// 			return errors.New("XML file format does not match Glade requirement ...")
	// 		}
	// 		newObjectCount = 0
	// 	}
	// 	return nil
	// }
	// return newString, detectValidGladeXmlFile()
	return newString, nil
}

var platforms = [][]string{
	{"darwin", "386", "\n"},
	{"darwin", "amd64", "\n"},
	{"dragonfly", "amd64", "\n"},
	{"freebsd", "386", "\n"},
	{"freebsd", "amd64", "\n"},
	{"freebsd", "arm", "\n"},
	{"linux", "386", "\n"},
	{"linux", "amd64", "\n"},
	{"linux", "arm", "\n"},
	{"linux", "arm64", "\n"},
	{"linux", "ppc64", "\n"},
	{"linux", "ppc64le", "\n"},
	{"linux", "mips", "\n"},
	{"linux", "mipsle", "\n"},
	{"linux", "mips64", "\n"},
	{"linux", "mips64le", "\n"},
	{"linux", "s390x", "\n"},
	{"nacl", "386", "\n"},
	{"nacl", "amd64p32", "\n"},
	{"nacl", "arm", "\n"},
	{"netbsd", "386", "\n"},
	{"netbsd", "amd64", "\n"},
	{"netbsd", "arm", "\n"},
	{"openbsd", "386", "\n"},
	{"openbsd", "amd64", "\n"},
	{"openbsd", "arm", "\n"},
	{"plan9", "386", "\n"},
	{"plan9", "amd64", "\n"},
	{"plan9", "arm", "\n"},
	{"solaris", "amd64", "\n"},
	{"windows", "386", "\r\n"},
	{"windows", "amd64", "\r\n"}}

// GetOsLineEnd: Get current OS line-feed
func GetOsLineEnd() string {
	return platforms[GetStrIndex2dCol(platforms, runtime.GOOS, 0)][2]
}

// GetTextEOL: Get EOL from text bytes (CR, LF, CRLF) > string or get OS line end.
func GetTextEOL(inTextBytes []byte) (outString string) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	if bytes.Contains(inTextBytes, bCRLF) {
		return string(bCRLF)
	} else if bytes.Contains(inTextBytes, bCR) {
		return string(bCR)
	}
	if bytes.Contains(inTextBytes, bLF) {
		return string(bLF)
	}
	return GetOsLineEnd()
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

type timeStamp struct {
	Year          string
	YearCopyRight string
	Month         string
	MonthWord     string
	Day           string
	DayWord       string
	Date          string
	Time          string
	Full          string
}

// Get current timestamp
func TimeStamp() *timeStamp {
	ts := new(timeStamp)
	timed := time.Now()
	regD := regexp.MustCompile("([^[:digit:]])")
	regA := regexp.MustCompile("([^[:alpha:]])")
	splitedNum := regD.Split(timed.Format(time.RFC3339), -1)
	splitedWrd := regA.Split(timed.Format(time.RFC850), -1)
	ts.Year = splitedNum[0]
	ts.Month = splitedNum[1]
	ts.Day = splitedNum[2]
	ts.Time = splitedNum[3] + `:` + splitedNum[4] + `:` + splitedNum[5]
	ts.DayWord = splitedWrd[0]
	ts.MonthWord = splitedWrd[5]
	ts.YearCopyRight = `©` + ts.Year
	ts.Full = strings.Join(strings.Split(timed.Format(time.RFC1123), " ")[:5], " ")
	return ts
}

// Read Text Controls from file
func (iFace *GtkInterface) ReadFile(filename string) (err error) {
	err = JsonRead(filename, iFace)
	if err == nil {
		iFace.objectsLoaded = true
		return err
	} else {
		iFace.objectsLoaded = false
		return err
	}
}

// Write Text Controls to file
func (iFace *GtkInterface) WriteFile(filename string) error {
	iFace.ObjectsCount = len(iFace.Objects)
	iFace.UpdatedOn = TimeStamp().Full
	return JsonWrite(filename, iFace)
}

// JsonRead: datas from file to given interface / structure
// i.e: err := ReadJson(filename, &person)
// remember to put upper char at left of variables names to be saved.
func JsonRead(filename string, interf interface{}) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(textFileBytes, &interf)
	}
	return err
}

// JsonWrite: datas to file from given interface / structure
// i.e: err := WriteJson(filename, &person)
// remember to put upper char at left of variables names to be saved.
func JsonWrite(filename string, interf interface{}) (err error) {
	var out bytes.Buffer
	var jsonData []byte
	if jsonData, err = json.Marshal(&interf); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			if err = ioutil.WriteFile(filename, out.Bytes(), 0644); err == nil {
			}
		}
	}
	return err
}
