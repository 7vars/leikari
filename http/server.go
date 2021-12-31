package http

import (
	"context"
	"net/http"
	"time"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/route"
	"github.com/gorilla/mux"
)

const (
	DEFAULT_HTTP_ADDRESS = ":9000"
	DEFAULT_HTTP_READ_TIMEOUT = 5 * time.Second
	DEFAULT_HTTP_WRITE_TIMEOUT = 10 * time.Second
	DEFAULT_HTTP_STOP_TIMEOUT = 5 * time.Second
)

func Address(addr string) leikari.Option {
	return leikari.Option{
		Name: "address",
		Value: addr,
	}
}

func ReadTimeout(t time.Duration) leikari.Option {
	return leikari.Option{
		Name: "readTimeout",
		Value: t,
	}
}

func WriteTimeout(t time.Duration) leikari.Option {
	return leikari.Option{
		Name: "writeTimeout",
		Value: t,
	}
}

func StopTimeout(t time.Duration) leikari.Option {
	return leikari.Option{
		Name: "stopTimeout",
		Value: t,
	}
}

type server struct {
	settings leikari.Settings
	server *http.Server
	def route.Route
}

func newServer(sys leikari.System, route route.Route, opts ...leikari.Option) *server {
	return &server{
		settings: sys.Settings().GetSub("http"),
		def: route,
	}
}

func (srv *server) Receive(ctx leikari.ActorContext, msg leikari.Message) {
	
}

func (srv *server) PreStart(ctx leikari.ActorContext) error {
	router := mux.NewRouter()
	ctx.Log().Debug("preStarting http-server")
	if _, err := ctx.Execute(newRouteActor(router, srv.def), srv.def.RouteName()); err != nil {
		ctx.Log().Error("could not initialize route actor for ", srv.def.RouteName(), err)
		return err
	}

	addr := srv.settings.GetDefaultString("address", DEFAULT_HTTP_ADDRESS)
	srv.server = &http.Server{
		Addr: addr,
		ReadTimeout: srv.settings.GetDefaultDuration("readTimeout", DEFAULT_HTTP_READ_TIMEOUT),
		WriteTimeout: srv.settings.GetDefaultDuration("writeTimeout", DEFAULT_HTTP_WRITE_TIMEOUT),
		Handler: router,
	}
	ctx.Log().Infof("http listen on %s", addr)
	go func(){
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctx.Log().Errorf("failed to listen on %s: %v", addr, err)
		}
	}()
	return nil
}

func (srv *server) PostStop(ctx leikari.ActorContext) error {
	c, cancel := context.WithTimeout(context.Background(), srv.settings.GetDefaultDuration("stopTimeout", DEFAULT_HTTP_STOP_TIMEOUT))
	defer cancel()
	return srv.server.Shutdown(c)
}

func Http(system leikari.System, route route.Route, opts ...leikari.Option) (leikari.ActorHandler, error) {
	if err := route.Validate(); err != nil {
		return nil, err
	}
	return system.ExecuteService(newServer(system, route), "http", opts...)
}
