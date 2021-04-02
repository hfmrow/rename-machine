// log2File.go

/*
	Copyright ©2020 H.F.M - A Log2File package v1.0 github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package log2file

import (
	"fmt"
	"log"
	"os"
	"strings"

	glfs "github.com/hfmrow/genLib/files"
	gltses "github.com/hfmrow/genLib/tools/errors"
)

// Log2FileStruct: Structure that hold methods to manage logging to file
// (errors and/or strings). On main program exit, logger must be
// closed using "CloseLogger" method: defer Logger.CloseLogger()
type Log2FileStruct struct {
	// Skip previous caller, usually the 1st is the good
	CallerSkip int

	DevMode bool
	logger  *log.Logger
	file    *os.File
	LogFilename,
	startLog,
	endLog string
}

// Log: Record "messageIn" to file, can be an error and/or a string
// NOTICE: first argument must be the error to speed up treatment.
// return a cleaned 'nil' error
func (l2f *Log2FileStruct) Log(messageIn ...interface{}) (err error) {
	var (
		outMess,
		outErrMess string
		errMess []string
	)
	for _, elem := range messageIn {
		switch elem.(type) {
		case nil: // No error, go back
			return
		case error:
			errMess = append(errMess, elem.(error).Error())
		default:
			if txt := fmt.Sprintf("%v", elem); len(txt) != 0 {
				outMess += "[" + txt + "]"
			}
		}
	}
	caller := gltses.WarnCallerMess(l2f.CallerSkip)
	outErrMess = strings.Join(errMess, " ")
	if l2f.DevMode {
		// Output to console if devmode is true
		fmt.Printf("%s %s %s\n", caller, outMess, outErrMess)
	}
	// Output to file
	l2f.logger.Printf("%s %s %s\n", caller, outMess, outErrMess)
	return
}

// CloseLogger: Whan main program EXIT this method must to be called
// to properly close lo file.
func (l2f *Log2FileStruct) CloseLogger() {
	if err := l2f.file.Close(); err != nil {
		log.Println(err)
	}
}

// Log2FileNew: Create a new Log2File structure
func Log2FileStructNew(logFilename string, devMode bool) (l2f *Log2FileStruct) {
	var err error

	l2f = new(Log2FileStruct)
	l2f.DevMode = devMode
	// l2f.CallerSkip = 0

	l2f.LogFilename = glfs.ExtEnsure(logFilename, ".log")

	if l2f.file, err = os.OpenFile(l2f.LogFilename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {

		l2f.logger = log.New(l2f.file, "["+glfs.BaseNoExt(l2f.LogFilename)+"] ", log.LstdFlags)
	}
	return
}
