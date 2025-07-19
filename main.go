package main

import (
	"fmt"

	"github.com/fummbly/tom-mato/timer"
)

func main() {

	fmt.Println("App started")

	p := timer.NewPomodoro()

	p.Update()

}
