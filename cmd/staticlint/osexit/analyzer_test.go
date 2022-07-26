package osexit

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestRun validates that ExitAnalyzer can find os.Exit call
func TestRun(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ExitAnalyzer, "./...")
}
