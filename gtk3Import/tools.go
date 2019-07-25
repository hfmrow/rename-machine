// tools.go

/// +build ignore

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"fmt"
	"regexp"

	"github.com/gotk3/gotk3/gtk"
)

// Fill / Clean comboBoxText
func ComboBoxTextFill(cbxEntry *gtk.ComboBoxText, entries []string, removAll ...bool) {
	if len(removAll) == 0 {
		for _, word := range entries {
			cbxEntry.PrependText(word)
		}
	} else if removAll[0] {
		cbxEntry.RemoveAll()
	}
}

// Get text from gtk.comboBoxTextEntry object
func ComboBoxTextGetEntry(cbxEntry *gtk.ComboBoxText) string {
	return cbxEntry.GetActiveText()
}

// SpinbuttonSetValues: Configure spin button
func SpinbuttonSetValues(sb *gtk.SpinButton, min, max, value int, step ...int) (err error) {
	incStep := 1
	if len(step) != 0 {
		incStep = step[0]
	}
	if ad, err := gtk.AdjustmentNew(float64(value), float64(min), float64(max), float64(incStep), 0, 0); err == nil {
		sb.Configure(ad, 1, 0)
	}
	return err
}

// MarkupHttpClickable:
func MarkupHttpClickable(inString string) (outString string) {
	pm := PangoMarkup{}
	outString = inString // Search for http adress to be treated as clickable link
	reg := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(\/\S*)?`)
	indexes := reg.FindAllIndex([]byte(inString), 1)

	if len(indexes) != 0 {
		pm.Init(inString)
		pm.AddPosition([]int{indexes[0][0], indexes[0][1]})
		mtype := [][]string{{"url", reg.FindString(inString)}}
		pm.AddTypes(mtype...)
		outString = pm.MarkupAtPos()
	}
	return outString
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

// LowercaseAtFirst: true if 1st char is lowercase (Exist into GenLib too !)
func LowercaseAtFirst(inString string) bool {
	if len(inString) != 0 {
		charType, _ := regexp.Compile("[[:lower:]]")
		return charType.MatchString(inString[:1])
	}
	return true
}
