// Package goropool implements a dead-simple bounded goroutine pool. It is
// mostly useful when dealing with blocking I/O calls.
package goropool // ゴロプル

import (
	"runtime"
	"sync"
)

// NewPool creates and starts a bounded goroutine pool with numWorkers goroutine
// accepting work on a queue with queueSize elements. Each pool is made up of
// numWorkers+1 goroutines. NewPool returns immediately.
//
// Two channels are returned: queue and done. The former is the channel jobs
// should be submitted on. A job is simply a func(). Closing the queue channel
// signals to the pool that no more jobs will be enqueued. Once the queue
// channel has been closed and all queued jobs have been completed the pool will
// close the second channel (done) to signal that the pool has shut down.
//
// The done channel is of type error for future extensibility, currently no
// error will be returned under any circumstances. The pool makes no attempt to
// recover panicking jobs.
func NewPool(numWorkers, queueSize int) (chan<- func(), <-chan error) {
	queue := make(chan func(), queueSize)
	done := make(chan error)
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for job := range queue {
				job()
			}
		}()
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	return queue, done
}

// NewDefaultPool creates a new goroutine pool with runtime.NumCPU() workers and
// a queue of size 0 (i.e. a worker must be idle for the send to the queue to
// succeed). See NewPool() for details about the returned values.
func NewDefaultPool() (chan<- func(), <-chan error) {
	return NewPool(runtime.NumCPU(), 0)
}
