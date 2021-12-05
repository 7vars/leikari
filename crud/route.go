package crud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/7vars/leikari/query"
	"github.com/7vars/leikari/route"
)

func HandleQuery(ref CrudRef) func(r route.Request) route.Response {
	return func(r route.Request) route.Response {
		query := query.Query{
			From: 0,
			Size: 10, // TODO configure
		}
		if url := r.URL(); url != nil {
			var qvals []string
			qry := url.Query()
			for key := range qry {
				switch key {
				case "from":
					if i, err := strconv.Atoi(qry.Get(key)); err == nil && i >= 0 {
						query.From = i
					}
				case "size":
					if i, err := strconv.Atoi(qry.Get(key)); err == nil && i > 0 {
						query.Size = i
					}
				default:
					value := qry.Get(key)
					// TODO check int, float and bool
					qvals = append(qvals, fmt.Sprintf("%s EQ %s", key, value))
				}
			}
			query.Query = strings.Join(qvals, " AND ")
		}

		result, err := ref.ListContext(r.Context(), query)
		if err != nil {
			return route.ErrorResponse(err)
		}
		return route.Response{
			Data: result,
		}
	}
}

func HandlePostQuery(ref CrudRef) func(r route.Request) route.Response {
	return func(r route.Request) route.Response {
		var query query.Query
		if err := r.Encode(&query); err != nil {
			return route.ErrorResponseWithStatus(400, err)
		}
		if query.From < 0 {
			query.From = 0
		}
		if query.Size <= 0 {
			query.Size = 10 // TODO configure
		}
		result, err := ref.ListContext(r.Context(), query)
		if err != nil {
			return route.ErrorResponse(err)
		}
		return route.Response{
			Data: result,
		}
	}
}

func HandleCreate(unmarshal func([]byte) (interface{}, error)) func(CrudRef) func(r route.Request) route.Response {
	return func(ref CrudRef) func(r route.Request) route.Response {
		return func(r route.Request) route.Response {
			buf, err := r.Body()
			if err != nil {
				return route.ErrorResponseWithStatus(400, err)
			}
			entity, err := unmarshal(buf)
			if err != nil {
				return route.ErrorResponseWithStatus(400, err)
			}
			env, err := ref.CreateContext(r.Context(), entity)
			if err != nil {
				return route.ErrorResponse(err)
			}
			return route.Response{
				Status: 201,
				Data: env.Entity,
			}
		}
	}
}

func HandleRead(ref CrudRef) func(r route.Request) route.Response {
	return func(r route.Request) route.Response {
		id := r.GetVar("id")
		evt, err := ref.ReadContext(r.Context(), id)
		if err != nil {
			return route.ErrorResponse(err)
		}
		return route.Response{
			Data: evt.Entity,
		}
	}
}

func HandleUpdate(unmarshal func([]byte) (interface{}, error)) func(CrudRef) func(r route.Request) route.Response {
	return func(ref CrudRef) func(r route.Request) route.Response {
		return func(r route.Request) route.Response {
			id := r.GetVar("id")
			buf, err := r.Body()
			if err != nil {
				return route.ErrorResponseWithStatus(400, err)
			}
			entity, err := unmarshal(buf)
			if err != nil {
				return route.ErrorResponseWithStatus(400, err)
			}
			evt, err := ref.UpdateContext(r.Context(), id, entity)
			if err != nil {
				return route.ErrorResponse(err)
			}
			return route.Response{
				Status: 200,
				Data: evt.Entity,
			}
		}
	}
}

func HandleDelete(ref CrudRef) func(r route.Request) route.Response {
	return func(r route.Request) route.Response {
		id := r.GetVar("id")
		evt, err := ref.DeleteContext(r.Context(), id)
		if err != nil {
			return route.ErrorResponse(err)
		}
		return route.Response{
			Data: evt.Entity,
		}
	}
}

func HandleUnmarshal(b []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}

type QueryRouteHandler interface {
	HandleQuery(CrudRef) func(route.Request) route.Response
}

type PostQueryRouteHandler interface {
	HandlePostQuery(CrudRef) func(route.Request) route.Response
}

type CreateRouteHandler interface {
	HandleCreate(CrudRef) func(route.Request) route.Response
}

type ReadRouteHandler interface {
	HandleRead(CrudRef) func(route.Request) route.Response
}

type UpdateRouteHandler interface {
	HandleUpdate(CrudRef) func(route.Request) route.Response
}

type DeleteRouteHandler interface {
	HandleDelete(CrudRef) func(route.Request) route.Response
}

type UnmarshalHandler interface {
	HandleUnmarshal([]byte) (interface{}, error)
}

func newCrudRoute(name string, ref CrudRef, handler interface{}) route.Route {
	handleUnmarshal := HandleUnmarshal
	if hdl, ok := handler.(UnmarshalHandler); ok {
		handleUnmarshal = hdl.HandleUnmarshal
	}

	handleQuery := HandleQuery
	if hdl, ok := handler.(QueryRouteHandler); ok {
		handleQuery = hdl.HandleQuery
	}

	handlePostQuery := HandlePostQuery
	if hdl, ok := handler.(PostQueryRouteHandler); ok {
		handlePostQuery = hdl.HandlePostQuery
	}

	handleCreate := HandleCreate(handleUnmarshal)
	if hdl, ok := handler.(CreateRouteHandler); ok {
		handleCreate = hdl.HandleCreate
	}

	handleRead := HandleRead
	if hdl, ok := handler.(ReadRouteHandler); ok {
		handleRead = hdl.HandleRead
	}

	handleUpdate := HandleUpdate(handleUnmarshal)
	if hdl, ok := handler.(UpdateRouteHandler); ok {
		handleUpdate = hdl.HandleUpdate
	}

	handleDelete := HandleDelete
	if hdl, ok := handler.(DeleteRouteHandler); ok {
		handleDelete = hdl.HandleDelete
	}

	return route.Route{
		Name: name,
		Path: "/" + strings.ToLower(name),
		Method: "GET",
		Handle: handleQuery(ref),
		Routes: []route.Route{
			{
				Name: name + "_query",
				Path: "/_query",
				Method: "POST",
				Handle: handlePostQuery(ref),
			},
			{
				Name: name + "_create",
				Path: "",
				Method: "POST",
				Handle: handleCreate(ref),
			},
			{
				Name: name + "_read",
				Path: "/{id}",
				Method: "GET",
				Handle: handleRead(ref),
			},
			{
				Name: name + "_update",
				Path: "/{id}",
				Method: "PUT",
				Handle: handleUpdate(ref),
			},
			// TODO for future versions implement PATCH to update specific attributes
			{
				Name: name + "_delete",
				Path: "/{id}",
				Method: "DELETE",
				Handle: handleDelete(ref),
			},
		},
	}
}