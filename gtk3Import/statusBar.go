// statusBar.go

/*
*	©2019 H.F.M. MIT license
*	Handle Statusbar messages.
 */

package gtk3Import

import (
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type StatusBar struct {
	Messages  []string /* Each row contain associated strings refere to contextId number */
	statusbar *gtk.Statusbar
	contextId uint
	Prefix    []string
}

/* Add: add new element and return his own position. */
func (bar *StatusBar) Add(prefix, inString string) (position int) {
	bar.Prefix = append(bar.Prefix, prefix)
	bar.Messages = append(bar.Messages, prefix+" "+inString)
	bar.Disp()
	return len(bar.Messages) - 1
}

/* Add: set element at desired position. */
func (bar *StatusBar) Set(inString string, pos int) {
	if pos > len(bar.Messages)-1 || pos < 0 {
		inString = "Statusbar error: Invalid range to setting this message -> " + inString
		pos = len(bar.Messages) - 1
	}
	bar.Messages[pos] = inString
	bar.Disp()
}

/* Del: remove element at defined position and get the new length of elements. */
func (bar *StatusBar) Del(pos int) (newLength int) {
	copy(bar.Messages[pos:], bar.Messages[pos+1:])
	bar.Messages = bar.Messages[:cap(bar.Messages)-2]
	copy(bar.Prefix[pos:], bar.Prefix[pos+1:])
	bar.Prefix = bar.Prefix[:cap(bar.Prefix)-2]
	bar.Disp()
	return len(bar.Messages)
}

/* CleanAll: remove all elements (set to empty string) from the messages list. */
func (bar *StatusBar) CleanAll() {
	for idx, _ := range bar.Messages {
		bar.Messages[idx] = ""
	}
}

/* Disp: display content of stored elements into statusbar */
func (bar *StatusBar) Disp() {
	var dispMessages []string
	for idxMessage, message := range bar.Messages {
		if len(message) != 0 {
			dispMessages = append(dispMessages, bar.Prefix[idxMessage]+" "+message)
		}
	}
	bar.statusbar.Push(bar.contextId, strings.Join(dispMessages, " | "))
}

/* Init: Initialise structure to handle elements to be displayed. */
func (bar *StatusBar) Init(originStatusbar *gtk.Statusbar, prefix []string, stackId ...int) {
	var stack int
	if len(stackId) == 0 {
		stack = 0
	} else {
		stack = stackId[0]
	}
	bar.statusbar = originStatusbar
	bar.contextId = bar.statusbar.GetContextId(strconv.Itoa(stack)) /* get contextId of stack */
	bar.Messages = make([]string, len(prefix))
	for _, pre := range prefix {
		bar.Prefix = append(bar.Prefix, pre)
	}
}

/* Init: Initialise structure to handle elements to be displayed. */
func StatusBarNew(originStatusbar *gtk.Statusbar, prefix []string, stackId ...int) *StatusBar {
	bar := new(StatusBar)
	var stack int
	if len(stackId) == 0 {
		stack = 0
	} else {
		stack = stackId[0]
	}
	bar.statusbar = originStatusbar
	bar.contextId = bar.statusbar.GetContextId(strconv.Itoa(stack)) /* get contextId of stack */
	bar.Messages = make([]string, len(prefix))
	for _, pre := range prefix {
		bar.Prefix = append(bar.Prefix, pre)
	}
	return bar
}
