package countdown

import (
	"encoding/json"
	"sync"
	"time"
)

const (
	Running   string = "running"
	Paused    string = "paused"
	Cancelled string = "cancelled"
)

// The timer can be expressed as a channel of Messages. The client will read from
// this channel to get updated time, and will send messages to it.
type Timer struct {
	// A unique identifier
	Id int `json:"id"`
	// How long it should run
	Duration int `json:"duration"`
	// Time remaining (last stored)
	timeRemaining int `json:"timeRemaining"`
	// Channel containing times remaining
	channel chan int
	// Handles stored timers
	store Store
	//
	mux sync.RWMutex
	// Manages running/stopping the timer
	status string `json:"status"`
}

func NewTimer(id, duration int, store Store) *Timer {
	channel := make(chan int)
	return &Timer{
		id,
		duration,
		duration,
		channel,
		store,
		sync.RWMutex{},
		Paused,
	}
}

func (t *Timer) Start(service chan int) error {
	t.mux.Lock()
	t.status = Running
	t.mux.Unlock()

	if t.store != nil {
		if err := t.store.AddTimer(t); err != nil {
			return err
		}
	}
	go t.countdown(service)

	return nil
}

func (t *Timer) Cancel() error {
	t.mux.Lock()
	t.status = Cancelled
	t.mux.Unlock()

	if t.store != nil {
		if err := t.store.RemoveTimer(t.Id); err != nil {
			return err
		}
	}

	return nil
}

func (t *Timer) Pause() {
	t.mux.Lock()
	t.status = Paused
	t.mux.Unlock()
}

func (t *Timer) Resume(service chan int) {
	t.mux.Lock()
	t.status = Running
	t.mux.Unlock()

	go t.countdown(service)
}

// Gets the timer channel
func (t *Timer) Channel() chan int {
	return t.channel
}

func (t *Timer) Status() string {
	t.mux.RLock()
	defer t.mux.RUnlock()

	return t.status
}

func (t *Timer) serialize() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Timer) deserialize(b []byte) error {
	return json.Unmarshal(b, t)
}

func (t *Timer) countdown(service chan int) error {
	for t.timeRemaining > 0 {
		switch t.Status() {
		case Running:
			t.channel <- t.timeRemaining
			t.timeRemaining--
			time.Sleep(1 * time.Second)
			if t.store != nil {
				if err := t.store.UpdateTimer(t); err != nil {
					return err
				}
			}
		case Paused:
			return nil
		case Cancelled:
			break
		}
	}
	t.complete(service)

	return nil
}

func (t *Timer) complete(service chan int) {
	if service != nil {
		service <- t.Id
	} else {
		close(t.channel)
	}
}
