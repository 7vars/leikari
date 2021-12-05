package repository

import (
	"reflect"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
)

var (
	queryType = reflect.TypeOf((*query.Query)(nil)).Elem()
	queryResultType = reflect.TypeOf((*query.QueryResult)(nil))
)

func insertFunc(v interface{}) InsertFunc {
	if meth, ok := leikari.MethodByName(v, "Insert"); ok {
		mt := meth.Type()
		if mt.NumIn() == 2 && leikari.CompareType(mt.In(0), leikari.ActorContextType) && mt.NumOut() == 2 && leikari.CompareType(mt.Out(1), leikari.ErrorType) {
			entitytype := mt.In(1)
			return func(ctx leikari.ActorContext, i interface{}) (interface{}, error) {
				entityvalue := reflect.ValueOf(i)
				if leikari.CompareType(entitytype, entityvalue.Type()) {
					result := meth.Call([]reflect.Value{ reflect.ValueOf(ctx), entityvalue })
					if !result[1].IsNil() {
						return nil, result[1].Interface().(error)
					}
					return result[0].Interface(), nil
				}
				return nil, ErrNotFound
			}
		}
	}
	return func(leikari.ActorContext, interface{}) (interface{}, error) { return nil, ErrNotFound }
}

func selectFunc(v interface{}) SelectFunc {
	if meth, ok := leikari.MethodByName(v, "Select"); ok {
		mt := meth.Type()
		if mt.NumIn() == 2 && leikari.CompareType(mt.In(0), leikari.ActorContextType) && mt.NumOut() == 2 && leikari.CompareType(mt.Out(1), leikari.ErrorType) {
			idtype := mt.In(1)
			return func(ctx leikari.ActorContext, i interface{}) (interface{}, error) {
				ivalue := reflect.ValueOf(i)
				if leikari.CompareType(idtype, ivalue.Type()) {
					result := meth.Call([]reflect.Value{reflect.ValueOf(ctx), ivalue})
					if !result[1].IsNil() {
						return nil, result[1].Interface().(error)
					}
					return result[0].Interface(), nil
				}
				return nil, ErrNotFound
			}
		}
	}
	return func(leikari.ActorContext, interface{}) (interface{}, error) { return nil, ErrNotFound }
}

func updateFunc(v interface{}) UpdateFunc {
	if meth, ok := leikari.MethodByName(v, "Update"); ok {
		mt := meth.Type()
		if mt.NumIn() == 3 && leikari.CompareType(mt.In(0), leikari.ActorContextType) && mt.NumOut() == 1 && leikari.CompareType(mt.Out(0), leikari.ErrorType) {
			idtype := mt.In(1)
			entitytype := mt.In(2)
			return func(ctx leikari.ActorContext, id, entity interface{}) error {
				idvalue := reflect.ValueOf(id)
				entityvalue := reflect.ValueOf(entity)
				if leikari.CompareType(idtype, idvalue.Type()) && leikari.CompareType(entitytype, entityvalue.Type()) {
					result := meth.Call([]reflect.Value{ reflect.ValueOf(ctx), idvalue, entityvalue })
					if !result[0].IsNil() {
						return result[0].Interface().(error)
					}
					return nil
				}
				return ErrNotFound
			}
		}
	}
	return func(ac leikari.ActorContext, i1, i2 interface{}) error { return ErrNotFound }
}

func deleteFunc(v interface{}) DeleteFunc {
	if meth, ok := leikari.MethodByName(v, "Delete"); ok {
		mt := meth.Type()
		if mt.NumIn() == 2 && leikari.CompareType(mt.In(0), leikari.ActorContextType) && mt.NumOut() == 2 && leikari.CompareType(mt.Out(1), leikari.ErrorType) {
			idtype := mt.In(1)
			return func(ctx leikari.ActorContext, id interface{}) (interface{}, error) {
				idvalue := reflect.ValueOf(id)
				if leikari.CompareType(idtype, idvalue.Type()) {
					result := meth.Call([]reflect.Value{ reflect.ValueOf(ctx), idvalue })
					if !result[1].IsNil() {
						return nil, result[1].Interface().(error)
					}
					return result[0].Interface(), nil
				}
				return nil, ErrNotFound
			}
		}
	}
	return func(ac leikari.ActorContext, i interface{}) (interface{}, error) { return nil, ErrNotFound }
}

func queryFunc(v interface{}) QueryFunc {
	if meth, ok := leikari.MethodByName(v, "Query"); ok {
		mt := meth.Type()
		if leikari.CheckIn(mt, leikari.ActorContextType, queryType) && leikari.CheckOut(mt, queryResultType, leikari.ErrorType) {
			return func(ctx leikari.ActorContext, qry query.Query) (*query.QueryResult, error) {
				result := meth.Call([]reflect.Value{ reflect.ValueOf(ctx), reflect.ValueOf(qry) })
				if !result[1].IsNil() {
					return nil, result[1].Interface().(error)
				}
				return result[0].Interface().(*query.QueryResult), nil
			}
		}
	}
	return func(leikari.ActorContext, query.Query) (*query.QueryResult, error) { return nil, ErrNotFound }
}

func wrap(v interface{}) *RepositoryHandler {
	repo := &RepositoryHandler{
		OnInsert: insertFunc(v),
		OnSelect: selectFunc(v),
		OnUpdate: updateFunc(v),
		OnDelete: deleteFunc(v),
		OnQuery: queryFunc(v),
		OnStart: leikari.PreStartFunc(v),
		OnStop: leikari.PostStopFunc(v),
	}

	return repo
}