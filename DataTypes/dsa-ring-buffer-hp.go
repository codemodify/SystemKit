package DataTypes

import "sync/atomic"
import "runtime"

// RingBufferHP -
type RingBufferHP struct {
	// The padding members 1 to 5 below are here to ensure each item is on a separate cache line.
	// This prevents false sharing and hence improves performance.
	padding1           [8]uint64
	lastCommittedIndex uint64
	padding2           [8]uint64
	nextFreeIndex      uint64
	padding3           [8]uint64
	readerIndex        uint64
	padding4           [8]uint64
	contents           [][]byte
	padding5           [8]uint64
	bufferSize         uint64
	indexMask          uint64 // Masking is faster than division
}

// NewRingBufferHP -
func NewRingBufferHP(bufferSize uint64) *RingBufferHP {
	return &RingBufferHP{
		lastCommittedIndex: 0,
		nextFreeIndex:      1,
		readerIndex:        1,
		contents:           make([][]byte, bufferSize),
		bufferSize:         bufferSize,
		indexMask:          bufferSize - 1,
	}
}

// Write -
func (thisRef *RingBufferHP) Write(value []byte) {
	var myIndex = atomic.AddUint64(&thisRef.nextFreeIndex, 1) - 1
	//Wait for reader to catch up, so we don't clobber a slot which it is (or will be) reading
	for myIndex > (thisRef.readerIndex + thisRef.bufferSize - 2) {
		runtime.Gosched()
	}
	//Write the item into it's slot
	thisRef.contents[myIndex&thisRef.indexMask] = value
	//Increment the lastCommittedIndex so the item is available for reading
	for !atomic.CompareAndSwapUint64(&thisRef.lastCommittedIndex, myIndex-1, myIndex) {
		runtime.Gosched()
	}
}

// Read -
func (thisRef *RingBufferHP) Read() []byte {
	var myIndex = atomic.AddUint64(&thisRef.readerIndex, 1) - 1
	//If reader has out-run writer, wait for a value to be committed
	for myIndex > thisRef.lastCommittedIndex {
		runtime.Gosched()
	}
	return thisRef.contents[myIndex&thisRef.indexMask]
}

// func (thisRef *RingBufferHP) Dump() {
// 	fmt.Printf("lastCommitted: %3d, nextFree: %3d, readerIndex: %3d, content:", thisRef.lastCommittedIndex, self.nextFreeIndex, self.readerIndex)
// 	for index, value := range thisRef.contents {
// 		fmt.Printf("%5v : %5v", index, value)
// 	}
// 	fmt.Print("\n")
// }
