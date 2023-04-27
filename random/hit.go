package random

import "math/rand"

func Hit(in, base uint32) bool {
	if base == 0 {
		return true
	}
	return in >= base || in >= UintU(base)
}

func HitMillion(in uint32) bool {
	return in >= 10000 || in >= UintU(10000)
}

func Hit64(in, base int64) bool {
	if base == 0 {
		return true
	}
	if base < 0 {
		return false
	}
	return in >= base || in >= rand.Int63n(base)
}
