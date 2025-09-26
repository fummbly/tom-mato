package tui

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	defaultFPS = 60
	maxFPS     = 120
)

type renderer interface {
	Start()
	Stop()
	Write(string)
}

type standardRenderer struct {
	mtx *sync.Mutex
	out io.Writer

	buf                bytes.Buffer
	queuedMessageLines []string
	framerate          time.Duration
	ticker             *time.Ticker
	done               chan struct{}
	lastRender         string
	lastRenderedLines  []string
	linesRendered      int
	once               sync.Once
}

func NewRenderer(out io.Writer, fps int) renderer {

	if fps < 1 {
		fps = defaultFPS
	} else if fps > maxFPS {
		fps = maxFPS
	}

	r := &standardRenderer{
		out:                out,
		mtx:                &sync.Mutex{},
		done:               make(chan struct{}),
		framerate:          time.Second / time.Duration(fps),
		queuedMessageLines: []string{},
	}
	return r
}

func (r *standardRenderer) Start() {
	if r.ticker == nil {
		r.ticker = time.NewTicker(r.framerate)
	} else {
		r.ticker.Reset(r.framerate)
	}

	r.once = sync.Once{}

	go r.listen()
}

func (r *standardRenderer) Stop() {
	r.once.Do(func() {
		r.done <- struct{}{}
	})

	r.flush()

	r.mtx.Lock()
	defer r.mtx.Unlock()

}

func (r *standardRenderer) listen() {
	for {
		select {
		case <-r.done:
			r.ticker.Stop()
			return

		case <-r.ticker.C:
			r.flush()
		}
	}

}

func (r *standardRenderer) clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = r.out
	cmd.Run()
}

func (r *standardRenderer) clearLine() {
	r.out.Write([]byte("\033[2K\r"))
}

func (r *standardRenderer) upLine() {
	r.out.Write([]byte("\033[1A"))
}

func (r *standardRenderer) flush() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.buf.Len() == 0 || r.buf.String() == r.lastRender {
		return
	}

	buf := &bytes.Buffer{}

	newLines := strings.Split(r.buf.String(), "\n")

	flushQueuedMessages := len(r.queuedMessageLines) > 0

	if flushQueuedMessages {
		fmt.Println("Writing queued messages")
		for _, line := range r.queuedMessageLines {
			buf.WriteString(line)
			buf.WriteString("\r\n")
		}

		r.queuedMessageLines = []string{}
	}

	for i := range len(newLines) {

		if i == 0 && r.lastRender == "" {
			buf.WriteString("\r")
		}

		r.clearLine()

		line := newLines[i]

		buf.WriteString(line)

		if i < len(newLines)-1 {
			buf.WriteString("\r\n")
		}

	}

	r.out.Write(buf.Bytes())
	r.lastRender = r.buf.String()

	r.lastRenderedLines = newLines
	r.buf.Reset()

}

func (r *standardRenderer) Write(s string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.buf.Reset()

	if s == "" {
		s = " "
	}

	r.buf.WriteString(s)
}
