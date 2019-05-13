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

type Page struct {
	Start int `json:"start"`
	Limit int `json:"limit"`
}

type Query struct {
	Paging      Page   `json:"page"`
	QueryString string `json:"query_string"`
}

const (
	CMDBINDEX = "cmdb"
	INDICES   = "indices"
)

var (
	testData = []SearchResult{
		{
			Source: map[string]interface{}{
				"bk_classification_id": "test",
				"bk_ispaused":          false,
				"bk_obj_icon":          "icon-cc-default",
				"bk_obj_id":            "test",
				"bk_obj_name":          "test",
				"bk_supplier_account":  "0",
				"create_time":          "2019-03-06T09:42:55.540000",
				"creator":              "",
				"description":          "",
				"id":                   16,
				"ispre":                false,
				"last_time":            "2019-03-06T09:42:55.540000",
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": "2",
					},
				},
				"modifier": "",
				"position": "",
			},
			Highlight: map[string][]string{
				"bk_classification_id": {
					"<em>test</em>",
				},
				"bk_classification_id.keyword": {
					"<em>test</em>",
				},
				"bk_obj_id": {
					"<em>test</em>",
				},
				"bk_obj_id.keyword": {
					"<em>test</em>",
				},
				"bk_obj_name": {
					"<em>test</em>",
				},
				"bk_obj_name.keyword": {
					"<em>test</em>",
				},
			},
			Type:  "model",
			Score: 1.3420606,
		},
		{
			Source: map[string]interface{}{
				"bk_classification_id": "test",
				"bk_ispaused":          false,
				"bk_obj_icon":          "icon-cc-default",
				"bk_obj_id":            "test2",
				"bk_obj_name":          "test2",
				"bk_supplier_account":  "0",
				"create_time":          "2019-03-13T12:23:06.539000",
				"creator":              "",
				"description":          "",
				"id":                   17,
				"ispre":                false,
				"last_time":            "2019-03-13T12:23:06.539000",
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": "2",
					},
				},
				"modifier": "",
				"position": "",
			},
			Highlight: map[string][]string{
				"bk_classification_id": {
					"<em>test</em>",
				},
				"bk_classification_id.keyword": {
					"<em>test</em>",
				},
			},
			Type:  "model",
			Score: 1.3378919,
		},
		{
			Source: map[string]interface{}{
				"bk_inst_id":          5,
				"bk_inst_name":        "1",
				"bk_obj_id":           "test",
				"bk_supplier_account": "0",
				"jw_test_1":           "1",
				"jw_test_2":           12,
				"jw_test_3":           "2019-03-06",
				"jw_test_4":           1,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": "2",
					},
				},
			},
			Highlight: map[string][]string{
				"bk_obj_id": {
					"<em>test</em>",
				},
				"bk_obj_id.keyword": {
					"<em>test</em>",
				},
			},
			Type:  "object",
			Score: 0.9642885,
		},
		{
			Source: map[string]interface{}{
				"bk_classification_id": "test",
				"bk_ispaused":          false,
				"bk_obj_icon":          "icon-cc-default",
				"bk_obj_id":            "ljp_test",
				"bk_obj_name":          "ljp测试",
				"bk_supplier_account":  "0",
				"create_time":          "2019-03-15T03:18:36.141000",
				"creator":              "cc_system",
				"description":          "",
				"id":                   18,
				"ispre":                false,
				"last_time":            "2019-03-15T03:18:36.141000",
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": "2",
					},
				},
				"modifier": "",
				"position": "",
			},
			Highlight: map[string][]string{
				"bk_classification_id": {
					"<em>test</em>",
				},
				"bk_classification_id.keyword": {
					"<em>test</em>",
				},
			},
			Type:  "model",
			Score: 0.8293917,
		},
		{
			Source: map[string]interface{}{
				"bk_classification_id": "test",
				"bk_ispaused":          false,
				"bk_obj_icon":          "icon-cc-mongodb",
				"bk_obj_id":            "test_search",
				"bk_obj_name":          "测试搜索a",
				"bk_supplier_account":  "0",
				"create_time":          "2019-03-05T03:32:43.407000",
				"creator":              "",
				"description":          "",
				"id":                   14,
				"ispre":                false,
				"last_time":            "2019-04-25T08:41:10.645000",
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": "2",
					},
				},
				"modifier": "admin",
				"position": "",
			},
			Highlight: map[string][]string{
				"bk_classification_id": {
					"<em>test</em>",
				},
				"bk_classification_id.keyword": {
					"<em>test</em>",
				},
			},
			Type:  "model",
			Score: 0.77660173,
		},
	}
)

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

		result, err := s.Es.Search(query.QueryString, CMDBINDEX, query.Paging.Start, query.Paging.Limit)
		if err != nil {
			blog.Errorf("full_text_find failed, find failed, err: %+v", err)
			return nil, params.Err.Error(common.CCErrorTopoFullTextFindErr)
		}

		// result is list
		searchResults := make([]SearchResult, 0)
		for _, hit := range result {
			// ignore not correct cmdb table data
			if hit.Index == CMDBINDEX && hit.Id != INDICES {
				sr := SearchResult{}
				sr.setHit(hit)
				searchResults = append(searchResults, sr)
			}
		}

		// test data
		//searchResults := testData
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
