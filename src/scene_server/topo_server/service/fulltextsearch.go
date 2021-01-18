package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"github.com/olivere/elastic/v7"
)

const (
	ESIndexPrefix = "cmdb"

	TypeHost        = "host"
	TypeObject      = "object"
	TypeApplication = "biz"

	BkBizMetaKey = "metadata.label.bk_biz_id"
)

var (
	SpecialChar = []string{"`", "~", "!", "@", "#", "$", "%", "^", "&", "*",
		"(", ")", "-", "_", "=", "+", "[", "{", "]", "}",
		"\\", "|", ";", ":", "'", "\"", ",", "<", ".", ">", "/", "?"}
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

func (s *Service) FullTextFind(ctx *rest.Contexts) {
	if s.Es.Client == nil {
		blog.Errorf("FullTextFind failed, es client is nil, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextClientNotInitialized))
		return
	}

	query := NewQuery()
	if err := ctx.DecodeInto(query); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if query.QueryString == "" {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	// check query string
	rawString, ok := query.checkQueryString()
	if !ok {
		blog.Errorf("full_text_find failed, query string [%s] large than 32, rid: %s", rawString, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "query_string"))
		return
	}

	// get query and search indexs
	esQuery, indexs := query.toEsBoolQueryAndIndexs()

	result, err := s.Es.Search(ctx.Kit.Ctx, esQuery, indexs, query.Paging.Start, query.Paging.Limit)
	if err != nil {
		blog.Errorf("full_text_find failed, es search failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}

	// result is hits and aggregations
	searchResults := new(SearchResults)

	searchResults.Total = result.Hits.TotalHits.Value
	// set hits
	for _, hit := range result.Hits.Hits {
		sr := SearchResult{}
		sr.setHit(ctx.Kit.Ctx, hit, query.BkBizId, rawString)
		searchResults.Hits = append(searchResults.Hits, sr)
	}

	keyMap := make(map[string]int64)
	notFoundKey := make(map[string]int64)
	indexArggr, found := result.Aggregations.Terms(common.IndexAggName)
	if found == true && indexArggr != nil {
		for _, bucket := range indexArggr.Buckets {
			// only cc_HostBase, cc_ApplicationBase currently
			if bucket.Key == getESIndexByCollection(common.BKTableNameBaseHost) ||
				bucket.Key == getESIndexByCollection(common.BKTableNameBaseApp) {
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
	ctx.RespEntity(searchResults)
}

func (query Query) checkQueryString() (string, bool) {
	// if query string is single string in SpecialChar, make it to null string
	for i := range SpecialChar {
		if query.QueryString == SpecialChar[i] {
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

func (query Query) toEsBoolQueryAndIndexs() (elastic.Query, []string) {
	qBool := elastic.NewBoolQuery()

	// if set bk_biz_id
	qBool.MinimumNumberShouldMatch(1)
	qBizRegex := elastic.NewRegexpQuery(BkBizMetaKey, "[0-9]*")
	qBizBool := elastic.NewBoolQuery()
	qBizBool.MustNot(qBizRegex)
	qBool.Should(qBizBool)
	if query.BkBizId != "" {
		qBizTerm := elastic.NewTermQuery(BkBizMetaKey, query.BkBizId)
		qBool.Should(qBizTerm)
	}

	// biz name
	resourcePool := elastic.NewMatchQuery(common.BKAppNameField, "资源池")
	qBool.MustNot(resourcePool)

	// ignore bk_supplier_account
	qSupplierMatch := elastic.NewMatchQuery(common.BkSupplierAccount, query.QueryString)
	qBool.MustNot(qSupplierMatch)

	qString := elastic.NewQueryStringQuery(query.QueryString)
	qString.Escape(true)
	qBool.Must(qString)

	if query.BkObjId == "" {
		// if bk_obj_id is "", we search all indexs
		indexs := make([]string, 0)
		indexs = append(indexs, getESIndexByCollection(common.BKTableNameBaseApp))
		indexs = append(indexs, getESIndexByCollection(common.BKTableNameBaseHost))
		indexs = append(indexs, getESIndexByCollection(common.BKTableNameBaseInst))
		return qBool, indexs
	} else if query.BkObjId == TypeHost {
		// if bk_obj_id is host, we search only from type cc_HostBase
		indexs := []string{getESIndexByCollection(common.BKTableNameBaseHost)}
		return qBool.Must(qString), indexs
	} else if query.BkObjId == TypeApplication {
		// if bk_obj_id is biz, we search only from type cc_ApplicationBase
		indexs := []string{getESIndexByCollection(common.BKTableNameBaseApp)}
		return qBool.Must(qString), indexs
	} else {
		// if define bk_obj_id, we use bool query include must(bk_obj_id=xxx) and should(query string)
		qBool.Must(elastic.NewTermQuery("bk_obj_id", query.BkObjId))
		qBool.Must(qString)
		indexs := []string{getESIndexByCollection(common.BKTableNameBaseInst)}
		return qBool, indexs
	}
}

func (agg *Aggregation) setAgg(bucket *elastic.AggregationBucketKeyItem) {
	if bucket.Key == getESIndexByCollection(common.BKTableNameBaseHost) {
		agg.Key = TypeHost
	} else if bucket.Key == getESIndexByCollection(common.BKTableNameBaseApp) {
		agg.Key = TypeApplication
	} else {
		agg.Key = bucket.Key
	}

	agg.Count = bucket.DocCount
}

func (sr *SearchResult) setHit(ctx context.Context, searchHit *elastic.SearchHit, bkBizId, rawString string) {
	rid := util.ExtractRequestIDFromContext(ctx)
	sr.Score = *searchHit.Score

	// sr.Highlight = searchHit.Highlight
	err := json.Unmarshal(searchHit.Source, &(sr.Source))
	if err != nil {
		blog.Warnf("full_text_find unmarshal search result source err: %+v, rid: %s", err, rid)
		sr.Source = nil
	}

	switch searchHit.Index {
	case getESIndexByCollection(common.BKTableNameBaseApp):
		sr.Type = TypeApplication
	case getESIndexByCollection(common.BKTableNameBaseHost):
		sr.Type = TypeHost
	case getESIndexByCollection(common.BKTableNameBaseInst):
		sr.Type = TypeObject
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

// getESIndexByCollection get the index of es through ESIndexPrefix and collection's name
func getESIndexByCollection(collectionName string) string {
	collectionName = strings.ToLower(collectionName)
	return fmt.Sprintf("%s.%s", strings.ToLower(ESIndexPrefix), collectionName)
}
