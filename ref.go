package leikari

import (
	"context"
)

type Ref interface {
	Send(interface{}) error

	RequestChan(interface{}) <-chan interface{}
	Request(interface{}) (interface{}, error)
	RequestContext(context.Context, interface{}) (interface{}, error)
}

type ref struct {
	messages chan<- Message
}

func newRef(messages chan<- Message) Ref {
	return &ref{
		messages: messages,
	}
}

func (r *ref) send(msg Message) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = Errorf("", "message-channel is closed: %v", rec)
		}
	}()
	r.messages <- msg
	return
}

func (r *ref) Send(v interface{}) error {
	return r.send(Send(v))
}

func (r *ref) RequestChan(v interface{}) <-chan interface{} {
	reply := make(chan interface{}, 1)
	go func() {
		if err := r.send(Request(reply, v)); err != nil {
			reply <- err
		}
	}()
	return reply
}

func (r *ref) RequestContext(ctx context.Context, v interface{}) (interface{}, error) {
	select{
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-r.RequestChan(v):
		if err, ok := res.(error); ok {
			return nil, err
		}	
		return res, nil	
	}
}

func (r *ref) Request(v interface{}) (interface{}, error) {
	return r.RequestContext(context.Background(), v)
}