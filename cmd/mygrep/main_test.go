package main

import (
	"github.com/codecrafters-io/grep-starter-go/cmd/mygrep/internal"
	"testing"
)

func TestMatcher(t *testing.T) {
	type MatchTest struct {
		line     []byte
		pattern  string
		expected bool
	}

	var matchLineTests = []MatchTest{

		{[]byte("dog"), "d", true},
		{[]byte("dog"), "d.g", true},
	}

	for _, tt := range matchLineTests {
		m := internal.NewMatcher().ScanPattern(tt.pattern).Match(tt.line)

		if m != tt.expected {
			t.Errorf("line=%v pattern=%v Expected %v, but got %v\n", string(tt.line), tt.pattern, tt.expected, m)
		}
	}

}
