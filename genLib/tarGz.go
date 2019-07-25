// tarGz.go

// Source file auto-generated on Sun, 31 Mar 2019 19:42:54 using Gotk3ObjHandler v1.3 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

// CreateTarball: create tar.gz file from given filenames.
// UntarGzip: unpack tar.gz files list. Return len(filesList)=0 if all files has been restored.

package genLib

import (
	"archive/tar"
	"compress/flate"
	"sort"

	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	gzip "github.com/klauspost/pgzip"
	"github.com/ulikunitz/xz"
)

var countFiles int

// CreateTarball: create tar.gz file from given filenames.
// -2 = HuffmanOnly (linear compression, low gain, fast compression)
// -1 = DefaultCompression
//  0 = NoCompression
//  1 -> 9 = BestSpeed -> BestCompression
func CreateTarballGzip(tarballFilePath string, filePaths []string, compressLvl int) (countedWrittenFiles int, err error) {
	countFiles = 0
	file, err := os.Create(tarballFilePath)
	if err != nil {
		return countFiles, errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarballFilePath, err.Error()))
	}
	defer file.Close()

	zipWriter, err := gzip.NewWriterLevel(file, compressLvl)
	if err != nil {
		return countFiles, errors.New(fmt.Sprintf("Bad compression level '%d', got error '%s'", int(flate.BestCompression), err.Error()))
	}
	defer zipWriter.Close()

	tarWriter := tar.NewWriter(zipWriter)
	defer tarWriter.Close()
	// currentStoreFiles.inUse = true
	for _, filePath := range filePaths {
		// select {
		// case <-quitGoRoutine:
		// quitGoRoutine = make(chan struct{})
		// err = changeFileOwner(tarballFilePath)
		// currentStoreFiles.inUse = false
		// return countFiles, errors.New(sts["userCancelled"])
		// default:
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			return countFiles, errors.New(fmt.Sprintf("Could not add file '%s', to tarball, got error '%s'", filePath, err.Error()))
		}
		// }
	}
	err = changeFileOwner(tarballFilePath)
	// currentStoreFiles.inUse = false
	return countFiles, err
}

// addFileToTarWriter:
func addFileToTarWriter(filePath string, tarWriter *tar.Writer) (err error) {
	var stat os.FileInfo
	var linkname string

	stat, err = os.Lstat(filePath)
	modeType := stat.Mode() & os.ModeType
	switch {
	case modeType&os.ModeNamedPipe > 0:
		return nil
	case modeType&os.ModeSocket > 0:
		return nil
	case modeType&os.ModeDevice > 0:
		return nil
	case err != nil:
		return nil
	}

	file, err := os.Open(filePath)
	switch {
	case os.IsPermission(err):
		return nil
	default:
		defer file.Close()
	}

	if link, err := os.Readlink(filePath); err == nil {
		linkname = link
	}

	header, err := tar.FileInfoHeader(stat, filepath.ToSlash(linkname))
	if err != nil {
		return errors.New(fmt.Sprintf("Could not build header for '%s', got error '%s'", filePath, err.Error()))
	}
	header.Name = filePath

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}
	countFiles++
	if header.Typeflag == tar.TypeReg {
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
		}
	}
	return nil
}

// UntarGzip: unpack tar.gz files list. Return len(filesList)=0 if all files has been restored.
// removeDir mean that, before restore folder content, the existing dir will be removed.
func UntarGzip(sourcefile, dst string, filesList *[]string, removeDir bool) (countedWrittenFiles int, err error) {
	var storePath, target string
	var header *tar.Header
	var tmpFilesList = *filesList
	var skipReadHeader bool
	// Initialise readers
	file, err := os.Open(sourcefile)
	if err != nil {
		return countedWrittenFiles, err
	}
	defer file.Close()

	var r io.ReadCloser = file

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return countedWrittenFiles, err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	sort.SliceStable(tmpFilesList, func(i, j int) bool {
		return tmpFilesList[i] < tmpFilesList[j]
	})

	// ownThis: set file/dir owner owner stored in tar archive.
	var ownThis = func(target string) (err error) {
		err = os.Chown(target, header.Uid, header.Gid)
		if !os.IsPermission(err) {
			if err != nil {
				return err
			}
		}
		return nil
	}
	// buildDirStruct: and set owner respective for each dir created.
	var buildDirStruct = func(target string) (err error) {
		if _, err := os.Stat(target); os.IsNotExist(err) {
			if err := os.MkdirAll(target, os.ModePerm); err == nil {
				for {
					if err = ownThis(target); err != nil {
						return err
					}
					target = filepath.Dir(target)
					if target == filepath.Dir(dst) {
						break
					}
				}
			}
		}
		return err
	}
	// restoreFile: Handling Directory, regular file, symlink and own them. Others
	// kind of files are ignored cause my packing function don't handle them.
	var restoreFile = func(target string) (err error) {
		countedWrittenFiles++
		switch header.Typeflag {
		case tar.TypeDir:
			if removeDir {
				if _, err := os.Stat(target); !os.IsNotExist(err) {
					err = os.RemoveAll(target)
					if err != nil {
						return err
					}
				}
			}
			err = buildDirStruct(target)
			if err != nil {
				return err
			}
			// fmt.Println(header.Name)
		case tar.TypeSymlink:
			if _, err := os.Lstat(target); err == nil {
				os.Remove(target)
			}
			err = os.Symlink(header.Linkname, target)
			if err != nil {
				return err
			}
			// fmt.Println(header.Name)
		case tar.TypeReg:
			err = buildDirStruct(filepath.Dir(target)) // in case of regular file where there is no directory to receive it.
			if err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			if err = ownThis(target); err == nil {
				f.Close()
				return err
			}
			f.Close()
			// fmt.Println(header.Name)
		}
		return err
	}
	// readHeader:
	var readHeader = func() (target string, err error) {
		header, err = tr.Next()
		switch {
		case err == io.EOF:
			return "", nil
		case err != nil:
			return "", nil
		}
		target = filepath.Join(dst, header.Name)
		return target, err
	}
	// writeFile: and dir, set perms and owner stored in tar archive.
	var writeFile = func(target *string) (err error) {
		if header != nil {
			if header.Typeflag == tar.TypeDir {
				storePath = header.Name
				for {
					// select {
					// case <-quitGoRoutine:
					// 	quitGoRoutine = make(chan struct{})
					// 	// currentStoreFiles.inUse = false
					// 	header = nil
					// 	return errors.New(sts["userCancelled"])
					// default:
					err = restoreFile(*target)
					if err != nil {
						return err
					}
					*target, err = readHeader()
					if err != nil {
						return err
					}
					if header == nil || !strings.Contains(filepath.Dir(header.Name), storePath) {
						return nil
					}
					// }
				}
			} else {
				err = restoreFile(*target)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Parse desired files and restore them.
	// currentStoreFiles.inUse = true
	for {
		if !skipReadHeader {
			target, err = readHeader()
		}
		if err != nil || header == nil || len(target) == 0 {
			return countedWrittenFiles, err
		}
		if len(tmpFilesList) != 0 {
			for idx := 0; idx < len(tmpFilesList); idx++ {
				source := filepath.Join(dst, tmpFilesList[idx])
				if source == target && len(tmpFilesList) != 0 {
					tmpFilesList = append(tmpFilesList[:idx], tmpFilesList[idx+1:]...)
					idx--
					if err = writeFile(&target); err != nil {
						return countedWrittenFiles, err
					}
					if len(tmpFilesList) == 0 {
						*filesList = tmpFilesList
						return countedWrittenFiles, err
					}
				}
				skipReadHeader = false
				if idx == -1 || filepath.Join(dst, tmpFilesList[idx]) == target {
					skipReadHeader = true
				}
			}
		} else {
			if err = writeFile(&target); err != nil {
				return countedWrittenFiles, err
			}
		}
	}
	// End of game ...
	*filesList = tmpFilesList
	// currentStoreFiles.inUse = false
	return countedWrittenFiles, err
}

// UnGzip: Unpack
func UnGzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

// Gzip: Pack file
func Gzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(source)
	//	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

// Tar: Make standalone tarball
func Tar(source, target string) error {
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Lstat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		if baseDir != "" {
			//			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			header.Name = path
		}

		if link, err := os.Readlink(path); err == nil {
			header.Linkname = link
		}

		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if header.Typeflag == tar.TypeReg {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
		}
		return err
	})
}

// Untar: extract file from tarball
func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open tar.gz file '%s', got error '%s'", tarball, err.Error()))
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.New(fmt.Sprintf("Could not read file in tar.gz archive, got error '%s'", err.Error()))
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTarballLXz:  create tar.xz file from given filenames. Just for fun test ...
// but it's really slow, no multi-threading and low size gain.
func CreateTarballXz(tarballFilePath string, filePaths []string, compressLvl int) (countedWrittenFiles int, err error) {
	countFiles = 0
	file, err := os.Create(tarballFilePath)
	if err != nil {
		return countFiles, errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarballFilePath, err.Error()))
	}
	defer file.Close()

	xzWriter, err := xz.NewWriter(file)
	if err != nil {
		return countFiles, errors.New(fmt.Sprintf("Could not create xzWriter, got error '%s'", err.Error()))
	}
	defer xzWriter.Close()
	tarWriter := tar.NewWriter(xzWriter)
	defer tarWriter.Close()
	// currentStoreFiles.inUse = true
	for _, filePath := range filePaths {
		// select {
		// case <-quitGoRoutine:
		// 	quitGoRoutine = make(chan struct{})
		// 	err = changeFileOwner(tarballFilePath)
		// 	currentStoreFiles.inUse = false
		// 	return countFiles, errors.New(sts["userCancelled"])
		// default:
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			return countFiles, errors.New(fmt.Sprintf("Could not add file '%s', to tarball, got error '%s'", filePath, err.Error()))
			// }
		}
	}
	err = changeFileOwner(tarballFilePath)
	// currentStoreFiles.inUse = false
	return countFiles, err
}

// UntarGzip: unpack tar.gz files list. Return len(filesList)=0 if all files has been restored.
// func UntarGz(sourcefile, dst string, filesList *[]string) (err error) {
// 	var ok bool
// 	tmpFilesList := *filesList
// 	var restoreFilesList = func() {
// 		*filesList = tmpFilesList
// 	}
// 	defer restoreFilesList()

// 	var createPath = func(filename string, header *tar.Header) (err error) {
// 		_, err = os.Stat(filepath.Dir(filename))
// 		if os.IsNotExist(err) {
// 			perm := header.FileInfo().Mode().Perm()
// 			err = os.MkdirAll(filepath.Dir(filename), os.FileMode(os.ModePerm|perm))
// 			if err != nil {
// 				return errors.New(fmt.Sprintf("Could not create directory '%s', got error '%s'", header.Name, err.Error()))
// 			}
// 		}
// 		return err
// 	}

// 	file, err := os.Open(sourcefile)
// 	if err != nil {
// 		return errors.New(fmt.Sprintf("Could not open tar.gz file '%s', got error '%s'", sourcefile, err.Error()))
// 	}
// 	defer file.Close()

// 	var fileReader io.ReadCloser = file
// 	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
// 	if mimeFile(sourcefile) == "gzip" {
// 		if fileReader, err = gzip.NewReader(file); err != nil {
// 			return errors.New(fmt.Sprintf("Could not create gz reader for '%s', got error '%s'", sourcefile, err.Error()))
// 		}
// 		defer fileReader.Close()
// 	} else if mimeFile(sourcefile) != "tar" {
// 		return errors.New(fmt.Sprintf("Could not read this kind of file '%s'", mimeFile(sourcefile)))
// 	}
// 	tarBallReader := tar.NewReader(fileReader)
// 	// Extracting tarred files
// 	var header = new(tar.Header)
// 	for {
// 		header, err = tarBallReader.Next()
// 		if err != nil {
// 			if err == io.EOF {
// 				err = nil
// 				break
// 			}
// 			return errors.New(fmt.Sprintf("Could not read archived file '%s', got error '%s'", header.Name, err.Error()))
// 		}

// 		filename := filepath.Join(dst, header.Name)
// 		if len(tmpFilesList) != 0 {
// 			ok = false
// 			for idx, file := range tmpFilesList {
// 				if strings.Contains(file, filename) || strings.Contains(filename, file) {
// 					ok = true
// 					if filepath.Join(dst, file) == filename {
// 						tmpFilesList = append(tmpFilesList[:idx], tmpFilesList[idx+1:]...)
// 					}
// 					break
// 				}
// 			}
// 			if !ok {
// 				header = nil
// 			}
// 		}

// 		if header != nil {
// 			switch header.Typeflag {
// 			case tar.TypeDir:
// 				if err = os.MkdirAll(filename, os.FileMode(header.Mode)); err == nil {
// 					err = os.Chown(filename, header.Uid, header.Gid)
// 					if err != nil {
// 						return errors.New(fmt.Sprintf("Could not copy data, chmod, chown file '%s', got error '%s'", filename, err.Error()))
// 					}
// 				}
// 			case tar.TypeReg:

// 				err = createPath(filename, header)
// 				if err != nil {
// 					return errors.New(fmt.Sprintf("Could not create path to hold file '%s', got error '%s'", filename, err.Error()))
// 				}

// 				writer, err := os.Create(filename)
// 				if err != nil {
// 					return errors.New(fmt.Sprintf("Could not create file '%s', got error '%s'", filename, err.Error()))
// 				}
// 				if _, err = io.Copy(writer, tarBallReader); err == nil {
// 					if err = os.Chmod(filename, os.FileMode(header.Mode)); err == nil {
// 						err = os.Chown(filename, header.Uid, header.Gid)
// 						if err != nil {
// 							return errors.New(fmt.Sprintf("Could not copy data, chmod, chown file '%s', got error '%s'", filename, err.Error()))
// 						}
// 					}
// 				}
// 				writer.Close()
// 			case tar.TypeSymlink:

// 				err = createPath(filename, header)
// 				if err != nil {
// 					return errors.New(fmt.Sprintf("Could not create path to hold file '%s', got error '%s'", filename, err.Error()))
// 				}

// 				if _, err := os.Lstat(filename); err == nil {
// 					os.Remove(filename)
// 				}
// 				err = os.Symlink(header.Linkname, filename)
// 				if err != nil {
// 					return errors.New(fmt.Sprintf("Could not create symlink '%s' to %s, got error '%s'", filename, header.Linkname, err.Error()))
// 				}
// 				if err = os.Chmod(filename, os.FileMode(header.Mode)); err == nil {
// 					err = os.Chown(filename, header.Uid, header.Gid)
// 					if err != nil {
// 						return errors.New(fmt.Sprintf("Could not chmod, chown file '%s', got error '%s'", filename, err.Error()))
// 					}
// 				}
// 			default:
// 				fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
// 			}
// 		}
// 	}
// 	return err
// }

// // magic nulber mime detection
// var magicTable = map[string]string{
// 	"\x37\x7A\xBC\xAF\x27\x1C\x00\x04": "7zip",
// 	"\xFD\x37\x7A\x58\x5A\x00\x00":     "xz",
// 	"\x1F\x8B\x08\x00\x00\x09\x6E\x88": "gzip",
// 	"\x75\x73\x74\x61\x72":             "tar",
// }

// // mimeFile: scan first bytes to detect mime type of file.
// func mimeFile(filename string) string {
// 	if file, err := os.Open(filename); err == nil {
// 		defer file.Close()
// 		buffReader := bufio.NewReader(file)
// 		for magic, mime := range magicTable {
// 			if peeked, err := buffReader.Peek(len([]byte(magic))); err == nil {
// 				tmpMagic := []byte(magic)
// 				if bytes.Index(peeked, tmpMagic) == 0 {
// 					return mime
// 				}
// 			}
// 		}
// 	}
// 	return "unknowen"
// }
