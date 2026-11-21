package create_place_test

import (
	"context"
	"testing"
	"time"

	"github.com/auchon/auchon-api/internal/inventory/application/create_place"
	"github.com/auchon/auchon-api/internal/inventory/domain"
	"github.com/auchon/auchon-api/internal/inventory/domain/name"
	inventorytest "github.com/auchon/auchon-api/internal/inventory/infra/test"
	"github.com/auchon/auchon-api/lib/clock"
	"github.com/auchon/auchon-api/lib/test"
)

func TestCreatePlace(t *testing.T) {
	t.Run("Create the basement", func(t *testing.T) {
		reset()
		whenPlaceCreated(basementID, "Basement")
		thenSaved(t, domain.LocationDTO{ID: basementID, Name: "Basement", Kind: "place"})
		thenPublished(t, domain.LocationCreatedEventName)
	})

	t.Run("An id that is not a UUID is rejected", func(t *testing.T) {
		reset()
		whenPlaceCreated("basement", "Basement")
		thenRejected[domain.InvalidLocationIDError](t)
		thenNothingSaved(t)
		thenPublished(t)
	})

	t.Run("A place needs a name", func(t *testing.T) {
		reset()
		whenPlaceCreated(basementID, "")
		thenRejected[name.NameRequiredError](t)
		thenNothingSaved(t)
		thenPublished(t)
	})

	t.Run("An id already taken is rejected", func(t *testing.T) {
		reset()
		givenLocation(basementID)
		whenPlaceCreated(basementID, "Basement")
		thenRejected[domain.LocationAlreadyExistsError](t)
		thenPublished(t)
	})
}

// ---- steps ----

const basementID = "01937a4e-0000-7000-8000-000000000001"

var (
	ctx     = context.Background()
	repo    = inventorytest.NewLocationRepositoryMock()
	events  = test.NewEventPublisherSpy()
	handler = create_place.NewHandler(domain.NewLocationStorer(repo), events,
		clock.NewStub(time.Date(2026, 11, 21, 18, 30, 0, 0, time.UTC)))
	err error
)

func reset() { repo.Reset(); events.Reset(); err = nil }

func givenLocation(id string) {
	repo.WithLocation(inventorytest.NewLocationBuilder().WithID(id).DTO())
}

func whenPlaceCreated(id, n string) {
	err = handler.Execute(ctx, create_place.CreatePlaceCommand{LocationID: id, Name: n})
}

func thenSaved(t *testing.T, want domain.LocationDTO) {
	t.Helper()
	test.AssertNoError(t, err)
	repo.AssertSaved(t, want)
}

func thenRejected[E error](t *testing.T) {
	t.Helper()
	test.AssertErrorIsOfType[E](t, err)
}

func thenNothingSaved(t *testing.T) { t.Helper(); repo.AssertNothingSaved(t) }

func thenPublished(t *testing.T, names ...string) { t.Helper(); events.AssertPublished(t, names...) }
