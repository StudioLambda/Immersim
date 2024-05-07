package resource

import (
	"errors"
	"fmt"
	"sync"

	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type Static[T storage.Supported] struct {
	*Constant[T]
	mutex sync.RWMutex
}

var ErrMissmatchedTypes = errors.New("missatched types")

func NewStatic[T storage.Supported](value T) *Static[T] {
	return &Static[T]{
		Constant: NewConstant(value),
		mutex:    sync.RWMutex{},
	}
}

func (static *Static[T]) Read() (any, error) {
	static.mutex.RLock()
	defer static.mutex.RUnlock()

	return static.Constant.Read()
}

func (static *Static[T]) Write(value any) error {
	static.mutex.Lock()
	defer static.mutex.Unlock()

	if val, ok := value.(T); ok {
		static.current = val
		static.events.Emit(event.Changed(static.name), nil)

		return nil
	}

	return fmt.Errorf("%w: expected %T, go %T", ErrMissmatchedTypes, *new(T), value)
}
