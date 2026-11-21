package move_location

import "github.com/auchon/auchon-api/internal/inventory/domain"

const CommandName = "inventory.MoveLocation"

type MoveLocationCommand struct {
	LocationID string
	ParentID   string
}

func (MoveLocationCommand) CommandName() string { return CommandName }

// GuardMoveLocationIsValid validates the input: both ids are UUIDs.
func GuardMoveLocationIsValid(cmd MoveLocationCommand) error {
	if err := domain.GuardLocationID(cmd.LocationID); err != nil {
		return err
	}
	return domain.GuardLocationID(cmd.ParentID)
}
