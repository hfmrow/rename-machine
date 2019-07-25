// miscFunc.go

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	g "github.com/hfmrow/renMachine/genLib"
	gi "github.com/hfmrow/renMachine/gtk3Import"
)

// receiveDnd: Called on drag & drop
func receiveDnd() {
	dnd = true
	mainOptions.primeFilesList = mainOptions.primeFilesList[:0]
	resetFilenamesStorage()
	// Init files list
	for _, value := range mainOptions.treeViewDropSet.FilesList {
		mainOptions.primeFilesList = append(mainOptions.primeFilesList, value)
	}
	restorePrimeFilesList()
	mainOptions.ScanSubDir = false
	mainOptions.ShowDirectory = false
	mainOptions.UpdateObjects()
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

// BuildExtSlice
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
	switch mainOptions.currentTab {
	case 0:
		ext = mainOptions.ExtMaskRen
	case 2:
		ext = mainOptions.ExtMask
	}

	if dnd {
		for _, file := range mainOptions.primeFilesList {
			if fi, err = os.Stat(file); fi.IsDir() && err == nil {
				if mainOptions.ShowDirectory {
					mainOptions.rawFilesList = append(mainOptions.rawFilesList, file)
				}
				if mainOptions.ScanSubDir {
					err = FindDir(file, ext, &mainOptions.rawFilesList,
						mainOptions.ScanSubDir, mainOptions.ShowDirectory, true)
				}
			} else if err == nil {
				for _, msk := range ext {
					ok, err = filepath.Match(msk, filepath.Ext(file))
					if err != nil {
						return err
					}
					if ok {
						break
					}
				}
				if ok {
					mainOptions.rawFilesList = append(mainOptions.rawFilesList, file)
				}
			}
		}
		return err
	}

	return FindDir(mainOptions.baseDirectory, ext,
		&mainOptions.rawFilesList, mainOptions.ScanSubDir,
		mainOptions.ShowDirectory, true)

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
	mainObjects.TitleAddAEntry.SetText("")
	mainObjects.TitleAddBEntry.SetText("")
	mainObjects.TitleAddBFileEntry.SetText("")
	mainObjects.TitleSepEntry.SetText("")
	mainObjects.TitleSpinbutton.SetValue(0)
	err = scanDirs()
	resetFilenamesStorage()
	refreshLists()
	RenRemEntryFocusOut()
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
	var splited g.Filepath
	var outDir, outFiles [][]string
	if len(inList) != 0 {
		if mainOptions.PreserveExt {
			for _, file := range inList {
				splited = g.SplitFilepath(file)
				if !splited.IsDir {
					outFiles = append(outFiles,
						[]string{splited.Path,
							splited.BaseNoExt,
							splited.Ext, "File"})
				} else if mainOptions.ShowDirectory {
					outDir = append(outDir,
						[]string{filepath.Dir(file),
							splited.BaseNoExt,
							splited.Ext, "Dir"})
				}
			}
		} else {
			for _, file := range inList {
				splited = g.SplitFilepath(file)
				if !splited.IsDir {
					outFiles = append(outFiles,
						[]string{splited.Path,
							splited.Base,
							"", "File"})
				} else if mainOptions.ShowDirectory {
					outDir = append(outDir,
						[]string{splited.Path,
							splited.Base,
							"", "Dir"})
				}
			}
		}
		// Sort string preserving order ascend
		sort.SliceStable(outDir, func(i, j int) bool { return strings.ToLower(outDir[i][2]) > strings.ToLower(outDir[j][2]) })
		sort.SliceStable(outFiles, func(i, j int) bool { return strings.ToLower(outFiles[i][1]) < strings.ToLower(outFiles[j][1]) })
		outList = append(outList, outFiles...)
		outList = append(outList, outDir...)
		// sort.SliceStable(outList, func(i, j int) bool { return strings.ToLower(outList[i][1]) > strings.ToLower(outList[j][1]) })
		// sort.SliceStable(outList, func(i, j int) bool { return strings.ToLower(outList[i][3]) > strings.ToLower(outList[j][3]) })
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
		if fullpath[0] {
			seeFullPath = true
		}
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

	displayStatusBar(fmt.Sprint(filesCount, wordFile), fmt.Sprint("Base directory: ", mainOptions.baseDirectory))
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
