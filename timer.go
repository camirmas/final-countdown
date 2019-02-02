package countdown

import (
	"encoding/json"
)

// The timer can be expressed as a channel of Messages. The client will read from
// this channel to get updated time, and will send messages to it.
type Timer struct {
	// A unique identifier
	Id int `json:"id"`
	// How long it should run
	Duration int `json:"duration"`
	// Time remaining (last stored)
	TimeRemaining int `json:"timeRemaining"`
	// Channel containing times remaining
	channel chan int
	// Channel for communicating timer completion to parent service
	service chan int
	// Handles stored timers
	store Store
}

// Start runs the Timer
func (t *Timer) Start() error {
	return t.store.AddTimer(t)
	// TODO: run countdown, pass time-remaining values to channel
}

// Stop ends the Timer
func (t *Timer) Stop() error {
	if err := t.store.RemoveTimer(t.Id); err != nil {
		return err
	}
	close(t.channel)

	return nil
}

// Gets the timer
func (t *Timer) Read() chan int {
	return t.channel
}

func (t *Timer) serialize() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Timer) deserialize(b []byte) error {
	return json.Unmarshal(b, t)
}

func (t *Timer) complete() {
	t.service <- t.Id
}
