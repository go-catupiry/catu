package helpers

import (
	"math/rand"

	"github.com/microcosm-cc/bluemonday"
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

var stripTagsPolicy = bluemonday.StripTagsPolicy()

// TruncateString - Truncate one string with with x words
func TruncateString(str string, length int, omission string) string {
	if length <= 0 {
		return ""
	}

	orgLen := len(str)
	if orgLen <= length {
		return str
	}

	if orgLen > length {
		return str[:length] + omission
	}

	return str[:length]

	// // Support Japanese
	// // Ref: Range loops https://blog.golang.org/strings
	// truncated := ""
	// count := 0
	// for _, char := range str {
	// 	truncated += string(char)
	// 	count++
	// 	if count >= length {
	// 		break
	// 	}
	// }

	// return truncated
}

// StripTags - Remove tags from html text
func StripTags(str string) string {
	return stripTagsPolicy.Sanitize(str)
}

// StripTagsAndTruncate - Remove all tags from string and then truncate if need with omission
func StripTagsAndTruncate(str string, length int, omission string) string {
	return TruncateString(StripTags(str), length, omission)
}
