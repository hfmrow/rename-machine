// miscFunc.go

package main

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	gi "github.com/hfmrow/renMachine/gtk3Import"
)

// receiveDnd: Called on drag & drop
func receiveDnd() {
	if !mainOptions.CumulativeDND {
		resetFilenamesStorage()
		mainOptions.primeFilesList = mainOptions.primeFilesList[:0]
	}
	// Init files list
	for _, value := range mainOptions.treeViewDropSet.FilesList {
		if !IsExistSl(mainOptions.primeFilesList, value) {
			mainOptions.primeFilesList = append(mainOptions.primeFilesList, value)
		}
	}
	restorePrimeFilesList()
	mainOptions.ScanSubDir = false
	mainOptions.ShowDirectory = false
	// mainOptions.UpdateObjects()
	Check(scanDirs())
	refreshLists()
}

// Increment: Add incrementation to filename
func Increment(inSL []string, separator string, startAt int, atLeft ...bool) (outSL []string) {
	var tmpStr string
	var left bool
	if len(atLeft) != 0 {
		left = atLeft[0]
	}
	digits := fmt.Sprintf("%d", len(inSL)+startAt)
	regRepl := regexp.MustCompile(`[[:digit:]]`)
	digits = regRepl.ReplaceAllString(digits, "0")
	for idx, line := range inSL {
		inc := fmt.Sprintf("%d", idx+startAt)
		inc = digits[len(inc):] + inc
		switch left {
		case true:
			tmpStr = inc + separator + line
		case false:
			tmpStr = line + separator + inc
		}
		outSL = append(outSL, tmpStr)
	}
	return outSL
}

// BuildExtSlice: Convert Entry mask to string slice
func BuildExtSlice(entry *gtk.Entry) (ExtMask []string) {
	tmpSliceStrings := strings.Split(getEntryText(entry), mainOptions.ExtSep)
	if len(tmpSliceStrings[0]) == 0 {
		ExtMask = append(ExtMask, "*")
	}
	for _, str := range tmpSliceStrings {
		str = strings.TrimSpace(str)
		if len(str) != 0 {
			ExtMask = append(ExtMask, str)
		}
	}
	entry.SetText(strings.Join(ExtMask, mainOptions.ExtSep+" "))
	return ExtMask
}

// Scan directories
func scanDirs() (err error) {
	var ext = []string{"*"}
	var fi os.FileInfo
	var ok bool

	mainOptions.rawFilesList = mainOptions.rawFilesList[:0]
	// According extension mask to tab position
	switch mainOptions.currentTab {
	case 0:
		ext = mainOptions.ExtMaskRen
	case 2:
		ext = mainOptions.ExtMask
	}
	for _, file := range mainOptions.primeFilesList {
		if fi, err = os.Stat(file); err == nil {
			if fi.IsDir() {
				if mainOptions.ShowDirectory {
					mainOptions.rawFilesList = append(mainOptions.rawFilesList, file)
				}
				if mainOptions.ScanSubDir {
					err = FindDir(file, ext, &mainOptions.rawFilesList,
						mainOptions.ScanSubDir, mainOptions.ShowDirectory, true)
				}
			} else {
				for _, msk := range ext {
					if ok, err = filepath.Match(msk, filepath.Ext(file)); ok {
						mainOptions.rawFilesList = append(mainOptions.rawFilesList, file)
						break
					}
					if err != nil {
						return err
					}
				}
			}
		} /* else {
			Check(err)
		}*/
	}
	return err
}

// FindDir retrieve file in a specific directory with more options.
func FindDir(dir string, masks []string, returnedStrSlice *[]string, scanSub, showDir, followSymlinkDir bool) (err error) {
	var ok bool
	var fName string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("%s\nDir-Name: %s\n", err.Error(), dir))
	}
	for _, file := range files {
		fName = filepath.Join(dir, file.Name())
		if followSymlinkDir { // Check for symlink ..
			file, err = os.Lstat(fName)
			if err != nil {
				return errors.New(fmt.Sprintf("%s\nFilename: %s\n", err.Error(), fName))
			}
			if file.Mode()&os.ModeSymlink != 0 { // Is a symlink ?
				fName, err := os.Readlink(fName) // Then read it...
				if err != nil {
					return errors.New(fmt.Sprintf("%s\nFilename: %s\n", err.Error(), fName))
				}
				file, err = os.Stat(fName) // Get symlink infos.
				if err != nil {
					return errors.New(fmt.Sprintf("%s\nFilename: %s\n", err.Error(), fName))
				}
				fName = filepath.Join(dir, file.Name())
			}
		}
		// Recursive play if it's a directory
		if file.IsDir() && scanSub {
			tmpFileList := new([]string)
			err = FindDir(fName, masks, tmpFileList, scanSub, showDir, followSymlinkDir)
			if err != nil {
				return errors.New(fmt.Sprintf("%s\nFilename: %s\n", err.Error(), fName))
			}
			*returnedStrSlice = append(*returnedStrSlice, *tmpFileList...)
		}

		for _, msk := range masks {
			ok, err = filepath.Match(msk, file.Name())
			if err != nil {
				return err
			}
			if ok {
				break
			}
		}
		if ok {
			if showDir { // Limit display directories if requested
				*returnedStrSlice = append(*returnedStrSlice, fName)
			} else {
				_, err = ioutil.ReadDir(fName)
				if err != nil {
					*returnedStrSlice = append(*returnedStrSlice, fName)
				}
			}
		}
	}
	return nil
}

// Restoring originals selected files
func restorePrimeFilesList() {
	mainOptions.rawFilesList = mainOptions.rawFilesList[:0]
	for _, file := range mainOptions.primeFilesList {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			mainOptions.rawFilesList = append(mainOptions.rawFilesList, file)
		}
	}
}

// resetEntriesAndRe: Clean entries from rename and title pages and refresh treeviews
func resetEntriesAndReBuildTreeView() (err error) {
	mainObjects.RenRemEntry.SetText("")
	mainObjects.RenRemEntry1.SetText("")
	mainObjects.RenRemEntry2.SetText("")
	mainObjects.RenReplEntry.SetText("")
	mainObjects.RenReplEntry1.SetText("")
	mainObjects.RenReplEntry2.SetText("")
	mainObjects.RenWthEntry.SetText("")
	mainObjects.RenWthEntry1.SetText("")
	mainObjects.RenWthEntry2.SetText("")
	mainObjects.RenIncrementChk.SetActive(false)

	if err = scanDirs(); err == nil {
		resetFilenamesStorage()
		refreshLists()
		RenRemEntryFocusOut()
	}
	return err
}

func resetFilenamesStorage() {
	mainOptions.renFilenames = mainOptions.renFilenames[:0]
	mainOptions.renFilenamesBak = mainOptions.renFilenamesBak[:0]
	mainOptions.titlFilenames = mainOptions.titlFilenames[:0]
	mainOptions.titlList = mainOptions.titlList[:0]
	mainOptions.moveFilesList = mainOptions.moveFilesList[:0]
	mainOptions.orgFilenames = mainOptions.orgFilenames[:0]
}

// Convert classic filelist into decomposed version path, name, ext
func decomposeFileList(inList []string) (outList [][]string) {
	var outDir, outFiles [][]string
	if len(inList) != 0 {
		if mainOptions.PreserveExt {
			for _, file := range inList {
				if fileInfo, _ := os.Stat(file); !fileInfo.IsDir() {
					outFiles = append(outFiles,
						[]string{filepath.Dir(file),
							BaseNoExt(filepath.Base(file)),
							filepath.Ext(file), "File"})
				} else if mainOptions.ShowDirectory {
					outDir = append(outDir,
						[]string{filepath.Dir(file),
							BaseNoExt(filepath.Base(file)),
							filepath.Ext(file), "Dir"})
				}
			}
		} else {
			for _, file := range inList {
				if fileInfo, _ := os.Stat(file); !fileInfo.IsDir() {
					outFiles = append(outFiles,
						[]string{filepath.Dir(file),
							filepath.Base(file),
							"", "File"})
				} else if mainOptions.ShowDirectory {
					outDir = append(outDir,
						[]string{filepath.Dir(file),
							filepath.Base(file),
							"", "Dir"})
				}
			}
		}
		// Sort string preserving order ascend
		sort.SliceStable(outDir, func(i, j int) bool { return strings.ToLower(outDir[i][2]) > strings.ToLower(outDir[j][2]) })
		sort.SliceStable(outFiles, func(i, j int) bool { return strings.ToLower(outFiles[i][1]) < strings.ToLower(outFiles[j][1]) })
		outList = append(outList, outFiles...)
		outList = append(outList, outDir...)
	}
	return outList

}

// Refresh treeviews
func refreshLists() {
	mainOptions.orgFilenames = decomposeFileList(mainOptions.rawFilesList)
	// Display filelist
	if mainObjects.fileListstore != nil {
		fillListstore(mainObjects.fileListstore, mainOptions.orgFilenames, true)
		mainObjects.renListstore.Clear()
		mainObjects.titleListstore.Clear()
		mainObjects.moveListstore.Clear()
	}
	updateStatusBar()
}

// Populate liststore to display filelist
func fillListstore(listStore *gtk.ListStore, inList [][]string, fullpath ...bool) {
	var seeFullPath bool
	var element []string
	if len(fullpath) > 0 {
		seeFullPath = fullpath[0]
	}
	// Cleaning liststore
	listStore.Clear()
	for _, file := range inList {
		if seeFullPath {
			element = append(element, file[1])
			element = append(element, file[3])
			element = append(element, filepath.Join(file[0], file[1]+file[2]))
		} else {
			element = append(element, file[1])
		}
		gi.ListStoreAddRow(listStore, element)
		element = element[:0]
	}
	updateStatusBar()
}

// Compile data to make new titled filename
func makeTitleFilename(inFiles [][]string, inTitle []string, toAdd ...string) (outFiles [][]string, err error) {
	var addB, addA string
	for idx, add := range toAdd {
		switch idx {
		case 0:
			addB = add
		case 1:
			addA = add
		}
	}

	lengthTitle := len(inTitle)
	lengthFiles := len(inFiles)
	if lengthTitle <= lengthFiles {
		err = errors.New("Too many files: " + strconv.Itoa(lengthFiles) + "\nNot enought titles: " + strconv.Itoa(lengthTitle))

	}
	if len(inTitle) != 0 {
		for idx, title := range inTitle {
			if idx >= lengthFiles {
				return outFiles, errors.New("Too many titles: " + strconv.Itoa(lengthTitle) + "\nNot enought files: " + strconv.Itoa(lengthFiles))
			}
			outFiles = append(outFiles,
				[]string{inFiles[idx][0],
					addB + inFiles[idx][1] + title + addA,
					inFiles[idx][2]})
		}
	} else if len(addB) != 0 || len(addA) != 0 {
		for idx, _ := range inFiles {
			outFiles = append(outFiles,
				[]string{inFiles[idx][0],
					addB + inFiles[idx][1] + addA,
					inFiles[idx][2]})
		}
	}
	return outFiles, err
}

// StatusBar function ...
func updateStatusBar() {
	wordFile := " File"
	filesCount := len(mainOptions.orgFilenames)
	if filesCount > 2 {
		wordFile = wordFile + "s"
	}
	displayStatusBar(fmt.Sprint(filesCount, wordFile))
}

// StatusBar function ...
func displayStatusBar(str ...string) {
	var outText []string
	if len(str) != 0 {
		for _, toDisp := range str {
			outText = append(outText, toDisp)
		}
		contextId1 := mainObjects.Statusbar.GetContextId("part1")
		mainObjects.Statusbar.Push(contextId1, strings.Join(outText, " | "))
	}
}

// getEntryText: retrieve value of an entry control.
func getEntryText(e *gtk.Entry) (outTxt string) {
	outTxt, err = e.GetText()
	Check(err, "getEntryText")
	return outTxt
}

// Check: Display error messages in HR version with onClickJump enabled in
// my favourite Golang IDE editor. Return true if error exist.
func Check(err error, message ...string) (state bool) {
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`) //	to match 2 or more whitespace symbols inside a string
	var msgs string
	if err != nil {
		state = true
		if len(message) != 0 { // Make string with messages if exists
			for _, mess := range message {
				msgs += `[` + mess + `]`
			}
		}
		pc, file, line, ok := runtime.Caller(1) //	(pc uintptr, file string, line int, ok bool)
		if ok == false {                        // Remove "== false" if needed
			fName := runtime.FuncForPC(pc).Name()
			fmt.Printf("[%s][%s][File: %s][Func: %s][Line: %d]\n", msgs, err.Error(), file, fName, line)
		} else {
			stack := strings.Split(fmt.Sprintf("%s", debug.Stack()), "\n")
			for idx := 5; idx < len(stack)-1; idx = idx + 2 {
				//	To match 2 or more whitespace leading/ending/inside a string (include \t, \n)
				mess1 := strings.Join(strings.Fields(stack[idx]), " ")
				mess2 := strings.TrimSpace(remInside.ReplaceAllString(stack[idx+1], " "))
				fmt.Printf("%s[%s][%s]\n", msgs, err.Error(), strings.Join([]string{mess1, mess2}, "]["))
			}
		}
	}
	return state
}

// BaseNoExt: get only the name without ext.
func BaseNoExt(filename string) (outFilename string) {
	outFilename = filepath.Base(filename)
	ext := filepath.Ext(outFilename)
	return outFilename[:len(outFilename)-len(ext)]
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
	for _, row := range platforms {
		if row[0] == runtime.GOOS {
			return row[2]
		}
	}
	return "\n"
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

// TrimSpace: Some multiple way to trim strings. cmds is optionnal or accept multiples args
func TrimSpace(inputString string, cmds ...string) (newstring string, err error) {

	osForbiden := regexp.MustCompile(`[<>:"/\\|?*]`)
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)    //	to match 2 or more whitespace symbols inside a string
	remInsideNoTab := regexp.MustCompile(`[\p{Zs}]{2,}`) //	(preserve \t) to match 2 or more space symbols inside a string

	if len(cmds) != 0 {
		for _, command := range cmds {
			switch command {
			case "+h": //	Escape html
				inputString = html.EscapeString(inputString)
			case "-h": //	UnEscape html
				inputString = html.UnescapeString(inputString)
			case "+e": //	Escape specials chars
				inputString = fmt.Sprintf("%q", inputString)
			case "-e": //	Un-Escape specials chars
				tmpString, err := strconv.Unquote(`"` + inputString + `"`)
				if err != nil {
					return inputString, err
				}
				inputString = tmpString
			case "-w": //	Change all illegals chars (for path in linux and windows) into "-"
				inputString = osForbiden.ReplaceAllString(inputString, "-")
			case "+w": //	clean all illegals chars (for path in linux and windows)
				inputString = osForbiden.ReplaceAllString(inputString, "")
			case "-c": //	Trim [[:space:]] and clean multi [[:space:]] inside
				inputString = strings.TrimSpace(remInside.ReplaceAllString(inputString, " "))
			case "-ct": //	Trim [[:space:]] and clean multi [[:space:]] inside (preserve TAB)
				inputString = strings.Trim(remInsideNoTab.ReplaceAllString(inputString, " "), " ")
			case "-s": //	To match 2 or more whitespace leading/ending/inside a string (include \t, \n)
				inputString = strings.Join(strings.Fields(inputString), " ")
			case "-&": //	Replace ampersand CHAR with ampersand HTML code
				inputString = strings.Replace(inputString, "&", "&amp;", -1)
			case "+&": //	Replace ampersand HTML code with ampersand CHAR
				inputString = strings.Replace(inputString, "&amp;", "&", -1)
			default:
				return inputString, errors.New("TrimSpace, " + command + ", does not exist")
			}
		}
	}
	return inputString, nil
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

// CmpRemSl2d: Compare slice1 and slice2 if row exist on both,
// they will be returned.
func CmpSl2d(sl1, sl2 [][]string) (outSlice []string) {
	var joined string
	for _, row2 := range sl2 {
		for _, row1 := range sl1 {
			if reflect.DeepEqual(row1, row2) {
				joined = filepath.Join(row2[0], row2[1]) + row2[2]
				if !IsExistSl(outSlice, joined) {
					outSlice = append(outSlice, joined)
				}
			}
		}
	}
	return outSlice
}

// CheckForDupSl2d: Return list of duplicated row is found
func CheckForDupSl2d(sl [][]string) (outSlice []string) {
	for idx := 0; idx < len(sl)-1; idx++ {
		for secIdx := idx + 1; secIdx < len(sl); secIdx++ {
			if reflect.DeepEqual(sl[idx], sl[secIdx]) {
				outSlice = append(outSlice, filepath.Join(sl[secIdx][0], sl[secIdx][1])+sl[secIdx][2])
			}
		}
	}
	return outSlice
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
