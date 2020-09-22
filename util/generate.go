package util

import "math/rand"

// RandString ...
func RandString() string {
	const l = 10
	const dict = "abcdefghijklmnopqrstuvwxyz"
	var res string
	for i := 0; i < l; i++ {
		r := dict[rand.Intn(len(dict))]
		res += string(r)
	}
	return res
}

// RandInt ...
func RandInt() int32 {
	return int32(rand.Intn(100))
}
