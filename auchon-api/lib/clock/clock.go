// Package clock is the kernel's port for "now". Domain and application code never call
// time.Now().
package clock

import "time"

type Clock interface {
	Now() time.Time
}

// System is the production clock (UTC).
type System struct{}

func (System) Now() time.Time { return time.Time{} }

// Stub is the test kit clock: a fixed, settable time.
type Stub struct{ now time.Time }

func NewStub(at time.Time) *Stub { return &Stub{} }

func (s *Stub) Now() time.Time           { return time.Time{} }
func (s *Stub) Set(at time.Time)         {}
func (s *Stub) Advance(d time.Duration) {}
