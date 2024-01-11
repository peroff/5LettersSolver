package main

const wordLen = 5

const (
	fixedCharAnsw = '+'
	badCharAnsw   = '.'
	deadCharAnsw  = '-'
)

// Предположительная функция генерации ответа игры (информации об угаданных
// буквах). С ее учетом строится алгоритм фильтрации слов.
func getGameAnswer(secret, try string) string {
	// результат работы
	answerChars := make([]rune, wordLen)

	// массивы символов загаданного слова и предположения
	secretChars := []rune(secret)
	tryChars := []rune(try)
	if len(secretChars) != wordLen || len(tryChars) != wordLen {
		panic("bad word length")
	}

	// число вхождений каждой буквы в загаданное слово
	secretCharCount := make(map[rune]int)
	for _, sc := range secretChars {
		secretCharCount[sc]++
	}

	// сколько раз каждую букву уже раскрыли текущим предположением
	openCharCount := make(map[rune]int)

	for i, tc := range tryChars {
		if tc == secretChars[i] {
			answerChars[i] = fixedCharAnsw
			openCharCount[tc]++
		} else {
			if openCharCount[tc] < secretCharCount[tc] {
				answerChars[i] = badCharAnsw
				openCharCount[tc]++
			} else {
				answerChars[i] = deadCharAnsw
			}
		}
	}

	return string(answerChars)
}
