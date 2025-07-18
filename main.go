package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fummbly/tom-mato/timer"
)

func main() {

	fmt.Println("App started")

	t := timer.Timer{
		Limit:  5,
		Paused: false,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		t.Update()
	}()

	time.Sleep(2 * time.Second)
	t.Stop()
	//time.Sleep(3 * time.Second)
	//t.Paused = false

	wg.Wait()
}
