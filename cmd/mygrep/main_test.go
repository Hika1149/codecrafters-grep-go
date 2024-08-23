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
		//{[]byte("dog"), "d", true},
		//{[]byte("dog"), "d.g", true},
		//{[]byte("dog"), "(cat|dog)", true},
		//{[]byte("a cat"), "a (cat|dog)", true},
		//{[]byte("caats"), `ca+ts`, true},
		//{[]byte("cat"), `\w+`, true},
		//{[]byte("cat and cat"), `(cat) and \1`, true},
		//{[]byte("cat and cat"), `(\w+) and \1`, true},
		//{[]byte("cat and dog"), `(\w+) and \1`, false},
		//{[]byte("abcd is"), `[abcd]+ is`, true},
		//{[]byte("abcd is abcd, not efg"), `([abcd]+) is \1, not [^xyz]+`, true},
		//{[]byte("this starts and ends with this"), `^(\w+) starts and ends with \1$`, true},
		// //alternation with wildcard
		//{[]byte("bugs here"), `(b..s|c..e) here`, true},
		//{[]byte("bugs here and bugs there"), `(b..s|c..e) here and \1 there`, true},
		//{[]byte("3 red squares and 3 red circles"), `(\d+) (\w+) squares and \1 \2 circles`, true},
		//{[]byte("3 red squares and 4 red circles"), `(\d+) (\w+) squares and \1 \2 circles`, false},
		{[]byte("cat and fish, cat with fish"), `(c.t|d.g) and (f..h|b..d), \1 with \2`, true},
	}

	for _, tt := range matchLineTests {
		matcher := internal.NewMatcher().ScanPattern(tt.pattern)

		m := matcher.Match(tt.line)
		if m != tt.expected {

			t.Errorf("line=%v pattern=%v Expected %v, but got %v\nmatcher:\n%v\n", string(tt.line), tt.pattern, tt.expected, m, matcher.String())

		}
	}

}
