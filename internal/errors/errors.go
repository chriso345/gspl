package errors

import (
	"fmt"
)

type ErrorKind int

const (
	ErrUnknown ErrorKind = iota
	ErrInfeasible
	ErrUnbounded
	ErrNumericalFailure
	ErrInvalidInput
)

type Error struct {
	Kind    ErrorKind
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Cause }

func New(kind ErrorKind, msg string, cause error) *Error {
	return &Error{Kind: kind, Message: msg, Cause: cause}
}
