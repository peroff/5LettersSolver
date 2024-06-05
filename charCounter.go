package main

import (
	"fmt"
	"sort"
)

type charCounter map[rune]int

func newCharCounter() charCounter {
	return make(map[rune]int)
}

func (cc charCounter) String() string {
	if len(cc) == 0 {
		return ""
	}

	chars := make([]rune, 0, len(cc))
	for c, _ := range cc {
		chars = append(chars, c)
	}
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

	s := ""
	for i, c := range chars {
		if i < len(chars)-1 {
			s += fmt.Sprintf("%c:%d, ", c, cc[c])
		} else {
			s += fmt.Sprintf("%c:%d", c, cc[c])
		}
	}

	return s
}

func (cc charCounter) clear() {
	for c := range cc {
		delete(cc, c)
	}
}

func (cc charCounter) count() int {
	return len(cc)
}
