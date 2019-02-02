package countdown

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
	"time"
)

// General interface for storing Timer info. This allows for multiple backend
// implementations.
type Store interface {
	// ListTimers returns all timers
	ListTimers() []*Timer
	// AddTimer adds a new Timer to the store
	AddTimer(timer *Timer) error
	// GetTimer gets a timer by id
	GetTimer(id int) (*Timer, error)
	// UpdateTimer updates an existing Timer
	UpdateTimer(timer *Timer) error
	// Remove removes a Timer from the store
	RemoveTimer(id int) error
}

type BoltStore struct {
	*bolt.DB
}

func (db BoltStore) ListTimers() []*Timer {
	var timers []*Timer
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("timers"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			timer := &Timer{}
			timer.deserialize(v)
			timers = append(timers, timer)
		}

		return nil
	})

	return timers
}

func (db BoltStore) AddTimer(t *Timer) error {
	if _, err := db.GetTimer(t.Id); err == nil {
		return TimerExistsError{}
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("timers"))
		tSer, err := t.serialize()
		if err != nil {
			return err
		}
		id := intToBytes(t.Id)
		return b.Put(id, tSer)
	}); err != nil {
		return err
	}
	return nil
}

func (db BoltStore) GetTimer(id int) (*Timer, error) {
	var timer Timer
	if err := db.View(func(tx *bolt.Tx) error {
		id := intToBytes(id)
		value := tx.Bucket([]byte("timers")).Get(id)

		if value == nil {
			return TimerNotFoundError{}
		} else {
			timer = Timer{}
			err := timer.deserialize(value)
			timer.store = db

			return err
		}
	}); err != nil {
		return nil, err
	}
	return &timer, nil
}

func (db BoltStore) UpdateTimer(t *Timer) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("timers"))
		tSer, err := t.serialize()
		if err != nil {
			return err
		}
		id := intToBytes(t.Id)
		return b.Put(id, tSer)
	}); err != nil {
		return err
	}
	return nil
}

func (db BoltStore) RemoveTimer(id int) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("timers"))
		id := intToBytes(id)
		return b.Delete(id)
	}); err != nil {
		return err
	}
	return nil
}

func ConnectBoltStore(filepath string) (*BoltStore, error) {
	if filepath == "" {
		filepath = "countdown.db"
	}
	db, err := bolt.Open(
		filepath,
		0600,
		&bolt.Options{Timeout: 1 * time.Second},
	)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("timers"))
		return nil
	}); err != nil {
		return nil, err
	}

	return &BoltStore{db}, nil
}

func intToBytes(i int) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, uint64(i))

	return res
}
