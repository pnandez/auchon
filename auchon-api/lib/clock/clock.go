// Package clock is the kernel's port for "now". Domain and application code never call
// time.Now().
package clock

import "time"

type Clock interface {
	Now() time.Time
}

// System is the production clock (UTC).
type System struct{}

func (System) Now() time.Time { return time.Now().UTC() }

// Stub is the test kit clock: a fixed, settable time.
type Stub struct{ now time.Time }

func NewStub(at time.Time) *Stub { return &Stub{now: at} }

func (s *Stub) Now() time.Time          { return s.now }
func (s *Stub) Set(at time.Time)        { s.now = at }
func (s *Stub) Advance(d time.Duration) { s.now = s.now.Add(d) }
