package errors

import "fmt"
import "errors"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common domain errors
var (
	// ErrInvalidEventKey is returned when an event key is empty
	ErrInvalidEventKey = errors.New("invalid event key: key cannot be empty")

	// ErrInvalidEventData is returned when event data is empty
	ErrInvalidEventData = errors.New("invalid event data: data cannot be empty")

	// ErrEventNotFound is returned when an event is not found
	ErrEventNotFound = errors.New("event not found")

	// ErrServiceClosed is returned when trying to use a closed service
	ErrServiceClosed = errors.New("service is closed")
)
