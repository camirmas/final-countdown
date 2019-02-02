package countdown

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"testing"
)

var store Store

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown(store)

	os.Exit(code)
}

func TestService(t *testing.T) {
	t.Run("with Bolt DB", boltDBService)
}

func boltDBService(t *testing.T) {
	service := Service{}
	if err := service.Start(store, nil); err != nil {
		t.Fatal(err)
	}
	id := 1
	duration := 3

	if _, err := service.StartTimer(id, duration); err != nil {
		t.Fatal(err)
	}

	if _, err := service.StartTimer(id, duration); err == nil {
		expectErr(TimerExistsError{}, t)
	}

	if err := service.PauseTimer(id); err != nil {
		t.Fatal(err)
	}

	if err := service.ResumeTimer(id); err != nil {
		t.Fatal(err)
	}

	if err := service.CancelTimer(id); err != nil {
		t.Fatalf("Expected CancelTimer, got error: %s", err.Error())
	}

	// Make sure Timer got removed by Service after completion
	if _, err := service.GetTimer(id); err == nil {
		expectErr(TimerNotFoundError{}, t)
	}
}

func teardown(store Store) {
	store.(BoltStore).Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("timers"))

		if err != nil {
			log.Println(err)
		}

		return nil
	})
}

func expectErr(err error, t *testing.T) {
	t.Fatalf("Expected error: %s", err.Error())
}

func setup() {
	s, err := ConnectBoltStore("/tmp/countdown_test.db")
	if err != nil {
		log.Fatal(err)
	}
	store = *s
}
