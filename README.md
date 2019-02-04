# :alarm_clock: final-countdown
> Simple service for managing countdown timers

### Purpose

This library is intended to supplement applications that need to maintain timers, e.g. a web-based board game like [GoStop](https://github.com/camirmas/go-stop-go). 

### Installation

`go get github.com/camirmas/final-countdown`

### Usage

There are three main entities of note in this library: 
- `Service`, which is is a higher-level abstraction that can manage concurrent `Timer`s
- `Timer`, which runs a countdown sequence
- `Store`, which is an interface for storing multiple timers. The default for a `Service` is `BoltStore`, which uses key/value database [Bolt](https://github.com/boltdb/bolt#lmdb)

##### Starting a `Service`
```go
package main

import "github.com/camirmas/final-countdown"

service := Service{}

// The first argument is an optional `Store` implementation.
// If you do not want to supply your own, you must include an
// Options value.
service.Start(nil, Options{DbPath: "countdown.db"}) // is the same as service.Start(nil, Options{})
```

##### Manage `Timer`s
```go
package main

import "github.com/camirmas/final-countdown"

service := countdown.Service{}
service.Start(nil, countdown.Options{DbPath: "countdown.db"}) // is the same as service.Start(nil, countdown.Options{})

id := 1
duration := 10

timer, err := service.StartTimer(id, duration)

if err != nil {
  // handle err
}

if err := service.PauseTimer(id); err != nil {
  // handle err
}

if err := service.ResumeTimer(id); err != nil {
  // handle err
}

if err := service.CancelTimer(id); err != nil {
  // handle err
}
```

##### Interacting with a `Timer`
Now, the client wants to be notified of timer changes, and it uses a provided channel to do so:
```go
timer, _ := service.GetTimer(id)

/// This will keep looping through the channel (chan int) until it closes
for tr := range timer.Channel() {
  // do stuff here
}
```
Here's the lifecycle of a managed `Timer`:
1. created by a `Service` via `StartTimer`
2. runs a countdown in a separate goroutine
3. Either:
  - finishes its countdown, then notifies the parent `Service`, which deletes it
  - gets cancelled by the `Service` before it finishes the countdown
  
The intended behavior of a managed `Timer` is such that it should spin up once, then either get cancelled by the `Service`, or complete its countdown and tell the `Service` to destroy it.

##### Using a standalone `Timer`
While a `Service` is intended to handle most aspects of using `Timer`s, especially many at a time, a `Timer` can still be run on its own. This can be used in a case where an application has its own way that it wishes to manage them, or simply doesn't need the extra overhead of a `Service`.

```go
id := 1
duration := 3

// can optionally provide a Store here
timer := NewTimer(id, duration, nil)

// can optionally provide a chan int here, the Timer will send its Id to that channel upon countdown completion
timer.Start(nil)
timer.Pause()
timer.Resume()
timer.Cancel()
```

##### Using a different `Store`
The library supports using different backends to manage `Timer`s. This can be achieved by simply implementing the following interface:
```go
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
```
That implementation can then be passed into a new `Service`:
```service.Start(myStore, nil)```
