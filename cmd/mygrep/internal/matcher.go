package internal

import (
	"bytes"
	"fmt"
	"strings"
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
func (m *Matcher) String() string {
	str := ""
	for _, ch := range m.Chs {
		str += fmt.Sprintf("charType: %s value: %s\n", ch.CharType, ch.Value)
	}
	return str

}

// ScanPattern scans the reg pattern string and convert it to a slice of Ch
func (m *Matcher) ScanPattern(pattern string) *Matcher {

	chs := make([]*Ch, 0)

	var (
		i = 0
	)

	// detect start of string line anchor
	if strings.HasPrefix(pattern, "^") {
		chs = append(chs, &Ch{
			CharType: CharStartAnchor,
			Value:    "",
		})
		i++
	}

	for i < len(pattern) {
		var (
			c  = pattern[i]
			nc byte
		)
		if i+1 < len(pattern) {
			nc = pattern[i+1]
		}

		// handle end of string line anchor
		if c == '$' && i == len(pattern)-1 {
			chs = append(chs, &Ch{
				CharType: CharEndAnchor,
				Value:    "",
			})
			break
		}

		// handle char class escape
		if c == '\\' && nc != '\\' {
			chs = append(chs, &Ch{
				CharType: CharClassEscape,
				Value:    fmt.Sprintf("%c%c", c, nc),
			})
			i += 2
			continue
		}
		// handle char positive/negative group
		if c == '[' {
			endPos := strings.Index(pattern[i:], "]")
			// found group
			if endPos != -1 {
				charGroup := pattern[i+1 : i+endPos]
				charType := CharPositiveGroup
				if charGroup[0] == '^' {
					charType = CharNegativeGroup
					charGroup = charGroup[1:]
				}
				chs = append(chs, &Ch{
					CharType: charType,
					Value:    charGroup,
				})
				//advanced
				i = i + endPos + 1
				continue
			}
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

	// should match from beginning of text
	if m.Chs[0].CharType == CharStartAnchor {
		return m.MatchHere(text, m.Chs[1:])
	}

	// try match at each position text[i:] with pattern []chs
	for i := 0; i < len(text); i++ {
		if m.MatchHere(text[i:], m.Chs) {
			return true
		}
	}
	return false
}

func (m *Matcher) MatchHere(text []byte, Chs []*Ch) bool {

	i := 0

	for _, ch := range Chs {

		// handle input text reaches end
		if i >= len(text) {

			// -> pattern also reaches end
			if ch.CharType == CharEndAnchor {
				return true
			}
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
		case CharPositiveGroup:
			// simple
			// - no range separator
			// - no class escape
			if !bytes.ContainsAny([]byte{tc}, ch.Value) {
				return false
			}

		case CharNegativeGroup:
			if bytes.ContainsAny([]byte{tc}, ch.Value) {
				return false
			}

		case CharEndAnchor:
			// If the regular expression is a $ at the end of the expression,
			// then the text matches only if it too is at its end.

			// previous matched xxx is not at the end of text

			return false

		}

		// advance i
		i++

	}
	return true
}
