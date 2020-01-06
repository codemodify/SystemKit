package Helpers

// RingBuffer -
type RingBuffer struct {
	In   chan []byte
	Data chan []byte
}

// NewRingBuffer -
func NewRingBuffer(in chan []byte, out chan []byte) *RingBuffer {
	var ringBuffer = &RingBuffer{
		In:   in,
		Data: out,
	}

	go ringBuffer.run()

	return ringBuffer
}

func (thisRef RingBuffer) run() {
	for v := range thisRef.In {
		select {
		case thisRef.Data <- v:
		default:
			<-thisRef.Data
			thisRef.Data <- v
		}
	}
}
