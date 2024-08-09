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

	matcher := internal.NewMatcher().ScanPattern(pattern)
	fmt.Println(matcher.String())
	ok := matcher.Match(line)

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}
