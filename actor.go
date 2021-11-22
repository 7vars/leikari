package leikari

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
		return MapError("", a.OnStart(ctx))
	}
	return nil
}

func (a Actor) PostStop(ctx ActorContext) error {
	if a.OnStop != nil {
		return MapError("", a.OnStop(ctx))
	}
	return nil
}

func (a Actor) ActorName() string {
	return a.Name
}

func (a Actor) AsyncActor() bool {
	return a.Async
}