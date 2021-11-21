package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/route"
	"github.com/gorilla/mux"
)

func httpHandlerFunc(ref leikari.Ref, log leikari.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(r)
		res, err := ref.RequestContext(r.Context(), req)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error")
			return
		}
		if response, ok := res.(route.Response); ok {
			if len(response.Header) > 0 {
				for key, value := range response.Header {
					w.Header().Add(key, value)
				}
			}

			w.Header().Add("Content-Type", response.ContentType())

			w.WriteHeader(response.StatusCode())
			buf, err := response.Decode()
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "internal server error")
				return 
			}
			w.Write(buf)
			return
		}
		log.Error("no response received")
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, "not implemented")
	}
}

type routeActor struct {
	router *mux.Router
	def route.Route
	middleware []route.Middleware
}

func newRouteActor(router *mux.Router, def route.Route, middleware ...route.Middleware) leikari.Receiver {
	return &routeActor{
		router: router,
		def: def,
		middleware: middleware,
	}
}

func (ra *routeActor) ActorName() string {
	return ra.def.RouteName()
}

func (ra *routeActor) middlewares() []route.Middleware {
	return append(ra.middleware, ra.def.RouteMiddleware()...)
}
 
func (ra *routeActor) PreStart(ctx leikari.ActorContext) error {
	if ra.def.Handle != nil {
		method := "GET"
		if ra.def.Method != "" {
			method = ra.def.Method
		}
		ra.router.HandleFunc(ra.def.Path, httpHandlerFunc(ctx.Self(), ctx.Log())).Methods(method)
	}
	
	if len(ra.def.Routes) > 0 {
		subrouter := ra.router.PathPrefix(ra.def.Path).Subrouter()
		for _, childRoute := range ra.def.Routes {
			if _, err := ctx.Execute(newRouteActor(subrouter, childRoute, ra.middlewares()...)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ra *routeActor) Receive(ctx leikari.ActorContext, msg leikari.Message) {
	if ra.def.Handle == nil {
		msg.Reply(errors.New("route handler not defined"))
		return
	}
	if request, ok := msg.Value().(route.Request); ok {
		handle := ra.def.Handle
		for _, mw := range ra.middlewares() {
			handle = mw(handle)
		}

		msg.Reply(handle(request))
		return
	}
	msg.Reply(fmt.Errorf("unkonwn type %T for Request", msg.Value()))
}

func (ra *routeActor) AsyncActor() bool {
	return true
}