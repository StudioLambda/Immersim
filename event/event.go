package event

import "fmt"

type Event string

type ChangedPayload struct {
	Resource string
	Value    any
}

func Changed(resource string) Event {
	return Event(resource)
}

func Action(resource string, action string) Event {
	return Event(fmt.Sprintf("%s:%s", resource, action))
}
