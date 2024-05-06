package resource

import (
	"fmt"
	"sync"

	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type Action struct {
	name     string
	storage  *storage.Storage
	events   *event.Events
	current  bool
	mutex    sync.RWMutex
	callback func(storage *storage.Storage, events *event.Events) bool
}

func NewAction(callback func(storage *storage.Storage, events *event.Events) bool) *Action {
	return &Action{
		name:     "",
		storage:  nil,
		events:   nil,
		current:  false,
		mutex:    sync.RWMutex{},
		callback: callback,
	}
}

func (action *Action) Start(name string, storage *storage.Storage, events *event.Events) {
	action.current = false
	action.name = name
	action.storage = storage
	action.events = events
}

func (action *Action) Stop() {
	action.name = ""
	action.storage = nil
	action.events = nil
}

func (action *Action) Read() (any, error) {
	action.mutex.RLock()
	defer action.mutex.RUnlock()

	return action.current, nil
}

func (action *Action) Write(value any) error {
	if val, ok := value.(bool); ok {
		action.mutex.Lock()
		action.current = val
		action.mutex.Unlock()

		if action.current {
			if action.callback(action.storage, action.events) {
				action.mutex.Lock()
				action.current = false
				action.mutex.Unlock()
			}
		}

		return nil
	}

	return fmt.Errorf("%w: expected %T, go %T", ErrMissmatchedTypes, action.current, value)
}
