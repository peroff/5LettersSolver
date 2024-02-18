package main

import (
	"testing"
)

func TestEmptyWordFilter(t *testing.T) {
	filter := newWordFilter()
	for _, w := range base.items {
		if !checkWord(t, filter, w) {
			t.Fatalf("word %q doesn't pass empty filter", w)
		}
	}
}

func TestWordFilter(t *testing.T) {
	for _, ts := range testSessions {
		filter := newWordFilter()

		if !checkWord(t, filter, ts.secret) {
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
			if w != ts.secret && checkWord(t, filter, w) {
				t.Fatalf("word %q passes filter after %q game session",
					w, ts.secret)
			}
		}
	}
}

func testTryWordFilter(t *testing.T, filter *wordFilter, secret, try,
	resp string, passedWords *int) {

	updateFilter(t, filter, try, resp)

	if !checkWord(t, filter, secret) {
		t.Fatalf("word %q doesn't pass filter after try %q", secret, try)
	}

	// после каждой попытки число подходящих слов должно быть не больше прежнего!
	p := getPassedWords(t, filter)
	if p > *passedWords {
		t.Fatalf("filter misbehavior. Last passed words: %d, now: %d. "+
			"Secret: %q, try: %q", *passedWords, p, secret, try)
	}
	*passedWords = p
}

func updateFilter(t *testing.T, filter *wordFilter, try, resp string) {
	err := filter.update(try, resp)
	if err != nil {
		t.Fatalf("filter update error: %s", err)
	}
}

func checkWord(t *testing.T, filter *wordFilter, word string) bool {
	ok, err := filter.checkWord(word)
	if err != nil {
		t.Fatalf("word checking error: %s", err)
	}
	return ok
}

func getPassedWords(t *testing.T, filter *wordFilter) int {
	n := 0
	for _, w := range base.items {
		if checkWord(t, filter, w) {
			n++
		}
	}

	return n
}

func TestWordFilterNegatives(t *testing.T) {
	var cases = []struct {
		try  string
		resp string
		test string
	}{
		{
			"azzzz",
			"+----",
			"zcccc",
		},
		{
			"azzzz",
			".----",
			"acccc",
		},
		{
			"azzzz",
			"-----",
			"acccc",
		},
		{
			"azzzz",
			".----",
			"ccccc",
		},
		{
			"aazzz",
			".----",
			"ccccc",
		},
		{
			"aazzz",
			".----",
			"cccaa",
		},
	}

	for _, c := range cases {
		filter := newWordFilter()
		updateFilter(t, filter, c.try, c.resp)
		if checkWord(t, filter, c.test) {
			t.Fatalf("wrong word passed: %q (after %q, %q)",
				c.test, c.try, c.resp)
		}
	}
}
