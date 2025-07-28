package main

import (
	"fmt"
	"time"

	"github.com/fummbly/tom-mato/timer"
)

func main() {

	fmt.Println("App started")

	p := timer.NewPomodoro(1, 3, 30)

	go p.Update()

	time.Sleep(6 * time.Minute)

	p.Stop()

}
