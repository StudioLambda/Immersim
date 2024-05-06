package immersim

import (
	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/storage"
)

type Application struct {
	storage *storage.Storage
	events  *event.Events
}

func NewApplication(storage *storage.Storage, events *event.Events) *Application {
	return &Application{
		storage: storage,
		events:  events,
	}
}

func (application *Application) Start() {
	application.storage.Start(application.events)
}

func (application *Application) Stop() {
	application.storage.Stop()
}
