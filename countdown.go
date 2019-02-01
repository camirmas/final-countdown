/*
The countdown package provides a simple service for managing countdown timers.
*/
package countdown

// The timer can be expressed as a channel of Messages. The client will read from
// this channel to get updated time, and will send messages to it.
type Timer struct {
	// A unique identifier
	Id int
	// How long it should run
	Duration int
	// Time remaining (last stored)
	TimeRemaining int
	// Channel containing times remaining
	Channel chan int
	// Channel for communicating timer completion to parent service
	service chan int
	// Handles stored timers
	store Store
}

// Start runs the Timer
func (t *Timer) Start() {
	t.store.Add(t)
	// TODO: run countdown, pass time-remaining values to channel
}

// Stop cancels the Timer and deletes it
func (t *Timer) Stop() {
	t.store.Remove(t.Id)
	t.service <- t.Id
	close(t.Channel)
}

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

// General interface for storing Timer info. This allows for multiple backend
// implementations.
type Store interface {
	// List returns all timers
	List() []*Timer
	// Add adds a new Timer to the store
	Add(timer *Timer) error
	// Get gets a timer by id
	Get(id int) (*Timer, bool)
	// Update updates an existing Timer
	Update(timer *Timer) error
	// Remove removes a Timer from the store
	Remove(id int) error
}

// Default to simple map for storage.
type timerStore map[int]*Timer

func (ts timerStore) List() []*Timer {
	var timers []*Timer

	for _, v := range ts {
		timers = append(timers, v)
	}

	return timers
}

func (ts timerStore) Add(timer *Timer) error {
	ts[timer.Id] = timer

	return nil
}

func (ts timerStore) Get(id int) (*Timer, bool) {
	timer, ok := ts[id]

	return timer, ok
}

func (ts timerStore) Update(timer *Timer) error {
	ts[timer.Id] = timer

	return nil
}

func (ts timerStore) Remove(id int) error {
	delete(ts, id)

	return nil
}
