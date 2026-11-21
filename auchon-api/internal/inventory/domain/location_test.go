package domain_test

import (
	"strings"
	"testing"

	"github.com/auchon/auchon-api/internal/inventory/domain"
	"github.com/auchon/auchon-api/internal/inventory/domain/name"
	"github.com/auchon/auchon-api/lib/test"
)

func TestCreatePlace(t *testing.T) {
	t.Run("The basement is a place at the top", func(t *testing.T) {
		reset()
		whenPlaceCreated(basementID, "Basement")
		thenItIs(t, domain.LocationDTO{ID: basementID, Name: "Basement", Kind: "place"})
		thenLocationCreatedRecorded(t, domain.LocationCreatedPayload{
			LocationID: basementID, Name: "Basement", Kind: "place", ParentID: "",
		})
	})

	t.Run("A place needs a name", func(t *testing.T) {
		reset()
		whenPlaceCreated(basementID, "  ")
		thenRejected[name.NameRequiredError](t)
	})

	t.Run("A name that is too long is rejected", func(t *testing.T) {
		reset()
		whenPlaceCreated(basementID, strings.Repeat("a", name.MaxLength+1))
		thenRejected[name.NameTooLongError](t)
	})
}

func TestCreateLocationIn(t *testing.T) {
	t.Run("Put a shelf in the basement", func(t *testing.T) {
		reset()
		whenCreatedIn(shelfAID, "Metal shelf", domain.KindArea, basementID)
		thenItIs(t, domain.LocationDTO{ID: shelfAID, Name: "Metal shelf", Kind: "area", ParentID: basementID})
		thenLocationCreatedRecorded(t, domain.LocationCreatedPayload{
			LocationID: shelfAID, Name: "Metal shelf", Kind: "area", ParentID: basementID,
		})
	})

	t.Run("Put a box on a shelf", func(t *testing.T) {
		reset()
		whenCreatedIn(toolsID, "Tools", domain.KindContainer, shelfAID)
		thenItIs(t, domain.LocationDTO{ID: toolsID, Name: "Tools", Kind: "container", ParentID: shelfAID})
	})

	t.Run("A place cannot be created inside a location", func(t *testing.T) {
		reset()
		whenCreatedIn(garageID, "Garage", domain.KindPlace, basementID)
		thenRejected[domain.InvalidLocationKindError](t)
	})

	t.Run("A box needs somewhere to be", func(t *testing.T) {
		reset()
		whenCreatedIn(toolsID, "Tools", domain.KindContainer, "")
		thenRejected[domain.LocationParentRequiredError](t)
	})

	t.Run("A box cannot be inside itself", func(t *testing.T) {
		reset()
		whenCreatedIn(toolsID, "Tools", domain.KindContainer, toolsID)
		thenRejected[domain.LocationCannotContainItselfError](t)
	})

	t.Run("A box needs a name", func(t *testing.T) {
		reset()
		whenCreatedIn(toolsID, "", domain.KindContainer, shelfAID)
		thenRejected[name.NameRequiredError](t)
	})
}

func TestParseKind(t *testing.T) {
	for _, k := range []domain.Kind{domain.KindPlace, domain.KindArea, domain.KindContainer} {
		t.Run(string(k)+" is a kind", func(t *testing.T) {
			reset()
			whenKindParsed(string(k))
			thenKindIs(t, k)
		})
	}

	t.Run("drawer is not a kind", func(t *testing.T) {
		reset()
		whenKindParsed("drawer")
		thenRejected[domain.InvalidLocationKindError](t)
	})
}

func TestRename(t *testing.T) {
	t.Run("Rename a box", func(t *testing.T) {
		reset()
		givenABox("Box 1", shelfAID)
		whenRenamed("Christmas lights")
		thenItIs(t, domain.LocationDTO{ID: toolsID, Name: "Christmas lights", Kind: "container", ParentID: shelfAID})
		thenLocationRenamedRecorded(t, domain.LocationRenamedPayload{
			LocationID: toolsID, From: "Box 1", To: "Christmas lights",
		})
	})

	t.Run("Renaming to the same name records nothing", func(t *testing.T) {
		reset()
		givenABox("Tools", shelfAID)
		whenRenamed(" Tools ")
		thenSucceeded(t)
		thenNothingRecorded(t)
	})

	t.Run("A box cannot lose its name", func(t *testing.T) {
		reset()
		givenABox("Tools", shelfAID)
		whenRenamed("")
		thenRejected[name.NameRequiredError](t)
		thenItIs(t, domain.LocationDTO{ID: toolsID, Name: "Tools", Kind: "container", ParentID: shelfAID})
		thenNothingRecorded(t)
	})
}

func TestMoveTo(t *testing.T) {
	t.Run("Move a box to another shelf", func(t *testing.T) {
		reset()
		givenABox("Tools", shelfAID)
		whenMovedTo(shelfBID)
		thenItIs(t, domain.LocationDTO{ID: toolsID, Name: "Tools", Kind: "container", ParentID: shelfBID})
		thenLocationMovedRecorded(t, domain.LocationMovedPayload{LocationID: toolsID, From: shelfAID, To: shelfBID})
	})

	t.Run("Moving a box where it already is records nothing", func(t *testing.T) {
		reset()
		givenABox("Tools", shelfAID)
		whenMovedTo(shelfAID)
		thenSucceeded(t)
		thenNothingRecorded(t)
	})

	t.Run("A box cannot go inside itself", func(t *testing.T) {
		reset()
		givenABox("Tools", shelfAID)
		whenMovedTo(toolsID)
		thenRejected[domain.LocationCannotContainItselfError](t)
		thenParentIs(t, shelfAID)
		thenNothingRecorded(t)
	})

	t.Run("A box needs somewhere to go", func(t *testing.T) {
		reset()
		givenABox("Tools", shelfAID)
		whenMovedTo("")
		thenRejected[domain.LocationParentRequiredError](t)
	})

	t.Run("A place stays at the top", func(t *testing.T) {
		reset()
		givenTheGarage(t)
		whenMovedTo(basementID)
		thenRejected[domain.PlaceCannotBeMovedError](t)
		thenParentIs(t, "")
		thenNothingRecorded(t)
	})
}

func TestRemove(t *testing.T) {
	t.Run("Remove an old box", func(t *testing.T) {
		reset()
		givenABox("Old box", shelfAID)
		whenRemoved()
		thenItIs(t, domain.LocationDTO{ID: toolsID, Name: "Old box", Kind: "container", ParentID: shelfAID, Removed: true})
		thenLocationRemovedRecorded(t, domain.LocationRemovedPayload{LocationID: toolsID})
	})
}

func TestGuardLocationID(t *testing.T) {
	t.Run("A UUID is a valid id", func(t *testing.T) {
		reset()
		whenIDGuarded(basementID)
		thenSucceeded(t)
	})

	for _, id := range []string{"", "basement", "01937a4e-0000-7000-8000"} {
		t.Run(strings.TrimSpace(id+" (invalid)")+" is rejected", func(t *testing.T) {
			reset()
			whenIDGuarded(id)
			thenRejected[domain.InvalidLocationIDError](t)
		})
	}
}

// ---- steps ----

// givenABox is the Tools box (by id) with the given name, inside parentID, with no events.
func givenABox(n, parentID string) {
	location = tools().WithName(n).Inside(parentID).Build()
}

func givenTheGarage(t *testing.T) {
	t.Helper()
	location, err = domain.CreatePlace(garageID, "Garage", at)
	test.AssertNoError(t, err)
	location.PullEvents()
}

func whenPlaceCreated(id, n string) { location, err = domain.CreatePlace(id, n, at) }

func whenCreatedIn(id, n string, k domain.Kind, parentID string) {
	location, err = domain.CreateLocationIn(id, n, k, parentID, at)
}

func whenKindParsed(s string) { kind, err = domain.ParseKind(s) }

func whenRenamed(n string) { err = location.Rename(n, at) }

func whenMovedTo(parentID string) { err = location.MoveTo(parentID, at) }

func whenRemoved() { location.Remove(at) }

func whenIDGuarded(id string) { err = domain.GuardLocationID(id) }

func thenKindIs(t *testing.T, want domain.Kind) {
	t.Helper()
	test.AssertNoError(t, err)
	test.AssertEqual(t, kind, want)
}

func thenParentIs(t *testing.T, want string) {
	t.Helper()
	test.AssertEqual(t, location.ParentID(), want)
}

func thenLocationCreatedRecorded(t *testing.T, want domain.LocationCreatedPayload) {
	t.Helper()
	test.AssertEqual(t, recordedOnly[domain.LocationCreated](t).LocationCreatedPayload, want)
}

func thenLocationRenamedRecorded(t *testing.T, want domain.LocationRenamedPayload) {
	t.Helper()
	test.AssertEqual(t, recordedOnly[domain.LocationRenamed](t).LocationRenamedPayload, want)
}

func thenLocationMovedRecorded(t *testing.T, want domain.LocationMovedPayload) {
	t.Helper()
	test.AssertEqual(t, recordedOnly[domain.LocationMoved](t).LocationMovedPayload, want)
}

func thenLocationRemovedRecorded(t *testing.T, want domain.LocationRemovedPayload) {
	t.Helper()
	test.AssertEqual(t, recordedOnly[domain.LocationRemoved](t).LocationRemovedPayload, want)
}
