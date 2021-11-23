package repository

import "github.com/7vars/leikari"

var (
	ErrNotFound = leikari.Errorln("", "entity not found").WithStatusCode(404)
	ErrUnknownCommand = leikari.Errorln("", "unknown command")
)