// gohStartGtk.go

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
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

/******************************/
/* Gtk3 Window Initialisation */
/******************************/
func mainStartGtk(winTitle string, width, height int, center, show bool) {
	mainObjects = new(MainControlsObj)
	gtk.Init(nil)
	if newBuilder(mainGlade) == nil {
		// Init tempDir and Remove it on quit if requested.
		if doTempDir {
			tempDir = tempMake(Name)
			defer os.RemoveAll(tempDir)
		}
		// Parse Gtk objects
		gladeObjParser()
		// Objects Signals initialisations
		signalsPropHandler()
		// Init objects images
		assignImages()
		// Set Window Properties
		if center {
			mainObjects.MainWindow.SetPosition(gtk.WIN_POS_CENTER)
		}
		mainObjects.MainWindow.SetTitle(winTitle)
		mainObjects.MainWindow.SetDefaultSize(width, height)
		mainObjects.MainWindow.Connect("delete-event", windowDestroy)
		if show {
			mainObjects.MainWindow.ShowAll()
		} else {
			mainObjects.SingleWindow.SetTitle(winTitle)
			mainObjects.SingleWindow.Connect("delete-event", windowDestroy)
			mainObjects.SingleWindow.SetKeepAbove(true)
			mainObjects.SingleWindow.SetSizeRequest(400, 10)
			mainObjects.SingleWindow.SetResizable(false)
			mainObjects.SingleWindow.SetModal(true)
			mainObjects.SingleWindow.ShowAll()
		}
		// Start main application ...
		mainApplication()
		//	Start Gui loop
		gtk.Main()
	} else {
		log.Fatal("Builder initialisation error.")
	}
}

// windowDestroy: on closing/destroying the gui window.
func windowDestroy() {
	if onShutdown() {
		gtk.MainQuit()
	}
}
