package http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/7vars/leikari/route"
	"github.com/gorilla/mux"
)

type request struct {
	req *http.Request
	vars map[string]string

	loaded bool
	body []byte
}

func NewRequest(r *http.Request) route.Request {
	return &request{
		req: r,
		vars: mux.Vars(r),
	}
}

func (r *request) Context() context.Context {
	return r.req.Context()
}

func (r *request) URL() *url.URL {
	return r.req.URL
}

func (r *request) GetHeader(key string) string {
	return r.req.Header.Get(key)
}

func (r *request) GetVar(key string) string {
	if v, ok := r.vars[key]; ok {
		return v
	}
	return ""
}

func (r *request) Body() ([]byte, error) {
	if !r.loaded {
		var err error
		r.body, err = ioutil.ReadAll(r.req.Body)
		if err != nil {
			return nil, err
		}
		defer r.req.Body.Close()
		r.loaded = true
	}
	return r.body, nil
}

func (r *request) Encode(v interface{}) error {
	switch r.GetHeader("Content-Type") {
	case "application/xml":
		return r.Unmarshal(v, xml.Unmarshal)
	// TODO other content-types here
	default:
		return r.Unmarshal(v, json.Unmarshal)
	}
}

func (r *request) Unmarshal(v interface{}, f func([]byte, interface{}) error) error {
	body, err := r.Body()
	if err != nil {
		return err
	}
	return f(body, v)
}