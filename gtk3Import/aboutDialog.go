// aboutDialog.go

//  Source file auto-generated on Thu, 03 Feb 2019 00:34:15 using Gotk3ObjHandler v1.0 ©2019 H.F.M

package gtk3Import

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Control over all used objects from glade.
var aboutDlgControlObj *aboutDlgControlsObj

// Assets var declarations
var aboutboxGlade interface{} // gtk3Import/aboutBox.glade

/******************************/
/* Main structure Declaration */
/******************************/

type aboutDlgControlsObj struct {
	mainUiBuilder               *gtk.Builder
	AboutboxButtonOk            *gtk.Button
	AboutboxLabelAppName        *gtk.Label
	AboutboxLabelDescription    *gtk.Label
	AboutboxLabelLicense        *gtk.Label
	AboutboxLabelRepolink       *gtk.Label
	AboutboxLabelTxtDescription *gtk.Label
	AboutboxLabelTxtSourceRepo  *gtk.Label
	AboutboxLabelVersion        *gtk.Label
	AboutboxLabelYearCreator    *gtk.Label
	AboutboxWindow              *gtk.Window
	AboutTopImage               *gtk.Image
}

/*
*	Aboutbox implementation
 */
type AboutInfos struct {
	titleBox       string
	AppName        string
	AppVers        string
	AppCreats      string
	YearCreat      string
	LicenseAbrv    string
	LicenseShort   string
	Repository     string
	Description    string
	imageTop       interface{}
	imageOkButton  interface{}
	AboutDialogBox *gtk.Window
	ready          bool
	initialised    bool
}

func (ab *AboutInfos) InitFillInfos(
	TitleBox,
	AppName,
	AppVers,
	AppCreat,
	YearCreat,
	LicenseAbrv,
	LicenseShort,
	Repository,
	Description string,
	topImage,
	okBtnIcon interface{}) {
	ab.titleBox = TitleBox
	ab.AppName = AppName
	ab.AppVers = AppVers
	ab.AppCreats = AppCreat
	ab.YearCreat = "©" + YearCreat
	ab.LicenseAbrv = LicenseAbrv
	ab.LicenseShort = LicenseShort
	ab.Repository = Repository
	ab.Description = Description
	ab.imageTop = topImage
	ab.imageOkButton = okBtnIcon
	ab.ready = true
}

// Display label with lower size fonts
func markupLabel(label *gtk.Label, text string) {
	pm := PangoMarkup{}
	pm.Init(text)
	mtype := [][]string{{"fsz", "small"}, {"fgc", pm.Colors.Brown}}
	pm.AddTypes(mtype...)
	label.SetMarkup(pm.Markup())
}

// Display label with lower size fonts
func markupLabels(appName, repo, License string) (outAppName, outRepo, outLicense string) {
	pm := PangoMarkup{}
	pm.Init(appName)
	mtype := [][]string{{"fsz", "x-large"}, {"bld"}, {"fgc", pm.Colors.Brown}}
	pm.AddTypes(mtype...)
	outAppName = pm.Markup()

	pm.Init(repo)
	mtype = [][]string{{"url", "https://" + repo}}
	pm.AddTypes(mtype...)
	outRepo = pm.Markup()

	outLicense = License // Search for http adress to be treated as clickable link
	reg := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(\/\S*)?`)
	indexes := reg.FindAllIndex([]byte(License), 1)

	if len(indexes) != 0 {
		pm.Init(License)
		pm.AddPosition([]int{indexes[0][0], indexes[0][1]})
		mtype = [][]string{{"url", reg.FindString(License)}}
		pm.AddTypes(mtype...)
		outLicense = pm.MarkupAtPos()
	}
	return outAppName, outRepo, outLicense
}

func (infos *AboutInfos) Show() {
	if infos.ready && !infos.initialised {
		aboutDlgStart(infos.titleBox, 100, 100, true)
		infos.AboutDialogBox = aboutDlgControlObj.AboutboxWindow
		//		aboutDlgControlObj.AboutboxWindow.SetResizable(false)

		name, repo, lic := markupLabels(infos.AppName, infos.Repository, infos.LicenseShort)
		aboutDlgControlObj.AboutboxLabelAppName.SetMarkup("\n" + name)
		aboutDlgControlObj.AboutboxLabelVersion.SetLabel(infos.AppVers + "\n")
		aboutDlgControlObj.AboutboxLabelYearCreator.SetLabel(infos.YearCreat + " " + infos.AppCreats + "\n")
		aboutDlgControlObj.AboutboxLabelDescription.SetLabel(infos.Description + "\n")
		aboutDlgControlObj.AboutboxLabelDescription.SetProperty("wrap", true)
		aboutDlgControlObj.AboutboxLabelRepolink.SetMarkup(repo + "\n")
		aboutDlgControlObj.AboutboxLabelLicense.SetMarkup("\n" + lic + "\n")
		/*	Image/Icons assignation functions */
		SetImage(aboutDlgControlObj.AboutTopImage, infos.imageTop)
		SetButtonImage(aboutDlgControlObj.AboutboxButtonOk, infos.imageOkButton)
		infos.initialised = true
	}
	if infos.initialised {
		infos.AboutDialogBox.Show()
	}
}

func (infos *AboutInfos) Hide() {
	hideAboutDialogWindow()
}

/******************************/
/* Gtk3 Window Initialisation */
/******************************/
func aboutDlgStart(winTitle string, width, height int, center bool) {
	// initMainGladeGuiVar: initialise embedded content.
	initMainGladeGuiVar()
	// Assign main ctrl objects
	aboutDlgControlObj = new(aboutDlgControlsObj)
	// buil interface
	if newBuilder(aboutboxGlade) == nil {
		// Parse Gtk objects
		gladeObjParser()
		// Objects Signals initialisations
		signalsPropHandler()
		// Set Window Properties
		if center {
			aboutDlgControlObj.AboutboxWindow.SetPosition(gtk.WIN_POS_CENTER)
		}
		aboutDlgControlObj.AboutboxWindow.SetTitle(winTitle)
		aboutDlgControlObj.AboutboxWindow.SetSizeRequest(width, height)
	} else {
		log.Fatal("Builder initialisation error.")
	}
}

/***************************/
/* Signals Implementations */
/***************************/
func signalsPropHandler() {
	aboutDlgControlObj.AboutboxButtonOk.Connect("clicked", hideAboutDialogWindow)
	aboutDlgControlObj.AboutboxWindow.Connect("delete-event", hideAboutDialogWindow)
}

/***************************/
/* Controls Initialisation */
/***************************/

func gladeObjParser() {
	aboutDlgControlObj.AboutboxButtonOk = loadObject("AboutboxButtonOk").(*gtk.Button)
	aboutDlgControlObj.AboutboxLabelAppName = loadObject("AboutboxLabelAppName").(*gtk.Label)
	aboutDlgControlObj.AboutboxLabelDescription = loadObject("AboutboxLabelDescription").(*gtk.Label)
	aboutDlgControlObj.AboutboxLabelLicense = loadObject("AboutboxLabelLicense").(*gtk.Label)
	aboutDlgControlObj.AboutboxLabelRepolink = loadObject("AboutboxLabelRepolink").(*gtk.Label)
	aboutDlgControlObj.AboutboxLabelTxtDescription = loadObject("AboutboxLabelTxtDescription").(*gtk.Label)
	aboutDlgControlObj.AboutboxLabelTxtSourceRepo = loadObject("AboutboxLabelTxtSourceRepo").(*gtk.Label)
	aboutDlgControlObj.AboutboxLabelVersion = loadObject("AboutboxLabelVersion").(*gtk.Label)
	aboutDlgControlObj.AboutboxLabelYearCreator = loadObject("AboutboxLabelYearCreator").(*gtk.Label)
	aboutDlgControlObj.AboutboxWindow = loadObject("AboutboxWindow").(*gtk.Window)
	aboutDlgControlObj.AboutTopImage = loadObject("AboutTopImage").(*gtk.Image)
}

/*******************************************************/
/* Functions declarations, used to initialize objects */
/*****************************************************/
// newBuilder: initialise builder with glade xml string
func newBuilder(varPath interface{}) (err error) {
	if aboutDlgControlObj.mainUiBuilder, err = gtk.BuilderNew(); err == nil {
		if Gtk3Interface, err := getBytesFromVarAsset(varPath); err == nil {
			err = aboutDlgControlObj.mainUiBuilder.AddFromString(string(Gtk3Interface))
		}
	}
	return err
}

// loadObject: Load GtkObject to be transtyped ...
func loadObject(name string) (newObj glib.IObject) {
	var err error
	newObj, err = aboutDlgControlObj.mainUiBuilder.GetObject(name)
	if err != nil {
		log.Panic(err)
	}
	return newObj
}

// WindowDestroy is the triggered handler when closing/destroying the gui window.
func WindowDestroy() {
	// Bye ...
	gtk.MainQuit()
}

// Signal handler delete_event (hidding window)
func hideAboutDialogWindow() bool {
	if aboutDlgControlObj.AboutboxWindow.GetVisible() {
		aboutDlgControlObj.AboutboxWindow.Hide()
	}
	return true
}

// getBytesFromVarAsset: Get []byte representation from file or asset, depending on type
func getBytesFromVarAsset(varPath interface{}) (outBytes []byte, err error) {
	//	outBytes = new([]byte)
	var rBytes []byte
	switch reflect.TypeOf(varPath).String() {
	case "string":
		rBytes, err = ioutil.ReadFile(varPath.(string))
	case "[]uint8":
		rBytes = varPath.([]byte)
	}
	return rBytes, err
}

// HexToBytes: Convert Gzip Hex to []byte used for embedded binary in source code
func HexToBytes(varPath string, gzipData []byte) (outByte []byte) {
	r, err := gzip.NewReader(bytes.NewBuffer(gzipData))
	if err == nil {
		var bBuffer bytes.Buffer
		_, err = io.Copy(&bBuffer, r)
		if err == nil {
			err = r.Close()
			if err == nil {
				return bBuffer.Bytes()
			}
		}
	}
	if err != nil {
		fmt.Printf("An error occurred while reading: %s\n%v\n", varPath, err.Error())
	}
	return outByte
}

// initMainGladeGuiVar: initialise embedded content.
func initMainGladeGuiVar() {
	aboutboxGlade = HexToBytes("aboutboxGlade", []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xec\x9a\x6d\x8f\x9b\x46\x10\xc7\xdf\xdf\xa7\xd8\xee\xab\x56\x95\xed\xb3\xaf\x49\xd3\xca\x26\x4a\x5b\x25\x3a\x29\x69\xa4\x9e\xdb\xaa\x52\x24\x6b\x81\x31\x4c\x58\x76\xc9\xee\xe0\x87\x7e\xfa\x08\x63\xeb\xcc\x05\x9f\x81\xcb\x45\x70\xe2\x9d\xbd\xcc\x1f\x66\x66\x7f\x33\x3c\xec\x4e\x5f\x6e\x62\xc9\x56\x60\x2c\x6a\x35\xe3\xe3\xe1\x25\x67\xa0\x3c\xed\xa3\x0a\x66\xfc\xef\xf9\xeb\xc1\x0b\xfe\xd2\xb9\x98\x7e\x37\x18\xb0\x37\xa0\xc0\x08\x02\x9f\xad\x91\x42\x16\x48\xe1\x03\xbb\x1a\x4e\x26\xc3\x31\x1b\x0c\x9c\x8b\x29\x2a\x02\xb3\x14\x1e\x38\x17\x8c\x4d\x0d\x7c\x4a\xd1\x80\x65\x12\xdd\x19\x0f\x28\xfa\x91\xdf\x5e\xe8\x6a\x38\x7e\xce\x47\x3b\x3b\xed\x7e\x04\x8f\x98\x27\x85\xb5\x33\xfe\x86\xa2\x7f\x51\xf9\x7a\xcd\x19\xfa\x33\xfe\xca\xd5\x29\xb9\x7a\xb3\x1f\xcb\x04\x8c\x4d\x13\xa3\x13\x30\xb4\x65\x4a\xc4\x30\xe3\x2b\xb4\xe8\x4a\xe0\xce\xdc\xa4\x30\x1d\x1d\x8e\x96\x1b\x7b\x42\x2d\x96\xda\x4b\x6d\x35\xf3\x10\x36\x89\x50\x7e\x35\xe3\x55\x1d\xe3\x58\xfb\x42\x56\x33\x5d\xef\xa2\x5f\x24\xda\x22\xa1\x56\xdc\xf1\x20\x4b\xf5\x39\x19\x6d\x13\x58\x84\xa8\x88\x3b\x3e\x0a\xa9\x83\x73\x02\x1b\x61\xb2\x20\x61\x23\x57\x98\xbd\xb0\xd4\x3d\x2f\x44\xe9\xe7\xbf\xb3\x93\x48\xe1\x41\xa8\xa5\x0f\x66\xb4\x37\x18\x1d\x59\xdc\xb1\xfe\x62\xba\x7f\xd3\x9b\x7c\xae\x5d\xbd\x19\xf3\x83\x5d\xcd\x59\x3e\x33\xd3\xaf\x85\xb4\x95\x34\xda\x20\x28\x12\x79\x96\x57\x60\x08\x3d\x21\x4b\x85\x85\xa8\xca\x23\xbb\x8e\x45\x00\x47\x1c\xcf\x75\x92\x0f\x1d\xcb\x1a\x04\xda\x34\xd8\x52\xbe\x85\xc4\xe0\x24\x51\x27\x5d\x6c\xa4\xb2\xa4\xbd\x88\x3b\x01\x45\x83\x18\xad\x45\x15\x0c\x30\xcb\x47\xb9\x7e\x3a\xca\x13\x5a\x18\x4b\x84\x17\xa1\x0a\xee\xbf\xce\x7d\x55\x78\x4a\xb3\x44\x29\xeb\x29\x6e\xab\xf1\xf2\x54\x00\x5f\x78\x5b\x28\x8c\xaa\x18\x1d\x17\xc8\xa4\x4d\xec\x54\x2e\x96\xf2\x48\xcb\xa3\x7d\x2b\x5c\x90\xc5\xe6\xbf\x1b\x7a\x95\x24\x7f\x8a\xf8\x6e\xed\x34\xce\xc1\x43\xf2\x50\xa6\x95\xb9\xdb\x64\x84\xb2\x52\x90\x70\x25\xcc\xf8\x16\x2c\x77\xf6\x8e\xd7\x39\xd9\xc7\xd4\x12\x2e\xb7\xe7\xea\xab\xbc\x46\x4e\xd7\xc9\x7d\xb5\x52\x3b\xe0\xf3\x05\xd3\xa8\x68\x4e\x14\x4e\x69\xf1\x7c\x0d\xac\xfe\xc9\x9f\x48\xba\x87\xd5\xde\xf1\x1e\xab\x22\x56\xe3\x76\x60\xf5\x1f\x08\xf3\xbb\x01\x41\xda\x74\x0f\xad\xcc\x79\x36\x62\x7b\xff\x7b\xc2\x8a\x84\x4d\xbe\x2d\x61\x37\x90\x08\xb3\xe3\x68\x47\x99\x3d\xfc\x1d\xb7\x95\xab\x58\x98\x00\xd5\x42\xc2\x92\xb8\xf3\xac\x81\xd2\x60\x10\x36\x94\x92\x4e\x9a\x09\x5d\x4d\xa4\xe3\x7b\xb5\x1d\x46\xf6\xaa\x1d\x4d\x71\xbe\xa1\x3f\xc0\x7a\x06\x13\xea\xe4\x2d\xf7\xc8\xf9\x5f\xfb\xae\x58\x44\xec\xa7\x76\x20\xf6\x54\xf8\x62\x04\x1b\x7a\x14\xc6\x4a\xbe\x2e\x19\x91\x9c\x4b\x41\x87\xc9\x7c\xd6\x9a\xe6\x77\xa3\x53\xe3\xc1\x5f\x90\xe8\xee\xb1\x99\xfb\xce\x0c\xec\x12\xab\xcd\xb6\xef\x80\x77\x38\x7b\xde\x0e\xce\x32\xbc\x24\xaa\xa8\x7b\x88\x65\x9e\xb3\xcc\xf5\x9e\xac\x22\x59\x3f\xb7\xe3\x8d\x63\xd2\xbf\x71\xf4\x6f\x1c\x15\x91\x7d\xd1\x8e\x66\xf8\x16\x3d\x50\xb6\x83\x1f\x8d\x3f\xa8\x79\x88\x96\x25\x46\x07\x46\xc4\xcc\xd3\x31\xd8\x7c\xd9\x55\xb8\x56\xcb\x94\x40\x6e\x99\xd2\x6c\x2d\x8c\x11\x8a\xb6\xc3\x0f\xea\x06\x80\x51\x08\x6c\x1e\x02\x7b\x77\x3d\x67\xfb\xd8\xd9\xf7\xef\xae\xe7\x3f\xb0\xa5\x36\xcc\x07\x12\x28\xed\xf0\x1b\x3d\x57\x82\x94\x98\x58\xfc\x1f\xb8\x13\xa3\xef\xcb\xa7\xf9\x78\xf9\x4b\x3b\x9a\xf3\x55\xdf\x9c\xfb\xe6\x5c\xf5\x1b\xf9\x57\x59\x7b\xe9\xce\xfa\xec\xf8\x91\xd7\x67\x53\x22\xad\x8a\x77\x9f\x7c\xec\x7d\x74\x66\xb9\xf6\xf4\x0d\xe0\x7d\x54\x35\xce\x35\xfa\x14\x2e\x0c\x7c\x4a\xc1\x52\x36\xb9\x97\x95\x57\xf1\x1f\xb8\x58\x5c\x47\x66\xc0\x03\x5c\x81\x5d\xf8\xb0\x14\xa9\xa4\x7a\xea\xc3\x2e\x05\x50\x7e\xdd\x2d\x0a\x35\x24\x85\x56\x35\xa9\xa9\xda\xb7\xa9\xba\xb2\x5d\x8b\xaa\x2b\x3a\xb4\xa7\xc9\x63\xed\x9c\xa8\xb5\x0d\xe0\x21\xa5\x39\x69\x5a\x9a\xc5\x18\x8f\x0e\xde\x1e\x98\x8e\x8e\x36\xa4\x7d\x0e\x00\x00\xff\xff\x82\x40\x8d\xbe\xe9\x26\x00\x00"))
}
