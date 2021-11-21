package crud

import "errors"

var (
	ErrNotCreated = errors.New("entity not created")
	ErrNotFound = errors.New("entity not found")
	ErrNotUpdated = errors.New("entity not updated")
	ErrNotDeleted = errors.New("entity not deleted")
	ErrUnknownCommand = errors.New(("unknown command"))
)