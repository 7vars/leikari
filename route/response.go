package route

import (
	"encoding/json"
	"encoding/xml"
)

type Response struct {
	Header map[string]string
	Status int
	Data interface{}
}

func (r Response) SetHeader(key, value string) {
	r.Header[key] = value
}

func (r Response) GetHeader(key string) string {
	if len(r.Header) > 0 {
		return r.Header[key]
	}
	return ""
}

func (r Response) StatusCode() int {
	if r.Status == 0 {
		return 200
	}
	return r.Status
}

func (r Response) ContentType() string {
	if ct := r.GetHeader("Content-Type"); ct != "" {
		return ct
	}
	return "application/json"
}

func (r Response) Decode() ([]byte, error) {
	switch r.ContentType() {
	case "application/xml":
		return r.Marshal(xml.Marshal)
	// TODO other encodings here
	default:
		return r.Marshal(json.Marshal)
	}
}

func (r Response) Marshal(f func(interface{}) ([]byte, error)) ([]byte, error) {
	return f(r.Data)
}

func ErrorResponse(err error) Response {
	// TODO leikari specific errors
	return ErrorResponseWithStatus(500, err)
}

func ErrorResponseWithStatus(status int, err error) Response {
	// TODO leikari specific errors
	return Response{
		Status: status,
		Data: map[string]interface{}{
			"error": err.Error(),
		},
	}
}