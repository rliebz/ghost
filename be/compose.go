package be

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/ghostlib"
)

// All asserts that every one of the provided assertions is true.
func All(results ...ghost.Result) ghost.Result {
	args := ghostlib.ArgsFromAST(results)
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
	args := ghostlib.ArgsFromAST(results)
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
	args []string,
) ghost.Result {
	if len(results) == 0 {
		return ghost.Result{
			Ok:      initial,
			Message: "no assertions were provided",
		}
	}

	out := ghost.Result{Ok: initial}
	for i, result := range results {
		out.Ok = apply(out.Ok, result.Ok)

		// Not sure why AST parsing would fail, but sometimes it does. Seems to be
		// environment dependent rather than code dependent.
		var arg string
		if len(args) > i {
			arg = fmt.Sprintf("`%s`", args[i])
		} else {
			arg = strconv.Itoa(i)
		}

		var b strings.Builder
		if i != 0 {
			b.WriteString("\n\n")
		}
		fmt.Fprintf(&b, "assertion %s is %t", arg, result.Ok)
		b.WriteString("\n\t")
		b.WriteString(indentString(result.Message))

		out.Message += b.String()
	}

	return out
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
