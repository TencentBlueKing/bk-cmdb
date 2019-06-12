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
	Hits         []SearchResult `json:"hits"`
}

type Page struct {
	Start int `json:"start"`
	Limit int `json:"limit"`
}

type Query struct {
	Paging      Page     `json:"page"`
	QueryString string   `json:"query_string"`
	TypeFilter  []string `json:"filter"`
	BkObjId     string   `json:"bk_obj_id"`
	BkBizId     string   `json:"bk_biz_id"`
}

func (s *Service) FullTextFind(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	if data.Exists("query_string") {
		query := new(Query)
		// set paging default
		query.Paging.Start = -1
		query.Paging.Limit = -1
		query.BkObjId = ""
		query.BkBizId = ""
		if err := data.MarshalJSONInto(query); err != nil {
			blog.Errorf("full_text_find failed, import query params, but got invalid query params:[%v], err: %+v", query, err)
			return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
		}

		// check query string
		checkQueryString(query)
		// get query and search types
		queryEs, types := getEsQueryAndSearchTypes(query)

		result, err := s.Es.CmdbSearch(queryEs, types, query.Paging.Start, query.Paging.Limit)
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
				// only cc_HostBase currently
				if bucket.Key == common.BKTableNameBaseHost {
					agg := Aggregation{}
					agg.setAgg(bucket)
					searchResults.Aggregations = append(searchResults.Aggregations, agg)
				}
			}
		}

		return searchResults, nil
	}

	return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
}

func checkQueryString(query *Query) {
	// if query string is single string in SpecialChar, make it to null string
	for i := range common.SpecialChar {
		if query.QueryString == common.SpecialChar[i] {
			query.QueryString = ""
			break
		}
	}
}

func getEsQueryAndSearchTypes(query *Query) (elastic.Query, []string) {
	qBool := elastic.NewBoolQuery()

	// if set bk_biz_id
	qBool.MinimumNumberShouldMatch(1)
	qBizRegex := elastic.NewRegexpQuery(common.BkBizMetaKey, "[0-9]*")
	qBizBool := elastic.NewBoolQuery()
	qBizBool.MustNot(qBizRegex)
	qBool.Should(qBizBool)
	if query.BkBizId != "" {
		qBizTerm := elastic.NewTermQuery(common.BkBizMetaKey, query.BkBizId)
		qBool.Should(qBizTerm)
	}

	// ignore bk_supplier_account
	qSupplierMatch := elastic.NewMatchQuery(common.BkSupplierAccount, query.QueryString)
	qBool.MustNot(qSupplierMatch)

	// if set bk_obj_id
	qString := elastic.NewQueryStringQuery(query.QueryString)
	if query.BkObjId == "" {
		// get search types from filter
		types := getEsIndexTypes(query.TypeFilter)
		return qBool.Must(qString), types
	} else if query.BkObjId == common.TypeHost {
		// if bk_obj_id is host, we search only from type cc_HostBase
		types := []string{common.BKTableNameBaseHost}
		return qBool.Must(qString), types
	} else {
		// if define bk_obj_id, we use bool query include must(bk_obj_id=xxx) and should(query string)
		qBool.Must(elastic.NewTermQuery("bk_obj_id", query.BkObjId))
		qBool.Must(qString)
		types := getEsIndexTypes(query.TypeFilter)
		return qBool, types
	}
}

func getEsIndexTypes(typesFilter []string) []string {
	typesMap := make([]string, 0)
	for _, filter := range typesFilter {
		switch filter {
		case common.TypeHost:
			typesMap = append(typesMap, common.BKTableNameBaseHost)
		case common.TypeModel:
			typesMap = append(typesMap, common.BKTableNameObjDes)
		case common.TypeObject:
			typesMap = append(typesMap, common.BKTableNameBaseInst)
		case common.TypeApplication:
			typesMap = append(typesMap, common.BKTableNameBaseApp)
		case common.TypeProcess:
			typesMap = append(typesMap, common.BKTableNameBaseProcess)
		}
	}

	types := make([]string, 0)
	for _, value := range common.CmdbFindTypes {
		if !inTypes(value, typesMap) {
			types = append(types, value)
		}
	}

	return types
}

func inTypes(val string, types []string) bool {
	for _, v := range types {
		if v == val {
			return true
		}
	}
	return false
}

func (agg *Aggregation) setAgg(bucket *elastic.AggregationBucketKeyItem) {
	if bucket.Key == common.BKTableNameBaseHost {
		agg.Key = common.TypeHost
	} else {
		agg.Key = bucket.Key
	}

	agg.Count = bucket.DocCount
}

func (sr *SearchResult) setHit(searchHit *elastic.SearchHit) {
	sr.Score = *searchHit.Score
	switch searchHit.Type {
	case common.BKTableNameBaseInst:
		sr.Type = common.TypeObject
	case common.BKTableNameBaseHost:
		sr.Type = common.TypeHost
	case common.BKTableNameBaseProcess:
		sr.Type = common.TypeProcess
	case common.BKTableNameBaseApp:
		sr.Type = common.TypeApplication
	case common.BKTableNameObjDes:
		sr.Type = common.TypeModel
	}

	sr.Highlight = searchHit.Highlight
	err := json.Unmarshal(*searchHit.Source, &(sr.Source))
	if err != nil {
		blog.Warnf("full_text_find unmarshal search result source err: %+v", err)
		sr.Source = nil
	}
}
