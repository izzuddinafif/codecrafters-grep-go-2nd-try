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
	var found, squareBrackets bool
	var in []rune
	for i := 0; i < len(pattern); i++ {
		pt := rune(pattern[i])

		if pt == '\\' && i+1 < len(pattern) {
			if pattern[i+1] == 'd' {
				i++
				found = matchFunc(line, 0, func(r ...rune) bool {
					return unicode.IsDigit(r[0])
				})
				continue
			} else if pattern[i+1] == 'w' {
				i++
				found = matchFunc(line, 0, func(r ...rune) bool {
					return unicode.IsDigit(r[0]) || unicode.IsLetter(r[0]) || r[0] == '_'
				})
				continue
			}
		} else if pt == '[' {
			squareBrackets = true
			continue
		}
		if squareBrackets {
			if pt == ']' {
				squareBrackets = false
				if len(in) > 0 {
					if in[0] == '^' {
						negatedChars := string(in[1:])
						for _, b := range line {
							li := rune(b)
							if !bytes.ContainsRune([]byte(negatedChars), li) {
								matches = append(matches, li)
								found = true
							}
						}
					} else {
						found = bytes.ContainsAny(line, string(in))
						matches = append(matches, findIndices(line, in)...)
					}
				}
				in = nil
				continue
			} else {
				in = append(in, pt)
			}
		} else {
			found = matchFunc(line, pt, nil)
		}
		continue
	}
	return found
}

func matchFunc(line []byte, pt rune, f func(...rune) bool) bool {
	var found bool
	for _, b := range line {
		li := rune(b)
		if f != nil {
			fmt.Printf("Checking: %c, with function\n", li)
			if f(li) { // check if f is not nil and evaluates to true
				fmt.Println("Function returned true")
				matches = append(matches, li)
				found = true
				continue
			}
		} else if pt == li {
			fmt.Println("Direct match")
			matches = append(matches, li)
			found = true
		}
	}
	return found
}

func findIndices(line []byte, chars []rune) []rune {
	var result []rune
	for _, b := range line {
		li := rune(b)
		for _, char := range chars {
			if li == char {
				result = append(result, li)
				break
			}
		}
	}
	return result
}

func contains(chars []rune, r rune) bool {
	for _, c := range chars {
		if r == c {
			return true
		}
	}
	return false
}
