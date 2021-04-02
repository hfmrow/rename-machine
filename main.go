// main.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 14:58:19 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 hfmrow - Rename Machine v1.6.1 github.com/hfmrow/rename-machine
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

func main() {

	/* Build options */

	// devMode: is used in some functions to control the behavior of the program
	// When software is ready to be published, this flag must be set at "false"
	// that means:
	// - options file will be stored in $HOME/.config/[Creat]/[softwareName],
	// - translate function if used, will no more auto-update "sts" map sentences,
	// - all built-in assets will be used instead of the files themselves.
	//   Be aware to update assets via "Goh" and translations with "Got" before all.
	devMode = true

	// Create temp directory .. or not
	doTempDir = false

	// Assets initialization according to the chosen mode (devMode).
	assetsDeclarationsUseEmbedded(!devMode)

	absoluteRealPath, optFilename = getAbsRealPath()

	/* Init & read options file */
	mainOptions = new(MainOpt) // Assignate options' structure.
	mainOptions.Read()

	/* Logger init. */
	Logger = Log2FileStructNew(optFilename, devMode)
	defer Logger.CloseLogger()

	/* Init gtk display */
	mainStartGtk(fmt.Sprintf("%s %s  %s %s %s",
		mainOptions.About.AppName,
		mainOptions.About.AppVers,
		mainOptions.About.YearCreat,
		mainOptions.About.AppCreats,
		mainOptions.About.LicenseAbrv),
		mainOptions.MainWinWidth,
		mainOptions.MainWinHeight, true,
		// If there is only 2 argument we start single file gui
		len(os.Args) != 2)
}

func mainApplication() {
	/*	Fill app infos	*/
	mainOptions.About.InitFillInfos(
		mainObjects.MainWindow,
		"About "+Name,
		Name,
		Vers,
		YearCreat,
		Creat,
		LicenseAbrv,
		LicenseShort,
		Repository,
		Descr,
		"",
		checked18x18)

	/* Translate init. */
	translate = MainTranslateNew(absoluteRealPath+mainOptions.LanguageFilename, devMode)

	/* Command line arguments handling */
	if len(os.Args) == 2 {
		singleEntry = []string{filepath.Dir(os.Args[1]), BaseNoExt(filepath.Base(os.Args[1])), filepath.Ext(os.Args[1])}
		if fileInfo, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
			Check(err)
			os.Exit(1)
		} else if fileInfo.IsDir() {
			// backup previous state of preserve extension in case of folder name modification
			bakPresExt = mainOptions.PreserveExtSingle
			mainOptions.PreserveExtSingle = false
			mainObjects.SinglePresExtChk.SetVisible(false)
			mainObjects.SingleEntry.SetText(singleEntry[1] + singleEntry[2])
			mainOptions.baseDirectory = os.Args[1]
		} else {
			mainObjects.SinglePresExtChk.SetActive(mainOptions.PreserveExtSingle)
			if mainOptions.PreserveExtSingle {
				mainObjects.SingleEntry.SetText(singleEntry[1])
			} else {
				mainObjects.SingleEntry.SetText(singleEntry[1] + singleEntry[2])
			}
		}
		// Select whole entry content (filename)
		mainObjects.SingleEntry.SelectRegion(0, int(mainObjects.SingleEntry.GetTextLength()))
		mainObjects.SingleEntry.GrabFocus()
	} else {
		// Build files list
		for idx := 1; idx < len(os.Args); idx++ {
			mainOptions.primeFilesList = append(mainOptions.primeFilesList, os.Args[idx])
		}
	}

	restorePrimeFilesList()

	if err = initSomeControls(); err == nil {

		/* display files */
		refreshLists()

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

		mainOptions.UpdateObjects(true)
	}

	gidg.DialogError(mainObjects.MainWindow, sts["mstk"], sts["issueWith"], err, devMode, true)

}

func initSomeControls() (err error) {

	// Drag & drop Init.
	mainOptions.treeViewDropSet = gimc.DragNDropNew(mainObjects.FileListTreeview, &mainOptions.dndFilesList, receiveDnd)

	// Initialiste liststore Columns
	if tvsFiles, err = gitw.TreeViewStructureNew(mainObjects.FileListTreeview, false, false); err == nil {
		tvsFiles.AddColumns(oriFileListstoreCol, true, true, false, true, true, true)
		if err = tvsFiles.StoreSetup(new(gtk.ListStore)); err == nil {
			// For title function
			if tvsTitle, err = gitw.TreeViewStructureNew(mainObjects.TitlePrevTreeview, false, false); err == nil {
				tvsTitle.AddColumns(modFileListstoreCol, false, true, false, true, true, true)
				if err = tvsTitle.StoreSetup(new(gtk.ListStore)); err == nil {
					// For rename function
					if tvsRenam, err = gitw.TreeViewStructureNew(mainObjects.RenPrevTreeview, false, false); err == nil {
						tvsRenam.AddColumns(modFileListstoreCol, false, true, false, true, true, true)
						if err = tvsRenam.StoreSetup(new(gtk.ListStore)); err == nil {
							// For move file function
							if tvsMoves, err = gitw.TreeViewStructureNew(mainObjects.MovePrevTreeview, false, false); err == nil {
								tvsMoves.AddColumns(oriFileListstoreCol, false, true, false, true, true, true)
								if err = tvsMoves.StoreSetup(new(gtk.ListStore)); err == nil {
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

/*************************************\
/* Executed just before closing all. */
/************************************/
func onShutdown() bool {
	var err error
	// Update mainOptions with GtkObjects and save it
	if err = mainOptions.Write(); err == nil {
		// What you want to execute before closing the app.
		// Return:
		// true for exit applicaton
		// false does not exit application
	}
	if err != nil {
		log.Fatalf("Unexpected error on exit: %s", err.Error())
	}
	return true
}
