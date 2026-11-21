package clock_test

import (
	"testing"
	"time"

	"github.com/auchon/auchon-api/lib/clock"
)

func TestStub(t *testing.T) {
	t.Run("Given a stub set to a moment", func(t *testing.T) {
		at := time.Date(2026, 11, 20, 18, 0, 0, 0, time.UTC)
		c := clock.NewStub(at)

		t.Run("Then Now returns that moment", func(t *testing.T) {
			if !c.Now().Equal(at) {
				t.Errorf("Now() = %v, want %v", c.Now(), at)
			}
		})

		t.Run("When it advances an hour, Then Now moves an hour", func(t *testing.T) {
			c.Advance(time.Hour)
			if want := at.Add(time.Hour); !c.Now().Equal(want) {
				t.Errorf("Now() = %v, want %v", c.Now(), want)
			}
		})

		t.Run("When it is set to another moment, Then Now returns it", func(t *testing.T) {
			other := at.AddDate(0, 0, 1)
			c.Set(other)
			if !c.Now().Equal(other) {
				t.Errorf("Now() = %v, want %v", c.Now(), other)
			}
		})
	})
}

func TestSystem(t *testing.T) {
	t.Run("Given the system clock, Then Now is current and in UTC", func(t *testing.T) {
		now := clock.System{}.Now()
		if now.Location() != time.UTC {
			t.Errorf("location = %v, want UTC", now.Location())
		}
		if time.Since(now) > time.Minute || time.Since(now) < 0 {
			t.Errorf("Now() = %v is not current", now)
		}
	})
}
