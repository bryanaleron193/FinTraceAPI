package testutils

import (
	"math/rand"
)

// GenerateRandomString generates a random string of a given length
func GenerateRandomString(n int) string {
	lettersAndNumbers := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = lettersAndNumbers[rand.Intn(len(lettersAndNumbers))]
	}
	return string(b)
}
