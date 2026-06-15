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

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool, len(slice))
	unique := make([]T, 0, len(slice))
	for _, item := range slice {
		if !seen[item] {
			unique = append(unique, item)
			seen[item] = true
		}
	}
	return unique
}
