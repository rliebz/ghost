package be

import (
	"fmt"
	"time"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/ghostlib"
)

// Eventually asserts that a function eventually returns an Ok [ghost.Result].
func Eventually(
	f func() ghost.Result,
	timeout time.Duration,
	interval time.Duration,
) ghost.Result {
	args := ghostlib.ArgsFromAST(f, timeout, interval)
	argF := args[0]

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	lastRun := ghost.Result{
		Ok:      false,
		Message: fmt.Sprintf("%s did not return value within %s timeout", argF, timeout),
	}

	ch := make(chan ghost.Result, 1)

	for tick := ticker.C; ; {
		select {
		case <-timer.C:
			return lastRun
		case <-tick:
			tick = nil
			go func() { ch <- f() }()
		case lastRun = <-ch:
			if lastRun.Ok {
				return lastRun
			}
			tick = ticker.C
		}
	}
}

// Not negates a [ghost.Result].
func Not(result ghost.Result) ghost.Result {
	result.Ok = !result.Ok
	return result
}
