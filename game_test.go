package main

import (
	"strings"
	"testing"
)

var nonExistentWord = strings.Repeat(" ", wordLen)

func TestGetGameResponse(t *testing.T) {
	for _, ts := range testSessions {
		for _, try := range ts.tries {
			testTryResponse(t, ts.secret, try.try, try.resp)
		}

		testTryResponse(t, ts.secret, ts.secret, allFixedCharsResp)
		testTryResponse(t, ts.secret, nonExistentWord, allDeadCharsResp)
	}
}

func testTryResponse(t *testing.T, secret, try, want string) {
	resp, err := getGameResponse(secret, try)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if resp != want {
		t.Fatalf("%q, %q: want %q, got %q", secret, try, want, resp)
	}
}
