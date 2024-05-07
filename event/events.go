package event

import (
	"slices"
	"sync"
	"time"
)

type Events struct {
	mutex     sync.RWMutex
	timeout   time.Duration
	listeners map[Event][]chan<- any
}

func NewEvents(timeout time.Duration) *Events {
	return &Events{
		mutex:     sync.RWMutex{},
		listeners: make(map[Event][]chan<- any),
	}
}

func (events *Events) Subscribe(event Event, listener chan<- any) {
	events.mutex.Lock()
	defer events.mutex.Unlock()

	events.listeners[event] = append(events.listeners[event], listener)
}

func (events *Events) Unsubscribe(event Event, listener chan<- any) {
	events.mutex.Lock()
	defer events.mutex.Unlock()

	events.listeners[event] = slices.DeleteFunc(events.listeners[event], func(lis chan<- any) bool {
		return listener == lis
	})
}

func (events *Events) Emit(event Event, payload any) {
	events.mutex.RLock()
	defer events.mutex.RUnlock()

	for _, listener := range events.listeners[event] {
		select {
		case listener <- payload:
			continue
		case <-time.After(events.timeout):
			continue
		}
	}
}
