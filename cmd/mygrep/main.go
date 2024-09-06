package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
)

var matches []rune

var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}
	pattern := os.Args[2]
	line, _ := io.ReadAll(os.Stdin) // assume we're only dealing with a single line

	if ok := matchLine(line, pattern); !ok {
		fmt.Println("no match found")
		os.Exit(1)
	}
	fmt.Printf("match(es) found: %q", matches)
}

func matchLine(line []byte, pattern string) bool {
	var found bool
	for i := 0; i < len(pattern); i++ {
		pt := rune(pattern[i])
		fmt.Println(string(pt))
		if pt == '\\' && i+1 < len(pattern) {
			if pattern[i+1] == 'd' {
				i++
				fmt.Println("here", i)
				found = matchFunc(line, 0, unicode.IsDigit)
			} else if pattern[i+1] == 'w' {
				i++
				found = matchFunc(line, 0, func(r rune) bool {
					return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_'
				})
			}
		}
		if matchFunc(line, pt, nil) {
			found = true
		} else {
			continue
		}
	}
	return found
}

// type findMatch func(rune)

func matchFunc(line []byte, pt rune, f func(rune) bool) bool {
	var found bool
	for _, b := range line {
		li := rune(b)
		if f != nil && f(li) { // check if f is not nil and evaluates to true
			matches = append(matches, li)
			found = true
			continue
		}
		if pt == li {
			matches = append(matches, li)
			found = true
		}
		continue
	}
	return found
}
