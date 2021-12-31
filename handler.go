package leikari

import (
	"fmt"
	"strings"
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

func (ctx *actorContext) System() System {
	return ctx.Handler().System()
}

func (ctx *actorContext) Log() Logger {
	return ctx.log
}

func (ctx *actorContext) Settings() Settings {
	return ctx.handler.Settings()
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

func (ctx *actorContext) At(path string) (Ref, bool) {
	if hdl, ok := ctx.handler.At(path); ok {
		return hdl.CreateRef(), ok
	}
	return nil, false
}

func (ctx *actorContext) Subscribe(ref Ref, f func(interface{}) bool) {
	ctx.System().Subscribe(ref, f)
}

func (ctx *actorContext) Unsubscribe(ref Ref) {
	ctx.System().Unsubscribe(ref)
}

func (ctx *actorContext) Publish(v interface{}) {
	ctx.System().Publish(v)
}

func (ctx *actorContext) Execute(receiver Receiver, name string, opts ...Option) (Ref, error) {
	hdl, err := ctx.Handler().ExecuteHandler(receiver, name, opts...)
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
	At(string) (ActorHandler, bool)
	ExecuteHandler(Receiver, string, ...Option) (ActorHandler, error)
}

type ActorHandler interface {
	ActorHandlerExecutor
	Name() string
	Close()

	Root() ActorHandler
	Parent() (ActorHandler, bool)
	Path() string
	System() System
	Log() Logger
	Settings() Settings

	Child(string) (ActorHandler, bool)
	Children() []ActorHandler

	CreateRef() Ref
	
	Cache() Cache
}

type handler struct {
	sync.RWMutex
	name string
	settings ActorSettings
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

func newHandler(system System, parent ActorHandler, receiver Receiver, name string, options ...Option) *handler {
	if actor, ok := receiver.(AsyncActor); ok && actor.AsyncActor() {
		options = append(options, Async())
	}
	settings := system.Settings().GetActorSettings(name, options...)

	hdl :=  &handler{
		name: name,
		settings: settings,
		messages: make(chan Message, settings.MessageQueueSize()),
		receiver: receiver,
		system: system,
		parent: parent,
		children: make(map[string]ActorHandler),
		cache: NewCache(),
	}

	var log Logger
	if parent != nil {
		log = parent.Log().ForName(hdl.Path())
	} else {
		log = system.Log().ForName(hdl.Path())
	}
	hdl.log = log

	log.Debug("actor", hdl.name, "with", "message-queue-size:", settings.MessageQueueSize(), ", worker-pool:", settings.WorkerPoolSize(), "created")

	return hdl
}

func (hdl *handler) createContext(name string, log Logger) ActorContext {
	hdl.Lock()
	defer hdl.Unlock()

	ctx := &actorContext{
		name: hdl.Name(),
		log: log,
		handler: hdl,
		self: hdl.CreateRef(),
		done: make(chan struct{}),
	}

	hdl.contextes = append(hdl.contextes, ctx)
	return ctx
}

func (hdl *handler) startup() error {
	pool := hdl.settings.WorkerPoolSize()
	for i := 0; i < pool; i++ {
		path := hdl.Path()
		log := hdl.Log()
		if pool > 1 {
			path = fmt.Sprintf("%s-%d", path, i)
			log = log.ForName(path)
		}

		ctx := hdl.createContext(hdl.Name(), log)

		if starter, ok := hdl.receiver.(Startable); ok {
			if err := starter.PreStart(ctx); err != nil {
				hdl.Close()
				return err
			}
		}
		go worker(ctx, hdl.messages, hdl.receiver, hdl.settings.Async())
	}
	return nil
}

func (hdl *handler) Name() string {
	return hdl.name
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

func (hdl *handler) Path() string {
	p, ok := hdl.Parent()
	if !ok {
		return "/"
	}
	ppath := p.Path()
	if strings.HasSuffix(ppath, "/") {
		return ppath + hdl.Name()
	}
	return ppath + "/" + hdl.Name()
}

func (hdl *handler) System() System {
	return hdl.system
}

func (hdl *handler) Log() Logger {
	return hdl.log
}

func (hdl *handler) Settings() Settings {
	return hdl.settings
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
	return newRef(hdl.messages)
}

func (hdl *handler) At(path string) (ActorHandler, bool) {
	if len(path) > 0 {
		if path[0] == '/' {
			if len(path) == 1 {
				return hdl.Root(), true
			}
			return hdl.Root().At(path[1:])
		} else if len(path) >= 2 && path[0:2] == ".." && hdl.parent != nil {
			if len(path) == 2 || path == "../" {
				return hdl.parent, true
			}
			return hdl.parent.At(path[3:])
		} else if path[0] == '.' || strings.HasPrefix(path, hdl.name) {
			if i := strings.IndexRune(path, '/'); i > -1 {
				if i < len(path)-1 {
					return hdl.At(path[i+1:])
				}
			} else {
				return hdl, true
			}
		} else {
			if i := strings.IndexRune(path, '/'); i != -1 {
				child, ok := hdl.Child(path[:i])
				if !ok {
					return nil, false
				}
				if i < len(path)-1 {
					return child.At(path[i+1:])
				}
				return child, true
			}
			return hdl.Child(path)
		}
	}
	return nil, false
}

func (hdl *handler) Cache() Cache {
	return hdl.cache
}

func (hdl *handler) ExecuteHandler(receiver Receiver, name string, opts ...Option) (ActorHandler, error) {
	hdl.Lock()
	defer hdl.Unlock()

	child := newHandler(hdl.System(), hdl, receiver, name, opts...)
	
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