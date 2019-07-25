// msgDlg.go

/*
*	©2018 H.F.M. MIT license
*	Handle gtk3 Dialogs.
*
*	Message, Question/Response dialogs, file/dir dialogs, Notifications.
 */

package gtk3Import

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Notify: Display a notify message at the top right of screen.
func Notify(title, text string) {
	const appID = "h.f.m"
	app, _ := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	//Shows an application as soon as the app starts
	app.Connect("activate", func() {
		notif := glib.NotificationNew(title)
		notif.SetBody(text)
		app.SendNotification(appID, notif)
	})
	app.Run(nil)
	app.Quit()
}

var FileChooserAction = map[string]gtk.FileChooserAction{
	"create-folder": gtk.FILE_CHOOSER_ACTION_CREATE_FOLDER,
	"open":          gtk.FILE_CHOOSER_ACTION_OPEN,
	"save":          gtk.FILE_CHOOSER_ACTION_SAVE,
	"select-folder": gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
}

// FileChooser: Display a file chooser dialog.
// "open", "save", "create-folder", "select-folder" as dlgType.
// Set title "" permit auto choice from dialog type.
// options are 1-keepAbove and 2-enablePreviewImages
func FileChooser(window *gtk.Window, dlgType, title, filename string, options ...bool) (outFilename string, result bool, err error) {
	var firstBtn, scndBtn string
	var kpAbove, preview bool // fileNotExist bool

	switch len(options) {
	case 1:
		kpAbove = options[0]
	case 2:
		kpAbove = options[0]
		preview = options[1]
	}

	firstBtn = "Cancel"
	scndBtn = "Ok"

	if len(title) == 0 {
		switch dlgType {
		case "create-folder":
			title = "Create folder"
		case "open":
			title = "Select file to open"
		case "save":
			title = "Select file to save"
		case "select-folder":
			title = "Select directory"
		}
	}
	fileChooser, err := gtk.FileChooserDialogNewWith2Buttons(title,
		window, FileChooserAction[dlgType],
		firstBtn, gtk.RESPONSE_CANCEL,
		scndBtn, gtk.RESPONSE_ACCEPT)
	if err != nil {
		return outFilename, result, err
	}

	if preview {
		if previewImage, err := gtk.ImageNew(); err == nil {
			previewImage.Show()
			var pixbuf *gdk.Pixbuf
			fileChooser.SetPreviewWidget(previewImage)
			fileChooser.Connect("update-preview", func(fc *gtk.FileChooserDialog) {
				if _, err = os.Stat(fc.GetFilename()); !os.IsNotExist(err) {
					if pixbuf, err = gdk.PixbufNewFromFile(fc.GetFilename()); err == nil {
						fileChooser.SetPreviewWidgetActive(true)
						if pixbuf.GetWidth() > 640 || pixbuf.GetHeight() > 480 {
							if pixbuf, err = gdk.PixbufNewFromFileAtScale(fc.GetFilename(), 200, 200, true); err != nil {
								fmt.Printf("Image '%s' cannot be loaded, got error: %s", fc.GetFilename(), err.Error())
							}
						}
						previewImage.SetFromPixbuf(pixbuf)
					} else {
						fileChooser.SetPreviewWidgetActive(false)
					}
				}
			})
		}
	}

	if dlgType == "save" {
		fileChooser.SetCurrentName(filepath.Base(filename))
	}
	fileChooser.SetCurrentFolder(filepath.Dir(filename))
	fileChooser.SetDoOverwriteConfirmation(true)
	fileChooser.SetModal(true)
	fileChooser.SetSkipPagerHint(true)
	fileChooser.SetSkipTaskbarHint(true)
	fileChooser.SetKeepAbove(kpAbove)

	switch int(fileChooser.Run()) {
	case -3:
		result = true
		outFilename = fileChooser.GetFilename()
	}

	fileChooser.Destroy()
	return outFilename, result, err
}

var dialogType = map[string]gtk.MessageType{
	"info": gtk.MESSAGE_INFO, "inf": gtk.MESSAGE_INFO,
	"warning": gtk.MESSAGE_WARNING, "wrn": gtk.MESSAGE_WARNING,
	"question": gtk.MESSAGE_QUESTION, "qst": gtk.MESSAGE_QUESTION,
	"error": gtk.MESSAGE_ERROR, "err": gtk.MESSAGE_ERROR,
	"other": gtk.MESSAGE_OTHER, "oth": gtk.MESSAGE_OTHER,
}

// DlgMessage: Display message dialog whit multiples buttons, text inside accet markup format,
// return get <0 for cross closed or >-1 correspondig to buttons order representation.
// "info", "warning", "question", "error", "other" as dlgType.
// iconFileName: "" = No image
func DlgMessage(window *gtk.Window, dlgType, title, text string, iconFileName interface{}, buttons ...string) (value int) {
	// Build dialog
	msgDialog := gtk.MessageDialogNew(window,
		gtk.DIALOG_MODAL,
		dialogType[dlgType],
		gtk.BUTTONS_NONE,
		"")

	box, err := msgDialog.GetContentArea()
	if err == nil {
		setBoxImage(box, iconFileName)
	}

	msgDialog.SetSkipTaskbarHint(true)
	msgDialog.SetKeepAbove(true)
	// Check for link to ake it clickable
	reg := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(\/\S*)?`)
	if reg.MatchString(text) {
		msgDialog.SetProperty("use-markup", true)
	}
	msgDialog.SetProperty("text", text)
	msgDialog.SetTitle(title)
	// Add button(s)
	for idx, btn := range buttons {
		button, err := msgDialog.AddButton(btn, gtk.ResponseType(idx))
		if err != nil {
			log.Fatal(btn+" button could not be created: ", err)
		}
		parent, _ := button.GetParent()
		parent.SetHAlign(gtk.ALIGN_END)
		button.SetSizeRequest(100, 1)
		button.SetBorderWidth(2)
	}
	result := msgDialog.Run()
	msgDialog.Destroy()
	return int(result)
}
