// searchAndReplace.go

package genLib

import (
	"os"
	"regexp"
	"strings"
)

type Find_s struct {
	TextBytes       []byte
	FileName        string
	LineEnd         string
	ToSearch        string
	ReplaceWith     string
	CaseSensitive   bool
	UseEscapeChar   bool
	POSIXcharClass  bool
	POSIXstrictMode bool
	Regex           bool
	Wildcard        bool
	WholeWord       bool
	DoReplace       bool
	Positions       Pos_s
}

type Pos_s struct {
	Line     []int
	AllLines []int
	WordsPos [][]int
}

type lineInfos_s struct {
	startPos   int
	endPos     int
	lineLength int
}

// Search in multiples text files , using "the Find_s struct" to fill needed informations about search preferences.
// return a slice type []Find_s that contain all informations about found patterns.
// Output can be saved with backup option.
func SearchAndReplaceInMultipleFiles(filenames []string, toSearch, replaceWith string,
	caseSensitive, POSIXcharClass, POSIXstrictMode, regex, wildcard, useEscapeChar, wholeWord,
	doReplace, doSave, doBackup, acceptBinary, removeEmptyResult bool) (founds []Find_s, err error) {

	founds = make([]Find_s, len(filenames))
	for idxFile, file := range filenames {
		// File exist ? ...
		if _, err = os.Stat(file); !os.IsNotExist(err) {
			founds[idxFile] = Find_s{}
			founds[idxFile].FileName = file
			// Check for text file
			isTxt, _ := IsTextFile(file)
			if !acceptBinary && !strings.Contains(isTxt, "utf-") {
				// Is binary file & not allowed
				founds[idxFile].TextBytes = []byte(isTxt)                                  // Put type of file in TextBytes field
				founds[idxFile].Positions.Line = append(founds[idxFile].Positions.Line, 0) // Adding a fake line to keep this entry
			} else {
				founds[idxFile].TextBytes, err = ReadFile(file)
				if err != nil {
					return founds, err
				}
				founds[idxFile].LineEnd = GetTextEOL(founds[idxFile].TextBytes)
				founds[idxFile].ToSearch = toSearch
				founds[idxFile].ReplaceWith = replaceWith
				founds[idxFile].CaseSensitive = caseSensitive
				founds[idxFile].DoReplace = doReplace
				founds[idxFile].POSIXcharClass = POSIXcharClass
				founds[idxFile].POSIXstrictMode = POSIXstrictMode
				founds[idxFile].Regex = regex
				founds[idxFile].UseEscapeChar = useEscapeChar
				founds[idxFile].WholeWord = wholeWord
				founds[idxFile].Wildcard = wildcard

				//				err = SearchAndReplace(&founds[idxFile])
				err = founds[idxFile].SearchAndReplace()
				if err != nil {
					return founds, err
				}
				// Saving file if one or more modifications was done
				if doSave && doReplace && len(founds[idxFile].Positions.WordsPos) > 0 {
					err = WriteFile(founds[idxFile].FileName, founds[idxFile].TextBytes, doBackup)
					if err != nil {
						return founds, err
					}
				}
			}
		} else { // File does not exist
			return founds, err
		}
	}
	// Removing empty structures if requested
	if removeEmptyResult {
		for idx := len(founds) - 1; idx >= 0; idx-- {
			if len(founds[idx].Positions.Line) == 0 {
				founds = append(founds[:idx], founds[idx+1:]...)
			}
		}
	}
	return founds, nil
}

// Search in plain text, use "the Find_s struct" to fill needed informations about search preferences.
// return a structure type Find_s that contain all informations about found patterns in the text
func (s *Find_s) SearchAndReplace() (err error) {

	// Get the number of lines and their absolutes positions in the text.
	var eolPos []lineInfos_s
	var regX *regexp.Regexp
	var search string
	tmpStrings := strings.Split(string(s.TextBytes), s.LineEnd)

	// Create lines indexes
	previousLineEndIdx := 0
	for idx, line := range tmpStrings {

		if len(line) != 0 {
			eolPos = append(eolPos,
				lineInfos_s{previousLineEndIdx + (len(s.LineEnd) * idx),
					len(line) + previousLineEndIdx + (len(s.LineEnd) * idx),
					len(line)})
			previousLineEndIdx = len(line) + previousLineEndIdx
		} else {
			eolPos = append(eolPos,
				lineInfos_s{previousLineEndIdx + (len(s.LineEnd) * idx),
					0,
					0})
		}
	}
	search = s.ToSearch

	if !s.Regex {
		switch {
		case s.POSIXcharClass:
			search = StringToCharacterClasses(search, s.CaseSensitive, s.POSIXstrictMode)
		case s.Wildcard:
			if !s.UseEscapeChar {
				search = strings.Replace(search, "?", "¤¤¤¤¤¤", -1)
				search = strings.Replace(search, "*", "¤¤¤¤¤", -1)
				search = regexp.QuoteMeta(search)
				search = strings.Replace(search, "¤¤¤¤¤¤", "?", -1)
				search = strings.Replace(search, "¤¤¤¤¤", "*", -1)
			}
			search = strings.Replace(search, "*", `.*`, -1)
			search = strings.Replace(search, "?", `.?`, -1)
		case !s.UseEscapeChar:
			search = regexp.QuoteMeta(search)
		}
		search = `(` + search + `)`

		if s.WholeWord {
			search = `\b` + search + `\b`
		}
		if !s.CaseSensitive && !s.POSIXcharClass {
			search = `(?i)` + search
		}
	}
	regX, err = regexp.Compile(search)
	if err != nil {
		return err
	}
	// Do the search/Replace job
	if s.DoReplace { // Replace requested .. do it and recursively call search function again.
		s.TextBytes = []byte(regX.ReplaceAllString(string(s.TextBytes), s.ReplaceWith))
		s.ToSearch = s.ReplaceWith
		s.DoReplace = false
		err = s.SearchAndReplace()
		if err != nil {
			return err
		}
		return nil
	} else { // Only search ... and store postitions
		if location := regX.FindAllStringIndex(string(s.TextBytes), -1); len(location) > 0 {

			// Proceed to compile pattern's positions to get lines numbers corresponding to search results,
			// the purpose of this step is to display (into control) only lines that contains the pattern founds.
			s.Positions.Line = append(s.Positions.Line, -1) // Add fake line to controle previous line
			for idxEol := len(eolPos) - 1; idxEol >= 0; idxEol-- {

				for idxLoc := 0; idxLoc < len(location); idxLoc++ {

					if eolPos[idxEol].startPos <= location[idxLoc][0] && eolPos[idxEol].endPos >= location[idxLoc][1] {
						// End of word IS in the line
						if idxEol != s.Positions.Line[len(s.Positions.Line)-1] { // Skip line if allready exist
							s.Positions.Line = append(s.Positions.Line, idxEol)
						}
						//						s.Positions.AllLines = append(s.Positions.AllLines, idxEol) // Store all line without skipping
					} else if eolPos[idxEol].startPos > location[idxLoc][0] && eolPos[idxEol].startPos < location[idxLoc][1] {
						// start of word is NOT in the line
						if idxEol != s.Positions.Line[len(s.Positions.Line)-1] { // Skip line if allready exist
							s.Positions.Line = append(s.Positions.Line, idxEol)
						}
						//						s.Positions.AllLines = append(s.Positions.AllLines, idxEol) // Store all line without skipping
						var tmpLines []int
						for eolPos[idxEol].startPos > location[idxLoc][0] && eolPos[idxEol].startPos < location[idxLoc][1] {

							idxEol--
							tmpLines = append(tmpLines, idxEol)
						}
						s.Positions.Line = append(s.Positions.Line, tmpLines...)
					}
				}
			}
			s.Positions.Line = s.Positions.Line[1:len(s.Positions.Line)] // Remove fake line
			s.Positions.WordsPos = append(s.Positions.WordsPos, location...)
		}
	}

	return err
}
