package leikari

import "context"

type Ref interface {
	Send(interface{}) error

	RequestChan(interface{}) <-chan interface{}
	Request(interface{}) (interface{}, error)
	RequestContext(context.Context, interface{}) (interface{}, error)
}

type ref struct {
	pusher Pusher
}

func NewRef(pusher Pusher) Ref {
	return &ref{
		pusher: pusher,
	}
}

func (r *ref) Send(v interface{}) error {
	return r.pusher.Push(Send(v))
}

func (r *ref) RequestChan(v interface{}) <-chan interface{} {
	reply := make(chan interface{})
	if err := r.pusher.Push(Request(reply, v)); err != nil {
		reply <- err
	}
	return reply
}

func (r *ref) RequestContext(ctx context.Context, v interface{}) (interface{}, error) {
	select{
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-r.RequestChan(v):
		if err, ok := res.(error); ok {
			return nil, MapError("", err)
		}	
		return res, nil	
	}
}

func (r *ref) Request(v interface{}) (interface{}, error) {
	return r.RequestContext(context.Background(), v)
}