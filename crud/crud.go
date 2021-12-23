package crud

import (
	"reflect"
	"time"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
	"github.com/7vars/leikari/route"
)

func CrudService(system leikari.ActorExecutor, handler interface{}, name string, opts ...leikari.Option) (CrudRef, route.Route, error) {
	if name == "" {
		return nil, route.Route{}, leikari.Errorln("", "name is not defined")
	}
	if handler == nil {
		return nil, route.Route{}, leikari.Errorln("", "handler is nil")
	}
	if reflect.TypeOf(handler).Kind() != reflect.Ptr {
		return nil, route.Route{}, leikari.Errorln("", "handler must be a pointer")
	}
	
	ref, err := system.Execute(newCrudActor(handler), name, opts...)
	if err != nil {
		return nil, route.Route{}, err
	}
	
	crudRef := newCrudRef(ref)
	return crudRef, newCrudRoute(name, crudRef, handler), nil
}

type QueryFunc func(leikari.ActorContext, query.Query) (*query.QueryResult, error)
type CreateFunc func(leikari.ActorContext, interface{}) (string, interface{}, error)
type ReadFunc func(leikari.ActorContext, string) (interface{}, error)
type UpdateFunc func(leikari.ActorContext, string, interface{}) error
type DeleteFunc func(leikari.ActorContext, string) (interface{}, error)

type CrudHandler struct {
	OnCreate CreateFunc
	OnRead ReadFunc
	OnUpdate UpdateFunc
	OnDelete DeleteFunc
	OnQuery QueryFunc
	OnReceive func(leikari.ActorContext, leikari.Message)
	OnStart func(leikari.ActorContext) error
	OnStop func(leikari.ActorContext) error
	OnUnmarshal func([]byte) (interface{}, error)
	Sync bool
}

func (a *CrudHandler) Create(ctx leikari.ActorContext, cmd CreateCommand) (*CreatedEvent, error) {
	if a.OnCreate != nil {
		start := time.Now()
		id, entity, err := a.OnCreate(ctx, cmd.Entity)
		if err != nil {
			return nil, err
		}
		return &CreatedEvent{
			Id: id,
			Entity: entity,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (a *CrudHandler) Read(ctx leikari.ActorContext, cmd ReadCommand) (*ReadEvent, error) {
	if a.OnRead != nil {
		start := time.Now()
		entity, err := a.OnRead(ctx, cmd.Id)
		if err != nil {
			return nil, err
		}
		return &ReadEvent{
			Id: cmd.Id,
			Entity: entity,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (a *CrudHandler) Update(ctx leikari.ActorContext, cmd UpdateCommand) (*UpdatedEvent, error) {
	if a.OnUpdate != nil {
		start := time.Now()
		if err := a.OnUpdate(ctx, cmd.Id, cmd.Entity); err != nil {
			return nil, err
		}
		return &UpdatedEvent{
			Id: cmd.Id,
			Entity: cmd.Entity,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (a *CrudHandler) Delete(ctx leikari.ActorContext, cmd DeleteCommand) (*DeletedEvent, error) {
	if a.OnDelete != nil {
		start := time.Now()
		entity, err := a.OnDelete(ctx, cmd.Id)
		if err != nil {
			return nil, err
		}
		return &DeletedEvent{
			Id: cmd.Id,
			Entity: entity,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (a *CrudHandler) Query(ctx leikari.ActorContext, qry query.Query) (*query.QueryResult, error) {
	if a.OnQuery != nil {
		result, err := a.OnQuery(ctx, qry)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, ErrNotFound
}

func (a *CrudHandler) Receive(ctx leikari.ActorContext, msg leikari.Message) {
	if a.OnReceive != nil {
		a.OnReceive(ctx, msg)
		return
	}
	msg.Reply(leikari.ErrUnknownCommand)
}

func (a *CrudHandler) PreStart(ctx leikari.ActorContext) error {
	if a.OnStart != nil {
		return a.OnStart(ctx)
	}
	return nil
}

func (a *CrudHandler) PostStop(ctx leikari.ActorContext) error {
	if a.OnStop != nil {
		return a.OnStop(ctx)
	}
	return nil
}

func (a *CrudHandler) HandleUnmarshal(data []byte) (interface{}, error) {
	if a.OnUnmarshal != nil {
		return a.OnUnmarshal(data)
	}
	return HandleUnmarshal(data)
}

func (a *CrudHandler) AsyncActor() bool { return !a.Sync }