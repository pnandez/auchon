package remove_location

import (
	"context"

	"github.com/auchon/auchon-api/internal/inventory/domain"
	"github.com/auchon/auchon-api/lib/clock"
	"github.com/auchon/auchon-api/lib/event"
)

type Handler struct {
	finder    *domain.LocationFinder
	storer    *domain.LocationStorer
	contents  domain.LocationContentsChecker
	publisher event.Publisher
	clock     clock.Clock
}

func NewHandler(finder *domain.LocationFinder, storer *domain.LocationStorer, contents domain.LocationContentsChecker, publisher event.Publisher, clock clock.Clock) *Handler {
	return &Handler{finder: finder, storer: storer, contents: contents, publisher: publisher, clock: clock}
}

func (h *Handler) Execute(ctx context.Context, cmd RemoveLocationCommand) error {
	if err := GuardRemoveLocationIsValid(cmd); err != nil {
		return err
	}
	location, err := h.finder.FindByID(ctx, cmd.LocationID)
	if err != nil {
		return err
	}
	hasContents, err := h.contents.HasContents(ctx, cmd.LocationID)
	if err != nil {
		return err
	}
	if hasContents {
		return NewLocationNotEmptyError(cmd.LocationID)
	}
	location.Remove(h.clock.Now())
	if err := h.storer.Store(ctx, location); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, location.PullEvents()...)
}
