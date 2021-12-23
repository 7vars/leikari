package crud

import (
	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
	"github.com/7vars/leikari/repository"
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

type Crud interface {
	CreateHandler
	ReadHandler
	UpdateHandler
	DeleteHandler
	repository.QueryHandler
}

func newCrudActor(handler interface{}) leikari.Actor {
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

	qryfunc := func(leikari.ActorContext, query.Query) (*query.QueryResult, error) { return nil, ErrNotFound }
	if qry, ok := handler.(repository.QueryHandler); ok {
		qryfunc = qry.Query
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
		case query.Query:
			return qryfunc(ctx, cmd)
		default:
			return nil, leikari.ErrUnknownCommand
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

	receiver := func(ctx leikari.ActorContext, msg leikari.Message) { msg.Reply(leikari.ErrUnknownCommand) }
	if hdl, ok := handler.(leikari.Receiver); ok {
		receiver = hdl.Receive
	}

	async := true
	if hdl, ok := handler.(leikari.AsyncActor); ok {
		async = hdl.AsyncActor()
	}


	return leikari.Actor{
		OnReceive: func(ctx leikari.ActorContext, msg leikari.Message) {
			result, err := receive(ctx, msg.Value())
			if err != nil {
				if err == leikari.ErrUnknownCommand {
					receiver(ctx, msg)
					return
				}
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
