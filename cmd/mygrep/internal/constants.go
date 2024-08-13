package internal

type CharType string

const (
	CharLiteral CharType = "charLiteral"
	// CharClassEscape eg: \d, \w, \s
	CharClassEscape CharType = "charClassEscape"

	// CharPositiveGroup eg: [xyz] [a-z]
	CharPositiveGroup CharType = "charPositiveGroup"
	// CharNegativeGroup eg: [^a-z]
	CharNegativeGroup CharType = "charNegativeGroup"

	// CharStartAnchor should begin the match at the beginning of the input string
	CharStartAnchor CharType = "^"

	// CharEndAnchor preceding pattern must occur at the end of the input string
	CharEndAnchor CharType = "$"

	// CharQuantifierOneOrMore +, eg: a+
	CharQuantifierOneOrMore CharType = "char.quantifier.oneOrMore"

	// CharQuantifierZeroOrOne ?, eg: a?
	CharQuantifierZeroOrOne CharType = "char.quantifier.zeroOrOne"

	// CharWildcard ., eg: a.
	CharWildcard CharType = "char.wildcard"

	// CharAlternation |, eg: a|b

	CharAlternation CharType = "char.alternation"

	// CharCaptureGroup eg: (a) can be used by backreference
	CharCaptureGroup CharType = "char.captureGroup"

	// CharBackReference eg: \1, \2, \3
	CharBackReference CharType = "char.backReference"
)

const (
	Digits            = "0123456789"
	AlphanumericChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)
