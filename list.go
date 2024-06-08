package main

import (
	"errors"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

type wordList struct {
	items           []string
	charsFreq       map[rune]int   // в скольки словах встречается каждая буква
	itemFreqIndexes map[string]int // сумма частот букв для каждого слова
}

func (wl *wordList) contains(word string) bool {
	for _, w := range wl.items {
		if w == word {
			return true
		}
	}
	return false
}

func (wl *wordList) remove(word string) bool {
	for i := range wl.items {
		if wl.items[i] == word {
			l := len(wl.items)
			copy(wl.items[i:l-1], wl.items[i+1:l])
			wl.items = wl.items[:l-1]
			return true
		}
	}
	return false
}

func (wl *wordList) count() int {
	return len(wl.items)
}

func (wl *wordList) load(fileName string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	text := normalizeWord(string(b))
	words := strings.Split(text, "\n")
	for i := range words {
		words[i] = strings.TrimSpace(words[i])
		if wlen := utf8.RuneCountInString(words[i]); wlen != wordLen {
			return wordLenError(words[i])
		}
	}

	wl.items = words
	wl.charsFreq = make(map[rune]int)
	wl.itemFreqIndexes = make(map[string]int)

	if len(wl.items) == 0 {
		return errors.New("в файле нет ни одного слова")
	}

	wordChars := newCharSet()
	for _, word := range wl.items {
		wordChars.clear()
		for _, c := range word {
			if !wordChars.has(c) {
				wl.charsFreq[c]++
				wordChars.add(c)
			}
		}
	}

	for _, word := range wl.items {
		wordChars.clear()
		for _, c := range word {
			if !wordChars.has(c) {
				wl.itemFreqIndexes[word] += wl.charsFreq[c]
				wordChars.add(c)
			}
		}
	}

	return nil
}

func (wl *wordList) save(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	words := make([]string, wl.count())
	copy(words, wl.items)
	sort.Strings(words)
	_, err = f.WriteString(strings.Join(words, "\r\n"))
	if err != nil {
		os.Remove(fileName)
		return err
	}

	return nil
}

func newWordList() *wordList {
	return &wordList{
		items:           make([]string, 0),
		charsFreq:       make(map[rune]int),
		itemFreqIndexes: make(map[string]int),
	}
}
