// formatDate.go

// TODO Rewrite on next use ....

package genLib

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DateEntry struct {
	InPos  int
	OutPos int
	Val    int
	Str    string
}

func NewDateFormat() []string {
	tmpString := []string{}
	tmpString = append(tmpString, "%d-%m-%y") // Day-Month-Year
	tmpString = append(tmpString, "%y-%m-%d") // Year-Month-Day
	tmpString = append(tmpString, "%m-%d-%y") // Month-Day-Year
	return tmpString
}

// Find the date in relation to the given format. (French format like %d, %m, %y or Y%). input fmtDate: "%d-%m-%y %H:%M:%S"
func FindDate(inputString, fmtDate string) []string {
	// Remove %H:%M:%S
	splitSpace := strings.Split(fmtDate, " ")
	// Split d, m, y
	splitAlNum := regexp.MustCompile("[[:punct:]]")
	outstr := splitAlNum.Split(strings.ToLower(splitSpace[0]), -1)
	var outShort, outLong, outSingleShort, outSingleLong, out2Single1Short, out2Single1Long, out1Single2Short, out1Single2Long []string
	// Replace d with 21, m with 12  and y or Y with 73 or 1973, there will be converted into character class equivalent
	for _, val := range outstr {
		switch val {
		case "d":
			out1Single2Short = append(out1Single2Short, "1")
			out1Single2Long = append(out1Single2Long, "1")
			out2Single1Short = append(out2Single1Short, "21")
			out2Single1Long = append(out2Single1Long, "21")
			outSingleShort = append(outSingleShort, "1")
			outSingleLong = append(outSingleLong, "1")
			outShort = append(outShort, "21")
			outLong = append(outLong, "21")
		case "m":
			out1Single2Short = append(out1Single2Short, "21")
			out1Single2Long = append(out1Single2Long, "21")
			out2Single1Short = append(out2Single1Short, "1")
			out2Single1Long = append(out2Single1Long, "1")
			outSingleShort = append(outSingleShort, "2")
			outSingleLong = append(outSingleLong, "2")
			outShort = append(outShort, "12")
			outLong = append(outLong, "12")
		case "y":
			out1Single2Short = append(out1Single2Short, "73")
			out1Single2Long = append(out1Single2Long, "1973")
			out2Single1Short = append(out2Single1Short, "73")
			out2Single1Long = append(out2Single1Long, "1973")
			outSingleShort = append(outSingleShort, "73")
			outSingleLong = append(outSingleLong, "1973")
			outShort = append(outShort, "73")
			outLong = append(outLong, "1973")
		}
	}
	dateLongFmt := regexp.MustCompile(StringToCharacterClasses(strings.Join(outLong, "/"), false, false))
	dateShrtFmt := regexp.MustCompile(StringToCharacterClasses(strings.Join(outShort, "/"), false, false))
	dateSingleLongFmt := regexp.MustCompile(StringToCharacterClasses(strings.Join(outSingleLong, "/"), false, false))
	dateSingleShrtFmt := regexp.MustCompile(StringToCharacterClasses(strings.Join(outSingleShort, "/"), false, false))
	date1Single2LongFmt := regexp.MustCompile(StringToCharacterClasses(strings.Join(out1Single2Long, "/"), false, false))
	date1Single2ShrtFmt := regexp.MustCompile(StringToCharacterClasses(strings.Join(out1Single2Short, "/"), false, false))
	date2Single1LongFmt := regexp.MustCompile(StringToCharacterClasses(strings.Join(out2Single1Long, "/"), false, false))
	date2Single1ShrtFmt := regexp.MustCompile(StringToCharacterClasses(strings.Join(out2Single1Short, "/"), false, false))
	// Search for character class replacement of integer.
	result := dateLongFmt.FindAllString(inputString, 1)
	if result == nil {
		result = append(result, dateShrtFmt.FindAllString(inputString, 1)...)
		if result == nil {
			result = append(result, date1Single2LongFmt.FindAllString(inputString, 1)...)
			if result == nil {
				result = append(result, date1Single2ShrtFmt.FindAllString(inputString, 1)...)
				if result == nil {
					result = append(result, date2Single1LongFmt.FindAllString(inputString, 1)...)
					if result == nil {
						result = append(result, date2Single1ShrtFmt.FindAllString(inputString, 1)...)
						if result == nil {
							result = append(result, dateSingleLongFmt.FindAllString(inputString, 1)...)
							if result == nil {
								result = append(result, dateSingleShrtFmt.FindAllString(inputString, 1)...)
							}
						}
					}
				}
			}
		}
	}
	return result
}

// Change date format to unix model with pattern like: fmtDate:='%d-%m-%y %H:%M:%S' --> '$y-%m-%d %H:%M:%S'
// Note: year format must be in lowercase ...
func FormatDate(fmtDate string, inpDate string) time.Time {
	// Make variable with defaults values (Unix date)
	lnxlistDate := []DateEntry{}
	fmtListDate := []DateEntry{}
	tmpListDate := [6]DateEntry{}
	tmpInt := 0
	timeUnix := []string{`y`, `m`, `d`, `H`, `M`, `S`}
	for idx, str := range timeUnix {
		lnxlistDate = append(lnxlistDate, DateEntry{0, idx, 0, str})
	} // Split fmtDate elements
	regExp := `[[:alpha:]]`
	re := regexp.MustCompile(regExp)
	submatchall := re.FindAllString(fmtDate, -1)
	for idx, str := range submatchall {
		fmtListDate = append(fmtListDate, DateEntry{idx, 0, 0, str})
	} // Split inpDate elements
	regExp = `[^[:alnum:]]`
	re = regexp.MustCompile(regExp)
	submatchall = re.Split(inpDate, -1)
	for idx, str := range submatchall {
		tmpInt, _ = strconv.Atoi(str)
		fmtListDate[idx] = DateEntry{fmtListDate[idx].InPos, 0, tmpInt, fmtListDate[idx].Str}
	} // Mixing data
	for lnxIdx, lnx := range lnxlistDate {
		for fmtIdx, frmt := range fmtListDate {
			if frmt.Str == timeUnix[0] && len(strconv.Itoa(frmt.Val)) == 2 { // Deal with short year format to convert to 4 digits
				if fmtListDate[fmtIdx].Val <= 70 { // Using EPOCH year date (1970) as reference
					fmtListDate[fmtIdx].Val = fmtListDate[fmtIdx].Val + 2000
				} else {
					fmtListDate[fmtIdx].Val = fmtListDate[fmtIdx].Val + 1900
				}
			} // Create output datas: `[{0 2 9 d} {1 1 6 m} {2 0 2018 y} {3 3 6 H} {4 4 54 M} {5 5 5 S}]`
			if lnx.Str == frmt.Str { // from this one: `9/06/18 6-54-05` for exemple
				fmtListDate[fmtIdx].OutPos = lnxlistDate[lnxIdx].OutPos
			}
		}
	}
	// Sorting result array to conforming unix std: [{2 0 2018 y} {1 1 6 m} {0 2 9 d} {3 3 6 H} {4 4 54 M} {5 5 5 S}]
	for idx := 0; idx < len(fmtListDate); idx++ {
		tmpListDate[fmtListDate[idx].OutPos] = fmtListDate[idx]
	}
	//	fmt.Println(lnxlistDate)	//	fmt.Println(fmtListDate)	//	fmt.Println(tmpListDate)
	return time.Date(tmpListDate[0].Val, time.Month(tmpListDate[1].Val), tmpListDate[2].Val, tmpListDate[3].Val, tmpListDate[4].Val, tmpListDate[5].Val, 0, time.Local)
}
