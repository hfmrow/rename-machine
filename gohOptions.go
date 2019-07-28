// gohOptions.go

// Source file auto-generated on Sun, 28 Jul 2019 07:02:22 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	gi "github.com/hfmrow/renMachine/gtk3Import"
)

// App infos
var Name = "RenameMachine"
var Vers = "v1.5"
var Descr = "Rename and add titles to files"
var Creat = "H.F.M"
var YearCreat = "2018-19"
var LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
var LicenseAbrv = "License (MIT)"
var Repository = "github.com/hfmrow/renMachine"

// Vars declarations
var absoluteRealPath, optFilename = getAbsRealPath()
var mainOptions *MainOpt
var err error
var tempDir string
var devMode bool

var singleEntry []string
var oriFileListstoreCol = [][]string{{"Original filename", "text"}, {"Type", "text"}, {"Full path", "text"}}
var modFileListstoreCol = [][]string{{"Modified filename", "text"}}

var errNbFileDoesNotMatch = "source and destination number files does not match !"
var errCancel = "cancelled"
var modifName ModifName
var bakPresExt bool

var tempScanSubDir bool
var tempShowDirectory bool
var tempPreserveExt bool
var tempDnd bool

type ModifName struct {
	remove  []string
	replace []string
	with    []string
}

type MainOpt struct {
	/* Public, will be saved and restored */
	AboutOptions      *gi.AboutInfos
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
	treeViewDropSet gi.DropSet
	currentTab      uint
}

// Main options initialisation
func (opt *MainOpt) Init() {
	opt.AboutOptions = new(gi.AboutInfos)

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
	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), 0644)
		}
	}
	return err
}
