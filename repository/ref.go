package repository

import (
	"context"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
)

type RepositoryRef interface {
	leikari.Ref
	
	Insert(string, interface{}) (*InsertedEvent, error)
	InsertContext(context.Context, string, interface{}) (*InsertedEvent, error)

	Select(string) (*SelectedEvent, error) 
	SelectContext(context.Context, string) (*SelectedEvent, error)

	Update(string, interface{}) (*UpdatedEvent, error)
	UpdateContext(context.Context, string, interface{}) (*UpdatedEvent, error)

	Delete(string) (*DeletedEvent, error)
	DeleteContext(context.Context, string) (*DeletedEvent, error)

	Query(query.Query) (*query.QueryResult, error)
	QueryContext(context.Context, query.Query) (*query.QueryResult, error)
}

type repoRef struct {
	leikari.Ref
}

func newRepositoryRef(ref leikari.Ref) RepositoryRef {
	return &repoRef{
		ref,
	}
}

func (r *repoRef) Insert(id string, entity interface{}) (*InsertedEvent, error) {
	return r.InsertContext(context.Background(), id, entity)
}

func (r *repoRef) InsertContext(ctx context.Context, id string, entity interface{}) (*InsertedEvent, error) {
	res, err := r.RequestContext(ctx, InsertCommand{id, entity})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*InsertedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (r *repoRef) Select(id string) (*SelectedEvent, error)  {
	return r.SelectContext(context.Background(), id)
}

func (r *repoRef) SelectContext(ctx context.Context, id string) (*SelectedEvent, error) {
	res, err := r.RequestContext(ctx, SelectCommand{id})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*SelectedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (r *repoRef) Update(id string, entity interface{}) (*UpdatedEvent, error) {
	return r.UpdateContext(context.Background(), id, entity)
}

func (r *repoRef) UpdateContext(ctx context.Context, id string, entity interface{}) (*UpdatedEvent, error) {
	res, err := r.RequestContext(ctx, UpdateCommand{id, entity})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*UpdatedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (r *repoRef) Delete(id string) (*DeletedEvent, error) {
	return r.DeleteContext(context.Background(), id)
}

func (r *repoRef) DeleteContext(ctx context.Context, id string) (*DeletedEvent, error) {
	res, err := r.RequestContext(ctx, DeleteCommand{id})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*DeletedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (r *repoRef) Query(qry query.Query) (*query.QueryResult, error) {
	return r.QueryContext(context.Background(), qry)
}

func (r *repoRef) QueryContext(ctx context.Context, qry query.Query) (*query.QueryResult, error) {
	res, err := r.RequestContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*query.QueryResult); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}
