// gohAssets.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 14:55:50 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 H.F.M - Rename Machine v1.6.1 github.com/hfmrow/rename-machine
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"embed"
	"log"
)

//go:embed assets/glade
//go:embed assets/images
var embeddedFiles embed.FS

// This functionality does not require explicit encoding of the files, at each
// compilation, the files are inserted into the resulting binary. Thus, updating
// assets is only required when new files are added to be embedded in order to
// create and declare the variables to which the files are linked.
// assetsDeclarationsUseEmbedded: Use native Go 'embed' package to include files
// content at runtime.
func assetsDeclarationsUseEmbedded(embedded ...bool) {
	mainGlade = readEmbedFile("assets/glade/main.glade")
	archiveFolder18x18 = readEmbedFile("assets/images/archive-folder-18x18.png")
	cancel18x18 = readEmbedFile("assets/images/cancel-18x18.png")
	checked18x18 = readEmbedFile("assets/images/checked-18x18.png")
	keepBetween18 = readEmbedFile("assets/images/keep-between-18.png")
	keepBetween48 = readEmbedFile("assets/images/keep-between-48.png")
	regex18x18 = readEmbedFile("assets/images/regex-18x18.png")
	regex48x48 = readEmbedFile("assets/images/regex-48x48.png")
	renameDocument18x18 = readEmbedFile("assets/images/rename-document-18x18.png")
	renameMachine400x27 = readEmbedFile("assets/images/rename-machine-400x27.png")
	renameMachine700x48 = readEmbedFile("assets/images/rename-machine-700x48.png")
	renameMultiDocuments18x18 = readEmbedFile("assets/images/rename-multi-documents-18x18.png")
	renameMultiDocuments48x48 = readEmbedFile("assets/images/rename-multi-documents-48x48.png")
	reset18x18 = readEmbedFile("assets/images/reset-18x18.png")
	substract18 = readEmbedFile("assets/images/substract-18.png")
	substract48 = readEmbedFile("assets/images/substract-48.png")
	toggleOn18x18 = readEmbedFile("assets/images/toggle-on-18x18.png")
}

// readEmbedFile: read 'embed' file system and return []byte data.
func readEmbedFile(filename string) (out []byte) {
	var err error
	out, err = embeddedFiles.ReadFile(filename)
	if err != nil {
		log.Printf("Unable to read embedded file: %s, %v\n", filename, err)
	}
	return
}
