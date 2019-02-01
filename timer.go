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
