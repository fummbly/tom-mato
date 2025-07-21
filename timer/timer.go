package timer

import (
	"fmt"
	"time"
)

type Timer struct {
	stop    chan bool
	limit   int
	elapsed int
	paused  bool
	ticker  *time.Ticker
}

func NewTimer() *Timer {
	return &Timer{
		stop:   make(chan bool),
		paused: false,
	}
}

func NewLimitedTimer(limit int) *Timer {
	return &Timer{
		stop:   make(chan bool),
		limit:  limit,
		paused: false,
	}
}

func (t *Timer) Pause() {
	t.paused = true
}

func (t *Timer) Resume() {
	t.paused = false
}

func (t *Timer) Stop() {
	t.stop <- true
}

func (t *Timer) GetTickerChan() <-chan time.Time {
	return t.ticker.C
}

func (t *Timer) GetStopChan() <-chan bool {
	return t.stop
}

func (t *Timer) GetElapsedTime() int {
	return t.elapsed
}

func (t Timer) PrintTimeRemaining() {
	timeRemaining := t.limit - t.elapsed
	minutes := timeRemaining / 60
	seconds := timeRemaining - (minutes * 60)
	fmt.Printf("\r%02d:%02d", minutes, seconds)
}

func (t *Timer) Update() {
	defer close(t.stop)
	t.ticker = time.NewTicker(time.Second)
	for {
		if !t.paused {
			select {
			case <-t.ticker.C:
				t.elapsed++
				t.PrintTimeRemaining()
			case <-t.stop:
				t.ticker = nil
				return
			}
		}

		if t.elapsed >= t.limit && t.limit != 0 {
			t.ticker = nil
			return
		}
	}
}
