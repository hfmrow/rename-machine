// popupMenu.go

/*
*	©2019 H.F.M. MIT license
*	Make popup menu
 */
/* Usage:

var err error
var tstF = func(a, b string) {
	fmt.Printf("a: %s, b: %s", a, b)
}

popup := new(gi.PopupMenu)
popup.WithIcons = false
popup.PopupAddItem("_tst", func() {
	tstF("10", "20")
}, "")
if p, err = popup.PopupMenuBuild(); err != nil {
	fmt.Println(err)
}

// Called with: p.PopupAtPointer(event)

*/

package gtk3Import

import (
	"github.com/gotk3/gotk3/gtk"
)

type PopupMenu struct {
	WithIcons  bool
	items      []*gtk.MenuItem
	separators []*gtk.SeparatorMenuItem
}

// PopupAddItem:
func (p *PopupMenu) PopupAddItem(lbl string, activateFunction interface{}, icon ...interface{}) (err error) {
	var menuItem *gtk.MenuItem
	var image interface{}
	if len(icon) != 0 {
		image = icon[0]
	}
	if p.WithIcons {
		menuItem, err = p.menuItemNewWithImage(lbl, image)
	} else {
		menuItem, err = gtk.MenuItemNewWithMnemonic(lbl)
	}
	if err == nil {
		menuItem.Connect("activate", activateFunction.(func()))
		p.items = append(p.items, menuItem)
		p.separators = append(p.separators, nil)
	}
	return err
}

// PopupAddSeparator:
func (p *PopupMenu) PopupAddSeparator() (err error) {
	if separatorItem, err := gtk.SeparatorMenuItemNew(); err == nil {
		p.items = append(p.items, nil)
		p.separators = append(p.separators, separatorItem)
	}
	return err
}

// PopupMenuBuild: Build popupmenu
func (p *PopupMenu) PopupMenuBuild() (menu *gtk.Menu, err error) {
	if menu, err = gtk.MenuNew(); err == nil {
		for idx, menuItem := range p.items {
			if p.separators[idx] != nil {
				menu.Append(p.separators[idx])
			} else {
				menu.Append(menuItem)
			}
		}
		menu.ShowAll()
	}
	return menu, err
}

func (p *PopupMenu) menuItemNewWithImage(label string, icon interface{}) (menuItem *gtk.MenuItem, err error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	if err == nil {
		image, err := gtk.ImageNew()
		if err == nil {
			SetImage(image, icon, 14)
			label, err := gtk.LabelNewWithMnemonic(label)
			if err == nil {
				menuItem, err = gtk.MenuItemNew()
				if err == nil {
					label.SetHAlign(gtk.ALIGN_START)
					box.Add(image)
					box.PackEnd(label, true, true, 8)
					box.SetHAlign(gtk.ALIGN_START)
					menuItem.Container.Add(box)
					menuItem.ShowAll()
				}
			}
		}
	}
	return menuItem, err
}
