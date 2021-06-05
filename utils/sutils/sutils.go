package sutils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Fmt(s string, a ...interface{}) string {
	return fmt.Sprintf(s, a...)
}

func Is(cond bool, afirmative string, negative string) string {
	if cond {
		return afirmative
	} else {
		return negative
	}
}

func ToHHMMSS(baseSeconds float64) string {
	hours := math.Floor(baseSeconds / 3600)
	minutes := math.Floor((baseSeconds - hours*3600) / 60)
	seconds := baseSeconds - hours*3600 - minutes*60
	return Fmt("%s%d:%d",
		Is(hours < 1, "", Fmt("%d:", int64(hours))),
		int64(minutes),
		int64(seconds))
}

func ToPrettyMonth(m int) string {
	switch m {
	case 1:
		return "Jan"
	case 2:
		return "Fev"
	case 3:
		return "Mar"
	case 4:
		return "Abr"
	case 5:
		return "Mai"
	case 6:
		return "Jun"
	case 7:
		return "Jul"
	case 8:
		return "Ago"
	case 9:
		return "Set"
	case 10:
		return "Out"
	case 11:
		return "Nov"
	case 12:
		return "Dez"
	default:
		return ""
	}
}

func ToLower(s interface{}) string {
	return strings.ToLower(Fmt("%v", s))
}

func ToHexNumber(hex string) (int, error) {

	hex = strings.TrimPrefix(hex, "#")

	result, err := strconv.ParseUint(hex, 16, 64)

	if err != nil {
		return -1, err
	}

	return int(result), nil
}
