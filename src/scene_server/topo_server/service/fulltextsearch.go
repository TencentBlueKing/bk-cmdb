package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
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

	TypeApplication = "biz"
	TypeHost        = "host"
	TypeInstance    = "instance"
	TypeModel       = "model"

	KindModel    = "model"
	KindInstance = "instance"

	BkBizMetaKey = "metadata.label.bk_biz_id"
)

var (
	SpecialChar = []string{"`", "~", "!", "@", "#", "$", "%", "^", "&", "*",
		"(", ")", "-", "_", "=", "+", "[", "{", "]", "}",
		"\\", "|", ";", ":", "'", "\"", ",", "<", ".", ">", "/", "?"}
)

// esSpecialCharactersRegex can match es speical characters which need to be escaped
var esSpecialCharactersRegex = regexp.MustCompile(`([+\-=&|><(){}\[\]\^"~'?!:*\/])`)

// all es index name
var (
	esBizIndex      = getESIndexByCollection(common.BKTableNameBaseApp)
	esHostIndex     = getESIndexByCollection(common.BKTableNameBaseHost)
	esInstanceIndex = getESIndexByCollection(common.BKTableNameBaseInst)
	esModelIndex    = getESIndexByCollection(common.BKTableNameObjDes)
)

var esIndexesNameTypeMap = map[string]string{
	esBizIndex:      TypeApplication,
	esHostIndex:     TypeHost,
	esInstanceIndex: TypeInstance,
	esModelIndex:    TypeModel,
}

var esIndexesTypeNameMap = map[string]string{
	TypeApplication: esBizIndex,
	TypeHost:        esHostIndex,
	// the wildcard * can match multiple indexes whose name start with "cc_ObjectBase"
	TypeInstance: esInstanceIndex + "*",
	TypeModel:    esModelIndex,
}

var esIndexesTypeKindMap = map[string]string{
	TypeApplication: KindInstance,
	TypeHost:        KindInstance,
	TypeInstance:    KindInstance,
	TypeModel:       KindModel,
}

type SearchResult struct {
	Source    map[string]interface{} `json:"source"` // data from mongo, key/value
	Highlight map[string][]string    `json:"highlight"`
	Type      string                 `json:"type"` // object, host, process, model
	Score     float64                `json:"score"`
}

type Aggregation struct {
	Key string `json:"key"`
	// Kind value only can be model or instance
	Kind  string `json:"kind"`
	Count int64  `json:"count"`
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
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "query_string"))
		return
	}

	// check the format of query string and adjust it when needed
	rawString, ok := query.checkQueryString()
	if !ok {
		blog.Errorf("full_text_find failed, query string [%s] large than 32, rid: %s", rawString, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "query_string"))
		return
	}
	// get es query indexs
	indexs, err := query.getQueryIndexes()
	if err != nil {
		blog.Errorf("full_text_find failed, get query indexes failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, err.Error()+", filter"))
		return
	}

	// get es query
	esQuery := query.toEsBoolQuery()

	result, err := s.Es.Search(ctx.Kit.Ctx, esQuery, indexs, query.Paging.Start, query.Paging.Limit)
	if err != nil {
		blog.Errorf("full_text_find failed, es search failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}

	plainQuery, _ := esQuery.Source()
	plainResult, _ := json.Marshal(result)
	blog.V(4).Infof("es search, esQuery:%s, indexs:%v, start:%d, limit:%d, result:%s", plainQuery, indexs,
		query.Paging.Start, query.Paging.Limit, plainResult)

	// result is hits and aggregations
	searchResults := new(SearchResults)

	searchResults.Total = result.Hits.TotalHits.Value
	// set hits
	for _, hit := range result.Hits.Hits {
		sr := SearchResult{}
		sr.setHit(ctx.Kit.Ctx, hit, query.BkBizId, rawString)
		searchResults.Hits = append(searchResults.Hits, sr)
	}

	setAggregationInfo(searchResults)

	ctx.RespEntity(searchResults)
}

// setAggregationInfo set aggregation info for result
func setAggregationInfo(searchResults *SearchResults) error {
	bizAggregation := &Aggregation{
		Key:   TypeApplication,
		Kind:  KindInstance,
		Count: 0,
	}
	hostAggregation := &Aggregation{
		Key:   TypeHost,
		Kind:  KindInstance,
		Count: 0,
	}
	modelAggregation := &Aggregation{
		Key:   TypeModel,
		Kind:  KindModel,
		Count: 0,
	}

	instanceAggregationMap := map[string]*Aggregation{}

	// count for different type
	for _, hit := range searchResults.Hits {
		switch hit.Type {
		case TypeApplication:
			bizAggregation.Count++
		case TypeHost:
			hostAggregation.Count++
		case TypeModel:
			modelAggregation.Count++
		case TypeInstance:
			if val, ok := hit.Source[common.BKObjIDField]; ok == true {
				objID := util.GetStrByInterface(val)
				if _, ok := instanceAggregationMap[objID]; ok == false {
					instanceAggregationMap[objID] = &Aggregation{
						Key:   objID,
						Kind:  KindInstance,
						Count: 0,
					}
				}
				instanceAggregationMap[objID].Count++
			}
		default:
			blog.Warnf("unsupported hit type:%s", hit.Type)
		}
	}

	if modelAggregation.Count > 0 {
		searchResults.Aggregations = append(searchResults.Aggregations, *modelAggregation)
	}
	if bizAggregation.Count > 0 {
		searchResults.Aggregations = append(searchResults.Aggregations, *bizAggregation)
	}
	if hostAggregation.Count > 0 {
		searchResults.Aggregations = append(searchResults.Aggregations, *hostAggregation)
	}
	for _, aggr := range instanceAggregationMap {
		searchResults.Aggregations = append(searchResults.Aggregations, *aggr)
	}

	return nil
}

// checkQueryString check the format of query string and adjust it when needed
func (query *Query) checkQueryString() (string, bool) {
	// if query string is single string in SpecialChar, make it to null string
	if len(query.QueryString) == 1 {
		for i := range SpecialChar {
			if query.QueryString == SpecialChar[i] {
				query.QueryString = ""
				return "", true
			}
		}
	}

	query.QueryString = strings.Trim(query.QueryString, "*")
	rawString := query.QueryString

	// escape special characters
	query.QueryString = "*" + esSpecialCharactersRegex.ReplaceAllString(query.QueryString, `\$1`) + "*"

	// judge string if large than 32
	if utf8.RuneCountInString(rawString) > 32 {
		return rawString, false
	} else {
		return rawString, true
	}
}

// getQueryIndexes get needed indexes to query es
func (query Query) getQueryIndexes() ([]string, error) {
	queryIndexes := []string{}
	if len(query.TypeFilter) == 0 {
		queryIndexes = []string{getESIndexByCollection("cc_*")}
	} else {
		for _, indexType := range query.TypeFilter {
			indexName, ok := esIndexesTypeNameMap[indexType]
			if !ok {
				return nil, fmt.Errorf("unsupoprted fitler type:%s, must be one of biz, host, instance, model",
					indexType)
			}
			queryIndexes = append(queryIndexes, indexName)
		}
	}

	return queryIndexes, nil
}

func (query Query) toEsBoolQuery() elastic.Query {
	qBool := elastic.NewBoolQuery()

	if query.BkBizId != "" {
		qBizTerm := elastic.NewTermQuery(common.BKAppIDField, query.BkBizId)
		qBool.Should(qBizTerm)
	}

	// biz name
	resourcePool := elastic.NewMatchQuery(common.BKAppNameField, "资源池")
	qBool.MustNot(resourcePool)

	// ignore bk_supplier_account
	qSupplierMatch := elastic.NewMatchQuery(common.BkSupplierAccount, query.QueryString)
	qBool.MustNot(qSupplierMatch)

	qString := elastic.NewQueryStringQuery(query.QueryString)
	qBool.Must(qString)

	if query.BkObjId != "" {
		qBool.Must(elastic.NewTermQuery("bk_obj_id", query.BkObjId))

	}

	return qBool
}

func (agg *Aggregation) setAgg(bucket *elastic.AggregationBucketKeyItem) {
	if key, ok := bucket.Key.(string); ok {
		agg.Key = esIndexesNameTypeMap[key]
		agg.Kind = esIndexesTypeKindMap[agg.Key]
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

	sr.Type = esIndexesNameTypeMap[searchHit.Index]
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
		} else if key == "bk_biz_id" || key == "bk_biz_id.keyword" {
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
