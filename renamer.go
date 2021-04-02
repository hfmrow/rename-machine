// renamer.go

/// +build ignore

/*
	©2018 H.F.M. MIT license
*/

package main

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	glss "github.com/hfmrow/genLib/slices"
	glsscc "github.com/hfmrow/genLib/strings/cClass"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gitl "github.com/hfmrow/gtk3Import/tools"
)

func doRename(list *[][]string) {
	var removed, renamed string
	var name = make([]string, 3)
	var inList []string
	List := *list
	// Remove part
	for idx := 0; idx < len(List); idx++ {
		name = List[idx]
		removed = ranameIt(name[1], modifName.remove, []string{"", "", ""})
		List[idx] = []string{name[0], removed, name[2]}
	}
	// Rename part
	for idx := 0; idx < len(List); idx++ {
		name = List[idx]
		renamed = ranameIt(name[1], modifName.replace, modifName.with)
		List[idx] = []string{name[0], renamed, name[2]}
	}
	// Increment part
	if mainObjects.RenIncrementChk.GetActive() {
		for idx := 0; idx < len(List); idx++ {
			name = List[idx]
			inList = append(inList, name[1])
		}
		separator, _ := mainObjects.RenIncSepEntry.GetText()
		inList = Increment(inList, separator, mainObjects.RenIncSpinbutton.GetValueAsInt(), !mainObjects.RenIncrementRightChk.GetActive())
		for idx := 0; idx < len(List); idx++ {
			List[idx] = []string{name[0], inList[idx], name[2]}
		}
	}
	fillListstore(tvsRenam, List)
}

// Initialize modifName structure wth entry boxes
func getEntryDatas() {
	//	osForbiden := regexp.MustCompile(`[<>:"/\\|?*]`)
	osForbReplacement := "-"
	// Remove
	modifName.remove = modifName.remove[:0]
	modifName.remove = append(modifName.remove, OsForbiden(gitl.GetEntryText(mainObjects.RenRemEntry), osForbReplacement))
	modifName.remove = append(modifName.remove, OsForbiden(gitl.GetEntryText(mainObjects.RenRemEntry1), osForbReplacement))
	modifName.remove = append(modifName.remove, OsForbiden(gitl.GetEntryText(mainObjects.RenRemEntry2), osForbReplacement))

	// Replace
	modifName.replace = modifName.replace[:0]
	modifName.replace = append(modifName.replace, OsForbiden(gitl.GetEntryText(mainObjects.RenReplEntry), osForbReplacement))
	modifName.replace = append(modifName.replace, OsForbiden(gitl.GetEntryText(mainObjects.RenReplEntry1), osForbReplacement))
	modifName.replace = append(modifName.replace, OsForbiden(gitl.GetEntryText(mainObjects.RenReplEntry2), osForbReplacement))

	// With
	modifName.with = modifName.with[:0]
	modifName.with = append(modifName.with, OsForbiden(gitl.GetEntryText(mainObjects.RenWthEntry), osForbReplacement))
	modifName.with = append(modifName.with, OsForbiden(gitl.GetEntryText(mainObjects.RenWthEntry1), osForbReplacement))
	modifName.with = append(modifName.with, OsForbiden(gitl.GetEntryText(mainObjects.RenWthEntry2), osForbReplacement))
}

// Remove function
func ranameIt(inStr string, replace, with []string) (outStr string) {
	var removerRegexp *regexp.Regexp
	outStr = inStr
	for idx, rem := range replace {
		if len(rem) != 0 {
			rem = regexp.QuoteMeta(rem)
			if mainOptions.CaseSensitive {
				removerRegexp = regexp.MustCompile(rem)
			} else {
				removerRegexp = regexp.MustCompile("(?i)" + rem)
			}
			outStr = removerRegexp.ReplaceAllString(outStr, with[idx])
		}
	}
	return outStr
}

// Substract function
func substractIt(inStr, toSub string, caseSensitive, posixCC, posixCCStrictMode bool) (outStr string) {
	var toSubReg *regexp.Regexp
	var err error
	outStr = inStr

	if !posixCC {
		toSub = regexp.QuoteMeta(toSub)

		if !caseSensitive {
			toSubReg, err = regexp.Compile("(?i)" + toSub)
		} else {
			toSubReg, err = regexp.Compile(toSub)
		}
	} else {
		toSubReg, err = regexp.Compile(glsscc.StringToCharacterClasses(toSub, caseSensitive, posixCCStrictMode))
	}
	Check(err)
	// Do the substract operation
	tmpOutStr := toSubReg.FindStringSubmatch(inStr)
	if len(tmpOutStr) != 0 {
		outStr = tmpOutStr[0]
	}
	return outStr
}

// Regex function
func regexIt(inStr, regexStr, replWith string) (outStr string) {
	outStr = inStr
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		gidg.DialogMessage(mainObjects.MainWindow,
			"error",
			"Regex error !",
			err.Error(),
			"", "ok")
	}
	// Do the regex operation
	return regex.ReplaceAllString(inStr, replWith)
}

// Keep between function
func keepBetweenIt(inStr string, toRemove []string, caseSensitive, posixCC, posixCCStrictMode bool) (outStr string) {
	var toKeepBtwReg *regexp.Regexp
	var err error
	outStr = inStr
	after := toRemove[0]
	before := toRemove[1]

	if !posixCC {
		after = regexp.QuoteMeta(after)
		before = regexp.QuoteMeta(before)

		if !caseSensitive {
			toKeepBtwReg, err = regexp.Compile("(?i)" + after + "(.*?)" + before)
		} else {
			toKeepBtwReg, err = regexp.Compile(after + "(.*?)" + before)
		}
	} else {
		toKeepBtwReg, err = regexp.Compile(glsscc.StringToCharacterClasses(after,
			caseSensitive, posixCCStrictMode) + "(.*?)" + glsscc.StringToCharacterClasses(before,
			caseSensitive, posixCCStrictMode))
	}
	Check(err)
	// Do the Keep between operation
	tmpOutStr := toKeepBtwReg.FindStringSubmatch(inStr)

	if len(tmpOutStr) == 2 {
		outStr = tmpOutStr[1]
	}
	return outStr
}

// Renamer function
func renameMe(from, to [][]string, fromTitle ...bool) (errList []string, err error) {
	var title bool
	var newName, oldName string
	var lengthTo int
	if len(fromTitle) > 0 {
		title = fromTitle[0]
	}

	if title {
		lengthTo = len(glss.CmpRemSl2d(from, to))
	} else {
		lengthTo = len(from)
	}

	if len(from) >= lengthTo && gidg.DialogMessage(mainObjects.MainWindow, "alert", sts["confirm"],
		"\n"+sts["proceed"]+strconv.Itoa(lengthTo)+sts["files"], "", sts["cancel"], sts["ok"]) == 1 {

		for idx := 0; idx < lengthTo; idx++ {
			if len(to[idx][1]+to[idx][2]) > 0 {
				newName = filepath.Join(to[idx][0], to[idx][1]+to[idx][2])
				oldName = filepath.Join(from[idx][0], from[idx][1]+from[idx][2])
				if oldName != newName {
					if _, err = os.Stat(newName); !os.IsNotExist(err) {
						// The case where samba is used over Window Os, when we have case sensitive
						// comparison that say identical filename, we check for content to see if
						// they're the same.

						// if glco.Md5File(newName) == glco.Md5File(oldName) {
						// 	if err = os.Rename(oldName, oldName+"~-~"); err == nil {
						// 		err = os.Rename(oldName+"~-~", newName)
						// 	}
						// } else {
						err = errors.New(newName + sts["fileExist"])
						// }
					} else {
						err = os.Rename(oldName, newName)
					}
				}

				if err != nil {
					errList = append(errList, filepath.Join(from[idx][0], from[idx][1]+from[idx][2]))
				}
			} else {
				errList = append(errList, sts["emptyname"]+": "+filepath.Join(from[idx][0], from[idx][1]+from[idx][2]))
			}
		}
	} else {
		return errList, errors.New(errCancel)
	}
	return errList, err
}

// Combined rename function
func renameAndRefresh(from, to [][]string) {
	var outString string
	var tmpSliceDup, tmpSliceExist []string
	var tmp2dSl [][]string
	if len(from) != 0 && len(to) != 0 {
		// Check for duplicate or existing files.
		tmpSliceDup = CheckForDupSl2dReturnDup(to)
		FindDir(to[0][0], []string{"*"}, &tmpSliceExist, false, false, false)
		tmp2dSl = decomposeFileList(tmpSliceExist)
		tmpSliceExist = CmpSl2dReturnDup(tmp2dSl, to)
		if len(tmpSliceExist) != 0 {
			outString = "\n" + sts["alreadyExist"] + "\n" + strings.Join(tmpSliceExist, "\n")
		}
		if len(tmpSliceDup) != 0 {
			outString += sts["dupFile"] + "\n" + strings.Join(tmpSliceDup, "\n")
		}

		if len(outString) == 0 {
			errList, err := renameMe(from, to)
			if err != nil {
				if err.Error() == errCancel {
					return
				}
				gidg.DialogMessage(mainObjects.MainWindow, "error", sts["mstk"],
					"\n"+sts["renErr"]+"\n"+err.Error()+"\n"+strings.Join(errList, "\n"),
					"", sts["ok"])
			} else {
				if len(errList) != 0 {
					gidg.DialogMessage(mainObjects.MainWindow, "error", sts["mstk"],
						"\n"+sts["renErr"]+"\n"+strings.Join(errList, "\n"),
						"", sts["ok"])
				}
				mainOptions.rawFilesList = mainOptions.rawFilesList[:0]
				mainOptions.primeFilesList = mainOptions.primeFilesList[:0]
				for _, file := range to {
					mainOptions.rawFilesList = append(mainOptions.rawFilesList, filepath.Join(file[0], file[1]+file[2]))
					mainOptions.primeFilesList = append(mainOptions.primeFilesList, filepath.Join(file[0], file[1]+file[2]))
				}
				if err := resetEntriesAndReBuildTreeView(); err != nil {
					gidg.DialogMessage(mainObjects.MainWindow, "alert", sts["mstk"], sts["errFiles"], "", sts["ok"])
				}
			}
		}
	}
	if len(outString) > 0 {
		gidg.DialogMessage(mainObjects.MainWindow, "alert", sts["mstk"], outString, "", sts["ok"])
	}
}
