// gohOptions.go

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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	glfs "github.com/hfmrow/genLib/files"
	gltses "github.com/hfmrow/genLib/tools/errors"
	gltsle "github.com/hfmrow/genLib/tools/log2file"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

// App infos
var (
	Name         = "Rename Machine"
	Vers         = "v1.6.1"
	Descr        = "Rename and add titles to files.\nRemember that on Samba shared networks with\nWindows Os, uppercase and lowercase are not\ntreated as under Linux."
	Creat        = "hfmrow"
	YearCreat    = "2018-21"
	LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
	LicenseAbrv  = "License (MIT)"
	Repository   = "github.com/hfmrow/rename-machine"
)

// Vars declarations
var (
	absoluteRealPath, optFilename = getAbsRealPath()
	mainOptions                   *MainOpt
	err                           error
	tempDir                       string
	doTempDir                     bool
	devMode                       bool

	// Functions mapping
	Check     = gltses.Check
	BaseNoExt = glfs.BaseNoExt

	Logger            *gltsle.Log2FileStruct
	Log2FileStructNew = gltsle.Log2FileStructNew

	singleEntry         []string
	oriFileListstoreCol = [][]string{{"Original filename", "text"}, {"Type", "text"}, {"Full path", "text"}}
	modFileListstoreCol = [][]string{{"Modified filename", "text"}}

	errNbFileDoesNotMatch = "source and destination number files does not match !"
	errCancel             = "cancelled"
	modifName             ModifName
	bakPresExt            bool

	tempScanSubDir    bool
	tempShowDirectory bool
	tempPreserveExt   bool
	tempDnd           bool

	tvsFiles,
	tvsTitle,
	tvsRenam,
	tvsMoves *gitw.TreeViewStructure
)

type ModifName struct {
	remove  []string
	replace []string
	with    []string
}

// to match 2 or more whitespace symbols inside a string
var RemDblSpace = func(inputString string) string {
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	return remInside.ReplaceAllString(inputString, " ")
}

// Change all illegals chars (for path in linux and windows) into "-"
var OsForbiden = func(inputString, replace string) string {
	osForbiden := regexp.MustCompile(`[<>:"/\\|?*]`)
	return RemDblSpace(osForbiden.ReplaceAllString(inputString, replace))
}

type MainOpt struct {
	/* Public, will be saved and restored */
	About             *gidg.AboutInfos
	MainWinWidth      int
	MainWinHeight     int
	LanguageFilename  string
	ExtMask           []string
	ExtMaskRen        []string
	ExtSep            string
	PreserveExt       bool
	PreserveExtSingle bool
	CaseSensitive     bool
	ShowDirectory     bool
	ScanSubDir        bool

	OvercaseSensitive  bool
	OverPosixCharClass bool
	OverPosixStrictMde bool
	IncAtRight         bool

	CumulativeDND bool

	/* Private, will NOT be saved */
	baseDirectory   string
	primeFilesList  []string
	rawFilesList    []string
	orgFilenames    [][]string
	titlFilenames   [][]string
	renFilenames    [][]string
	renFilenamesBak [][]string
	moveFilesList   [][]string
	titlList        []string
	overText        []string
	dndFilesList    []string
	treeViewDropSet *gimc.DragNDropStruct
	currentTab      uint
}

// Main options initialisation
func (opt *MainOpt) Init() {
	opt.About = new(gidg.AboutInfos)

	opt.MainWinWidth = 800
	opt.MainWinHeight = 600
	opt.LanguageFilename = "assets/lang/eng.lang"

	opt.ExtMask = []string{"*"}
	opt.ExtMaskRen = []string{"*"}

	opt.ExtSep = ";"
}

// Variables -> Objects.
func (opt *MainOpt) UpdateObjects(sizeToo ...bool) {
	var All bool
	if len(sizeToo) > 0 {
		All = sizeToo[0]
	}
	if All {
		mainObjects.MainWindow.Resize(opt.MainWinWidth, opt.MainWinHeight)
	}
	mainObjects.RenCumulativeDndChk.SetActive(opt.CumulativeDND)
	mainObjects.MoveCumulativeDndChk.SetActive(opt.CumulativeDND)
	mainObjects.TitleCumulativeDndChk.SetActive(opt.CumulativeDND)

	mainObjects.RenCaseSensChk.SetActive(opt.CaseSensitive)
	mainObjects.RenPresExtChk.SetActive(opt.PreserveExt)

	mainObjects.MoveEntryExtMask.SetText(strings.Join(opt.ExtMask, ";"))
	mainObjects.RenEntryExtMask.SetText(strings.Join(opt.ExtMaskRen, ";"))

	mainObjects.SinglePresExtChk.SetActive(opt.PreserveExtSingle)

	mainObjects.RenShowDirChk.SetActive(opt.ShowDirectory)
	mainObjects.TitleScanSubDirChk.SetActive(opt.ScanSubDir)

	mainObjects.OverCaseSensChk.SetActive(opt.OvercaseSensitive)
	mainObjects.OverCharClassChk.SetActive(opt.OverPosixCharClass)
	mainObjects.OverCharClassStrictModeChk.SetActive(opt.OverPosixStrictMde)

	mainObjects.RenIncrementRightChk.SetActive(opt.IncAtRight)
	mainObjects.RenIncrementRightChk.SetSensitive(false)
	mainObjects.RenIncSepEntry.SetSensitive(false)
	mainObjects.RenIncSpinbutton.SetSensitive(false)
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {

	opt.MainWinWidth, opt.MainWinHeight = mainObjects.MainWindow.GetSize()

	opt.CumulativeDND = mainObjects.RenCumulativeDndChk.GetActive()

	opt.CaseSensitive = mainObjects.RenCaseSensChk.GetActive()
	opt.PreserveExt = mainObjects.RenPresExtChk.GetActive()

	opt.PreserveExtSingle = mainObjects.SinglePresExtChk.GetActive()

	opt.ShowDirectory = mainObjects.RenShowDirChk.GetActive()
	opt.ScanSubDir = mainObjects.TitleScanSubDirChk.GetActive()

	opt.OvercaseSensitive = mainObjects.OverCaseSensChk.GetActive()
	opt.OverPosixCharClass = mainObjects.OverCharClassChk.GetActive()
	opt.OverPosixStrictMde = mainObjects.OverCharClassStrictModeChk.GetActive()

	opt.IncAtRight = mainObjects.RenIncrementRightChk.GetActive()
}

// AboutInfos holder.
type AboutInfos struct {
	AppName      string
	AppVers      string
	AppCreats    string
	YearCreat    string
	LicenseShort string
	LicenseAbrv  string
	Repository   string
	Description  string
}

// Read Options from file
func (opt *MainOpt) Read() (err error) {
	var textFileBytes []byte
	opt.Init()
	if textFileBytes, err = ioutil.ReadFile(optFilename); err == nil {
		err = json.Unmarshal(textFileBytes, &opt)
	}
	tempScanSubDir = mainOptions.ScanSubDir
	tempShowDirectory = mainOptions.ShowDirectory
	tempPreserveExt = mainOptions.PreserveExt
	return err
}

// Write Options to file
func (opt *MainOpt) Write() (err error) {
	var out bytes.Buffer
	var jsonData []byte
	opt.UpdateOptions()

	opt.About.DlgBoxStruct = nil // remove dialog object before saving

	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), 0644)
		}
	}
	return err
}
