package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	alph := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result = append(result, alph[rnd.Intn(len(alph))])
	}

	return string(result)
}
