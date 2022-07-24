package sequence

import (
	"errors"
	"math/rand"
)

type Generator interface {
	// Generate - creates a random string of lettersNumber length.
	Generate(lettersNumber int) (string, error)
}

// if Sequence struct will no longer complains with Storage interface, code will be broken on building stage
var _ Generator = (*Sequence)(nil)

var letters []rune

type Sequence struct{}

func init() {
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func NewSequence() *Sequence {
	return &Sequence{}
}

// Generate will create string which contains of random letters with length of lettersNumber provided
func (s *Sequence) Generate(lettersNumber int) (string, error) {
	if lettersNumber < 0 {
		return "", errors.New("to generate random sequence positive number of letters must be provided")
	}

	b := make([]rune, lettersNumber)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b), nil
}
