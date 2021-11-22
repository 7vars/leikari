package crud

import (
	"reflect"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/route"
)

func CRUD(system leikari.System, handler interface{}, name string, opts ...leikari.Option) (CrudRef, route.Route, error) {
	if name == "" {
		return nil, route.Route{}, leikari.Errorln("", "name is not defined")
	}
	if handler == nil {
		return nil, route.Route{}, leikari.Errorln("", "handler is nil")
	}
	if reflect.TypeOf(handler).Kind() != reflect.Ptr {
		return nil, route.Route{}, leikari.Errorln("", "handler must be a pointer")
	}
	
	ref, err := system.Execute(newCrudActor(name, handler), opts...)
	if err != nil {
		return nil, route.Route{}, err
	}
	
	crudRef := newCrudRef(ref)
	return crudRef, newCrudRoute(name, crudRef, handler), nil
}