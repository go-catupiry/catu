package helpers

import "strings"

func ParseUrlQueryOrder(order string) (string, bool, bool) {
	if order == "" {
		return "", true, false
	}

	words := strings.Fields(order)
	if len(words) != 2 {
		return "", true, false
	}

	operatorStr := strings.ToUpper(words[1])
	if operatorStr == "ASC" {
		return words[0], true, true
	} else {
		return words[0], false, true
	}
}
