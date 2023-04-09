package ghost_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/rliebz/ghost"
)

func TestRunner_NoError(t *testing.T) {
	g := ghost.New(t)

	mockT := newMockT()
	testG := ghost.New(mockT)

	myErr := errors.New("oh no")
	testG.NoError(myErr)

	if g.Should(ghost.Len(1, mockT.logCalls)) {
		g.Should(ghost.DeepEqual(
			[]any{"myErr has error value: oh no"},
			mockT.logCalls[0],
		))
	}

	g.Should(ghost.Len(0, mockT.failCalls))
	g.Should(ghost.Len(1, mockT.failNowCalls))
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
