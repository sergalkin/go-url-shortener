package certificate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Certificate can be generated in provided path",
		},
	}
	for _, tt := range tests {
		tmpDir := os.TempDir()
		t.Run(tt.name, func(t *testing.T) {
			Generate(tmpDir)
			assert.DirExists(t, tmpDir)
			assert.FileExists(t, tmpDir+"/cert.key")
		})
	}
}
