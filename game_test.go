package main

import (
	"strings"
	"testing"
)

func TestGetGameAnswer(t *testing.T) {
	for _, ts := range testSessions {
		for _, try := range ts.tries {
			testTryAnswer(t, ts.secret, try.try, try.res)
		}

		testTryAnswer(t, ts.secret, ts.secret, allFixedCharsAnsw)
		testTryAnswer(t, ts.secret, strings.Repeat(" ", wordLen),
			allDeadCharsAnsw)
	}
}

func testTryAnswer(t *testing.T, secret, try, want string) {
	if res := getGameAnswer(secret, try); res != want {
		t.Fatalf("%q, %q: want %q, got %q", secret, try, want, res)
	}
}
