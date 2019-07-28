// objectsHandlers.go

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	gi "github.com/hfmrow/renMachine/gtk3Import"
)

// Notebook page changed
func NotebookPageChanged(nb *gtk.Notebook, page *gtk.Widget, pageNum uint) {
	mainOptions.currentTab = pageNum
	// Restore options
	mainOptions.ScanSubDir = tempScanSubDir
	mainOptions.ShowDirectory = tempShowDirectory
	mainOptions.PreserveExt = tempPreserveExt
	mainOptions.UpdateObjects()

	switch pageNum {
	case 1: // TitlePage
		mainOptions.PreserveExt = true
		mainOptions.ShowDirectory = false
	case 2: // MovePage
		mainOptions.PreserveExt = false
		mainOptions.ShowDirectory = false
		mainOptions.ScanSubDir = true
	}

	Check(scanDirs())
	resetFilenamesStorage()
	refreshLists()
}

// MoveEntryExtMaskEnterKeyPressed: Update data and disp on enter pressed
func MoveEntryExtMaskEnterKeyPressed() bool {
	mainOptions.ExtMask = BuildExtSlice(mainObjects.MoveEntryExtMask)
	Check(scanDirs())
	resetFilenamesStorage()
	refreshLists()
	return false // GDK_EVENT_PROPAGATE signal
}

// RenEntryExtMaskEnterKeyPressed: Update data and disp on enter pressed
func RenEntryExtMaskEnterKeyPressed() bool {
	mainOptions.ExtMaskRen = BuildExtSlice(mainObjects.RenEntryExtMask)
	Check(scanDirs())
	resetFilenamesStorage()
	refreshLists()
	return false // GDK_EVENT_PROPAGATE signal
}

// Chk handling for preserve ext
func RenPresExtChkChanged() {
	mainOptions.PreserveExt = mainObjects.RenPresExtChk.GetActive()
	tempPreserveExt = mainOptions.PreserveExt
	resetFilenamesStorage()
	refreshLists()
	RenRemEntryFocusOut()
}

// OverCaseSensChkChanged checkboxes handling
func OverCaseSensChkChanged() {
	mainOptions.OvercaseSensitive = mainObjects.OverCaseSensChk.GetActive()
	if mainOptions.OvercaseSensitive {
		mainObjects.OverCharClassChk.SetActive(false)
	} else {
		mainObjects.OverCharClassChk.SetSensitive(true)
	}
}

// OverCharClassChkChanged checkboxes handling
func OverCharClassChkChanged() {
	mainOptions.OverPosixCharClass = mainObjects.OverCharClassChk.GetActive()
	if mainOptions.OverPosixCharClass {
		mainObjects.OverCharClassStrictModeChk.SetSensitive(true)
		mainObjects.OverCaseSensChk.SetActive(false)
	} else {
		mainObjects.OverCharClassStrictModeChk.SetSensitive(false)
		mainObjects.OverCharClassStrictModeChk.SetActive(false)
	}
}

// OverCharClassStrictModeChkChanged checkboxes handling
func OverCharClassStrictModeChkChanged() {
	mainOptions.OverPosixStrictMde = mainObjects.OverCharClassStrictModeChk.GetActive()
}

// RenIncrementChkClicked:
func RenIncrementChkClicked(chk *gtk.CheckButton) {
	mainObjects.RenIncrementRightChk.SetSensitive(chk.GetActive())
	mainObjects.RenIncSepEntry.SetSensitive(chk.GetActive())
	mainObjects.RenIncSpinbutton.SetSensitive(chk.GetActive())
	RenRemEntryFocusOut()
}

// Handling rename apply button
func RenApplyButtonClicked() {
	renameAndRefresh(mainOptions.orgFilenames, mainOptions.renFilenamesBak)
}

// RenIncrementRightChkClicked:
func RenIncrementRightChkClicked(chk *gtk.CheckButton) {
	RenRemEntryFocusOut()
}

// Handling CaseSensitve toggled
func RenCaseSensChkChanged() {
	mainOptions.CaseSensitive = mainObjects.RenCaseSensChk.GetActive()
	RenRemEntryFocusOut()
}

// Show dircetories toggled
func RenShowDirChkChanged(chk *gtk.CheckButton) {
	mainOptions.ShowDirectory = chk.GetActive()
	tempShowDirectory = mainOptions.ShowDirectory
	Check(scanDirs())
	resetFilenamesStorage()
	refreshLists()
	RenRemEntryFocusOut()
}

// Handling check subdir
func ScanSubDirChkChanged(chk *gtk.CheckButton) {
	mainOptions.ScanSubDir = chk.GetActive()
	mainObjects.RenScanSubDirChk.SetActive(mainOptions.ScanSubDir)
	mainObjects.TitleScanSubDirChk.SetActive(mainOptions.ScanSubDir)
	tempScanSubDir = mainOptions.ScanSubDir
	Check(scanDirs())
	resetFilenamesStorage()
	refreshLists()
	RenRemEntryFocusOut()
}

// Handling button for choosing directory to move files to
func MoveFilechooserButtonClicked() {
	mainOptions.moveFilesList = mainOptions.moveFilesList[:0]
	// mainOptions.primeFilesList = mainOptions.primeFilesList[:0]
	for _, file := range mainOptions.orgFilenames {
		mainOptions.moveFilesList = append(mainOptions.moveFilesList,
			[]string{mainObjects.MoveFilechooserButton.GetFilename(), file[1], file[2], file[3]})

		/* = append(mainOptions.primeFilesList,
		filepath.Join(mainObjects.MoveFilechooserButton.GetFilename(), file[1]+file[2], file[3]))*/
	}
	fillListstore(mainObjects.moveListstore, mainOptions.moveFilesList, true)
}

// Handling button for apply moving files
func MoveApplyButtonClicked() {
	renameAndRefresh(mainOptions.orgFilenames, mainOptions.moveFilesList)
}

// Clear Files list and all slices content
func CumulativeDndChkClicked(chk *gtk.CheckButton) {
	mainOptions.CumulativeDND = chk.GetActive()
	mainObjects.RenCumulativeDndChk.SetActive(mainOptions.CumulativeDND)
	mainObjects.MoveCumulativeDndChk.SetActive(mainOptions.CumulativeDND)
	mainObjects.TitleCumulativeDndChk.SetActive(mainOptions.CumulativeDND)
}

// Handling title apply button
func TitleApplyButtonClicked() {
	var value int
	var textFiles, textTitles, alertMessage string
	lengthListTitle := len(mainOptions.titlList)
	lengthTitle := len(mainOptions.titlFilenames)
	lengthFiles := len(mainOptions.orgFilenames)
	if lengthTitle != 0 && lengthFiles != 0 {

		if lengthFiles != lengthListTitle {
			textFiles = "There is : " + strconv.Itoa(lengthFiles) + " file\n"
			textTitles = "And : " + strconv.Itoa(lengthListTitle) + "  title"
			if lengthFiles > 1 {
				textTitles = "And : " + strconv.Itoa(lengthListTitle) + "  titles"
			}
			if lengthListTitle > 1 {
				textFiles = "There are : " + strconv.Itoa(lengthFiles) + " files\n"
			}
			if lengthFiles > lengthListTitle {
				alertMessage = "\nTiteling alert.\nNot enought titles\n\n"
			} else {
				alertMessage = "\nTiteling alert.\nToo many titles\n\n"
			}
			value = gi.DlgMessage(mainObjects.MainWindow, "alert", "Attention !",
				alertMessage+textFiles+textTitles+"\n\nProceed anyway ?",
				"", "Cancel", "ok")
		} else {
			value = 1
		}
		if value == 1 {
			errList, err := renameMe(mainOptions.orgFilenames, mainOptions.titlFilenames, true)
			if err != nil {
				if err.Error() == errCancel {
					return
				}
				if !strings.Contains(err.Error(), errNbFileDoesNotMatch) {
					gi.DlgMessage(mainObjects.MainWindow, "alert", "Attention !",
						"\nTiteling error.\n"+err.Error()+"\n"+strings.Join(errList, "\n")+"\n", "", "ok")
				}
			}

			mainOptions.rawFilesList = mainOptions.rawFilesList[:0]
			mainOptions.primeFilesList = mainOptions.primeFilesList[:0]
			for _, file := range mainOptions.titlFilenames {
				mainOptions.rawFilesList = append(mainOptions.rawFilesList, filepath.Join(file[0], file[1]+file[2]))
				mainOptions.primeFilesList = append(mainOptions.primeFilesList, filepath.Join(file[0], file[1]+file[2]))
			}
			// Clean textview
			txtBuff, err := mainObjects.TitleTextview.GetBuffer()
			Check(err)
			txtBuff.SetText("")

			resetFilenamesStorage()

			// Cleaning entries
			mainObjects.TitleAddAEntry.SetText("")
			mainObjects.TitleAddBEntry.SetText("")
			mainObjects.TitleAddBFileEntry.SetText("")
			mainObjects.TitleSepEntry.SetText("")
			mainObjects.TitleSpinbutton.SetValue(float64(0))

			RenPresExtChkChanged() // Include refreshlist
		}
	}
}

// Record TitleTextview content when focus out
func TitleTextviewFocusOut() {
	var newList []string
	txtBuff, err := mainObjects.TitleTextview.GetBuffer()
	Check(err)
	txt, err := txtBuff.GetText(txtBuff.GetStartIter(), txtBuff.GetEndIter(), false)
	Check(err)
	mainOptions.titlList = strings.Split(txt, GetTextEOL([]byte(txt)))

	for _, line := range mainOptions.titlList {
		tmpStr, err := TrimSpace(line, "-c")
		Check(err)
		if len(tmpStr) != 0 {
			newList = append(newList, tmpStr)
		}
	}
	mainOptions.titlList = newList
}

// Handling TitleAddToFileEntryEvent
func TitleAddToFileEntryEvent() {
	textB, err := mainObjects.TitleAddBFileEntry.GetText()
	Check(err)
	textB, err = TrimSpace(textB, "+w")
	Check(err)
	mainObjects.TitleAddBFileEntry.SetText(textB)

	if len(textB) != 0 {
		if len(mainOptions.titlFilenames) == 0 {
			mainOptions.titlFilenames = make([][]string, len(mainOptions.orgFilenames))
			copy(mainOptions.titlFilenames, mainOptions.orgFilenames)
		}
		// mainOptions.titlList = mainOptions.titlList[:0]
		// for idx := 0; idx < len(mainOptions.titlFilenames); idx++ {
		// 	mainOptions.titlList = append(mainOptions.titlList, textB)
		// }
	}
	// Display modified titlefilelist
	mainOptions.titlFilenames, _ = makeTitleFilename(mainOptions.orgFilenames,
		mainOptions.titlList,
		textB)
	fillListstore(mainObjects.titleListstore, mainOptions.titlFilenames)
}

// Copy filelist content to textview
func TitlCopyListButtonClicked() {
	if len(mainOptions.orgFilenames) != 0 {
		var tmpText string
		txtBuff, err := mainObjects.TitleTextview.GetBuffer()
		Check(err)
		txtBuff.SetText("")
		for _, name := range mainOptions.orgFilenames {
			tmpText += name[1] + "\n"
		}
		txtBuff.SetText(tmpText)
	}
}

// Title format calculation
func TitleEntryFocusOut() {
	// To allways get newest entries
	TitleTextviewFocusOut()

	var tmpFinalList []string
	strSep, err := mainObjects.TitleSepEntry.GetText()
	Check(err)
	strAddA, err := mainObjects.TitleAddAEntry.GetText()
	Check(err)
	strAddB, err := mainObjects.TitleAddBEntry.GetText()
	Check(err)
	spinVal := mainObjects.TitleSpinbutton.GetValueAsInt()

	if len(strSep) != 0 {
		for _, line := range mainOptions.titlList {
			tmpStrSl := strings.Split(line, strSep)
			if spinVal < len(tmpStrSl) {
				// Sanitize entries and remove useless spaces
				tmpStr, err := TrimSpace(tmpStrSl[spinVal], "-c", "-w")
				Check(err)
				tmpFinalList = append(tmpFinalList, strAddB+tmpStr+strAddA)
			} else {
				fmt.Println("field out of range")
			}
		}
	} else {
		for _, line := range mainOptions.titlList {
			tmpStr, err := TrimSpace(line, "-c", "-w")
			Check(err)
			tmpFinalList = append(tmpFinalList, strAddB+tmpStr+strAddA)
		}
	}
	mainOptions.titlList = tmpFinalList
	// Display modified titlefilelist
	mainOptions.titlFilenames, _ = makeTitleFilename(mainOptions.orgFilenames, mainOptions.titlList)

	fillListstore(mainObjects.titleListstore, mainOptions.titlFilenames)
	// Update Add before filname if entry exist

	TitleAddToFileEntryEvent()
}

// Renamer entry focus out
func RenRemEntryFocusOut() {
	// Retreive Entry values
	getEntryDatas()
	if len(mainOptions.renFilenames) == 0 {
		mainOptions.renFilenames = make([][]string, len(mainOptions.orgFilenames))
		copy(mainOptions.renFilenames, mainOptions.orgFilenames)
	}
	mainOptions.renFilenamesBak = make([][]string, len(mainOptions.renFilenames))
	copy(mainOptions.renFilenamesBak, mainOptions.renFilenames)

	doRename(&mainOptions.renFilenamesBak)
}

// Handling rename KEEPBETWEEN button
func RenKeepBtwButtonClicked() {
	displayOverWin("Keep Between")
	mainObjects.OverEntry.GrabFocus()
	mainObjects.OverGrid1.SetVisible(true)
	mainObjects.OverKeepBeforeLbl.SetText("Keep Before")
	mainObjects.OverKeepAfterLbl.SetText("Keep After")
	mainObjects.OverKeepAfterLbl.SetSizeRequest(100, 24)
	mainObjects.OverKeepBeforeLbl.SetSizeRequest(100, 24)
	setImage(mainObjects.OverImageTop, keepBetween48)

	mainObjects.OverOkButton.Connect("clicked", func() {
		txt, err := mainObjects.OverEntry.GetText()
		Check(err)
		txt1, err := mainObjects.OverEntry1.GetText()
		Check(err)
		if len(txt) != 0 && len(txt1) != 0 {
			mainOptions.overText = mainOptions.overText[:0]
			mainOptions.overText = append(mainOptions.overText, txt)
			mainOptions.overText = append(mainOptions.overText, txt1)
			genericHideWindow(mainObjects.OverWindow)

			// Keep Between part
			var name = make([]string, 3)
			var keepBetween string
			mainOptions.renFilenames = make([][]string, len(mainOptions.orgFilenames))
			copy(mainOptions.renFilenames, mainOptions.orgFilenames)

			for idx := 0; idx < len(mainOptions.renFilenames); idx++ {
				name = mainOptions.renFilenames[idx]

				keepBetween = keepBetweenIt(name[1], mainOptions.overText,
					mainOptions.OvercaseSensitive,
					mainOptions.OverPosixCharClass,
					mainOptions.OverPosixStrictMde)

				mainOptions.renFilenames[idx] = []string{name[0], keepBetween, name[2]}
			}
			// Display Keep Between result
			RenRemEntryFocusOut()

		} else {
			genericHideWindow(mainObjects.OverWindow)
			fmt.Println("nothing done")
		}
	})
}

// Handling REGEX button
func RenRegexButtonClicked() {
	displayOverWin("Regular expression")
	mainObjects.OverEntry.GrabFocus()
	mainObjects.OverGrid7.SetVisible(false)
	mainObjects.OverKeepBeforeLbl.SetText("Replace with")
	mainObjects.OverKeepAfterLbl.SetText("Regexp")
	mainObjects.OverKeepAfterLbl.SetSizeRequest(100, 24)
	mainObjects.OverKeepBeforeLbl.SetSizeRequest(100, 24)
	setImage(mainObjects.OverImageTop, regex48x48)

	mainObjects.OverOkButton.Connect("clicked", func() {
		regexStr, err := mainObjects.OverEntry.GetText()
		Check(err)
		replaceWith, err := mainObjects.OverEntry1.GetText()
		Check(err)
		if len(regexStr) != 0 {
			genericHideWindow(mainObjects.OverWindow)

			// Regexp part
			var name = make([]string, 3)
			var substracted string
			mainOptions.renFilenames = make([][]string, len(mainOptions.orgFilenames))
			copy(mainOptions.renFilenames, mainOptions.orgFilenames)

			for idx := 0; idx < len(mainOptions.renFilenames); idx++ {
				name = mainOptions.renFilenames[idx]

				substracted = regexIt(name[1], regexStr, replaceWith)
				mainOptions.renFilenames[idx] = []string{name[0], substracted, name[2]}
			}
			// Display substracted result
			RenRemEntryFocusOut()

		} else {
			genericHideWindow(mainObjects.OverWindow)
			fmt.Println("nothing done")
		}
	})
}

// Handling rename SUBSTRACT button
func RenSubButtonClicked() {
	displayOverWin("Substract")
	mainObjects.OverEntry.GrabFocus()
	mainObjects.OverGrid1.SetVisible(false)
	mainObjects.OverKeepAfterLbl.SetText("To Substract")
	mainObjects.OverKeepAfterLbl.SetSizeRequest(100, 24)
	mainObjects.OverKeepBeforeLbl.SetSizeRequest(100, 24)
	setImage(mainObjects.OverImageTop, substract48)

	mainObjects.OverOkButton.Connect("clicked", func() {
		txt, err := mainObjects.OverEntry.GetText()
		Check(err)
		if len(txt) != 0 {

			mainOptions.overText = mainOptions.overText[:0]
			mainOptions.overText = append(mainOptions.overText, txt)
			genericHideWindow(mainObjects.OverWindow)

			// Substract part
			var name = make([]string, 3)
			var substracted string
			mainOptions.renFilenames = make([][]string, len(mainOptions.orgFilenames))
			copy(mainOptions.renFilenames, mainOptions.orgFilenames)

			for idx := 0; idx < len(mainOptions.renFilenames); idx++ {
				name = mainOptions.renFilenames[idx]

				substracted = substractIt(name[1], txt,
					mainOptions.OvercaseSensitive,
					mainOptions.OverPosixCharClass,
					mainOptions.OverPosixStrictMde)

				mainOptions.renFilenames[idx] = []string{name[0], substracted, name[2]}
			}
			// Display substracted result
			RenRemEntryFocusOut()

		} else {
			genericHideWindow(mainObjects.OverWindow)
			fmt.Println("nothing done")
		}
	})
}

/**************************/
/* Handling SINGLE window */
/**************************/
func SingleEntryEnterKeyPressed() {
	SingleOkButtonClicked()
}

func SingleSwMultiButtonClicked() {
	filename := filepath.Join(singleEntry[0], singleEntry[1]) + singleEntry[2]
	fi, err := os.Stat(filename)
	if err != nil {
		gi.DlgMessage(mainObjects.MainWindow, "error", sts["mstk"],
			"\n"+sts["filenameErr"],
			"", sts["ok"])
		return
	}
	if fi.IsDir() {
		if err = FindDir(filename, mainOptions.ExtMaskRen, &mainOptions.primeFilesList,
			mainOptions.ScanSubDir, mainOptions.ShowDirectory, true); err != nil {
			gi.DlgMessage(mainObjects.MainWindow, "error", sts["mstk"], sts["\nerrDir\n"]+err.Error(), "", sts["ok"])
		}
	} else {
		mainOptions.primeFilesList = append(mainOptions.primeFilesList, filename)
		restorePrimeFilesList()
	}
	resetFilenamesStorage()
	refreshLists()

	mainOptions.UpdateObjects()
	mainObjects.SingleWindow.Hide()
	mainObjects.MainWindow.ShowAll()
}

func SingleOkButtonClicked() {
	var err error
	var newName, oldName string
	text, _ := mainObjects.SingleEntry.GetText()

	if !mainObjects.SinglePresExtChk.GetActive() {
		newName = filepath.Join(singleEntry[0], text)
	} else {
		newName = filepath.Join(singleEntry[0], text+singleEntry[2])
	}
	oldName = filepath.Join(singleEntry[0], singleEntry[1]+singleEntry[2])

	if oldName != newName {
		if _, err := os.Stat(newName); !os.IsNotExist(err) {
			err = errors.New(newName + ": file exists\n\n" + oldName + "\n")
		} else {
			err = os.Rename(oldName, newName)
		}
	}
	if err != nil {
		gi.DlgMessage(mainObjects.MainWindow, "error", "Attention !", "\nRenaming error: \n"+err.Error(), "", "ok")
	} else {
		windowDestroy()
	}
}

func SingleResetButtonClicked() {
	if mainObjects.SinglePresExtChk.GetActive() {
		mainObjects.SingleEntry.SetText(singleEntry[1])
	} else {
		mainObjects.SingleEntry.SetText(singleEntry[1] + singleEntry[2])
	}
}

func SinglePresExtChkClicked() {
	mainOptions.PreserveExtSingle = mainObjects.SinglePresExtChk.GetActive()
	text, _ := mainObjects.SingleEntry.GetText()
	if !mainObjects.SinglePresExtChk.GetActive() {
		mainObjects.SingleEntry.SetText(text + singleEntry[2])
	} else {
		mainObjects.SingleEntry.SetText(BaseNoExt(filepath.Base(text)))
	}
}

func SingleEntryChanged() {
	text, _ := mainObjects.SingleEntry.GetText()
	text, _ = TrimSpace(text, "+w")
	mainObjects.SingleEntry.SetText(text)
}

/************************/
/* Handling OVER window */
/************************/
func displayOverWin(title string) {
	mainObjects.OverWindow.SetTitle(title)
	// mainObjects.OverImageTop.SetFromFile(imgName.(string))
	mainObjects.OverWindow.SetSkipTaskbarHint(true)
	mainObjects.OverWindow.SetKeepAbove(true)
	mainObjects.OverWindow.SetSizeRequest(400, 10)
	mainObjects.OverWindow.SetResizable(false)
	mainObjects.OverWindow.SetModal(true)
	mainObjects.OverWindow.Connect("delete_event", genericHideWindow)
	mainObjects.OverWindow.ShowAll()
	mainObjects.OverResetButton.Connect("clicked", func() {
		mainObjects.OverEntry.SetText("")
		mainObjects.OverEntry1.SetText("")
		mainOptions.overText = mainOptions.overText[:0]
		resetFilenamesStorage()
		refreshLists()
		RenRemEntryFocusOut()
		genericHideWindow(mainObjects.OverWindow)
	})
}

// Signal handler delete_event (hidding window)
func genericHideWindow(w *gtk.Window) bool {
	if w.GetVisible() {
		w.Hide()
	}
	return true
}

/*****************************************************/
/* imgTop handler release signal (Display about box) */
/*****************************************************/
func imgTopReleaseEvent() {
	mainOptions.AboutOptions.Show()
}
