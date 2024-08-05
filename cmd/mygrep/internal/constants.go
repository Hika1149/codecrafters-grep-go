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
)

const (
	Digits            = "0123456789"
	AlphanumericChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)
