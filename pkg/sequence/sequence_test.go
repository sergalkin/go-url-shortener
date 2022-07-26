package sequence_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sergalkin/go-url-shortener.git/pkg/sequence"
)

func TestNewSequence(t *testing.T) {
	tests := []struct {
		name string
		want *sequence.Sequence
	}{
		{
			name: "Sequence object can be created",
			want: &sequence.Sequence{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, sequence.NewSequence())
		})
	}
}

func TestGenerate(t *testing.T) {
	type args struct {
		lettersNumber int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Non empty string can be generated from positive value",
			args: args{8},
			want: 8,
		},
		{
			name: "Empty string will be generated from zero value",
			args: args{0},
			want: 0,
		},
		{
			name: "Empty string will be generated from negative zero value",
			args: args{-0},
			want: 0,
		},
		{
			name: "Error will be thrown on providing negative number to sequence generator",
			args: args{-1},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sequence.Sequence{}
			v, err := s.Generate(tt.args.lettersNumber)
			if err != nil {
				require.Errorf(t, err, "to generate random sequence positive number of letters must be provided")
			}
			require.Len(t, v, tt.want)
		})
	}
}

func ExampleSequence_Generate() {
	// Creating new Sequence
	s := sequence.NewSequence()

	// Generating random sequence of letters based on provided int
	s.Generate(5)
}
