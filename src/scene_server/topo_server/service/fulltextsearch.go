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
	"strings"
	"sync"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"

	"github.com/olivere/elastic"
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
	if enabled, _ := specialCharacters[r.QueryString]; enabled {
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
	var (
		searchFilterConditions []*FullTextSearchCondition
	)

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
		// build elastic count aggregations conditions for search.
		subCountQueries []*FullTextSearchESQuery
		// build elastic main query and indexes for search.
		indexes []string
		// main  search condition flag
		bObjIdFlag bool
		// objCond if necessary
		subResource []*FullTextSearchCondition
	)

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
	query.Must(elastic.NewQueryStringQuery(r.QueryString))

	//ignore biz name
	resourcePool := elastic.NewMatchQuery(common.BKAppNameField, "资源池")
	query.MustNot(resourcePool)

	//  main search select objIdCond first
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

		// handle sub aggregation conditions.
		boolQuery.Must(elastic.NewQueryStringQuery(r.QueryString))
		subCountQueries = append(subCountQueries, &FullTextSearchESQuery{Query: boolQuery, Condition: condFilter})
	}

	// add all sub aggregation conditions to main query.
	for property, value := range queryConditions {
		query.Must(elastic.NewTermsQuery(property, value...))
	}

	return query, indexes, subCountQueries
}

// fullTextAggregation count aggregations in multi gcoroutines mode.
func (s *Service) fullTextAggregation(ctx *rest.Contexts, esQueries []*FullTextSearchESQuery) ([]Aggregation, error) {
	// elastic pipeline sub aggregation search.
	var (
		pipelineErr             error
		wg                      sync.WaitGroup
		aggregationQueryResults []Aggregation
	)

	// control max gcoroutines num.
	pipeline := make(chan struct{}, esPipelineConcurrency)
	// pipeline results.
	aggregationQueryTmp := make([]Aggregation, len(esQueries))

	// search in multi gcoroutines.
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
				blog.Errorf("fulltext search count failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
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
func (s *Service) fullTextMetadata(ctx *rest.Contexts, hits []*elastic.SearchHit, request FullTextSearchReq) ([]SearchResult, error) {
	var (
		objectIDs     []string
		searchResults []SearchResult
	)

	instMetadataConditions := make(map[string][]int64)

	insHights := make(map[string]map[string][]string)
	objHights := make(map[string]map[string][]string)

	// set read preference.
	ctx.SetReadPreference(common.SecondaryPreferredMode)
	// build metadata models and instances conditions.
	for _, hit := range hits {
		source := make(map[string]interface{})
		if err := json.Unmarshal(hit.Source, &source); err != nil {
			blog.Warnf("fulltext handle search result source data failed, err: %+v,  rid: %s", err, ctx.Kit.Rid)
			continue
		}
		objectID := util.GetStrByInterface(source[metadata.IndexPropertyBKObjID])
		dataKind := util.GetStrByInterface(source[metadata.IndexPropertyDataKind])
		metaID, err := util.GetInt64ByInterface(source[metadata.IndexPropertyID])
		if err != nil {
			blog.Errorf(" query meta data fail,objectID[%s],err=[%v] rid: %s", objectID, err, ctx.Kit.Rid)
			continue
		}
		for k := range hit.Highlight {
			if k == "meta_bk_obj_id" {
				delete(hit.Highlight, k)
			}
		}

		// parse meta fields.
		if dataKind == metadata.DataKindModel {
			objectIDs = append(objectIDs, objectID)
			if objHights[objectID] == nil {
				objHights[objectID] = make(map[string][]string)
			}

			objHights[objectID] = hit.Highlight
			blog.Errorf("00000000000000000 Highlight %v", hit.Highlight)

		} else if dataKind == metadata.DataKindInstance {
			instMetadataConditions[objectID] = append(instMetadataConditions[objectID], metaID)
			if insHights[objectID] == nil {
				insHights[objectID] = make(map[string][]string)
			}
			insHights[objectID] = hit.Highlight
			blog.Errorf("1111111111111111111 Highlight %v", hit.Highlight)
		} else {
			blog.Warnf("fulltext handle search source, unknown data kind[%s], rid: %s", dataKind, ctx.Kit.Rid)
		}
	}

	blog.V(5).Infof("fulltext metadata query models[%s], instances[%s], rid: %s",
		objectIDs, instMetadataConditions, ctx.Kit.Rid)

	// query metadata model.
	modelCondition := condition.CreateCondition()
	modelCondition.Field(common.BKObjIDField).In(objectIDs)
	objects, err := s.Core.ObjectOperation().FindObject(ctx.Kit, modelCondition)
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		searchRes := SearchResult{}
		searchRes.Highlight = make(map[string][]string)
		searchRes.Kind = metadata.DataKindModel
		searchRes.Key = object.Object().ObjectID
		searchRes.Source = object.Object()
		searchRes.Highlight = objHights[object.Object().ObjectID]
		searchResults = append(searchResults, searchRes)
	}

	// query metadata instance.
	input := &metadata.CommonSearchFilter{}
	for objectID, ids := range instMetadataConditions {
		switch objectID {
		case common.BKInnerObjIDApp:
			input = &metadata.CommonSearchFilter{
				Conditions: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							&querybuilder.AtomRule{
								Field:    common.BKAppIDField,
								Operator: querybuilder.OperatorIn,
								Value:    ids,
							},
						},
					},
				},
				Page: metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit},
			}
		case common.BKInnerObjIDSet:
			input = &metadata.CommonSearchFilter{
				Conditions: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							&querybuilder.AtomRule{
								Field:    common.BKSetIDField,
								Operator: querybuilder.OperatorIn,
								Value:    ids,
							},
						},
					},
				},
				Page: metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit},
			}
		case common.BKInnerObjIDModule:
			input = &metadata.CommonSearchFilter{
				Conditions: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							&querybuilder.AtomRule{
								Field:    common.BKModuleIDField,
								Operator: querybuilder.OperatorIn,
								Value:    ids,
							},
						},
					},
				},
				Page: metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit},
			}

		case common.BKInnerObjIDHost:
			input = &metadata.CommonSearchFilter{
				Conditions: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							&querybuilder.AtomRule{
								Field:    common.BKHostIDField,
								Operator: querybuilder.OperatorIn,
								Value:    ids,
							},
						},
					},
				},
				Page: metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit},
			}
		default:

			input = &metadata.CommonSearchFilter{
				Conditions: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							&querybuilder.AtomRule{
								Field:    common.BKInstIDField,
								Operator: querybuilder.OperatorIn,
								Value:    ids,
							},
						},
					},
				},
				Page:     metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit},
				ObjectID: objectID,
			}
		}

		// search object instances.
		result, err := s.Core.InstOperation().SearchObjectInstances(ctx.Kit, objectID, input)
		if err != nil {
			return nil, err
		}

		for _, instance := range result.Info {
			searchRes := SearchResult{}
			searchRes.Highlight = make(map[string][]string)
			searchRes.Kind = metadata.DataKindInstance
			searchRes.Key = objectID
			searchRes.Source = instance
			searchRes.Highlight = insHights[objectID]
			searchResults = append(searchResults, searchRes)

		}
	}

	return searchResults, nil
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

	mainESQuery, _ := esQuery.Source()
	blog.V(5).Infof("fulltext main query[%s], indexes[%s], rid: %s", mainESQuery, indexes, ctx.Kit.Rid)

	// main search.
	mainSearchResult, err := s.Es.Search(ctx.Kit.Ctx, esQuery, indexes, request.Page.Start, request.Page.Limit)
	if err != nil {
		blog.Errorf("fulltext main search failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
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
	if mainSearchResult.Hits.TotalHits.Value != 0 {
		response = FullTextSearchResp{Total: total}
	} else {
		ctx.RespEntity(response)
		return
	}
	response.Aggregations = aggregations

	// metadata search.
	metadatas, err := s.fullTextMetadata(ctx, mainSearchResult.Hits.Hits, request)
	if err != nil {
		blog.Errorf("fulltext metadata search failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}
	response.Hits = metadatas

	ctx.RespEntity(response)
	return
}

func (sr *SearchResult) setHit(ctx context.Context, searchHit *elastic.SearchHit, bkBizId, rawString string) {
	rid := util.ExtractRequestIDFromContext(ctx)
	sourceTmp := make(map[string]interface{})

	err := json.Unmarshal(searchHit.Source, &(sourceTmp))
	if err != nil {
		blog.Warnf("full_text_find unmarshal search result source err: %+v, rid: %s", err, rid)
		sr.Source = nil
	}
	blog.Errorf("000000000000000000000000000 sourceTmp: %+v,Highlight: %v", sourceTmp, searchHit.Highlight)

	sr.dealHighlight(sourceTmp, searchHit.Highlight, bkBizId, rawString)
	return
}

func (sr *SearchResult) dealHighlight(source map[string]interface{}, highlight elastic.SearchHitHighlight, bkBizId, rawString string) {

	isObject := true
	var bkObjId, oldHighlightObjId string
	if _, ok := source["meta_bk_obj_id"]; ok {
		bkObjId = source["meta_bk_obj_id"].(string)
		oldHighlightObjId = "<em>" + bkObjId + "</em>"
	} else {
		isObject = false
	}
	oldHighlightBizId := "<em>" + bkBizId + "</em>"

	blog.Errorf("1111111111111111111111 highlight: %v,oldHighlightObjId: %v,oldHighlightBizId: %v", highlight, oldHighlightObjId, oldHighlightBizId)

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
		} else if key == "meta_bk_biz_id" || key == "meta_bk_obj_id" {
			delete(highlight, key)
		} else {
			// we don't need highlight with bk_obj_id and bk_biz_id, just like <em>bk_obj_id</em>, <em>bk_biz_id</em>
			// replace it <em>bk_obj_id</em> be bk_obj_id (do not need <em>)
			for i := range values {
				blog.Errorf("2222222222222222222222222 i: %v", i)

				if isObject && strings.Contains(values[i], oldHighlightObjId) {
					blog.Errorf("33333333333333333333333 values: %v,oldHighlightObjId: %v", values[i], oldHighlightObjId)

					values[i] = strings.Replace(values[i], oldHighlightObjId, bkObjId, -1)
				}
				if strings.Contains(values[i], oldHighlightBizId) {
					blog.Errorf("4444444444444444444444444 values: %v，oldHighlightBizId: %v", values[i], oldHighlightBizId)

					values[i] = strings.Replace(values[i], oldHighlightBizId, bkBizId, -1)
				}
			}
		}
	}
	blog.Errorf("5555555555555555555 highlight: %v,oldHighlightObjId: %v,oldHighlightBizId: %v", highlight, oldHighlightObjId, oldHighlightBizId)
	sr.Highlight = highlight
}
