// Package test is the kernel's test kit: assertions and spies shared by every context.
package test

import (
	"errors"
	"reflect"
	"testing"
)

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertErrorIsOfType fails unless err (or an error it wraps) is a T.
func AssertErrorIsOfType[T error](t testing.TB, err error) T {
	t.Helper()
	var target T
	if !errors.As(err, &target) {
		t.Fatalf("error %v (%T) is not a %T", err, err, target)
	}
	return target
}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
