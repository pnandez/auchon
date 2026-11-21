package domain

import "github.com/auchon/auchon-api/lib/event"

const LocationRemovedEventName = "inventory.LocationRemoved"

// LocationRemoved: the location was taken out of the inventory.
type LocationRemoved struct {
	event.Base
	LocationRemovedPayload
}

type LocationRemovedPayload struct {
	LocationID string
}

func (LocationRemoved) EventName() string { return LocationRemovedEventName }
