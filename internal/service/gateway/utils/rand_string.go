package utils

import (
	"crypto/rand"
	"math/big"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	maxlen := big.NewInt(int64(len(letterRunes)))

	for i := range b {
		j, _ := rand.Int(rand.Reader, maxlen)
		b[i] = letterRunes[j.Uint64()]
	}

	return string(b)
}
