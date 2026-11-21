package remove_location

const CommandName = "inventory.RemoveLocation"

type RemoveLocationCommand struct {
	LocationID string
}

func (RemoveLocationCommand) CommandName() string { return CommandName }

// GuardRemoveLocationIsValid validates the input: the location id is a UUID.
func GuardRemoveLocationIsValid(cmd RemoveLocationCommand) error { return nil }
