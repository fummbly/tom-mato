package main

import (
	"fmt"
	"time"

	"github.com/fummbly/tom-mato/timer"
)

func main() {

	fmt.Println("App started")

	p := timer.NewPomodoro()

	go p.Update()

	time.Sleep(6 * time.Second)

	p.Pause()

	time.Sleep(4 * time.Second)

	p.Resume()

	time.Sleep(10 * time.Second)

	p.Stop()

}
