// scanDir.go

package genLib

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	times "gopkg.in/djherbis/times.v1"
)

// Scan dir and subdir to get symlinks with specified endpoint.
func RecurseScanSymlink(path, linkEndPoint string) (fileList []string, err error) {
	err = filepath.Walk(path, func(filePath string, fileInfos os.FileInfo, err error) error {
		if fileInfos.Mode()&os.ModeSymlink != 0 {
			realPath, err := os.Readlink(filePath)
			if err != nil {
				return err
			}
			if strings.Contains(realPath, linkEndPoint) {
				fileList = append(fileList, strings.Replace(filePath, path, "", -1))
			}
		}
		return nil
	})
	if err != nil {
		return fileList, err
	}
	return fileList, nil
}

// Used in ScanFiles function to store file informations.
type fileInfos struct {
	IsExists         bool
	PathBase         string
	Base             string
	Path             string
	Ext              string
	Mtime            time.Time
	Atime            time.Time
	MtimeYMDhm       string
	AtimeYMDhm       string
	MtimeYMDhms      string
	AtimeYMDhms      string
	MtimeYMDhmsShort string
	AtimeYMDhmsShort string
	MtimeFriendlyHR  string
	AtimeFriendlyHR  string
	Type             string
	Size             int64
	SizeHR           string
}

// ScanFiles: Scan given files and retreive informations about them stored in a []fileInfos structure.
func ScanFiles(inFiles []string) (outFiles []fileInfos) {
	for _, file := range inFiles {
		outFiles = append(outFiles, ScanFile(file))
	}
	return outFiles
}

// ScanFile: Scan a file and retreive informations about it stored in a fileInfos structure.
func ScanFile(file string) (fi fileInfos) {
	var tmpStr string
	fi.PathBase = file

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fi.IsExists = false
		return fi
	}
	fi.IsExists = true

	infos, err := os.Stat(file)
	Check(err)

	switch {
	case (infos.Mode()&os.ModeDir != 0):
		fi.Type = "Dir"
	case (infos.Mode()&os.ModeSymlink != 0):
		fi.Type = "Link"
	case (infos.Mode()&os.ModeAppend != 0):
		fi.Type = "Append" // a: append-only
	case (infos.Mode()&os.ModeExclusive != 0):
		fi.Type = "Exclusive" // l: exclusive use
	case (infos.Mode()&os.ModeTemporary != 0):
		fi.Type = "Temp" // T: temporary file; Plan 9 only
	case (infos.Mode()&os.ModeDevice != 0):
		fi.Type = "Device" // D: device file
	case (infos.Mode()&os.ModeNamedPipe != 0):
		fi.Type = "Pipe" // p: named pipe (FIFO)
	case (infos.Mode()&os.ModeSocket != 0):
		fi.Type = "Socket" // S: Unix domain socket
	case (infos.Mode()&os.ModeCharDevice != 0):
		fi.Type = "CharDev" // c: Unix character device, when ModeDevice is set
	case (infos.Mode()&os.ModeSticky != 0):
		fi.Type = "Sticky" // t: sticky
	case (infos.Mode()&os.ModeIrregular != 0):
		fi.Type = "Unknowen" // ?: non-regular file; nothing else is known about this file
	default:
		fi.Type = "File"
	}

	if fi.Type == "Dir" {
		fi.Base = filepath.Base(file)
		fi.Path = file
		fi.Ext = "Dir"
	} else {
		fi.Base = filepath.Base(file)
		fi.Path = filepath.Dir(file)
		fi.Ext = filepath.Ext(file)
	}

	fi.Size = infos.Size()
	fi.SizeHR = humanize.Bytes(uint64(fi.Size))

	fi.Atime = times.Get(infos).AccessTime()
	fi.Mtime = times.Get(infos).ModTime()
	fi.MtimeFriendlyHR = humanize.Time(fi.Mtime)
	fi.AtimeFriendlyHR = humanize.Time(fi.Atime)
	fi.MtimeYMDhm = fi.Mtime.String()[:16]
	fi.AtimeYMDhm = fi.Atime.String()[:16]
	fi.MtimeYMDhms = fi.Mtime.String()[:19]
	fi.AtimeYMDhms = fi.Atime.String()[:19]
	tmpStr = RemoveNonAlNum(fi.MtimeYMDhms)
	fi.MtimeYMDhmsShort = tmpStr[2:len(tmpStr)]
	tmpStr = RemoveNonAlNum(fi.AtimeYMDhms)
	fi.AtimeYMDhmsShort = tmpStr[2:len(tmpStr)]

	return fi
}

// ScanSubDir retrieve files in a specific directory and his sub-directory
// depending on depth argument. depth = -1 mean infinite.
func ScanDirDepth(root string, depth int, showDirs ...bool) (files []string, err error) {
	var listDir bool
	var depthRecurse int
	var tmpFiles []string
	if len(showDirs) != 0 {
		listDir = showDirs[0]
	}
	// osSep := string(os.PathSeparator)
	// root = strings.TrimSuffix(root, osSep)
	filesInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, file := range filesInfos {
		depthRecurse = depth
		if !file.IsDir() {
			files = append(files, filepath.Join(root, file.Name()))
		} else {
			if depth != 0 {
				depthRecurse--
				tmpFiles, err = ScanDirDepth(filepath.Join(root, file.Name()), depthRecurse, listDir)
			}
			if err != nil {
				return files, err
			}
			files = append(files, tmpFiles...)
			if listDir {
				files = append(files, filepath.Join(root, file.Name()))
			}
		}
	}
	return files, err
}

// ScanSubDir retrieve files in a specific directory and his sub-directory.
// don't follow symlink (walk)
func ScanSubDir(root string, showDirs ...bool) (files []string, err error) {
	var listDir bool
	if len(showDirs) != 0 {
		listDir = showDirs[0]
	}
	err = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			} else if listDir {
				files = append(files, path)
			}
			return nil
		})
	return files, err
}

// ScanDir retrieve files in a specific directory
func ScanDir(root string, showDirs ...bool) (files []string, err error) {
	var listDir bool
	if len(showDirs) != 0 {
		listDir = showDirs[0]
	}
	// osSep := string(os.PathSeparator)
	// root = strings.TrimSuffix(root, osSep)
	filesInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, file := range filesInfos {
		if !file.IsDir() {
			files = append(files, filepath.Join(root, file.Name()))
		} else if listDir {
			files = append(files, filepath.Join(root, file.Name()))
		}
	}
	return files, err
}
