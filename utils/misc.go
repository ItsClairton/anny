package utils

import (
	"math"
	"math/rand"
)

func RandomBool() bool {
	return rand.Int63()&1 == 0
}

func ToDisplayTime(baseSeconds float64) string {
	hours := math.Floor(baseSeconds / 3600)
	minutes := math.Floor((baseSeconds - hours*3600) / 60)
	seconds := baseSeconds - hours*3600 - minutes*60

	if hours >= 1 {
		return Fmt("%02d:%02d:%02d", int(hours), int(minutes), int(seconds))
	}

	return Fmt("%02v:%02v", int(minutes), int(seconds))
}
