// dnd.go

// Source file auto-generated on Tue, 23 Jul 2019 04:14:20 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
	Drag & drop handling .

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"bytes"
	"log"
	"net/url"
	"reflect"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type DropSet struct {
	Object    interface{} // gtkObject that receive DND
	FilesList []string    // Contain the files list received
	CallBack  func()      // Function called after receiving datas
}

// initDropSets: configure controls to receive dndcontent.
func (ds *DropSet) InitDropSet(objects interface{}, filesList []string, callBack func()) {
	ds.Object = objects
	ds.CallBack = callBack
	ds.FilesList = filesList
	var targets []gtk.TargetEntry // Build dnd context
	te, err := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 0)
	if err != nil {
		log.Fatal(err)
	}
	targets = append(targets, *te)
	objType := reflect.TypeOf(ds.Object).String()
	switch objType {

	case "*gtk.TreeView":
		ds.Object.(*gtk.TreeView).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.TreeView).Connect("drag-data-received", ds.dnDFilesReceived)
	case "gtk.Button":
		ds.Object.(*gtk.Button).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.Button).Connect("drag-data-received", ds.dnDFilesReceived)
	}
}

// ButtonInFilesReceived: Store in files list
func (ds *DropSet) dnDFilesReceived(object interface{}, context *gdk.DragContext, x, y int, data_ptr uintptr, info, time uint) {
	ds.FilesList = ds.FilesList[:0]
	data := gtk.GetData(data_ptr)
	list := strings.Split(string(data), getTextEOL(data))
	for _, file := range list {
		if len(file) != 0 {
			if u, err := url.PathUnescape(file); err == nil {
				ds.FilesList = append(ds.FilesList, strings.TrimPrefix(u, "file://"))
			}
		}
	}
	ds.CallBack()
}

// GetTextEOL: Get EOL from text bytes (CR, LF, CRLF) > string
func getTextEOL(inTextBytes []byte) (outString string) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	if bytes.Contains(inTextBytes, bCRLF) {
		return string(bCRLF)
	} else if bytes.Contains(inTextBytes, bCR) {
		return string(bCR)
	}
	return string(bLF)
}
