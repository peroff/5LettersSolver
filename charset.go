package main

type charSet map[rune]struct{}

func (cs charSet) clear() {
	for c, _ := range cs {
		delete(cs, c)
	}
}

func (cs charSet) add(char rune) {
	cs[char] = struct{}{}
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

func newCharSet() charSet {
	return make(map[rune]struct{})
}
