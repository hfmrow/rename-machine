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
	"strings"

	"github.com/gotk3/gotk3/gtk"

	gipf "github.com/hfmrow/gtk3Import/pixbuff"
)

/***************************
* DlgMessage implementation.
 ***************************/
var dialogType = map[string]gtk.MessageType{
	"info": gtk.MESSAGE_INFO, "inf": gtk.MESSAGE_INFO,
	"warning": gtk.MESSAGE_WARNING, "wrn": gtk.MESSAGE_WARNING,
	"question": gtk.MESSAGE_QUESTION, "qst": gtk.MESSAGE_QUESTION,
	"error": gtk.MESSAGE_ERROR, "err": gtk.MESSAGE_ERROR,
	"other": gtk.MESSAGE_OTHER, "oth": gtk.MESSAGE_OTHER,
	"infoWithMarkup": gtk.MESSAGE_INFO, "infMU": gtk.MESSAGE_INFO,
	"warningWithMarkup": gtk.MESSAGE_WARNING, "wrMUn": gtk.MESSAGE_WARNING,
	"questionWithMarkup": gtk.MESSAGE_QUESTION, "qstMU": gtk.MESSAGE_QUESTION,
	"errorWithMarkup": gtk.MESSAGE_ERROR, "errMU": gtk.MESSAGE_ERROR,
	"otherWithMarkup": gtk.MESSAGE_OTHER, "othMU": gtk.MESSAGE_OTHER,
}

// DlgMessage: Display message dialog with multiples buttons.
// return get <0 for cross closed or >-1 correspondig to buttons order representation.
// dlgType accepte: "info", "warning", "question", "error", "other", by adding "WithMarkup" you enable it.
// iconFileName: can be a []byte or a string. '' or nil -> No image
func DialogMessage(window *gtk.Window, dlgType, title, text string, iconFileName interface{}, buttons ...string) (value int) {
	var msgDialog *gtk.MessageDialog
	var box *gtk.Box
	var err error

	msgDialog = gtk.MessageDialogNew(window,
		gtk.DIALOG_MODAL,
		dialogType[dlgType],
		gtk.BUTTONS_NONE,
		"")

	// Image
	switch iconFileName.(type) {
	case string:
		if len(iconFileName.(string)) != 0 {
			if box, err = msgDialog.GetContentArea(); err == nil {
				gipf.SetPict(box, iconFileName, 18)
			} else {
				fmt.Println(fmt.Sprintf("DlgMessage, could not get content area: %s", err))
			}
		}
	}

	msgDialog.SetSkipTaskbarHint(true)
	msgDialog.SetKeepAbove(true)
	// Check for link to make it clickable
	// reg := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(\/\S*)?`)
	// if reg.MatchString(text) {
	// }
	if strings.Contains(dlgType, "WithMarkup") {
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
		if iWidget, err := button.GetParent(); err == nil {
			parent := iWidget.ToWidget()
			parent.SetHAlign(gtk.ALIGN_END)
		}
		button.SetSizeRequest(100, 1)
		button.SetBorderWidth(2)
	}
	result := msgDialog.Run()
	msgDialog.Destroy()
	return int(result)
}
