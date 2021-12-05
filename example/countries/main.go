package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/7vars/leikari"
	"github.com/7vars/leikari/crud"
	lhttp "github.com/7vars/leikari/http"
	"github.com/7vars/leikari/mapper"
	"github.com/7vars/leikari/query"
	"github.com/7vars/leikari/repository"
	"github.com/7vars/leikari/route"
)

type dhCountry struct {
	ISO string `json:"ISO3166-1-Alpha-2"`
	ISO3 string `json:"ISO3166-1-Alpha-3"`
	Name string `json:"official_name_en"`
	Capital string `json:"Capital"`
	DS string `json:"DS"`
	TLD string `json:"TLD"`
	Continent string `json:"Continent"`
	Currency string `json:"ISO4217-currency_alphabetic_code"`
	Languages string `json:"Languages"`
}

func (dh *dhCountry) ToCountry() *Country {
	return &Country{
		ISO: dh.ISO,
		ISO3: dh.ISO3,
		Name: dh.Name,
		Capital: dh.Capital,
		DS: dh.DS,
		TLD: dh.TLD,
	}
}

type Country struct {
	ISO string `json:"iso"`
	ISO3 string `json:"iso3"`
	Name string `json:"name"`
	Capital string `json:"capital"`
	DS string `json:"ds"`
	TLD string `json:"tld"`
}

func (c *Country) String() string {
	return fmt.Sprintf("%s - %s", c.ISO, c.Name)
}

type CountryRepo struct {
	sync.RWMutex
	url string
	countries []*Country
}

func newCountryRepo() *CountryRepo {
	return &CountryRepo{
		url: "https://datahub.io/core/country-codes/r/0.json",
		countries: make([]*Country, 0),
	}
}

func (repo *CountryRepo) PreStart(ctx leikari.ActorContext) error {
	go func() {
		start := time.Now()
		defer func() {
			ctx.Log().Infof("countries loaded in %d ms", (time.Now().UnixMilli()-start.UnixMilli()))
		}()
		ctx.Log().Infof("load countries from %s", repo.url)
		resp, err := http.Get(repo.url)
		if err != nil {
			ctx.Log().Errorf("could not load countries from %v: %v", repo.url, err)
		}
		defer resp.Body.Close()

		var countries []*dhCountry
		if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
			ctx.Log().Errorf("could not unmarshal countries: %v", err)
		}

		repo.Lock()
		defer repo.Unlock()
		
		for _, c := range countries {
			if c.ISO != "" {
				repo.countries = append(repo.countries, c.ToCountry())
			}
		}
		sort.SliceStable(repo.countries, func(i, j int) bool {
			return repo.countries[i].ISO < repo.countries[j].ISO
		})
	}()

	return nil
}

func (repo *CountryRepo) Query(ctx leikari.ActorContext, qry query.Query) (*query.QueryResult, error) {
	node, err := qry.Parse()
	if err != nil {
		return nil, err
	}

	repo.RLock()
	var result []interface{}
	for _, c := range repo.countries {
		if mapper.ApplyFilter(node, c) {
			result = append(result, c)
		}
	}
	repo.RUnlock()

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
	}, nil
}

func (repo *CountryRepo) Insert(ctx leikari.ActorContext, entity *Country) (string, error) {
	id := entity.ISO
	if id == "" {
		return "", repository.ErrIdNotPresent
	}
	if e, _ := repo.Select(ctx, id); e != nil {
		return "", repository.ErrEntityExists
	}
	repo.Lock()
	defer repo.Unlock()
	repo.countries = append(repo.countries, entity)
	return id, nil
}

func (repo *CountryRepo) Select(ctx leikari.ActorContext, id string) (*Country, error) {
	res, err := repo.Query(ctx, query.Query{ From: 0, Size: 1, Query: fmt.Sprintf("iso EQ '%v'", id) })
	if err != nil {
		return nil, err
	}
	if res.Size > 0 {
		return res.Result[0].(*Country), nil
	}
	return nil, repository.ErrNotFound
}

func (repo *CountryRepo) Update(ctx leikari.ActorContext, id string, entity *Country) error {
	repo.Lock()
	defer repo.Unlock()
	for i, c := range repo.countries {
		if c.ISO == id {
			entity.ISO = id
			repo.countries[i] = entity
			return nil
		}
	}
	return repository.ErrNotFound
}

func (repo *CountryRepo) Delete(ctx leikari.ActorContext, id string) (*Country, error) {
	repo.Lock()
	defer repo.Unlock()
	for i, c := range repo.countries {
		if c.ISO == id {
			repo.countries = append(repo.countries[:i], repo.countries[i+1:]...)
			return c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func newCrudHandler(ref repository.RepositoryRef) *crud.CrudHandler {
	return &crud.CrudHandler{
		OnCreate: func(ac leikari.ActorContext, entity interface{}) (string, interface{}, error) {
			evt, err := ref.Insert(entity)
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("%v", evt.Id), evt.Entity, nil
		},
		OnQuery: func(ac leikari.ActorContext, qry query.Query) (*query.QueryResult, error) {
			return ref.Query(qry)
		},
		OnRead: func(ac leikari.ActorContext, id string) (interface{}, error) {
			res, err := ref.Select(strings.ToUpper(id))
			if err != nil {
				return nil, err
			}
			return res.Entity, nil
		},
		OnUpdate: func(ac leikari.ActorContext, id string, entity interface{}) error {
			if _, err := ref.Update(strings.ToUpper(id), entity); err != nil {
				return err
			}
			return nil
		},
		OnDelete: func(ac leikari.ActorContext, id string) (interface{}, error) {
			res, err := ref.Delete(strings.ToUpper(id))
			if err != nil {
				return nil, err
			}
			return res.Entity, nil
		},
		OnUnmarshal: func(b []byte) (interface{}, error) {
			var country Country
			if err := json.Unmarshal(b, &country); err != nil {
				return nil, err
			}
			return &country, nil
		},
	}
}

func main() {
	sys := leikari.NewSystem()

	repoRef, err := repository.RepositoryService(sys, newCountryRepo(), "country-repo")
	if err != nil {
		panic(err)
	}

	_, coutryRoute, err := crud.CrudService(sys, newCrudHandler(repoRef), "country")
	if err != nil {
		panic(err)
	}

	route := route.Route{
		Name: "v1",
		Path: "/api/v1",
		Routes: []route.Route{ coutryRoute },
	}

	_, err = lhttp.HttpServer(sys, route)
	if err != nil {
		panic(err)
	}

	sys.Run()
}