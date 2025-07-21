package main

import (
	"fmt"
	"time"

	"github.com/fummbly/tom-mato/timer"
)

func main() {

	fmt.Println("App started")

	p := timer.NewPomodoro(20, 5, 30)

	go p.Update()

	time.Sleep(6 * time.Second)

	p.Pause()

	time.Sleep(4 * time.Second)

	fmt.Println("\nApp Stopped")

}
