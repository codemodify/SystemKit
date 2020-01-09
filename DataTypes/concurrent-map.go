package DataTypes

import "sync"

// ConcurrentMap -
type ConcurrentMap struct {
	sync.RWMutex
	internal map[interface{}]bool
}

// NewConcurrentMap -
func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		internal: make(map[interface{}]bool),
	}
}

// Add -
func (thisRef *ConcurrentMap) Add(key interface{}) {
	thisRef.Lock()
	defer thisRef.Unlock()

	thisRef.internal[key] = true
}

// Remove -
func (thisRef *ConcurrentMap) Remove(key interface{}) {
	thisRef.Lock()
	defer thisRef.Unlock()

	delete(thisRef.internal, key)
}

// Contains -
func (thisRef *ConcurrentMap) Contains(key interface{}) bool {
	thisRef.RLock()
	defer thisRef.RUnlock()

	_, ok := thisRef.internal[key]
	return ok
}
