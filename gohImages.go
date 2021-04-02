// gohImages.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 14:58:19 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 hfmrow - Rename Machine v1.6.1 github.com/hfmrow/rename-machine
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/**********************************************************/
/* This section preserve user modifications on update.   */
/* Images declarations, used to initialize objects with */
/* The SetPict() func, accept both kind of variables:  */
/* filename or []byte content in case of using        */
/* embedded binary data. The variables names are the */
/* same. "assetsDeclarationsUseEmbedded(bool)" func */
/* could be used to toggle between filenames and   */
/* embedded binary type. See SetPict()            */
/* declaration to learn more on how to use it.   */
/************************************************/
func assignImages() {
	SetPict(mainObjects.MainWindow, renameMultiDocuments48x48)
	SetPict(mainObjects.MoveApplyButton, toggleOn18x18)
	SetPict(mainObjects.OverImageTop, renameMachine400x27)
	SetPict(mainObjects.OverOkButton, checked18x18)
	SetPict(mainObjects.OverResetButton, cancel18x18)
	SetPict(mainObjects.OverWindow, renameMultiDocuments48x48)
	SetPict(mainObjects.RenApplyButton, toggleOn18x18)
	SetPict(mainObjects.RenKeepBtwButton, keepBetween18)
	SetPict(mainObjects.RenRegexButton, regex18x18)
	SetPict(mainObjects.RenSubButton, substract18)
	SetPict(mainObjects.SingleCancelButton, cancel18x18)
	SetPict(mainObjects.SingleImageTop, renameMachine400x27)
	SetPict(mainObjects.SingleOkButton, renameDocument18x18)
	SetPict(mainObjects.SingleResetButton, reset18x18)
	SetPict(mainObjects.SingleSwMultiButton, renameMultiDocuments18x18)
	SetPict(mainObjects.SingleWindow, renameMultiDocuments48x48)
	SetPict(mainObjects.TitleApplyButton, toggleOn18x18)
	SetPict(mainObjects.TopImage, renameMachine700x48)
}

/**********************************************************/
/* This section is rewritten on assets update.           */
/* Assets var declarations, this step permit to make a  */
/* bridge between the differents types used, string or */
/* []byte, and to simply switch from one to another.  */
/*****************************************************/
var mainGlade interface{}                 // assets/glade/main.glade
var archiveFolder18x18 interface{}        // assets/images/archive-folder-18x18.png
var cancel18x18 interface{}               // assets/images/cancel-18x18.png
var checked18x18 interface{}              // assets/images/checked-18x18.png
var keepBetween18 interface{}             // assets/images/keep-between-18.png
var keepBetween48 interface{}             // assets/images/keep-between-48.png
var regex18x18 interface{}                // assets/images/regex-18x18.png
var regex48x48 interface{}                // assets/images/regex-48x48.png
var renameDocument18x18 interface{}       // assets/images/rename-document-18x18.png
var renameMachine400x27 interface{}       // assets/images/rename-machine-400x27.png
var renameMachine700x48 interface{}       // assets/images/rename-machine-700x48.png
var renameMultiDocuments18x18 interface{} // assets/images/rename-multi-documents-18x18.png
var renameMultiDocuments48x48 interface{} // assets/images/rename-multi-documents-48x48.png
var reset18x18 interface{}                // assets/images/reset-18x18.png
var substract18 interface{}               // assets/images/substract-18.png
var substract48 interface{}               // assets/images/substract-48.png
var toggleOn18x18 interface{}             // assets/images/toggle-on-18x18.png
