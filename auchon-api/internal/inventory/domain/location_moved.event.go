package domain

import "github.com/auchon/auchon-api/lib/event"

const LocationMovedEventName = "inventory.LocationMoved"

// LocationMoved: the location was put inside another one.
type LocationMoved struct {
	event.Base
	LocationMovedPayload
}

type LocationMovedPayload struct {
	LocationID string
	From       string
	To         string
}

func (LocationMoved) EventName() string { return LocationMovedEventName }
