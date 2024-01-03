package ghost_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestGhost_Should(t *testing.T) {
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
		g.Should(be.SliceLen(mockT.logCalls, 0))
		g.Should(be.SliceLen(mockT.failCalls, 0))
		g.Should(be.SliceLen(mockT.failNowCalls, 0))
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
		g.Should(be.SliceLen(mockT.failNowCalls, 0))

		g.Should(be.DeepEqual(
			[][]any{{msg}},
			mockT.logCalls,
		))
		g.Should(be.SliceLen(mockT.failCalls, 1))
	})
}

func TestGhost_Must(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := ghost.New(t)

		mockT := newMockT()
		testG := ghost.New(mockT)
		msg := "some message"

		testG.Must(ghost.Result{
			Ok:      true,
			Message: msg,
		})

		g.Should(be.SliceLen(mockT.logCalls, 0))
		g.Should(be.SliceLen(mockT.failCalls, 0))
		g.Should(be.SliceLen(mockT.failNowCalls, 0))
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

		g.Should(be.SliceLen(mockT.failNowCalls, 1))
		g.Should(be.DeepEqual(
			mockT.logCalls,
			[][]any{{msg}},
		))
	})
}

func TestGhost_NoError(t *testing.T) {
	g := ghost.New(t)

	mockT := newMockT()
	testG := ghost.New(mockT)

	myErr := errors.New("oh no")
	testG.NoError(myErr)

	if g.Should(be.SliceLen(mockT.logCalls, 1)) {
		g.Should(be.DeepEqual(
			mockT.logCalls[0],
			[]any{"myErr has error value: oh no"},
		))
	}

	g.Should(be.SliceLen(mockT.failCalls, 0))
	g.Should(be.SliceLen(mockT.failNowCalls, 1))
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
