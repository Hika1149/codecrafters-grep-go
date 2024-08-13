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
		{[]byte("dog"), "(cat|dog)", true},
		{[]byte("a cat"), "a (cat|dog)", true},
	}

	for _, tt := range matchLineTests {
		matcher := internal.NewMatcher().ScanPattern(tt.pattern)

		m := matcher.Match(tt.line)
		if m != tt.expected {

			t.Errorf("line=%v pattern=%v Expected %v, but got %v\nmatcher=%v\n", string(tt.line), tt.pattern, tt.expected, m, matcher.String())

		}
	}

}
