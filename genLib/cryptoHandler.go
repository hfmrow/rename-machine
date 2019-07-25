// cryptoHandler.go

/// +build OMIT
package genLib

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// Get MD5 checksum from file.
func Md5Sum(filename string) string {
	f, err := os.Open(filename)
	defer f.Close()
	Check(err, "os.Open!")

	h := md5.New()
	_, err = io.Copy(h, f)
	Check(err, "io.Copy!")

	return fmt.Sprintf("%x", h.Sum(nil))
}

// Get MD5 checksum from string.
func Md5String(inString string) string {
	data := []byte(inString)
	return fmt.Sprintf("%x", md5.Sum(data))
}
