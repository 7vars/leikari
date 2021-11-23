package crud

import (
	"github.com/7vars/leikari"
	"github.com/7vars/leikari/repository"
)

var (
	ErrNotCreated = leikari.Errorln("", "entity not created")
	ErrNotFound = repository.ErrNotFound
	ErrNotUpdated = leikari.Errorln("", "entity not updated")
	ErrNotDeleted = leikari.Errorln("", "entity not deleted")
	ErrUnknownCommand = repository.ErrUnknownCommand
)