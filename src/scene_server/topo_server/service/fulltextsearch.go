/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"

	"github.com/olivere/elastic/v7"
)

var (
	// specialCharacters query_string keyword special characters, character -> flag.
	specialCharacters = map[string]bool{
		"`":  true,
		"~":  true,
		"!":  true,
		"@":  true,
		"#":  true,
		"$":  true,
		"%":  true,
		"^":  true,
		"&":  true,
		"*":  true,
		"(":  true,
		")":  true,
		"-":  true,
		"_":  true,
		"=":  true,
		"+":  true,
		"[":  true,
		"{":  true,
		"]":  true,
		"}":  true,
		"\\": true,
		"\"": true,
		"|":  true,
		";":  true,
		":":  true,
		"'":  true,
		",":  true,
		".":  true,
		"<":  true,
		">":  true,
		"/":  true,
		"?":  true,
	}

	// esSpecialCharactersRegex elastic special characters regexp.
	esSpecialCharactersRegex = regexp.MustCompile(`([+\-=&|><(){}\[\]\^"~'?!:*\/])`)

	// esPipelineConcurrency elastic pipeline query concurrency.
	esPipelineConcurrency = 50

	// esQueryStringLengthLimit query_string length limit(utf-8).
	esQueryStringLengthLimit = 32
)

// SearchResult fulltext search result.
type SearchResult struct {
	// Kind data kind model or instance.
	Kind string `json:"kind"`

	// Key data model key biz/set/module/host/{common object id}.
	Key string `json:"key"`

	// Source mongodb metadata.
	Source interface{} `json:"source"`

	// Highlight highlight keywords
	Highlight map[string][]string `json:"highlight"`
}

// Aggregation fulltext search aggregation.
type Aggregation struct {
	// Kind data kind model or instance.
	Kind string `json:"kind"`

	// Key data model key or instance biz/set/module/host/{common object id}.
	Key string `json:"key"`

	// Count hits data count.
	Count int64 `json:"count"`
}

// FullTextSearchResp is fulltext search response.
type FullTextSearchResp struct {
	// Total total number.
	Total int64 `json:"total"`

	// Aggregations fulltext search aggregations
	Aggregations []Aggregation `json:"aggregations"`

	// Hits search result.
	Hits []SearchResult `json:"hits"`
}

// Page search page settings.
type Page struct {
	// Start search start.
	Start int `json:"start"`

	// Limit search limit.
	Limit int `json:"limit"`
}

// Validate validate fulltext search request page settings.
func (p *Page) Validate() error {
	if p == nil {
		return errors.New("page is missing")
	}

	if p.Start < 0 {
		return errors.New("page start must great than or equal to 0")
	}

	if p.Limit <= 0 {
		return errors.New("page limit must great than 0")
	}
	if p.Limit > 100 {
		return errors.New("page limit must less than or equal to 100")
	}

	return nil
}

// FullTextSearchCondition fulltext search condition.
type FullTextSearchCondition struct {
	// IndexName es index name.
	IndexName string

	// Conditions es search conditions.
	Conditions map[string]interface{}
}

// FullTextSearchESQuery fulltext search elastic query.
type FullTextSearchESQuery struct {
	// Query elastic query.
	Query elastic.Query

	// Condition elastic search condition.
	Condition *FullTextSearchCondition
}

// FullTextSearchFilter is fulltext filter.
type FullTextSearchFilter struct {
	// Models model filter. e.g. biz/set/module/host/{common object id}.
	Models []string `json:"models"`

	// Instances model filter. e.g. biz/set/module/host/{common object id}.
	Instances []string `json:"instances"`
}

// Validate validate the fulltext search request filter.
func (f *FullTextSearchFilter) Validate() error {
	if f == nil {
		return errors.New("filter is missing")
	}

	// validate models.
	modelsMap := make(map[string]struct{})
	for _, model := range f.Models {
		if _, exist := modelsMap[model]; exist {
			return fmt.Errorf("repeated model[%s]", model)
		}
		modelsMap[model] = struct{}{}
	}

	// validate instances.
	instancesMap := make(map[string]struct{})
	for _, instance := range f.Instances {
		if _, exist := instancesMap[instance]; exist {
			return fmt.Errorf("repeated instance[%s]", instance)
		}
		instancesMap[instance] = struct{}{}
	}

	if len(f.Models) == 0 && len(f.Instances) == 0 {
		return errors.New("empty models and instances filter")
	}

	return nil
}

// FullTextSearchReq is fulltext search request.
type FullTextSearchReq struct {
	// OwnerID supplier account.
	OwnerID string `json:"bk_supplier_account"`

	// BizID business id.
	BizID string `json:"bk_biz_id"`

	// Filter search filter.
	Filter *FullTextSearchFilter `json:"filter"`

	// Fields search object filter
	Fields []string `json:"fields"`

	// SubResource main search list
	SubResource *FullTextSearchFilter `json:"sub_resource"`
	// QueryString elastic query_string keyword.
	QueryString string `json:"query_string"`

	// Page search page settings.
	Page *Page `json:"page"`
}

// Validate validate the fulltext search request.
func (r *FullTextSearchReq) Validate() error {
	if len(r.QueryString) == 0 {
		return errors.New("can't search with the empty keyword")
	}

	// check single special character.
	if specialCharacters[r.QueryString] {
		return fmt.Errorf("can't search with the special character: %s", r.QueryString)
	}

	// escape special characters.
	rawString := strings.Trim(r.QueryString, "*")
	r.QueryString = "*" + esSpecialCharactersRegex.ReplaceAllString(rawString, `\$1`) + "*"

	// check query_string length in UTF-8 encoding.
	utf8Length := utf8.RuneCountInString(rawString)

	if utf8Length > esQueryStringLengthLimit {
		return fmt.Errorf("invalid search string[%s], length[%d] in UTF-8 encoding is too large, max: %d",
			rawString, utf8Length, esQueryStringLengthLimit)
	}

	// check filter.
	if err := r.Filter.Validate(); err != nil {
		return fmt.Errorf("invalid search request filter, %+v", err)
	}

	// check page.
	if err := r.Page.Validate(); err != nil {
		return fmt.Errorf("invalid search request page, %+v", err)
	}

	return nil
}

// generateESQueryConditions parse and handle models/instances filter, generate main elastic
// query, sub-count aggregations query.
func (r *FullTextSearchFilter) generateESQueryConditions() []*FullTextSearchCondition {

	searchFilterConditions := make([]*FullTextSearchCondition, 0)
	for _, model := range r.Models {
		searchFilterConditions = append(searchFilterConditions, &FullTextSearchCondition{
			IndexName: metadata.IndexNameModel,
			Conditions: map[string]interface{}{
				metadata.IndexPropertyBKObjID: model,
			},
		})
	}

	for _, instance := range r.Instances {
		switch instance {
		case common.BKInnerObjIDBizSet:
			searchFilterConditions = append(searchFilterConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameBizSet,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDBizSet,
				},
			})

		case common.BKInnerObjIDApp:
			searchFilterConditions = append(searchFilterConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameBiz,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDApp,
				},
			})

		case common.BKInnerObjIDSet:
			searchFilterConditions = append(searchFilterConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameSet,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDSet,
				},
			})

		case common.BKInnerObjIDModule:
			searchFilterConditions = append(searchFilterConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameModule,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDModule,
				},
			})

		case common.BKInnerObjIDHost:
			searchFilterConditions = append(searchFilterConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameHost,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDHost,
				},
			})

		default:
			searchFilterConditions = append(searchFilterConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameObjectInstance,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: instance,
				},
			})
		}
	}

	return searchFilterConditions
}

// GenerateESQuery returns the elastic query for main search and sub-count searches.
func (r *FullTextSearchReq) GenerateESQuery() (elastic.Query, []string, []*FullTextSearchESQuery) {
	// elastic query conditions for each model or instance.
	var (
		// main search use objIdCond firstly.
		bObjIdFlag bool

		// objCond if necessary.
		subResource []*FullTextSearchCondition
	)

	// build elastic count aggregations conditions for search.
	subCountQueries := make([]*FullTextSearchESQuery, 0)

	// build elastic main query and indexes for search.
	indexes := make([]string, 0)

	filterCond := r.Filter.generateESQueryConditions()

	if r.SubResource != nil {
		subResource = r.SubResource.generateESQueryConditions()
	}
	indexMap := make(map[string]struct{})

	// main query.
	query := elastic.NewBoolQuery()
	queryConditions := make(map[string][]interface{})
	if len(r.OwnerID) != 0 {
		query.Must(elastic.NewMatchQuery(metadata.IndexPropertyBKSupplierAccount, r.OwnerID))
	}
	if len(r.BizID) != 0 {
		query.Must(elastic.NewMatchQuery(metadata.IndexPropertyBKBizID, r.BizID))
	}
	query.Must(elastic.NewQueryStringQuery(r.QueryString).Field(metadata.IndexPropertyKeywords))

	//  main search select objIdCond firstly.
	if len(subResource) > 0 {
		for _, cond := range subResource {
			// build elastic main query condition.
			for property, value := range cond.Conditions {
				queryConditions[property] = append(queryConditions[property], value)
			}
			// build elastic main query indexes.
			if _, exist := indexMap[cond.IndexName]; !exist {
				indexes = append(indexes, cond.IndexName)
				indexMap[cond.IndexName] = struct{}{}
			}
		}
		bObjIdFlag = true
	}

	// sub aggregations query.
	for _, condFilter := range filterCond {
		boolQuery := elastic.NewBoolQuery()
		if len(r.OwnerID) != 0 {
			boolQuery.Must(elastic.NewMatchQuery(metadata.IndexPropertyBKSupplierAccount, r.OwnerID))
		}
		if len(r.BizID) != 0 {
			boolQuery.Must(elastic.NewMatchQuery(metadata.IndexPropertyBKBizID, r.BizID))
		}

		// handle filter conditions.
		for property, value := range condFilter.Conditions {
			if !bObjIdFlag {
				//  when objId is nil  main query use filterCond
				queryConditions[property] = append(queryConditions[property], value)
			}
			boolQuery.Must(elastic.NewMatchQuery(property, value))
		}

		// if objId is nil use condFilter build main query indexes.
		if !bObjIdFlag {
			if _, exist := indexMap[condFilter.IndexName]; !exist {
				indexes = append(indexes, condFilter.IndexName)
				indexMap[condFilter.IndexName] = struct{}{}
			}
		}

		// assign sub aggregation conditions and field.
		boolQuery.Must(elastic.NewQueryStringQuery(r.QueryString).Field(metadata.IndexPropertyKeywords))
		subCountQueries = append(subCountQueries, &FullTextSearchESQuery{Query: boolQuery, Condition: condFilter})
	}

	// add all sub aggregation conditions to main query.
	for property, value := range queryConditions {
		query.Must(elastic.NewTermsQuery(property, value...))
	}

	return query, indexes, subCountQueries
}

// fullTextAggregation count aggregations in multi goroutines mode.
func (s *Service) fullTextAggregation(ctx *rest.Contexts, esQueries []*FullTextSearchESQuery) ([]Aggregation, error) {
	// elastic pipeline sub aggregation search.
	var (
		pipelineErr error
		wg          sync.WaitGroup
	)

	// control max gcoroutines num.
	pipeline := make(chan struct{}, esPipelineConcurrency)
	// pipeline results.
	aggregationQueryTmp := make([]Aggregation, len(esQueries))
	aggregationQueryResults := make([]Aggregation, 0)

	// search with multi goroutines.
	for idx, query := range esQueries {
		// try to start one search.
		pipeline <- struct{}{}
		wg.Add(1)

		// start one search gcoroutine.
		go func(ctx *rest.Contexts, idx int, esQuery *FullTextSearchESQuery) {
			defer func() {
				// one search gcoroutine done.
				wg.Done()
				<-pipeline
			}()

			count, err := s.Es.Count(ctx.Kit.Ctx, esQuery.Query, []string{esQuery.Condition.IndexName})
			if err != nil {
				blog.Errorf("fulltext search count failed,query cond: %s err: %+v, rid: %s", esQuery.Query, err, ctx.Kit.Rid)
				pipelineErr = err
				return
			}

			aggregation := Aggregation{
				Kind:  metadata.DataKindInstance,
				Key:   util.GetStrByInterface(esQuery.Condition.Conditions[metadata.IndexPropertyBKObjID]),
				Count: count,
			}
			if esQuery.Condition.IndexName == metadata.IndexNameModel {
				aggregation.Kind = metadata.DataKindModel
			}
			aggregationQueryTmp[idx] = aggregation

		}(ctx, idx, query)
	}

	// wait for searches done.
	wg.Wait()

	if pipelineErr != nil {
		return nil, pipelineErr
	}
	for _, v := range aggregationQueryTmp {
		if v.Count != 0 {
			aggregationQueryResults = append(aggregationQueryResults, v)
		}
	}

	return aggregationQueryResults, nil
}

// fullTextMetadata returns metadata base on the elastic hits.
func (s *Service) fullTextMetadata(ctx *rest.Contexts, hits []*elastic.SearchHit, request FullTextSearchReq) (
	[]SearchResult, error) {

	// for meta_bk_obj_id in es.
	objectIDs := make([]string, 0)
	searchResults := make([]SearchResult, 0)

	instMetadataConditions := make(map[string][]int64)

	// search the highlight fields for instance.
	insHits := make(map[string]map[int64]*elastic.SearchHit)

	// search the highlight fields for model.
	objHits := make(map[string]*elastic.SearchHit)

	for _, hit := range hits {
		source := make(map[string]interface{})
		if err := json.Unmarshal(hit.Source, &source); err != nil {
			blog.Warnf("fulltext handle search result source data failed, err: %+v,  rid: %s", err, ctx.Kit.Rid)
			continue
		}
		objectID := util.GetStrByInterface(source[metadata.IndexPropertyBKObjID])
		dataKind := util.GetStrByInterface(source[metadata.IndexPropertyDataKind])
		metaID, err := strconv.ParseInt(util.GetStrByInterface(source[metadata.IndexPropertyID]), 10, 64)
		if err != nil {
			blog.Errorf(" query meta data fail,objectID[%s],err=[%v] rid: %s", objectID, err, ctx.Kit.Rid)
			continue
		}

		if dataKind == metadata.DataKindModel {
			objectIDs = append(objectIDs, objectID)
			objHits[objectID] = hit
		} else if dataKind == metadata.DataKindInstance {
			instMetadataConditions[objectID] = append(instMetadataConditions[objectID], metaID)
			if insHits[objectID] == nil {
				insHits[objectID] = make(map[int64]*elastic.SearchHit)
			}
			insHits[objectID][metaID] = hit

		} else {
			blog.Errorf("fulltext handle search source, unknown data kind: %s, rid: %s", dataKind, ctx.Kit.Rid)
		}
	}

	blog.V(5).Infof("fulltext metadata query models: %s, instances: %s, rid: %s",
		objectIDs, instMetadataConditions, ctx.Kit.Rid)
	// set read preference.
	ctx.SetReadPreference(common.SecondaryPreferredMode)

	// query metadata object.
	searchObjectResults := s.fullTextSearchForObject(ctx, objectIDs, objHits, request)
	searchResults = append(searchResults, searchObjectResults...)

	// query metadata instance.
	searchInstanceResults := s.fullTextSearchForInstance(ctx, instMetadataConditions, insHits, request)
	searchResults = append(searchResults, searchInstanceResults...)
	return searchResults, nil
}

// initCommonSearchFilter init common search filter cond
func initCommonSearchFilter(field string, ids []int64) *metadata.CommonSearchFilter {
	return &metadata.CommonSearchFilter{
		Conditions: &querybuilder.QueryFilter{
			Rule: &querybuilder.AtomRule{
				Field:    field,
				Operator: querybuilder.OperatorIn,
				Value:    ids,
			},
		},
		Page: metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit},
	}
}

// fullTextSearchForInstanceCond composition query instance condition.
func fullTextSearchForInstanceCond(objectID string, ids []int64) *metadata.CommonSearchFilter {
	input := &metadata.CommonSearchFilter{}
	switch objectID {
	case common.BKInnerObjIDBizSet:
		input = initCommonSearchFilter(common.BKBizSetIDField, ids)
	case common.BKInnerObjIDApp:
		input = initCommonSearchFilter(common.BKAppIDField, ids)
	case common.BKInnerObjIDSet:
		input = initCommonSearchFilter(common.BKSetIDField, ids)
	case common.BKInnerObjIDModule:
		input = initCommonSearchFilter(common.BKModuleIDField, ids)
	case common.BKInnerObjIDHost:
		input = initCommonSearchFilter(common.BKHostIDField, ids)
	default:
		input = initCommonSearchFilter(common.BKInstIDField, ids)
	}
	return input
}

// fullTextSearchForInstance search instance result.
func (s *Service) fullTextSearchForInstance(ctx *rest.Contexts, instMetadataConditions map[string][]int64,
	insHits map[string]map[int64]*elastic.SearchHit, request FullTextSearchReq) []SearchResult {

	searchResults := make([]SearchResult, 0)
	if len(instMetadataConditions) == 0 {
		blog.Errorf("inst metadata cond is invalid, cond: %+v, rid: %s", instMetadataConditions, ctx.Kit.Rid)
		return nil
	}

	// query metadata instance.
	input := &metadata.CommonSearchFilter{}
	var (
		wg       sync.WaitGroup
		rwLock   sync.RWMutex
		firstErr error
	)
	pipeline := make(chan bool, 10)
	for objectID, ids := range instMetadataConditions {
		pipeline <- true
		wg.Add(1)
		go func(objectID string, ids []int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			input = fullTextSearchForInstanceCond(objectID, ids)
			// search object instances.
			result, err := s.Logics.InstOperation().SearchObjectInstances(ctx.Kit, objectID, input)
			if err != nil {
				blog.Errorf("search obj instances fail, objID: %s, ids: %v, rid: %s", objectID, ids, ctx.Kit.Rid)
				firstErr = err
				return
			}

			for _, instance := range result.Info {

				if _, ok := instance.(*mapstr.MapStr); !ok {
					blog.Errorf("get inst struct fail, objectID: %v, rid: %s", objectID, ctx.Kit.Rid)
					continue
				}
				var idStr string
				inst := instance.(*mapstr.MapStr)
				switch objectID {
				case common.BKInnerObjIDBizSet:
					idStr, err = inst.String(common.BKBizSetIDField)
				case common.BKInnerObjIDApp:
					idStr, err = inst.String(common.BKAppIDField)
				case common.BKInnerObjIDSet:
					idStr, err = inst.String(common.BKSetIDField)
				case common.BKInnerObjIDModule:
					idStr, _ = inst.String(common.BKModuleIDField)
				case common.BKInnerObjIDHost:
					idStr, err = inst.String(common.BKHostIDField)
				default:
					idStr, err = inst.String(common.BKInstIDField)
				}
				if err != nil {
					blog.Errorf("get instId fail, objectID: %v,rid: %s", objectID, ctx.Kit.Rid)
					continue
				}

				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					blog.Errorf("parse instId fail, objectID: %v, rid: %s", objectID, ctx.Kit.Rid)
					continue
				}

				// instance result
				searchRes := SearchResult{}
				rawString := strings.Trim(request.QueryString, "*")
				searchRes.setHit(ctx.Kit.Ctx, insHits[objectID][id], request.BizID, rawString)
				searchRes.Kind = metadata.DataKindInstance
				searchRes.Key = objectID
				searchRes.Source = instance
				rwLock.Lock()
				searchResults = append(searchResults, searchRes)
				rwLock.Unlock()
			}

		}(objectID, ids)
	}
	wg.Wait()

	if firstErr != nil {
		return nil
	}
	return searchResults
}

// fullTextSearchForObject search object result.
func (s *Service) fullTextSearchForObject(ctx *rest.Contexts, objectIDs []string,
	objHits map[string]*elastic.SearchHit, request FullTextSearchReq) []SearchResult {

	modelCondition := &metadata.QueryCondition{
		Fields:         request.Fields,
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		Condition:      mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objectIDs}},
		DisableCounter: true,
	}

	objects, err := s.Engine.CoreAPI.CoreService().Model().ReadModel(ctx.Kit.Ctx, ctx.Kit.Header, modelCondition)
	if err != nil {
		blog.Errorf("get objects(%+v) failed, err: %v, rid: %s", objectIDs, err, ctx.Kit.Rid)
		return nil
	}

	if objects.Count == 0 {
		blog.Errorf("meet the modelCond object is empty, modelCond: %+v, rid: %s", modelCondition, ctx.Kit.Rid)
		return nil
	}

	searchResults := make([]SearchResult, 0)
	// model result.
	for _, object := range objects.Info {
		searchRes := SearchResult{}
		rawString := strings.Trim(request.QueryString, "*")
		searchRes.setHit(ctx.Kit.Ctx, objHits[object.ObjectID], request.BizID, rawString)
		searchRes.Kind = metadata.DataKindModel
		searchRes.Key = object.ObjectID
		searchRes.Source = object
		searchResults = append(searchResults, searchRes)
	}

	return searchResults
}

// FullTextSearch fulltext search service.
func (s *Service) FullTextSearch(ctx *rest.Contexts) {
	// check elastic client.
	if s.Es.Client == nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextClientNotInitialized))
		return
	}
	// decode input parameter.
	request := FullTextSearchReq{}
	if err := ctx.DecodeInto(&request); err != nil {
		ctx.RespAutoError(err)
		return
	}
	// validate request data.
	if err := request.Validate(); err != nil {
		blog.Errorf("validate fulltext search input parameters failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, err.Error()))
		return
	}

	// generate elastic query.
	esQuery, indexes, subCountQueries := request.GenerateESQuery()

	mainESQuery, err := esQuery.Source()
	if err != nil {
		blog.Errorf("fulltext parse mainESQuery fail: err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}
	blog.V(5).Infof("fulltext main query[%s], indexes[%s], rid: %s", mainESQuery, indexes, ctx.Kit.Rid)

	// main search.
	mainSearchResult, err := s.Es.Search(ctx.Kit.Ctx, esQuery, indexes, request.Page.Start, request.Page.Limit)
	if err != nil {
		blog.Errorf("fulltext main search failed,mainESQuery: %s err: %+v, rid: %s", mainESQuery, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}

	if mainSearchResult.Hits == nil || mainSearchResult.Hits.TotalHits == nil {
		blog.Errorf("fulltext main search failed, invalid search result, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}

	// aggregation search.
	aggregations, err := s.fullTextAggregation(ctx, subCountQueries)
	if err != nil {
		blog.Errorf("fulltext sub-count search failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}

	var total int64
	for _, agg := range aggregations {
		total += agg.Count
	}
	// build response data.
	// when objId is not nil,mainSearchResult.Hits.TotalHits.Value is inaccurate ,
	// so we must use sum of each subCountQueries result
	response := FullTextSearchResp{}
	if mainSearchResult.Hits.TotalHits.Value == 0 {
		ctx.RespEntity(response)
		return
	}
	response = FullTextSearchResp{Total: total}
	response.Aggregations = aggregations

	// metadata search.
	metadata, err := s.fullTextMetadata(ctx, mainSearchResult.Hits.Hits, request)
	if err != nil {
		blog.Errorf("fulltext metadata search failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}
	response.Hits = metadata
	ctx.RespEntity(response)
	return
}

// setHit get highlight words.
func (sr *SearchResult) setHit(ctx context.Context, searchHit *elastic.SearchHit, bkBizId, rawString string) {

	if searchHit == nil {
		return
	}

	rid := util.ExtractRequestIDFromContext(ctx)
	sourceTmp := make(map[string]interface{})

	err := json.Unmarshal(searchHit.Source, &(sourceTmp))
	if err != nil {
		blog.Warnf("full_text_find unmarshal search result source err: %+v, rid: %s", err, rid)
		sr.Source = nil
		return
	}

	sr.dealHighlight(sourceTmp, searchHit.Highlight, bkBizId, rawString)
	return
}

// dealHighlight 此函数会处理掉一些不需要展示出来的内部关系id，防止高亮出一些用户原本不希望高亮的字段
func (sr *SearchResult) dealHighlight(source map[string]interface{}, highlight elastic.SearchHitHighlight,
	bkBizId, rawString string) {

	isObject := true
	var bkObjId, oldHighlightObjId string
	if _, ok := source[metadata.IndexPropertyBKObjID]; ok {
		bkObjId = source[metadata.IndexPropertyBKObjID].(string)
		oldHighlightObjId = "<em>" + bkObjId + "</em>"
	} else {
		isObject = false
	}
	oldHighlightBizId := "<em>" + bkBizId + "</em>"

	inputKey := strings.ToLower(rawString)

	for key, values := range highlight {

		if key == metadata.IndexPropertyBKObjID {
			// judge if raw query string in bk_obj_id, if not, ignore bk_obj_id highlight
			rawStringInObjId := false
			for _, value := range values {
				v := strings.ToLower(value)
				if strings.Contains(v, inputKey) {
					rawStringInObjId = true
					break
				} else {
					continue
				}
			}
			if !rawStringInObjId {
				delete(highlight, key)
			}
		} else {
			// we don't need highlight with meta_bk_obj_id and meta_bk_biz_id, just like <em>meta_bk_obj_id</em>,
			// <em>meta_bk_biz_id</em>. Replace it <em>meta_bk_obj_id</em> be bk_obj_id (do not need <em>)
			for i := range values {
				vLower := strings.ToLower(values[i])
				if isObject && strings.Contains(vLower, oldHighlightObjId) && !strings.Contains(vLower, inputKey) {
					values[i] = strings.Replace(vLower, oldHighlightObjId, bkObjId, -1)
				}
				if strings.Contains(vLower, oldHighlightBizId) {
					values[i] = strings.Replace(vLower, oldHighlightBizId, bkBizId, -1)
				}
			}
		}
	}
	sr.Highlight = highlight
	return
}
