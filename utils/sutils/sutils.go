package sutils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

var converter = md.NewConverter("", true, nil)

func GetRegex(reg string) *regexp.Regexp {
	result, _ := regexp.Compile(reg)
	return result
}

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
		return Fmt("%02d:%02d:%02d", int(hours), int(minutes), int(seconds))
	}

	return Fmt("%02v:%02v", int(minutes), int(seconds))
}

func SplitString(r rune) bool {
	return r == ' ' || r == '\n'
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
