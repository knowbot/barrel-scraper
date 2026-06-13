package utils

import (
	"strings"
	"unicode"
)

func CleanText(s string) string {
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
