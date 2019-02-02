package countdown

import (
	"fmt"
	"github.com/boltdb/bolt"
	"testing"
)

func TestService(t *testing.T) {
	t.Run("with Bolt DB", boltDBService)
}

func boltDBService(t *testing.T) {
	s := Service{}
	if err := s.Start(nil, true); err != nil {
		t.Fatal(err)
	}
	runTests(s, t)
	teardown(s, t)
}

func runTests(s Service, t *testing.T) {
	id := 1
	duration := 3
	timer, err := s.StartTimer(id, duration)

	if err != nil {
		t.Fatal(err)
	}

	if timer.Id != id {
		t.Fatalf("Expected Id to be %d, got %d", id, timer.Id)
	}
	if timer.Duration != duration {
		t.Fatalf("Expected Duration to be %d, got %d", duration, timer.Duration)
	}
	if timer.TimeRemaining != duration {
		t.Fatalf("Expected TimeRemaining to be %d, got %d", duration, timer.TimeRemaining)
	}

	if _, err := s.StartTimer(id, duration); err == nil {
		expectErr(TimerExistsError{}, t)
	}

	// Collect times remaining
	expected := [3]int{3, 2, 1}
	timesRemaining := make([]int, 0)
	for tr := range timer.Read() {
		fmt.Println(tr)
		timesRemaining = append(timesRemaining, tr)
	}

	// Make sure values are the same
	for i, v := range expected {
		if v != timesRemaining[i] {
			t.Fatalf("Expected %v, got %v", expected, timesRemaining)
		}
	}

	// Make sure Timer got removed by Service after completion
	if _, err := s.GetTimer(id); err == nil {
		expectErr(TimerNotFoundError{}, t)
	}
}

func teardown(s Service, t *testing.T) {
	s.store.(BoltStore).Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("timers"))

		if err != nil {
			t.Error(err)
		}

		return nil
	})
}

func expectErr(err error, t *testing.T) {
	t.Fatalf("Expected error: %s", err.Error())
}
