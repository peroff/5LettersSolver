package main

import (
	"fmt"
)

type wordFilter struct {
	fixedChars [wordLen]rune    // отгаданные буквы (на своих позициях)
	badChars   [wordLen]charSet // буквы, которые не подходят для i-ой позиции
	charCnt    charCounter      // точное число вхождений каждой известной буквы
	minCharCnt charCounter      // минимально необходимое число вхождений каждой известной буквы
}

func (f *wordFilter) update(try, response string) error {
	tryChars := []rune(try)
	if len(tryChars) != wordLen {
		return wordLenError(try)
	}
	respChars := []rune(response)
	if len(respChars) != wordLen {
		return wordLenError(response)
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
			return fmt.Errorf("неизвестный символ \"%c\"", rc)
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
		return false, wordLenError(word)
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

	for c, min := range f.minCharCnt {
		if charCount[c] < min {
			return false, nil
		}
	}

	for c, cnt := range f.charCnt {
		if charCount[c] != cnt {
			return false, nil
		}
	}

	return true, nil
}

func (f *wordFilter) String() string {
	fixed := ""
	for _, f := range f.fixedChars {
		if f != 0 {
			fixed += string(f)
		} else {
			fixed += "_"
		}
	}

	bad := ""
	for i, b := range f.badChars {
		if i < wordLen-1 {
			bad += b.String() + ", "
		} else {
			bad += b.String()
		}
	}

	return fmt.Sprintf("fixed: \"%s\"; bad: %s; min: %s; cnt: %s",
		fixed, bad, f.minCharCnt, f.charCnt)
}

func newWordFilter() *wordFilter {
	f := &wordFilter{
		charCnt:    newCharCounter(),
		minCharCnt: newCharCounter(),
	}
	for i := 0; i < wordLen; i++ {
		f.badChars[i] = newCharSet()
	}
	return f
}
