package leikari

import "reflect"

type ActorExecutor interface {
	Execute(Receiver, ...Option) (Ref, error)
}

type ActorContext interface {
	ActorExecutor
	Name() string
	Log() Logger
	Done() <-chan struct{}

	Self() Ref

	Handler() ActorHandler

	Set(string, interface{})
	Add(string, interface{}) error
	Replace(string, interface{}) error
	Get(string) (interface{}, bool)
}

type Receiver interface {
	Receive(ActorContext, Message)
}

type ReceiverFunc func(ActorContext, Message) 

func (f ReceiverFunc) Receive(ctx ActorContext, msg Message) {
	f(ctx, msg)
}

type Startable interface {
	PreStart(ActorContext) error
}

type Stopable interface {
	PostStop(ActorContext) error
}

type NamedActor interface {
	ActorName() string
}

type AsyncActor interface {
	AsyncActor() bool
}

type Actor struct {
	Name string
	OnReceive func(ActorContext, Message)
	OnStart func(ActorContext) error
	OnStop func(ActorContext) error
	Async bool
}

func (a Actor) Receive(ctx ActorContext, msg Message) {
	a.OnReceive(ctx, msg)
}

func (a Actor) PreStart(ctx ActorContext) error {
	if a.OnReceive == nil {
		return Errorln("", "receiver is nil")
	}
	if a.OnStart != nil {
		return a.OnStart(ctx)
	}
	return nil
}

func (a Actor) PostStop(ctx ActorContext) error {
	if a.OnStop != nil {
		return a.OnStop(ctx)
	}
	return nil
}

func (a Actor) ActorName() string {
	return a.Name
}

func (a Actor) AsyncActor() bool {
	return a.Async
}

func NewActor(v interface{}, name string) Actor {
	actor := Actor{
		Name: name,
	}
	if ps, ok := v.(Startable); ok {
		actor.OnStart = ps.PreStart
	}
	if ps, ok := v.(Stopable); ok {
		actor.OnStop = ps.PostStop
	}
	if aa, ok := v.(AsyncActor); ok {
		actor.Async = aa.AsyncActor()
	}
	if rc, ok := v.(Receiver); ok {
		actor.OnReceive = rc.Receive
		return actor
	}

	vact := reflect.ValueOf(v)

	actor.OnReceive = func(ctx ActorContext, msg Message) {
		valv := reflect.ValueOf(msg.Value())
		valt := valv.Type()
		for i := 0; i < vact.NumMethod(); i++ {
			m := vact.Method(i)
			mt := m.Type()

			var result []reflect.Value
			if mt.NumIn() == 2 && CheckIn(mt, ActorContextType, valt) {
				result = m.Call([]reflect.Value{reflect.ValueOf(ctx), valv})
			} else if mt.NumIn() == 1 && CheckIn(mt, valt) {
				result = m.Call([]reflect.Value{valv})
			}

			switch len(result) {
			case 1:
				msg.Reply(result[0].Interface())
				return
			case 2:
				secv := result[1]
				if !secv.IsNil() {
					if err, ok := secv.Interface().(error); ok {
						msg.Reply(err)
						return
					}
				}
				msg.Reply(result[0].Interface())
				return
			}
		}
		msg.Reply(Done())
	}

	return actor
}