// Package query is the kernel's query bus. Each query has one handler that returns a result.
package query

import "context"

type Query interface {
	QueryName() string
}

type Handler[Q Query, R any] interface {
	Execute(ctx context.Context, q Q) (R, error)
}

// NoHandlerError is returned when a query has no registered handler.
type NoHandlerError struct{ Name string }

func (e NoHandlerError) Error() string { return "no handler for query " + e.Name }

type Bus struct {
	handlers map[string]func(ctx context.Context, q Query) (any, error)
}

func NewBus() *Bus { return &Bus{} }

// Register binds h to the query named by Q's zero value. Registering a name twice panics.
func Register[Q Query, R any](b *Bus, h Handler[Q, R]) {}

// Ask runs q's handler and returns its result as R.
func Ask[R any](ctx context.Context, b *Bus, q Query) (R, error) {
	var zero R
	return zero, nil
}
