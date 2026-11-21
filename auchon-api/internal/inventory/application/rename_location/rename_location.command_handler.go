package rename_location

import (
	"context"

	"github.com/auchon/auchon-api/internal/inventory/domain"
	"github.com/auchon/auchon-api/lib/clock"
	"github.com/auchon/auchon-api/lib/event"
)

type Handler struct {
	finder    *domain.LocationFinder
	storer    *domain.LocationStorer
	publisher event.Publisher
	clock     clock.Clock
}

func NewHandler(finder *domain.LocationFinder, storer *domain.LocationStorer, publisher event.Publisher, clock clock.Clock) *Handler {
	return &Handler{finder: finder, storer: storer, publisher: publisher, clock: clock}
}

func (h *Handler) Execute(ctx context.Context, cmd RenameLocationCommand) error {
	if err := GuardRenameLocationIsValid(cmd); err != nil {
		return err
	}
	location, err := h.finder.FindByID(ctx, cmd.LocationID)
	if err != nil {
		return err
	}
	if err := location.Rename(cmd.Name, h.clock.Now()); err != nil {
		return err
	}
	events := location.PullEvents()
	if len(events) == 0 {
		return nil // same name
	}
	if err := h.storer.Store(ctx, location); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, events...)
}
