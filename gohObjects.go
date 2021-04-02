// gohObjects.go

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
	"github.com/gotk3/gotk3/gtk"
)

// Control over all used objects from glade.
var mainObjects *MainControlsObj

/******************************/
/* Main structure Declaration */
/******************************/
type MainControlsObj struct {
	FileListTreeview           *gtk.TreeView
	mainUiBuilder              *gtk.Builder
	MainWindow                 *gtk.Window
	MoveApplyButton            *gtk.Button
	MoveCumulativeDndChk       *gtk.CheckButton
	MoveEntryExtMask           *gtk.Entry
	MoveFilechooserButton      *gtk.FileChooserButton
	MoveLabelExtMask           *gtk.Label
	MovePrevTreeview           *gtk.TreeView
	Notebook                   *gtk.Notebook
	OverCaseSensChk            *gtk.CheckButton
	OverCharClassChk           *gtk.CheckButton
	OverCharClassStrictModeChk *gtk.CheckButton
	OverEntry                  *gtk.Entry
	OverEntry1                 *gtk.Entry
	OverGrid                   *gtk.Grid
	OverGrid1                  *gtk.Grid
	OverGrid7                  *gtk.Grid
	OverImageTop               *gtk.Image
	OverKeepAfterLbl           *gtk.Label
	OverKeepBeforeLbl          *gtk.Label
	OverOkButton               *gtk.Button
	OverResetButton            *gtk.Button
	OverWindow                 *gtk.Window
	RenApplyButton             *gtk.Button
	RenCaseSensChk             *gtk.CheckButton
	RenCumulativeDndChk        *gtk.CheckButton
	RenEntryExtMask            *gtk.Entry
	RenIncrementChk            *gtk.CheckButton
	RenIncrementRightChk       *gtk.CheckButton
	RenIncSepEntry             *gtk.Entry
	RenIncSpinbutton           *gtk.SpinButton
	RenKeepBtwButton           *gtk.Button
	RenLabelExtMask            *gtk.Label
	RenPresExtChk              *gtk.CheckButton
	RenPrevTreeview            *gtk.TreeView
	RenRegexButton             *gtk.Button
	RenRemEntry                *gtk.Entry
	RenRemEntry1               *gtk.Entry
	RenRemEntry2               *gtk.Entry
	RenReplEntry               *gtk.Entry
	RenReplEntry1              *gtk.Entry
	RenReplEntry2              *gtk.Entry
	RenScanSubDirChk           *gtk.CheckButton
	RenShowDirChk              *gtk.CheckButton
	RenSubButton               *gtk.Button
	RenWthEntry                *gtk.Entry
	RenWthEntry1               *gtk.Entry
	RenWthEntry2               *gtk.Entry
	SingleCancelButton         *gtk.Button
	SingleEntry                *gtk.Entry
	SingleImageTop             *gtk.Image
	SingleOkButton             *gtk.Button
	SinglePresExtChk           *gtk.CheckButton
	SingleResetButton          *gtk.Button
	SingleSwMultiButton        *gtk.Button
	SingleWindow               *gtk.Window
	Statusbar                  *gtk.Statusbar
	TitleAddAEntry             *gtk.Entry
	TitleAddBEntry             *gtk.Entry
	TitleAddBFileEntry         *gtk.Entry
	TitleApplyButton           *gtk.Button
	TitleCumulativeDndChk      *gtk.CheckButton
	TitlePrevTreeview          *gtk.TreeView
	TitleScanSubDirChk         *gtk.CheckButton
	TitleSepEntry              *gtk.Entry
	TitleSpinbutton            *gtk.SpinButton
	TitleTextview              *gtk.TextView
	TopImage                   *gtk.Image
	TopImageEventbox           *gtk.EventBox
}

/******************************/
/* GtkObjects  Initialisation */
/******************************/
// gladeObjParser: Initialise Gtk3 Objects into mainObjects structure.
func gladeObjParser() {
	mainObjects.FileListTreeview = loadObject("FileListTreeview").(*gtk.TreeView)
	mainObjects.MainWindow = loadObject("MainWindow").(*gtk.Window)
	mainObjects.MoveApplyButton = loadObject("MoveApplyButton").(*gtk.Button)
	mainObjects.MoveCumulativeDndChk = loadObject("MoveCumulativeDndChk").(*gtk.CheckButton)
	mainObjects.MoveEntryExtMask = loadObject("MoveEntryExtMask").(*gtk.Entry)
	mainObjects.MoveFilechooserButton = loadObject("MoveFilechooserButton").(*gtk.FileChooserButton)
	mainObjects.MoveLabelExtMask = loadObject("MoveLabelExtMask").(*gtk.Label)
	mainObjects.MovePrevTreeview = loadObject("MovePrevTreeview").(*gtk.TreeView)
	mainObjects.Notebook = loadObject("Notebook").(*gtk.Notebook)
	mainObjects.OverCaseSensChk = loadObject("OverCaseSensChk").(*gtk.CheckButton)
	mainObjects.OverCharClassChk = loadObject("OverCharClassChk").(*gtk.CheckButton)
	mainObjects.OverCharClassStrictModeChk = loadObject("OverCharClassStrictModeChk").(*gtk.CheckButton)
	mainObjects.OverEntry = loadObject("OverEntry").(*gtk.Entry)
	mainObjects.OverEntry1 = loadObject("OverEntry1").(*gtk.Entry)
	mainObjects.OverGrid = loadObject("OverGrid").(*gtk.Grid)
	mainObjects.OverGrid1 = loadObject("OverGrid1").(*gtk.Grid)
	mainObjects.OverGrid7 = loadObject("OverGrid7").(*gtk.Grid)
	mainObjects.OverImageTop = loadObject("OverImageTop").(*gtk.Image)
	mainObjects.OverKeepAfterLbl = loadObject("OverKeepAfterLbl").(*gtk.Label)
	mainObjects.OverKeepBeforeLbl = loadObject("OverKeepBeforeLbl").(*gtk.Label)
	mainObjects.OverOkButton = loadObject("OverOkButton").(*gtk.Button)
	mainObjects.OverResetButton = loadObject("OverResetButton").(*gtk.Button)
	mainObjects.OverWindow = loadObject("OverWindow").(*gtk.Window)
	mainObjects.RenApplyButton = loadObject("RenApplyButton").(*gtk.Button)
	mainObjects.RenCaseSensChk = loadObject("RenCaseSensChk").(*gtk.CheckButton)
	mainObjects.RenCumulativeDndChk = loadObject("RenCumulativeDndChk").(*gtk.CheckButton)
	mainObjects.RenEntryExtMask = loadObject("RenEntryExtMask").(*gtk.Entry)
	mainObjects.RenIncrementChk = loadObject("RenIncrementChk").(*gtk.CheckButton)
	mainObjects.RenIncrementRightChk = loadObject("RenIncrementRightChk").(*gtk.CheckButton)
	mainObjects.RenIncSepEntry = loadObject("RenIncSepEntry").(*gtk.Entry)
	mainObjects.RenIncSpinbutton = loadObject("RenIncSpinbutton").(*gtk.SpinButton)
	mainObjects.RenKeepBtwButton = loadObject("RenKeepBtwButton").(*gtk.Button)
	mainObjects.RenLabelExtMask = loadObject("RenLabelExtMask").(*gtk.Label)
	mainObjects.RenPresExtChk = loadObject("RenPresExtChk").(*gtk.CheckButton)
	mainObjects.RenPrevTreeview = loadObject("RenPrevTreeview").(*gtk.TreeView)
	mainObjects.RenRegexButton = loadObject("RenRegexButton").(*gtk.Button)
	mainObjects.RenRemEntry = loadObject("RenRemEntry").(*gtk.Entry)
	mainObjects.RenRemEntry1 = loadObject("RenRemEntry1").(*gtk.Entry)
	mainObjects.RenRemEntry2 = loadObject("RenRemEntry2").(*gtk.Entry)
	mainObjects.RenReplEntry = loadObject("RenReplEntry").(*gtk.Entry)
	mainObjects.RenReplEntry1 = loadObject("RenReplEntry1").(*gtk.Entry)
	mainObjects.RenReplEntry2 = loadObject("RenReplEntry2").(*gtk.Entry)
	mainObjects.RenScanSubDirChk = loadObject("RenScanSubDirChk").(*gtk.CheckButton)
	mainObjects.RenShowDirChk = loadObject("RenShowDirChk").(*gtk.CheckButton)
	mainObjects.RenSubButton = loadObject("RenSubButton").(*gtk.Button)
	mainObjects.RenWthEntry = loadObject("RenWthEntry").(*gtk.Entry)
	mainObjects.RenWthEntry1 = loadObject("RenWthEntry1").(*gtk.Entry)
	mainObjects.RenWthEntry2 = loadObject("RenWthEntry2").(*gtk.Entry)
	mainObjects.SingleCancelButton = loadObject("SingleCancelButton").(*gtk.Button)
	mainObjects.SingleEntry = loadObject("SingleEntry").(*gtk.Entry)
	mainObjects.SingleImageTop = loadObject("SingleImageTop").(*gtk.Image)
	mainObjects.SingleOkButton = loadObject("SingleOkButton").(*gtk.Button)
	mainObjects.SinglePresExtChk = loadObject("SinglePresExtChk").(*gtk.CheckButton)
	mainObjects.SingleResetButton = loadObject("SingleResetButton").(*gtk.Button)
	mainObjects.SingleSwMultiButton = loadObject("SingleSwMultiButton").(*gtk.Button)
	mainObjects.SingleWindow = loadObject("SingleWindow").(*gtk.Window)
	mainObjects.Statusbar = loadObject("Statusbar").(*gtk.Statusbar)
	mainObjects.TitleAddAEntry = loadObject("TitleAddAEntry").(*gtk.Entry)
	mainObjects.TitleAddBEntry = loadObject("TitleAddBEntry").(*gtk.Entry)
	mainObjects.TitleAddBFileEntry = loadObject("TitleAddBFileEntry").(*gtk.Entry)
	mainObjects.TitleApplyButton = loadObject("TitleApplyButton").(*gtk.Button)
	mainObjects.TitleCumulativeDndChk = loadObject("TitleCumulativeDndChk").(*gtk.CheckButton)
	mainObjects.TitlePrevTreeview = loadObject("TitlePrevTreeview").(*gtk.TreeView)
	mainObjects.TitleScanSubDirChk = loadObject("TitleScanSubDirChk").(*gtk.CheckButton)
	mainObjects.TitleSepEntry = loadObject("TitleSepEntry").(*gtk.Entry)
	mainObjects.TitleSpinbutton = loadObject("TitleSpinbutton").(*gtk.SpinButton)
	mainObjects.TitleTextview = loadObject("TitleTextview").(*gtk.TextView)
	mainObjects.TopImage = loadObject("TopImage").(*gtk.Image)
	mainObjects.TopImageEventbox = loadObject("TopImageEventbox").(*gtk.EventBox)
}
