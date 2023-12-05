package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

const version = "0.4"

const (
	wordLen     = 5
	maxWords    = 200
	wordsInLine = 10
)

type wordsInfo struct {
	words []string
	base  *wordsBase
}

func (wi *wordsInfo) Len() int { return len(wi.words) }

func (wi *wordsInfo) Less(i, j int) bool {
	f1 := wi.base.itemFreqIndexes[wi.words[i]]
	f2 := wi.base.itemFreqIndexes[wi.words[j]]
	return f1 >= f2
}

func (wi *wordsInfo) Swap(i, j int) {
	wi.words[i], wi.words[j] = wi.words[j], wi.words[i]
}

func sortWordsByCharsFreq(words []string, base *wordsBase) {
	info := &wordsInfo{words, base}
	sort.Sort(info)
}

func getStartWord(base *wordsBase) string {
	words := make([]string, len(base.items))
	copy(words, base.items)
	sortWordsByCharsFreq(words, base)
	// printWords(words)
	return words[0]
}

func selectWords(base *wordsBase, filter *wordFilter) []string {
	res := []string{}
	for _, word := range base.items {
		if filter.checkWord(word) {
			res = append(res, word)
		}
	}
	sortWordsByCharsFreq(res, base)
	return res
}

func printWords(words []string) {
	total := len(words)
	if total > maxWords {
		words = words[:maxWords]
	}
	for offs := 0; offs < len(words); offs += wordsInLine {
		cnt := len(words) - offs
		if cnt > wordsInLine {
			cnt = wordsInLine
		}
		fmt.Printf("  %s\n", strings.Join(words[offs:offs+cnt], ", "))
	}
	fmt.Printf("(%d total, %d shown)\n", total, len(words))
}

func main() {
	fmt.Printf("Words5Solver v%s (c) Dan Peroff, 2022-2023\n", version)
	fmt.Println()

	base, err := loadBase("words.txt")
	if err != nil {
		fmt.Printf("Words base loading error: %s\n", err)
		return
	}
	fmt.Printf("Loaded words: %d\n\n", base.count())

	filter := newWordFilter()
	input := bufio.NewScanner(os.Stdin)

	move := 1
	currentWord := ""
	waitingForResponse := false

mainLp:
	for {
		if waitingForResponse {
			fmt.Printf("%d. Enter app's response, 5 symbols: '+' - correct letter, '-' - wrong letter,\n", move)
			fmt.Printf("   '.' - misplaced letter. Response (empty for exit): ")
		} else {
			if move == 1 {
				fmt.Printf("%d. Enter your first word (recommended: \"%s\"): ",
					move, getStartWord(base))
			} else {
				fmt.Printf("%d. Enter your next word (same there and in the app): ", move)
			}
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
			move++
			words := selectWords(base, filter)
			switch len(words) {
			case 0:
				fmt.Printf("\nNo possible words found :( Sorry...\n\n")
				fmt.Print("Press ENTER for exit")
				input.Scan()
				break mainLp
			case 1:
				fmt.Printf("\nFOUND! Your word: [%s]\n\n", words[0])
				fmt.Print("Press ENTER for exit")
				input.Scan()
				break mainLp
			default:
				fmt.Printf("\n%d. Possible words:\n", move)
				printWords(words)
				fmt.Println()
			}
			waitingForResponse = false
		} else {
			if !base.hasWord(s) {
				fmt.Printf("Unknown word \"%s\"\n\n", s)
				continue
			}
			currentWord = s
			waitingForResponse = true
		}
	}
	if err := input.Err(); err != nil {
		panic(fmt.Sprintf("input scanning error: %s", err))
	}
}
