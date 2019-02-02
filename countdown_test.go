package countdown

import (
	"github.com/boltdb/bolt"
	"testing"
	// "bytes"
	// "os"
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
	timer, err := s.StartTimer(id, 10)

	if err != nil {
		t.Fatal(err)
	}

	if timer.Id != id {
		t.Fatalf("Expected Id to be %d, got %d", id, timer.Id)
	}
	if timer.Duration != 10 {
		t.Fatalf("Expected Duration to be 10, got %d", timer.Duration)
	}
	if timer.TimeRemaining != 10 {
		t.Fatalf("Expected TimeRemaining to be 10, got %d", timer.TimeRemaining)
	}

	s.StopTimer(id)

	if _, err := s.GetTimer(id); err == nil {
		t.Fatalf("Expected error: %s", TimerNotFoundError{}.Error())
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
