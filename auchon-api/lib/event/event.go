// Package event is the kernel's event recorder, base event and publisher port.
package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Event is a domain event. Concrete events embed Base and implement EventName.
type Event interface {
	EventName() string
	EventID() string
	AggregateID() string
	OccurredAt() time.Time
}

// Base is embedded in every event. Its fields are unexported, so they stay out of the JSON
// payload.
type Base struct {
	id          string
	aggregateID string
	occurredAt  time.Time
}

// NewBase creates a base event with a fresh UUIDv7 id.
func NewBase(aggregateID string, occurredAt time.Time) Base {
	return Base{id: uuid.Must(uuid.NewV7()).String(), aggregateID: aggregateID, occurredAt: occurredAt}
}

func (b Base) EventID() string       { return b.id }
func (b Base) AggregateID() string   { return b.aggregateID }
func (b Base) OccurredAt() time.Time { return b.occurredAt }

// Recorder is embedded in entities; every state change records an event through it.
type Recorder struct {
	events []Event
}

func (r *Recorder) Record(e Event) { r.events = append(r.events, e) }

// PullEvents returns the recorded events in order and clears them.
func (r *Recorder) PullEvents() []Event {
	events := r.events
	r.events = nil
	return events
}

// Publisher is the port use cases call after persisting.
type Publisher interface {
	Publish(ctx context.Context, events ...Event) error
}
