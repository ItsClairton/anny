package utils

import "math/rand"

func RandomBool() bool {
	return rand.Int63()&1 == 0
}
