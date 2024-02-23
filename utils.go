package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func wordLenError(word string) error {
	return fmt.Errorf("неверная длина слова: \"%s\" (%d)",
		word, utf8.RuneCountInString(word))
}

func normalizeWord(word string) string {
	return strings.ReplaceAll(strings.ToLower(word), "ё", "е")
}
