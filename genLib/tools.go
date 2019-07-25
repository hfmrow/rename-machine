// tools.go

/*
*	©2019 H.F.M. MIT license
*	Some functions facility
 */

package genLib

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// GetUser: retrieve realUser and currentUser.
func GetUser() (realUser, currentUser *user.User, err error) {
	if currentUser, err = user.Current(); err == nil {
		realUser, err = user.Lookup(os.Getenv("SUDO_USER"))
	}
	return realUser, currentUser, err
}

// changeOwner: set file owner to real user instead of root.
func changeFileOwner(filename string) (err error) {
	var realUser *user.User
	cmd := exec.Command("id", "-u")
	output, _ := cmd.Output()
	if string(output[:len(output)-1]) == "0" {
		if realUser, err = user.Lookup(os.Getenv("SUDO_USER")); err == nil {
			if uid, err := strconv.Atoi(realUser.Uid); err == nil {
				if gid, err := strconv.Atoi(realUser.Gid); err == nil {
					err = os.Chown(filename, uid, gid)
				}
			}
		}
	}
	return err
}

// UrlGet: find into sentence the url available part.
func UrlsGet(inString string) []string {
	reg := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(\/\S*)?`)
	return reg.FindAllString(inString, -1)
}

// ExecCommand: launch commandline application with arguments
// return output and error.
func ExecCommand(commands string, args ...string) (output []byte, err error) {
	output, err = exec.Command(commands, args...).CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, err
}

// Breakpoint check
// TODO adding function caller
func Bpcheck(val1, val2 int, pos ...int) {
	if len(pos) == 0 {
		pos = append(pos, -1)
	}
	if val1 == val2 {
		fmt.Printf("Position: %v\n", pos[0])
	}
}

// Measuring lapse (may be multiples) between operations
type Bench struct {
	lapse   []time.Time
	label   []string
	totalNs int64
	Results []string
	Average string
	Display bool
}

func (b *Bench) Lapse(label ...string) {
	b.lapse = append(b.lapse, time.Now())
	if len(label) == 0 {
		label = append(label, fmt.Sprintf("%d", len(b.lapse)))
	}
	b.label = append(b.label, label...)
}

func (b *Bench) Stop() {
	b.Lapse("Total")
	lapseCount := len(b.lapse) - 1
	var getMSmsnano = func(diff int64) (min, sec, ms, ns int64) {
		min = (diff / 1000000000) / 60
		sec = diff/1000000000 - (min * 60)
		ms = (diff / 1000000) - (sec * 1000)
		ns = diff - ((diff / 1000000) * 1000000)
		return min, sec, ms, ns
	}
	var calculateLapse = func(count int, start, stop time.Time) {
		diff := stop.Sub(start).Nanoseconds()
		m, s, ms, ns := getMSmsnano(diff)
		b.Results = append(b.Results, fmt.Sprintf("%s: %v m, %v s, %v ms, %v ns",
			b.label[count], m, s, ms, ns))
		b.totalNs += diff
		if b.Display {
			fmt.Println(b.Results[len(b.Results)-1] + GetOsLineEnd())
		}
	}

	b.Results = b.Results[:0]
	if lapseCount > 1 {
		for idx := 0; idx < len(b.lapse)-1; idx++ {
			calculateLapse(idx, b.lapse[idx], b.lapse[idx+1])
		}
	}
	calculateLapse(lapseCount, b.lapse[0], b.lapse[lapseCount])

	m, s, ms, ns := getMSmsnano(b.totalNs / int64(len(b.lapse)))
	b.Average = fmt.Sprintf("%v m, %v s, %v ms, %v ns", m, s, ms, ns)

	b.totalNs = 0
	b.label = b.label[:0]
	b.lapse = b.lapse[:0]
}

// Get input from commandline stdin
func GetStdin(ask string) (input string) {
	fmt.Print(ask + ": ")
	fmt.Scanln(&input)
	return input
}

type timeStamp struct {
	Year          string
	YearCopyRight string
	Month         string
	MonthWord     string
	Day           string
	DayWord       string
	Date          string
	Time          string
	Full          string
	FullNum       string
}

// Get current timestamp
func TimeStamp() *timeStamp {
	ts := new(timeStamp)
	timed := time.Now()
	regD := regexp.MustCompile("([^[:digit:]])")
	regA := regexp.MustCompile("([^[:alpha:]])")
	splitedNum := regD.Split(timed.Format(time.RFC3339), -1)
	splitedWrd := regA.Split(timed.Format(time.RFC850), -1)
	ts.Year = splitedNum[0]
	ts.Month = splitedNum[1]
	ts.Day = splitedNum[2]
	ts.Time = splitedNum[3] + `:` + splitedNum[4] + `:` + splitedNum[5]
	ts.DayWord = splitedWrd[0]
	ts.MonthWord = splitedWrd[5]
	ts.YearCopyRight = `©` + ts.Year
	ts.Full = strings.Join(strings.Split(timed.Format(time.RFC1123), " ")[:5], " ")

	nonAlNum := regexp.MustCompile(`[[:punct:][:alpha:]]`)
	date := nonAlNum.ReplaceAllString(time.Now().Format(time.RFC3339), "")[:14]
	ts.FullNum = date[:8] + "-" + date[8:]
	return ts
}

// Check: Display error messages in HR version with onClickJump enabled in
// my favourite Golang IDE editor. Return true if error exist.
func Check(err error, message ...string) (state bool) {
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`) //	to match 2 or more whitespace symbols inside a string
	var msgs string
	if err != nil {
		state = true
		if len(message) != 0 { // Make string with messages if exists
			for _, mess := range message {
				msgs += `[` + mess + `]`
			}
		}
		pc, file, line, ok := runtime.Caller(1) //	(pc uintptr, file string, line int, ok bool)
		if ok == false {                        // Remove "== false" if needed
			fName := runtime.FuncForPC(pc).Name()
			fmt.Printf("[%s][%s][File: %s][Func: %s][Line: %d]\n", msgs, err.Error(), file, fName, line)
		} else {
			stack := strings.Split(fmt.Sprintf("%s", debug.Stack()), "\n")
			for idx := 5; idx < len(stack)-1; idx = idx + 2 {
				//	To match 2 or more whitespace leading/ending/inside a string (include \t, \n)
				mess1 := strings.Join(strings.Fields(stack[idx]), " ")
				mess2 := strings.TrimSpace(remInside.ReplaceAllString(stack[idx+1], " "))
				fmt.Printf("%s[%s][%s]\n", msgs, err.Error(), strings.Join([]string{mess1, mess2}, "]["))
			}
		}
	}
	return state
}

// Check error function, and exit if error, input message is optional or accept multiple arguments.
func CheckE(err error, message ...string) {
	var msgs string
	if err != nil {
		if len(message) != 0 { // Make string with messages if exists
			for _, mess := range message {
				msgs += "[ " + mess + " ]"
			}
			fmt.Println("Error: " + msgs)
		}
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

// Check error and return true if error exist
func IsError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

// Use function to avoid "Unused variable ..." msgs
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
