package test

import (
	"context"
	"testing"

	"github.com/auchon/auchon-api/lib/command"
	"github.com/auchon/auchon-api/lib/event"
)

// EventPublisherSpy records published events.
type EventPublisherSpy struct {
	Events []event.Event
	Fail   error
}

func NewEventPublisherSpy() *EventPublisherSpy { return &EventPublisherSpy{} }

func (s *EventPublisherSpy) Publish(_ context.Context, events ...event.Event) error {
	if s.Fail != nil {
		return s.Fail
	}
	s.Events = append(s.Events, events...)
	return nil
}

// AssertPublished fails unless exactly the named events were published, in order.
func (s *EventPublisherSpy) AssertPublished(t testing.TB, names ...string) {
	t.Helper()
	got := make([]string, len(s.Events))
	for i, e := range s.Events {
		got[i] = e.EventName()
	}
	if len(got) == 0 && len(names) == 0 {
		return
	}
	AssertEqual(t, got, names)
}

// LastEvent returns the last published event as T.
func LastEvent[T event.Event](t testing.TB, s *EventPublisherSpy) T {
	t.Helper()
	if len(s.Events) == 0 {
		t.Fatal("no events published")
	}
	e, ok := s.Events[len(s.Events)-1].(T)
	if !ok {
		t.Fatalf("last event is %T", s.Events[len(s.Events)-1])
	}
	return e
}

// CommandBusSpy records dispatched commands; controller and event-handler tests use it.
type CommandBusSpy struct {
	Commands []command.Command
	// Results, when set, returns an error per command name.
	Results map[string]error
}

func NewCommandBusSpy() *CommandBusSpy { return &CommandBusSpy{} }

func (s *CommandBusSpy) Dispatch(_ context.Context, cmd command.Command) error {
	s.Commands = append(s.Commands, cmd)
	return s.Results[cmd.CommandName()]
}

// AssertDispatched fails unless exactly the given commands were dispatched, in order.
func (s *CommandBusSpy) AssertDispatched(t testing.TB, cmds ...command.Command) {
	t.Helper()
	if len(s.Commands) == 0 && len(cmds) == 0 {
		return
	}
	AssertEqual(t, s.Commands, cmds)
}
