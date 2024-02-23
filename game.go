package main

import "fmt"

const wordLen = 5

const (
	fixedCharResp = '+'
	badCharResp   = '.'
	deadCharResp  = '-'
)

// Предположительная функция генерации ответа игры (информации об угаданных
// буквах). С ее учетом строится алгоритм фильтрации слов.
func getGameResponse(secret, try string) (string, error) {
	// результат работы
	respChars := make([]rune, wordLen)

	// массивы символов загаданного слова и предположения
	secretChars := []rune(secret)
	if len(secretChars) != wordLen {
		return "", fmt.Errorf("неверная длина слова: %q (%d)",
			secret, len(secretChars))
	}
	tryChars := []rune(try)
	if len(tryChars) != wordLen {
		return "", fmt.Errorf("неверная длина слова: %q (%d)",
			try, len(tryChars))
	}

	// число вхождений каждой буквы в загаданное слово
	secretCharCount := make(map[rune]int)
	for _, sc := range secretChars {
		secretCharCount[sc]++
	}

	// сначала обозначаем угаданные буквы, уменьшая счетчик оставшихся в слове букв
	for i, tc := range tryChars {
		if tc == secretChars[i] {
			respChars[i] = fixedCharResp
			secretCharCount[tc]--
		}
	}

	// обозначаем буквы не на своих местах и отсутствующие, в зависимости от счетчика
	for i, tc := range tryChars {
		if tc != secretChars[i] {
			if secretCharCount[tc] > 0 {
				respChars[i] = badCharResp
				secretCharCount[tc]--
			} else {
				respChars[i] = deadCharResp
			}
		}
	}

	return string(respChars), nil
}
