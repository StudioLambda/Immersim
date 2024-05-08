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
	initial  T
	mutex    sync.RWMutex
	quit     chan struct{}
	wg       sync.WaitGroup
	step     T
	interval time.Duration
	reset    chan any
	pause    chan any
	resume   chan any
}

func NewIncrement[T storage.SupportedNumeric](initial T, step T, interval time.Duration) *Increment[T] {
	return &Increment[T]{
		name:     "",
		storage:  nil,
		events:   nil,
		current:  initial,
		initial:  initial,
		mutex:    sync.RWMutex{},
		quit:     nil,
		wg:       sync.WaitGroup{},
		step:     step,
		interval: interval,
		reset:    nil,
		pause:    nil,
		resume:   nil,
	}
}

func (increment *Increment[T]) loop() {
	defer increment.wg.Done()

	ticker := time.NewTicker(increment.interval)
	defer ticker.Stop()

	for {
		select {
		case <-increment.pause:
			ticker.Stop()
		case <-increment.resume:
			ticker.Reset(increment.interval)
		case <-increment.reset:
			increment.mutex.Lock()
			increment.current = increment.initial
			increment.events.Emit(event.Changed(increment.name), event.ChangedPayload{
				Resource: increment.name,
				Value:    increment.current,
			})
			increment.mutex.Unlock()
		case <-ticker.C:
			increment.mutex.Lock()
			increment.current += increment.step
			increment.events.Emit(event.Changed(increment.name), event.ChangedPayload{
				Resource: increment.name,
				Value:    increment.current,
			})
			increment.mutex.Unlock()
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
	increment.reset = make(chan any)
	increment.pause = make(chan any)
	increment.resume = make(chan any)

	increment.wg.Add(1)
	go increment.loop()

	increment.events.Subscribe(event.Action(increment.name, "reset"), increment.reset)
	increment.events.Subscribe(event.Action(increment.name, "resume"), increment.resume)
	increment.events.Subscribe(event.Action(increment.name, "pause"), increment.pause)
}

func (increment *Increment[T]) Stop() {
	increment.events.Unsubscribe(event.Action(increment.name, "reset"), increment.reset)
	increment.events.Unsubscribe(event.Action(increment.name, "resume"), increment.resume)
	increment.events.Unsubscribe(event.Action(increment.name, "pause"), increment.pause)

	close(increment.quit)

	increment.wg.Wait()

	close(increment.reset)
	close(increment.resume)
	close(increment.pause)

	increment.name = ""
	increment.storage = nil
	increment.events = nil
	increment.quit = nil
	increment.reset = nil
	increment.resume = nil
	increment.pause = nil
}

func (increment *Increment[T]) Read() (any, error) {
	increment.mutex.RLock()
	defer increment.mutex.RUnlock()

	return increment.current, nil
}
