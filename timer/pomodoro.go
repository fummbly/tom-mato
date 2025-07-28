package timer

import (
	"fmt"
	"sync"
)

type pomoState int

const (
	Work pomoState = iota
	ShortRest
	LongRest
)

var stateNames = map[pomoState]string{
	Work:      "work",
	ShortRest: "short rest",
	LongRest:  "long rest",
}

type Pomodoro struct {
	timers    map[pomoState]*Timer
	durations [3]int
	stop      chan bool
	state     pomoState
	cycles    int
	total     int
}

func NewPomodoro(workDuration, shortDuration, longDuration int) *Pomodoro {
	return &Pomodoro{
		durations: [3]int{workDuration, shortDuration, longDuration},
		timers: map[pomoState]*Timer{
			Work:      NewLimitedTimer(workDuration * 60),
			ShortRest: NewLimitedTimer(shortDuration * 60),
			LongRest:  NewLimitedTimer(longDuration * 60),
		},
		stop:  make(chan bool),
		state: Work,
	}
}

func (p *Pomodoro) Pause() {
	p.timers[p.state].Pause()
}

func (p *Pomodoro) Resume() {
	p.timers[p.state].Resume()
}

func (p *Pomodoro) Stop() {
	p.timers[p.state].Stop()

}

func (p *Pomodoro) PrintTime() {

}

func (p *Pomodoro) Update() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			fmt.Printf("\r%v\n", stateNames[p.state])
			p.timers[p.state].Update()
			switch p.state {
			case Work:
				p.total += p.timers[p.state].GetElapsedTime()
				p.timers[p.state] = NewLimitedTimer(p.durations[0] * 60)
				p.state = ShortRest
				continue
			case ShortRest:
				p.timers[p.state] = NewLimitedTimer(p.durations[1] * 60)
				if p.cycles == 3 {
					p.state = LongRest
					continue
				}
				p.cycles++
				p.state = Work
				continue
			case LongRest:
				return
			}
			return

		}

	}()

	wg.Wait()
}
