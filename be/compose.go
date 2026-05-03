package be

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/ghostlib"
)

// All asserts that every one of the provided assertions is true.
func All(results ...ghost.Result) ghost.Result {
	argsOfAny := make([]any, 0, len(results))
	for _, result := range results {
		argsOfAny = append(argsOfAny, result)
	}
	args := ghostlib.ArgsFromAST(argsOfAny...)

	return applyVariadicBooleanLogic(
		true,
		func(acc, val bool) bool {
			return acc && val
		},
		results,
		args,
	)
}

// Any asserts that at least one of the provided assertions is true.
func Any(results ...ghost.Result) ghost.Result {
	argsOfAny := make([]any, 0, len(results))
	for _, result := range results {
		argsOfAny = append(argsOfAny, result)
	}
	args := ghostlib.ArgsFromAST(argsOfAny...)

	return applyVariadicBooleanLogic(
		false,
		func(acc, val bool) bool {
			return acc || val
		},
		results,
		args,
	)
}

func applyVariadicBooleanLogic(
	initial bool,
	apply func(acc, val bool) bool,
	results []ghost.Result,
	args ghostlib.Args,
) ghost.Result {
	if len(results) == 0 {
		return ghost.Result{
			Ok:      initial,
			Message: func() string { return "no assertions were provided" },
		}
	}

	ok := initial
	for _, result := range results {
		ok = apply(ok, result.Ok)
	}

	return ghost.Result{
		Ok: ok,
		Message: func() string {
			var b strings.Builder
			for i, result := range results {
				if i != 0 {
					b.WriteString("\n\n")
				}
				fmt.Fprintf(&b, "assertion `%s` is %t", args.Get(i), result.Ok)
				b.WriteString("\n\t")
				b.WriteString(indentString(result.Message()))
			}
			return b.String()
		},
	}
}

var reWhitespaceLine = regexp.MustCompile(`\n[ \t]+\n`)

func indentString(s string) string {
	s = strings.ReplaceAll(s, "\n", "\n\t")
	s = reWhitespaceLine.ReplaceAllString(s, "\n\n")
	s = strings.TrimSpace(s)
	return s
}

// Eventually asserts that a function eventually returns an Ok [ghost.Result].
func Eventually(
	f func() ghost.Result,
	timeout time.Duration,
	interval time.Duration,
) ghost.Result {
	args := ghostlib.ArgsFromAST(f, timeout, interval)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	lastRun := ghost.Result{
		Ok: false,
		Message: func() string {
			return fmt.Sprintf("%s did not return value within %s timeout", args.Get(0), timeout)
		},
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
