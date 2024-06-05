package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	releaseVersion = "0.6"
	releaseYear    = "2024"
)

const (
	wordsFile   = "words.txt"
	maxWords    = 50
	wordsInLine = 10
	removingCmd = "!!!" // префикс, активирующий функцию удаления слов из базы, например можно ввести: "!!! ёпрст"
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
	return "норка"
}

func selectWords(base *wordsBase, filter *wordFilter) ([]string, error) {
	res := []string{}
	for _, word := range base.items {
		ok, err := filter.checkWord(word)
		if err != nil {
			return nil, err
		}
		if ok {
			res = append(res, word)
		}
	}
	sortWordsByCharsFreq(res, base)
	return res, nil
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
	fmt.Printf("(%d всего, %d показано)\n", total, len(words))
}

func main() {
	fmt.Printf("Words5Solver v%s (c) Dan Peroff, 2022-%s\n",
		releaseVersion, releaseYear)
	fmt.Println()

	base, err := loadBase(wordsFile)
	if err != nil {
		fmt.Printf("Ошибка при загрузке базы слов: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Загружено слов: %d\n\n", base.count())

	filter := newWordFilter()
	input := bufio.NewScanner(os.Stdin)

	move := 1
	currentWord := ""
	defaultWord := getStartWord(base)
	waitingForResponse := false

mainLp:
	for {
		if !waitingForResponse {
			if move == 1 {
				fmt.Printf("%d. Введите слово, с которого начинаем (по умолчанию: \"%s\"): ",
					move, defaultWord)
			} else {
				fmt.Printf("%d. Введите выбранное вами слово (по умолчанию: \"%s\"): ",
					move, defaultWord)
			}
		} else {
			fmt.Printf("%d. Введите ответ приложения (5 символов: '+' - буква отгадана,\n"+
				"  '-' - отсутствует, '.' - не на своем месте), пустой для выхода: ", move)
		}

		if !input.Scan() {
			break // ошибка
		}
		s := strings.TrimSpace(input.Text())
		if s == "" {
			if !waitingForResponse {
				s = defaultWord
			} else {
				break
			}
		}

		if strings.HasPrefix(s, removingCmd) {
			removeWordsFromBase(base, strings.TrimPrefix(s, removingCmd))
			continue
		}

		if utf8.RuneCountInString(s) != wordLen {
			fmt.Printf("Неверное число символов\n\n")
			continue
		}

		if !waitingForResponse {
			s = normalizeWord(s)
			if !base.hasWord(s) {
				fmt.Printf("Неизвестное слово \"%s\"\n\n", s)
				continue
			}
			currentWord = s
			waitingForResponse = true
			fmt.Printf("Выбрано слово \"%s\", введите его в приложении игры.\n",
				currentWord)
		} else {
			if err := filter.update(currentWord, s); err != nil {
				fmt.Printf("Некорректный ответ: %s\n\n", err)
				continue
			}
			move++
			words, err := selectWords(base, filter)
			if err != nil {
				fmt.Printf("Упс! Непредвиденная ошибка: %s\n", err)
				os.Exit(1)
			}
			switch len(words) {
			case 0:
				fmt.Printf("\nНе найдено подходящих слов :( Сожалею...\n\n")
				fmt.Print("Нажмите ENTER для выхода")
				input.Scan()
				break mainLp
			case 1:
				fmt.Printf("\nНАШЛИ! Ваше слово: \"%s\"\n\n", words[0])
				fmt.Print("Нажмите ENTER для выхода")
				input.Scan()
				break mainLp
			default:
				fmt.Printf("\n%d. Возможные слова:\n", move)
				printWords(words)
				fmt.Println()
				defaultWord = words[0]
			}

			waitingForResponse = false
		}
	}
	if err := input.Err(); err != nil {
		panic(fmt.Sprintf("input scanning error: %s", err))
	}
}

func removeWordsFromBase(base *wordsBase, words string) {
	n := 0
	for _, w := range strings.Split(normalizeWord(words), " ") {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		if base.removeWord(w) {
			fmt.Printf("%s Удалено: \"%s\"\n", removingCmd, w)
			n++
		} else {
			fmt.Printf("%s Не найдено: \"%s\"\n", removingCmd, w)
		}
	}

	err := base.save(wordsFile)
	if err != nil {
		fmt.Printf("%s Ошибка сохранения базы: %s\n", removingCmd, err)
		return
	}

	fmt.Printf("%s Успешно удалено слов: %d\n", removingCmd, n)
	fmt.Println()
}
