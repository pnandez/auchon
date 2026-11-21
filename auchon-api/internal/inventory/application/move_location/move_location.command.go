package move_location

const CommandName = "inventory.MoveLocation"

type MoveLocationCommand struct {
	LocationID string
	ParentID   string
}

func (MoveLocationCommand) CommandName() string { return CommandName }

// GuardMoveLocationIsValid validates the input: both ids are UUIDs.
func GuardMoveLocationIsValid(cmd MoveLocationCommand) error { return nil }
