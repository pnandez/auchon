package create_location

import (
	"context"
	"errors"

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

func (h *Handler) Execute(ctx context.Context, cmd CreateLocationCommand) error {
	if err := GuardCreateLocationIsValid(cmd); err != nil {
		return err
	}
	if _, err := h.finder.FindByID(ctx, cmd.ParentID); err != nil {
		return parentError(cmd.ParentID, err)
	}
	location, err := domain.CreateLocationIn(cmd.LocationID, cmd.Name, domain.Kind(cmd.Kind), cmd.ParentID, h.clock.Now())
	if err != nil {
		return err
	}
	if err := h.storer.Store(ctx, location); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, location.PullEvents()...)
}

// parentError reports a missing parent as ParentLocationNotFoundError, not as the location.
func parentError(parentID string, err error) error {
	var notFound domain.LocationNotFoundError
	if errors.As(err, &notFound) {
		return domain.NewParentLocationNotFoundError(parentID)
	}
	return err
}
