package timer

import (
	"fmt"
	"time"
)

type Timer struct {
	stop    chan bool
	Limit   int
	Elapsed int
	Paused  bool
}

func (t *Timer) GetElapsedTime() int {
	return t.Elapsed
}

func (t *Timer) Pause() {
	t.Paused = true
}

func (t *Timer) Resume() {
	t.Paused = false
}

func (t *Timer) Stop() {
	t.stop <- true
}

func (t *Timer) Update() {
	t.stop = make(chan bool)
	ticker := time.NewTicker(time.Second)
	for {
		if !t.Paused {
			select {
			case <-ticker.C:
				t.Elapsed++
				fmt.Printf("Time since: %d\n", t.Elapsed)
			case <-t.stop:
				ticker = nil
				return
			}
		}

		if t.Elapsed >= t.Limit {
			ticker = nil
			return
		}
	}
}
