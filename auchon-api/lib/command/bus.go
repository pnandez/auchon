// Package command is the kernel's command bus. Each command has one handler; every dispatch
// runs in a database transaction.
package command

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/auchon/auchon-api/lib/db"
)

type Command interface {
	CommandName() string
}

// Handler is the typed handler a use case implements.
type Handler[C Command] interface {
	Execute(ctx context.Context, cmd C) error
}

// Dispatcher is what controllers depend on (and what the command bus spy implements).
type Dispatcher interface {
	Dispatch(ctx context.Context, cmd Command) error
}

// NoHandlerError is returned when a command has no registered handler.
type NoHandlerError struct{ Name string }

func (e NoHandlerError) Error() string { return "no handler for command " + e.Name }

type Bus struct {
	conn     *sql.DB
	handlers map[string]func(ctx context.Context, cmd Command) error
}

func NewBus(conn *sql.DB) *Bus {
	return &Bus{conn: conn, handlers: map[string]func(context.Context, Command) error{}}
}

// Register binds h to the command named by C's zero value. Registering a name twice panics.
func Register[C Command](b *Bus, h Handler[C]) {
	var zero C
	name := zero.CommandName()
	if _, exists := b.handlers[name]; exists {
		panic(fmt.Sprintf("command %s already has a handler", name))
	}
	b.handlers[name] = func(ctx context.Context, cmd Command) error {
		typed, ok := cmd.(C)
		if !ok {
			return fmt.Errorf("command %s: got %T", name, cmd)
		}
		return h.Execute(ctx, typed)
	}
}

// Dispatch runs the command's handler inside a transaction.
func (b *Bus) Dispatch(ctx context.Context, cmd Command) error {
	handle, ok := b.handlers[cmd.CommandName()]
	if !ok {
		return NoHandlerError{Name: cmd.CommandName()}
	}
	return db.WithTx(ctx, b.conn, func(ctx context.Context) error {
		return handle(ctx, cmd)
	})
}
