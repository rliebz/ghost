package ghost_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestRunner_Should(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		ok := testG.Should(ghost.Result{
			Ok:      true,
			Message: msg,
		})

		g.Should(be.True(ok))
		g.Should(be.SliceLen(0, mockT.logCalls))
		g.Should(be.SliceLen(0, mockT.failCalls))
		g.Should(be.SliceLen(0, mockT.failNowCalls))
	})

	t.Run("not ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		ok := testG.Should(ghost.Result{
			Ok:      false,
			Message: msg,
		})

		g.Should(be.False(ok))
		g.Should(be.SliceLen(0, mockT.failNowCalls))

		g.Should(be.DeepEqual(
			[][]any{{msg}},
			mockT.logCalls,
		))
		g.Should(be.SliceLen(1, mockT.failCalls))
	})
}

func TestRunner_ShouldNot(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		ok := testG.ShouldNot(ghost.Result{
			Ok:      true,
			Message: msg,
		})

		g.Should(be.False(ok))
		g.Should(be.SliceLen(0, mockT.failNowCalls))

		g.Should(be.DeepEqual(
			[][]any{{msg}},
			mockT.logCalls,
		))
		g.Should(be.SliceLen(1, mockT.failCalls))
	})

	t.Run("not ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		ok := testG.ShouldNot(ghost.Result{
			Ok:      false,
			Message: msg,
		})

		g.Should(be.True(ok))
		g.Should(be.SliceLen(0, mockT.logCalls))
		g.Should(be.SliceLen(0, mockT.failCalls))
		g.Should(be.SliceLen(0, mockT.failNowCalls))
	})
}

func TestRunner_Must(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		testG.Must(ghost.Result{
			Ok:      true,
			Message: msg,
		})

		g.Should(be.SliceLen(0, mockT.logCalls))
		g.Should(be.SliceLen(0, mockT.failCalls))
		g.Should(be.SliceLen(0, mockT.failNowCalls))
	})

	t.Run("not ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		testG.Must(ghost.Result{
			Ok:      false,
			Message: msg,
		})

		g.Should(be.SliceLen(1, mockT.failNowCalls))
		g.Should(be.DeepEqual(
			[][]any{{msg}},
			mockT.logCalls,
		))
	})
}

func TestRunner_MustNot(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		testG.MustNot(ghost.Result{
			Ok:      true,
			Message: msg,
		})

		g.Should(be.SliceLen(1, mockT.failNowCalls))
		g.Should(be.DeepEqual(
			[][]any{{msg}},
			mockT.logCalls,
		))
	})

	t.Run("not ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		testG.MustNot(ghost.Result{
			Ok:      false,
			Message: msg,
		})

		g.Should(be.SliceLen(0, mockT.logCalls))
		g.Should(be.SliceLen(0, mockT.failCalls))
		g.Should(be.SliceLen(0, mockT.failNowCalls))
	})
}

func TestRunner_NoError(t *testing.T) {
	g := ghost.New(t)

	mockT := newMockT()
	testG := ghost.New(mockT)

	myErr := errors.New("oh no")
	testG.NoError(myErr)

	if g.Should(be.SliceLen(1, mockT.logCalls)) {
		g.Should(be.DeepEqual(
			[]any{"myErr has error value: oh no"},
			mockT.logCalls[0],
		))
	}

	g.Should(be.SliceLen(0, mockT.failCalls))
	g.Should(be.SliceLen(1, mockT.failNowCalls))
}

type mockT struct {
	m sync.Mutex

	logCalls     [][]any
	failCalls    []struct{}
	failNowCalls []struct{}
}

var _ ghost.T = (*mockT)(nil)

func newMockT() *mockT {
	return &mockT{}
}

func (t *mockT) Log(args ...any) {
	t.m.Lock()
	defer t.m.Unlock()

	t.logCalls = append(t.logCalls, args)
}

func (t *mockT) Fail() {
	t.m.Lock()
	defer t.m.Unlock()

	t.failCalls = append(t.failCalls, struct{}{})
}

func (t *mockT) FailNow() {
	t.m.Lock()
	defer t.m.Unlock()

	t.failNowCalls = append(t.failNowCalls, struct{}{})
}
