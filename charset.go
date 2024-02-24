package main

import (
	"fmt"
	"sort"
)

type charSet map[rune]struct{}

func (cs charSet) clear() {
	for c := range cs {
		delete(cs, c)
	}
}

func (cs charSet) add(char rune) {
	cs[char] = struct{}{}
}

func (cs charSet) count() int {
	return len(cs)
}

func (cs charSet) has(char rune) bool {
	_, ok := cs[char]
	return ok
}

func (cs charSet) hasAll(chars charSet) bool {
	for c, _ := range chars {
		if !cs.has(c) {
			return false
		}
	}
	return true
}

func (cs charSet) String() string {
	if len(cs) == 0 {
		return "[]"
	}

	chars := make([]rune, 0, len(cs))
	for c, _ := range cs {
		chars = append(chars, c)
	}
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

	s := "["
	for i, c := range chars {
		if i < len(chars)-1 {
			s += fmt.Sprintf("%c, ", c)
		} else {
			s += fmt.Sprintf("%c]", c)
		}
	}

	return s
}

func newCharSet() charSet {
	return make(map[rune]struct{})
}
