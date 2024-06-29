package main

import (
	"bytes"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"sync"
	"time"
)

type stopWatch struct {
	isRunning bool
	isExit    bool
	drawBuf   *bytes.Buffer
	time      int
	wg        *sync.WaitGroup
	mu        sync.Mutex
}

func newStopwatch() *stopWatch {
	return &stopWatch{
		isExit:  false,
		drawBuf: new(bytes.Buffer),
		wg:      &sync.WaitGroup{},
	}
}
func (w *stopWatch) start() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isRunning {
		return
	}

	w.isRunning = true
	w.wg.Add(1)
	go w.loop()
}

func (w *stopWatch) stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.isRunning {
		return
	}

	w.isRunning = false
	w.wg.Wait()
}
func (w *stopWatch) loop() {
	defer w.wg.Done()

	for w.isRunning {
		w.render()
		w.timeCounter()
		time.Sleep(time.Second * 1)
	}
}

func (w *stopWatch) render() {
	w.drawBuf.Reset()
	fmt.Fprintf(os.Stdout, "\033[2J\033[1;1H")
	w.renderTime()
	fmt.Fprint(os.Stdout, w.drawBuf.String())

}

func (w *stopWatch) renderTime() {
	h, m, s := w.formatTime()

	hour_left := h / 10
	hour_right := h % 10
	min_left := m / 10
	min_right := m % 10
	sec_left := s / 10
	sec_right := s % 10

	for i := 0; i < len(Hour); i++ {
		w.drawBuf.WriteString(fmt.Sprintf("%s%s  %s　 %s%s  %s   %s%s  %s",
			AsciiNumMap[hour_left][i][0], AsciiNumMap[hour_right][i][0], Hour[i][0],
			AsciiNumMap[min_left][i][0], AsciiNumMap[min_right][i][0], Min[i][0],
			AsciiNumMap[sec_left][i][0], AsciiNumMap[sec_right][i][0], Sec[i][0]))
		w.drawBuf.WriteString("\n")
	}
}

func (w *stopWatch) formatTime() (uint8, uint8, uint8) {
	hour := uint8(w.time / 3600)
	minute := uint8((w.time % 3600) / 60)
	sec := uint8(w.time % 60)
	return hour, minute, sec
}

func (w *stopWatch) controller() {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Failed to initialize termbox")
		return
	}

	defer termbox.Close()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 'r' { // restart
				w.stop()
				w.time = 0
				w.start()
			} else if ev.Ch == 'p' { // pause
				if w.isRunning {
					w.stop()
				} else if !(w.isRunning) {
					w.start()
				}
			} else if ev.Ch == 'q' { // quit
				w.stop()
				w.isExit = true
				return
			}

		case termbox.EventError:
			fmt.Println("Error:", ev.Err)
			return
		}
	}
}

func (w *stopWatch) timeCounter() {
	w.time++
}

func main() {
	s := newStopwatch()
	s.start()
	go s.controller()
	for !(s.isExit) {
	}
}
