package internal

import (
	"bytes"
	"fmt"
)

type Ch struct {
	CharType CharType
	Value    string
}

type Matcher struct {
	Chs []*Ch
}

func NewMatcher() *Matcher {
	return &Matcher{}
}

// ScanPattern scans the reg pattern string and convert it to a slice of Ch
func (m *Matcher) ScanPattern(pattern string) *Matcher {

	chs := make([]*Ch, 0)
	for i := 0; i < len(pattern); {
		var (
			c  = pattern[i]
			nc byte
		)
		if i+1 < len(pattern) {
			nc = pattern[i+1]
		}
		if c == '\\' && nc != '\\' {
			chs = append(chs, &Ch{
				CharType: CharClassEscape,
				Value:    fmt.Sprintf("%c%c", c, nc),
			})
			i += 2
			continue
		}

		chs = append(chs, &Ch{
			CharType: CharLiteral,
			Value:    string(c),
		})
		i++

	}
	m.Chs = chs
	return m

}

func (m *Matcher) Match(text []byte) bool {
	// try match at each position text[i:] with pattern []chs
	for i := 0; i < len(text); i++ {
		if m.MatchHere(text[i:]) {
			return true
		}
	}
	return false
}

func (m *Matcher) MatchHere(text []byte) bool {

	i := 0

	for _, ch := range m.Chs {
		//
		if i >= len(text) {
			return false
		}

		tc := text[i]

		switch ch.CharType {

		case CharLiteral:
			if string(tc) != ch.Value {
				return false
			}

		case CharClassEscape:
			switch ch.Value {
			case "\\w":
				if !bytes.ContainsAny([]byte{tc}, AlphanumericChars) {
					return false
				}
			case "\\d":
				if !bytes.ContainsAny([]byte{tc}, Digits) {
					return false
				}
			}
		}
		// advance i
		i++

	}
	return true
}
