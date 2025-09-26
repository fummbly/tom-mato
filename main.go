package main

import (
	"fmt"
	"os"
	"time"

	//"github.com/fummbly/tom-mato/timer"
	"github.com/fummbly/tom-mato/tui"
)

func main() {

	fmt.Println("App started")

	r := tui.NewRenderer(os.Stdout, 60)

	r.Start()

	for i := range 10 {
		r.Write(fmt.Sprintf("Rendering: %d\nTesting multi line\nRenderings", i))
		time.Sleep(time.Second)
	}

	//p := timer.NewPomodoro(1, 3, 30)

	//go p.Update()

	//time.Sleep(6 * time.Minute)

	//p.Stop()

}
