package main

import (
	"errors"
	"fmt"
)

type wordFilter struct {
	deadChars  charSet
	badChars   [wordLen]charSet
	reqChars   charSet
	fixedChars [wordLen]rune
}

func (f *wordFilter) update(lastWord, answer string) error {
	lwChars := []rune(lastWord)
	aChars := []rune(answer)
	if len(lwChars) != wordLen || len(aChars) != wordLen {
		return errors.New("wrong word length")
	}
	for i, lwc := range lwChars {
		switch aChars[i] {
		case '+':
			f.fixedChars[i] = lwc
		case '-':
			f.deadChars.add(lwc)
		case '?', '*', '.':
			f.badChars[i].add(lwc)
			f.reqChars.add(lwc)
		default:
			return fmt.Errorf("unknown char \"%c\"", aChars[i])
		}
	}
	return nil
}

func (f *wordFilter) checkWord(word string) bool {
	chars := newCharSet()
	for i, c := range []rune(word) {
		if fc := f.fixedChars[i]; fc != 0 && c != fc {
			return false
		}
		if f.deadChars.has(c) {
			return false
		}
		if f.badChars[i].has(c) {
			return false
		}
		chars.add(c)
	}
	if !chars.hasAll(f.reqChars) {
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
