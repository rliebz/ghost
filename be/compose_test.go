package be_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestAll(t *testing.T) {
	t.Run("no arguments passed", func(t *testing.T) {
		g := ghost.New(t)

		result := be.All()
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "no assertions were provided"))
	})

	t.Run("one valid", func(t *testing.T) {
		g := ghost.New(t)

		result := be.All(
			be.Equal(1, 0),
			be.Equal(1, 1),
			be.Equal(1, 2),
		)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			fmt.Sprintf(`assertion %s is false
	1 != 0
	got:  1
	want: 0

assertion %s is true
	1 == 1

assertion %s is false
	1 != 2
	got:  1
	want: 2`,
				"`be.Equal(1, 0)`",
				"`be.Equal(1, 1)`",
				"`be.Equal(1, 2)`",
			),
		))
	})

	t.Run("all valid", func(t *testing.T) {
		g := ghost.New(t)

		result := be.All(
			be.Equal(1, 1),
			be.Equal(2, 2),
		)

		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			fmt.Sprintf("assertion %s is true"+`
	1 == 1

assertion %s is true
	2 == 2`,
				"`be.Equal(1, 1)`",
				"`be.Equal(2, 2)`",
			),
		))
	})

	t.Run("nested", func(t *testing.T) {
		g := ghost.New(t)

		result := be.Any(
			be.Any(
				be.Equal(1, 0),
				be.Equal(1, 2),
			),
		)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			fmt.Sprintf(`assertion %s is false
	assertion %s is false
		1 != 0
		got:  1
		want: 0

	assertion %s is false
		1 != 2
		got:  1
		want: 2`,
				"`be.Any(be.Equal(1, 0), be.Equal(1, 2))`",
				"`be.Equal(1, 0)`",
				"`be.Equal(1, 2)`",
			),
		))
	})
}

func TestAny(t *testing.T) {
	t.Run("no arguments passed", func(t *testing.T) {
		g := ghost.New(t)

		result := be.Any()
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "no assertions were provided"))
	})

	t.Run("one valid", func(t *testing.T) {
		g := ghost.New(t)

		result := be.Any(
			be.Equal(1, 0),
			be.Equal(1, 1),
			be.Equal(1, 2),
		)

		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			fmt.Sprintf(`assertion %s is false
	1 != 0
	got:  1
	want: 0

assertion %s is true
	1 == 1

assertion %s is false
	1 != 2
	got:  1
	want: 2`,
				"`be.Equal(1, 0)`",
				"`be.Equal(1, 1)`",
				"`be.Equal(1, 2)`",
			),
		))
	})

	t.Run("none valid", func(t *testing.T) {
		g := ghost.New(t)

		result := be.Any(
			be.Equal(1, 0),
			be.Equal(1, 2),
		)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			fmt.Sprintf("assertion %s is false"+`
	1 != 0
	got:  1
	want: 0

assertion %s is false
	1 != 2
	got:  1
	want: 2`,
				"`be.Equal(1, 0)`",
				"`be.Equal(1, 2)`",
			),
		))
	})

	t.Run("nested", func(t *testing.T) {
		g := ghost.New(t)

		result := be.Any(
			be.Any(
				be.Equal(1, 0),
				be.Equal(1, 2),
			),
		)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			fmt.Sprintf(`assertion %s is false
	assertion %s is false
		1 != 0
		got:  1
		want: 0

	assertion %s is false
		1 != 2
		got:  1
		want: 2`,
				"`be.Any(be.Equal(1, 0), be.Equal(1, 2))`",
				"`be.Equal(1, 0)`",
				"`be.Equal(1, 2)`",
			),
		))
	})
}

func TestEventually(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := ghost.New(t)

		count := 0
		result := be.Eventually(func() ghost.Result {
			count++
			return be.Equal(count, 3)
		}, 100*time.Millisecond, 5*time.Millisecond)

		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `count == 3`))
	})

	t.Run("not ok", func(t *testing.T) {
		g := ghost.New(t)

		count := 0
		result := be.Eventually(func() ghost.Result {
			count++
			return be.Equal(count, -1)
		}, 100*time.Millisecond, 5*time.Millisecond)

		g.Should(be.False(result.Ok))
		// TODO: This is a good signal for a native regexp assertion
		matched, err := regexp.MatchString(`count != -1
got:  \d+
want: -1
`, result.Message)
		g.NoError(err)
		g.Should(be.True(matched))
	})

	t.Run("timeout", func(t *testing.T) {
		g := ghost.New(t)

		result := be.Eventually(func() ghost.Result {
			time.Sleep(100 * time.Millisecond)
			return be.True(true)
		}, 10*time.Millisecond, 100*time.Millisecond)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `func() ghost.Result {
	time.Sleep(100 * time.Millisecond)
	return be.True(true)
} did not return value within 10ms timeout`))
	})
}

func TestNot(t *testing.T) {
	g := ghost.New(t)

	message := "some message"

	result := ghost.Result{
		Ok:      true,
		Message: message,
	}

	negated := be.Not(result)
	g.Should(be.False(negated.Ok))
	g.Should(be.Equal(negated.Message, message))

	doubleNegated := be.Not(negated)
	g.Should(be.Equal(doubleNegated, result))
}
