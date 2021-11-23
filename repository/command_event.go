package repository

import "time"

type InsertCommand struct {
	Id string
	Entity interface{}
}

type SelectCommand struct {
	Id string
}

type UpdateCommand struct {
	Id string
	Entity interface{}
}

type DeleteCommand struct {
	Id string
}

type InsertedEvent struct {
	Id string `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}

type SelectedEvent struct {
	Id string `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}

type UpdatedEvent struct {
	Id string `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}

type DeletedEvent struct {
	Id string `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}