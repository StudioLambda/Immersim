package resource

import (
	"math/rand/v2"
	"sync"
	"time"

	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type Random[T storage.Supported] struct {
	name     string
	current  T
	min      T
	max      T
	interval time.Duration
	mutex    sync.RWMutex
	events   *event.Events
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewRandom[T storage.Supported](min T, max T, interval time.Duration) *Random[T] {
	return &Random[T]{
		name:     "",
		events:   nil,
		current:  *new(T),
		min:      min,
		max:      max,
		interval: interval,
		mutex:    sync.RWMutex{},
	}
}

func (random *Random[T]) loop() {
	defer random.wg.Done()

	ticker := time.NewTicker(random.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var value any

			switch any(random.current).(type) {
			case int32:
				min := any(random.min).(int32)
				max := any(random.max).(int32)
				value = min + rand.Int32N(max-min+1)
			case float32:
				min := any(random.min).(float32)
				max := any(random.max).(float32)
				value = min + rand.Float32()*(max-min)
			case bool:
				value = rand.Int32N(2) == 1
			}

			random.mutex.Lock()
			random.current = value.(T)
			random.mutex.Unlock()

			random.events.Emit(random.name)
		case <-random.quit:
			return
		}
	}
}

func (random *Random[T]) Start(name string, storage *storage.Storage, events *event.Events) {
	random.name = name
	random.events = events
	random.quit = make(chan struct{})

	random.wg.Add(1)
	go random.loop()
}

func (random *Random[T]) Stop() {
	close(random.quit)
	random.wg.Wait()

	random.name = ""
	random.events = nil
	random.quit = nil
}

func (random *Random[T]) Read() (any, error) {
	random.mutex.RLock()
	defer random.mutex.RUnlock()

	return random.current, nil
}
