package utils

import (
	"strconv"

	"github.com/cuducos/go-cnpj"
)

// FormatCNPJ -
func FormatCNPJ(cnpjNumbers string) string {
	return cnpj.Mask(cnpjNumbers)
}

func UnmaskCNPJ(CNPJ string) string {
	intCNPJ, err := strconv.Atoi(cnpj.Unmask(CNPJ))
	if err != nil {
		return ""
	}

	result := strconv.Itoa(intCNPJ)
	if len(result) < 3 {
		return ""
	}

	return result
}
