// Package life manages the lifecycle of a set of goroutines, and their associated resources.
package life

import "sync"

// A Cycle tracks the lifecycle of one or more goroutines
type Cycle struct {
	wg sync.WaitGroup
	mu sync.Mutex // protects following fields

	procs []chan struct{}
	errc  chan error
}

func (lc *Cycle) Start(fn func(<-chan struct{}) error) {
	lc.mu.Lock()
	if lc.errc == nil {
		lc.errc = make(chan error, 1)
	}
	stop := make(chan struct{})
	lc.procs = append(lc.procs, stop)
	lc.mu.Unlock()

	lc.wg.Add(1)
	go func() {
		defer lc.wg.Done()

		err := fn(stop)
		select {
		case lc.errc <- err:
			// sent our error
		default:
			// someone else has sent their error first
		}
	}()
}

func (lc *Cycle) Stop() {
	lc.mu.Lock()
	for _, g := range lc.procs {
		close(g)
	}
	lc.mu.Unlock()
}

func (lc *Cycle) Wait() error {
	lc.wg.Wait()
	close(lc.errc)
	return <-lc.errc
}
