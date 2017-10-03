// Package life manages the lifecycle of a set of goroutines.
// life draws strong inspiration from gopkg.in/tomb.v1 from Gustavo Niemeyer
// and github.com/oklog/oklog/pkg/group from Peter Bourgon.
package life

import "sync"

// A Cycle tracks the lifecycle of one or more goroutines.
type Cycle struct {
	wg sync.WaitGroup
	mu sync.Mutex // protects following fields

	procs []chan struct{}
	errc  chan error
}

// Start starts fn in a new goroutine.
func (lc *Cycle) Start(fn func(stop <-chan struct{}) error) {
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

// Stop signals to each function started by Start to exit by closing
// its respective stop channel then returns immediately.
func (lc *Cycle) Stop() {
	lc.mu.Lock()
	for _, g := range lc.procs {
		close(g)
	}
	lc.mu.Unlock()
}

// Wait waits until all functions passed to Start have exited.
// If any function returns an error, that error is propagated back to
// the caller of Wait. Subsequent errors are discared.
// If Start or Stop is called after Wait, they will panic.
func (lc *Cycle) Wait() error {
	lc.wg.Wait()
	close(lc.errc)
	return <-lc.errc
}
