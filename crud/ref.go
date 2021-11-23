package crud

import (
	"context"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
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

	List(query.Query) (*query.QueryResult, error)
	ListContext(context.Context, query.Query) (*query.QueryResult, error)
}

type crudRef struct {
	leikari.Ref
}

func newCrudRef(ref leikari.Ref) CrudRef {
	return &crudRef{
		ref,
	}
}

func (ref *crudRef) Create(v interface{}) (*CreatedEvent, error) {
	return ref.CreateContext(context.Background(), v)
}

func (ref *crudRef) CreateContext(ctx context.Context, v interface{}) (*CreatedEvent, error) {
	res, err := ref.RequestContext(ctx, CreateCommand{v})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*CreatedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (ref *crudRef) Read(id string) (*ReadEvent, error) {
	return ref.ReadContext(context.Background(), id)
}

func (ref *crudRef) ReadContext(ctx context.Context, id string) (*ReadEvent, error) {
	res, err := ref.RequestContext(ctx, ReadCommand{id})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*ReadEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (ref *crudRef) Update(id string, v interface{}) (*UpdatedEvent, error) {
	return ref.UpdateContext(context.Background(), id, v)
}

func (ref *crudRef) UpdateContext(ctx context.Context, id string, v interface{}) (*UpdatedEvent, error) {
	res, err := ref.RequestContext(ctx, UpdateCommand{id, v})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*UpdatedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (ref *crudRef) Delete(id string) (*DeletedEvent, error) {
	return ref.DeleteContext(context.Background(), id)
}

func (ref *crudRef) DeleteContext(ctx context.Context, id string) (*DeletedEvent, error) {
	res, err := ref.RequestContext(ctx, DeleteCommand{id})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*DeletedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (ref *crudRef) List(qry query.Query) (*query.QueryResult, error) {
	return ref.ListContext(context.Background(), qry)
}

func (ref *crudRef) ListContext(ctx context.Context, qry query.Query) (*query.QueryResult, error) {
	res, err := ref.RequestContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*query.QueryResult); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}