package repository

import (
	"sync"
	"time"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/mapper"
	"github.com/7vars/leikari/query"
)

type mapRepo struct {
	sync.RWMutex
	data map[interface{}]interface{}
	keyField string
}

func MapRepo(keyField string) Repository {
	return &mapRepo{
		data: make(map[interface{}]interface{}),
		keyField: keyField,
	}
}

func (mr *mapRepo) Insert(ctx leikari.ActorContext, cmd InsertCommand) (*InsertedEvent, error) {
	start := time.Now()
	if id, ok := mapper.Value(mr.keyField, cmd.Entity); ok {
		mr.Lock()
		defer mr.Unlock()

		if _, exists := mr.data[id]; exists {
			return nil, ErrEntityExists
		}

		mr.data[id] = cmd.Entity
		return &InsertedEvent{
			Id: id,
			Entity: cmd.Entity,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrIdNotPresent
}

func (mr *mapRepo) Select(ctx leikari.ActorContext, cmd SelectCommand) (*SelectedEvent, error) {
	start := time.Now()
	mr.RLock()
	defer mr.RUnlock()

	if val, exists := mr.data[cmd.Id]; exists {
		return &SelectedEvent{
			Id: cmd.Id,
			Entity: val,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (mr *mapRepo) Update(ctx leikari.ActorContext, cmd UpdateCommand) (*UpdatedEvent, error) {
	start := time.Now()
	mr.Lock()
	defer mr.Unlock()

	if _, exists := mr.data[cmd.Id]; exists {
		mr.data[cmd.Id] = cmd.Entity
		return &UpdatedEvent{
			Id: cmd.Id,
			Entity: cmd.Entity,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (mr *mapRepo) Delete(ctx leikari.ActorContext, cmd DeleteCommand) (*DeletedEvent, error) {
	start := time.Now()
	mr.Lock()
	defer mr.Unlock()

	if val, exists := mr.data[cmd.Id]; exists {
		delete(mr.data, cmd.Id)
		return &DeletedEvent{
			Id: cmd.Id,
			Entity: val,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli() - start.UnixMilli(),
		}, nil
	}
	return nil, ErrNotFound
}

func (mr *mapRepo) Query(ctx leikari.ActorContext, qry query.Query) (*query.QueryResult, error) {
	start := time.Now()

	node, err := qry.Parse()
	if err != nil {
		return nil, err
	}

	mr.RLock()
	result := make([]interface{}, 0)
	for _, val := range mr.data {
		if mapper.ApplyFilter(node, val) {
			result = append(result, val)
		}
	}
	mr.Unlock()

	cnt := len(result)
	if qry.From > cnt {
		result = make([]interface{}, 0)
		goto return_result
	}
	result = result[qry.From:]

	if qry.Size < len(result) {
		result = result[:qry.Size]
	}

	return_result:
	return &query.QueryResult{
		From: qry.From,
		Size: len(result),
		Count: cnt,
		Result: result,
		Timestamp: time.Now(),
		Took: time.Now().UnixMilli() - start.UnixMilli(),
	}, nil
}