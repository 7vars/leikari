package leikari

import (
	"sync"
	"time"
)

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return Errorln("", "timeout reached")
	}
}