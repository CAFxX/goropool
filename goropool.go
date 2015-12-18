package goropool // ゴロプル

import (
	"runtime"
	"sync"
)

// NewDefaultPool creates a new goroutine pool with runtime.NumCPU() workers and
// a queue of size 0 (i.e. a worker must be idle for the send to the queue to
// succeed). See NewPool() for details about the returned values.
func NewDefaultPool() (chan<- func(), <-chan error) {
	return NewPool(runtime.NumCPU(), 0)
}

// NewPool creates a bounded goroutine pool with workers goroutine accepting
// work on a queue with queueSize elements. The first channel returned is the
// queue jobs should be submitted on. Closing this channel signals that no more
// jobs will be sent to the pool. The second channel signals completion of all
// queued jobs (this can only happen after the queue channel has been closed).
// Each pool is made up of workers+1 goroutines. NewPool returns immediately.
func NewPool(workers, queueSize int) (chan<- func(), <-chan error) {
	queue := make(chan func(), queueSize)
	done := make(chan error)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
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
