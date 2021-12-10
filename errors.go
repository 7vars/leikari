package leikari

import "fmt"

var (
	ErrUnknownCommand = Errorln("", "unknown command")
)

type Error struct {
	Code string `json:"code,omitempty"`
	Message string `json:"error"`
	Description string `json:"description,omitempty"`
	Status int `json:"-"`
}

func NewOf(code string, err error) *Error {
	return &Error{
		Code: code,
		Message: err.Error(),
	}
}

func New(code string, msg string) *Error {
	return &Error{
		Code: code,
		Message: msg,
	}
}

func MapError(code string, err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		if e.Code == "" {
			e.Code = code
		}
		return e
	}
	return NewOf(code, err)
}

func Errorln(code string, v ...interface{}) *Error {
	msg := []rune(fmt.Sprintln(v...))
	if len(msg) > 0 {
		return New(code, string(msg[:len(msg)-1]))
	}
	return New(code, "")
}

func Errorf(code string, format string, v ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, v...))
}

func (e *Error) Error() string {
	if e.Code == "" {
		return e.Message
	}
	return fmt.Sprint(e.Code, " - ", e.Message)
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) WithDescription(desc string) *Error {
	e.Description = desc
	return e
}

func (e *Error) WithStatusCode(status int) *Error {
	e.Status = status
	return e
}

func (e *Error) StatusCode() int {
	if e.Status == 0 {
		return 500
	}
	return e.Status
}
