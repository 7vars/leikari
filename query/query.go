package query

import "time"

type Query struct {
	From int `json:"from"`
	Size int `json:"size,omitempty"`
	Query string `json:"query"`
}

func (qry Query) Parse() (Node, error) {
	return Parse(qry.Query)
}

type QueryResult struct {
	From int `json:"from"`
	Size int `json:"size"`
	Count int `json:"count"`
	Result []interface{} `json:"result,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Took int64 `json:"millis,omitempty"`
}

func NewQueryResult(query Query, result []interface{}, count int) *QueryResult {
	return &QueryResult{
		From: query.From,
		Size: len(result),
		Count: count,
		Result: result,
		Timestamp: time.Now(),
	}
}