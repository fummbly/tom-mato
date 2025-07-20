package timer

import "sync"

type pomoState int

const (
	Work pomoState = iota
	ShortRest
	LongRest
)

type Pomodoro struct {
	timers map[pomoState]*Timer
	stop   chan bool
	state  pomoState
	cycles int
	total  int
}

func NewPomodoro() *Pomodoro {
	return &Pomodoro{
		timers: map[pomoState]*Timer{
			Work:      NewLimitedTimer(20),
			ShortRest: NewLimitedTimer(5),
			LongRest:  NewLimitedTimer(30),
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
	p.stop <- true

}

func (p *Pomodoro) Update() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-p.stop:
				p.timers[p.state].Stop()
				return
			default:
			}
			p.timers[p.state].Update()
			switch p.state {
			case Work:
				p.total += p.timers[p.state].GetElapsedTime()
				p.timers[p.state] = NewLimitedTimer(20)
				p.state = ShortRest
				continue
			case ShortRest:
				p.timers[p.state] = NewLimitedTimer(5)
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

		}

	}()

	wg.Wait()
}
