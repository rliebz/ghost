package ghost_test

import (
	"testing"

	"github.com/rliebz/ghost"
)

func TestBeNil(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var v interface{}

		result := ghost.BeNil(v)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is nil", result.Message))

		result = ghost.BeNil(nil)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("nil is nil", result.Message))
	})

	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		var v int

		result := ghost.BeNil(v)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is 0, not nil", result.Message))

		result = ghost.BeNil(-1 + 1)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("-1 + 1 is 0, not nil", result.Message))
	})
}

func TestBeTrue(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		g := ghost.New(t)

		v := true
		result := ghost.BeTrue(v)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is true", result.Message))

		result = ghost.BeTrue(true)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("true is true", result.Message))
	})

	t.Run("false", func(t *testing.T) {
		g := ghost.New(t)

		v := false
		result := ghost.BeTrue(v)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is false", result.Message))

		result = ghost.BeTrue(false)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("false is false", result.Message))
	})
}

func TestBeZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		g := ghost.New(t)

		var v int
		result := ghost.BeZero(v)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is the zero value", result.Message))

		result = ghost.BeZero(0)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("0 is the zero value", result.Message))
	})

	t.Run("non-zero", func(t *testing.T) {
		g := ghost.New(t)

		v := 1
		result := ghost.BeZero(v)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is non-zero\nvalue: 1", result.Message))

		result = ghost.BeZero(1)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("1 is non-zero", result.Message))
	})
}

func TestContain(t *testing.T) {
	t.Run("contains <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 2

		result := ghost.Contain(slice, elem)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`slice contains elem
slice:   [1 2 3]
element: 2
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3}, 2)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`[]int{1, 2, 3} contains 2
slice:   [1 2 3]
element: 2
`, result.Message))
	})

	t.Run("contains > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 2

		result := ghost.Contain(slice, elem)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`slice contains elem
slice:   [
	1
>	2
	3
	4
]
element: 2
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3, 4}, 2)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`[]int{1, 2, 3, 4} contains 2
slice:   [
	1
>	2
	3
	4
]
element: 2
`, result.Message))
	})

	t.Run("does not contain <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 5

		result := ghost.Contain(slice, elem)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`slice does not contain elem
slice:   [1 2 3]
element: 5
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3}, 5)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`[]int{1, 2, 3} does not contain 5
slice:   [1 2 3]
element: 5
`, result.Message))
	})

	t.Run("does not contain > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 5

		result := ghost.Contain(slice, elem)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`slice does not contain elem
slice:   [
	1
	2
	3
	4
]
element: 5
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3, 4}, 5)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`[]int{1, 2, 3, 4} does not contain 5
slice:   [
	1
	2
	3
	4
]
element: 5
`, result.Message))
	})
}
