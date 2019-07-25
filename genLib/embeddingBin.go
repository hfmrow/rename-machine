// embeddingBin.go

/*
	©2019 H.F.M
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	The source-code below is derived from the work of:
	[github.com/jteeuwen/go-bindata], his work is subject to the CC0 1.0 Universal
	(CC0 1.0) Public Domain Dedication. http://creativecommons.org/publicdomain/zero/1.0/
	which I thank the author (jteeuwen) for his great work.
*/

package genLib

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

const lowerHex = "0123456789abcdef"

type StringWriter struct {
	io.Writer
	c int
}

type VarFile struct {
	VarName  string
	Filename string
}

// BinToHex: Convert binary file to gzipped []byte
func BinToHexString(filename string) (outString string, err error) {
	var byteToString = func(data []byte) (outString string) {
		var inByte byte
		buffer := []byte(`\x00`)
		for _, inByte = range data {
			buffer[2] = lowerHex[inByte/16]
			buffer[3] = lowerHex[inByte%16]
			outString += string(buffer)
		}
		return outString
	}
	fdIn, err := os.Open(filename)
	if err != nil {
		return outString, err
	}
	inf, err := os.Stat(filename)
	if err != nil {
		return outString, err
	}
	buff := make([]byte, inf.Size())
	_, err = fdIn.Read(buff)
	if err != nil {
		return outString, err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(buff)
	err = w.Close()
	if err != nil {
		return outString, err
	}
	return byteToString(b.Bytes()), nil
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

// BinFilesToHexFile: Convert binary files to gzipped []byte in specific file.
// Much faster than the previous version that only deals with one file at a time.
func BinFilesToHexFile(outFilename string, assets []VarFile, doBackup bool) (err error) {
	// Create output file.
	if doBackup {
		err = os.Remove(outFilename + "~")
		if err == nil {
			err = os.Rename(outFilename, outFilename+"~")
		}
	}
	fd, err := os.Create(outFilename)
	defer fd.Close()
	if err == nil {
		// Create a buffered writer for better performance.
		w := bufio.NewWriter(fd)
		defer w.Flush()
		if err == nil {
			for _, asset := range assets {
				_, err = fmt.Fprintf(w, `var _%s = HexToBytes("_%s", []byte("`, asset.VarName, asset.VarName)
				if err == nil {
					// Read asset content
					fd, err = os.Open(asset.Filename)
					if err == nil {
						// Compress and write
						gz := gzip.NewWriter(&StringWriter{Writer: w})
						_, err = io.Copy(gz, fd)
						gz.Close()
						if err != nil {
							return err
						}
						_, err = fmt.Fprintf(w, `"))
`)
					}
				}
			}
		}
	}
	return err
}

func (w *StringWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	buf := []byte(`\x00`)
	var b byte
	for n, b = range p {
		buf[2] = lowerHex[b/16]
		buf[3] = lowerHex[b%16]
		w.Writer.Write(buf)
		w.c++
	}
	n++
	return
}
