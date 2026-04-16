package walletpay

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidSignature is returned when webhook signature verification fails.
	ErrInvalidSignature = errors.New("walletpay: invalid webhook signature")
)

// RequestError represents a 400 Bad Request error.
type RequestError struct {
	Code       int
	Message    string
	StatusCode int
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("walletpay: request error (status %d): %s", e.StatusCode, e.Message)
}

// AuthError represents a 401 Unauthorized error (invalid API key).
type AuthError struct {
	Message    string
	StatusCode int
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("walletpay: authentication error: %s", e.Message)
}

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	Message    string
	StatusCode int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("walletpay: not found: %s", e.Message)
}

// RateLimitError represents a 429 Too Many Requests error.
type RateLimitError struct {
	Message    string
	StatusCode int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("walletpay: rate limit exceeded: %s", e.Message)
}

// ServerError represents a 500 Internal Server Error.
type ServerError struct {
	Message    string
	StatusCode int
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("walletpay: server error (status %d): %s", e.StatusCode, e.Message)
}

// APIError represents a generic API error.
type APIError struct {
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("walletpay: api error (status %d): %s", e.StatusCode, e.Message)
}
