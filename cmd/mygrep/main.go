package main

import (
	"bytes"
	// Uncomment this to pass the first stage
	// "bytes"
	"fmt"
	"io"
	"os"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

const (
	Digits            = "0123456789"
	AlphanumericChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func matchLine(line []byte, pattern string) (bool, error) {
	//if utf8.RuneCountInString(pattern) != 1 {
	//	return false, fmt.Errorf("unsupported pattern: %q", pattern)
	//}

	var ok bool

	// positive character groups '[abc]', [a-z], [0-9], [a-zA-Z0-9_] ...
	if len(pattern) > 2 && pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
		for i := 1; i < len(pattern)-1; i++ {
			if ok, _ := matchLine(line, string(pattern[i])); ok {
				return ok, nil
			}
		}
	}

	switch pattern {
	case "\\w":
		ok = bytes.ContainsAny(line, AlphanumericChars)
	case "\\d":
		ok = bytes.ContainsAny(line, Digits)
	default:
		ok = bytes.ContainsAny(line, pattern)
	}

	return ok, nil
}
