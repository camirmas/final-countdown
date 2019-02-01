package countdown

import "testing"

func TestService(t *testing.T) {
	s := Service{}
	s.Start(nil)

	timer := s.StartTimer(1, 10)

	if timer.Id != 1 {
		t.Fatalf("Expected Id to be 1, got %d", timer.Id)
	}
	if timer.Duration != 10 {
		t.Fatalf("Expected Duration to be 10, got %d", timer.Duration)
	}
	if timer.TimeRemaining != 10 {
		t.Fatalf("Expected TimeRemaining to be 10, got %d", timer.TimeRemaining)
	}

	if _, ok := s.store.Get(1); !ok {
		t.Fatal("Expected Timer")
	}

	err := s.StopTimer(1)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := s.store.Get(1); ok {
		t.Fatalf("Expected error: %s", TimerNotFoundError{}.Error())
	}
}
