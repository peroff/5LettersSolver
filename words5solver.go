package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

const wordLen = 5

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

func main() {
	base, err := loadBase("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded words: %d\n", len(base.items))

	filter := newWordFilter()

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
			err := filter.update(lastWord, s)
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
