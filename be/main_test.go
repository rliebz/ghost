package be

import (
	"os"
	"testing"
)

// Avoid dealing with ANSI escape sequences.
func TestMain(m *testing.M) {
	if err := os.Setenv("NO_COLOR", "1"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
