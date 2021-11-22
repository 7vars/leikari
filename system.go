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

type SystemSettings struct {
	NoSignature bool
}

func newSystemSettings(opts ...Option) SystemSettings {
	settings := SystemSettings{
	}

	for _, opt := range opts {
		switch opt.Name {
		case "noSignature":
			if nos, _ := opt.Bool(); nos {
				settings.NoSignature = true
			}
		}
	}

	return settings
}

type System interface {
	ActorExecutor
	ServiceExecutor
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
	usr ActorHandler
	svc ActorHandler
}

func NewSystem(opts ... Option) System {
	sys := &system{
		settings: newSystemSettings(opts...),
		exitChan: make(chan int, 1),
		log: SysLogger(),
	}

	if !sys.settings.NoSignature {
		fmt.Println(signature)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		sys.Log().Infof("receive signal: %v", sig.String())
		sys.terminate(0)
	}()

	root := newHandler(sys, nil, root(), Log(sys.Log()))
	if err := root.startup(); err != nil {
		panic(err)
	}
	sys.root = root

	usr, err := root.ExecuteHandler(usr())
	if err != nil {
		panic(err)
	}
	sys.usr = usr

	svc, err := root.ExecuteHandler(svc())
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

func (sys *system) Execute(receiver Receiver, opts ...Option) (Ref, error) {
	hdl, err := sys.usr.ExecuteHandler(receiver, opts...)
	if err != nil {
		return nil, err
	}
	return hdl.CreateRef(), nil
}

func (sys *system) ExecuteService(receiver Receiver, opts ...Option) (ActorHandler, error) {
	return sys.svc.ExecuteHandler(receiver, opts...)
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

func root() Actor {
	return Actor{
		Name: "root",
		OnReceive: func(ac ActorContext, m Message) {
			
		},
		OnStart: func(ac ActorContext) error {
			return nil
		},
		OnStop: func(ac ActorContext) error {
			return nil
		},
	}
}

func usr() Actor {
	return Actor{
		Name: "usr",
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
		Name: "svc",
		OnReceive: func(ac ActorContext, m Message) {},
		OnStart: func(ac ActorContext) error {
			return nil
		},
		OnStop: func(ac ActorContext) error {
			return nil
		},
	}
}