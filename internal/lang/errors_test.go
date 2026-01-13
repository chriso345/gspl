package lang

import (
	"errors"
	"testing"
)

func TestParseErrorNilAndWrapped(t *testing.T) {
	var pe *ParseError
	if pe.Error() != "" {
		t.Fatalf("expected empty string from nil ParseError.Error(), got %q", pe.Error())
	}
	if pe.Unwrap() != nil {
		t.Fatalf("expected nil Unwrap from nil ParseError")
	}

	orig := errors.New("boom")
	pe2 := &ParseError{Message: "failed", Err: orig}
	if got := pe2.Error(); got != "failed: boom" {
		t.Fatalf("unexpected Error() output: %q", got)
	}
	if pe2.Unwrap() != orig {
		t.Fatalf("Unwrap did not return inner error")
	}
}
