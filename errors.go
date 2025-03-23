package mediasoupgo

import "fmt"

// UnsupportedError represents an error for unsupported operations.
type UnsupportedError struct {
	message string
}

// NewUnsupportedError creates a new UnsupportedError with the given message.
func NewUnsupportedError(message string) *UnsupportedError {
	return &UnsupportedError{
		message: message,
	}
}

// Error implements the error interface for UnsupportedError.
func (e *UnsupportedError) Error() string {
	return fmt.Sprintf("UnsupportedError: %s", e.message)
}

// InvalidStateError represents an error produced when calling a method in an invalid state.
type InvalidStateError struct {
	message string
}

// NewInvalidStateError creates a new InvalidStateError with the given message.
func NewInvalidStateError(message string) *InvalidStateError {
	return &InvalidStateError{
		message: message,
	}
}

// Error implements the error interface for InvalidStateError.
func (e *InvalidStateError) Error() string {
	return fmt.Sprintf("InvalidStateError: %s", e.message)
}
