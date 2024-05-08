package resource

import (
	"sync"

	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type Computed[T storage.Supported] struct {
	current      T
	name         string
	storage      *storage.Storage
	events       *event.Events
	callback     func(name string, storage *storage.Storage) T
	dependencies []string
	mutex        sync.RWMutex
	listener     chan any
	waitGroup    sync.WaitGroup
}

func NewComputed[T storage.Supported](callback func(name string, storage *storage.Storage) T, dependencies []string) *Computed[T] {
	return &Computed[T]{
		name:         "",
		storage:      nil,
		events:       nil,
		callback:     callback,
		dependencies: dependencies,
		listener:     nil,
		current:      *new(T),
		mutex:        sync.RWMutex{},
		waitGroup:    sync.WaitGroup{},
	}
}

func (computed *Computed[T]) Start(name string, storage *storage.Storage, events *event.Events) {
	computed.name = name
	computed.storage = storage
	computed.events = events
	computed.listener = make(chan any, len(computed.dependencies))
	computed.current = computed.callback(computed.name, computed.storage)

	computed.waitGroup.Add(1)
	go computed.loop()

	for _, dependency := range computed.dependencies {
		computed.events.Subscribe(event.Changed(dependency), computed.listener)
	}
}

func (computed *Computed[T]) Stop() {
	for _, dependency := range computed.dependencies {
		computed.events.Unsubscribe(event.Changed(dependency), computed.listener)
	}

	close(computed.listener)
	computed.waitGroup.Wait()

	computed.name = ""
	computed.storage = nil
	computed.events = nil
	computed.listener = nil
	computed.current = *new(T)
}

func (computed *Computed[T]) Read() (any, error) {
	computed.mutex.RLock()
	defer computed.mutex.RUnlock()

	return computed.current, nil
}

func (computed *Computed[T]) loop() {
	defer computed.waitGroup.Done()

	for range computed.listener {
		computed.mutex.Lock()
		computed.current = computed.callback(computed.name, computed.storage)
		computed.events.Emit(event.Changed(computed.name), event.ChangedPayload{
			Resource: computed.name,
			Value:    computed.current,
		})
		computed.mutex.Unlock()
	}
}
