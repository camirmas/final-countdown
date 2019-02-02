/*
The countdown package provides a simple service for managing countdown timers.
*/
package countdown

import (
	"log"
)

// The primary entry point for this library. It is responsible for maintaining
// the Timers and interacting with the client.
type Service struct {
	// Maintains current timer channels
	timers map[int]*Timer
	// Stores timer information
	store Store
	// Waits for messages from channels that are going to close.
	channel chan int
	// Store options for running the Service
	Options *Options
}

// Options for running the service
type Options struct {
	DbPath string
}

// Start runs the service
func (s *Service) Start(store Store, options *Options) error {
	s.Options = options
	if store == nil {
		db, err := ConnectBoltStore(options.DbPath)

		if err != nil {
			return err
		}

		s.store = *db
	} else {
		s.store = store
	}
	s.channel = make(chan int)
	s.timers = make(map[int]*Timer)

	if err := s.restoreTimers(); err != nil {
		log.Printf("Countdown package failed to restore timers: %s", err.Error())
	}

	go s.manageTimers()

	return nil
}

// StartTimer creates a new Timer for the Service, and Starts it.
func (s *Service) StartTimer(id, duration int) (*Timer, error) {
	timer := NewTimer(id, duration, s.store)
	if err := timer.Start(s.channel); err != nil {
		return nil, err
	}
	s.timers[id] = timer

	return timer, nil
}

// CancelTimer cancels an existing Timer if it exists.
func (s *Service) CancelTimer(id int) error {
	if timer, ok := s.timers[id]; ok {
		timer.Cancel()
		delete(s.timers, id)

		return nil
	} else {
		return TimerNotFoundError{}
	}
}

// PauseTimer temporarily pauses a Timer if it exists.
func (s *Service) PauseTimer(id int) error {
	if timer, err := s.GetTimer(id); err == nil {
		timer.Pause()

		return nil
	} else {
		return err
	}
}

// ResumeTimer temporarily pauses a Timer if it exists.
func (s *Service) ResumeTimer(id int) error {
	if timer, err := s.GetTimer(id); err == nil {
		timer.Resume(s.channel)

		return nil
	} else {
		return err
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
		timer.Start(s.channel)
		s.timers[id] = timer

		return timer, nil
	} else {
		return nil, err
	}
}

func (s *Service) manageTimers() {
	for id := range s.channel {
		s.CancelTimer(id)
	}
}

func (s *Service) restoreTimers() error {
	for _, timer := range s.store.ListTimers() {
		timerChannel := make(chan int)
		timer.store = s.store
		timer.channel = timerChannel
		if err := timer.Start(s.channel); err != nil {
			return err
		}
		s.timers[timer.Id] = timer
	}
	return nil
}
