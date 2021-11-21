package http

import (
	"context"
	"net/http"
	"time"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/route"
	"github.com/gorilla/mux"
)

func Address(addr string) leikari.Option {
	return leikari.Option{
		Name: "address",
		Value: addr,
	}
}

func ReadTimeout(t int) leikari.Option {
	return leikari.Option{
		Name: "readTimeout",
		Value: t,
	}
}

func WriteTimeout(t int) leikari.Option {
	return leikari.Option{
		Name: "writeTimeout",
		Value: t,
	}
}

func StopTimeout(t int) leikari.Option {
	return leikari.Option{
		Name: "stopTimeout",
		Value: t,
	}
}

type HttpServerSettings struct {
	Address string
	ReadTimeout int
	WriteTimeout int
	StopTimeout int
}

func newHttpServerSettings(opts ...leikari.Option) HttpServerSettings {
	settings := HttpServerSettings{
		Address: ":9000",
		ReadTimeout: 5,
		WriteTimeout: 10,
		StopTimeout: 5,
	}

	for _, opt := range opts {
		switch opt.Name {
		case "address":
			if addr := opt.String(); addr != "" {
				settings.Address = addr
			}
		case "readTimeout":
			if t, _ := opt.Int(); t > 0 {
				settings.ReadTimeout = t
			}
		case "writeTimeout":
			if t, _ := opt.Int(); t > 0 {
				settings.WriteTimeout = t
			}
		case "stopTimeout":
			if t, _ := opt.Int(); t > 0 {
				settings.StopTimeout = t
			}
		}
	}

	return settings
}

type server struct {
	settings HttpServerSettings
	server *http.Server
	def route.Route
}

func newServer(route route.Route, opts ...leikari.Option) *server {
	return &server{
		settings: newHttpServerSettings(opts...),
		def: route,
	}
}

func (srv *server) Receive(ctx leikari.ActorContext, msg leikari.Message) {
	
}

func (srv *server) ActorName() string {
	return "http"
}

func (srv *server) PreStart(ctx leikari.ActorContext) error {
	router := mux.NewRouter()
	ctx.Log().Debug("preStarting http-server")
	if _, err := ctx.Execute(newRouteActor(router, srv.def)); err != nil {
		ctx.Log().Error("could not initialize route actor for ", srv.def.RouteName(), err)
		return err
	}

	srv.server = &http.Server{
		Addr: srv.settings.Address,
		ReadTimeout: time.Duration(srv.settings.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(srv.settings.WriteTimeout) * time.Second,
		Handler: router,
	}
	ctx.Log().Infof("http listen on %s", srv.settings.Address)
	go func(){
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctx.Log().Errorf("failed to listen on %s: %v", srv.settings.Address, err)
		}
	}()
	return nil
}

func (srv *server) PostStop(ctx leikari.ActorContext) error {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(srv.settings.StopTimeout) * time.Second)
	defer cancel()
	return srv.server.Shutdown(c)
}

func HttpServer(system leikari.System, route route.Route, opts ...leikari.Option) (leikari.ActorHandler, error) {
	if err := route.Validate(); err != nil {
		return nil, err
	}
	server := newServer(route, opts...)
	return system.ExecuteService(server, opts...)
}
