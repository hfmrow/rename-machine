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

	g "github.com/hfmrow/renMachine/genLib"
	gi "github.com/hfmrow/renMachine/gtk3Import"
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
	fillListstore(mainObjects.renListstore, List)
}

// Initialize modifName structure wth entry boxes
func getEntryDatas() {
	var entry string
	var err error
	// Remove
	modifName.remove = modifName.remove[:0]
	entry, err = mainObjects.RenRemEntry.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.remove = append(modifName.remove, entry)
	entry, err = mainObjects.RenRemEntry1.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.remove = append(modifName.remove, entry)
	entry, err = mainObjects.RenRemEntry2.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.remove = append(modifName.remove, entry)
	// Replace
	modifName.replace = modifName.replace[:0]
	entry, err = mainObjects.RenReplEntry.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.replace = append(modifName.replace, entry)
	entry, err = mainObjects.RenReplEntry1.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.replace = append(modifName.replace, entry)
	entry, err = mainObjects.RenReplEntry2.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.replace = append(modifName.replace, entry)
	// With
	modifName.with = modifName.with[:0]
	entry, err = mainObjects.RenWthEntry.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.with = append(modifName.with, entry)
	entry, err = mainObjects.RenWthEntry1.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.with = append(modifName.with, entry)
	entry, err = mainObjects.RenWthEntry2.GetText()
	g.Check(err)
	entry, _ = g.TrimSpace(entry, "-w")
	modifName.with = append(modifName.with, entry)
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
		toSubReg, err = regexp.Compile(g.StringToCharacterClasses(toSub, caseSensitive, posixCCStrictMode))
	}
	g.Check(err)
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
		gi.DlgMessage(mainObjects.MainWindow,
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
		toKeepBtwReg, err = regexp.Compile(g.StringToCharacterClasses(after,
			caseSensitive, posixCCStrictMode) + "(.*?)" + g.StringToCharacterClasses(before,
			caseSensitive, posixCCStrictMode))
	}
	g.Check(err)
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
		lengthTo = len(g.CmpRemSl2d(from, to))
	} else {
		lengthTo = len(from)
	}

	if len(from) >= lengthTo && gi.DlgMessage(mainObjects.MainWindow, "alert", sts["confirm"],
		"\n"+sts["proceed"]+strconv.Itoa(lengthTo)+sts["files"], "", sts["cancel"], sts["ok"]) == 1 {

		for idx := 0; idx < lengthTo; idx++ {
			if len(to[idx][1]+to[idx][2]) > 0 {
				newName = filepath.Join(to[idx][0], to[idx][1]+to[idx][2])
				oldName = filepath.Join(from[idx][0], from[idx][1]+from[idx][2])
				if oldName != newName {
					if _, err = os.Stat(newName); !os.IsNotExist(err) {
						err = errors.New(newName + sts["fileExist"])
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
	if len(from) != 0 && len(to) != 0 {
		if !g.CheckForDupSl2d(to) {
			errList, err := renameMe(from, to)
			if err != nil {
				if err.Error() == errCancel {
					return
				}
				gi.DlgMessage(mainObjects.MainWindow, "error", sts["mstk"],
					"\n"+sts["renErr"]+"\n"+err.Error()+"\n"+strings.Join(errList, "\n"),
					"", sts["ok"])
			} else {
				if len(errList) != 0 {
					gi.DlgMessage(mainObjects.MainWindow, "error", sts["mstk"],
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
					gi.DlgMessage(mainObjects.MainWindow, "alert", sts["mstk"], sts["errFiles"], "", sts["ok"])
				}
			}
		} else {
			gi.DlgMessage(mainObjects.MainWindow, "alert", sts["mstk"], sts["couldnotproceed"], "", sts["ok"])
		}
	}
}
