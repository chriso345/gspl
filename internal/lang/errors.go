package lang

import "fmt"

var ErrUnsupportedLanguage = fmt.Errorf("unsupported language")
var ErrConflict = fmt.Errorf("registration conflict")

// ParseError represents a parsing error with an optional message and cause.
type ParseError struct {
	Message string
	Err     error
}

func (e *ParseError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *ParseError) Unwrap() error { return e.Err }
