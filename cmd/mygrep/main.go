package main

import (
	"github.com/codecrafters-io/grep-starter-go/cmd/mygrep/internal"

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

	// patten "\d apple"
	// First turn pattern into a slice of structs
	//

	ok := internal.NewMatcher().ScanPattern(pattern).Match(line)

	//ok, err := matchLine(line, pattern)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	//	os.Exit(2)
	//}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

//func matchLine(line []byte, pattern string) (bool, error) {
//	//if utf8.RuneCountInString(pattern) != 1 {
//	//	return false, fmt.Errorf("unsupported pattern: %q", pattern)
//	//}
//
//	var ok bool
//
//	// positive/negative character groups
//	if len(pattern) > 2 && pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
//
//		// negative character groups
//		// [^abc], [^a-z], [^0-9], [^a-zA-Z0-9_] ...
//		if pattern[1] == '^' {
//			for i := 2; i < len(pattern)-1; i++ {
//				if ok, _ := matchLine(line, string(pattern[i])); ok {
//					return false, nil
//				}
//			}
//			return true, nil
//		} else {
//			// positive character groups '[abc]', [a-z], [0-9], [a-zA-Z0-9_] ...
//			for i := 1; i < len(pattern)-1; i++ {
//				if ok, _ := matchLine(line, string(pattern[i])); ok {
//					return ok, nil
//				}
//			}
//		}
//
//	}
//
//	switch pattern {
//	case "\\w":
//		ok = bytes.ContainsAny(line, AlphanumericChars)
//	case "\\d":
//		ok = bytes.ContainsAny(line, Digits)
//	default:
//		ok = bytes.ContainsAny(line, pattern)
//	}
//
//	return ok, nil
//}
