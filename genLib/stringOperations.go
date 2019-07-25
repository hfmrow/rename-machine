// stringOperations.go

package genLib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// LowercaseAtFirst: true if 1st char is lowercase
func LowercaseAtFirst(inString string) bool {
	if len(inString) != 0 {
		charType, _ := regexp.Compile("[[:lower:]]")
		return charType.MatchString(inString[:1])
	}
	return true
}

// ToCamel: Turn string into camel case
func ToCamel(inString string, lowerAtFirst ...bool) (outString string) {
	var laf bool
	if len(lowerAtFirst) != 0 {
		laf = lowerAtFirst[0]
	}
	nonAlNum := regexp.MustCompile(`[[:punct:][:space:]]`)
	tmpString := nonAlNum.Split(inString, -1)

	for idx, word := range tmpString {
		if laf && idx < 1 {
			outString += strings.ToLower(word)
		} else {
			outString += strings.Title(word)
		}
	}
	return outString
}

// removeEmptyLines:
func removeEmptyLines(inString string) (ouString string) {
	tmpSplitted := strings.Split(inString, GetTextEOL([]byte(inString)))
	for idx := len(tmpSplitted) - 1; idx >= 0; idx-- {
		if len(tmpSplitted[idx]) == 0 {
			tmpSplitted = append(tmpSplitted[:idx], tmpSplitted[idx+1:]...)
		}
	}
	return strings.Join(tmpSplitted, GetTextEOL([]byte(inString)))
}

var platforms = [][]string{
	{"darwin", "386", "\n"},
	{"darwin", "amd64", "\n"},
	{"dragonfly", "amd64", "\n"},
	{"freebsd", "386", "\n"},
	{"freebsd", "amd64", "\n"},
	{"freebsd", "arm", "\n"},
	{"linux", "386", "\n"},
	{"linux", "amd64", "\n"},
	{"linux", "arm", "\n"},
	{"linux", "arm64", "\n"},
	{"linux", "ppc64", "\n"},
	{"linux", "ppc64le", "\n"},
	{"linux", "mips", "\n"},
	{"linux", "mipsle", "\n"},
	{"linux", "mips64", "\n"},
	{"linux", "mips64le", "\n"},
	{"linux", "s390x", "\n"},
	{"nacl", "386", "\n"},
	{"nacl", "amd64p32", "\n"},
	{"nacl", "arm", "\n"},
	{"netbsd", "386", "\n"},
	{"netbsd", "amd64", "\n"},
	{"netbsd", "arm", "\n"},
	{"openbsd", "386", "\n"},
	{"openbsd", "amd64", "\n"},
	{"openbsd", "arm", "\n"},
	{"plan9", "386", "\n"},
	{"plan9", "amd64", "\n"},
	{"plan9", "arm", "\n"},
	{"solaris", "amd64", "\n"},
	{"windows", "386", "\r\n"},
	{"windows", "amd64", "\r\n"}}

// GetOsLineEnd: Get current OS line-feed
func GetOsLineEnd() string {
	for _, row := range platforms {
		if row[0] == runtime.GOOS {
			return row[2]
		}
	}
	return "\n"
}

// GetTextEOL: Get EOL from text bytes (CR, LF, CRLF) > string or get OS line end.
func GetTextEOL(inTextBytes []byte) (outString string) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	if bytes.Contains(inTextBytes, bCRLF) {
		return string(bCRLF)
	} else if bytes.Contains(inTextBytes, bCR) {
		return string(bCR)
	}
	if bytes.Contains(inTextBytes, bLF) {
		return string(bLF)
	}
	return GetOsLineEnd()
}

// SetTextEOL: Get EOL from text bytes and convert it to another EOL (CR, LF, CRLF)
func SetTextEOL(inTextBytes []byte, eol string) (outTextBytes []byte, err error) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	var outEol []byte
	switch eol {
	case "CR":
		outEol = bCR
	case "LF":
		outEol = bLF
	case "CRLF":
		outEol = bCRLF
	default:
		return outTextBytes, errors.New("EOL convert error: Undefined end of line")
	}
	// Handle end of line
	outTextBytes = bytes.Replace(inTextBytes, []byte(GetTextEOL(inTextBytes)), outEol, -1)
	return outTextBytes, nil
}

// WordsFrequency: Counts the frequency with which words appear in a text.
func WordsFrequency(text string) (list [][]string) {
	tmpText := RemoveNonAlNum(text)
	tmpText, _ = TrimSpace(tmpText, "-c")
	tmpLines := strings.Fields(tmpText)

	for _, word := range tmpLines {
		if !IsExist2dCol(list, word, 0) {
			list = append(list, []string{word, ""})
		}
	}
	tmpText = strings.Join(tmpLines, " ")
	for idx, word := range list {
		// Compile for matching whole word
		regX, err := regexp.Compile(`\b(` + word[0] + `)\b`)
		if err == nil {
			list[idx][1] = fmt.Sprintf("%d", len(regX.FindAllString(tmpText, -1)))
		}
	}
	// Sorting result ...
	sort.SliceStable(list, func(i, j int) bool {
		return fmt.Sprintf("%d", list[i][1]) > fmt.Sprintf("%d", list[j][1])
	})
	return list
}

// IsFloat: Check if string is float
func IsFloat(inString string) bool {
	_, err := strconv.ParseFloat(strings.Replace(strings.Replace(inString, " ", "", -1), ",", ".", -1), 64)
	if err == nil {
		return true
	}
	return false
}

// IsDate: Check if string is date
func IsDate(inString string) bool {
	dateFormats := NewDateFormat()
	for _, dteFmt := range dateFormats {
		if len(FindDate(inString, dteFmt+" %H:%M:%S")) != 0 {
			return true
		}
	}
	return false
}

// ByteToHexStr: Convert []byte to hexString
func ByteToHexStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// GenFileName: Generate a randomized file name
func GenFileName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Md5String(fmt.Sprint(r.Int63n(time.Now().UnixNano())))
}

// RemoveNonAlNum: Remove all non alpha-numeric char
func RemoveNonAlNum(inString string) string {
	nonAlNum := regexp.MustCompile(`[[:punct:]]`)
	return nonAlNum.ReplaceAllString(inString, "")
}

// RemoveNonNum: Remove all non numeric char
func RemoveNonNum(inString string) string {
	nonAlNum := regexp.MustCompile(`[[:punct:][:alpha:]]`)
	return nonAlNum.ReplaceAllString(inString, "")
}

// ReplaceSpace: replace all [[:space::]] with underscore "_"
func ReplaceSpace(inString string) string {
	spaceRegex := regexp.MustCompile(`[[:space:]]`)
	return spaceRegex.ReplaceAllString(inString, "_")
}

// ReplaceSpace: replace all [[:space::]] with underscore "_"
func RemoveSpace(inString string) string {
	spaceRegex := regexp.MustCompile(`[[:space:]]`)
	return spaceRegex.ReplaceAllString(inString, "")
}

// ReplacePunct: replace all [[:punct::]] with underscore "_"
func ReplacePunct(inString string) string {
	spaceRegex := regexp.MustCompile(`[[:punct:]]`)
	return spaceRegex.ReplaceAllString(inString, "_")
}

// SplitNumeric: Split and keep all numeric values in a string
func SplitNumeric(inString string) (outText []string, err error) {
	toSplit := regexp.MustCompile(`[[:alpha:][:punct:]]`)
	spaceSepared := string(toSplit.ReplaceAll([]byte(inString), []byte(" ")))
	spaceSepared, err = TrimSpace(spaceSepared, "-c")
	if err != nil {
		return outText, err
	}
	outText = strings.Split(spaceSepared, " ")
	return outText, err
}

// An unwanted behavior may occur on string where word's length > max...
func FormatTextQuoteBlankLines(str string) (out string) {
	outStrings := strings.Split(str, GetTextEOL([]byte(str)))
	for _, line := range outStrings {
		if len(line) != 0 {
			out += strings.TrimSuffix(strings.TrimPrefix(strconv.Quote(line+"\n"), `"`), `"`)
		} else {
			out += `\n`
		}
	}
	return `"` + out + `"`
}

// FormatTextParagraphFormatting:
func FormatTextParagraphFormatting(str string, max int, truncateWordOverMaxSize bool, indentFirstLinetr ...string) string {
	var tmpParag, tmpFinal, indent string
	eol := GetTextEOL([]byte(str))
	var addLine = func() {
		// if indent != "" {
		// 	tmpFinal += indent
		// }
		tmpFinal += FormatText(tmpParag, max, truncateWordOverMaxSize, indent)
	}
	if len(indentFirstLinetr) != 0 {
		indent = indentFirstLinetr[0]
	}
	lines := strings.Split(str, eol)
	for _, line := range lines {
		if len(line) != 0 {
			tmpParag += line
		} else {
			addLine()
			tmpFinal += eol + eol
			tmpParag = tmpParag[:0]
		}
	}
	addLine()
	return strings.TrimSuffix(tmpFinal, eol)
}

// FormatText: Format words text to fit (column/windows with limited width) "max" chars.
// An unwanted behavior may occur on string where word's length > max...
func FormatText(str string, max int, truncateWordOverMaxSize bool, indenT ...string) string {
	var indent string
	if len(indenT) != 0 {
		indent = indenT[0]
	}
	if len(str) > 0 {
		eol := GetTextEOL([]byte(str))
		var tmpLines, tmpWordTooLong []string
		space := regexp.MustCompile(`[[:space:]]`)
		var outText []string
		var countChar, length int
		// in case where string does not contain any [[:space:]] char
		var cutLongString = func(inStr string, inMax int, truncate bool) (outSliceText []string) {
			if !truncate {
				length = len(inStr)
				for count := 0; count < length; count = count + inMax {
					if count+inMax < length {
						outSliceText = append(outSliceText, inStr[count:count+inMax])
					} else {
						outSliceText = append(outSliceText, inStr[count:length])
					}
				}
			} else {
				outSliceText = append(outSliceText, TruncateString(inStr, "...", inMax, 1))
			}
			return outSliceText
		}
		text := space.Split(str, -1) // Split str at each blank char
		for idxWord := 0; idxWord < len(text); idxWord++ {
			length = len(text[idxWord]) + 1
			if length >= max {
				tmpWordTooLong = cutLongString(text[idxWord], max-1, truncateWordOverMaxSize)
				text = append(text[:idxWord], text[idxWord+1:]...) // Remove slice entry
				// Insert slices entries
				text = append(text[:idxWord], append(tmpWordTooLong, text[idxWord:]...)...)
				length = len(text[idxWord]) // Calculate new length
			}

			if countChar+length <= max {
				tmpLines = append(tmpLines, text[idxWord])
				countChar += length
			} else {
				outText = append(outText, indent+strings.Join(tmpLines, " "))
				tmpLines = tmpLines[:0] // Clear slice
				countChar = 0
				idxWord--
			}
		}
		// Get the rest of the text.
		outText = append(outText, indent+strings.Join(tmpLines, " "))
		return strings.Join(outText, eol)
	}
	return ""
}

// ReducePath: Reduce patch length by preserving count element from the end
func TruncatePath(fullpath string, count ...int) (reduced string) {
	elemCnt := 2
	if len(count) != 0 {
		elemCnt = count[0]
	}
	splited := strings.Split(fullpath, string(os.PathSeparator))
	if len(splited) > elemCnt+1 {
		return "..." + string(os.PathSeparator) + filepath.Join(splited[len(splited)-elemCnt:]...)
	}
	return fullpath
}

// TruncateString: Reduce string length for display (prefix is separator like: "...", option=0 -> put separator at the begening
// of output string. Option=1 -> center, is where separation is placed. option=2 -> line feed, trunc the whole string using LF
// without shorting it. Max, is max char length of the output string.
func TruncateString(inString, prefix string, max, option int) string {
	var center, cutAt bool
	var outText string
	switch option {
	case 1:
		center = true
		cutAt = false
		max = max - len(prefix)
	case 2:
		center = false
		cutAt = true
	default:
		center = false
		cutAt = false
		max = max - len(prefix)
	}
	length := len(inString)
	if length > max {
		if cutAt {
			for count := 0; count < length; count = count + max {
				if count+max < length {
					outText += fmt.Sprintln(inString[count : count+max])
				} else {
					outText += fmt.Sprintln(inString[count:length])
				}
			}
			return outText
		} else if center {
			midLength := max / 2
			inString = inString[:midLength] + prefix + inString[length-midLength-1:]
		} else {
			inString = prefix + inString[length-max:]
		}
	}
	return inString
}

// TrimSpace: Some multiple way to trim strings. cmds is optionnal or accept multiples args
func TrimSpace(inputString string, cmds ...string) (newstring string, err error) {

	osForbiden := regexp.MustCompile(`[<>:"/\\|?*]`)
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)    //	to match 2 or more whitespace symbols inside a string
	remInsideNoTab := regexp.MustCompile(`[\p{Zs}]{2,}`) //	(preserve \t) to match 2 or more space symbols inside a string

	if len(cmds) != 0 {
		for _, command := range cmds {
			switch command {
			case "+h": //	Escape html
				inputString = html.EscapeString(inputString)
			case "-h": //	UnEscape html
				inputString = html.UnescapeString(inputString)
			case "+e": //	Escape specials chars
				inputString = fmt.Sprintf("%q", inputString)
			case "-e": //	Un-Escape specials chars
				tmpString, err := strconv.Unquote(`"` + inputString + `"`)
				if err != nil {
					return inputString, err
				}
				inputString = tmpString
			case "-w": //	Change all illegals chars (for path in linux and windows) into "-"
				inputString = osForbiden.ReplaceAllString(inputString, "-")
			case "+w": //	clean all illegals chars (for path in linux and windows)
				inputString = osForbiden.ReplaceAllString(inputString, "")
			case "-c": //	Trim [[:space:]] and clean multi [[:space:]] inside
				inputString = strings.TrimSpace(remInside.ReplaceAllString(inputString, " "))
			case "-ct": //	Trim [[:space:]] and clean multi [[:space:]] inside (preserve TAB)
				inputString = strings.Trim(remInsideNoTab.ReplaceAllString(inputString, " "), " ")
			case "-s": //	To match 2 or more whitespace leading/ending/inside a string (include \t, \n)
				inputString = strings.Join(strings.Fields(inputString), " ")
			case "-&": //	Replace ampersand CHAR with ampersand HTML code
				inputString = strings.Replace(inputString, "&", "&amp;", -1)
			case "+&": //	Replace ampersand HTML code with ampersand CHAR
				inputString = strings.Replace(inputString, "&amp;", "&", -1)
			default:
				return inputString, errors.New("TrimSpace, " + command + ", does not exist")
			}
		}
	}
	return inputString, nil
}
