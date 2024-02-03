package main

import (
	"strings"
	"testing"
)

func TestGetGameResponse(t *testing.T) {
	for _, ts := range testSessions {
		for _, try := range ts.tries {
			testTryResponse(t, ts.secret, try.try, try.resp)
		}

		testTryResponse(t, ts.secret, ts.secret, allFixedCharsResp)
		testTryResponse(t, ts.secret, strings.Repeat(" ", wordLen),
			allDeadCharsResp)
	}
}

func testTryResponse(t *testing.T, secret, try, want string) {
	if resp := getGameResponse(secret, try); resp != want {
		t.Fatalf("%q, %q: want %q, got %q", secret, try, want, resp)
	}
}
