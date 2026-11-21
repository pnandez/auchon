package create_location

import "github.com/auchon/auchon-api/internal/inventory/domain"

const CommandName = "inventory.CreateLocation"

type CreateLocationCommand struct {
	LocationID string
	ParentID   string
	Name       string
	Kind       string // "area" or "container"
}

func (CreateLocationCommand) CommandName() string { return CommandName }

// GuardCreateLocationIsValid validates the input: both ids are UUIDs and the kind exists.
func GuardCreateLocationIsValid(cmd CreateLocationCommand) error {
	if err := domain.GuardLocationID(cmd.LocationID); err != nil {
		return err
	}
	if err := domain.GuardLocationID(cmd.ParentID); err != nil {
		return err
	}
	_, err := domain.ParseKind(cmd.Kind)
	return err
}
