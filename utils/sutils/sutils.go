package sutils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

var converter = md.NewConverter("", true, nil)

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

	if hours >= 1 {
		return Fmt("%02v:%02v:%02v", hours, minutes, seconds)
	}

	return Fmt("%02v:%02v", minutes, seconds)
}

func SplitString(r rune) bool {
	return r == ' ' || r == '\n'
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

func ToHexNumber(hex string) int {
	hex = strings.TrimPrefix(hex, "#")
	result, _ := strconv.ParseUint(hex, 16, 64)

	return int(result)
}

func ToMD(html string) string {
	md, _ := converter.ConvertString(html)

	return md
}
