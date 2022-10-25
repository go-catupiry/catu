package helpers

import (
	"strings"
)

func ParseUrlQueryOrder(order, sort, sortDirection string) (string, bool, bool) {
	if sort != "" {
		return parseUrlQuerySort(sort, sortDirection)
	}

	if order == "" {
		return "", true, false
	}

	words := strings.Fields(order)
	if len(words) != 2 {
		return "", true, false
	}

	operatorStr := strings.ToUpper(words[1])
	if operatorStr == "ASC" {
		return words[0], false, true
	} else {
		return words[0], true, true
	}
}

func parseUrlQuerySort(sort, sortDirection string) (string, bool, bool) {
	if sortDirection == "" || (sortDirection != "DESC" && sortDirection != "ASC") {
		sortDirection = "DESC"
	}

	sort = strings.Replace(strings.TrimSpace(sort), `"`, "", -1)

	if sortDirection == "ASC" {
		return sort, false, true
	} else {
		return sort, true, true
	}
}
