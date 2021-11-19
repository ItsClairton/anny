package utils

import "fmt"

func Contains(array []string, term string) bool {
	for _, item := range array {
		if item == term {
			return true
		}
	}

	return false
}

func Fmt(s string, a ...interface{}) string {
	return fmt.Sprintf(s, a...)
}

func Is(value bool, a, n string) string {
	if value {
		return a
	} else {
		return n
	}
}
