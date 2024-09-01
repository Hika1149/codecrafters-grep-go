package internal

import (
	"bytes"
	"fmt"
	"strconv"
)

type MatchedResult struct {
	Matched bool
	// all possible matched endPos
	EndPosList []int
}

type Matcher struct {
	// Chs: split pattern string to slice of Ch
	Chs []*Ch

	// CaptureGroups storing capturing groups matched value and will be used by backReferences
	CaptureGroups []string

	CaptureGroupCount int
}

func NewMatcher() *Matcher {
	return &Matcher{}
}
func (m *Matcher) String() string {
	str := ""
	for _, ch := range m.Chs {
		str += fmt.Sprintf("%v", ch.String())
	}
	return str

}

func (m *Matcher) Match(text []byte) bool {

	// should match from beginning of text
	if m.Chs[0].CharType == CharStartAnchor {
		r := m.MatchHere(text, m.Chs[1:])
		return r.Matched
	}

	// try match at each position text[i:] with pattern []chs
	for i := 0; i < len(text); i++ {
		if r := m.MatchHere(text[i:], m.Chs); r.Matched {
			return true
		}
	}
	return false
}
func (m *Matcher) MatchBasePattern(tc byte, ch *Ch) bool {
	switch ch.CharType {
	case CharLiteral:
		if string(tc) == ch.Value {
			return true
		}
	case CharClassEscape:
		switch ch.Value {
		case "\\w":
			if bytes.ContainsAny([]byte{tc}, AlphanumericChars) {
				return true
			}
		case "\\d":
			if bytes.ContainsAny([]byte{tc}, Digits) {
				return true
			}
		}
	case CharPositiveGroup:
		// simple
		// - no range separator
		// - no class escape
		if bytes.ContainsAny([]byte{tc}, ch.Value) {
			return true
		}

	case CharNegativeGroup:
		if !bytes.ContainsAny([]byte{tc}, ch.Value) {
			return true
		}

	}
	return false

}

func (m *Matcher) MatchHere(text []byte, Chs []*Ch) *MatchedResult {

	i := 0
	res := &MatchedResult{
		Matched:    true,
		EndPosList: make([]int, 0),
	}

	for pi, ch := range Chs {

		// handle input text reaches end
		if i >= len(text) {

			// -> pattern also reaches end
			if ch.CharType == CharEndAnchor {
				res.EndPosList = append(res.EndPosList, i)
				break
			}
			res.Matched = false
			break
		}

		var (
			// previous char
			//pc byte
			tc = text[i]
		)
		//if i > 0 {
		//	pc = text[i-1]
		//}

		switch ch.CharType {

		case CharLiteral, CharClassEscape,
			CharPositiveGroup,
			CharNegativeGroup:
			if !m.MatchBasePattern(tc, ch) {
				res.Matched = false
				break
			}

		case CharEndAnchor:
			// If the regular expression is a $ at the end of the expression,
			// then the text matches only if it too is at its end.

			// previous matched xxx is not at the end of text
			res.Matched = false
			break

		case CharQuantifierOneOrMore:
			// should match precedingElement one or more times
			// recursive try one or more times
			sr := false
			for j := i; j < len(text) && m.MatchBasePattern(text[j], ch.PrecedingElement); j++ {

				if mr := m.MatchHere(text[j+1:], Chs[pi+1:]); mr.Matched {
					// store all possible matched endPos
					for _, endPos := range mr.EndPosList {
						res.EndPosList = append(res.EndPosList, j+1+endPos)
					}
					sr = true
				}
			}
			if !sr {
				res.Matched = false
			}
			return res

		case CharQuantifierZeroOrOne:
			// should match ch.Value zero or one times

			// todo support wildcard with quantifier

			// zero times
			if mr := m.MatchHere(text[i:], Chs[pi+1:]); mr.Matched {
				res = mr
				return res
			}
			// one times
			if string(tc) == ch.Value {
				if mr := m.MatchHere(text[i+1:], Chs[pi+1:]); mr.Matched {
					res = mr
					break
				}
				res.Matched = false
				return res
			}
			res.Matched = false
			return res

		case CharWildcard:
			i++
			continue
		case CharAlternation:
			// try each alternation
			mq := false
			for _, alterValue := range ch.AlterValues {

				//  try match current alterValue
				if ma := m.MatchHere(text[i:], alterValue); ma.Matched {
					// matched current alterValue

					for _, endPosAlter := range ma.EndPosList {
						nextI := i + endPosAlter
						// each matched alter value can be capture group
						m.CaptureGroups[ch.GroupIndex] = string(text[i:nextI])
						//fmt.Printf("matched alterValue len=%v nextI=%v\n", len(alterValue), nextI)
						if mr := m.MatchHere(text[nextI:], Chs[pi+1:]); mr.Matched {
							// store all possible matched endPos
							for _, endPos := range mr.EndPosList {
								res.EndPosList = append(res.EndPosList, nextI+endPos)
							}
							mq = true
						}
					}

				}
			}
			if !mq {
				res.Matched = false
			}
			//fmt.Println("alteration end mq=", mq)
			return res

		case CharCaptureGroup:
			// if text[i:] and ch.Groups matched
			// - 1. need to know the all possible matched text endIndex, then we can advance i
			// - 2. need to know index of current match group, store current matched group value for backreference
			sr := false
			fmt.Printf("capture group groupIndex=%v lenGroupElements=%v\n", ch.GroupIndex, len(ch.GroupElements))
			if mg := m.MatchHere(text[i:], ch.GroupElements); mg.Matched {
				for _, mgEnd := range mg.EndPosList {
					// store matched group value
					m.CaptureGroups[ch.GroupIndex] = string(text[i : i+mgEnd])
					fmt.Println("matched group value", ch.GroupIndex, m.CaptureGroups[ch.GroupIndex])

					nextI := i + mgEnd
					// advanced to i+mgEnd
					if mr := m.MatchHere(text[nextI:], Chs[pi+1:]); mr.Matched {
						sr = true
						// store all possible matched endPos
						for _, endPos := range mr.EndPosList {
							res.EndPosList = append(res.EndPosList, nextI+endPos)
						}
					}
				}
			}
			fmt.Println("capture group end sr=", sr)
			res.Matched = sr
			return res
		case CharBackReference:
			groupIndex, _ := strconv.Atoi(ch.Value)
			groupValue := m.CaptureGroups[groupIndex]
			nextI := i + len(groupValue)
			//fmt.Printf("backReferenceValue=%v groupValue=%v i=%v nextI=%v lenT=%v\n", groupIndex, groupValue, i, nextI, len(text))
			if string(text[i:nextI]) == groupValue {
				i = nextI
				continue
			}
			res.Matched = false
			return res
		}

		// advance i
		i++

	}
	if res.Matched {
		res.EndPosList = append(res.EndPosList, i)
	}
	return res
}
