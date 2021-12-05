package repository

import "time"

type InsertCommand struct {
	Entity interface{}
}

type SelectCommand struct {
	Id interface{}
}

type UpdateCommand struct {
	Id interface{}
	Entity interface{}
}

type DeleteCommand struct {
	Id interface{}
}

type InsertedEvent struct {
	Id interface{} `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}

type SelectedEvent struct {
	Id interface{} `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}

type UpdatedEvent struct {
	Id interface{} `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}

type DeletedEvent struct {
	Id interface{} `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}