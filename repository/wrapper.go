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
		OnSelect: selectFunc(v),
		OnQuery: queryFunc(v),
		OnStart: leikari.PreStartFunc(v),
		OnStop: leikari.PostStopFunc(v),
	}

	return repo
}