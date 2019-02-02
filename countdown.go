/*
The countdown package provides a simple service for managing countdown timers.
*/
package countdown

import "log"

// The primary entry point for this library. It is responsible for maintaining
// the Timers and interacting with the client.
type Service struct {
	// Maintains current timer channels
	timers map[int]*Timer
	// Stores timer information
	store Store
	// Waits for messages from channels that are going to close.
	channel chan int
}

// Start runs the service
func (s *Service) Start(store *Store, test bool) error {
	if store == nil {
		var path string
		if test {
			path = "/tmp/countdown_test.db"
		}

		db, err := connectBoltStore(path)

		if err != nil {
			return err
		}

		s.store = *db
	}
	s.channel = make(chan int)
	s.timers = make(map[int]*Timer)

	if err := s.restoreTimers(); err != nil {
		log.Printf("Countdown package failed to restore timers: %s", err.Error())
	}

	go s.manageTimers()

	return nil
}

// StartTimer creates a new Timer for the Service.
func (s *Service) StartTimer(id int, duration int) (*Timer, error) {
	timerChannel := make(chan int)
	timer := &Timer{id, duration, duration, timerChannel, s.channel, s.store}
	if err := timer.Start(); err != nil {
		return nil, err
	}
	s.timers[id] = timer

	return timer, nil
}

// StopTimer stops an existing Timer if it exists.
func (s *Service) StopTimer(id int) error {
	if timer, ok := s.timers[id]; ok {
		timer.Stop()
		delete(s.timers, id)

		return nil
	} else {
		return TimerNotFoundError{}
	}
}

// GetTimer retrieves an existing Timer. First it checks if there is a Timer in
// memory. Then it checks the DB, and if it finds one there, it sets up a new one
// because the old channels are gone.
func (s *Service) GetTimer(id int) (*Timer, error) {
	if timer, ok := s.timers[id]; ok {
		return timer, nil
	}
	if timer, err := s.store.GetTimer(id); err == nil {
		timer.store = s.store
		timer.channel = make(chan int)
		timer.Start()
		s.timers[id] = timer

		return timer, nil
	} else {
		return nil, err
	}
}

func (s *Service) manageTimers() {
	for id := range s.channel {
		s.StopTimer(id)
	}
}

func (s *Service) restoreTimers() error {
	for _, timer := range s.store.ListTimers() {
		timerChannel := make(chan int)
		timer.store = s.store
		timer.channel = timerChannel
		if err := timer.Start(); err != nil {
			return err
		}
		s.timers[timer.Id] = timer
	}
	return nil
}
