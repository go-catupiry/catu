package helpers

import "strings"

func CleanupSymbol(symbol string) string {
	if strings.Contains(symbol, ":") {
		splitedS := strings.Split(symbol, ":")
		if len(splitedS) == 0 {
			// invalid symbol value
			return ""
		} else if len(splitedS) == 1 {
			return splitedS[0]
		} else {
			return splitedS[1]
		}
	}

	return symbol
}
