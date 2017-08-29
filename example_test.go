package life_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/pkg/life"
)

func ExampleCycle_Start() {
	var lc life.Cycle
	lc.Start(func(<-chan struct{}) error {
		fmt.Println("A")
		return nil
	})
	lc.Start(func(<-chan struct{}) error {
		fmt.Println("B")
		return errors.New("B failed")
	})
	lc.Start(func(<-chan struct{}) error {
		fmt.Println("C")
		return nil
	})
	err := lc.Wait()
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleCycle_Stop() {
	var lc life.Cycle
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
		fmt.Println(err)
	}
}

func ExampleCycle_Wait() {
	var lc life.Cycle
	lc.Start(func(stop <-chan struct{}) error {
		<-stop
		return nil
	})

	time.AfterFunc(100*time.Millisecond, lc.Stop)

	if err := lc.Wait(); err != nil {
		fmt.Println(err)
	}
}
