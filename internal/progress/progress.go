package progress

import (
	"fmt"
	"sync"
	"time"
)

// Spinner shows a "..." animation while a task is running
type Spinner struct {
	msg  string
	done chan struct{}
	wg   sync.WaitGroup
}

func New(msg string) *Spinner {
	s := &Spinner{
		msg:  msg,
		done: make(chan struct{}),
	}
	s.wg.Add(1)
	go s.run()
	return s
}

func (s *Spinner) run() {
	defer s.wg.Done()
	frames := []string{"   ", ".  ", ".. ", "..."}
	i := 0
	for {
		select {
		case <-s.done:
			// clear the spinner line before returning
			fmt.Printf("\r%-60s\r", "")
			return
		case <-time.After(300 * time.Millisecond):
			fmt.Printf("\r[*]%s%s", s.msg, frames[i%len(frames)])
			i++
		}
	}
}

func (s *Spinner) Stop(result string) {
	close(s.done)
	s.wg.Wait()
	fmt.Printf("[*]%s %s\n", s.msg, result)
}
