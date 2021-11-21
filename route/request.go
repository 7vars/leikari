package route

import (
	"context"
	"net/url"
)

type Request interface {
	Context() context.Context
	URL() *url.URL
	GetHeader(string) string
	GetVar(string) string
	Body() ([]byte, error)
	Encode(interface{}) error
	Unmarshal(interface{}, func([]byte, interface{}) error) error
}
