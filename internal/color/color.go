// Package color exports utilities for applying color to terminal outputs.
package color

import (
	"os"
	"regexp"
	"runtime"
	"sync"
)

const (
	// ANSIReset is the ANSI escape sequence for reset.
	ANSIReset = "\033[0m"
	// ANSIRed is the ANSI escape sequence for red.
	ANSIRed = "\033[31m"
	// ANSIGreen is the ANSI escape sequence for green.
	ANSIGreen = "\033[32m"
	// ANSIYellow is the ANSI escape sequence for yellow.
	ANSIYellow = "\033[33m"
)

// Red colors the string red.
func Red(s string) string {
	return apply(ANSIRed, s)
}

// Green colors the string green.
func Green(s string) string {
	return apply(ANSIGreen, s)
}

// Yellow colors the string yellow.
func Yellow(s string) string {
	return apply(ANSIYellow, s)
}

var (
	// reReset captures any reset sequence that is followed by at least one
	// character. The extra character check prevents us from writing escape
	// squences at the end of lines.
	reReset = regexp.MustCompile(`(\033\[0m)(.)`)
	// reANSI identifies any non-reset ANSI escape sequence.
	reANSI = regexp.MustCompile(`\033\[[1-9][\d;]*m`)
)

func apply(color string, s string) string {
	if !Enabled() {
		return s
	}

	// Preserve existing colors by resetting before ANSI escape sequences
	// and re-applying after the reset sequence.
	s = reReset.ReplaceAllString(s, "$1"+color+"$2")
	s = reANSI.ReplaceAllString(s, ANSIReset+"$0")
	return color + s + ANSIReset
}

// Enabled returns whether colors are enabled.
var Enabled = sync.OnceValue(func() bool {
	if runtime.GOOS == "windows" {
		return false
	}

	_, ok := os.LookupEnv("NO_COLOR")
	return !ok
})
