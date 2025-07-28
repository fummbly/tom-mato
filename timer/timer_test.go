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
		t.Run(fmt.Sprintf("Timer:%d", expected), func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			var timer = NewTimer()
			go timer.Update()
			waitTime := expected
			time.Sleep(time.Duration(waitTime) * time.Second)
			timer.Stop()
			got := timer.GetElapsedTime()
			elapsed := int(time.Since(start).Seconds())
			assert.Equal(t, elapsed, got)
		})
	}
}

func TestPauseResumeTimer(t *testing.T) {
	cases := []struct {
		startWait, pauseWait, resumeWait, expected int
	}{
		{2, 3, 4, 7},
		{4, 1, 6, 11},
		{6, 8, 2, 9},
		{5, 9, 2, 8},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Run:%d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()
			var timer = NewTimer()
			go timer.Update()
			time.Sleep(time.Duration(tc.startWait) * time.Second)
			timer.Pause()
			time.Sleep(time.Duration(tc.pauseWait) * time.Second)
			timer.Resume()
			time.Sleep(time.Duration(tc.resumeWait) * time.Second)
			timer.Stop()
			got := timer.GetElapsedTime()
			assert.Equal(t, tc.expected, got)
		})
	}
}
