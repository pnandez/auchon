package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/auchon/auchon-api/lib/query"
	"github.com/auchon/auchon-api/lib/test"
)

type countBoxes struct{ InPlace string }

func (countBoxes) QueryName() string { return "test.CountBoxes" }

type countBoxesHandler struct{ fail error }

func (h countBoxesHandler) Execute(_ context.Context, q countBoxes) (int, error) {
	if h.fail != nil {
		return 0, h.fail
	}
	if q.InPlace == "Basement" {
		return 3, nil
	}
	return 0, nil
}

func TestBus(t *testing.T) {
	ctx := context.Background()

	t.Run("Given a registered handler, When asked, Then it returns the typed result", func(t *testing.T) {
		bus := query.NewBus()
		query.Register[countBoxes, int](bus, countBoxesHandler{})
		got, err := query.Ask[int](ctx, bus, countBoxes{InPlace: "Basement"})
		test.AssertNoError(t, err)
		if got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})

	t.Run("Given the handler fails, When asked, Then the error is returned", func(t *testing.T) {
		bus := query.NewBus()
		boom := errors.New("boom")
		query.Register[countBoxes, int](bus, countBoxesHandler{fail: boom})
		_, err := query.Ask[int](ctx, bus, countBoxes{})
		if !errors.Is(err, boom) {
			t.Errorf("err = %v, want boom", err)
		}
	})

	t.Run("Given no handler, When asked, Then it returns NoHandlerError", func(t *testing.T) {
		_, err := query.Ask[int](ctx, query.NewBus(), countBoxes{})
		test.AssertErrorIsOfType[query.NoHandlerError](t, err)
	})

	t.Run("Given a registered query, When registered again, Then it panics", func(t *testing.T) {
		bus := query.NewBus()
		query.Register[countBoxes, int](bus, countBoxesHandler{})
		defer func() {
			if recover() == nil {
				t.Error("expected a panic")
			}
		}()
		query.Register[countBoxes, int](bus, countBoxesHandler{})
	})
}
