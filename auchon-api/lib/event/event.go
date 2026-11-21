// Package event is the kernel's event recorder, base event and publisher port.
package event

import (
	"context"
	"time"
)

// Event is a domain event. Concrete events embed Base and implement EventName.
type Event interface {
	EventName() string
	EventID() string
	AggregateID() string
	OccurredAt() time.Time
}

// Base is embedded in every event.
type Base struct {
	id          string
	aggregateID string
	occurredAt  time.Time
}

// NewBase creates a base event with a fresh UUIDv7 id.
func NewBase(aggregateID string, occurredAt time.Time) Base {
	return Base{}
}

func (b Base) EventID() string       { return "" }
func (b Base) AggregateID() string   { return "" }
func (b Base) OccurredAt() time.Time { return time.Time{} }

// Recorder is embedded in entities; every state change records an event through it.
type Recorder struct {
	events []Event
}

func (r *Recorder) Record(e Event) {}

// PullEvents returns the recorded events in order and clears them.
func (r *Recorder) PullEvents() []Event { return nil }

// Publisher is the port use cases call after persisting.
type Publisher interface {
	Publish(ctx context.Context, events ...Event) error
}
