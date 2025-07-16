package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fummbly/tom-mato/timer"
)

func main() {

	fmt.Println("App started")
	fmt.Print("Enter Command:")

	myTimer := timer.NewPomodoro()

	myTimer.Start()

	var wg sync.WaitGroup
	wg.Add(2)

	live := true

	go func() {
		defer wg.Done()
		updateChannel := myTimer.GetUpdateChannel()
		for elapsedTime := range updateChannel {
			if live {
				fmt.Printf("[Live] time: %d\n", elapsedTime)
			}
		}
		fmt.Println("[Live] update channel closed")
	}()

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(os.Stdin)

		userInputCh := make(chan string)

		go func() {
			for scanner.Scan() {
				userInputCh <- scanner.Text()
			}
			close(userInputCh)
		}()

		timerDoneCh := myTimer.GetDoneChannel()

		for {

			for {
				select {
				case input, ok := <-userInputCh:
					if !ok {
						fmt.Println("[Input handler] user input channel closed")
						return
					}
					cmd := strings.TrimSpace(strings.ToLower(input))
					switch cmd {
					case "p":
						myTimer.Pause()
					case "r":
						myTimer.Resume()
					case "l":
						live = !live
					case "s":
						myTimer.Stop()
						return
					default:
						fmt.Println("Invalid command")
					}
					if input != "s" {
						fmt.Print("Enter Command:")
					}
				case <-timerDoneCh:
					fmt.Println("[Input handler] timer finished exiting")
					return
				}
			}

		}
		if err := scanner.Err(); err != nil {
			fmt.Println(os.Stderr, "reading standard input", err)
		}
	}()

	wg.Wait()

}
