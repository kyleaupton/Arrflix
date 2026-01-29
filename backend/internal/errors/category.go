// Package errors provides error categorization for retry logic.
// Transient errors are retried; permanent errors fail immediately.
package errors

import (
	"errors"
	"fmt"
)

// Category indicates whether an error is transient or permanent.
type Category string

const (
	// Transient errors should be retried (network timeouts, 5xx, temp failures)
	Transient Category = "transient"
	// Permanent errors should not be retried (invalid input, auth failures)
	Permanent Category = "permanent"
)

// categorizedError wraps an error with a category.
type categorizedError struct {
	category Category
	err      error
}

func (e *categorizedError) Error() string {
	return e.err.Error()
}

func (e *categorizedError) Unwrap() error {
	return e.err
}

// Category returns the error's category.
func (e *categorizedError) Category() Category {
	return e.category
}

// AsTransient wraps an error as transient (should be retried).
func AsTransient(err error) error {
	if err == nil {
		return nil
	}
	return &categorizedError{
		category: Transient,
		err:      err,
	}
}

// AsPermanent wraps an error as permanent (should not be retried).
func AsPermanent(err error) error {
	if err == nil {
		return nil
	}
	return &categorizedError{
		category: Permanent,
		err:      err,
	}
}

// AsTransientf creates a transient error with a formatted message.
func AsTransientf(format string, args ...any) error {
	return AsTransient(fmt.Errorf(format, args...))
}

// AsPermanentf creates a permanent error with a formatted message.
func AsPermanentf(format string, args ...any) error {
	return AsPermanent(fmt.Errorf(format, args...))
}

// CategoryOf extracts the category from an error.
// Returns Transient if the error is not categorized (safe default: retry).
func CategoryOf(err error) Category {
	if err == nil {
		return Transient
	}

	var ce *categorizedError
	if errors.As(err, &ce) {
		return ce.category
	}

	// Default: assume transient, allow retry
	return Transient
}

// IsPermanent returns true if the error is categorized as permanent.
func IsPermanent(err error) bool {
	return CategoryOf(err) == Permanent
}

// IsTransient returns true if the error is categorized as transient.
func IsTransient(err error) bool {
	return CategoryOf(err) == Transient
}
