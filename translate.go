// translate.go

// File generated on Sun, 28 Jul 2019 07:02:49 using Gotk3ObjectsTranslate v1.3 2019 H.F.M

/*
* 	This program comes with absolutely no warranty.
*	See the The MIT License (MIT) for details:
*	https://opensource.org/licenses/mit-license.php
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// initGtkObjectsText: read translations from structure and set them to objects.
func (trans *MainTranslate) initGtkObjectsText() {
	trans.setTextToGtkObjects(&mainObjects.FileListTreeview.Widget, "FileListTreeview")
	trans.setTextToGtkObjects(&mainObjects.MoveApplyButton.Widget, "MoveApplyButton")
	trans.setTextToGtkObjects(&mainObjects.MoveCumulativeDndChk.Widget, "MoveCumulativeDndChk")
	trans.setTextToGtkObjects(&mainObjects.MoveEntryExtMask.Widget, "MoveEntryExtMask")
	trans.setTextToGtkObjects(&mainObjects.MoveFilechooserButton.Widget, "MoveFilechooserButton")
	trans.setTextToGtkObjects(&mainObjects.MoveLabelExtMask.Widget, "MoveLabelExtMask")
	trans.setTextToGtkObjects(&mainObjects.MovePrevTreeview.Widget, "MovePrevTreeview")
	trans.setTextToGtkObjects(&mainObjects.OverCaseSensChk.Widget, "OverCaseSensChk")
	trans.setTextToGtkObjects(&mainObjects.OverCharClassChk.Widget, "OverCharClassChk")
	trans.setTextToGtkObjects(&mainObjects.OverCharClassStrictModeChk.Widget, "OverCharClassStrictModeChk")
	trans.setTextToGtkObjects(&mainObjects.OverEntry.Widget, "OverEntry")
	trans.setTextToGtkObjects(&mainObjects.OverEntry1.Widget, "OverEntry1")
	trans.setTextToGtkObjects(&mainObjects.OverImageTop.Widget, "OverImageTop")
	trans.setTextToGtkObjects(&mainObjects.OverKeepAfterLbl.Widget, "OverKeepAfterLbl")
	trans.setTextToGtkObjects(&mainObjects.OverKeepBeforeLbl.Widget, "OverKeepBeforeLbl")
	trans.setTextToGtkObjects(&mainObjects.OverOkButton.Widget, "OverOkButton")
	trans.setTextToGtkObjects(&mainObjects.OverResetButton.Widget, "OverResetButton")
	trans.setTextToGtkObjects(&mainObjects.RenApplyButton.Widget, "RenApplyButton")
	trans.setTextToGtkObjects(&mainObjects.RenCaseSensChk.Widget, "RenCaseSensChk")
	trans.setTextToGtkObjects(&mainObjects.RenCumulativeDndChk.Widget, "RenCumulativeDndChk")
	trans.setTextToGtkObjects(&mainObjects.RenEntryExtMask.Widget, "RenEntryExtMask")
	trans.setTextToGtkObjects(&mainObjects.RenIncrementChk.Widget, "RenIncrementChk")
	trans.setTextToGtkObjects(&mainObjects.RenIncrementRightChk.Widget, "RenIncrementRightChk")
	trans.setTextToGtkObjects(&mainObjects.RenIncSepEntry.Widget, "RenIncSepEntry")
	trans.setTextToGtkObjects(&mainObjects.RenIncSpinbutton.Widget, "RenIncSpinbutton")
	trans.setTextToGtkObjects(&mainObjects.RenKeepBtwButton.Widget, "RenKeepBtwButton")
	trans.setTextToGtkObjects(&mainObjects.RenLabelExtMask.Widget, "RenLabelExtMask")
	trans.setTextToGtkObjects(&mainObjects.RenPresExtChk.Widget, "RenPresExtChk")
	trans.setTextToGtkObjects(&mainObjects.RenPrevTreeview.Widget, "RenPrevTreeview")
	trans.setTextToGtkObjects(&mainObjects.RenRegexButton.Widget, "RenRegexButton")
	trans.setTextToGtkObjects(&mainObjects.RenRemEntry.Widget, "RenRemEntry")
	trans.setTextToGtkObjects(&mainObjects.RenRemEntry1.Widget, "RenRemEntry1")
	trans.setTextToGtkObjects(&mainObjects.RenRemEntry2.Widget, "RenRemEntry2")
	trans.setTextToGtkObjects(&mainObjects.RenReplEntry.Widget, "RenReplEntry")
	trans.setTextToGtkObjects(&mainObjects.RenReplEntry1.Widget, "RenReplEntry1")
	trans.setTextToGtkObjects(&mainObjects.RenReplEntry2.Widget, "RenReplEntry2")
	trans.setTextToGtkObjects(&mainObjects.RenScanSubDirChk.Widget, "RenScanSubDirChk")
	trans.setTextToGtkObjects(&mainObjects.RenShowDirChk.Widget, "RenShowDirChk")
	trans.setTextToGtkObjects(&mainObjects.RenSubButton.Widget, "RenSubButton")
	trans.setTextToGtkObjects(&mainObjects.RenWthEntry.Widget, "RenWthEntry")
	trans.setTextToGtkObjects(&mainObjects.RenWthEntry1.Widget, "RenWthEntry1")
	trans.setTextToGtkObjects(&mainObjects.RenWthEntry2.Widget, "RenWthEntry2")
	trans.setTextToGtkObjects(&mainObjects.SingleCancelButton.Widget, "SingleCancelButton")
	trans.setTextToGtkObjects(&mainObjects.SingleEntry.Widget, "SingleEntry")
	trans.setTextToGtkObjects(&mainObjects.SingleImageTop.Widget, "SingleImageTop")
	trans.setTextToGtkObjects(&mainObjects.SingleOkButton.Widget, "SingleOkButton")
	trans.setTextToGtkObjects(&mainObjects.SinglePresExtChk.Widget, "SinglePresExtChk")
	trans.setTextToGtkObjects(&mainObjects.SingleResetButton.Widget, "SingleResetButton")
	trans.setTextToGtkObjects(&mainObjects.SingleSwMultiButton.Widget, "SingleSwMultiButton")
	trans.setTextToGtkObjects(&mainObjects.TitleAddAEntry.Widget, "TitleAddAEntry")
	trans.setTextToGtkObjects(&mainObjects.TitleAddBEntry.Widget, "TitleAddBEntry")
	trans.setTextToGtkObjects(&mainObjects.TitleAddBFileEntry.Widget, "TitleAddBFileEntry")
	trans.setTextToGtkObjects(&mainObjects.TitleApplyButton.Widget, "TitleApplyButton")
	trans.setTextToGtkObjects(&mainObjects.TitleCumulativeDndChk.Widget, "TitleCumulativeDndChk")
	trans.setTextToGtkObjects(&mainObjects.TitlePrevTreeview.Widget, "TitlePrevTreeview")
	trans.setTextToGtkObjects(&mainObjects.TitleScanSubDirChk.Widget, "TitleScanSubDirChk")
	trans.setTextToGtkObjects(&mainObjects.TitleSepEntry.Widget, "TitleSepEntry")
	trans.setTextToGtkObjects(&mainObjects.TitleSpinbutton.Widget, "TitleSpinbutton")
	trans.setTextToGtkObjects(&mainObjects.TitleTextview.Widget, "TitleTextview")
	trans.setTextToGtkObjects(&mainObjects.TopImage.Widget, "TopImage")
}

// Translations structure declaration. To be used in main application.
var translate = new(MainTranslate)

// sts: some sentences/words used in the application. Mostly used in Development mode.
// You must add there all sentences used in your application. Or not ...
// They'll be added to language file each time application started
// when "devMode" is set at true.
var sts = map[string]string{
	`no`:        `No`,
	`ok`:        `Ok`,
	`emptyname`: `Empty filename`,
	`file`:      ` file ? `,
	`confirm`:   `Comfirmation !`,
	`savef`:     `Save file`,
	`cancel`:    `Cancel`,
	`deny`:      `Deny`,
	`allow`:     `Allow`,
	`openf`:     `Open file`,
	`errDir`:    `Error reading directory content.`,
	`proceed`:   `Proceed with : `,
	`renErr`:    `Renaming error: `,
	`yes`:       `Yes`,
	`retry`:     `Retry`,
	`mstk`:      `A mistake ...`,
	`filenameErr`: `You got an issue with:
`,
	`none`:         `None`,
	`all`:          `All`,
	`fileExist`:    `: file exists`,
	`alreadyExist`: `Filename(s) already exists:`,
	`dupFile`: `
Duplicate filename(s):`,
}

// Translations structure with methods
type MainTranslate struct {
	// Public
	ProgInfos    progInfo
	Language     language
	Options      parsingFlags
	ObjectsCount int
	Objects      []object
	Sentences    map[string]string
	// Private
	objectsLoaded bool
}

// MainTranslateNew: Initialise new translation structure and assign language file content to GtkObjects.
// devModeActive, indicate that the new sentences must be added to previous language file.
func MainTranslateNew(filename string, devModeActive ...bool) (mt *MainTranslate) {
	mt = new(MainTranslate)
	if _, err := os.Stat(filename); err == nil {
		mt.read(filename)
		mt.initGtkObjectsText()
		if len(devModeActive) != 0 {
			if devModeActive[0] {
				mt.Sentences = sts
				err := mt.write(filename)
				if err != nil {
					fmt.Printf("%s\n%s\n", "Cannot write actual sentences to language file.", err.Error())
				}
			}
		}
	} else {
		fmt.Printf("%s\n%s\n", "Error loading language file !\nNot an error when is just creating from glade Xml or GOH project file.", err.Error())
	}
	return mt
}

// sortId: sort by numbering methode
func (trans *MainTranslate) sortId() {
	var tmpWordList []string
	for _, wrd := range trans.Objects {
		tmpWordList = append(tmpWordList, wrd.Id)
	}
	numberedWords := new(WordWithDigit)
	numberedWords.Init(tmpWordList)
	sort.SliceStable(trans.Objects, func(i, j int) bool {
		return numberedWords.FillWordToMatchMaxLength(trans.Objects[i].Id) < numberedWords.FillWordToMatchMaxLength(trans.Objects[j].Id)
	})
}

// buildIdx: build index for each object.
func (trans *MainTranslate) buildIdx() {
	trans.sortId()
	for idx, _ := range trans.Objects {
		trans.Objects[idx].Idx = idx
	}
}

// readFile: language file.
func (trans *MainTranslate) read(filename string) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		if err = json.Unmarshal(textFileBytes, &trans); err == nil {
			trans.objectsLoaded = true
		}
	}
	return err
}

// Write json datas to file
func (trans *MainTranslate) write(filename string) (err error) {
	var out bytes.Buffer
	var jsonData []byte
	if jsonData, err = json.Marshal(&trans); err == nil && trans.objectsLoaded {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), 0644)
		}
	}
	return err
}

type parsingFlags struct {
	SkipLowerCase  bool
	SkipEmptyLabel bool
	SkipEmptyName  bool
	DoBackup       bool
}

type progInfo struct {
	Name              string
	Version           string
	Creat             string
	MainObjStructName string
	GladeXmlFilename  string
	TranslateFilename string
}

type language struct {
	LangNameLong string
	LangNameShrt string
	Author       string
	Date         string
	Updated      string
	Contributors []string
}

type object struct {
	Class         string
	Id            string
	Label         string
	LabelMarkup   bool
	LabelWrap     bool
	Tooltip       string
	TooltipMarkup bool
	Text          string
	Uri           string
	Comment       string
	Idx           int
}

// Define available property within objects
type propObject struct {
	Class         string
	Label         bool
	LabelMarkup   bool
	LabelWrap     bool
	Tooltip       bool
	TooltipMarkup bool
	Text          bool
	Uri           bool
}

// Property that exists for Gtk3 Object ...	(Used for Class capability)
var propPerObjects = []propObject{
	{Class: "GtkButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkToggleButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkLabel", Label: true, LabelMarkup: true, Tooltip: true, TooltipMarkup: true, LabelWrap: true},
	{Class: "GtkSpinButton", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkEntry", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkCheckButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkProgressBar", Tooltip: true, TooltipMarkup: true, Text: true},
	{Class: "GtkSearchBar", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkImage", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkRadioButton", Label: true, LabelMarkup: false, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkComboBoxText", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkComboBox", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkLinkButton", Label: true, Tooltip: true, TooltipMarkup: true, Uri: true},
	{Class: "GtkSwitch", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkTreeView", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkFileChooserButton", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkTextView", Tooltip: true, TooltipMarkup: true},
}

// setTextToGtkObjects: read translations from structure and set them to object.
// like this: setTextToGtkObjects(&mainObjects.TransLabelHint.Widget, "TransLabelHint")
func (trans *MainTranslate) setTextToGtkObjects(obj *gtk.Widget, objectId string) {
	for _, currObject := range trans.Objects {
		if currObject.Id == objectId {
			for _, props := range propPerObjects {
				if currObject.Class == props.Class {
					if props.Label {
						obj.SetProperty("label", currObject.Label)
						if props.LabelMarkup {
							obj.SetProperty("use-markup", currObject.LabelMarkup)
							obj.SetProperty("label", strings.ReplaceAll(currObject.Label, "&", "&amp;"))
						}
					}
					if props.LabelWrap {
						obj.SetProperty("wrap", currObject.LabelWrap)
					}
					if props.Tooltip && !currObject.TooltipMarkup {
						obj.SetProperty("tooltip_text", currObject.Tooltip)
					}
					if props.Tooltip && currObject.TooltipMarkup {
						obj.SetProperty("tooltip_markup", strings.ReplaceAll(currObject.Tooltip, "&", "&amp;"))
					}
					if props.Text {
						obj.SetProperty("text", currObject.Text)
					}
					if props.Uri {
						obj.SetProperty("uri", currObject.Uri)
					}
				}
			}
		}
	}
}

// Digit sorting functions
type WordWithDigit struct {
	maxLength, maxLengthLeft int
	zeroMask                 string
	ForceRightDigit          bool
}

func (w *WordWithDigit) Init(words []string) {
	for _, word := range words {
		if len(word) > w.maxLength {
			w.maxLength = len(word)
			if digitsPosition(word) == 0 {
				digits := getDigits(word)
				if len(digits) > w.maxLengthLeft {
					w.maxLengthLeft = len(digits)
				}
			}
		}
	}
}

// FillWordToMatchMaxLength: Convert word(s) into numbered one, like "label1" -> "label000001" etc...
// results are based on list of words that determine max length of them to determinate
// the final length of the modified word. This is used in case of sorting list
// of words that contains numeric value to avoid the disorder result
// like "1label", "10label", "2label" etc ...
func (w *WordWithDigit) FillWordToMatchMaxLength(inString string) (outString string) {
	var word string

	inString = strings.ToLower(strings.TrimSpace(inString))
	zeroCount := w.maxLength - len(inString)
	for idx := 0; idx < zeroCount; idx++ {
		w.zeroMask += "0"
	}
	wordPos := digitsPosition(inString)
	digits := getDigits(inString)
	switch wordPos {
	case 0: // Left
		word = inString[len(digits):len(inString)]
		outString = word + w.zeroMask + digits
	case 1: // Right
		word = inString[:len(inString)-len(digits)]
		outString = word + w.zeroMask + digits
	case -1: // None
		outString = inString + w.zeroMask
	}
	w.zeroMask = ""
	return outString
}

// numPosition: detect position of digit part: 0=left, 1=right, -1=none
func digitsPosition(inString string) int {
	digitS := regexp.MustCompile("^[[:digit:]]")
	digitE := regexp.MustCompile("[[:digit:]]$")
	switch {
	case digitS.MatchString(inString):
		return 0 // Left
	case digitE.MatchString(inString):
		return 1 // Right
	}
	return -1 // None
}

// getDigits: return digit part of string prior at start or at end, -1=none
func getDigits(inString string) (value string) {
	digitS := regexp.MustCompile("(^[0-9]*)")
	digitE := regexp.MustCompile("([0-9]*$)")
	start := digitS.FindString(inString)
	end := digitE.FindString(inString)
	switch {
	case len(start) != 0: // Left
		value = start
	case len(end) != 0: // Right
		value = end
	}
	return value
}
