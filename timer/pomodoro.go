package timer

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
	work      *Timer
	shortRest *Timer
	longRest  *Timer
	state     PomodoroState
	cycles    int
}
