package timer

import "fmt"

type PomodoroState int

const (
	TimerWork PomodoroState = iota
	TimerShortRest
	TimerLongRest
)

var stateName = map[PomodoroState]string{
	TimerWork:      "work",
	TimerShortRest: "short rest",
	TimerLongRest:  "long rest",
}

func (ps PomodoroState) String() string {
	return stateName[ps]
}

type Pomodoro struct {
	timers    map[PomodoroState]*Timer
	stopCh    chan bool
	updateCh  chan int
	done      chan bool
	state     PomodoroState
	totalTime int
	cycles    int
	active    bool
}

func NewPomodoro() *Pomodoro {
	return &Pomodoro{
		timers: map[PomodoroState]*Timer{
			TimerWork:      NewLimited(30),
			TimerShortRest: NewLimited(5),
			TimerLongRest:  NewLimited(20),
		},
		stopCh:   make(chan bool),
		updateCh: make(chan int, 1),
		done:     make(chan bool),
		state:    TimerWork,
		cycles:   0,
	}
}

func (p *Pomodoro) Start() {
	if p.state != TimerWork && p.cycles != 0 {
		fmt.Println("Pomodoro already started")
		return
	}

	p.timers[TimerWork].Start()
	p.active = true

	fmt.Println("Podooro started")

	go func() {
		defer close(p.done)
		defer close(p.updateCh)

		for {
			currTimer := p.timers[p.state]
			currUpdateCh := currTimer.GetUpdateChannel()
			currDoneCh := currTimer.GetDoneChannel()
			select {
			case <-p.stopCh:
				currTimer.Stop()
				return

			case <-currDoneCh:
				switch p.state {
				case TimerShortRest:
					if p.cycles == 3 {
						p.state = TimerLongRest
						p.timers[TimerLongRest].Start()
					} else {
						p.state = TimerWork
						p.timers[TimerWork].Start()
						p.cycles++
					}
					continue
				case TimerLongRest:
					return

				case TimerWork:
					p.state = TimerShortRest
					p.timers[TimerShortRest].Start()
					continue

				default:
					fmt.Println("Invalid state")
					return

				}
			case <-currUpdateCh:
				if p.active {
					p.totalTime++
					select {
					case p.updateCh <- currTimer.seconds:

					default:

					}
				}
			}

		}
	}()

}

func (p *Pomodoro) Stop() {
	if p.timers[p.state].ticker == nil {
		fmt.Println("Pomodoro already stopped")
		return
	}

	p.stopCh <- true
	<-p.done
}

func (p *Pomodoro) Pause() {
	if !p.active {
		fmt.Println("Pomodoro already paused")
		return
	}

	p.active = false
	p.timers[p.state].Pause()
}

func (p *Pomodoro) Resume() {
	if p.active {
		fmt.Println("Pomodoro already playing")
		return
	}

	p.active = true
	p.timers[p.state].Resume()
}

func (p *Pomodoro) GetElapsedTime() int {
	return p.totalTime
}

func (p *Pomodoro) GetUpdateChannel() <-chan int {
	return p.updateCh
}

func (p *Pomodoro) GetDoneChannel() <-chan bool {
	return p.done
}
