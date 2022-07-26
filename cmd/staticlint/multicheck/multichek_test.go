package multicheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithBuiltin(t *testing.T) {
	tests := []struct {
		want     LintOptions
		name     string
		isActive bool
	}{
		{
			name:     "WithBultin will add basic analyzers to list if flag is provided",
			isActive: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := len(*NewWithOptions(WithBuiltin(tt.isActive))) > 0
			assert.True(t, result)
		})
	}
}

func TestWithStatic(t *testing.T) {
	tests := []struct {
		want     LintOptions
		name     string
		isActive bool
	}{
		{
			name:     "WithStatic will add basic analyzers to list if flag is provided",
			isActive: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := len(*NewWithOptions(WithStatic(tt.isActive))) > 0
			assert.True(t, result)
		})
	}
}

func TestWithExtra(t *testing.T) {
	tests := []struct {
		want     LintOptions
		name     string
		isActive bool
	}{
		{
			name:     "WithExtra will add basic analyzers to list if flag is provided",
			isActive: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := len(*NewWithOptions(WithExtra(tt.isActive))) > 0
			assert.True(t, result)
		})
	}
}
