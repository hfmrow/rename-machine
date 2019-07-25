// clipboard.go

/*
*	©2019 H.F.M. MIT license
*	This is a simple clipboard handler to use with gotk3 "https://github.com/gotk3/gotk3"
 */

package gtk3Import

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type Clipboard struct {
	Entity *gtk.Clipboard
}

// Initialise clipboard
func (c *Clipboard) Init() (err error) {
	c.Entity, err = gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	return err
}

func (c *Clipboard) GetText() (clipboardContent string, err error) {
	return c.Entity.WaitForText()
}

func (c *Clipboard) SetText(clipboardContent string) {
	c.Entity.SetText(clipboardContent)
}

// Stores the current clipboard data somewhere so that it will stay around after the application has quit.
func (c *Clipboard) Store() {
	c.Entity.Store()
}
