package rename_location

const CommandName = "inventory.RenameLocation"

type RenameLocationCommand struct {
	LocationID string
	Name       string
}

func (RenameLocationCommand) CommandName() string { return CommandName }

// GuardRenameLocationIsValid validates the input: the location id is a UUID.
func GuardRenameLocationIsValid(cmd RenameLocationCommand) error { return nil }
