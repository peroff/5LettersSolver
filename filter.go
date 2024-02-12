package main

import (
	"errors"
	"fmt"
)

type wordFilter struct {
	fixedChars [wordLen]rune    // отгаданные буквы (на своих позициях)
	badChars   [wordLen]charSet // буквы, которые не подходят для i-ой позиции
	charCnt    map[rune]int     // точное число вхождений каждой известной буквы
	minCharCnt map[rune]int     // минимально необходимое число вхождений каждой известной буквы
}

func (f *wordFilter) update(try, response string) error {
	tryChars := []rune(try)
	respChars := []rune(response)
	if len(tryChars) != wordLen || len(respChars) != wordLen {
		return errors.New("wrong word length")
	}

	openCharCnt := make(map[rune]int)

	for i, rc := range respChars {
		if rc == fixedCharResp {
			f.fixedChars[i] = tryChars[i]
			openCharCnt[tryChars[i]]++
		}
	}

	for i, rc := range respChars {
		switch rc {
		case fixedCharResp:
			// ничего не делаем
		case badCharResp:
			f.badChars[i].add(tryChars[i])
			openCharCnt[tryChars[i]]++
		case deadCharResp:
			f.badChars[i].add(tryChars[i])
			// Если нашли минус, значит ВСЕ вхождения данной буквы в загаданное
			// слово (если они есть) обозначены в ответе в виде плюсов и точек,
			// уже подсчитаны, и можно сделать вывод о фактическом их числе.
			// Не может быть чтобы одна и та же буква сначала была минусом,
			// а потом (правее) точкой.
			f.charCnt[tryChars[i]] = openCharCnt[tryChars[i]]
		default:
			return fmt.Errorf("unknown char \"%c\"", rc)
		}
	}

	for wc, cnt := range openCharCnt {
		if cnt > f.minCharCnt[wc] {
			f.minCharCnt[wc] = cnt
		}
	}

	return nil
}

func (f *wordFilter) checkWord(word string) (bool, error) {
	wordChars := []rune(word)
	if len(wordChars) != wordLen {
		return false, fmt.Errorf("wrong word length: %q (%d)",
			word, len(wordChars))
	}

	charCount := make(map[rune]int)

	for i, wc := range wordChars {
		if fixed := f.fixedChars[i]; fixed != 0 && wc != fixed {
			return false, nil
		}
		if f.badChars[i].has(wc) {
			return false, nil
		}
		charCount[wc]++
	}

	for wc, cnt := range charCount {
		if actCnt, ok := f.charCnt[wc]; ok && cnt != actCnt {
			return false, nil
		}
		if min, ok := f.minCharCnt[wc]; ok && cnt < min {
			return false, nil
		}
	}

	return true, nil
}

func newWordFilter() *wordFilter {
	f := &wordFilter{
		charCnt:    make(map[rune]int),
		minCharCnt: make(map[rune]int),
	}
	for i := 0; i < wordLen; i++ {
		f.badChars[i] = newCharSet()
	}
	return f
}
