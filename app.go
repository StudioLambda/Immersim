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

func (application *Application) Read(resource string) (any, error) {
	return application.storage.Read(resource)
}

func (application *Application) Write(resource string, value any) error {
	return application.storage.Write(resource, value)
}

func (application *Application) Action(resource string, action string, payload any) {
	application.events.Emit(event.Action(resource, action), payload)
}

func (application *Application) SubscribeChanges(resource string, listener chan any) {
	application.events.Subscribe(event.Changed(resource), listener)
}

func (application *Application) UnsubscribeChanges(resource string, listener chan<- any) {
	application.events.Unsubscribe(event.Changed(resource), listener)
}

func (application *Application) SubscribeAction(resource string, action string, listener chan any) {
	application.events.Subscribe(event.Action(resource, action), listener)
}

func (application *Application) UnsubscribeAction(resource string, action string, listener chan<- any) {
	application.events.Unsubscribe(event.Action(resource, action), listener)
}

func (application *Application) Start() {
	application.storage.Start(application.events)
}

func (application *Application) Stop() {
	application.storage.Stop()
}
