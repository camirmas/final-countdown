package countdown

import (
	"log"
	"testing"
)

func TestTimer(t *testing.T) {
	t.Run("NewTimer", newTimer)
	t.Run("till complete", completeTimer)
	t.Run("Stop/Resume", startStopTimer)
	t.Run("Cancel", cancelTimer)
}

func newTimer(t *testing.T) {
	id := 1
	duration := 3
	timer := NewTimer(id, duration, nil)

	if timer.Id != id {
		t.Errorf("Expected Id to be %d, got %d", id, timer.Id)
	}
	if timer.Duration != duration {
		t.Errorf("Expected Duration to be %d, got %d", duration, timer.Duration)
	}
	if status := timer.Status(); status != Paused {
		t.Errorf("Expected status %s, got %s", Paused, status)
	}
}

func completeTimer(t *testing.T) {
	id := 1
	duration := 3
	timer := NewTimer(id, duration, nil)

	timer.Start(nil)

	// Collect times remaining
	expected := [3]int{3, 2, 1}
	timesRemaining := make([]int, 0)
	for tr := range timer.Channel() {
		log.Printf("Time remaining: %d\n", tr)
		timesRemaining = append(timesRemaining, tr)
	}

	// Make sure values are the same
	for i, v := range expected {
		if v != timesRemaining[i] {
			t.Fatalf("Expected %v, got %v", expected, timesRemaining)
		}
	}
}

func startStopTimer(t *testing.T) {
	id := 1
	duration := 3
	timer := NewTimer(id, duration, nil)

	timer.Start(nil)

	timer.Pause()

	if status := timer.Status(); status != Paused {
		t.Fatalf("Expected status %s, got %s", Paused, status)
	}

	timer.Resume(nil)

	if status := timer.Status(); status != Running {
		t.Fatalf("Expected status %s, got %s", Running, status)
	}
}

func cancelTimer(t *testing.T) {
	id := 1
	duration := 3
	timer := NewTimer(id, duration, nil)

	timer.Start(nil)

	timer.Cancel()

	if status := timer.Status(); status != Cancelled {
		t.Fatalf("Expected status %s, got %s", Cancelled, status)
	}
}
