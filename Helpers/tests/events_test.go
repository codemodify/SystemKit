package tests

import (
	"sync"
	"testing"
	"time"

	helpersEvents "github.com/codemodify/SystemKit/Helpers"
)

func Test_Events_01(t *testing.T) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			helpersEvents.EventsWithDataOnce().Raise("PING", nil)
		}
	}()

	// 1 ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
	wg1 := sync.WaitGroup{}

	wg1.Add(1)
	pongHandler1 := func(data []byte) {
		wg1.Done()
	}
	helpersEvents.EventsWithDataOnce().On("PING", pongHandler1)
	wg1.Wait()

	// // 2 ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
	// wg2 := sync.WaitGroup{}

	// wg2.Add(1)
	// pongHandler2 := func(data []byte) {
	// 	wg2.Done()
	// }
	// helpersEvents.EventsWithData().On("PONG", pongHandler2)
	// wg2.Wait()
	// helpersEvents.EventsWithData().Off("PONG", pongHandler2)

	// 3 ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
	for {
		time.Sleep(5 * time.Second)
	}
}
