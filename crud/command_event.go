package crud

import "time"

type CreateCommand struct {
	Entity interface{}
}

type ReadCommand struct {
	Id string
}

type UpdateCommand struct {
	Id string
	Entity interface{}
}

type DeleteCommand struct {
	Id string
}

type CreatedEvent struct {
	Id string `json:"id"`
	Entity interface{} `json:"entity"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}

type ReadEvent struct {
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
	Timestamp time.Time `json:"timestamp"`
	Took int64 `json:"millis,omitempty"`
}