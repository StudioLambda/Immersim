package resource

import (
	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type Constant[T storage.Supported] struct {
	current T
	name    string
	events  *event.Events
}

func NewConstant[T storage.Supported](value T) *Constant[T] {
	return &Constant[T]{
		current: value,
		name:    "",
		events:  nil,
	}
}

func (static *Constant[T]) Read() (any, error) {
	return static.current, nil
}

func (static *Constant[T]) Start(name string, storage *storage.Storage, events *event.Events) {
	static.name = name
	static.events = events
}

func (static *Constant[T]) Stop() {
	static.name = ""
	static.events = nil
}
