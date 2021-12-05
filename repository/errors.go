package repository

import "github.com/7vars/leikari"

var (
	ErrNotFound = leikari.Errorln("", "entity not found").WithStatusCode(404)
	ErrEntityExists = leikari.Errorln("", "entity exists").WithStatusCode(400)
	ErrIdNotPresent = leikari.Errorln("", "id not present").WithStatusCode(400)
	ErrUnknownCommand = leikari.Errorln("", "unknown command")
)