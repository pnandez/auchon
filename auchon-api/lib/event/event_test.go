package event_test

import (
	"testing"
	"time"

	"github.com/auchon/auchon-api/lib/event"
)

type boxLabelled struct {
	event.Base
	Label string
}

func (boxLabelled) EventName() string { return "test.BoxLabelled" }

var at = time.Date(2026, 11, 20, 18, 0, 0, 0, time.UTC)

func TestBase(t *testing.T) {
	t.Run("Given a new base event", func(t *testing.T) {
		b := event.NewBase("box-1", at)

		t.Run("Then it carries the aggregate id and the moment", func(t *testing.T) {
			if b.AggregateID() != "box-1" {
				t.Errorf("AggregateID() = %q", b.AggregateID())
			}
			if !b.OccurredAt().Equal(at) {
				t.Errorf("OccurredAt() = %v", b.OccurredAt())
			}
		})

		t.Run("Then it has a unique id", func(t *testing.T) {
			other := event.NewBase("box-1", at)
			if b.EventID() == "" || b.EventID() == other.EventID() {
				t.Errorf("ids %q and %q must be non-empty and different", b.EventID(), other.EventID())
			}
		})
	})
}

func TestRecorder(t *testing.T) {
	t.Run("Given a recorder with two events", func(t *testing.T) {
		var r event.Recorder
		first := boxLabelled{Base: event.NewBase("box-1", at), Label: "Cables"}
		second := boxLabelled{Base: event.NewBase("box-1", at), Label: "Tools"}
		r.Record(first)
		r.Record(second)

		t.Run("When the events are pulled", func(t *testing.T) {
			got := r.PullEvents()

			t.Run("Then they come back in order", func(t *testing.T) {
				if len(got) != 2 || got[0] != event.Event(first) || got[1] != event.Event(second) {
					t.Fatalf("PullEvents() = %v", got)
				}
			})

			t.Run("And pulling again returns nothing", func(t *testing.T) {
				if again := r.PullEvents(); len(again) != 0 {
					t.Errorf("second PullEvents() = %v, want empty", again)
				}
			})
		})
	})
}
