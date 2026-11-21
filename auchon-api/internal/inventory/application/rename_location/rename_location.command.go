package rename_location

import "github.com/auchon/auchon-api/internal/inventory/domain"

const CommandName = "inventory.RenameLocation"

type RenameLocationCommand struct {
	LocationID string
	Name       string
}

func (RenameLocationCommand) CommandName() string { return CommandName }

// GuardRenameLocationIsValid validates the input: the location id is a UUID.
func GuardRenameLocationIsValid(cmd RenameLocationCommand) error {
	return domain.GuardLocationID(cmd.LocationID)
}
