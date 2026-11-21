package domain

import "github.com/google/uuid"

// GuardLocationID returns InvalidLocationIDError unless id is a UUID. Shared by every use case
// that receives a location id from outside.
func GuardLocationID(id string) error {
	// uuid.Parse also accepts the urn, braced and undashed forms; ids are always the dashed one.
	if len(id) != 36 {
		return NewInvalidLocationIDError(id)
	}
	if _, err := uuid.Parse(id); err != nil {
		return NewInvalidLocationIDError(id)
	}
	return nil
}
