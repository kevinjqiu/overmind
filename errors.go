package overmind

import "errors"

var (
	// ErrDatabaseNotFound is raised when an expected database is not found
	ErrDatabaseNotFound = errors.New("Database not found")
	// ErrInvalidCommand is raised when an invalid command is issued to a zergling
	ErrInvalidCommand = errors.New("Invalid command")
	// ErrUnknownFacing is raised when a zergling is facing an unknown direction
	ErrUnknownFacing = errors.New("Unknown zergling facing")
)
