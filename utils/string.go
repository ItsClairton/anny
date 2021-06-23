package utils

import (
	"fmt"
	"strconv"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

var converter = md.NewConverter("", true, nil)

func Fmt(s string, a ...interface{}) string {
	return fmt.Sprintf(s, a...)
}

func Is(cond bool, afirmative, negative string) string {
	if cond {
		return afirmative
	} else {
		return negative
	}
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
