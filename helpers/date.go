package helpers

import (
	"regexp"
	"time"

	"github.com/go-catupiry/catu/configuration"
)

func FormatDate(date *time.Time, format string) string {
	timeZone := configuration.GetEnv("SITE_TIMEZONE", "")
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		panic(err)
	}

	return date.In(loc).Format(format)
}

func ExtractYearFromText(text string) string {
	re := regexp.MustCompile(`\d{4}`)
	match := re.FindStringSubmatch(text)
	if len(match) > 0 {
		return match[0]
	}

	return ""
}
