package resource

import (
	"errors"
	"fmt"
	"sync"

	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type StaticReader[T storage.Supported] struct {
	current T
	name    string
	events  *event.Events
}

type StaticReadWriter[T storage.Supported] struct {
	*StaticReader[T]
	mutex sync.RWMutex
}

var ErrMissmatchedTypes = errors.New("missatched types")

func NewStaticReader[T storage.Supported](value T) *StaticReader[T] {
	return &StaticReader[T]{
		current: value,
		name:    "",
		events:  nil,
	}
}

func NewStaticReadWriter[T storage.Supported](value T) *StaticReadWriter[T] {
	return &StaticReadWriter[T]{
		StaticReader: NewStaticReader(value),
		mutex:        sync.RWMutex{},
	}
}

func (static *StaticReader[T]) Read() (any, error) {
	return static.current, nil
}

func (static *StaticReader[T]) Start(name string, storage *storage.Storage, events *event.Events) {
	static.name = name
	static.events = events
}

func (static *StaticReader[T]) Stop() {
	static.name = ""
	static.events = nil
}

func (static *StaticReadWriter[T]) Read() (any, error) {
	static.mutex.RLock()
	defer static.mutex.RUnlock()

	return static.StaticReader.Read()
}

func (static *StaticReadWriter[T]) Write(value any) error {
	static.mutex.Lock()
	defer static.mutex.Unlock()

	if val, ok := value.(T); ok {
		static.current = val
		static.events.Emit(static.name)

		return nil
	}

	return fmt.Errorf("%w: expected %T, go %T", ErrMissmatchedTypes, *new(T), value)
}
