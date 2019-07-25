// loremIpsum.go

package genLib

import (
	"math/rand"
	"strings"
	"time"
)

var words []string

// Generate Lorem Ipsum texte file
func LoremGenerateTextFile(filename string, maxWords, maxChars int) error {
	return WriteFile(filename, []byte(FormatText(Lorem(maxWords), maxChars, false)))
}

// Generate Lorem Ipsum text
func Lorem(wordCount ...int) string {
	if len(wordCount) == 0 {
		wordCount = append(wordCount, 50)
	}
	tmpSentence := make([]string, 0)
	histRand := make([]int, 0)
	words := newLoremWords()
	totalWords := len(words)
	punct := countPunct{}
	punct.Init()
	totalPunct := len(punct.Str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Get random word from list. Option=true get (non repeated) until limit of stored words
	rndWrd := func(noRepeat bool) string {
		var skip bool
		val := int(r.Int63n(int64(totalWords)))
		if noRepeat {
			for true {
				for _, hist := range histRand {
					// Check if this word has already been written
					if hist == val {
						skip = true
						break
					}
				}
				if !skip {
					histRand = append(histRand, val)
					break
				} else {
					// Reset word counter, limit reached ...
					if len(histRand) >= totalWords {
						histRand = histRand[:0]
					}
					skip = !skip
					val = int(r.Int63n(int64(totalWords)))
				}
			}
		}
		return words[val]
	}
	// Get random punctuation from list.
	rndPunct := func() string {
		val := int(r.Int63n(int64(totalPunct)))
		punct.Counted[val]++
		if punct.Max[val] <= punct.Counted[val] {
			punct.Counted[val] = 0
			return punct.Str[val]
		}
		return ""
	}
	for count := 0; count < wordCount[0]; count++ {
		tmpSentence = append(tmpSentence, rndWrd(false)+rndPunct())
	}

	// Upper case first letter of word following a dot "."
	for idxWord := 0; idxWord < len(tmpSentence); idxWord++ {
		if tmpSentence[idxWord][len(tmpSentence[idxWord])-1] == []byte(".")[0] {
			if idxWord+1 < len(tmpSentence) {
				tmpSentence[idxWord+1] = strings.Title(tmpSentence[idxWord+1])
			}
		}
	}
	// 1st letter of 1st word with uppercase
	tmpSentence[0] = strings.Title(tmpSentence[0])
	// Last word folowed by a dot "."
	tmpSentence[len(tmpSentence)-1] = RemoveNonAlNum(tmpSentence[len(tmpSentence)-1]) + "."
	return strings.Join(tmpSentence, " ")
}

type countPunct struct {
	Str     []string
	Counted []int
	Max     []int
}

func (p *countPunct) Init() {
	p.Str = []string{",", ".", ";", ":", " ?", " !", " ?!,", "...,", " ..."}
	p.Max = []int{1, 2, 25, 25, 20, 20, 25, 28, 28}
	p.Counted = make([]int, len(p.Str))
}

func newLoremWords() []string {
	return []string{"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", "mauris", "sed", "lacus", "in", "urna", "sollicitudin", "blandit", "id", "at",
		"vestibulum", "commodo", "nisi", "accumsan", "nunc", "a", "facilisis", "congue", "arcu", "bibendum", "vitae", "eleifend", "ante", "et", "lectus", "quis", "rhoncus",
		"ligula", "nam", "magna", "sodales", "malesuada", "eget", "sapien", "lacinia", "purus", "quisque", "eu", "placerat", "fusce", "orci", "tristique", "nisl", "vivamus",
		"libero", "nec", "ultricies", "fermentum", "efficitur", "duis", "ullamcorper", "pulvinar", "suscipit", "sagittis", "neque", "tincidunt", "tortor", "lobortis", "aliquam",
		"volutpat", "viverra", "primis", "faucibus", "luctus", "ultrices", "posuere", "cubilia", "curae", "gravida", "aliquet", "nibh", "donec", "porta", "dignissim", "sem",
		"pellentesque", "scelerisque", "est", "augue", "dapibus", "felis", "quam", "ut", "aenean", "nulla", "mollis", "diam", "ac", "elementum", "risus", "turpis", "iaculis",
		"mi", "morbi", "tempus", "pharetra", "maecenas", "mattis", "feugiat", "odio", "curabitur", "enim", "phasellus", "hendrerit", "justo", "maximus", "varius", "natoque",
		"penatibus", "magnis", "dis", "parturient", "montes", "nascetur", "ridiculus", "mus", "vehicula", "suspendisse", "potenti", "habitant", "senectus", "netus", "erat",
		"fames", "egestas", "proin", "porttitor", "non", "semper", "cras", "imperdiet", "ornare", "dui", "facilisi", "tempor", "praesent", "massa", "metus", "auctor", "cursus",
		"condimentum", "tellus", "vulputate", "euismod", "convallis", "finibus", "velit", "vel", "rutrum", "molestie", "etiam", "interdum", "fringilla", "venenatis", "leo",
		"ex", "eros", "laoreet", "nullam", "dictum", "integer", "hac", "habitasse", "platea", "dictumst", "consequat", "pretium", "class", "aptent", "taciti", "sociosqu", "ad",
		"litora", "torquent", "per", "conubia", "nostra", "inceptos", "himenaeos"}
}
