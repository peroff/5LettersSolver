package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

const wordLen = 5

type wordsBase struct {
	items    []string
	runeFreq map[rune]float64
}

type wordFilter struct {
	deadChars  charSet
	badChars   [wordLen]charSet
	reqChars   charSet
	fixedChars [wordLen]rune
}

func (f *wordFilter) checkWord(word string) bool {
	chars := newCharSet()
	for i, c := range []rune(word) {
		if fc := f.fixedChars[i]; fc != 0 && c != fc {
			return false
		}
		if f.deadChars.has(c) {
			return false
		}
		if f.badChars[i].has(c) {
			return false
		}
		chars.add(c)
	}
	if !chars.hasAll(f.reqChars) {
		return false
	}
	return true
}

func loadBase(fileName string) (*wordsBase, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	base := &wordsBase{
		items:    strings.Split(string(b), "\n"),
		runeFreq: make(map[rune]float64),
	}

	chars := make(map[rune]int)
	total := 0
	for _, word := range base.items {
		for _, r := range word {
			chars[r]++
			total++
		}
	}
	if total == 0 {
		return nil, errors.New("empty base")
	}

	for c, n := range chars {
		base.runeFreq[c] = float64(n) / float64(total)
	}

	return base, nil
}

func getFirstWord(base *wordsBase) string {
	return "буква"
}

func selectWords(base *wordsBase, filter *wordFilter) []string {
	res := []string{}
	for _, word := range base.items {
		if filter.checkWord(word) {
			res = append(res, word)
		}
	}
	return res
}

func updateWordFilter(filter *wordFilter, lastWord, answer string) error {
	lwChars := []rune(lastWord)
	for i, c := range []rune(answer) {
		lwc := lwChars[i]
		switch c {
		case '+':
			filter.fixedChars[i] = lwc
		case '-':
			filter.deadChars.add(lwc)
		case '?':
			filter.badChars[i].add(lwc)
			filter.reqChars.add(lwc)
		default:
			return fmt.Errorf("unknown char \"%c\"", c)
		}
	}
	return nil
}

func main() {
	base, err := loadBase("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded words: %d\n", len(base.items))

	filter := &wordFilter{}
	filter.deadChars = newCharSet()
	for i := 0; i < wordLen; i++ {
		filter.badChars[i] = newCharSet()
	}
	filter.reqChars = newCharSet()

	move := 1
	firstWord := getFirstWord(base)
	fmt.Printf("%d. Start word: %s\n", move, firstWord)
	lastWord := firstWord

	fmt.Printf("%d. Answer: ", move)
	waitingForAnswer := true

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		s := strings.TrimSpace(input.Text())
		if s == "" {
			break
		}
		if utf8.RuneCountInString(s) != wordLen {
			fmt.Println("Wrong input")
			continue
		}

		if waitingForAnswer {
			err := updateWordFilter(filter, lastWord, s)
			if err != nil {
				fmt.Printf("Wrong filter: %s\n", err)
				continue
			}
			words := selectWords(base, filter)
			move++
			fmt.Printf("%d. Possible words (%d):\n", move, len(words))
			for _, w := range words {
				fmt.Println(w)
			}
			fmt.Printf("%d. Your word: ", move)
			waitingForAnswer = false
		} else {
			lastWord = s
			fmt.Printf("%d. Answer: ", move)
			waitingForAnswer = true
		}
	}
	err = input.Err()
	if err != nil {
		panic(fmt.Sprintf("input scanning error: %s\n", err))
	}
}
