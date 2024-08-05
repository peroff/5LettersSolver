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
	items []string
}

func (wl *wordList) add(word string) {
	wl.items = append(wl.items, word)
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
	if len(wl.items) == 0 {
		return errors.New("в файле нет ни одного слова")
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
	return &wordList{}
}
