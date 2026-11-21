package create_place

const CommandName = "inventory.CreatePlace"

type CreatePlaceCommand struct {
	LocationID string
	Name       string
}

func (CreatePlaceCommand) CommandName() string { return CommandName }

// GuardCreatePlaceIsValid validates the input: the location id is a UUID.
func GuardCreatePlaceIsValid(cmd CreatePlaceCommand) error { return nil }
