package utils

import (
	"math/rand"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

// SearchForString -
func SearchForString(str string, substr string) (int, int) {
	m := search.New(language.English, search.IgnoreCase)
	return m.IndexString(str, substr)
}

// StringIsInText -
func StringIsInText(s, substr string) (bool, int, int) {
	var start, end int

	start, end = SearchForString(s, substr)
	if start != -1 && end != -1 {
		return true, start, end
	}

	return false, start, end
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
