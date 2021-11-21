package crud

import (
	"context"

	"github.com/7vars/leikari"
)

type CrudRef interface {
	leikari.Ref
	Create(interface{}) (*CreatedEvent, error)
	CreateContext(context.Context, interface{}) (*CreatedEvent, error)

	Read(string) (*ReadEvent, error)
	ReadContext(context.Context, string) (*ReadEvent, error)

	Update(string, interface{}) (*UpdatedEvent, error)
	UpdateContext(context.Context, string, interface{}) (*UpdatedEvent, error)

	Delete(string) (*DeletedEvent, error)
	DeleteContext(context.Context, string) (*DeletedEvent, error)

	List(Query) (*QueryResult, error)
	ListContext(context.Context, Query) (*QueryResult, error)
}

type crud struct {
	leikari.Ref
}

func newCrudRef(ref leikari.Ref) CrudRef {
	return &crud{
		ref,
	}
}

func (ref *crud) Create(v interface{}) (*CreatedEvent, error) {
	return ref.CreateContext(context.Background(), v)
}

func (ref *crud) CreateContext(ctx context.Context, v interface{}) (*CreatedEvent, error) {
	res, err := ref.RequestContext(ctx, CreateCommand{v})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*CreatedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (ref *crud) Read(id string) (*ReadEvent, error) {
	return ref.ReadContext(context.Background(), id)
}

func (ref *crud) ReadContext(ctx context.Context, id string) (*ReadEvent, error) {
	res, err := ref.RequestContext(ctx, ReadCommand{id})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*ReadEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (ref *crud) Update(id string, v interface{}) (*UpdatedEvent, error) {
	return ref.UpdateContext(context.Background(), id, v)
}

func (ref *crud) UpdateContext(ctx context.Context, id string, v interface{}) (*UpdatedEvent, error) {
	res, err := ref.RequestContext(ctx, UpdateCommand{id, v})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*UpdatedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (ref *crud) Delete(id string) (*DeletedEvent, error) {
	return ref.DeleteContext(context.Background(), id)
}

func (ref *crud) DeleteContext(ctx context.Context, id string) (*DeletedEvent, error) {
	res, err := ref.RequestContext(ctx, DeleteCommand{id})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*DeletedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (ref *crud) List(qry Query) (*QueryResult, error) {
	return ref.ListContext(context.Background(), qry)
}

func (ref *crud) ListContext(ctx context.Context, qry Query) (*QueryResult, error) {
	res, err := ref.RequestContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*QueryResult); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}