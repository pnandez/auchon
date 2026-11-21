package domain

import "github.com/auchon/auchon-api/lib/event"

const LocationRenamedEventName = "inventory.LocationRenamed"

// LocationRenamed: the location got a new name.
type LocationRenamed struct {
	event.Base
	LocationRenamedPayload
}

type LocationRenamedPayload struct {
	LocationID string
	From       string
	To         string
}

func (LocationRenamed) EventName() string { return LocationRenamedEventName }
