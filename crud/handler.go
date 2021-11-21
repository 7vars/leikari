package crud

import (
	"github.com/7vars/leikari"
)

type CreateHandler interface {
	Create(leikari.ActorContext, CreateCommand) (*CreatedEvent, error)
}

type ReadHandler interface {
	Read(leikari.ActorContext, ReadCommand) (*ReadEvent, error)
}

type UpdateHandler interface {
	Update(leikari.ActorContext, UpdateCommand) (*UpdatedEvent, error)
}

type DeleteHandler interface {
	Delete(leikari.ActorContext, DeleteCommand) (*DeletedEvent, error)
}

type ListHandler interface {
	List(leikari.ActorContext, Query) (*QueryResult, error)
}

func newCrudActor(name string, handler interface{}) leikari.Actor {
	crfunc := func(leikari.ActorContext,CreateCommand) (*CreatedEvent, error) { return nil, ErrNotFound }
	if cr, ok := handler.(CreateHandler); ok {
		crfunc = cr.Create
	}

	refunc := func(leikari.ActorContext, ReadCommand) (*ReadEvent, error) { return nil, ErrNotFound }
	if re, ok := handler.(ReadHandler); ok {
		refunc = re.Read
	}

	upfunc := func(leikari.ActorContext, UpdateCommand) (*UpdatedEvent, error) { return nil, ErrNotFound }
	if up, ok := handler.(UpdateHandler); ok {
		upfunc = up.Update
	}

	delfunc := func(leikari.ActorContext, DeleteCommand) (*DeletedEvent, error) { return nil, ErrNotFound }
	if del, ok := handler.(DeleteHandler); ok {
		delfunc = del.Delete
	}

	qryfunc := func(leikari.ActorContext, Query) (*QueryResult, error) { return nil, ErrNotFound }
	if qry, ok := handler.(ListHandler); ok {
		qryfunc = qry.List
	}

	receive := func(ctx leikari.ActorContext, v interface{}) (interface{}, error) {
		switch cmd := v.(type) {
		case CreateCommand:
			return crfunc(ctx, cmd)
		case ReadCommand:
			return refunc(ctx, cmd)
		case UpdateCommand:
			return upfunc(ctx, cmd)
		case DeleteCommand:
			return delfunc(ctx, cmd)
		case Query:
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
