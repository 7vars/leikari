package crud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/query"
)

type Client interface {
	Create(interface{}) (interface{}, error)
	CreateContext(context.Context, interface{}) (interface{}, error)

	Read(string) (interface{}, error)
	ReadContext(context.Context, string) (interface{}, error)

	Update(string, interface{}) (interface{}, error)
	UpdateContext(context.Context, string, interface{}) (interface{}, error)

	Delete(string) (interface{}, error)
	DeleteContext(context.Context, string) (interface{}, error)

	List(query.Query) (*query.QueryResult, error)
	ListContext(context.Context, query.Query) (*query.QueryResult, error)
}

func CrudClient(baseurl string) (Client, error) {
	return NewCrudClient(http.DefaultClient, baseurl)
}

func NewCrudClient(client *http.Client, baseurl string) (Client, error) {
	u := strings.TrimRight(baseurl, "/")
	if _, err := url.Parse(u); err != nil {
		return nil, err
	}

	return &crudClient{
		client: client,
		baseurl: u,
	}, nil
}

type crudClient struct {
	client *http.Client
	baseurl string
	marshal func(interface{}) ([]byte, error)
	unmarshal func([]byte) (interface{}, error)
}

func (cli *crudClient) request(ctx context.Context, method string, id string, entity interface{}) (interface{}, error) {
	var reader io.Reader
	if entity != nil {
		body, err := cli.marshal(entity)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(body)
	}

	u := cli.baseurl
	if id != "" {
		u = fmt.Sprintf("%s/%s", cli.baseurl, id)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reader)
	if err != nil {
		return nil, err
	}

	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return cli.unmarshal(body)
	}

	var e *leikari.Error
	if err := json.Unmarshal(body, e); err != nil {
		return nil, err
	}
	return nil, e.WithStatusCode(resp.StatusCode)
}

func (cli *crudClient) Create(entity interface{}) (interface{}, error) {
	return cli.CreateContext(context.Background(), entity)
}

func (cli *crudClient) CreateContext(ctx context.Context, entity interface{}) (interface{}, error) {
	return cli.request(ctx, "POST", "", entity)
}

func (cli *crudClient) Read(id string) (interface{}, error) {
	return cli.ReadContext(context.Background(), id)
}

func (cli *crudClient) ReadContext(ctx context.Context, id string) (interface{}, error) {
	return cli.request(ctx, "GET", id, nil)
}

func (cli *crudClient) Update(id string, entity interface{}) (interface{}, error) {
	return cli.UpdateContext(context.Background(), id, entity)
}

func (cli *crudClient) UpdateContext(ctx context.Context, id string, entity interface{}) (interface{}, error) {
	return cli.request(ctx, "PUT", id, entity)
}

func (cli *crudClient) Delete(id string) (interface{}, error) {
	return cli.DeleteContext(context.Background(), id)
}

func (cli *crudClient) DeleteContext(ctx context.Context, id string) (interface{}, error) {
	return cli.request(ctx, "POST", id, nil)
}

func (cli *crudClient) List(qry query.Query) (*query.QueryResult, error) {
	return cli.ListContext(context.Background(), qry)
}

type customResult struct {
	unmarshal func([]byte) (interface{}, error)
	From int `json:"from"`
	Size int `json:"size"`
	Count int `json:"count"`
	Result []interface{} `json:"result,omitempty"`
}

func newCustomResult(unmarshal func([]byte) (interface{}, error)) *customResult {
	return &customResult{
		unmarshal: unmarshal,
	}
}

func (cr *customResult) UnmarshalJSON(data []byte) error {
	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(data, &objmap); err != nil {
		return err
	}

	if item, ok := objmap["from"]; ok {
		if err := json.Unmarshal(*item, &cr.From); err != nil {
			return err
		}
	}

	if item, ok := objmap["size"]; ok {
		if err := json.Unmarshal(*item, &cr.Size); err != nil {
			return err
		}
	}

	if item, ok := objmap["count"]; ok {
		if err := json.Unmarshal(*item, &cr.Count); err != nil {
			return err
		}
	}

	if item, ok := objmap["result"]; ok {
		var raws []*json.RawMessage
		if err := json.Unmarshal(*item, &raws); err != nil {
			return err
		}
		for _, itm := range raws {
			item, err := cr.unmarshal(*itm)
			if err != nil {
				return err
			}
			cr.Result = append(cr.Result, item)
		}
	}

	return nil
}

func (cli *crudClient) ListContext(ctx context.Context, qry query.Query) (*query.QueryResult, error) {
	start := time.Now()

	body, err := json.Marshal(&qry)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/_query", cli.baseurl), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		res := newCustomResult(cli.unmarshal)
		if err := json.Unmarshal(respbody, res); err != nil {
			return nil, err
		}
		return &query.QueryResult{
			From: res.From,
			Size: len(res.Result),
			Count: res.Count,
			Result: res.Result,
			Timestamp: time.Now(),
			Took: time.Now().UnixMilli()-start.UnixMilli(),
		}, nil
	}

	var e *leikari.Error
	if err := json.Unmarshal(body, e); err != nil {
		return nil, err
	}
	return nil, e.WithStatusCode(resp.StatusCode)
}