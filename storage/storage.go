package storage

import (
	"errors"
	"fmt"

	"github.com/studiolambda/immersim/event"
)

type Reader interface {
	Read() (any, error)
}

type Writer interface {
	Write(value any) error
}

type Resource interface {
	Start(name string, storage *Storage, events *event.Events)
	Stop()
}

type SupportedNumeric interface {
	int32 | float32
}

type Supported interface {
	SupportedNumeric | bool
}

type Storage struct {
	memory map[string]Resource
}

var (
	ErrRead                = errors.New("failed to read resource")
	ErrWrite               = errors.New("failed to write resource")
	ErrResourceNotReadable = errors.New("resource is not readable")
	ErrResourceNotWritable = errors.New("resource is not writable")
)

func NewStorage(memory map[string]Resource) *Storage {
	return &Storage{
		memory: memory,
	}
}

func (storage *Storage) Start(events *event.Events) {
	for name, resource := range storage.memory {
		resource.Start(name, storage, events)
	}
}

func (storage *Storage) Stop() {
	for _, resource := range storage.memory {
		resource.Stop()
	}
}

func (storage *Storage) Read(resource string) (any, error) {
	if reader, ok := storage.memory[resource].(Reader); ok {
		result, err := reader.Read()

		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("%w: %s", ErrRead, resource),
				err,
			)
		}

		return result, nil
	}

	return nil, errors.Join(
		fmt.Errorf("%w: %s", ErrRead, resource),
		ErrResourceNotReadable,
	)
}

func (storage *Storage) Write(resource string, value any) error {
	if writer, ok := storage.memory[resource].(Writer); ok {
		if err := writer.Write(value); err != nil {
			return errors.Join(
				fmt.Errorf("%w: %s", ErrWrite, resource),
				err,
			)
		}

		return nil
	}

	return errors.Join(
		fmt.Errorf("%w: %s", ErrWrite, resource),
		ErrResourceNotWritable,
	)
}
