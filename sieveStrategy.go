package main

import "sort"

type SieveStrategy struct {
	base            *wordList
	charsFreq       map[rune]int   // в скольки словах встречается каждая буква
	itemFreqIndexes map[string]int // сумма частот букв для каждого слова
	filter          *wordFilter
}

func (st *SieveStrategy) Init(base *wordList) {
	st.base = base

	st.charsFreq = make(map[rune]int)
	wordChars := newCharSet()
	for _, word := range st.base.items {
		wordChars.clear()
		for _, c := range word {
			if !wordChars.contains(c) {
				st.charsFreq[c]++
				wordChars.add(c)
			}
		}
	}

	st.itemFreqIndexes = make(map[string]int)
	for _, word := range st.base.items {
		wordChars.clear()
		for _, c := range word {
			if !wordChars.contains(c) {
				st.itemFreqIndexes[word] += st.charsFreq[c]
				wordChars.add(c)
			}
		}
	}

	st.filter = newWordFilter()
}

func (st *SieveStrategy) GetFirstMoveTries() ([]string, error) {
	return []string{"норка"}, nil

	// words := make([]string, st.base.count())
	// copy(words, st.base.items)
	// sortWordsByCharsFreq(words, st)
	// return words, nil
}

func (st *SieveStrategy) GetNextMoveTries(move int, lastTry,
	gameResponse string) ([]string, error) {

	err := st.filter.update(lastTry, gameResponse)
	if err != nil {
		return nil, newStratError(err, false)
	}

	words := []string{}
	for _, word := range st.base.items {
		ok, err := st.filter.checkWord(word)
		if err != nil {
			return nil, newStratError(err, true)
		}
		if ok {
			words = append(words, word)
		}
	}
	sortWordsByCharsFreq(words, st)

	return words, nil
}

type wordsInfo struct {
	words    []string
	strategy *SieveStrategy
}

func (wi *wordsInfo) Len() int { return len(wi.words) }

func (wi *wordsInfo) Less(i, j int) bool {
	f1 := wi.strategy.itemFreqIndexes[wi.words[i]]
	f2 := wi.strategy.itemFreqIndexes[wi.words[j]]
	return f1 >= f2
}

func (wi *wordsInfo) Swap(i, j int) {
	wi.words[i], wi.words[j] = wi.words[j], wi.words[i]
}

func sortWordsByCharsFreq(words []string, strategy *SieveStrategy) {
	info := &wordsInfo{words, strategy}
	sort.Sort(info)
}

func NewSieveStrategy() *SieveStrategy {
	return &SieveStrategy{}
}
