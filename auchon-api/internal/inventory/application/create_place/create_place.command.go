package create_place

import "github.com/auchon/auchon-api/internal/inventory/domain"

const CommandName = "inventory.CreatePlace"

type CreatePlaceCommand struct {
	LocationID string
	Name       string
}

func (CreatePlaceCommand) CommandName() string { return CommandName }

// GuardCreatePlaceIsValid validates the input: the location id is a UUID.
func GuardCreatePlaceIsValid(cmd CreatePlaceCommand) error {
	return domain.GuardLocationID(cmd.LocationID)
}
