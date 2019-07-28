// gohFunctions.go

// Source file auto-generated on Sun, 28 Jul 2019 07:02:22 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

/*******************************************************/
/* Functions declarations, used to initialize objects */
/*****************************************************/
// newBuilder: initialise builder with glade xml string
func newBuilder(varPath interface{}) (err error) {
	if mainObjects.mainUiBuilder, err = gtk.BuilderNew(); err == nil {
		if Gtk3Interface, err := getBytesFromVarAsset(varPath); err == nil {
			err = mainObjects.mainUiBuilder.AddFromString(string(Gtk3Interface))
		}
	}
	return err
}

// loadObject: Load GtkObject to be transtyped ...
func loadObject(name string) (newObj glib.IObject) {
	var err error
	newObj, err = mainObjects.mainUiBuilder.GetObject(name)
	if err != nil {
		log.Panic(err)
	}
	return newObj
}

// WindowDestroy: is the triggered handler when closing/destroying the gui window.
func windowDestroy() {
	// Doing something before quit.

	// Restoring previous state of preserve extension in cas of directory filename
	visible := mainObjects.SinglePresExtChk.GetVisible()
	if !visible {
		mainOptions.PreserveExtSingle = bakPresExt
	}

	err = mainOptions.Write() /* Update mainOptions with values of gtk conctrols and write to file */
	if err != nil {
		fmt.Printf("%s\n%v\n", "Writing options error.", err)
	}
	// Bye ...
	gtk.MainQuit()
}

/*************************************************/
/* Images functions, used to initialize objects */
/***********************************************/
// setWinIcon: Set Icon to GtkWindow objects
func setWinIcon(object *gtk.Window, varPath interface{}) {
	if inPixbuf, err := getPixBuff(varPath); err == nil {
		object.SetIcon(inPixbuf)
	} else {
		if len(varPath.(string)) != 0 {
			fmt.Printf("An error occurred on image: %s\n%v\n", varPath, err.Error())
		}
	}
}

// setImage: Set Image to GtkImage objects
func setImage(object *gtk.Image, varPath interface{}) {
	if inPixbuf, err := getPixBuff(varPath); err == nil {
		object.SetFromPixbuf(inPixbuf)
		return
	} else {
		if len(varPath.(string)) != 0 {
			fmt.Printf("An error occurred on image: %s\n%v\n", varPath, err.Error())
		}
	}
	object.Hide()
}

// setButtonImage: Set Icon to GtkButton objects
func setButtonImage(object *gtk.Button, varPath interface{}) {
	if inPixbuf, err := getPixBuff(varPath); err == nil {
		if image, err := gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
			object.SetImage(image)
			object.SetAlwaysShowImage(true)
			return
		}
	}
	if err != nil {
		if len(varPath.(string)) != 0 {
			fmt.Printf("An error occurred on image: %s\n%v\n", varPath, err.Error())
		}
	}
}

// getPixBuff: Get gtk.Pixbuff image representation from file or []byte, depending on type
func getPixBuff(varPath interface{}) (outPixbuf *gdk.Pixbuf, err error) {
	switch reflect.TypeOf(varPath).String() {
	case "string":
		outPixbuf, err = gdk.PixbufNewFromFile(varPath.(string))
	case "[]uint8":
		pbLoader, err := gdk.PixbufLoaderNew()
		if err == nil {
			outPixbuf, err = pbLoader.WriteAndReturnPixbuf(varPath.([]byte))
		}
	}
	return outPixbuf, err
}

/***************************************/
/* Embedded data conversion functions */
/* Used to make variable content     */
/* available in go-source           */
/***********************************/
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
	if r, err := gzip.NewReader(bytes.NewBuffer(gzipData)); err == nil {
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
	return outByte
}

/*******************************/
/* Simplified files Functions */
/*****************************/
// tempMake: Make temporary directory
func tempMake(prefix string) (dir string) {
	if dir, err = ioutil.TempDir("", prefix+"-"); err != nil {
		log.Fatal(err)
	}
	return dir + string(os.PathSeparator)
}

// getAbsRealPath: Retrieve app current realpath and options filenme
func getAbsRealPath() (absoluteRealPath, optFilename string) {
	absoluteBaseName, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	base := filepath.Base(absoluteBaseName)
	splited := strings.Split(base, ".")
	length := len(splited)
	if length == 1 {
		optFilename = base
	} else {
		splited = splited[:length-1]
		optFilename = strings.Join(splited, ".")
	}
	return filepath.Dir(absoluteBaseName) + string(os.PathSeparator),
		filepath.Join(filepath.Dir(absoluteBaseName), optFilename+".opt")
}

// Used as fake function for signals section
func blankNotify() {}
