package main

import (
	"fmt"
	"unicode/utf8"
)

func wordLenError(word string) error {
	return fmt.Errorf("неверная длина слова: \"%s\" (%d)",
		word, utf8.RuneCountInString(word))
}
