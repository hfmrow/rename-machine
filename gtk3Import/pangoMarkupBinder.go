// PangoMarkupBinder.go

/*
*	©2019 H.F.M. MIT license
*	Handle Pango markup functions.

*	This is a pango markup binder to gotk3 pango library ...
*	I have made it to make more simple markup handling when i'm working on gtk objects.
*	It can be used in treeview, dialog, label, ... Each object where you ca.n set "markup" content
 */

package gtk3Import

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

var pangoEscapeChar = [][]string{{"<", "&lt;", string([]byte{0x15})}, {"&", "&amp;", string([]byte{0x16})}}

var markupType = map[string][]string{
	"bold": {"<b>", "</b>"}, "bld": {"<b>", "</b>"}, // Bold
	"big":   {"<big>", "</big>"},                                     // Makes font relatively larger, equivalent to <span size="larger">
	"small": {"<small>", "</small>"}, "sml": {"<small>", "</small>"}, // Makes font relatively smaller, equivalent to <span size="smaller">
	"italic": {"<i>", "</i>"}, "ita": {"<i>", "</i>"}, // Italic
	"subscript": {"<sub>", "</sub>"}, "sub": {"<sub>", "</sub>"}, // Subscript
	"supscript": {"<sup>", "</sup>"}, "sup": {"<sup>", "</sup>"}, // Superscript
	"monospace": {"<tt>", "</tt>"}, "msp": {"<tt>", "</tt>"}, // Monospace font

	"font_family": {`<span font_family="`, `">`, `</span>`},
	"ffy":         {`<span font_family="`, `">`, `</span>`}, // A font family name

	/*	A font description string, such as "Sans Italic 12". See pango_font_description_from_string() for a description of the format of the
		string representation. Note that any other span attributes will override this description.
		So if you have "Sans Italic" and also a style="normal" attribute, you will get Sans normal, not italic.*/
	"font": {`<span font="`, `">`, `</span>`},
	"fnt":  {`<span font="`, `">`, `</span>`},

	/*  Font size in 1024ths of a point, or one of the absolute sizes 'xx-small', 'x-small', 'small', 'medium', 'large', 'x-large', 'xx-large'
	    or one of the relative sizes 'smaller' or 'larger'. If you want to specify a absolute size, it's usually easier to take advantage of
		the ability to specify a partial font description using 'font'; you can use font='12.5' rather than size='12800'.*/
	"font_size": {`<span font_size="`, `">`, `</span>`},
	"fsz":       {`<span font_size="`, `">`, `</span>`},

	"strike": {"<s>", "</s>"}, "stk": {"<s>", "</s>"}, // Strikethrough
	"strikethrough_color": {`<span strikethrough="true" strikethrough_color="`, `">`, `</span>`}, // 'true' or 'false' whether to strike through the text
	"stc":                 {`<span strikethrough="true" strikethrough_color="`, `">`, `</span>`}, // An RGB color specification such as '#00FF00' or a color name such as 'blue'.

	"underline": {"<u>", "</u>"}, "und": {"<u>", "</u>"}, // Underline
	"underline_color": {`<span underline="`, `" underline_color="`, `">`, `</span>`}, // One of 'none', 'single', 'double', 'low', 'error'
	"udc":             {`<span underline="`, `" underline_color="`, `">`, `</span>`}, // An RGB color specification such as '#00FF00' or a color name such as 'red'.

	"foreground": {`<span foreground="`, `">`, `</span>`}, "fgc": {`<span foreground="`, `">`, `</span>`}, // An RGB color specification such as '#00FF00' or a color name such as 'red'.
	"background": {`<span background="`, `">`, `</span>`}, "bgc": {`<span background="`, `">`, `</span>`}, // An RGB color specification such as '#00FF00' or a color name such as 'red'.
	"fgalpha": {`<span fgalpha="`, `">`, `</span>`}, "fga": {`<span fgalpha="`, `">`, `</span>`}, // An alpha value for the background color, either a plain integer between 1 and 65536 or a percentage value like '50%'.
	"bgalpha": {`<span bgalpha="`, `">`, `</span>`}, "bga": {`<span bgalpha="`, `">`, `</span>`}, // An alpha value for the background color, either a plain integer between 1 and 65536 or a percentage value like '50%'.

	"url": {`<a href="`, `">`, `</a>`}, // Url clickable 1st arg: adress

	/*
	 Doesn't work on my system with std fonts/size ...
	*/
	"font_style": {`<span font_style="`, `">`, `</span>`}, // font_style: One of 'normal', 'oblique', 'italic'. N.b: 'oblique' seems to be the same as 'italic'
	"fst":        {`<span font_style="`, `">`, `</span>`}, // font_style: One of 'normal', 'oblique', 'italic'. So, look really useless ...

	"font_variant": {`<span font_variant="`, `">`, `</span>`},
	"fvt":          {`<span font_variant="`, `">`, `</span>`}, // One of 'normal' or 'smallcaps'

	"font_stretch": {`<span font_stretch="`, `">`, `</span>`},
	"fsh":          {`<span font_stretch="`, `">`, `</span>`}, // One of 'ultracondensed', 'extracondensed', 'condensed', 'semicondensed', 'normal', 'semiexpanded', 'expanded', 'extraexpanded', 'ultraexpanded'

	"font_weight": {`<span font_weight="`, `">`, `</span>`},
	"wgt":         {`<span font_weight="`, `">`, `</span>`}, // One of 'ultralight', 'light', 'normal', 'bold', 'ultrabold', 'heavy', or a numeric weight
}

type PangoColor struct {
	Black     string
	Brown     string
	White     string
	Red       string
	Green     string
	Blue      string
	Cyan      string
	Magenta   string
	Purple    string
	Turquoise string
	Violet    string
	Darkred   string
	Darkgreen string
	Darkblue  string
	Darkgray  string

	Darkcyan       string
	Lightblue      string
	Lightgray      string
	Lightgreen     string
	Lightturquoise string
	Lightred       string
	Lightyellow    string
}

func (pc *PangoColor) Init() {
	//	Colors initialisation
	pc.Black = "#000000"
	pc.Brown = "#7C2020"
	pc.White = "#FFFFFF"
	pc.Red = "#FF2222"
	pc.Green = "#22BB22"
	pc.Blue = "#0044FF"
	pc.Cyan = "#14FFFA"
	pc.Magenta = "#D72D6C"
	pc.Purple = "#8B0037"
	pc.Turquoise = "#009187"
	pc.Violet = "#7F00FF"
	pc.Darkred = "#300000"
	pc.Darkgreen = "#003000"
	pc.Darkblue = "#000030"
	pc.Darkcyan = "#003333"
	pc.Darkgray = "#303030"
	pc.Lightturquoise = "#80FFE7"
	pc.Lightblue = "#ADD8E6"
	pc.Lightgray = "#E4DDDD"
	pc.Lightgreen = "#87FF87"
	pc.Lightred = "#FF6666"
	pc.Lightyellow = "#FFFF6F"
}

type PangoMarkup struct {
	InString      string
	OutString     string
	markPositions [][]int
	markTypes     [][]string
	Colors        PangoColor
}

func (pm *PangoMarkup) Init(inString string) {
	pm.Colors.Init()
	// Object initialisation/cleaning
	pm.InString = inString
	pm.markPositions = [][]int{}
	pm.markTypes = [][]string{}
	pm.OutString = ""
}

// Add multiples positions, (where markup is applied)
func (pm *PangoMarkup) AddPosition(pos ...[]int) {
	pm.markPositions = append(pm.markPositions, pos...)
}

// Add multiples markup types, (the style applied at given positions)
func (pm *PangoMarkup) AddTypes(mType ...[]string) {
	pm.markTypes = append(pm.markTypes, mType...)
}

// Apply multiples pango markups to the whole text.
func (pm *PangoMarkup) Markup() string {
	pm.prepare()
	text := pm.InString
	for _, mType := range pm.markTypes {
		text = markup(text, mType...)
	}
	pm.OutString = text
	pm.finalize()
	return pm.OutString
}

// Apply multiples pango markups to text at specified positions given into 2d slices.
func (pm *PangoMarkup) MarkupAtPos() string {
	var eol = [][]byte{{0x0D, 0x0A}, {0x0D}, {0x0A}}
	var actEol string
	var multiMarks []string
	pm.prepare()
	// Sorting slice to get positions from the last to the first, (preserve positions in string)
	sort.SliceStable(pm.markPositions, func(i, j int) bool {
		return pm.markPositions[i][0] > pm.markPositions[j][0]
	})
	pm.OutString = pm.InString

	for _, pos := range pm.markPositions {
		prefix := pm.OutString[:pos[0]]
		toMark := pm.OutString[pos[0]:pos[1]]
		suffix := pm.OutString[pos[1]:]

		multiMarks = []string{toMark}
		for idx, val := range eol {
			if bytes.Contains([]byte(toMark), val) {
				actEol = string(eol[idx])
				multiMarks = strings.Split(toMark, actEol)
			}
		}
		for idx, _ := range multiMarks {
			for _, mType := range pm.markTypes {
				multiMarks[idx] = markup(multiMarks[idx], mType...)
			}
		}
		pm.OutString = prefix + strings.Join(multiMarks, actEol) + suffix
	}
	pm.finalize()
	return pm.OutString
}

// Prepare string with special characters to be marked ("<", "&")
func (pm *PangoMarkup) prepare() {
	pm.InString = strings.Replace(pm.InString, pangoEscapeChar[1][0], pangoEscapeChar[1][2], -1)
	pm.InString = strings.Replace(pm.InString, pangoEscapeChar[0][0], pangoEscapeChar[0][2], -1)
}

// Escape special characters after marking ("<", "&")
func (pm *PangoMarkup) finalize() {
	pm.OutString = strings.Replace(pm.OutString, pangoEscapeChar[1][2], pangoEscapeChar[1][1], -1)
	pm.OutString = strings.Replace(pm.OutString, pangoEscapeChar[0][2], pangoEscapeChar[0][1], -1)
}

// Apply pango markup format to text. They can be combined.
func markup(text string, mType ...string) string {
	switch len(mType) {
	case 1:
		// i.e: markup("display", "sub")
		return fmt.Sprint(markupType[mType[0]][0], text, markupType[mType[0]][1])
	case 2:
		// i.e: markup("display", "stc", "red")
		return fmt.Sprint(markupType[mType[0]][0], mType[1], markupType[mType[0]][1], text, markupType[mType[0]][2])
	case 3:
		// i.e: markup("display", "stc", "double", "red")
		return fmt.Sprint(markupType[mType[0]][0], mType[1], markupType[mType[0]][1], mType[2], markupType[mType[0]][2], text, markupType[mType[0]][3])
	default:
		return fmt.Sprint("Markup type error: ", mType)
	}
}
