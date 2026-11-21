package domain

import "github.com/auchon/auchon-api/lib/event"

const LocationCreatedEventName = "inventory.LocationCreated"

// LocationCreated: the location was added.
type LocationCreated struct {
	event.Base
	LocationCreatedPayload
}

type LocationCreatedPayload struct {
	LocationID string
	Name       string
	Kind       string
	ParentID   string
}

func (LocationCreated) EventName() string { return LocationCreatedEventName }
