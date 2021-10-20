package utils

import (
	"math/rand"
	"regexp"
	"time"
)

var URLRegex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

func RandomBool() bool {
	return rand.Int63()&1 == 0
}

func FormatTime(duration time.Duration) string {
	totalSeconds := int64(duration.Seconds())
	days := totalSeconds / 86400
	hours := totalSeconds % 86400 / 3600
	minutes := totalSeconds % 3600 / 60
	seconds := totalSeconds % 60

	if days > 0 {
		return Fmt("%02d:%02d:%02d:%02d", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return Fmt("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return Fmt("%02d:%02d", minutes, seconds)
}
