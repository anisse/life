package life_test

import (
	"errors"
	"fmt"

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
