package Helpers

import (
	"reflect"
	"sync"
)

// TODO: Add locking at some point for Add/Remove operations

var eventManagerInstance *EventManager
var eventManagerOnce sync.Once

// EventHandler - event handler prototype
type EventHandler func()

type event struct {
	subscribers  []EventHandler
	happenedOnce bool
}

// EventManager - event manager
type EventManager struct {
	events     map[string]*event
	eventsSync sync.RWMutex
}

// Events - Singleton to get the events manager instance
func Events() *EventManager {
	eventManagerOnce.Do(func() {
		eventManagerInstance = &EventManager{
			events:     make(map[string]*event),
			eventsSync: sync.RWMutex{},
		}
	})

	return eventManagerInstance
}

// On - Tells EventManager to add a subscriber for an event
func (thisRef *EventManager) On(eventName string, eventHandler EventHandler) {
	thisRef.addEventIfNotExists(eventName)
	thisRef.addSubscriberIfNotExists(eventName, eventHandler)

	thisRef.eventsSync.RLock()
	defer thisRef.eventsSync.RUnlock()

	if thisRef.events[eventName].happenedOnce {
		go eventHandler()
	}
}

// Off - Tells EventManager to remove a subscriber for an event
func (thisRef *EventManager) Off(eventName string, eventHandler EventHandler) {
	thisRef.addEventIfNotExists(eventName)

	thisRef.eventsSync.RLock()
	defer thisRef.eventsSync.RUnlock()

	var foundIndex = -1
	for index, existingEventHandler := range thisRef.events[eventName].subscribers {
		if reflect.ValueOf(eventHandler) == reflect.ValueOf(existingEventHandler) {
			foundIndex = index
			break
		}
	}

	if foundIndex != -1 {
		thisRef.events[eventName].subscribers = append(
			thisRef.events[eventName].subscribers[:foundIndex],
			thisRef.events[eventName].subscribers[foundIndex+1:]...,
		)
	}
}

// Raise - Informs all subscribers about the event
func (thisRef *EventManager) Raise(eventName string) {
	thisRef.addEventIfNotExists(eventName)

	thisRef.eventsSync.RLock()
	defer thisRef.eventsSync.RUnlock()

	thisRef.events[eventName].happenedOnce = true
	for _, eventHandler := range thisRef.events[eventName].subscribers {
		go eventHandler()
	}
}

func (thisRef *EventManager) addEventIfNotExists(eventName string) {
	thisRef.eventsSync.Lock()
	defer thisRef.eventsSync.Unlock()

	if _, ok := thisRef.events[eventName]; !ok {
		thisRef.events[eventName] = &event{
			subscribers:  []EventHandler{},
			happenedOnce: false,
		}
	}
}

func (thisRef *EventManager) addSubscriberIfNotExists(eventName string, eventHandler EventHandler) bool {
	thisRef.eventsSync.RLock()
	defer thisRef.eventsSync.RUnlock()

	// 1. Check if delegate for the event already there, assumes map-key exists
	var alreadyThere = false

	for _, existingEventHandler := range thisRef.events[eventName].subscribers {
		alreadyThere = (reflect.ValueOf(eventHandler) == reflect.ValueOf(existingEventHandler))
		if alreadyThere {
			break
		}
	}

	// 2. Add the delegate
	if !alreadyThere {
		thisRef.events[eventName].subscribers = append(
			thisRef.events[eventName].subscribers,
			eventHandler,
		)
	}

	return alreadyThere
}
