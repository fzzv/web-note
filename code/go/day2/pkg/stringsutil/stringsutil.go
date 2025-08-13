package stringsutil

import (
	"strings"
	"unicode"
)

func Upper(s string) string { return strings.ToUpper(s) }

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func IsPalindrome(s string) bool {
	buf := make([]rune, 0, len(s))
	for _, r := range []rune(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buf = append(buf, unicode.ToLower(r))
		}
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		if buf[i] != buf[j] {
			return false
		}
	}
	return true
} 
