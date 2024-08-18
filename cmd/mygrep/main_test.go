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
		{[]byte("caats"), `ca+ts`, true},
		{[]byte("cat"), `\w+`, true},
		{[]byte("cat and cat"), `(cat) and \1`, true},
		{[]byte("cat and cat"), `(\w+) and \1`, true},
		{[]byte("cat and dog"), `(\w+) and \1`, false},
		{[]byte("abcd is abcd, not efg"), `"([abcd]+) is \1, not [^xyz]+`, false},
	}

	for _, tt := range matchLineTests {
		matcher := internal.NewMatcher().ScanPattern(tt.pattern)

		m := matcher.Match(tt.line)
		if m != tt.expected {

			t.Errorf("line=%v pattern=%v Expected %v, but got %v\nmatcher:\n%v\n", string(tt.line), tt.pattern, tt.expected, m, matcher.String())

		}
	}

}
