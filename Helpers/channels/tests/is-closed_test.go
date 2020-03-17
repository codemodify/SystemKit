package tests

import (
	"fmt"
	"testing"

	helpersChannels "github.com/codemodify/SystemKit/Helpers/channels"
)

func Test_is_closed(t *testing.T) {

	// 1
	c := make(chan []byte)

	if helpersChannels.IsClosed(c) {
		fmt.Println("A-CLOSED-1")
	} else {
		fmt.Println("A-CLOSED-NOT-1")
	}

	if helpersChannels.IsClosed(c) {
		fmt.Println("A-CLOSED-2")
	} else {
		fmt.Println("A-CLOSED-NOT-2")
	}

	close(c)

	if helpersChannels.IsClosed(c) {
		fmt.Println("A-CLOSED-3")
	} else {
		fmt.Println("A-CLOSED-NOT-3")
	}

	// 2
	c = make(chan []byte)

	if helpersChannels.IsClosed(c) {
		fmt.Println("B-CLOSED-1")
	} else {
		fmt.Println("B-CLOSED-NOT-1")
	}

	if helpersChannels.IsClosed(c) {
		fmt.Println("B-CLOSED-2")
	} else {
		fmt.Println("B-CLOSED-NOT-2")
	}

	close(c)

	if helpersChannels.IsClosed(c) {
		fmt.Println("B-CLOSED-3")
	} else {
		fmt.Println("B-CLOSED-NOT-3")
	}
}
