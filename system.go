package leikari

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const signature = `
__/\\\\\\__________________________________________________________________________        
 _\////\\\__________________________/\\\____________________________________________       
  ____\/\\\____________________/\\\_\/\\\_______________________________________/\\\_      
   ____\/\\\________/\\\\\\\\__\///__\/\\\\\\\\_____/\\\\\\\\\_____/\\/\\\\\\\__\///__     
    ____\/\\\______/\\\/////\\\__/\\\_\/\\\////\\\__\////////\\\___\/\\\/////\\\__/\\\_    
     ____\/\\\_____/\\\\\\\\\\\__\/\\\_\/\\\\\\\\/_____/\\\\\\\\\\__\/\\\___\///__\/\\\_   
      ____\/\\\____\//\\///////___\/\\\_\/\\\///\\\____/\\\/////\\\__\/\\\_________\/\\\_  
       __/\\\\\\\\\__\//\\\\\\\\\\_\/\\\_\/\\\_\///\\\_\//\\\\\\\\/\\_\/\\\_________\/\\\_ 
        _\/////////____\//////////__\///__\///____\///___\////////\//__\///__________\///__
`

func NoSignature() Option {
	return Option{
		Name: "noSignature",
		Value: true,
	}
}

type System interface {
	ActorExecutor
	ServiceExecutor
	PubSub
	Settings() SystemSettings
	Log() Logger
	Terminate()
	Terminated() <-chan int
	Run()

	Timer(time.Duration, func(time.Time)) *time.Timer
	Ticker(time.Duration, func(time.Time)) *time.Ticker
}

type system struct {
	settings SystemSettings
	log Logger
	exitChan chan int
	root ActorHandler
	rootRef Ref
	usr ActorHandler
	svc ActorHandler
}

func NewSystem(opts ... Option) System {
	sys := &system{
		settings: newSystemSettings(opts...),
		exitChan: make(chan int, 1),
	}

	sys.log = newLogger(logLevel(sys.settings.GetDefaultString("loglevel", "INFO")))

	if !sys.settings.NoSignature() {
		fmt.Printf("%s\r\n", signature)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		sys.Log().Infof("receive signal: %v", sig.String())
		sys.terminate(0)
	}()

	root := newHandler(sys, nil, root(), "root")
	if err := root.startup(); err != nil {
		panic(err)
	}
	sys.root = root
	sys.rootRef = root.CreateRef()

	usr, err := root.ExecuteHandler(usr(), "usr")
	if err != nil {
		panic(err)
	}
	sys.usr = usr

	svc, err := root.ExecuteHandler(svc(), "svc")
	if err != nil {
		panic(err)
	}
	sys.svc =svc

	return sys
}

func (sys *system) terminate(sig int) {
	go func() {
		sys.root.Close()
		time.Sleep(100 * time.Millisecond)
		sys.exitChan <- sig
	}()
}

func (sys *system) Settings() SystemSettings {
	return sys.settings
}

func (sys *system) Log() Logger {
	return sys.log
}

func (sys *system) Terminate() {
	sys.terminate(0)
}

func (sys *system) Terminated() <-chan int {
	return sys.exitChan
}

func (sys *system) Run() {
	os.Exit(<-sys.Terminated())
}

func (sys *system) At(path string) (Ref, bool) {
	if hdl, ok := sys.root.At(path); ok {
		return hdl.CreateRef(), ok
	}
	return nil, false
}

func (sys *system) Subscribe(ref Ref, f func(interface{}) bool) {
	sys.rootRef.Send(Subscribe{
		Ref: ref,
		Filter: f,
	})
}

func (sys *system) Unsubscribe(ref Ref) {
	sys.rootRef.Send(Unsubscribe{
		Ref: ref,
	})
}

func (sys *system) Publish(v interface{}) {
	sys.rootRef.Send(Publish{v})
}

func (sys *system) Execute(receiver Receiver, name string, opts ...Option) (Ref, error) {
	hdl, err := sys.usr.ExecuteHandler(receiver, name, opts...)
	if err != nil {
		return nil, err
	}
	return hdl.CreateRef(), nil
}

func (sys *system) ExecuteService(receiver Receiver, name string, opts ...Option) (ActorHandler, error) {
	return sys.svc.ExecuteHandler(receiver, name, opts...)
}

func (sys *system) Timer(d time.Duration, f func(time.Time)) *time.Timer {
	timer := time.NewTimer(d)
	go func(tx *time.Timer) {
		t := <-timer.C
		f(t)
	}(timer)
	return timer
}

func (sys *system) Ticker(d time.Duration, f func(time.Time)) *time.Ticker {
	ticker := time.NewTicker(d)
	go func(tx *time.Ticker) {
		for t := range tx.C {
			f(t)
		}
	}(ticker)
	return ticker
}

func usr() Actor {
	return Actor{
		OnReceive: func(ac ActorContext, m Message) {},
		OnStart: func(ac ActorContext) error {
			return nil
		},
		OnStop: func(ac ActorContext) error {
			return nil
		},
	}
}

func svc() Actor {
	return Actor{
		OnReceive: func(ac ActorContext, m Message) {},
		OnStart: func(ac ActorContext) error {
			return nil
		},
		OnStop: func(ac ActorContext) error {
			return nil
		},
	}
}