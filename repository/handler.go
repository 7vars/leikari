package repository

import (
	"reflect"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
)

type InsertHandler interface {
	Insert(leikari.ActorContext, InsertCommand) (*InsertedEvent, error)
}

type SelectHandler interface {
	Select(leikari.ActorContext, SelectCommand) (*SelectedEvent, error)
}

type UpdateHandler interface {
	Update(leikari.ActorContext, UpdateCommand) (*UpdatedEvent, error)
}

type DeleteHandler interface {
	Delete(leikari.ActorContext, DeleteCommand) (*DeletedEvent, error)
}

type QueryHandler interface {
	Query(leikari.ActorContext, query.Query) (*query.QueryResult, error)
}

var (
	insertHandlerType = reflect.TypeOf((*InsertHandler)(nil)).Elem()
	updateHandlerType = reflect.TypeOf((*UpdateHandler)(nil)).Elem()
	selectHandlerType = reflect.TypeOf((*SelectHandler)(nil)).Elem()
	deleteHandlerType = reflect.TypeOf((*DeleteHandler)(nil)).Elem()
)

func newRepositoryActor(name string, v interface{}) leikari.Actor {
	handler := v
	if !leikari.CheckImplementsOneOf(reflect.TypeOf(v), insertHandlerType, updateHandlerType, selectHandlerType, deleteHandlerType) {
		handler = wrap(v)
	}

	insfunc := func(leikari.ActorContext, InsertCommand) (*InsertedEvent, error) { return nil, ErrNotFound }
	if hdl, ok := handler.(InsertHandler); ok {
		insfunc = hdl.Insert
	}

	selfunc := func(leikari.ActorContext, SelectCommand) (*SelectedEvent, error) { return nil, ErrNotFound }
	if hdl, ok := handler.(SelectHandler); ok {
		selfunc = hdl.Select
	}

	uptfunc := func(leikari.ActorContext, UpdateCommand) (*UpdatedEvent, error) { return nil, ErrNotFound }
	if hdl, ok := handler.(UpdateHandler); ok {
		uptfunc = hdl.Update
	}

	delfunc := func(leikari.ActorContext, DeleteCommand) (*DeletedEvent, error) { return nil, ErrNotFound }
	if hdl, ok := handler.(DeleteHandler); ok {
		delfunc = hdl.Delete
	}

	qryfunc := func(leikari.ActorContext, query.Query) (*query.QueryResult, error) { return nil, ErrNotFound }
	if hdl, ok := handler.(QueryHandler); ok {
		qryfunc = hdl.Query
	}

	receive := func(ctx leikari.ActorContext, v interface{}) (interface{}, error) {
		switch cmd := v.(type) {
		case InsertCommand:
			return insfunc(ctx, cmd)
		case SelectCommand:
			return selfunc(ctx, cmd)
		case UpdateCommand:
			return uptfunc(ctx, cmd)
		case DeleteCommand:
			return delfunc(ctx, cmd)
		case query.Query:
			return qryfunc(ctx, cmd)
		default:
			return nil, ErrUnknownCommand
		}
	}

	start := func(leikari.ActorContext) error { return nil }
	if hdl, ok := handler.(leikari.Startable); ok {
		start = hdl.PreStart
	}

	stop := func(leikari.ActorContext) error { return nil }
	if hdl, ok := handler.(leikari.Stopable); ok {
		stop = hdl.PostStop
	}

	async := true
	if hdl, ok := handler.(leikari.AsyncActor); ok {
		async = hdl.AsyncActor()
	}

	return leikari.Actor{
		Name: name,
		OnReceive: func(ctx leikari.ActorContext, msg leikari.Message) {
			result, err := receive(ctx, msg.Value())
			if err != nil {
				msg.Reply(err)
				return
			}
			msg.Reply(result)
		},
		OnStart: start,
		OnStop: stop,
		Async: async,
	}
}
