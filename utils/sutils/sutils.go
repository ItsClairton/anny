package sutils

import (
	"fmt"
	"math"
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

func ToLower(s interface{}) string {
	return strings.ToLower(Fmt("%v", s))
}
