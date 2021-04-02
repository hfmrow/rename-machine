// gohSignals.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 14:58:19 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 hfmrow - Rename Machine v1.6.1 github.com/hfmrow/rename-machine
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
	mainObjects.MoveCumulativeDndChk.Connect("clicked", CumulativeDndChkClicked)
	mainObjects.MoveEntryExtMask.Connect("activate", MoveEntryExtMaskEnterKeyPressed)
	mainObjects.MoveEntryExtMask.Connect("focus-out-event", MoveEntryExtMaskEnterKeyPressed)
	mainObjects.MoveFilechooserButton.Connect("selection-changed", MoveFilechooserButtonClicked)
	mainObjects.Notebook.Connect("switch-page", NotebookPageChanged)
	mainObjects.OverCaseSensChk.Connect("clicked", OverCaseSensChkChanged)
	mainObjects.OverCharClassChk.Connect("clicked", OverCharClassChkChanged)
	mainObjects.OverCharClassStrictModeChk.Connect("clicked", OverCharClassStrictModeChkChanged)
	mainObjects.RenApplyButton.Connect("clicked", RenApplyButtonClicked)
	mainObjects.RenCaseSensChk.Connect("clicked", RenCaseSensChkChanged)
	mainObjects.RenCumulativeDndChk.Connect("clicked", CumulativeDndChkClicked)
	mainObjects.RenEntryExtMask.Connect("activate", RenEntryExtMaskEnterKeyPressed)
	mainObjects.RenEntryExtMask.Connect("focus-out-event", RenEntryExtMaskEnterKeyPressed)
	mainObjects.RenIncrementChk.Connect("clicked", RenIncrementChkClicked)
	mainObjects.RenIncrementRightChk.Connect("clicked", RenIncrementRightChkClicked)
	mainObjects.RenIncSepEntry.Connect("changed", RenRemEntryFocusOut)
	mainObjects.RenIncSpinbutton.Connect("value-changed", RenRemEntryFocusOut)
	mainObjects.RenKeepBtwButton.Connect("clicked", RenKeepBtwButtonClicked)
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
	mainObjects.TitleAddAEntry.Connect("changed", TitleEntryFocusOut)
	mainObjects.TitleAddBEntry.Connect("changed", TitleEntryFocusOut)
	mainObjects.TitleAddBFileEntry.Connect("changed", TitleAddToFileEntryEvent)
	mainObjects.TitleApplyButton.Connect("clicked", TitleApplyButtonClicked)
	mainObjects.TitleCumulativeDndChk.Connect("clicked", CumulativeDndChkClicked)
	mainObjects.TitleScanSubDirChk.Connect("clicked", ScanSubDirChkChanged)
	mainObjects.TitleSepEntry.Connect("changed", TitleEntryFocusOut)
	mainObjects.TitleSpinbutton.Connect("value-changed", TitleEntryFocusOut)
	mainObjects.TitleTextview.Connect("event", TitleEntryFocusOut)
	mainObjects.TopImageEventbox.Connect("button-release-event", imgTopReleaseEvent)
}
