package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	if len(os.Args) != 2 {
		return
	}

	bytes := make([]byte, 2048)
	ptr := 0
	source := os.Args[1]

	for i := 0; i < len(source); i++ {
		switch source[i] {
		case '>':
			ptr++
			if ptr >= len(source) {
				ptr = 0
			}
		case '<':
			ptr--
			if ptr < 0 {
				ptr = 2047
			}
		case '+':
			bytes[ptr]++
		case '-':
			bytes[ptr]--
		case '.':
			z01.PrintRune(rune(bytes[ptr]))
		case '[':
			if bytes[ptr] == 0 {
				nest := 1
				for nest > 0 {
					i++
					if i >= len(source) {
						return
					}
					if source[i] == '[' {
						nest++
					} else if source[i] == ']' {
						nest--
					}
				}
			}
		case ']':
			if bytes[ptr] != 0 {
				nest := 1
				for nest > 0 {
					i--
					if i < 0 {
						return
					}
					if source[i] == ']' {
						nest++
					} else if source[i] == '[' {
						nest--
					}
				}
			}
		}
	}

}
