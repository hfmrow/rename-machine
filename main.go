// main.go

// Source file auto-generated on Thu, 25 Jul 2019 18:41:58 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
	RenameMachine v1.4.5 ©2018-19 H.F.M

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:

	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
	associated documentation files (the "Software"), to dealin the Software without restriction,
	including without limitation the rights to use, copy, modify, merge, publish, distribute,
	sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all copies or
	substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
	NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
	NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
	DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT
	OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
	g "github.com/hfmrow/renMachine/genLib"
	gi "github.com/hfmrow/renMachine/gtk3Import"
)

func main() {
	/* Be or not to be ... in dev mode ... */
	devMode = true

	/* Set to true when you choose using embedded assets functionality */
	assetsDeclarationsUseEmbedded(!devMode)

	var multi bool

	/* Init Options */
	mainOptions = new(MainOpt)
	mainOptions.Init()

	/* Read Options */
	err = mainOptions.Read()
	if err != nil {
		fmt.Printf("%s\n%v\n", "Reading options error.", err)
	}

	/*	Fill app infos	*/
	mainOptions.AboutOptions.InitFillInfos(
		"About "+Name,
		Name,
		Vers,
		YearCreat,
		Creat,
		LicenseAbrv,
		LicenseShort,
		Repository,
		Descr,
		renameMachine400x27,
		checked18x18)

	if len(os.Args) > 2 {
		multi = true
	}

	/* Init gtk display */
	mainStartGtk(fmt.Sprintf("%s %s  %s %s %s",
		mainOptions.AboutOptions.AppName,
		mainOptions.AboutOptions.AppVers,
		mainOptions.AboutOptions.YearCreat,
		mainOptions.AboutOptions.AppCreats,
		mainOptions.AboutOptions.LicenseAbrv),
		mainOptions.MainWinWidth,
		mainOptions.MainWinHeight, true, multi)
}

func mainApplication() {

	/* Translate init. */
	translate = MainTranslateNew(absoluteRealPath+mainOptions.LanguageFilename, devMode)

	// Init files list
	for idx := 1; idx < len(os.Args); idx++ {
		mainOptions.primeFilesList = append(mainOptions.primeFilesList, os.Args[idx])
		if len(mainOptions.baseDirectory) == 0 {
			name := mainOptions.primeFilesList[0]
			fileInfos := g.SplitFilepath(name)
			if fileInfos.IsDir {
				if fileInfos.SymLink {
					mainOptions.baseDirectory = fileInfos.SymLinkTo
				} else {
					mainOptions.baseDirectory = name
				}
			} else {
				mainOptions.baseDirectory = fileInfos.Path
			}
		}
	}
	dnd = true
	switch len(os.Args) {
	case 1:
		gi.DlgMessage(mainObjects.MainWindow, "alert", "Files selection", "\nPlease select some files to proceed ... ", "", "ok")
		fmt.Println("\nNo file to proceed !")
		os.Exit(0)
	case 2:
		dnd = false
		splited := g.SplitFilepath(mainOptions.primeFilesList[0])
		singleEntry = []string{splited.Path, splited.BaseNoExt, splited.Ext}
		if splited.IsDir {
			// backup previous state of preserve extension in case of modification of folder name
			bakPresExt = mainOptions.PreserveExtSingle
			mainOptions.PreserveExtSingle = false
			mainObjects.SinglePresExtChk.SetVisible(false)
			mainObjects.SingleEntry.SetText(singleEntry[1] + singleEntry[2])
		} else {
			mainObjects.SinglePresExtChk.SetActive(mainOptions.PreserveExtSingle)
			if mainOptions.PreserveExtSingle {
				mainObjects.SingleEntry.SetText(singleEntry[1])
			} else {
				mainObjects.SingleEntry.SetText(singleEntry[1] + singleEntry[2])
			}
		}
		// Select whole entry content
		mainObjects.SingleEntry.SelectRegion(0, int(mainObjects.SingleEntry.GetTextLength()))
		mainObjects.SingleEntry.GrabFocus()
	}

	restorePrimeFilesList()
	mainOptions.ScanSubDir = false
	mainOptions.ShowDirectory = false

	initTreeview()

	/* Init Spinbutton */
	if ad, err := gtk.AdjustmentNew(0, 0, 65534, 1, 0, 0); err == nil {
		mainObjects.RenIncSpinbutton.Configure(ad, 1, 0)
	} else {
		fmt.Println("Error on:", "RenIncSpinbutton", "Initialisation")
	}
	if ad, err := gtk.AdjustmentNew(0, 0, 30, 1, 0, 0); err == nil {
		mainObjects.TitleSpinbutton.Configure(ad, 1, 0)
	} else {
		fmt.Println("Error on:", "TitleSpinbutton", "Initialisation")
	}

	mainOptions.UpdateObjects()
}

func initTreeview() {

	// Drag & drop Init.
	mainOptions.treeViewDropSet.InitDropSet(mainObjects.FileListTreeview, mainOptions.dndFilesList, receiveDnd)

	// Initialiste liststore Columns
	var column *gtk.TreeViewColumn
	mainObjects.fileListstore = gi.TreeViewListStoreSetup(mainObjects.FileListTreeview, false, oriFileListstoreCol, true)
	// column = mainObjects.FileListTreeview.GetColumn(0)
	// column.SetSortColumnID(0)
	// column.SetSizing(gtk.TREE_VIEW_COLUMN_AUTOSIZE)
	// column = mainObjects.FileListTreeview.GetColumn(1)
	// column.SetSizing(gtk.TREE_VIEW_COLUMN_AUTOSIZE)

	// For title function
	mainObjects.titleListstore = gi.TreeViewListStoreSetup(mainObjects.TitlePrevTreeview, false, modFileListstoreCol, false)
	column = mainObjects.TitlePrevTreeview.GetColumn(0)
	column.SetSortColumnID(-1)

	// For rename function
	mainObjects.renListstore = gi.TreeViewListStoreSetup(mainObjects.RenPrevTreeview, false, modFileListstoreCol, false)
	column = mainObjects.RenPrevTreeview.GetColumn(0)
	column.SetSortColumnID(-1)

	// For move file function
	mainObjects.moveListstore = gi.TreeViewListStoreSetup(mainObjects.MovePrevTreeview, false, oriFileListstoreCol, false)
	column = mainObjects.MovePrevTreeview.GetColumn(0)
	column.SetSortColumnID(-1)

	refreshLists()
}
