package timer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartStopTimer(t *testing.T) {
	cases := []int{
		2, 4, 6, 3, 5, 7,
	}

	for _, expected := range cases {
		t.Run(fmt.Sprintf("Timer: %d seconds", expected), func(t *testing.T) {
			t.Parallel()
			var timer = New()
			timer.Start()
			waitTime := expected
			time.Sleep(time.Duration(waitTime) * time.Second)
			timer.Stop()
			got := timer.GetElapsedTime()
			assert.Equal(t, expected, got)
		})
	}
}

func TestPauseResumeTimer(t *testing.T) {
	var timer = New()
	t.Run("Pause Resume", func(t *testing.T) {
		timer.Start()
		time.Sleep(2 * time.Second)
		timer.Pause()
		time.Sleep(4 * time.Second)
		timer.Resume()
		time.Sleep(6 * time.Second)
		timer.Stop()
		ans := timer.GetElapsedTime()
		expected := 8
		assert.Equal(t, expected, ans)
	})
}
