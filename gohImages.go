// gohImages.go

// Source file auto-generated on Sun, 28 Jul 2019 07:02:22 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/************************************************************/
/* Images declarations, used to initialize objects with it */
/* The functions: setImage, setWinIcon and setButtonImage */
/* accept both kind of variables: filename or []byte     */
/* content in case of using embedded binary data. The   */
/* variables names are the same. You can use function  */
/* "func assetsDeclarationsUseEmbedded(bool)"         */
/* to toggle between filenames and embedded binary.  */
/****************************************************/
func assignImages() {
	setWinIcon(mainObjects.MainWindow, renameMultiDocuments48x48)
	setButtonImage(mainObjects.MoveApplyButton, toggleOn18x18)
	setImage(mainObjects.OverImageTop, renameMachine400x27)
	setButtonImage(mainObjects.OverOkButton, checked18x18)
	setButtonImage(mainObjects.OverResetButton, cancel18x18)
	setWinIcon(mainObjects.OverWindow, renameMultiDocuments48x48)
	setButtonImage(mainObjects.RenApplyButton, toggleOn18x18)
	setButtonImage(mainObjects.RenKeepBtwButton, keepBetween18)
	setButtonImage(mainObjects.RenRegexButton, regex18x18)
	setButtonImage(mainObjects.RenSubButton, substract18)
	setButtonImage(mainObjects.SingleCancelButton, cancel18x18)
	setImage(mainObjects.SingleImageTop, renameMachine400x27)
	setButtonImage(mainObjects.SingleOkButton, renameDocument18x18)
	setButtonImage(mainObjects.SingleResetButton, reset18x18)
	setButtonImage(mainObjects.SingleSwMultiButton, renameMultiDocuments18x18)
	setWinIcon(mainObjects.SingleWindow, renameMultiDocuments48x48)
	setButtonImage(mainObjects.TitleApplyButton, toggleOn18x18)
	setImage(mainObjects.TopImage, renameMachine700x48)
}

// Assets var declarations, this step permit to make a "bridge" between the differents
// types used: (string or []byte) and to simply switch from one to another.
var archiveFolder18x18 interface{}        // assets/images/archive-folder-18x18.png
var cancel18x18 interface{}               // assets/images/cancel-18x18.png
var checked18x18 interface{}              // assets/images/checked-18x18.png
var keepBetween18 interface{}             // assets/images/keep-between-18.png
var keepBetween48 interface{}             // assets/images/keep-between-48.png
var mainGlade interface{}                 // assets/glade/main.glade
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
