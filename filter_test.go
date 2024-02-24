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

func TestWordFilterClear(t *testing.T) {
	filter := newWordFilter()
	testWordFilterClearCase(t, filter, "abcde", "+++++")
	testWordFilterClearCase(t, filter, "abcde", ".....")
	testWordFilterClearCase(t, filter, "abcde", "-----")
	testWordFilterClearCase(t, filter, "abcde", "+-+-+")
	testWordFilterClearCase(t, filter, "abcde", "-.-.-")
	testWordFilterClearCase(t, filter, "abcde", ".+.+.")
}

func testWordFilterClearCase(t *testing.T, filter *wordFilter, try, resp string) {
	updateFilter(t, filter, try, resp)
	filter.clear()

	for i := 0; i < wordLen; i++ {
		if filter.fixedChars[i] != 0 {
			t.Fatalf("fixedChars[%d] isn't zero after clear()", i)
		}
		if filter.badChars[i].count() != 0 {
			t.Fatalf("badChars[%d] isn't empty after clear()", i)
		}
	}
	if filter.minCharCnt.count() != 0 {
		t.Fatalf("minCharCnt isn't empty after clear()")
	}
	if filter.charCnt.count() != 0 {
		t.Fatalf("charCnt isn't empty after clear()")
	}
}

func TestWordFilter(t *testing.T) {
	filter := newWordFilter()
	for _, ts := range testSessions {
		filter.clear()

		if !checkWord(t, filter, ts.secret) {
			t.Fatalf("word %q doesn't pass empty filter", ts.secret)
		}

		passedWords := base.count()

		for _, try := range ts.tries {
			testWordFilterTry(t, filter, ts.secret, try.try, try.resp,
				&passedWords)
		}

		testWordFilterTry(t, filter, ts.secret, ts.secret, allFixedCharsResp,
			&passedWords)

		for _, w := range base.items {
			if w != ts.secret && checkWord(t, filter, w) {
				t.Fatalf("word %q passes filter after %q game session",
					w, ts.secret)
			}
		}
	}
}

func testWordFilterTry(t *testing.T, filter *wordFilter, secret, try,
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

	filter := newWordFilter()
	for _, c := range cases {
		filter.clear()
		updateFilter(t, filter, c.try, c.resp)
		if checkWord(t, filter, c.test) {
			t.Fatalf("wrong word passed: %q (after %q, %q)",
				c.test, c.try, c.resp)
		}
	}
}
