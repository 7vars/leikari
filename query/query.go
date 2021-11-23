package query

import "time"

type Query struct {
	From int `json:"from"`
	Size int `json:"size,omitempty"`
	Query string `json:"query"`
}

type QueryResult struct {
	From int `json:"from"`
	Size int `json:"size,omitempty"`
	Count int `json:"count,omitempty"`
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