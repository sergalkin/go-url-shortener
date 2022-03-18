package utils

import "math/rand"

// Generate will create string which contains of random letters with length of lettersNumber provided
func Generate(lettersNumber int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, lettersNumber)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
