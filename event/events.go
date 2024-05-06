package event

import (
	"slices"
	"sync"
)

type Events struct {
	mutex     sync.RWMutex
	listeners map[string][]chan<- struct{}
}

func NewEvents() *Events {
	return &Events{
		mutex:     sync.RWMutex{},
		listeners: make(map[string][]chan<- struct{}),
	}
}

func (events *Events) Subscribe(event string, listener chan<- struct{}) {
	events.mutex.Lock()
	defer events.mutex.Unlock()

	events.listeners[event] = append(events.listeners[event], listener)
}

func (events *Events) Unsubscribe(event string, listener chan<- struct{}) {
	events.mutex.Lock()
	defer events.mutex.Unlock()

	events.listeners[event] = slices.DeleteFunc(events.listeners[event], func(lis chan<- struct{}) bool {
		return listener == lis
	})
}

func (events *Events) Emit(event string) {
	events.mutex.RLock()
	defer events.mutex.RUnlock()

	for _, listener := range events.listeners[event] {
		listener <- struct{}{}
	}
}
