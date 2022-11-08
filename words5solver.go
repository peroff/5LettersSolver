package main

import (
	"bufio"
	"bytes"
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

type charSet map[rune]struct{}

type wordFilter struct {
	deadChars  charSet
	badChars   [wordLen]charSet
	reqChars   charSet
	fixedChars [wordLen]rune
}

func (f *wordFilter) checkWord(word string) bool {
	chars := make(charSet)
	for i, c := range []rune(word) {
		if fc := f.fixedChars[i]; fc != 0 && c != fc {
			return false
		}
		if _, ok := f.deadChars[c]; ok {
			return false
		}
		if _, ok := f.badChars[i][c]; ok {
			return false
		}
		chars[c] = struct{}{}
	}
	for rc, _ := range f.reqChars {
		if _, ok := chars[rc]; !ok {
			return false
		}
	}
	return true
}

func rebuildBase() {
	b, err := ioutil.ReadFile("raw_words.txt")
	if err != nil {
		log.Fatal(err)
	}

	var base bytes.Buffer
	n := 0

	text := string(b)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		words := strings.Fields(line)
		for _, word := range words {
			word = strings.TrimSpace(word)
			if word == "" || word[len(word)-1] == '.' {
				continue
			}
			word = strings.ToLower(word)
			if n > 0 {
				base.WriteRune('\n')
			}
			base.WriteString(word)
			n++
		}
	}

	err = ioutil.WriteFile("words.txt", base.Bytes(), 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Words written: %d\n", n)
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
			filter.deadChars[lwc] = struct{}{}
		case '?':
			filter.badChars[i][lwc] = struct{}{}
			filter.reqChars[lwc] = struct{}{}
		default:
			return fmt.Errorf("unknown char \"%c\"", c)
		}
	}
	return nil
}

func main() {
	// rebuildBase()

	base, err := loadBase("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded words: %d\n", len(base.items))

	filter := &wordFilter{}
	filter.deadChars = make(charSet)
	for i := 0; i < wordLen; i++ {
		filter.badChars[i] = make(charSet)
	}
	filter.reqChars = make(charSet)

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
