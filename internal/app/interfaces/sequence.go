package interfaces

type Sequence interface {
	Generate(lettersNumber int) (string, error)
}
