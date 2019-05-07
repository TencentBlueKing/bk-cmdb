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
	Highlight map[string][]string	  `json:"highlight"`
	Type      string  `json:"type"` // object, host, process
	Score     float64 `json:"score"`
}

type Query struct {
	QueryString string `json:"query_string"`
}

const (
	CMDBINDEX = "cmdb"
)

func (s *Service) FullTextFind(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	if data.Exists("query_string") {
		query := new(Query)
		if err := data.MarshalJSONInto(query); err != nil {
			blog.Errorf("full_text_find failed, import query_string, but got invalid query_string:[%v], err: %+v", query, err)
			return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
		}

		result, err := s.Es.Search(query.QueryString, CMDBINDEX)
		if err != nil {
			blog.Errorf("full_text_find failed, find failed, err: %+v", err)
			return nil, params.Err.Error(common.CCErrorTopoFullTextFindErr)
		}

		// result is list
		searchResults := make([]SearchResult, 0)
		for _, hit := range result {
			if hit.Index == CMDBINDEX {
				sr := SearchResult{}
				sr.setHit(hit)
				searchResults = append(searchResults, sr)
			}
		}
		return searchResults, nil
	}

	return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
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
