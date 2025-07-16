package timer

import (
	"fmt"
	"time"
)

type Timer struct {
	ticker   *time.Ticker // ticker for ticking seconds
	seconds  int          // seconds variable to keep track of time
	limit    int          // optional limit to stop timer
	stopCh   chan bool    // stop channel to signal timer to stop
	done     chan bool    // done channel to signal that timer is done
	updateCh chan int     // live update channel
	active   bool         // active bool for pause and resume
}

// Basic timer constructor
func New() *Timer {
	return &Timer{
		stopCh:   make(chan bool),
		done:     make(chan bool),
		updateCh: make(chan int, 1),
		active:   false,
	}
}

// Limited timer constructor
func NewLimited(limit int) *Timer {
	timer := New()
	timer.limit = limit
	return timer
}

// function for starting timer
func (t *Timer) Start() {
	// checks if timer is already started
	if t.ticker != nil {
		fmt.Println("Timer already started")
		return
	}

	// initiates timer values for starting
	t.ticker = time.NewTicker(time.Second)
	t.seconds = 0
	t.active = true

	fmt.Println("Timer Started")

	// goroutine for live ticking
	go func() {
		// close done and updateCh once function returns
		defer close(t.done)
		defer close(t.updateCh)

		for {
			if t.limit != 0 && t.seconds >= t.limit {
				t.ticker = nil
				return
			}
			// if the seconds reach limit stop the timer
			// select between stopCh and ticker
			select {
			case <-t.stopCh:
				return

			case <-t.ticker.C:

				// if the timer is active increment seconds
				// and send seconds to updateCh
				if t.active {
					t.seconds++
					select {
					case t.updateCh <- t.seconds:

					default:

					}
				}
			}
		}

	}()
}

// stop timer function
func (t *Timer) Stop() {
	// checks if timer is already stopped
	if t.ticker == nil {
		fmt.Println("Timer already stopped")
		return
	}
	fmt.Println("Timer Stopped")

	fmt.Println("Total time elapsed:", t.seconds, "seconds")
	// send stop and done and nil ticker
	t.stopCh <- true
	<-t.done
	t.ticker = nil

}

// Pause function
func (t *Timer) Pause() {
	if !t.active {
		fmt.Println("Timer already paused")
		return
	}

	t.active = false
}

// Resume function
func (t *Timer) Resume() {
	if t.active {
		fmt.Println("Timer already resumed")
		return
	}

	t.active = true
}

// Getter for elapsed time
func (t *Timer) GetElapsedTime() int {
	return t.seconds
}

// Getter for update channel
func (t *Timer) GetUpdateChannel() <-chan int {
	return t.updateCh
}

// Getter for done channel
func (t *Timer) GetDoneChannel() <-chan bool {
	return t.done
}
