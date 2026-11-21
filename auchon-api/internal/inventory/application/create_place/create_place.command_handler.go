package create_place

import (
	"context"

	"github.com/auchon/auchon-api/internal/inventory/domain"
	"github.com/auchon/auchon-api/lib/clock"
	"github.com/auchon/auchon-api/lib/event"
)

type Handler struct {
	storer    *domain.LocationStorer
	publisher event.Publisher
	clock     clock.Clock
}

func NewHandler(storer *domain.LocationStorer, publisher event.Publisher, clock clock.Clock) *Handler {
	return &Handler{storer: storer, publisher: publisher, clock: clock}
}

func (h *Handler) Execute(ctx context.Context, cmd CreatePlaceCommand) error {
	if err := GuardCreatePlaceIsValid(cmd); err != nil {
		return err
	}
	place, err := domain.CreatePlace(cmd.LocationID, cmd.Name, h.clock.Now())
	if err != nil {
		return err
	}
	if err := h.storer.Store(ctx, place); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, place.PullEvents()...)
}
