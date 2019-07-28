// gohStartGtk.go

// Source file auto-generated on Sun, 28 Jul 2019 07:02:22 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
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
		// Init tempDir and Remove tempDirectory on exit
		tempDir = tempMake(Name)
		defer os.RemoveAll(tempDir)
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
			if len(os.Args) == 2 {
				mainObjects.SingleWindow.SetTitle(winTitle)
				mainObjects.SingleWindow.Connect("delete-event", windowDestroy)
				mainObjects.SingleWindow.SetKeepAbove(true)
				mainObjects.SingleWindow.SetSizeRequest(400, 10)
				mainObjects.SingleWindow.SetResizable(false)
				mainObjects.SingleWindow.SetModal(true)
				mainObjects.SingleWindow.ShowAll()
			}
		}
		// Start main application ...
		mainApplication()
		//	Start Gui loop
		gtk.Main()
	} else {
		log.Fatal("Builder initialisation error.")
	}
}
