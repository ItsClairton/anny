package date

import (
	"strings"

	"github.com/ItsClairton/Anny/utils/i18n"
	"github.com/ItsClairton/Anny/utils/sutils"
)

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

func ToPrettyMonth(lc *i18n.Locale, month int) string {
	return lc.GetString(sutils.Fmt("months.[%d]", month-1))[0:3]
}

func ToPrettyDate(lc *i18n.Locale, date *Date) string {

	if date.Year == 0 {
		return lc.GetString("notYetReleased")
	}

	if date.Month == 0 && date.Day == 0 {
		return sutils.Fmt("%d", date.Year)
	}

	if date.Day == 0 {
		return strings.TrimSpace(lc.GetString("prettyDate", lc, "", ToPrettyMonth(lc, date.Month), date.Year))
	}

	return lc.GetString("prettyDate", date.Day, ToPrettyMonth(lc, date.Month), date.Year)
}
