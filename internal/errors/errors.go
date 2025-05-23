/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/
package errors

import (
	"fmt"
)

// Error types
type (
	// ConfigError represents configuration-related errors
	ConfigError struct {
		Op  string
		Err error
	}

	// SessionError represents session-related errors
	SessionError struct {
		SessionName string
		Op          string
		Err         error
	}

	// ValidationError represents input validation errors
	ValidationError struct {
		Field string
		Value string
		Err   error
	}

	// GPTError represents GPT API-related errors
	GPTError struct {
		Op  string
		Err error
	}
)

// Error implementations
func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error in %s: %v", e.Op, e.Err)
}

func (e *SessionError) Error() string {
	if e.SessionName != "" {
		return fmt.Sprintf("session error for '%s' in %s: %v", e.SessionName, e.Op, e.Err)
	}
	return fmt.Sprintf("session error in %s: %v", e.Op, e.Err)
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s '%s': %v", e.Field, e.Value, e.Err)
}

func (e *GPTError) Error() string {
	return fmt.Sprintf("GPT error in %s: %v", e.Op, e.Err)
}

// Unwrap implementations
func (e *ConfigError) Unwrap() error     { return e.Err }
func (e *SessionError) Unwrap() error    { return e.Err }
func (e *ValidationError) Unwrap() error { return e.Err }
func (e *GPTError) Unwrap() error        { return e.Err }

// Constructor functions
func NewConfigError(op string, err error) error {
	return &ConfigError{Op: op, Err: err}
}

func NewSessionError(sessionName, op string, err error) error {
	return &SessionError{SessionName: sessionName, Op: op, Err: err}
}

func NewValidationError(field, value string, err error) error {
	return &ValidationError{Field: field, Value: value, Err: err}
}

func NewGPTError(op string, err error) error {
	return &GPTError{Op: op, Err: err}
}
