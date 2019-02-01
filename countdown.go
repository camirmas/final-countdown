/*
The countdown package provides a simple service for managing countdown timers.
*/
package countdown

// The primary entry point for this library. It is responsible for maintaining
// the Timers and interacting with the client.
type Service struct {
	// Stores timers
	store Store
	// Waits for messages from channels that are going to close.
	channel chan int
}

// Start runs the service
func (s *Service) Start(store *Store) {
	if store == nil {
		ts := timerStore{}
		s.store = ts
	} else {
		s.store = *store
	}
	s.channel = make(chan int)

	for _, timer := range s.store.List() {
		timer.Start()
	}

	go s.manageTimers()
}

// StartTimer creates a new Timer for the Service.
func (s *Service) StartTimer(id int, duration int) *Timer {
	clientChannel := make(chan int)
	timer := &Timer{id, duration, duration, clientChannel, s.channel, s.store}
	timer.Start()

	return timer
}

// StopTimer creates a new Timer for the Service.
func (s *Service) StopTimer(id int) error {
	if timer, ok := s.store.Get(id); ok {
		timer.Stop()
	} else {
		return TimerNotFoundError{}
	}

	return nil
}

func (s *Service) manageTimers() {
	for id := range s.channel {
		s.store.Remove(id)
	}
}
