package main

const wordLen = 5

const (
	fixedCharAnsw = '+'
	badCharAnsw   = '.'
	deadCharAnsw  = '-'
)

// Предположительная функция генерации ответа игры (информации об угаданных
// буквах). С ее учетом строится алгоритм фильтрации слов.
// TODO rename answer -> responce
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

	// сначала обозначаем угаданные буквы, уменьшая счетчик оставшихся в слове букв
	for i, tc := range tryChars {
		if tc == secretChars[i] {
			answerChars[i] = fixedCharAnsw
			secretCharCount[tc]--
		}
	}

	// обозначаем буквы не на своих местах и отсутствующие, в зависимости от счетчика
	for i, tc := range tryChars {
		if tc != secretChars[i] {
			if secretCharCount[tc] > 0 {
				answerChars[i] = badCharAnsw
				secretCharCount[tc]--
			} else {
				answerChars[i] = deadCharAnsw
			}
		}
	}

	return string(answerChars)
}
