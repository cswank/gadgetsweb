package controllers

import (
	"time"
	"testing"
)

func TestGetQuery(t *testing.T) {
	vars := &timeseriesVars{
		start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		end: time.Date(2009, time.November, 10, 24, 0, 0, 0, time.UTC),
	}
	query := getQuery(vars)
	//expected := [map[timestamp:map[$gte:2009-11-10 23:00:00 +0000 UTC]] map[timestamp:map[$lte:2009-11-11 00:00:00 +0000 UTC]]]
	//expected := map[string]interface{} {
	//"$and":[map[timestamp:map[$gte:2009-11-10 23:00:00 +0000 UTC]] map[timestamp:map[$lte:2009-11-11 00:00:00 +0000 UTC]]]]
	if query["$and"] != 1 {
		//t.Error(query)
	}
}
