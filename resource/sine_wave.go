package resource

import (
	"math"
	"sync"
	"time"

	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type SineWave struct {
	current   float32
	frequency float64
	amplitude float64
	offset    float64
	interval  time.Duration
	name      string
	events    *event.Events
	mutex     sync.RWMutex
	quit      chan struct{}
	waitGroup sync.WaitGroup
}

func NewSineWave(frequency float64, amplitude float64, offset float64, interval time.Duration) *SineWave {
	generator := &SineWave{
		current:   0,
		frequency: frequency,
		amplitude: amplitude,
		interval:  interval,
		offset:    offset,
		name:      "",
		events:    nil,
		quit:      nil,
		mutex:     sync.RWMutex{},
		waitGroup: sync.WaitGroup{},
	}

	return generator
}

func (sine *SineWave) Start(name string, storage *storage.Storage, events *event.Events) {
	sine.quit = make(chan struct{})
	sine.name = name
	sine.events = events

	sine.waitGroup.Add(1)

	go sine.loop()
}

func (sine *SineWave) loop() {
	defer sine.waitGroup.Done()

	ticker := time.NewTicker(sine.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t := float64(time.Now().UnixNano()) / 1e9
			sample := sine.amplitude*math.Sin(2*math.Pi*sine.frequency*t) + sine.offset
			sine.mutex.Lock()
			sine.current = float32(sample)
			sine.events.Emit(event.Changed(sine.name), event.ChangedPayload{
				Resource: sine.name,
				Value:    sine.current,
			})
			sine.mutex.Unlock()
		case <-sine.quit:
			return
		}
	}
}

func (sine *SineWave) Stop() {
	close(sine.quit)
	sine.quit = nil
	sine.name = ""
	sine.events = nil

	sine.waitGroup.Wait()
}

func (sine *SineWave) Read() (any, error) {
	sine.mutex.RLock()
	defer sine.mutex.RUnlock()

	return sine.current, nil
}
