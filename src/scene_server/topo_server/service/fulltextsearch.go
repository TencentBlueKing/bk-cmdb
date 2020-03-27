package service

import (
	"context"
	"encoding/json"
	"strings"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"

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

func NewQuery() *Query {
	query := new(Query)
	// set paging default
	query.Paging.Start = -1
	query.Paging.Limit = -1
	query.BkObjId = ""
	query.BkBizId = ""
	return query
}

func (s *Service) FullTextFind(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	if s.Es.Client == nil {
		blog.Errorf("FullTextFind failed, es client is nil, rid: %s", params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoFullTextClientNotInitialized)
	}

	if _, exist := data["query_string"]; exist == false {
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	query := NewQuery()
	if err := data.MarshalJSONInto(query); err != nil {
		blog.Errorf("full_text_find failed, import query params, but got invalid query params:[%v], err: %+v, rid: %s", query, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommJSONMarshalFailed)
	}

	// check query string
	rawString, ok := query.checkQueryString()
	if !ok {
		blog.Errorf("full_text_find failed, query string [%s] large than 32, rid: %s", rawString, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsIsInvalid, "query_string")
	}
	// get query and search types
	esQuery, searchTypes := query.toEsQueryAndSearchTypes()

	result, err := s.Es.Search(params.Context, esQuery, searchTypes, query.Paging.Start, query.Paging.Limit)
	if err != nil {
		blog.Errorf("full_text_find failed, es search failed, err: %+v, rid: %s", err, params.ReqID)
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
			sr.setHit(params.Context, hit, query.BkBizId, rawString)
			searchResults.Hits = append(searchResults.Hits, sr)
		}
	}

	keyMap := make(map[string]int64)
	notFoundKey := make(map[string]int64)
	// set aggregations
	bkObjIdAggr, found := result.Aggregations.Terms(common.BkObjIdAggName)
	if found == true && bkObjIdAggr != nil {
		for _, bucket := range bkObjIdAggr.Buckets {
			agg := Aggregation{}
			agg.setAgg(bucket)
			searchResults.Aggregations = append(searchResults.Aggregations, agg)
			keyMap[util.GetStrByInterface(agg.Key)] = agg.Count
		}
	}

	typeAggr, found := result.Aggregations.Terms(common.TypeAggName)
	if found == true && typeAggr != nil {
		for _, bucket := range typeAggr.Buckets {
			// only cc_HostBase, cc_ApplicationBase currently
			if bucket.Key == common.BKTableNameBaseHost || bucket.Key == common.BKTableNameBaseApp {
				agg := Aggregation{}
				agg.setAgg(bucket)
				searchResults.Aggregations = append(searchResults.Aggregations, agg)
				keyMap[util.GetStrByInterface(agg.Key)] = agg.Count
			}
		}
	}

	// fix aggregation data incomplete problem
	for _, hit := range searchResults.Hits {
		if val, ok := hit.Source[common.BKObjIDField]; ok == true {
			objID := util.GetStrByInterface(val)
			if _, exist := keyMap[objID]; exist == false {
				if _, ok := notFoundKey[objID]; ok == false {
					notFoundKey[objID] = 0
				}
				notFoundKey[objID]++
			}
		}
	}
	for key, count := range notFoundKey {
		agg := Aggregation{
			Key:   key,
			Count: count,
		}
		searchResults.Aggregations = append(searchResults.Aggregations, agg)
	}
	return searchResults, nil
}

func (query Query) checkQueryString() (string, bool) {
	// if query string is single string in SpecialChar, make it to null string
	for i := range common.SpecialChar {
		if query.QueryString == common.SpecialChar[i] {
			query.QueryString = ""
			return "", true
		}
	}
	rawString := strings.Replace(query.QueryString, "*", "", -1)
	// judge string if large than 32
	if utf8.RuneCountInString(rawString) > 32 {
		return rawString, false
	} else {
		return rawString, true
	}
}

func (query Query) toEsQueryAndSearchTypes() (elastic.Query, []string) {
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

	// biz name
	resourcePool := elastic.NewMatchQuery(common.BKAppNameField, "资源池")
	qBool.MustNot(resourcePool)

	// ignore bk_supplier_account
	qSupplierMatch := elastic.NewMatchQuery(common.BkSupplierAccount, query.QueryString)
	qBool.MustNot(qSupplierMatch)

	// if set bk_obj_id
	qString := elastic.NewQueryStringQuery(query.QueryString)
	if query.BkObjId == "" {
		// get search types from filter
		indexTypes := getEsIndexTypes(query.TypeFilter)
		// add search cc_ApplicationBase type
		indexTypes = append(indexTypes, common.BKTableNameBaseApp)
		return qBool.Must(qString), indexTypes
	} else if query.BkObjId == common.TypeHost {
		// if bk_obj_id is host, we search only from type cc_HostBase
		indexTypes := []string{common.BKTableNameBaseHost}
		return qBool.Must(qString), indexTypes
	} else if query.BkObjId == common.TypeApplication {
		// if bk_obj_id is biz, we search only from type cc_ApplicationBase
		indexTypes := []string{common.BKTableNameBaseApp}
		return qBool.Must(qString), indexTypes
	} else {
		// if define bk_obj_id, we use bool query include must(bk_obj_id=xxx) and should(query string)
		qBool.Must(elastic.NewTermQuery("bk_obj_id", query.BkObjId))
		qBool.Must(qString)
		indexTypes := getEsIndexTypes(query.TypeFilter)
		return qBool, indexTypes
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

	indexTypes := make([]string, 0)
	for _, value := range common.CmdbFindTypes {
		if !inTypes(value, typesMap) {
			indexTypes = append(indexTypes, value)
		}
	}

	return indexTypes
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
	} else if bucket.Key == common.BKTableNameBaseApp {
		agg.Key = common.TypeApplication
	} else {
		agg.Key = bucket.Key
	}

	agg.Count = bucket.DocCount
}

func (sr *SearchResult) setHit(ctx context.Context, searchHit *elastic.SearchHit, bkBizId, rawString string) {
	rid := util.ExtractRequestIDFromContext(ctx)
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

	// sr.Highlight = searchHit.Highlight
	err := json.Unmarshal(*searchHit.Source, &(sr.Source))
	if err != nil {
		blog.Warnf("full_text_find unmarshal search result source err: %+v, rid: %s", err, rid)
		sr.Source = nil
	}

	sr.dealHighlight(sr.Source, searchHit.Highlight, bkBizId, rawString)
}

func (sr *SearchResult) dealHighlight(source map[string]interface{}, highlight elastic.SearchHitHighlight, bkBizId, rawString string) {

	isObject := true
	var bkObjId, oldHighlightObjId string
	if _, ok := source["bk_obj_id"]; ok {
		bkObjId = source["bk_obj_id"].(string)
		oldHighlightObjId = "<em>" + bkObjId + "</em>"
	} else {
		isObject = false
	}
	oldHighlightBizId := "<em>" + bkBizId + "</em>"

	for key, values := range highlight {
		if key == "bk_obj_id" || key == "bk_obj_id.keyword" {
			// judge if raw query string in bk_obj_id, if not, ignore bk_obj_id highlight
			rawStringInObjId := false
			for _, value := range values {
				if strings.Contains(value, rawString) {
					rawStringInObjId = true
					break
				} else {
					continue
				}
			}
			if !rawStringInObjId {
				delete(highlight, key)
			}
		} else if key == "metadata.label.bk_biz_id" || key == "metadata.label.bk_biz_id.keyword" {
			delete(highlight, key)
		} else {
			// we don't need highlight with bk_obj_id and bk_biz_id, just like <em>bk_obj_id</em>, <em>bk_biz_id</em>
			// replace it <em>bk_obj_id</em> be bk_obj_id (do not need <em>)
			for i := range values {
				if isObject && strings.Contains(values[i], oldHighlightObjId) {
					values[i] = strings.Replace(values[i], oldHighlightObjId, bkObjId, -1)
				}
				if strings.Contains(values[i], oldHighlightBizId) {
					values[i] = strings.Replace(values[i], oldHighlightBizId, bkBizId, -1)
				}
			}
		}
	}

	sr.Highlight = highlight
}
