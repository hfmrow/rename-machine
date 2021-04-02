// entry.go

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"

	"github.com/gotk3/gotk3/gtk"
)

// GetSepEntry: Sanitize and get separated entries
func GetSepEntry(e *gtk.Entry, separator string) (out []string) {
	tmpOut := strings.Split(strings.TrimSpace(GetEntryText(e)), separator)
	if len(tmpOut) > 0 {
		for _, item := range tmpOut {
			if len(item) > 0 {
				out = append(out, strings.TrimSpace(item))
			}
		}
	}
	return
}

// SetSepEntry: set separated entries
func SetSepEntry(e *gtk.Entry, separator string, in []string) {
	e.SetText(strings.Join(in, separator+" "))
}

// GetExtEntry: Sanitize and get extension entries
func GetExtEntry(e *gtk.Entry, separator string) (out []string) {
	tmpOut := strings.Split(strings.TrimSpace(GetEntryText(e)), separator)
	if len(tmpOut) > 0 {
		for _, ext := range tmpOut {
			if len(ext) > 0 {
				tmp := strings.Split(strings.TrimSpace(ext), ".")
				var tmp1 []string
				for _, s := range tmp {
					if len(s) == 0 {
						s = "*"
					}
					tmp1 = append(tmp1, s)
				}
				out = append(out, strings.Join(tmp1, "."))
			}
		}
	}
	return
}

// SetExtEntry: set extension entries
func SetExtEntry(e *gtk.Entry, separator string, in []string) {
	if len(in) > 0 {
		e.SetText(strings.Join(in, separator+" "))
	} else {
		e.SetText("")
	}
}

// GetEntryText: retrieve value of an entry control.
func GetEntryText(entry *gtk.Entry) (outString string) {
	var err error
	if outString, err = entry.GetText(); err != nil {
		log.Printf("GetEntryText: %v", err)
	}
	return
}

// GetEntryTextAsInt: retrieve value of an entry control as integer
func GetEntryTextAsInt(entry *gtk.Entry) (outint int) {
	var err error
	var outString string
	if outString, err = entry.GetText(); err == nil {
		if outint, err = strconv.Atoi(outString); err == nil {
			return
		}
	}
	if err != nil {
		log.Printf("GetEntryTextAsInt: %v", err)
	}
	return
}

// GetEntryChangedEvent: Retrieval of the value entered in real time. This
// means that the '&storeIn' variable will always contain the content of the
// 'entryCtrl' object. Works with: GtkEntry, GtkSearchEntry and GtkSpinButton.
// The 'setDefault' flag determine if the given '&storeIn' that hold value,
// will be set as default. The 'callback' func (if used), controls whether
// the obtained value should be or not modified / recorded (returned 'bool').
// e.g:
// GetEntryChangedEvent(Entry, &entryValue, true, func(val interface{}) bool {
//		valPtr := val.(*string) // or valPtr := val.(*float64) for GtkSpinbutton
//		if *valPtr == "xxx" {
//			*valPtr = "xxx forbidden !"
//		}
//		return true
//	})
func GetEntryChangedEvent(
	entryCtrl interface{},
	storeIn interface{}, setDefault bool,
	callback ...func(value interface{}) bool) (signalHandle glib.SignalHandle) {

	switch eCtrl := entryCtrl.(type) {

	case *gtk.Entry: // GtkEntry
		signalHandle = eCtrl.Connect("changed",
			func(e *gtk.Entry) {

				value, err := e.GetText()
				if err != nil {
					log.Printf("GetEntryChangedEvent/Entry/GetText: %v", err)
					return
				}

				if len(callback) > 0 {
					if !callback[0](&value) {
						return
					}
				}

				*storeIn.(*string) = value
			})

		// Set default value if requested
		if setDefault {
			eCtrl.HandlerBlock(signalHandle)
			defer eCtrl.HandlerUnblock(signalHandle)
			eCtrl.SetText(*storeIn.(*string))
		}

	case *gtk.SearchEntry: // GtkSearchEntry
		signalHandle = eCtrl.Connect("changed",
			func(e *gtk.SearchEntry) {

				value, err := e.GetText()
				if err != nil {
					log.Printf("GetEntryChangedEvent/SearchEntry/GetText: %v", err)
					return
				}

				if len(callback) > 0 {
					if !callback[0](&value) {
						return
					}
				}

				*storeIn.(*string) = value
			})

		// Set default value if requested
		if setDefault {
			eCtrl.HandlerBlock(signalHandle)
			defer eCtrl.HandlerUnblock(signalHandle)
			eCtrl.SetText(*storeIn.(*string))
		}

	case *gtk.SpinButton: // GtkSpinButton
		signalHandle = eCtrl.Connect("changed",
			func(e *gtk.SpinButton) {

				value := e.GetValue()
				if len(callback) > 0 {
					if !callback[0](&value) {
						return
					}
				}

				*storeIn.(*float64) = value
			})

		// Set default value if requested
		if setDefault {
			eCtrl.HandlerBlock(signalHandle)
			defer eCtrl.HandlerUnblock(signalHandle)
			eCtrl.SetValue(*storeIn.(*float64))
		}

	case *gtk.CheckButton: // GtkCheckButton
		signalHandle = eCtrl.Connect("toggled",
			func(chk *gtk.CheckButton) {

				value := chk.GetActive()
				if len(callback) > 0 {
					if !callback[0](&value) {
						return
					}
				}

				*storeIn.(*bool) = value
			})

		// Set default value if requested
		if setDefault {
			eCtrl.HandlerBlock(signalHandle)
			defer eCtrl.HandlerUnblock(signalHandle)
			eCtrl.SetActive(*storeIn.(*bool))
		}

	case *gtk.RadioButton: // GtkRadioButton
		signalHandle = eCtrl.Connect("toggled",
			func(chk *gtk.RadioButton) {

				value := chk.GetActive()
				if len(callback) > 0 {
					if !callback[0](&value) {
						return
					}
				}

				*storeIn.(*bool) = value
			})

		// Set default value if requested
		if setDefault && *storeIn.(*bool) {
			eCtrl.HandlerBlock(signalHandle)
			defer eCtrl.HandlerUnblock(signalHandle)
			eCtrl.SetActive(*storeIn.(*bool))
		}
	}

	return
}
