package main

import (
	"errors"
	"fmt"
	"sort"
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
		return errors.New("неверная длина слова")
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
		return false, fmt.Errorf("неверная длина слова: %q (%d)",
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

	min := charCntToStr(f.minCharCnt)

	cnt := charCntToStr(f.charCnt)

	return fmt.Sprintf("fixed: \"%s\"; bad: %s; min: %s; cnt: %s",
		fixed, bad, min, cnt)
}

func charCntToStr(cnt map[rune]int) string {
	if len(cnt) == 0 {
		return ""
	}

	chars := make([]rune, 0, len(cnt))
	for c, _ := range cnt {
		chars = append(chars, c)
	}
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

	s := ""
	for i, c := range chars {
		if i < len(chars)-1 {
			s += fmt.Sprintf("%c:%d, ", c, cnt[c])
		} else {
			s += fmt.Sprintf("%c:%d", c, cnt[c])
		}
	}

	return s
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
