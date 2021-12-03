package utils

import (
	"math"
	"strings"

	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

var ac = accounting.Accounting{Symbol: "", Precision: 2, Thousand: ".", Decimal: ","}

// DecimalToPrice - Format a decimal number to string in market format Ex: 2000.55 to 2.000,55
func DecimalToPrice(v decimal.Decimal) string {
	return "R$ " + ac.FormatMoneyDecimal(v)
}

// DecimalToPriceNoSign -
func DecimalToPriceNoSign(v decimal.Decimal) string {
	return ac.FormatMoneyDecimal(v)
}

// DecimalToPercent -
func DecimalToPercent(v decimal.Decimal) string {
	v2 := strings.ReplaceAll(v.String(), ".", ",")
	return v2 + "%"
}

// Convert string to decimal without errors
func NewFromString(v string) decimal.Decimal {
	r, _ := decimal.NewFromString(v)

	return r
}

func RoundCurrency(n float64) float64 {
	return math.Round(n*100) / 100
}
