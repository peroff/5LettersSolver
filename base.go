package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

type wordsBase struct {
	items    []string
	runeFreq map[rune]float64
}

func loadBase(fileName string) (*wordsBase, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	base := &wordsBase{
		items:    strings.Split(string(b), "\n"),
		runeFreq: make(map[rune]float64),
	}

	chars := make(map[rune]int)
	total := 0
	for _, word := range base.items {
		for _, r := range word {
			chars[r]++
			total++
		}
	}
	if total == 0 {
		return nil, errors.New("empty base")
	}

	for c, n := range chars {
		base.runeFreq[c] = float64(n) / float64(total)
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
