package utils

import (
	"errors"
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
	"math/rand"
)

// if Sequence struct will no longer complains with Storage interface, code will be broken on building stage
var _ interfaces.Sequence = (*Sequence)(nil)

type Sequence struct{}

func NewSequence() *Sequence {
	return &Sequence{}
}

// Generate will create string which contains of random letters with length of lettersNumber provided
func (s *Sequence) Generate(lettersNumber int) (string, error) {
	if lettersNumber < 0 {
		return "", errors.New("to generate random sequence positive number of letters must be provided")
	}

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, lettersNumber)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b), nil
}
