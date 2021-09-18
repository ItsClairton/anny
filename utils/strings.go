package utils

import "fmt"

func Fmt(s string, a ...interface{}) string {
	return fmt.Sprintf(s, a...)
}
