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
	Highlight map[string][]string
	Type      string  `json:"type"` // object, host, process
	Score     float64 `json:"score"`
	UrlSuffix string  `json:"url_suffix"`
}

type Query struct {
	QueryString string `json:"query_string"`
}

const (
	CMDBINDEX = "cmdb"
)

var (
	testData = []SearchResult{
		{
			Type:      "cc_ObjectBase",
			Score:     3.566052,
			UrlSuffix: "/#/general-model/test_search",
			Source: map[string]interface{}{
				"jw_test_4":           1,
				"bk_inst_id":          5,
				"bk_supplier_account": "0",
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": "2",
					},
				},
				"bk_obj_id":    "test",
				"bk_inst_name": "1",
				"jw_test_1":    "1",
				"jw_test_2":    12,
				"jw_test_3":    "2019-03-06",
			},
		},
		{
			Type:      "cc_HostBase",
			Score:     2.2986379,
			UrlSuffix: "/#/resource?business=1&ip=10.0.0.6&outer=false&inner=true&exact=1&assigned=true",
			Source: map[string]interface{}{
				"bk_bak_operator":     nil,
				"bk_supplier_account": "0",
				"bk_disk":             50,
				"bk_host_innerip":     "10.0.0.6",
				"bk_os_name":          "windows",
				"import_from":         "1",
				"bk_state_name":       nil,
				"bk_cloud_id":         0,
				"bk_cpu_mhz":          2,
				"bk_mac":              "aa:aa:aa:aa:aa:aa",
				"bk_asset_id":         "",
				"bk_comment":          "this is test host",
				"bk_host_name":        "",
				"bk_host_outerip":     "175.0.0.6",
				"bk_outer_mac":        "aa:aa:aa:aa:aa:aa",
				"operator":            nil,
				"bk_isp_name":         nil,
				"bk_os_version":       "",
				"bk_service_term":     nil,
				"bk_sla":              nil,
				"bk_os_type":          nil,
				"bk_cpu_module":       "",
				"bk_mem":              nil,
				"bk_os_bit":           "32",
				"bk_sn":               "",
				"bk_province_name":    nil,
				"bk_cpu":              nil,
				"create_time":         nil,
				"bk_host_id":          2,
			},
		},
	}
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
