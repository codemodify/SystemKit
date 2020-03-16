package Helpers

import (
	"reflect"
	"sync"
)

var eventManagerWithDataOnceInstance *EventManagerWithDataOnce
var eventManagerWithDataOnceSync sync.Once

// EventManagerWithDataOnce - eventWithData manager
type EventManagerWithDataOnce struct {
	events     map[string]*eventWithData
	eventsSync sync.RWMutex
}

// EventsWithDataOnce - Singleton to get the events manager instance
func EventsWithDataOnce() *EventManagerWithDataOnce {
	eventManagerWithDataOnceSync.Do(func() {
		eventManagerWithDataOnceInstance = &EventManagerWithDataOnce{
			events:     make(map[string]*eventWithData),
			eventsSync: sync.RWMutex{},
		}
	})

	return eventManagerWithDataOnceInstance
}

// On - Tells EventManagerWithDataOnce to add a subscriber for an eventWithData
func (thisRef *EventManagerWithDataOnce) On(eventName string, eventHandler EventHandlerWithData) {
	thisRef.addEventIfNotExists(eventName)
	thisRef.addSubscriberIfNotExists(eventName, eventHandler)
}

// Off - Tells EventManagerWithDataOnce to remove a subscriber for an eventWithData
func (thisRef *EventManagerWithDataOnce) off(eventName string, eventHandler EventHandlerWithData) {
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

// Raise - Informs all subscribers about the eventWithData
func (thisRef *EventManagerWithDataOnce) Raise(eventName string, data []byte) {
	thisRef.addEventIfNotExists(eventName)

	thisRef.eventsSync.RLock()
	defer thisRef.eventsSync.RUnlock()

	thisRef.events[eventName].happenedOnce = true
	thisRef.events[eventName].lastData = data
	for _, eventHandler := range thisRef.events[eventName].subscribers {
		go eventHandler(data)
		go thisRef.off(eventName, eventHandler)
	}
}

func (thisRef *EventManagerWithDataOnce) addEventIfNotExists(eventName string) {
	thisRef.eventsSync.Lock()
	defer thisRef.eventsSync.Unlock()

	if _, ok := thisRef.events[eventName]; !ok {
		thisRef.events[eventName] = &eventWithData{
			subscribers:  []EventHandlerWithData{},
			happenedOnce: false,
		}
	}
}

func (thisRef *EventManagerWithDataOnce) addSubscriberIfNotExists(eventName string, eventHandler EventHandlerWithData) bool {
	thisRef.eventsSync.RLock()
	defer thisRef.eventsSync.RUnlock()

	// 1. Check if delegate for the eventWithData already there, assumes map-key exists
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
