package name_test

import (
	"strings"
	"testing"

	"github.com/auchon/auchon-api/internal/inventory/domain/name"
	"github.com/auchon/auchon-api/lib/test"
)

func TestNew(t *testing.T) {
	t.Run("Spaces around a name are trimmed", func(t *testing.T) {
		reset()
		whenNamed("  Metal shelf  ")
		thenNameIs(t, "Metal shelf")
	})

	t.Run("An empty name is required", func(t *testing.T) {
		reset()
		whenNamed("")
		thenRejected[name.NameRequiredError](t)
	})

	t.Run("Only spaces is no name", func(t *testing.T) {
		reset()
		whenNamed("   ")
		thenRejected[name.NameRequiredError](t)
	})

	t.Run("Exactly the maximum length is valid", func(t *testing.T) {
		reset()
		whenNamed(strings.Repeat("a", name.MaxLength))
		thenSucceeded(t)
	})

	t.Run("One character too many is too long", func(t *testing.T) {
		reset()
		whenNamed(strings.Repeat("a", name.MaxLength+1))
		thenRejected[name.NameTooLongError](t)
	})

	t.Run("The length counts characters, not bytes", func(t *testing.T) {
		reset()
		whenNamed(strings.Repeat("ñ", name.MaxLength))
		thenSucceeded(t)
	})

	t.Run("Names with the same text are equal", func(t *testing.T) {
		reset()
		thenEqual(t, "Tools", " Tools ")
		thenNotEqual(t, "Tools", "Screws")
	})
}

// ---- steps ----

var (
	got name.Name
	err error
)

func reset() { got = name.Name{}; err = nil }

func whenNamed(raw string) { got, err = name.New(raw) }

func thenSucceeded(t *testing.T) { t.Helper(); test.AssertNoError(t, err) }

func thenNameIs(t *testing.T, want string) {
	t.Helper()
	test.AssertNoError(t, err)
	test.AssertEqual(t, got.String(), want)
}

func thenRejected[E error](t *testing.T) {
	t.Helper()
	test.AssertErrorIsOfType[E](t, err)
}

func thenEqual(t *testing.T, a, b string) {
	t.Helper()
	if x, y := mustName(a), mustName(b); !x.Equals(y) {
		t.Errorf("%q Equals %q = false, want true", a, b)
	}
}

func thenNotEqual(t *testing.T, a, b string) {
	t.Helper()
	if x, y := mustName(a), mustName(b); x.Equals(y) {
		t.Errorf("%q Equals %q = true, want false", a, b)
	}
}

func mustName(raw string) name.Name { n, _ := name.New(raw); return n }
