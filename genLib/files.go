// files.go

package genLib

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// ByteCountDecimal: Format byte size to human readable format
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// ByteCountBinary: Format byte size to human readable format
func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// goFormatting: Format output file using gofmt function.
func GoFormatting(filename string) {
	output, err := exec.Command("gofmt", "-w", filename).CombinedOutput()
	if err != nil {
		Check(err, `gofmt ERROR !.`, string(output))
	} else {
		fmt.Println("Output file formated: ", filename)
	}
}

// magic nulber mime detection
var magicTable = map[string]string{
	"\x37\x7A\xBC\xAF\x27\x1C\x00\x04": "7zip",
	"\xFD\x37\x7A\x58\x5A\x00\x00":     "xz",
	"\x1F\x8B\x08\x00\x00\x09\x6E\x88": "gzip",
	"\x75\x73\x74\x61\x72":             "tar",
}

// GetFileMime: scan first bytes to detect mime type of file.
func GetFileMime(filename string) string {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		buffReader := bufio.NewReader(file)
		for magic, mime := range magicTable {
			if peeked, err := buffReader.Peek(len([]byte(magic))); err == nil {
				tmpMagic := []byte(magic)
				if bytes.Index(peeked, tmpMagic) == 0 {
					return mime
				}
			}
		}
	}
	return "unknowen"
}

// CheckMime: return the mime type of a file
func CheckMime(filename string) (mime string, err error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return mime, err
	}
	buff := make([]byte, 512)
	if _, err = file.Read(buff); err != nil {
		return mime, err
	}
	return http.DetectContentType(buff), err
}

// IsTextFile: check for first 512 bytes if they contains some bytes usually present
// in utf-32, utf-16 or ascii/utf-8 files. Return detected type, including binary.
func IsTextFile(filename string) (fileType string, err error) {
	var threshold = int(8)
	var length = int(512)
	var lf = []byte{0x0A}
	var cr = []byte{0x0D}

	// Look at data if contain bytes from 0x01 to 0x06 (usually not present in text files)
	var fCheckForbbiden = func(dt []byte) bool {
		for idx := 1; idx < 7; idx++ {
			if bytes.Contains(dt, []byte{byte(idx)}) {
				return false
			}
		}
		return true
	}
	// Look at data if contain bytes of utf32 + line-end
	var fCheck32 = func(dt, chk []byte) bool {
		return bytes.Contains(dt, []byte{0x00, 0x00, 0x00, chk[0]})
	}
	// Look at data if contain bytes of utf16 + line-end
	var fCheck16 = func(dt, chk []byte) bool {
		return bytes.Contains(dt, []byte{0x00, chk[0]})
	}

	// File operations
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return "Error opening file: " + file.Name(), err
	}
	stat, err := file.Stat()
	if err != nil {
		return "Error stat file: " + file.Name(), err
	}
	size := stat.Size()
	if size >= int64(threshold) {
		// Read datas
		data := make([]byte, length)
		length, err = file.Read(data)
		if err != nil {
			return "Error reading file: " + file.Name(), err
		}
		switch {
		case fCheck32(data, cr) && fCheckForbbiden(data):
			return "utf-32", nil
		case fCheck32(data, lf) && fCheckForbbiden(data):
			return "utf-32", nil

		case fCheck16(data, cr) && fCheckForbbiden(data):
			return "utf-16", nil
		case fCheck16(data, lf) && fCheckForbbiden(data):
			return "utf-16", nil

		case bytes.Contains(data, cr) && fCheckForbbiden(data):
			return "ascii/utf-8", nil
		case bytes.Contains(data, lf) && fCheckForbbiden(data):
			return "ascii/utf-8", nil

		case fCheckForbbiden(data):
			return "ascii/utf-8", nil // At this point, no line end exist.
		default:
			fileType = "binary"
			return fileType, errors.New(fileType)
		}
	}
	fileType = "File size < " + strconv.Itoa(threshold) + " bytes"
	return fileType, errors.New(fileType)
}

// CheckCmd: Check for command if exist
func CheckCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}
	return true
}

// FileExist: reports whether the named file or directory exists.
func FileExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDirEmpty:
func IsDirEmpty(name string) (empty bool, err error) {
	var f os.File
	if f, err := os.Open(name); err == nil {
		if _, err = f.Readdirnames(1); err == io.EOF {
			return true, nil
		}
	}
	defer f.Close()
	return false, err
}

// GetCurrentDir: Get current directory
func GetCurrentDir() (dir string, err error) {
	return os.Getwd()
}

// GetOsPathSep: Get OS path separator
func GetOsPathSep() string {
	return string(os.PathSeparator)
}

// TempMake: Make temporary directory
func TempMake(prefix string) string {
	dir, err := ioutil.TempDir("", prefix+"-")
	Check(err)
	return dir + string(os.PathSeparator)
}

// TempRemove: Remove directory recursively
func TempRemove(fName string) (err error) {
	if err = os.RemoveAll(fName); err != nil {
		return (err)
	}
	return nil
}

// GetFileEOL: Open file and get (CR, LF, CRLF) > string or get OS line end.
func GetFileEOL(filename string) (outString string, err error) {
	textFileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return outString, err
	}
	return GetTextEOL(textFileBytes), nil
}

// SetFileEOL: Open file and convert EOL (CR, LF, CRLF) then write it back.
func SetFileEOL(filename, eol string) error {
	textFileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// Handle end of line
	textFileBytes, err = SetTextEOL(textFileBytes, eol)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, textFileBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ExtEnsure: ensure the filename have desired extension
func ExtEnsure(filename, ext string) (outFilename string) {
	outFilename = filename
	if !strings.HasSuffix(filename, ext) {
		currExt := path.Ext(filename)
		outFilename = filename[:len(filename)-len(currExt)] + ext
	}
	return outFilename
}

// BaseNoExt: get only the name without ext.
func BaseNoExt(filename string) (outFilename string) {
	outFilename = filepath.Base(filename)
	ext := filepath.Ext(outFilename)
	return outFilename[:len(outFilename)-len(ext)]
}

// CopyFile:
func CopyFile(inFile, outFile string, doBackup ...bool) (err error) {
	var inBytes []byte
	if inBytes, err = ioutil.ReadFile(inFile); err == nil {
		if len(doBackup) != 0 {
			if doBackup[0] {
				if err = os.Rename(outFile, outFile+"~"); err != nil {
					return err
				}
			}
		}
		err = ioutil.WriteFile(outFile, inBytes, 0644)
	}
	return err
}

// renameProjectFiles: Mass rename function
func renameListFiles(fromFileList, toFileList []string) (err error) {
	for idx, file := range fromFileList {
		if file != toFileList[idx] {
			if err = os.Rename(file, toFileList[idx]); err != nil {
				return err
			}
		}
	}
	return err
}

// ReadFile:
func ReadFile(filename string) (data []byte, err error) {
	return ioutil.ReadFile(filename)
}

// writeFile: with file backup capability
func WriteFile(filename string, datas []byte, doBackup ...bool) (err error) {
	if len(doBackup) != 0 {
		if doBackup[0] {
			if _, err = os.Stat(filename); !os.IsNotExist(err) {
				if err = os.Rename(filename, filename+"~"); err != nil {
					return err
				}
			}
		}
	}
	return ioutil.WriteFile(filename, datas, os.ModePerm)
}

// replaceInFile: allow using regexp in argument.
func replaceInFile(filename, search, replace string, doBackup ...bool) (found bool, err error) {
	var findReg = regexp.MustCompile(`(` + regexp.QuoteMeta(search) + `)`)
	var data []byte
	data, err = ioutil.ReadFile(filename)
	if err == nil {
		if len(doBackup) != 0 {
			if doBackup[0] {
				if err = os.Rename(filename, filename+"~"); err != nil {
					return found, err
				}
			}
		}
		if !findReg.Match(data) {
			return false, err
		}
		found = true
		err = ioutil.WriteFile(filename, findReg.ReplaceAll(data, []byte(replace)), 0644)
	}
	return found, err
}

// GenericReader: return io.reader for standar files and .bz2, gz.
func GenericReader(filename string) (io.Reader, *os.File, error) {
	if filename == "" {
		return bufio.NewReader(os.Stdin), nil, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasSuffix(filename, "bz2") {
		return bufio.NewReader(bzip2.NewReader(bufio.NewReader(file))), file, err
	}

	if strings.HasSuffix(filename, "gz") {
		reader, err := gzip.NewReader(bufio.NewReader(file))
		if err != nil {
			return nil, nil, err
		}
		return bufio.NewReader(reader), file, err
	}
	return bufio.NewReader(file), file, err
}

// FindDir retrieve file in a specific directory with more options.
func FindDir(dir, mask string, returnedStrSlice *[][]string, scanSub, showDir, followSymlinkDir bool) error {
	var fName, time, size string
	// Remove unwanted os path separator if exist
	//	dir = strings.TrimSuffix(dir, string(os.PathSeparator))

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		fName = filepath.Join(dir, file.Name())
		if followSymlinkDir { // Check for symlink ..
			file, err = os.Lstat(fName)
			if err != nil {
				return err
			}
			if file.Mode()&os.ModeSymlink != 0 { // Is a symlink ?
				fName, err := os.Readlink(fName) // Then read it...
				if err != nil {
					return err
				}
				file, err = os.Stat(fName) // Get symlink infos.
				if err != nil {
					return err
				}
				fName = filepath.Join(dir, file.Name())
			}
		}
		// Recursive play if it's a directory
		if file.IsDir() && scanSub {
			tmpFileList := new([][]string)
			err = FindDir(fName, mask, tmpFileList, scanSub, true, followSymlinkDir)
			*returnedStrSlice = append(*returnedStrSlice, *tmpFileList...)
			if err != nil {
				return err
			}
		}
		// get information to be displayed.
		size = fmt.Sprintf("%s", humanize.Bytes(uint64(file.Size())))
		time = fmt.Sprintf("%s.", humanize.Time(file.ModTime()))
		// Check for ext matching well.
		ok, err := filepath.Match(mask, file.Name())
		if err != nil {
			return err
		}
		if ok {
			if showDir { // Limit display directories if requested
				*returnedStrSlice = append(*returnedStrSlice, []string{file.Name(), size, time, fName})
			} else {
				if !file.IsDir() {
					*returnedStrSlice = append(*returnedStrSlice, []string{file.Name(), size, time, fName})
				}
			}
		}
	}
	return nil
}

// File struct (SplitFilePath)
type Filepath struct {
	Absolute             string
	Relative             string
	Path                 string
	Base                 string
	BaseNoExt            string
	Ext                  string
	ExecFullName         string
	RealPath             string
	RealName             string
	OutputNewExt         string
	OutputAppendFilename string
	OsSeparator          string
	IsDir                bool
	SymLink              bool
	SymLinkTo            string
}

// Split full filename into path, ext, name, ... optionally add suffix before original extension or change extension
// Relative: SplitFilepath("wanted relative path", fullpath).Relative
// Absolute: SplitFilepath("relative path", fullpath).Absolute
func SplitFilepath(filename string, newExt ...string) Filepath {
	var dir, link bool
	var f = Filepath{}
	var newExtension, dot, addToFilename string
	if len(newExt) != 0 {
		addToFilename = newExt[0]
		if !strings.Contains(newExt[0], ".") {
			dot = "."
		}
		newExtension = dot + newExt[0]
	}
	// IsDir
	fileInfos, err := os.Lstat(filename)
	if err == nil {
		dir = (fileInfos.Mode()&os.ModeDir != 0)
		link = (fileInfos.Mode()&os.ModeSymlink != 0)
		f.IsDir = dir
		if link {
			// IsLink
			f.SymLink = link
			// Symlink endpoint
			f.SymLinkTo, _ = os.Readlink(filename)
			// Symlink and Directory
			ls, err := os.Lstat(f.SymLinkTo)
			if err == nil {
				f.IsDir = (ls.Mode()&os.ModeDir != 0)
			}
		}
	}
	// Absolute
	f.Absolute, _ = filepath.Abs(filename)
	// Relative - Use the optional argument to set as basepath ...
	f.Relative, _ = filepath.Rel(newExtension, filename)
	// OsSep
	f.OsSeparator = string(os.PathSeparator)
	// Path
	if f.Path = filepath.Dir(filename); f.Path == "." {
		f.Path = ""
	}
	// Base
	f.Base = filepath.Base(filename)
	// Ext
	f.Ext = filepath.Ext(filename)
	// BaseNoExt
	splited := strings.Split(f.Base, ".")
	length := len(splited)
	if length == 1 {
		f.BaseNoExt = f.Base

	} else {
		if f.Base[:1] == "." { // Case of hidden file starting with dot
			f.Ext = ""
			f.BaseNoExt = f.Base
		} else {
			splited = splited[:length-1]
			f.BaseNoExt = strings.Join(splited, ".")
		}
	}
	// ExecFullName
	f.ExecFullName, _ = os.Executable()
	// RealPath
	realPathName, _ := filepath.EvalSymlinks(filename)
	if f.RealPath = filepath.Dir(realPathName); f.RealPath == "." {
		f.RealPath = ""
	}
	// RealName
	if f.RealName = filepath.Base(realPathName); f.RealName == "." {
		f.RealName = ""
	}
	// OutNewExt
	if f.Path == "" {
		f.OutputNewExt = f.BaseNoExt + newExtension
	} else {
		f.OutputNewExt = f.Path + f.OsSeparator + f.BaseNoExt + newExtension
	}
	// OutputAppendFilename
	if f.Path == "" {
		f.OutputAppendFilename = f.BaseNoExt + addToFilename + f.Ext
	} else {
		f.OutputAppendFilename = f.Path + f.OsSeparator + f.BaseNoExt + addToFilename + f.Ext
	}
	return f
}

// GetFileBytesString: Retrieve 'from' 'to' bytes from file in string format.
func GetFileBytesString(filename string, from, to int) (outString string) {
	var WriteBytesString = func(p []byte) (data string) {
		const lowerHex = "0123456789abcdef"
		if len(p) == 0 {
			return data
		}
		buf := []byte(`\x00`)
		var b byte
		for _, b = range p {
			buf[2] = lowerHex[b/16]
			buf[3] = lowerHex[b%16]
			data += string(buf)
		}
		return data
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	buff := make([]byte, to-from)
	_, err = file.ReadAt(buff, int64(from))
	if err != nil {
		fmt.Println(err)
	}
	return WriteBytesString(buff)
}

/* Const for FileMode
Usage to Create any directories needed to put this file in them:
     var dir_file_mode os.FileMode
     dir_file_mode = os.ModeDir | (OS_USER_RWX | OS_ALL_R)
     os.MkdirAll(dir_str, dir_file_mode)

	fmt.Printf("%o\n%o\n%o\n%o\n",
		os.ModePerm&OS_ALL_RWX,
		os.ModePerm&OS_USER_RW|OS_GROUP_R|OS_OTH_R,
		os.ModePerm&OS_USER_RW|OS_GROUP_RW|OS_OTH_R,
		os.ModePerm&OS_USER_RWX|OS_GROUP_RWX|OS_OTH_R)
*/
const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_ALL_X
)

// DispRights: display right.
//i.e: g.DispRights(g.OS_USER_RWX | g.OS_GROUP_RWX | g.OS_OTH_R)
func DispRights(value int) {
	fmt.Printf("%o\n", value)
}
