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
			helpersEvents.EventsWithData().Raise("PING", nil)
		}
	}()

	// 1 ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
	wg1 := sync.WaitGroup{}

	pongHandler1 := func(data []byte) {
		wg1.Done()
	}

	wg1.Add(1)
	helpersEvents.EventsWithData().On("PING", pongHandler1)

	wg1.Add(1)
	helpersEvents.EventsWithData().On("PING", pongHandler1)

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
