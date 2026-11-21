package domain_test

import (
	"testing"

	"github.com/auchon/auchon-api/internal/inventory/domain"
	inventorytest "github.com/auchon/auchon-api/internal/inventory/infra/test"
	"github.com/auchon/auchon-api/lib/test"
)

func TestLocationStorer(t *testing.T) {
	t.Run("A new place is saved", func(t *testing.T) {
		reset()
		givenANew(basement())
		whenStored()
		thenSaved(t, domain.LocationDTO{ID: basementID, Name: "Basement", Kind: "place"})
	})

	t.Run("A new location whose id is taken already exists", func(t *testing.T) {
		reset()
		givenStored(basement())
		givenANew(basement())
		whenStored()
		thenRejected[domain.LocationAlreadyExistsError](t)
		thenNothingSaved(t)
	})

	t.Run("A box moved to another shelf is saved with its new parent", func(t *testing.T) {
		reset()
		givenStored(basement(), shelfA(), shelfB(), tools())
		givenLoaded(t, toolsID)
		givenMovedTo(t, shelfBID)
		whenStored()
		thenSaved(t, domain.LocationDTO{ID: toolsID, Name: "Tools", Kind: "container", ParentID: shelfBID})
	})

	t.Run("A box moved into a box it holds is not saved", func(t *testing.T) {
		reset()
		givenStored(basement(), shelfA(), tools(), screws())
		givenLoaded(t, toolsID)
		givenMovedTo(t, screwsID)
		whenStored()
		thenRejected[domain.LocationCannotContainItselfError](t)
		thenNothingSaved(t)
	})

	t.Run("A stored location stored again is updated, not a duplicate", func(t *testing.T) {
		reset()
		givenStored(basement(), shelfA(), tools())
		givenLoaded(t, toolsID)
		givenRenamed(t, "Hand tools")
		whenStored()
		thenSaved(t, domain.LocationDTO{ID: toolsID, Name: "Hand tools", Kind: "container", ParentID: shelfAID})
	})
}

// ---- steps ----

func givenANew(l *inventorytest.LocationBuilder) { location = l.Build() }

func givenLoaded(t *testing.T, id string) {
	t.Helper()
	location, err = finder.FindByID(ctx, id)
	test.AssertNoError(t, err)
}

func givenMovedTo(t *testing.T, parentID string) {
	t.Helper()
	test.AssertNoError(t, location.MoveTo(parentID, at))
}

func givenRenamed(t *testing.T, n string) {
	t.Helper()
	test.AssertNoError(t, location.Rename(n, at))
}

func whenStored() { err = storer.Store(ctx, location) }

func thenSaved(t *testing.T, want domain.LocationDTO) {
	t.Helper()
	test.AssertNoError(t, err)
	repo.AssertSaved(t, want)
}

func thenNothingSaved(t *testing.T) { t.Helper(); repo.AssertNothingSaved(t) }
