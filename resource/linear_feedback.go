package resource

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type LinearFeedback[T storage.SupportedNumeric] struct {
	current      T
	step         T
	events       *event.Events
	storage      *storage.Storage
	name         string
	stepInterval time.Duration
	setpoint     string
	mutex        sync.RWMutex
	listener     chan struct{}
	wg           sync.WaitGroup
}

var (
	ErrNotNumeric = errors.New("setpoint type must be int32 or float32")
)

func NewLinearFeedback[T storage.SupportedNumeric](step T, stepInterval time.Duration, setpoint string) *LinearFeedback[T] {
	return &LinearFeedback[T]{
		current:      *new(T),
		events:       nil,
		storage:      nil,
		name:         "",
		step:         step,
		stepInterval: stepInterval,
		setpoint:     setpoint,
		mutex:        sync.RWMutex{},
		listener:     make(chan struct{}),
		wg:           sync.WaitGroup{},
	}
}

func (feedback *LinearFeedback[T]) loop() {
	defer feedback.wg.Done()

	ticker := time.NewTicker(feedback.stepInterval)
	defer ticker.Stop()

	target, _ := feedback.readSetpoint()

	for {
		select {
		case _, ok := <-feedback.listener:
			if !ok {
				return
			}

			target, _ = feedback.readSetpoint()
		case <-ticker.C:
			feedback.mutex.Lock()

			if feedback.current < target {
				feedback.current += feedback.step

				if feedback.current > target {
					feedback.current = target
				}

				feedback.events.Emit(feedback.name)
			} else if feedback.current > target {
				feedback.current -= feedback.step

				if feedback.current < target {
					feedback.current = target
				}

				feedback.events.Emit(feedback.name)
			}

			feedback.mutex.Unlock()
		}
	}
}

func (feedback *LinearFeedback[T]) readSetpoint() (T, error) {
	value, err := feedback.storage.Read(feedback.setpoint)

	if err != nil {
		return *new(T), err
	}

	switch v := value.(type) {
	case int32:
		switch any(feedback.current).(type) {
		case int32:
			return any(v).(T), nil
		case float32:
			return any(float32(v)).(T), nil
		}
	case float32:
		switch any(feedback.current).(type) {
		case int32:
			return any(int32(v)).(T), nil
		case float32:
			return any(v).(T), nil
		}
	}

	return *new(T), fmt.Errorf("%w: %T", ErrNotNumeric, value)
}

func (feedback *LinearFeedback[T]) Start(name string, storage *storage.Storage, events *event.Events) {
	feedback.name = name
	feedback.storage = storage
	feedback.events = events
	feedback.listener = make(chan struct{})

	feedback.wg.Add(1)
	go feedback.loop()

	feedback.events.Subscribe(feedback.setpoint, feedback.listener)
}

func (feedback *LinearFeedback[T]) Stop() {
	feedback.events.Unsubscribe(feedback.setpoint, feedback.listener)
	close(feedback.listener)

	feedback.wg.Done()

	feedback.name = ""
	feedback.storage = nil
	feedback.events = nil
}

func (feedback *LinearFeedback[T]) Read() (any, error) {
	feedback.mutex.RLock()
	defer feedback.mutex.RUnlock()

	return feedback.current, nil
}
