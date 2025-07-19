package timer

import "sync"

type pomoState int

const (
	Work pomoState = iota
	ShortRest
	LongRest
)

type pomoTimers struct {
	work      *Timer
	shortRest *Timer
	longRest  *Timer
}

type Pomodoro struct {
	timers *pomoTimers
	stop   chan bool
	state  pomoState
	cycles int
	total  int
}

func NewPomodoro() *Pomodoro {
	return &Pomodoro{
		timers: &pomoTimers{
			work:      NewLimitedTimer(20),
			shortRest: NewLimitedTimer(5),
			longRest:  NewLimitedTimer(30),
		},
		stop:  make(chan bool),
		state: Work,
	}
}

func (p *Pomodoro) Update() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			switch p.state {
			case Work:
				p.timers.work.Update()
				p.total += p.timers.work.GetElapsedTime()
				p.timers.work = NewLimitedTimer(20)
				p.state = ShortRest
				continue
			case ShortRest:
				p.timers.shortRest.Update()
				p.timers.shortRest = NewLimitedTimer(5)
				if p.cycles == 3 {
					p.state = LongRest
					continue
				}
				p.cycles++
				p.state = Work
				continue
			case LongRest:
				p.timers.longRest.Update()
				return
			}

		}

	}()

	wg.Wait()
}
