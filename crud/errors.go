package crud

import (
	"github.com/7vars/leikari"
)

var (
	ErrNotCreated = leikari.Errorln("", "entity not created")
	ErrNotFound = leikari.Errorln("", "entity not found").WithStatusCode(404)
	ErrNotUpdated = leikari.Errorln("", "entity not updated")
	ErrNotDeleted = leikari.Errorln("", "entity not deleted")
	ErrUnknownCommand = leikari.Errorln("", "unknown command")
)