// pixbuffByte.go

/*
	Source file auto-generated on Sun, 20 Oct 2019 13:50:31 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	©2019 H.F.M

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/*************************************************/
/* Images functions, used to initialize objects */
/* You can use it to load your own embedded    */
/* images, icons ...                          */
/*********************************************/

// TryStockIcon: Convenience function to combine a stock-icons-theme
// and an asset/filename icon to send to "GetPixBuf()" function.
func TryStockIcon(iconLookupFlag gtk.IconLookupFlags, stockName, assetName interface{}) (names []interface{}) {
	names = append(names, stockName)
	names = append(names, iconLookupFlag)
	names = append(names, assetName)
	return
}

// OptPict: hold options for SetPict(), you can use OptPictStructNew()
// to set defaults values.
type OptPict struct {
	Size int
	// 0 = limit height, 1 = limit width, 2 = percents,
	// 3 = limit height or width depending on MaxSize value.
	ResizeMethode int
	// Rotate image / icon
	Rotate gdk.PixbufRotation
	// SpinButton / Entry ... icon positions
	Position gtk.EntryIconPosition
	// Opacity [0.01 to 1.0] work with: Box, ToggleButton, ToolButton, MenuButton, Button
	// a value of 0 means 1.0; 0.01 means 0
	Opacity float64
	// Ask for a GdkPixbufAnimation image
	Animated bool
	// Max size accepted relative to "ResizeMethode" = 3
	MaxSize int
}

// OptPictStructNew: 'values' options can be nil, ther is no obligation
// to use this function to create the structure because functions are
// designed to handle defaults options automatically, this is a convenient use.
// The options (each values can be nil) order (if used) must be:
// Size, ResizeMethode, Rotate, Position, Opacity, Animated, MaxSize.
func OptPictStructNew(values ...interface{}) OptPict {

	// Set default options
	opt := OptPict{
		Size:          0,
		ResizeMethode: 0,
		Rotate:        gdk.PIXBUF_ROTATE_NONE,
		Position:      gtk.ENTRY_ICON_PRIMARY,
		Opacity:       1.0,
		Animated:      false,
		MaxSize:       256,
	}

	var fillOptInt = func(i, index int) {
		switch index {

		case 0:
			opt.Size = i
		case 1:
			opt.ResizeMethode = i
		case 6:
			opt.MaxSize = i
		}
	}
	// Fill opt with given values
	if values != nil {
		for idx, val := range values {
			switch val.(type) {

			case nil:
				continue
			case gdk.PixbufRotation:
				opt.Rotate = val.(gdk.PixbufRotation)
			case gtk.EntryIconPosition:
				opt.Position = val.(gtk.EntryIconPosition)
			case float64:
				opt.Opacity = val.(float64)
			case bool:
				opt.Animated = val.(bool)
			case int:
				fillOptInt(val.(int), idx)
			}
		}
	}
	return opt
}

// getOptPict: Handling OptPict for many usages
func _getOptPict(options ...interface{}) OptPict {

	opt := OptPictStructNew()

	if len(options) > 0 {
		switch options[0].(type) {
		case int:
			opt.Size = options[0].(int)
		case interface{}:
			opt = options[0].(OptPict)
		}
	}
	return opt
}

// SetPic: Assign image to an Object depending on type, accept stock-icons-theme,
// filename or []byte. Use OptPict{} struct to set options, size can be applied
// directly as int or be nil.
// The options are described below:
// - Load a gif animated image and specify using animation (true), resizing not
//   allowed with animations.
//     SetPict(GtkImage, "linearProgressHorzBlue.gif", OptPict{Animated: true})
// - Resize to 32 pixels height, keep porportions & assign image to GtkButton.
//     SetPict(GtkButton, "stillImage.png", 32) or
//     SetPict(GtkButton, "stillImage.png", OptPict{Size: 32})
// - With default size, resizing not allowed for GtkSpinButton position at right-end.
//     SetPict(GtkSpinButton, "stillImage.png", OptPict{Position: gtk.ENTRY_ICON_SECONDARY})
// - Load a stock-icons-theme or an asset/filename icon to a GtkMenuButton where stock not available.
//     SetPict(gtk.MenuButton, TryStockIcon("open-menu-symbolic", "stillImage.png"), 18) or
//     SetPict(gtk.MenuButton, TryStockIcon("open-menu-symbolic", "stillImage.png"), OptPict{Size: 18})
// - Rotate upsidedown, resize 14px:
//     SetPict(gtk.MenuButton, TryStockIcon("open-menu-symbolic", "stillImage.png"), OptPict{Size: 14, Rotate: gdk.PIXBUF_ROTATE_UPSIDEDOWN})
func SetPict(iObject, varPath interface{}, options ...interface{}) (err error) {
	var inPixbuf *gdk.Pixbuf
	var inPixbufAnimation *gdk.PixbufAnimation
	var image *gtk.Image

	opt := _getOptPict(options...)

	// Since the default initialization value is set to 0,
	// we need to change it to show a visible image ...
	if opt.Opacity == 0 {
		opt.Opacity = 1
	}

	// Get pixbuff type (normal or animation for GtkImage)
	if opt.Animated {
		if inPixbufAnimation, err = GetPixBufAnimation(varPath); err == nil {
			image, err = gtk.ImageNewFromAnimation(inPixbufAnimation)
		}
	} else {
		if inPixbuf, err = GetPixBuf(varPath, opt); err == nil {
			image, err = gtk.ImageNewFromPixbuf(inPixbuf)
		}
	}

	// Handling the error if there is
	if err != nil {
		var filename string
		// Look for a stock icon or a filename or an ambedded asset
		switch path := varPath.(type) {
		case []interface{}:
			filename = path[2].(string)
		default:
			filename = path.(string)
		}
		if _, err = os.Stat(filename); len(filename) > 0 /*os.IsNotExist(err)*/ {
			err = fmt.Errorf("SetPict: [%v] %v", varPath, err)
			fmt.Println(err.Error()) // double output
		}
		return
	}

	// Objects parsing
	if image != nil {

		image.SetOpacity(opt.Opacity)

		switch object := iObject.(type) {

		case *gtk.Button: // Set Image to GtkButton
			object.SetImage(image)
			object.SetAlwaysShowImage(true)

		case *gtk.MenuButton: // Set Image to GtkMenuButton
			object.SetImage(image)
			object.SetAlwaysShowImage(true)

		case *gtk.ToggleButton: // Set Image to GtkToggleButton
			object.SetImage(image)
			object.SetAlwaysShowImage(true)

		case *gtk.ToolButton: // Set Image to GtkToolButton
			object.SetIconWidget(image)

		case *gtk.Box: // Add Image to GtkBox
			object.Add(image)

		case *gtk.SpinButton: // Set Icon to GtkSpinButton, No resize, No Animate.
			object.SetIconFromPixbuf(opt.Position, inPixbuf)

		case *gtk.ApplicationWindow: // Set Icon to GtkApplicationWindow, No Animate.
			object.SetIcon(inPixbuf)

		case *gtk.Window: // Set Icon to GtkWindow, No Animate.
			object.SetIcon(inPixbuf)

		case *gtk.Image: // Set Image to GtkImage
			if opt.Animated {
				object.SetFromAnimation(inPixbufAnimation)
			} else {
				object.SetFromPixbuf(inPixbuf)
			}
		}
	} else {
		err = fmt.Errorf("Image error: %v", err)
	}
	return
}

// GetPixBuf: Get gdk.PixBuf from stock, filename or []byte, depending
// on type. For an explanation of the options, see above, OptPict{}.
// note: If stock-icons-theme does not exist, the asset/filename icon is used.
func GetPixBuf(varPath interface{}, options ...interface{}) (outPixbuf *gdk.Pixbuf, err error) {
	var pixbufLoader *gdk.PixbufLoader

	opt := _getOptPict(options...)

	// Look for a stock icon and a filename or ambedded asset
	switch varPath.(type) {
	case []interface{}:
		var iconTheme *gtk.IconTheme
		if iconTheme, err = gtk.IconThemeGetDefault(); err == nil {
			stock := varPath.([]interface{})[0].(string)
			iconLookupFlags := varPath.([]interface{})[1].(gtk.IconLookupFlags)
			outPixbuf, err = iconTheme.LoadIcon(stock, opt.Size, iconLookupFlags)
		}
		if err != nil {
			varPath = varPath.([]interface{})[2]
			fmt.Println(err.Error() + ". try to load the secondary choice.")
		}
	}

	if outPixbuf == nil { // in this case, no stock icon was found or was not requested
		switch varPath.(type) {
		case string: // Its a filename
			outPixbuf, err = gdk.PixbufNewFromFile(varPath.(string))
		case []uint8: // Its a binary data (embedded asset)
			if pixbufLoader, err = gdk.PixbufLoaderNew(); err == nil {
				outPixbuf, err = pixbufLoader.WriteAndReturnPixbuf(varPath.([]byte))
			}
		}

		if err == nil && opt.Size != 0 {
			newWidth, newHeight := NormalizeSize(outPixbuf.GetWidth(), outPixbuf.GetHeight(), opt)
			if outPixbuf, err = outPixbuf.ScaleSimple(newWidth, newHeight, gdk.INTERP_HYPER); err == nil {
				outPixbuf, err = outPixbuf.RotateSimple(opt.Rotate)
			}
		}
	}
	return
}

// GetPixBufAnimation: Get gdk.PixBufAnimation from filename or []byte, depending on type
func GetPixBufAnimation(varPath interface{}) (outPixbufAnimation *gdk.PixbufAnimation, err error) {
	var pixbufLoader *gdk.PixbufLoader
	switch varPath.(type) {
	case string:
		outPixbufAnimation, err = gdk.PixbufAnimationNewFromFile(varPath.(string))
	case []uint8:
		if pixbufLoader, err = gdk.PixbufLoaderNew(); err == nil {
			outPixbufAnimation, err = pixbufLoader.WriteAndReturnPixbufAnimation(varPath.([]byte))
		}
	}
	return
}

// NormalizeSize: compute new size with kept proportions based on defined format.
// formats: 0 = limit height, 1 = limit width, 2 = percents,
// 3 = limit height or width depending on MaxSize value.
func NormalizeSize(oldWidth, oldHeight int, options ...interface{}) (outWidth, outHeight int) {

	opt := _getOptPict(options...)

	switch opt.ResizeMethode {
	case 0: // limit Height
		outWidth = int(float64(oldWidth) * (float64(opt.Size) / float64(oldHeight)))
		outHeight = opt.Size
	case 1: // limit Width
		outWidth = opt.Size
		outHeight = int(float64(oldHeight) * (float64(opt.Size) / float64(oldWidth)))
	case 2: // percent
		outWidth = int((float64(oldWidth) * float64(opt.Size)) / 100)
		outHeight = int((float64(oldHeight) * float64(opt.Size)) / 100)
	case 3: // limit Height or Width
		switch {
		case oldWidth >= opt.MaxSize:
			opt.Size = opt.MaxSize
			opt.ResizeMethode = 1
			return NormalizeSize(oldWidth, oldHeight, opt)
		case oldHeight >= opt.MaxSize:
			opt.Size = opt.MaxSize
			opt.ResizeMethode = 0
			return NormalizeSize(oldWidth, oldHeight, opt)
		}
		opt.ResizeMethode = 0
		return NormalizeSize(oldWidth, oldHeight, opt)
	}
	return
}

/***************************************/
/* Embedded data conversion functions */
/* Used to make variable content     */
/* available in go-source           */
/***********************************/
// GetBytesFromVarAsset: Get []byte representation from file or asset, depending on type
func GetBytesFromVarAsset(varPath interface{}) (outBytes []byte, err error) {
	switch reflect.TypeOf(varPath).String() {
	case "string":
		return ioutil.ReadFile(varPath.(string))
	case "[]uint8":
		return varPath.([]byte), err
	}
	return
}

// HexToBytes: Convert Gzip Hex to []byte used for embedded binary in source code
func HexToBytes(varPath string, gzipData []byte) (outByte []byte) {
	r, err := gzip.NewReader(bytes.NewBuffer(gzipData))
	if err == nil {
		var bBuffer bytes.Buffer
		if _, err = io.Copy(&bBuffer, r); err == nil {
			if err = r.Close(); err == nil {
				return bBuffer.Bytes()
			}
		}
	}
	if err != nil {
		fmt.Printf("An error occurred while reading: %s\n%v\n", varPath, err.Error())
	}
	return
}

/*
* Previous VERSION
 */

// /*************************************************/
// /* Images functions, used to initialize objects */
// /* You can use it to load your own embedded    */
// /* images, icons ...                          */
// /*********************************************/

// // SetPic: Assign image to an Object depending on type, accept filename or []byte.
// // options: 1- size (int), 2- enable animation (bool)
// // ie:
// // - Load a gif animated image and specify using animation (true), resizing not allowed with animations.
// //     SetPict(GtkImage, "linearProgressHorzBlue.gif", 0, true)
// // - Resize to 32 pixels height, keep porportions & assign image to GtkButton.
// //     SetPict(GtkButton, "stillImage.png", 32)
// // - With default size, resizing not allowed for GtkSpinButton.
// //     SetPict(GtkSpinButton, "stillImage.png")
// func SetPict(iObject, varPath interface{}, options ...interface{}) {
// 	var err error
// 	var sze int
// 	var inPixbuf *gdk.Pixbuf
// 	var inPixbufAnimation *gdk.PixbufAnimation
// 	var image *gtk.Image
// 	var isAnim bool
// 	pos := gtk.ENTRY_ICON_PRIMARY
// 	// Options parsing
// 	var getSizeOrPos = func() {
// 		switch s := options[0].(type) {
// 		case int:
// 			sze = s
// 		case string:
// 			if s == "right" {
// 				pos = gtk.ENTRY_ICON_SECONDARY
// 			}
// 		}
// 	}
// 	switch len(options) {
// 	case 1:
// 		getSizeOrPos()
// 	case 2: //
// 		getSizeOrPos()
// 		isAnim = options[1].(bool)
// 	}
// 	if isAnim { // PixbufAnimation case
// 		if inPixbufAnimation, err = GetPixBufAnimation(varPath); err == nil {
// 			image, err = gtk.ImageNewFromAnimation(inPixbufAnimation)
// 		}
// 	} else { // Pixbuf case
// 		if inPixbuf, err = GetPixBuf(varPath, sze); err == nil {
// 			image, err = gtk.ImageNewFromPixbuf(inPixbuf)
// 		}
// 	}
// 	if err != nil {
// 		if _, err = os.Stat(varPath.(string)); !os.IsNotExist(err) {
// 			log.Fatalf("SetPict: %v\n%s\n", varPath, err.Error())
// 			return
// 		}
// 	}
// 	// Objects parsing
// 	if image != nil {
// 		switch object := iObject.(type) {
// 		case *gtk.Image: // Set Image to GtkImage
// 			if isAnim {
// 				object.SetFromAnimation(inPixbufAnimation)
// 			} else {
// 				object.SetFromPixbuf(inPixbuf)
// 			}
// 		case *gtk.Window: // Set Icon to GtkWindow, No Animate.
// 			object.SetIcon(inPixbuf)
// 		case *gtk.Button: // Set Image to GtkButton
// 			object.SetImage(image)
// 			object.SetAlwaysShowImage(true)
// 		case *gtk.ToolButton: // Set Image to GtkToolButton
// 			object.SetIconWidget(image)
// 		case *gtk.ToggleButton: // Set Image to GtkToggleButton
// 			object.SetImage(image)
// 			object.SetAlwaysShowImage(true)
// 		case *gtk.SpinButton: // Set Icon to GtkSpinButton. options[0] = "left" or "right", No resize, No Animate.
// 			object.SetIconFromPixbuf(pos, inPixbuf)
// 		case *gtk.Box: // Add Image to GtkBox
// 			object.Add(image)
// 		}
// 	}
// 	return
// }

// // GetPixBuf: Get gdk.PixBuf from filename or []byte, depending on type
// // size: resize height keeping porportions. 0 = no change
// func GetPixBuf(varPath interface{}, size ...int) (outPixbuf *gdk.Pixbuf, err error) {
// 	var pixbufLoader *gdk.PixbufLoader
// 	sze := 0
// 	if len(size) != 0 {
// 		sze = size[0]
// 	}
// 	switch varPath.(type) {
// 	case string:
// 		outPixbuf, err = gdk.PixbufNewFromFile(varPath.(string))
// 	case []uint8:
// 		if pixbufLoader, err = gdk.PixbufLoaderNew(); err == nil {
// 			outPixbuf, err = pixbufLoader.WriteAndReturnPixbuf(varPath.([]byte))
// 		}
// 	}
// 	if err == nil && sze != 0 {
// 		newWidth, newHeight := NormalizeSize(outPixbuf.GetWidth(), outPixbuf.GetHeight(), sze, 2)
// 		outPixbuf, err = outPixbuf.ScaleSimple(newWidth, newHeight, gdk.INTERP_BILINEAR)
// 	}
// 	return
// }

// // GetPixBufAnimation: Get gdk.PixBufAnimation from filename or []byte, depending on type
// func GetPixBufAnimation(varPath interface{}) (outPixbufAnimation *gdk.PixbufAnimation, err error) {
// 	var pixbufLoader *gdk.PixbufLoader
// 	switch varPath.(type) {
// 	case string:
// 		outPixbufAnimation, err = gdk.PixbufAnimationNewFromFile(varPath.(string))
// 	case []uint8:
// 		if pixbufLoader, err = gdk.PixbufLoaderNew(); err == nil {
// 			outPixbufAnimation, err = pixbufLoader.WriteAndReturnPixbufAnimation(varPath.([]byte))
// 		}
// 	}
// 	return
// }

// // NormalizeSize: compute new size with kept proportions based on defined format.
// // format: 0 percent, 1 reducing width, 2 reducing height
// func NormalizeSize(oldWidth, oldHeight, newValue, format int) (outWidth, outHeight int) {
// 	switch format {
// 	case 0: // percent
// 		outWidth = int((float64(oldWidth) * float64(newValue)) / 100)
// 		outHeight = int((float64(oldHeight) * float64(newValue)) / 100)
// 	case 1: // Width
// 		outWidth = newValue
// 		outHeight = int(float64(oldHeight) * (float64(newValue) / float64(oldWidth)))
// 	case 2: // Height
// 		outWidth = int(float64(oldWidth) * (float64(newValue) / float64(oldHeight)))
// 		outHeight = newValue
// 	}
// 	return
// }

// // ResizeImage: Get Resized gtk.Pixbuff image representation from file or []byte, depending on type
// // interp: 0 GDK_INTERP_NEAREST, 1 GDK_INTERP_TILES, 2 GDK_INTERP_BILINEAR (default), 3 GDK_INTERP_HYPER.
// func ResizeImage(varPath interface{}, width, height int, interp ...int) (outPixbuf *gdk.Pixbuf, err error) {
// 	interpolation := gdk.INTERP_BILINEAR
// 	if len(interp) != 0 {
// 		switch interp[0] {
// 		case 0:
// 			interpolation = gdk.INTERP_NEAREST
// 		case 1:
// 			interpolation = gdk.INTERP_TILES
// 		case 3:
// 			interpolation = gdk.INTERP_HYPER
// 		}
// 	}
// 	if outPixbuf, err = GetPixBuf(varPath); err == nil {
// 		if width != outPixbuf.GetWidth() || height != outPixbuf.GetHeight() {
// 			outPixbuf, err = outPixbuf.ScaleSimple(width, height, interpolation)
// 		}
// 	}
// 	return
// }

// // RotateImage: Rotate by 90,180,270 degres and get gdk.PixBuf image representation from file or []byte, depending on type
// func RotateImage(varPath interface{}, angle gdk.PixbufRotation) (outPixbuf *gdk.Pixbuf, err error) {
// 	if outPixbuf, err = GetPixBuf(varPath); err == nil {
// 		switch angle {
// 		case 90:
// 			outPixbuf, err = outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_COUNTERCLOCKWISE)
// 		case 180:
// 			outPixbuf, err = outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_UPSIDEDOWN)
// 		case 270:
// 			outPixbuf, err = outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_CLOCKWISE)
// 		default:
// 			return nil, errors.New("Rotation not allowed: " + fmt.Sprintf("%d", angle))
// 		}
// 	}
// 	return
// }

// // FlipImage: Get Flipped gdk.PixBuf image representation from file or []byte, depending on type
// func FlipImage(varPath interface{}, horizontal bool) (outPixbuf *gdk.Pixbuf, err error) {
// 	if outPixbuf, err = GetPixBuf(varPath); err == nil {
// 		outPixbuf, err = outPixbuf.Flip(horizontal)
// 	}
// 	return
// }

// /*
// 	Old functions, not used in new programs written and
// 	will be deleted after updating all other programs
// */
// // setImage: Set Image to GtkImage objects
// func SetImage(object *gtk.Image, varPath interface{}, size ...int) {
// 	if inPixbuf, err := GetPixBuf(varPath, size...); err == nil {
// 		object.SetFromPixbuf(inPixbuf)
// 		return
// 	} else if len(varPath.(string)) != 0 {
// 		fmt.Printf("SetImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setWinIcon: Set Icon to GtkWindow objects
// func SetWinIcon(object *gtk.Window, varPath interface{}, size ...int) {
// 	if inPixbuf, err := GetPixBuf(varPath, size...); err == nil {
// 		object.SetIcon(inPixbuf)
// 	} else if len(varPath.(string)) != 0 {
// 		fmt.Printf("SetWinIcon: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setButtonImage: Set Icon to GtkButton objects
// func SetButtonImage(object *gtk.Button, varPath interface{}, size ...int) {
// 	var image *gtk.Image
// 	inPixbuf, err := GetPixBuf(varPath, size...)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
// 			object.SetImage(image)
// 			object.SetAlwaysShowImage(true)
// 			return
// 		}
// 	}
// 	if err != nil && len(varPath.(string)) != 0 {
// 		fmt.Printf("SetButtonImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setToolButtonImage: Set Icon to GtkToolButton objects
// func SetToolButtonImage(object *gtk.ToolButton, varPath interface{}, size ...int) {
// 	var image *gtk.Image
// 	inPixbuf, err := GetPixBuf(varPath, size...)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
// 			object.SetIconWidget(image)
// 			return
// 		}
// 	}
// 	if err != nil && len(varPath.(string)) != 0 {
// 		fmt.Printf("setToolButtonImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setToggleButtonImage: Set Icon to GtkToggleButton objects
// func SetToggleButtonImage(object *gtk.ToggleButton, varPath interface{}, size ...int) {
// 	var image *gtk.Image
// 	inPixbuf, err := GetPixBuf(varPath, size...)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
// 			object.SetImage(image)
// 			object.SetAlwaysShowImage(true)
// 			return
// 		}
// 	}
// 	if err != nil && len(varPath.(string)) != 0 {
// 		fmt.Printf("SetToggleButtonImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // SetSpinButtonImage: Set Icon to GtkSpinButton objects. Position = "left" or "right"
// func SetSpinButtonImage(object *gtk.SpinButton, varPath interface{}, position ...string) {
// 	var inPixbuf *gdk.Pixbuf
// 	var err error
// 	pos := gtk.ENTRY_ICON_PRIMARY
// 	if len(position) > 0 {
// 		if position[0] == "right" {
// 			pos = gtk.ENTRY_ICON_SECONDARY
// 		}
// 	}
// 	if inPixbuf, err = GetPixBuf(varPath); err == nil {
// 		object.SetIconFromPixbuf(pos, inPixbuf)
// 		return
// 	} else if len(varPath.(string)) != 0 {
// 		fmt.Printf("SetSpinButtonImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setBoxImage:  Set Image to GtkBox objects
// func SetBoxImage(object *gtk.Box, varPath interface{}, size ...int) {
// 	var image *gtk.Image
// 	inPixbuf, err := GetPixBuf(varPath, size...)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
// 			image.Show()
// 			object.Add(image)
// 			return
// 		}
// 	}
// 	if err != nil && len(varPath.(string)) != 0 {
// 		fmt.Printf("setBoxImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }
