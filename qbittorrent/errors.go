package qbittorrent

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError is an error returned by the qBittorrent Web API. It carries the
// HTTP status code and the message returned by the server.
type APIError struct {
	StatusCode int
	Message    string
	// sentinel is the sentinel error this APIError maps to, if any.
	sentinel error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("qbittorrent: status %d: %s", e.StatusCode, e.Message)
}

// Unwrap returns the sentinel error associated with the HTTP status code
func (e *APIError) Unwrap() error {
	return e.sentinel
}

var (
	// ErrUnauthorized is returned when the server responds with HTTP 401.
	ErrUnauthorized = errors.New("qbittorrent: unauthorized")
	// ErrForbidden is returned when the server responds with HTTP 403. This is
	// the status qBittorrent usually returns for authentication failures.
	ErrForbidden = errors.New("qbittorrent: forbidden")

	// ErrNotFound is returned when the server responds with HTTP 404.
	ErrNotFound = errors.New("qbittorrent: not found")

	// ErrInvalidCredentials is returned by Login when the server rejects the
	// provided username and password.
	ErrInvalidCredentials = errors.New("qbittorrent: invalid credentials")

	// ErrInvalidAPIKeyFormat is returned locally, before any request is made,
	// when an API key does not match the expected format (the qbt_ prefix
	// followed by 28 alphanumeric characters).
	ErrInvalidAPIKeyFormat = errors.New("qbittorrent: invalid api key format")
)

func newAPIError(statusCode int, format string, args ...any) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    fmt.Sprintf(format, args...),
		sentinel:   sentinelForStatus(statusCode),
	}
}

func sentinelForStatus(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return nil
	}
}
