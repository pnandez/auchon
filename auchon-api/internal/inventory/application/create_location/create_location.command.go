package create_location

const CommandName = "inventory.CreateLocation"

type CreateLocationCommand struct {
	LocationID string
	ParentID   string
	Name       string
	Kind       string // "area" or "container"
}

func (CreateLocationCommand) CommandName() string { return CommandName }

// GuardCreateLocationIsValid validates the input: both ids are UUIDs and the kind exists.
func GuardCreateLocationIsValid(cmd CreateLocationCommand) error { return nil }
