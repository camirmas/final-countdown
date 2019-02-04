package countdown

import (
	"testing"
)

func TestStore(t *testing.T) {
	timer := NewTimer(1, 1, nil)

	if err := store.AddTimer(timer); err != nil {
		t.Fatal(err)
	}

	if _, err := store.GetTimer(timer.Id); err != nil {
		t.Error(err)
	}

	timers := store.ListTimers()
	lenTimers := len(timers)
	if lenTimers != 1 {
		t.Errorf("Expected 1 timer, got %d", lenTimers)
	}

	if err := store.UpdateTimer(timer); err != nil {
		t.Error(err)
	}

	if err := store.RemoveTimer(timer.Id); err != nil {
		t.Error(err)
	}

	if _, err := store.GetTimer(timer.Id); err == nil {
		expectErr(TimerNotFoundError{}, t)
	}
}
