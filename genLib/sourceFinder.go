// sourceFinder.go

/*
	©2019 H.F.M. MIT license
*/

package genLib

import (
	"regexp"
	"strings"
)

type storeDecl struct {
	RowNb      int
	Definition string
}

type GoDeclarations struct {
	Filename       string
	rows           []string
	PackageName    string
	Functions      []storeDecl
	FunctNoComment []storeDecl
	FullFunc       []storeDecl
	Types          []storeDecl
	TypesNoComment []storeDecl
	Imports        []storeDecl
	FoundRows      []storeDecl
}

// getPackageName: retrieve package name of Go source file
func (d *GoDeclarations) GoSourceGetLines(filename string, wholeWord bool, terms ...string) (exist bool, err error) {
	var ww string
	var notMatch bool
	var regs []regexp.Regexp

	if d.readFile(filename) != nil {
		return exist, err
	}

	if wholeWord {
		ww = "\b"
	}

	for _, term := range terms {
		term = regexp.QuoteMeta(term)
		regs = append(regs, *regexp.MustCompile(`(` + ww + term + ww + `)`))
	}

	for idxRow, row := range d.rows {
		for _, reg := range regs {
			if !reg.MatchString(row) {
				notMatch = true

			}
		}
		if !notMatch {
			d.FoundRows = append(d.FoundRows, storeDecl{idxRow, row})
			exist = true
		}
		notMatch = false
	}
	return exist, err
}

// readFile: read file and make slice of strings
func (d *GoDeclarations) readFile(filename string) (err error) {
	data, err := ReadFile(filename)
	if err != nil {
		return err
	}
	d.rows = strings.Split(string(data), GetTextEOL(data))
	return err
}

// getPackageName: retrieve package name of Go source file
func (d *GoDeclarations) GoSourceGetInfos(filename string, funcName ...string) (err error) {
	var toFind string
	if len(funcName) != 0 {
		toFind = funcName[0]
	}

	if d.readFile(filename) != nil {
		return err
	}

	d.Types, err = d.getDecl("type", "}", "", true, toFind)
	if err == nil {
		d.TypesNoComment, err = d.getDecl("type", "}", "", false, toFind)
		if err == nil {
			d.FullFunc, err = d.getDecl("func", "}", "", false, toFind)
			if err == nil {
				d.Functions, err = d.getDecl("func", "", "{", true, toFind)
				if err == nil {
					d.FunctNoComment, err = d.getDecl("func", "", "{", false, toFind)
					if err == nil {
						err = d.getImports()
						if err == nil {
							err = d.getPackageName()
						}
					}
				}
			}
		}
	}
	return err
}

// getPackageName: retrieve package name of Go source file
func (d *GoDeclarations) getPackageName() error {
	pkgReg, err := regexp.Compile(`^(\bpackage\b)`)
	for idxRow := 0; idxRow < len(d.rows); idxRow++ {
		if pkgReg.MatchString(d.rows[idxRow]) {
			d.PackageName = strings.Split(d.rows[idxRow], " ")[1]
			break
		}
	}
	return err
}

// getImports: Scan go source file and return all requested imports.
func (d *GoDeclarations) getImports() error {
	var tempStrings []string

	importReg := regexp.MustCompile(`^(import .*)`)
	startMulti := regexp.MustCompile(`(.*\()$`)
	endMulti := regexp.MustCompile(`^(\))`)

	for idxRow := 0; idxRow < len(d.rows); idxRow++ {
		if importReg.MatchString(d.rows[idxRow]) {
			if !startMulti.MatchString(d.rows[idxRow]) {
				tempStrings = append(tempStrings, d.rows[idxRow])
				break
			}
			for !endMulti.MatchString(d.rows[idxRow]) {
				idxRow++
			}
			for idx := idxRow; idx >= 0; idx-- {
				if len(d.rows[idx]) != 0 {
					tempStrings = append(tempStrings, d.rows[idx])
				} else {
					tempStrings = append(tempStrings, "")
					break
				}
			}
		}
	}
	for idxRow := len(tempStrings) - 1; idxRow >= 0; idxRow-- {
		d.Imports = append(d.Imports, storeDecl{idxRow, tempStrings[idxRow]})
	}
	return nil
}

// getDecl: Scan go source file and return declarations.
func (d *GoDeclarations) getDecl(startDecl,
	endDeclAtStart,
	endDeclAtEnd string,
	wantComments bool,
	funcName ...string) (tmpStore []storeDecl, err error) {

	var stop bool
	var toFind string

	endDeclAtStart = regexp.QuoteMeta(endDeclAtStart)
	endDeclAtEnd = regexp.QuoteMeta(endDeclAtEnd)

	if len(funcName) != 0 {
		toFind = funcName[0]
	}

	toFindReg := regexp.MustCompile(`(\b` + toFind + `\b)`)
	startDeclReg := regexp.MustCompile(`^(\b` + startDecl + `\b)`)
	endDeclReg := regexp.MustCompile(`^(` + endDeclAtStart + `)`)
	if len(endDeclAtEnd) != 0 {
		endDeclReg = regexp.MustCompile(`(` + endDeclAtEnd + `)$`)
	}

	for idxRow := 0; idxRow < len(d.rows); idxRow++ {
		if startDeclReg.MatchString(d.rows[idxRow]) && toFindReg.MatchString(d.rows[idxRow]) {

			for len(d.rows[idxRow]) != 0 && wantComments && idxRow != 0 { //Go back to get comments
				idxRow--
			}

			if endDeclReg.MatchString(d.rows[idxRow]) {
				tmpStore = append(tmpStore, storeDecl{idxRow, d.rows[idxRow]})
				break
			} else {
				for !endDeclReg.MatchString(d.rows[idxRow]) {
					tmpStore = append(tmpStore, storeDecl{idxRow, d.rows[idxRow]})
					idxRow++
					stop = true
				}
			}
		}
		if stop {
			tmpStore = append(tmpStore, storeDecl{idxRow, d.rows[idxRow]})
			tmpStore = append(tmpStore, storeDecl{-1, ""})

			stop = false
		}
	}
	return tmpStore, nil
}
