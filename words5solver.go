package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const version = "0.1"

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
	fmt.Printf("Words5Solver v%s (c) Dan Peroff, 2022\n", version)
	fmt.Println()

	base, err := loadBase("words.txt")
	if err != nil {
		fmt.Printf("Words base loading error: %s\n", err)
		return
	}
	fmt.Printf("Loaded words: %d\n\n", len(base.items))

	filter := newWordFilter()
	input := bufio.NewScanner(os.Stdin)

	move := 1
	currentWord := getFirstWord(base)
	fmt.Printf("%d. Start with word: %s\n", move, currentWord)
	waitingForResponse := true

	for {
		if waitingForResponse {
			fmt.Printf("%d. Enter app's response, 5 symbols: '+' - correct letter, '-' - wrong letter,\n", move)
			fmt.Printf("   '?', '*' or '.' - misplaced letter. Response (empty for exit): ")
		} else {
			fmt.Printf("%d. Enter your next word (same there and in the app): ", move)
		}

		if !input.Scan() {
			break
		}
		s := strings.TrimSpace(input.Text())
		if s == "" {
			break
		}
		if utf8.RuneCountInString(s) != wordLen {
			fmt.Printf("Wrong input length\n\n")
			continue
		}

		if waitingForResponse {
			if err := filter.update(currentWord, s); err != nil {
				fmt.Printf("Wrong filter: %s\n\n", err)
				continue
			}
			words := selectWords(base, filter)
			move++
			fmt.Println()
			fmt.Printf("%d. Possible words (%d):\n", move, len(words))
			for _, w := range words {
				fmt.Println(w)
			}
			fmt.Println()
			waitingForResponse = false
		} else {
			currentWord = s
			waitingForResponse = true
		}
	}
	if err := input.Err(); err != nil {
		panic(fmt.Sprintf("input scanning error: %s", err))
	}
}
