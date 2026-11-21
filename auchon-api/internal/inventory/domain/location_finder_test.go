package domain_test

import (
	"errors"
	"testing"

	"github.com/auchon/auchon-api/internal/inventory/domain"
)

func TestLocationFinder(t *testing.T) {
	t.Run("A stored box comes back as it was stored", func(t *testing.T) {
		reset()
		givenStored(tools())
		whenFoundByID(toolsID)
		thenItIs(t, domain.LocationDTO{ID: toolsID, Name: "Tools", Kind: "container", ParentID: shelfAID})
		thenNothingRecorded(t)
	})

	t.Run("A location that was never stored is not found", func(t *testing.T) {
		reset()
		whenFoundByID(toolsID)
		thenRejected[domain.LocationNotFoundError](t)
	})

	t.Run("A removed location is not found", func(t *testing.T) {
		reset()
		givenStored(tools().Removed())
		whenFoundByID(toolsID)
		thenRejected[domain.LocationNotFoundError](t)
	})

	t.Run("A repository failure is returned", func(t *testing.T) {
		reset()
		givenTheRepositoryFails()
		whenFoundByID(toolsID)
		thenTheRepositoryFailureIsReturned(t)
	})
}

// ---- steps ----

var errDiskFull = errors.New("disk full")

func givenTheRepositoryFails() { repo.Fail = errDiskFull }

func whenFoundByID(id string) { location, err = finder.FindByID(ctx, id) }

func thenTheRepositoryFailureIsReturned(t *testing.T) {
	t.Helper()
	if !errors.Is(err, errDiskFull) {
		t.Errorf("err = %v, want disk full", err)
	}
}
