package countdown

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

// Default to simple map for storage. This should eventually
// be replaced with a serious default like BoltDB.
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
