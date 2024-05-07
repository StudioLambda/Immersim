package resource

import (
	"sync"
	"time"

	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type Increment[T storage.SupportedNumeric] struct {
	name     string
	storage  *storage.Storage
	events   *event.Events
	current  T
	mutex    sync.RWMutex
	quit     chan struct{}
	wg       sync.WaitGroup
	step     T
	interval time.Duration
}

func NewIncrement[T storage.SupportedNumeric](initial T, step T, interval time.Duration) *Increment[T] {
	return &Increment[T]{
		name:     "",
		storage:  nil,
		events:   nil,
		current:  initial,
		mutex:    sync.RWMutex{},
		quit:     nil,
		wg:       sync.WaitGroup{},
		step:     step,
		interval: interval,
	}
}

func (increment *Increment[T]) loop() {
	defer increment.wg.Done()

	ticker := time.NewTicker(increment.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			increment.mutex.Lock()
			increment.current += increment.step
			increment.mutex.Unlock()

			increment.events.Emit(increment.name)
		case <-increment.quit:
			return
		}
	}
}

func (increment *Increment[T]) Start(name string, storage *storage.Storage, events *event.Events) {
	increment.name = name
	increment.storage = storage
	increment.events = events
	increment.quit = make(chan struct{})

	increment.wg.Add(1)
	go increment.loop()
}

func (increment *Increment[T]) Stop() {
	close(increment.quit)

	increment.wg.Wait()

	increment.name = ""
	increment.storage = nil
	increment.events = nil
}

func (increment *Increment[T]) Read() (any, error) {
	increment.mutex.RLock()
	defer increment.mutex.RUnlock()

	return increment.current, nil
}
