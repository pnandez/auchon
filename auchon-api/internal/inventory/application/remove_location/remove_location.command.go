package remove_location

import "github.com/auchon/auchon-api/internal/inventory/domain"

const CommandName = "inventory.RemoveLocation"

type RemoveLocationCommand struct {
	LocationID string
}

func (RemoveLocationCommand) CommandName() string { return CommandName }

// GuardRemoveLocationIsValid validates the input: the location id is a UUID.
func GuardRemoveLocationIsValid(cmd RemoveLocationCommand) error {
	return domain.GuardLocationID(cmd.LocationID)
}
