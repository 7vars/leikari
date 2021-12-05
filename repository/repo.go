package repository

import (
	"reflect"
	"time"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
)

func RepositoryService(system leikari.ActorExecutor, handler interface{}, name string, opts ...leikari.Option) (RepositoryRef, error) {
	if name == "" {
		return nil, leikari.Errorln("", "name is not defined")
	}
	if handler == nil {
		return nil, leikari.Errorln("", "handler is nil")
	}
	if reflect.TypeOf(handler).Kind() != reflect.Ptr {
		return nil, leikari.Errorln("", "handler must be a pointer")
	}
	
	ref, err := system.Execute(newRepositoryActor(name, handler), opts...)
	if err != nil {
		return nil, err
	}
	
	crudRef := newRepositoryRef(ref)
	return crudRef, nil
}

type QueryFunc func(leikari.ActorContext, query.Query) (*query.QueryResult, error)
type InsertFunc func(leikari.ActorContext, interface{}) (interface{}, error)
type SelectFunc func(leikari.ActorContext, interface{}) (interface{}, error)
type UpdateFunc func(leikari.ActorContext, interface{}, interface{}) error
type DeleteFunc func(leikari.ActorContext, interface{}) (interface{}, error)

type RepositoryHandler struct {
	OnQuery QueryFunc
	OnInsert InsertFunc
	OnSelect SelectFunc
	OnUpdate UpdateFunc
	OnDelete DeleteFunc
	OnStart func(leikari.ActorContext) error
	OnStop func(leikari.ActorContext) error
	Sync bool
}

func (rh *RepositoryHandler) Insert(ctx leikari.ActorContext, cmd InsertCommand) (*InsertedEvent, error) {
	if rh.OnInsert != nil {
		start := time.Now()
		id, err := rh.OnInsert(ctx, cmd.Entity)
		if err != nil {
			return nil, err
		}
		return &InsertedEvent{
			Id: id,
			Entity: cmd.Entity,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (rh *RepositoryHandler) Select(ctx leikari.ActorContext, cmd SelectCommand) (*SelectedEvent, error) {
	if rh.OnSelect != nil {
		start := time.Now()
		entity, err := rh.OnSelect(ctx, cmd.Id)
		if err != nil {
			return nil, err
		}
		return &SelectedEvent{
			Id: cmd.Id,
			Entity: entity,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (rh *RepositoryHandler) Update(ctx leikari.ActorContext, cmd UpdateCommand) (*UpdatedEvent, error) {
	if rh.OnUpdate != nil {
		start := time.Now()
		if err := rh.OnUpdate(ctx, cmd.Id, cmd.Entity); err != nil {
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

func (rh *RepositoryHandler) Delete(ctx leikari.ActorContext, cmd DeleteCommand) (*DeletedEvent, error) {
	if rh.OnDelete != nil {
		start := time.Now()
		entity, err := rh.OnDelete(ctx, cmd.Id)
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

func (rh *RepositoryHandler) Query(ctx leikari.ActorContext, qry query.Query) (*query.QueryResult, error) {
	if rh.OnQuery != nil {
		start := time.Now()
		result, err := rh.OnQuery(ctx, qry)
		if err != nil {
			return nil, err
		}
		result.Timestamp = time.Now()
		result.Took = time.Now().UnixMilli() - start.UnixMilli()
		return result, nil
	}
	return nil, ErrNotFound
}

func (rh *RepositoryHandler) PreStart(ctx leikari.ActorContext) error {
	if rh.OnStart != nil {
		return rh.OnStart(ctx)
	}
	return nil
}

func (rh *RepositoryHandler) PostStop(ctx leikari.ActorContext) error {
	if rh.OnStop != nil {
		return rh.OnStop(ctx)
	}
	return nil
}

func (rh *RepositoryHandler) AsyncActor() bool { return !rh.Sync }