package goropool_test

import (
	"fmt"

	"github.com/CAFxX/goropool"
)

func Example() {
	queue, done := goropool.NewDefaultPool()
	for i := 0; i < 10; i++ {
		queue <- func() {
			fmt.Print("*")
		}
	}
	close(queue)
	<-done
	// Output: **********
}
