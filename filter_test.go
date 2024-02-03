package main

import (
	"testing"
)

func TestEmptyWordFilter(t *testing.T) {
	filter := newWordFilter()
	for _, w := range base.items {
		if !filter.checkWord(w) {
			t.Fatalf("word %q doesn't pass empty filter", w)
		}
	}
}

func TestWordFilter(t *testing.T) {
	for _, ts := range testSessions {
		filter := newWordFilter()

		if !filter.checkWord(ts.secret) {
			t.Fatalf("word %q doesn't pass empty filter", ts.secret)
		}

		passedWords := base.count()

		for _, try := range ts.tries {
			testTryWordFilter(t, filter, ts.secret, try.try, try.resp,
				&passedWords)
		}

		testTryWordFilter(t, filter, ts.secret, ts.secret, allFixedCharsResp,
			&passedWords)

		for _, w := range base.items {
			if w != ts.secret && filter.checkWord(w) {
				t.Fatalf("word %q passes filter after %q game session",
					w, ts.secret)
			}
		}
	}
}

func testTryWordFilter(t *testing.T, filter *wordFilter, secret, try, res string,
	passedWords *int) {

	err := filter.update(try, res)
	if err != nil {
		t.Fatalf("filter update error: %s", err)
	}

	if !filter.checkWord(secret) {
		t.Fatalf("word %q doesn't pass filter after try %q", secret, try)
	}

	// после каждой попытки число подходящих слов должно быть не больше прежнего!
	p := getPassedWords(filter)
	if p > *passedWords {
		t.Fatalf("filter misbehavior. Last passed words: %d, now: %d. "+
			"Secret: %q, try: %q", *passedWords, p, secret, try)
	}
	*passedWords = p
}

func getPassedWords(f *wordFilter) int {
	n := 0
	for _, w := range base.items {
		if f.checkWord(w) {
			n++
		}
	}

	return n
}
