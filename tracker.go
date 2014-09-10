package tracker

import (
	"sync"
	"sync/atomic"
	"time"
)

// Tracker is designed to signal a dynamic number of
// go routines to exit cleanly via a channel, boolean, or both.
type Tracker struct {
	toSignal []*chan bool
	lock     sync.Mutex
	count    int64
	stopping bool
}

// Join returns a channel that will fire when the tracker
// wants all go routines to exit. In order for the tracker
// to exit, there must be a matching call to Leave
func (t *Tracker) Join() *chan bool {
	ch := make(chan bool, 1)
	atomic.AddInt64(&t.count, int64(1))
	t.lock.Lock()
	defer t.lock.Unlock()
	t.toSignal = append(t.toSignal, &ch)
	return &ch
}

// Leave tells the tracker that a tracked go routine is exiting
// either because it was asked to do so or simply because it is finished.
func (t *Tracker) Leave() {
	atomic.AddInt64(&t.count, int64(-1))
}

// IsRunning returns ture if KillAll() has not been called and
// the tracked goroutines should continue to run.
func (t *Tracker) IsRunning() bool {
	return !t.stopping
}

// Count returns how many go routines are currently being tracked (those
// that have called Join() but not Leave())
func (t *Tracker) Count() int64 {
	return atomic.LoadInt64(&t.count)
}

// KillAll will signal tracked go routines to quite via the IsRunning() boolean
// and the channel that was returned to the thread.
func (t *Tracker) KillAll() {
	t.stopping = true
	for _, v := range t.toSignal {
		*v <- true
	}
	for atomic.LoadInt64(&t.count) != 0 {
		time.Sleep(5 * time.Millisecond)
	}
}
