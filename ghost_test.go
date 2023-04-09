package ghost_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/rliebz/ghost"
)

func TestRunner_NoErr(t *testing.T) {
	g := ghost.New(t)

	mockT := newMockT()
	testG := ghost.New(mockT)

	myErr := errors.New("oh no")
	testG.NoError(myErr)

	if g.Should(ghost.Equal(1, len(mockT.logCalls))) {
		g.Should(ghost.DeepEqual(
			[]any{"myErr has error value: oh no"},
			mockT.logCalls[0],
		))
	}

	g.Should(ghost.Equal(0, len(mockT.failCalls)))
	g.Should(ghost.Equal(1, len(mockT.failNowCalls)))
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
	t.logCalls = append(t.logCalls, args)
}

func (t *mockT) Fail() {
	t.failCalls = append(t.failCalls, struct{}{})
}

func (t *mockT) FailNow() {
	t.failNowCalls = append(t.failNowCalls, struct{}{})
}
