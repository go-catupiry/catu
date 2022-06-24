package catu

import (
	"html/template"

	"github.com/go-catupiry/catu/helpers"
	"github.com/go-catupiry/catu/pagination"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func noEscapeHTML(str string) template.HTML {
	return template.HTML(str)
}

func paginate(ctx *RequestContext, pager *pagination.Pager, queryString string) template.HTML {
	return renderPager(ctx, pager, queryString)
}

type ContentDates interface {
	GetTeaserDatesHTML(separator string) template.HTML
}

func contentDates(record ContentDates, separator string) template.HTML {
	return record.GetTeaserDatesHTML(separator)
}

func truncate(text string, length int, ellipsis string) template.HTML {
	html, err := helpers.Truncate(text, length, ellipsis)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"text":     text,
			"length":   length,
			"ellipsis": ellipsis,
		}).Error("truncate error on truncate text")
	}
	return html
}

func formatCurrency(value decimal.Decimal) string {
	return helpers.DecimalToPrice(value)
}

func formatDecimalWithDots(value decimal.Decimal) string {
	return helpers.FormatDecimalWithDots(value)
}
