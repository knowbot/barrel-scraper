package main

import (
	"strings"
	"unicode"
)

func cleanText(s string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsPunct(r) {
				return -1
			}
			return r
		},
		s,
	)
}
