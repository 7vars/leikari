package repository

import (
	"context"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
)

type RepositoryRef interface {
	leikari.Ref
	
	Insert(interface{}, interface{}) (*InsertedEvent, error)
	InsertContext(context.Context, interface{}, interface{}) (*InsertedEvent, error)

	Select(interface{}) (*SelectedEvent, error) 
	SelectContext(context.Context, interface{}) (*SelectedEvent, error)

	Update(interface{}, interface{}) (*UpdatedEvent, error)
	UpdateContext(context.Context, interface{}, interface{}) (*UpdatedEvent, error)

	Delete(interface{}) (*DeletedEvent, error)
	DeleteContext(context.Context, interface{}) (*DeletedEvent, error)

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

func (r *repoRef) Insert(id interface{}, entity interface{}) (*InsertedEvent, error) {
	return r.InsertContext(context.Background(), id, entity)
}

func (r *repoRef) InsertContext(ctx context.Context, id interface{}, entity interface{}) (*InsertedEvent, error) {
	res, err := r.RequestContext(ctx, InsertCommand{id, entity})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*InsertedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (r *repoRef) Select(id interface{}) (*SelectedEvent, error)  {
	return r.SelectContext(context.Background(), id)
}

func (r *repoRef) SelectContext(ctx context.Context, id interface{}) (*SelectedEvent, error) {
	res, err := r.RequestContext(ctx, SelectCommand{id})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*SelectedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (r *repoRef) Update(id interface{}, entity interface{}) (*UpdatedEvent, error) {
	return r.UpdateContext(context.Background(), id, entity)
}

func (r *repoRef) UpdateContext(ctx context.Context, id interface{}, entity interface{}) (*UpdatedEvent, error) {
	res, err := r.RequestContext(ctx, UpdateCommand{id, entity})
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*UpdatedEvent); ok {
		return result, nil
	}
	return nil, ErrUnknownCommand
}

func (r *repoRef) Delete(id interface{}) (*DeletedEvent, error) {
	return r.DeleteContext(context.Background(), id)
}

func (r *repoRef) DeleteContext(ctx context.Context, id interface{}) (*DeletedEvent, error) {
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
