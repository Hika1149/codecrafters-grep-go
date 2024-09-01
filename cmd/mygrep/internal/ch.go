package internal

import (
	"fmt"
	"strings"
)

type Ch struct {
	CharType CharType
	Value    string

	// AlterValues is used for alternation
	// eg: (a|b..c|d...e|xyz)
	// each alter value can be a slice of ch
	// each matched alter value can be capture group

	AlterValues [][]*Ch

	// PrecedingElement is used by quantifier
	PrecedingElement *Ch

	// GroupElements elements of capture group
	GroupElements []*Ch
	// GroupIndex index of capture group in mather pattern
	GroupIndex int
}

func (ch *Ch) String() string {

	//var precedingStr string
	//if ch.PrecedingElement != nil {
	//	precedingStr = ch.PrecedingElement.String()
	//}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("charType: %s value: %s ", ch.CharType, ch.Value))
	if ch.CharType == CharCaptureGroup {
		buf.WriteString(fmt.Sprintf("groupIndex: %d ", ch.GroupIndex))
	}
	buf.WriteString("\n")
	if len(ch.AlterValues) > 0 {
		buf.WriteString(" alterValues:\n")
		for _, alterValue := range ch.AlterValues {
			buf.WriteString("[\n")
			for _, ch := range alterValue {
				buf.WriteString("\t")
				buf.WriteString(ch.String())
			}
			buf.WriteString("]\n")
		}
	}

	if len(ch.GroupElements) > 0 {
		buf.WriteString(" groupElements:\n")
		for _, groupElement := range ch.GroupElements {
			buf.WriteString("\t")
			buf.WriteString(groupElement.String())
		}
	}

	return buf.String()
}

func popCh(chs []*Ch) ([]*Ch, *Ch) {
	length := len(chs)
	if length == 0 {
		return chs, nil
	}

	return chs[:length-1], chs[length-1]

}
