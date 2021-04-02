// miscFunc.go

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	glss "github.com/hfmrow/genLib/slices"

	gitl "github.com/hfmrow/gtk3Import/tools"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

// receiveDnd: Called on drag & drop
func receiveDnd() {
	if !mainOptions.CumulativeDND {
		resetFilenamesStorage()
		mainOptions.primeFilesList = mainOptions.primeFilesList[:0]
	}
	// Init files list
	for _, value := range *mainOptions.treeViewDropSet.FilesList {
		if !glss.IsExistSl(mainOptions.primeFilesList, value) {
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
	tmpSliceStrings := strings.Split(gitl.GetEntryText(entry), mainOptions.ExtSep)
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
		}
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
	if tvsFiles.ListStore != nil {
		fillListstore(tvsFiles, mainOptions.orgFilenames, true)
		tvsRenam.ListStore.Clear()
		tvsTitle.ListStore.Clear()
		tvsMoves.ListStore.Clear()
	}
	updateStatusBar()
}

// Populate liststore to display filelist
func fillListstore(tvs *gitw.TreeViewStructure, inList [][]string, fullpath ...bool) {
	if tvs.ListStore != nil {
		var seeFullPath bool
		var element []string
		if len(fullpath) > 0 {
			seeFullPath = fullpath[0]
		}
		// Cleaning liststore
		tvs.ListStore.Clear()
		for _, file := range inList {
			if seeFullPath {
				element = append(element, file[1])
				element = append(element, file[3])
				element = append(element, filepath.Join(file[0], file[1]+file[2]))
			} else {
				element = append(element, file[1])
			}
			tvs.AddRow(nil, tvs.ColValuesStringSliceToIfaceSlice(element...)...)
			element = element[:0]
		}
		updateStatusBar()
	}
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

// CmpSl2dReturnDup: Compare slice1 and slice2 if row exist on both,
// they will be returned.
func CmpSl2dReturnDup(sl1, sl2 [][]string) (outSlice []string) {
	var joined string
	for _, row2 := range sl2 {
		for _, row1 := range sl1 {
			if reflect.DeepEqual(row1, row2) {
				joined = filepath.Join(row2[0], row2[1]) + row2[2]
				if !glss.IsExistSl(outSlice, joined) {
					outSlice = append(outSlice, joined)
				}
			}
		}
	}
	return outSlice
}

// CheckForDupSl2dReturnDup: Return list of duplicated rows if found.
func CheckForDupSl2dReturnDup(sl [][]string) (outSlice []string) {
	for idx := 0; idx < len(sl)-1; idx++ {
		for secIdx := idx + 1; secIdx < len(sl); secIdx++ {
			if reflect.DeepEqual(sl[idx], sl[secIdx]) {
				outSlice = append(outSlice, filepath.Join(sl[secIdx][0], sl[secIdx][1])+sl[secIdx][2])
			}
		}
	}
	return outSlice
}
