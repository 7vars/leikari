package route

import (
	"errors"
	"strings"
)

type HandleRequest func(Request) Response

type Middleware func(HandleRequest) HandleRequest

type Route struct {
	Name string
	Path string
	Method string
	Handle HandleRequest
	Middleware []Middleware
	Routes []Route
}

func (r Route) RouteName() string {
	if r.Name == "" {
		return strings.ReplaceAll(r.Path, "/", "_")
	}
	return r.Name
}

func (r Route) RouteMiddleware() []Middleware {
	if r.Middleware == nil {
		return make([]Middleware, 0)
	}
	return r.Middleware
}

func (r Route) Validate() error {
	if r.Path == "" {
		return errors.New("path must not be empty")
	}
	return nil
}