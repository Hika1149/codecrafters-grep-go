package internal

import (
	"bytes"
	"fmt"
	"strings"
)

func (m *Matcher) peek(pattern string, i int) byte {
	if i+1 < len(pattern) {
		return pattern[i+1]
	}
	return 0

}

// ScanPattern scans the reg pattern string and convert it to a slice of Ch
func (m *Matcher) scanRawPattern(pattern string) []*Ch {

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
			nc = m.peek(pattern, i)
		)

		// handle end of string line anchor
		if c == '$' && i == len(pattern)-1 {
			chs = append(chs, &Ch{
				CharType: CharEndAnchor,
				Value:    "",
			})
			break
		}

		// handle quantifier one or more
		if c == '+' {
			poppedChs, lastElement := popCh(chs)
			chs = append(poppedChs, &Ch{
				CharType:         CharQuantifierOneOrMore,
				Value:            "",
				PrecedingElement: lastElement,
			})
			i++
			continue
		}

		// handle char class escape / backreference
		if c == '\\' && nc != '\\' {
			if bytes.ContainsAny([]byte{nc}, Digits) {
				chs = append(chs, &Ch{
					CharType: CharBackReference,
					Value:    fmt.Sprintf("%c", nc),
				})
			} else {
				chs = append(chs, &Ch{
					CharType: CharClassEscape,
					Value:    fmt.Sprintf("%c%c", c, nc),
				})
			}

			i += 2
			continue
		}

		// handle quantifier zero or one
		if nc == '?' {
			chs = append(chs, &Ch{
				CharType: CharQuantifierZeroOrOne,
				Value:    string(c),
			})
			i += 2
			continue
		}
		// handle wildcard
		if c == '.' {
			chs = append(chs, &Ch{
				CharType: CharWildcard,
				Value:    "",
			})
			i++
			continue
		}

		// try to found
		// - alternation
		// - capture group
		if c == '(' {

			//endPos := strings.Index(pattern[i:], ")")

			// nest capture group -> (a(b(c)d)e)
			// find the end of the capture group
			// ->  [s1,e1],[s2,e2],[sn,en]   s1<s2..<sn, e1>e2>...>en
			// just find e1, and recursively find e2, e3, en..
			endPos := -1
			hasNestedCaptureGroup := false
			for j := i + 1; j < len(pattern); {
				if pattern[j] == '(' {
					nextRight := strings.Index(pattern[j:], ")")
					j = j + nextRight + 1
					hasNestedCaptureGroup = true
					continue
				} else if pattern[j] == ')' {
					endPos = j
					// outer capture group
					break
				} else {
					j++
				}

			}

			if endPos != -1 {

				m.CaptureGroupCount = m.CaptureGroupCount + 1
				groupIndex := m.CaptureGroupCount

				// (a|b|c|d)
				// ((c.t|d.g) and (f..h|b..d))
				alterStrList := strings.Split(pattern[i+1:endPos], "|")
				// found alternation
				// each matched alter value can be capture group
				// if found nested capture group, can not just split into alternation group
				// for simple , each alter value don't have nested capture group, ((abc.af)|def) -> (abc|def)
				if len(alterStrList) > 1 && !hasNestedCaptureGroup {
					ch := &Ch{
						CharType:    CharAlternation,
						Value:       "",
						AlterValues: make([][]*Ch, 0),
						GroupIndex:  groupIndex,
					}
					for _, alterStr := range alterStrList {
						ch.AlterValues = append(ch.AlterValues, m.scanRawPattern(alterStr))
					}

					chs = append(chs, ch)

				} else {
					chs = append(chs, &Ch{
						CharType:      CharCaptureGroup,
						Value:         pattern[i+1 : endPos],
						AlterValues:   nil,
						GroupElements: m.scanRawPattern(pattern[i+1 : endPos]),
						GroupIndex:    groupIndex,
					})
				}

				// 2. store in CaptureGroups field for backreference
				m.CaptureGroups = append(m.CaptureGroups, "")
				//captureGroups = append(captureGroups, pattern[i+1:i+endPos])
				//i = i + endPos + 1
				i = endPos + 1

				continue

			}

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

	return chs

}

// ScanPattern scans the reg pattern string and convert it to a slice of Ch
func (m *Matcher) ScanPattern(pattern string) *Matcher {
	m.CaptureGroupCount = 0
	m.CaptureGroups = make([]string, 1) // 0 index is not used
	m.Chs = m.scanRawPattern(pattern)
	return m

}
