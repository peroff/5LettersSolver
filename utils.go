package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func newWordLenError(word string) error {
	return fmt.Errorf("неверная длина слова: \"%s\" (%d)",
		word, utf8.RuneCountInString(word))
}

func normalizeWord(word string) string {
	return strings.ReplaceAll(strings.ToLower(word), "ё", "е")
}

func capitalizeFirstWord(text string) string {
	words := strings.Split(text, " ")
	if len(words) > 0 {
		words[0] = strings.Title(words[0])
	}
	return strings.Join(words, " ")
}
