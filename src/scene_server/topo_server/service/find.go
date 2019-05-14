package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
	"encoding/json"
	"github.com/olivere/elastic"
)

type SearchResult struct {
	Source    map[string]interface{} `json:"source"` // data from mongo, key/value
	Highlight map[string][]string    `json:"highlight"`
	Type      string                 `json:"type"` // object, host, process, model
	Score     float64                `json:"score"`
}

type Aggregation struct {
	Key   interface{} `json:"key"`
	Count int64       `json:"count"`
}

type SearchResults struct {
	Total        int64          `json:"total"`
	Aggregations []Aggregation  `json:"aggregations"`
	Hits         []SearchResult `json:"hits,omitempty"`
}

type Page struct {
	Start int `json:"start"`
	Limit int `json:"limit"`
}

type Query struct {
	Paging      Page   `json:"page"`
	QueryString string `json:"query_string"`
}

func (s *Service) FullTextFind(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	if data.Exists("query_string") {
		query := new(Query)
		// set paging default
		query.Paging.Start = -1
		query.Paging.Limit = -1
		if err := data.MarshalJSONInto(query); err != nil {
			blog.Errorf("full_text_find failed, import query_string, but got invalid query_string:[%v], err: %+v", query, err)
			return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
		}

		result, err := s.Es.CmdbSearch(query.QueryString, query.Paging.Start, query.Paging.Limit)
		if err != nil {
			blog.Errorf("full_text_find failed, find failed, err: %+v", err)
			return nil, params.Err.Error(common.CCErrorTopoFullTextFindErr)
		}

		// result is hits and aggregations
		searchResults := new(SearchResults)

		searchResults.Total = result.Hits.TotalHits
		// set hits
		for _, hit := range result.Hits.Hits {
			// ignore not correct cmdb table data
			if hit.Index == common.CMDBINDEX && hit.Id != common.INDICES {
				sr := SearchResult{}
				sr.setHit(hit)
				searchResults.Hits = append(searchResults.Hits, sr)
			}
		}

		// set aggregations
		bkObjIdAggs, found := result.Aggregations.Terms(common.BkObjIdAggName)
		if found == true && bkObjIdAggs != nil {
			for _, bucket := range bkObjIdAggs.Buckets {
				agg := Aggregation{}
				agg.setAgg(bucket)
				searchResults.Aggregations = append(searchResults.Aggregations, agg)
			}
		}

		typeAggs, found := result.Aggregations.Terms(common.TypeAggName)
		if found == true && typeAggs != nil {
			for _, bucket := range typeAggs.Buckets {
				agg := Aggregation{}
				agg.setAgg(bucket)
				searchResults.Aggregations = append(searchResults.Aggregations, agg)
			}
		}
		// test data
		//searchResults := testData
		return searchResults, nil
	}

	return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
}

func (agg *Aggregation) setAgg(bucket *elastic.AggregationBucketKeyItem) {
	if bucket.Key == common.BKTableNameBaseHost {
		agg.Key = "host"
	} else {
		agg.Key = bucket.Key
	}

	agg.Count = bucket.DocCount
}

func (sr *SearchResult) setHit(searchHit *elastic.SearchHit) {
	sr.Score = *searchHit.Score
	switch searchHit.Type {
	case common.BKTableNameBaseInst:
		sr.Type = "object"
	case common.BKTableNameBaseHost:
		sr.Type = "host"
	case common.BKTableNameBaseProcess:
		sr.Type = "process"
	case common.BKTableNameBaseApp:
		sr.Type = "application"
	case common.BKTableNameObjDes:
		sr.Type = "model"
	}

	sr.Highlight = searchHit.Highlight
	err := json.Unmarshal(*searchHit.Source, &(sr.Source))
	if err != nil {
		blog.Warnf("full_text_find unmarshal search result source err: %+v", err)
		sr.Source = nil
	}
}
