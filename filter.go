package main

import (
	"errors"
	"fmt"
)

type wordFilter struct {
	deadChars  charSet          // буквы, которых точно нет в слове
	badChars   [wordLen]charSet // буквы, которые не подходят для i-ой позиции
	reqChars   charSet          // буквы, которые точно есть в слове
	fixedChars [wordLen]rune    // отгаданные буквы (на своих позициях)
}

func (f *wordFilter) update(lastWord, answer string) error {
	wordChars := []rune(lastWord)
	answChars := []rune(answer)
	if len(wordChars) != wordLen || len(answChars) != wordLen {
		return errors.New("wrong word length")
	}
	for i, curChar := range wordChars {
		switch answChars[i] {
		case '+':
			f.fixedChars[i] = curChar
		case '-':
			f.deadChars.add(curChar)
		case '.':
			f.badChars[i].add(curChar)
			f.reqChars.add(curChar)
		default:
			return fmt.Errorf("unknown char \"%c\"", answChars[i])
		}
	}
	return nil
}

func (f *wordFilter) checkWord(word string) bool {
	wordChars := newCharSet()
	for i, curChar := range []rune(word) {
		if fixed := f.fixedChars[i]; fixed != 0 && curChar != fixed {
			return false
		}
		if f.deadChars.has(curChar) {
			return false
		}
		if f.badChars[i].has(curChar) {
			return false
		}
		wordChars.add(curChar)
	}
	if !wordChars.hasAll(f.reqChars) {
		return false
	}
	return true
}

func newWordFilter() *wordFilter {
	f := &wordFilter{}
	f.deadChars = newCharSet()
	for i := 0; i < wordLen; i++ {
		f.badChars[i] = newCharSet()
	}
	f.reqChars = newCharSet()
	return f
}
