package leikari

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type actorContext struct {
	name string
	log Logger
	handler ActorHandler
	self Ref
	done chan struct{}
}

func (ctx *actorContext) Name() string {
	return ctx.name
}

func (ctx *actorContext) Log() Logger {
	return ctx.log
}

func (ctx *actorContext) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *actorContext) Self() Ref {
	return ctx.self
}

func (ctx *actorContext) Handler() ActorHandler {
	return ctx.handler
}

func (ctx *actorContext) Execute(receiver Receiver, opts ...Option) (Ref, error) {
	hdl, err := ctx.Handler().ExecuteHandler(receiver, opts...)
	if err != nil {
		return nil, err
	}
	return hdl.CreateRef(), nil
}

func (ctx *actorContext) terminate() {
	ctx.done <- struct{}{}
}

func (ctx *actorContext) Set(key string, value interface{}) {
	ctx.handler.Cache().Set(key, value)
}

func (ctx *actorContext) Add(key string, value interface{}) error {
	return ctx.handler.Cache().Add(key, value)
}

func (ctx *actorContext) Replace(key string, value interface{}) error {
	return ctx.handler.Cache().Replace(key, value)
}

func (ctx *actorContext) Get(key string) (interface{}, bool) {
	return ctx.handler.Cache().Get(key)
}

func Name(name string) Option {
	return Option{
		Name: "name",
		Value: name,
	}
}

func Log(log Logger) Option {
	return Option{
		Name: "logger",
		Value: log,
	}
}

func WorkerPool(size int) Option {
	return Option{
		Name: "workerPool",
		Value: size,
	}
}

func MessageQueue(size int) Option {
	return Option{
		Name: "messageQueue",
		Value: size,
	}
}

func Async() Option {
	return Option{
		Name: "async",
		Value: true,
	}
}

type HandlerSettings struct {
	Name string
	WorkerPool int
	MessageQueue int
	Log Logger
	Async bool
}

func newHandlerSettings(opts ...Option) HandlerSettings {
	settings := HandlerSettings{
		Name: GenerateName(),
		Log: SysLogger(),
		WorkerPool: 1,
		MessageQueue: 1000,
	}

	for _, opt := range opts {
		switch opt.Name {
		case "name":
			if nm := opt.String(); nm != "" {
				settings.Name = nm
			}
		case "logger":
			if log, ok := opt.Value.(Logger); ok {
				settings.Log = log
			}
		case "workerPool":
			if wp, _ := opt.Int(); wp > 0 {
				settings.WorkerPool = wp
			}
		case "messageQueue":
			if mq, _ := opt.Int(); mq > 0 {
				settings.MessageQueue = mq
			}
		case "async":
			if asb, _ := opt.Bool(); asb {
				settings.Async = true
			}
		}
	}

	return settings
}

func worker(ctx ActorContext, jobs <-chan Message, r Receiver, async bool) {
	defer func(){
		if stop, ok := r.(Stopable); ok {
			if err := stop.PostStop(ctx); err != nil {
				ctx.Log().Error(err)
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			ctx.Log().Debug("worker queue stopped")
			return 
		case msg := <-jobs:
			if async {
				go r.Receive(ctx, msg)
			} else {
				r.Receive(ctx, msg)
			}
		}
	}
}

type ActorHandlerExecutor interface {
	ExecuteHandler(Receiver, ...Option) (ActorHandler, error)
}

type ActorHandler interface {
	Pusher
	ActorHandlerExecutor
	Name() string
	Close()

	Root() ActorHandler
	Parent() (ActorHandler, bool)
	System() System
	Log() Logger

	Child(string) (ActorHandler, bool)
	Children() []ActorHandler

	CreateRef() Ref
	
	Cache() Cache
}

type handler struct {
	sync.RWMutex
	settings HandlerSettings
	messages chan Message
	receiver Receiver
	system System
	parent ActorHandler
	log Logger
	closed bool
	contextes []*actorContext

	children map[string]ActorHandler

	cache Cache
}

func newHandler(system System, parent ActorHandler, receiver Receiver, options ...Option) *handler {
	if actor, ok := receiver.(NamedActor); ok {
		options = append(options, Name(actor.ActorName()))
	}
	if actor, ok := receiver.(AsyncActor); ok && actor.AsyncActor() {
		options = append(options, Async())
	}
	settings := newHandlerSettings(options...)
	return &handler{
		settings: settings,
		messages: make(chan Message, settings.MessageQueue),
		receiver: receiver,
		system: system,
		parent: parent,
		log: settings.Log.ForName(settings.Name),
		children: make(map[string]ActorHandler),
		cache: NewCache(),
	}
}

func (hdl *handler) createContext(name string, log Logger) ActorContext {
	hdl.Lock()
	defer hdl.Unlock()

	ctx := &actorContext{
		name: hdl.Name(),
		log: hdl.Log(),
		handler: hdl,
		self: hdl.CreateRef(),
		done: make(chan struct{}),
	}

	hdl.contextes = append(hdl.contextes, ctx)
	return ctx
}

func (hdl *handler) startup() error {
	pool := hdl.settings.WorkerPool
	for i := 0; i < pool; i++ {
		name := hdl.Name()
		log := hdl.Log()
		if pool > 1 {
			name = fmt.Sprintf("%s-%d", name, i)
			log = log.ForName(name)
		}

		ctx := hdl.createContext(name, log)

		if starter, ok := hdl.receiver.(Startable); ok {
			if err := starter.PreStart(ctx); err != nil {
				hdl.Close()
				return err
			}
		}
		go worker(ctx, hdl.messages, hdl.receiver, hdl.settings.Async)
	}
	return nil
}

func (hdl *handler) Name() string {
	return hdl.settings.Name
}

func (hdl *handler) Push(msg Message) error {
	hdl.RLock()
	defer hdl.RUnlock()

	if hdl.closed {
		return io.EOF
	}
	hdl.messages <- msg
	return nil
}

func (hdl *handler) Close() {
	hdl.Lock()
	defer hdl.Unlock()
	hdl.closed = true

	var wg sync.WaitGroup

	for _, child := range hdl.children {
		wg.Add(1)
		
		go func(c ActorHandler){
			defer wg.Done()
			c.Close()
		}(child)
	}

	for _, ctx := range hdl.contextes {
		wg.Add(1)

		go func(c *actorContext){
			defer wg.Done()
			c.terminate()
		}(ctx)
	}
	
	if err := waitTimeout(&wg, 10 * time.Second); err != nil { // TODO configure
		hdl.Log().Warnf("could not close successfully: %v", err)
	}

	close(hdl.messages)
}

func (hdl *handler) Root() ActorHandler {
	if par, ok := hdl.Parent(); ok {
		return par.Root()
	}
	return hdl
}

func (hdl *handler) Parent() (ActorHandler, bool) {
	if hdl.parent != nil {
		return hdl.parent, true
	}
	return nil, false
}

func (hdl *handler) System() System {
	return hdl.system
}

func (hdl *handler) Log() Logger {
	return hdl.log
}

func (hdl *handler) Child(name string) (ActorHandler, bool) {
	hdl.RLock()
	defer hdl.RUnlock()
	child, ok := hdl.children[name]
	return child, ok
}

func (hdl *handler) Children() []ActorHandler {
	hdl.RLock()
	defer hdl.RUnlock()
	var children []ActorHandler
	for name := range hdl.children {
		children = append(children, hdl.children[name])
	}
	return children
}

func (hdl *handler) CreateRef() Ref {
	return NewRef(hdl)
}

func (hdl *handler) Cache() Cache {
	return hdl.cache
}

func (hdl *handler) ExecuteHandler(receiver Receiver, opts ...Option) (ActorHandler, error) {
	hdl.Lock()
	defer hdl.Unlock()

	child := newHandler(hdl.System(), hdl, receiver, append(opts, Log(hdl.log))...)
	
	if _, exists := hdl.children[child.Name()]; exists {
		child.Close()
		return nil, Errorf("", "child %v already exists", child.Name())
	}

	if err := child.startup(); err != nil {
		return nil, err
	}

	hdl.children[child.Name()] = child
	return child, nil
}