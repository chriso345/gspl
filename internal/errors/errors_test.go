package errors

import (
	"errors"
	"testing"
)

func TestErrorKindAndMessage(t *testing.T) {
	err := New(ErrInfeasible, "problem infeasible", nil)
	if err.Kind != ErrInfeasible {
		t.Errorf("expected kind %v, got %v", ErrInfeasible, err.Kind)
	}
	if err.Error() != "problem infeasible" {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}

func TestErrorWrapping(t *testing.T) {
	base := errors.New("base error")
	err := New(ErrNumericalFailure, "failed numerics", base)
	if !errors.Is(err, base) {
		t.Errorf("wrapped error not found with errors.Is")
	}
	if err.Unwrap() != base {
		t.Errorf("Unwrap did not return base error")
	}
	msg := err.Error()
	if msg != "failed numerics: base error" {
		t.Errorf("unexpected wrapped error message: %q", msg)
	}
}
