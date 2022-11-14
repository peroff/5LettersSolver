package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

type wordsBase struct {
	items           []string
	charsFreq       map[rune]float64
	itemFreqIndexes map[string]float64
}

func loadBase(fileName string) (*wordsBase, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	base := &wordsBase{
		items:           strings.Split(string(b), "\n"),
		charsFreq:       make(map[rune]float64),
		itemFreqIndexes: make(map[string]float64),
	}

	if len(base.items) == 0 {
		return nil, errors.New("empty base")
	}

	wordsByChar := make(map[rune]int)
	wordChars := newCharSet()
	for _, word := range base.items {
		wordChars.clear()
		for _, c := range word {
			if !wordChars.has(c) {
				wordsByChar[c]++
				wordChars.add(c)
			}
		}
	}

	for c, n := range wordsByChar {
		base.charsFreq[c] = float64(n) / float64(len(base.items))
	}

	for _, word := range base.items {
		wordChars.clear()
		for _, c := range word {
			if !wordChars.has(c) {
				base.itemFreqIndexes[word] += base.charsFreq[c]
				wordChars.add(c)
			}
		}
	}

	return base, nil
}

func (b *wordsBase) hasWord(word string) bool {
	for _, w := range b.items {
		if w == word {
			return true
		}
	}
	return false
}
