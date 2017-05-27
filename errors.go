package overmind

import "errors"

var (
	// ErrDatabaseNotFound is raised when an expected database is not found
	ErrDatabaseNotFound = errors.New("Database not found")
)
