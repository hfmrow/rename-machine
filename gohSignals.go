// gohSignals.go

// Source file auto-generated on Thu, 25 Jul 2019 18:41:58 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/***************************/
/* Signals and Properties  */
/*    Implementations      */
/***************************/
func signalsPropHandler() {
	mainObjects.MoveApplyButton.Connect("clicked", MoveApplyButtonClicked)
	mainObjects.MoveEntryExtMask.Connect("activate", MoveEntryExtMaskEnterKeyPressed)
	mainObjects.MoveEntryExtMask.Connect("focus-out-event", MoveEntryExtMaskEnterKeyPressed)
	mainObjects.MoveFilechooserButton.Connect("selection-changed", MoveFilechooserButtonClicked)
	mainObjects.MoveLabelExtMask.Connect("notify", blankNotify)
	mainObjects.Notebook.Connect("switch-page", NotebookPageChanged)
	mainObjects.OverCaseSensChk.Connect("clicked", OverCaseSensChkChanged)
	mainObjects.OverCharClassChk.Connect("clicked", OverCharClassChkChanged)
	mainObjects.OverCharClassStrictModeChk.Connect("clicked", OverCharClassStrictModeChkChanged)
	mainObjects.RenApplyButton.Connect("clicked", RenApplyButtonClicked)
	mainObjects.RenCaseSensChk.Connect("clicked", RenCaseSensChkChanged)
	mainObjects.RenEntryExtMask.Connect("activate", RenEntryExtMaskEnterKeyPressed)
	mainObjects.RenEntryExtMask.Connect("focus-out-event", RenEntryExtMaskEnterKeyPressed)
	mainObjects.RenIncrementChk.Connect("clicked", RenIncrementChkClicked)
	mainObjects.RenIncrementRightChk.Connect("clicked", RenIncrementRightChkClicked)
	mainObjects.RenIncSepEntry.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenIncSpinbutton.Connect("value-changed", RenRemEntryFocusOut)
	mainObjects.RenKeepBtwButton.Connect("clicked", RenKeepBtwButtonClicked)
	mainObjects.RenLabelExtMask.Connect("notify", blankNotify)
	mainObjects.RenOrgSelButton.Connect("clicked", OrgSelButtonClicked)
	mainObjects.RenPresExtChk.Connect("clicked", RenPresExtChkChanged)
	mainObjects.RenRegexButton.Connect("clicked", RenRegexButtonClicked)
	mainObjects.RenRemEntry.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenRemEntry1.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenRemEntry2.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenReplEntry.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenReplEntry1.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenReplEntry2.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenScanSubDirChk.Connect("clicked", ScanSubDirChkChanged)
	mainObjects.RenShowDirChk.Connect("clicked", RenShowDirChkChanged)
	mainObjects.RenSubButton.Connect("clicked", RenSubButtonClicked)
	mainObjects.RenWthEntry.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenWthEntry1.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenWthEntry2.Connect("changed", RenRemEntryFocusOut)
	mainObjects.SingleCancelButton.Connect("clicked", windowDestroy)
	mainObjects.SingleEntry.Connect("changed", SingleEntryChanged)
	mainObjects.SingleEntry.Connect("activate", SingleEntryEnterKeyPressed)
	mainObjects.SingleOkButton.Connect("clicked", SingleOkButtonClicked)
	mainObjects.SinglePresExtChk.Connect("clicked", SinglePresExtChkClicked)
	mainObjects.SingleResetButton.Connect("clicked", SingleResetButtonClicked)
	mainObjects.SingleSwMultiButton.Connect("clicked", SingleSwMultiButtonClicked)
	mainObjects.TitlCopyListButton.Connect("clicked", TitlCopyListButtonClicked)
	mainObjects.TitleAddAEntry.Connect("changed", TitleEntryFocusOut)
	mainObjects.TitleAddBEntry.Connect("changed", TitleEntryFocusOut)
	// mainObjects.TitleAddBFileEntry.Connect("focus-out-event", TitleAddToFileEntryEvent)
	mainObjects.TitleAddBFileEntry.Connect("changed", TitleAddToFileEntryEvent)
	mainObjects.TitleApplyButton.Connect("clicked", TitleApplyButtonClicked)
	mainObjects.TitleScanSubDirChk.Connect("clicked", ScanSubDirChkChanged)
	mainObjects.TitleSepEntry.Connect("changed", TitleEntryFocusOut)
	mainObjects.TitleSpinbutton.Connect("value-changed", TitleEntryFocusOut)
	mainObjects.TitleTextview.Connect("event", TitleEntryFocusOut)
	mainObjects.TitlOrgSelButton.Connect("clicked", OrgSelButtonClicked)
	mainObjects.TopImageEventbox.Connect("button-release-event", imgTopReleaseEvent)
}
