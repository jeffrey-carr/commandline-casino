package utils

import "unicode/utf8"

func PadRight(s string, n int) string {
	for RuneCount(s) < n {
		s += " "
	}
	// Truncate if too long (unlikely)
	return CutRunes(s, n)
}

func PadLeft(s string, n int) string {
	for RuneCount(s) < n {
		s = " " + s
	}
	return CutRunes(s, n)
}

func RuneCount(s string) int { return utf8.RuneCountInString(s) }

func CutRunes(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}
