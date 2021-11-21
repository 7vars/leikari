package crud

import (
	"time"

	"github.com/7vars/leikari"
)

type ListFunc func(leikari.ActorContext, Query) ([]interface{}, int, error)
type CreateFunc func(leikari.ActorContext, interface{}) (string, interface{}, error)
type ReadFunc func(leikari.ActorContext, string) (interface{}, error)
type UpdateFunc func(leikari.ActorContext, string, interface{}) (interface{}, error)
type DeleteFunc func(leikari.ActorContext, string) (interface{}, error)

type CrudHandler struct {
	OnCreate CreateFunc
	OnRead ReadFunc
	OnUpdate UpdateFunc
	OnDelete DeleteFunc
	OnList ListFunc
	OnStart func(leikari.ActorContext) error
	OnStop func(leikari.ActorContext) error
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
		entity, err := a.OnUpdate(ctx, cmd.Id, cmd.Entity)
		if err != nil {
			return nil, err
		}
		return &UpdatedEvent{
			Id: cmd.Id,
			Entity: entity,
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

func (a *CrudHandler) List(ctx leikari.ActorContext, qry Query) (*QueryResult, error) {
	if a.OnList != nil {
		start := time.Now()
		result, cnt, err := a.OnList(ctx, qry)
		if err != nil {
			return nil, err
		}
		return &QueryResult{
			From: qry.From,
			Size: len(result),
			Count: cnt,
			Result: result,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
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

func (a *CrudHandler) AsyncActor() bool { return !a.Sync }