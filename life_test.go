package life

import (
	"errors"
	"testing"
	"time"
)

func TestCycleStart(t *testing.T) {
	var lc Cycle
	lc.Start(func(<-chan struct{}) error {
		return nil
	})
	if err := lc.Wait(); err != nil {
		t.Fatal("unexpected", err)
	}
}

func TestCycleStartStop(t *testing.T) {
	var lc Cycle
	lc.Start(func(stop <-chan struct{}) error {
		select {
		case <-stop:
			return nil
		case <-time.After(100 * time.Millisecond):
		}
		return errors.New("timeout")
	})
	lc.Stop()
	if err := lc.Wait(); err != nil {
		t.Fatal("unexpected", err)
	}
}
