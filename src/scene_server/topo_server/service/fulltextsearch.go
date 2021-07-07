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
	// Source mongodb metadata.
	Source interface{} `json:"source"`

	// Kind data kind model or instance.
	Kind string `json:"kind"`

	// Key data model key biz/set/module/host/{common object id}.
	Key string `json:"key"`
}

// Aggregation fulltext search aggregation.
type Aggregation struct {
	// Kind data kind model or instance.
	Kind string `json:"kind"`

	// Key data model key biz/set/module/host/{common object id}.
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

// FullTextSearchESQuery fulltext search elastic query.
type FullTextSearchESQuery struct {
	// Query elastic query.
	Query elastic.Query

	// Condition elastic search condition.
	Condition *FullTextSearchCondition
}

// FullTextSearchCondition fulltext search condition.
type FullTextSearchCondition struct {
	// IndexName es index name.
	IndexName string

	// Conditions es search conditions.
	Conditions map[string]interface{}
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
		return fmt.Errorf("can't search with the special character[%s]", r.QueryString)
	}

	// escape special characters.
	rawString := strings.Trim(r.QueryString, "*")
	r.QueryString = "*" + esSpecialCharactersRegex.ReplaceAllString(rawString, `\$1`) + "*"

	// check query_string length in UTF-8 encoding.
	utf8Length := utf8.RuneCountInString(rawString)

	if utf8Length > esQueryStringLengthLimit {
		return fmt.Errorf("invalid search string[%s], length[%d] in UTF-8 encoding is too large, max[%d]",
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

func (r *FullTextSearchReq) generateModelQueryConditions() []*FullTextSearchCondition {
	var searchConditions []*FullTextSearchCondition

	for _, model := range r.Filter.Models {
		searchConditions = append(searchConditions, &FullTextSearchCondition{
			IndexName: metadata.IndexNameModel,
			Conditions: map[string]interface{}{
				metadata.IndexPropertyBKObjID: model,
			},
		})
	}

	return searchConditions
}

func (r *FullTextSearchReq) generateInstanceQueryConditions() []*FullTextSearchCondition {
	var searchConditions []*FullTextSearchCondition

	for _, instance := range r.Filter.Instances {
		switch instance {
		case common.BKInnerObjIDApp:
			searchConditions = append(searchConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameBiz,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDApp,
				},
			})

		case common.BKInnerObjIDSet:
			searchConditions = append(searchConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameSet,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDSet,
				},
			})

		case common.BKInnerObjIDModule:
			searchConditions = append(searchConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameModule,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDModule,
				},
			})

		case common.BKInnerObjIDHost:
			searchConditions = append(searchConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameHost,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: common.BKInnerObjIDHost,
				},
			})

		default:
			searchConditions = append(searchConditions, &FullTextSearchCondition{
				IndexName: metadata.IndexNameObjectInstance,
				Conditions: map[string]interface{}{
					metadata.IndexPropertyBKObjID: instance,
				},
			})
		}
	}

	return searchConditions
}

// GenerateESQuery returns the elastic query for main search and sub-count searches.
func (r *FullTextSearchReq) GenerateESQuery() (elastic.Query, []string, []*FullTextSearchESQuery) {
	// elastic query conditions.
	conditions := r.generateModelQueryConditions()
	conditions = append(conditions, r.generateInstanceQueryConditions()...)

	// build elastic bool query.
	var indexes []string
	var subCountQueries []*FullTextSearchESQuery

	indexMap := make(map[string]struct{})
	queryConditions := make(map[string][]interface{})

	query := elastic.NewBoolQuery()
	if len(r.OwnerID) != 0 {
		query.Should(elastic.NewMatchQuery(metadata.IndexPropertyBKSupplierAccount, r.OwnerID))
	}
	if len(r.BizID) != 0 {
		query.Should(elastic.NewMatchQuery(metadata.IndexPropertyBKBizID, r.BizID))
	}
	query.Must(elastic.NewQueryStringQuery(r.QueryString))

	for _, condition := range conditions {
		boolQuery := elastic.NewBoolQuery()

		// handle conditions.
		if len(r.OwnerID) != 0 {
			boolQuery.Should(elastic.NewMatchQuery(metadata.IndexPropertyBKSupplierAccount, r.OwnerID))
		}
		if len(r.BizID) != 0 {
			boolQuery.Should(elastic.NewMatchQuery(metadata.IndexPropertyBKBizID, r.BizID))
		}

		// handle filter conditions.
		for property, value := range condition.Conditions {
			queryConditions[property] = append(queryConditions[property], value)
			boolQuery.Should(elastic.NewMatchQuery(property, value))
		}

		// build elastic search.
		if _, exist := indexMap[condition.IndexName]; !exist {
			indexes = append(indexes, condition.IndexName)
			indexMap[condition.IndexName] = struct{}{}
		}

		// handle query string.
		boolQuery.Must(elastic.NewQueryStringQuery(r.QueryString))

		subCountQueries = append(subCountQueries, &FullTextSearchESQuery{Query: boolQuery, Condition: condition})
	}

	for property, value := range queryConditions {
		query.Should(elastic.NewMatchQuery(property, value))
	}

	return query, indexes, subCountQueries
}

func (s *Service) fullTextAggregationSearch(ctx *rest.Contexts,
	esQueries []*FullTextSearchESQuery) ([]Aggregation, error) {

	// elastic pipeline sub-count search.
	var pipelineErr error
	var wg sync.WaitGroup

	pipeline := make(chan struct{}, esPipelineConcurrency)

	// pipeline results.
	aggregationQueryResults := make([]Aggregation, len(esQueries))

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
			aggregationQueryResults[idx] = aggregation

		}(ctx, idx, query)
	}

	// wait for searches done.
	wg.Wait()

	if pipelineErr != nil {
		return nil, pipelineErr
	}

	return aggregationQueryResults, nil
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

	// build response data.
	response := FullTextSearchResp{}
	response.Total = mainSearchResult.Hits.TotalHits.Value

	if mainSearchResult.Hits.TotalHits.Value == 0 {
		ctx.RespEntity(response)
		return
	}

	// sub-count aggregation search.
	aggregations, err := s.fullTextAggregationSearch(ctx, subCountQueries)
	if err != nil {
		blog.Errorf("fulltext sub-count search failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoFullTextFindErr))
		return
	}
	response.Aggregations = aggregations

	// query result metadata.
	modelMetadataConditions := []string{}
	instMetadataConditions := make(map[string][]string)

	// set read preference.
	ctx.SetReadPreference(common.SecondaryPreferredMode)

	for _, hit := range mainSearchResult.Hits.Hits {
		source := make(map[string]interface{})
		if err := json.Unmarshal(hit.Source, &source); err != nil {
			blog.Warnf("fulltext handle search result source data failed, err: %+v,  rid: %s", err, ctx.Kit.Rid)
			continue
		}

		metaID := util.GetStrByInterface(source[metadata.IndexPropertyID])
		objectID := util.GetStrByInterface(source[metadata.IndexPropertyBKObjID])

		if hit.Type == metadata.DataKindModel {
			modelMetadataConditions = append(modelMetadataConditions, objectID)
		} else if hit.Type == metadata.DataKindInstance {
			instMetadataConditions[objectID] = append(instMetadataConditions[objectID], metaID)
		} else {
			blog.Warnf("fulltext handle search source, unknown data kind[%s], rid: %s", hit.Type, ctx.Kit.Rid)
		}
	}

	// query model.
	modelCondition := condition.CreateCondition()
	modelCondition.Field(common.BKObjIDField).In(modelMetadataConditions)

	objects, err := s.Core.ObjectOperation().FindObject(ctx.Kit, modelCondition)
	if err != nil {
		blog.Errorf("fulltext search model failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	for _, object := range objects {
		response.Hits = append(response.Hits, SearchResult{
			Kind:   metadata.DataKindModel,
			Key:    object.Object().ObjectID,
			Source: object.Object(),
		})
	}

	// query instance.
	for objectID, ids := range instMetadataConditions {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "_id", Operator: querybuilder.OperatorIn, Value: ids},
					},
				},
			},
			Page: metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit},
		}

		// search object instances.
		result, err := s.Core.InstOperation().SearchObjectInstances(ctx.Kit, objectID, input)
		if err != nil {
			blog.Errorf("fulltext search object[%s] instances failed, err: %+v, rid: %s", objectID, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		for _, instance := range result.Info {
			response.Hits = append(response.Hits, SearchResult{
				Kind:   metadata.DataKindInstance,
				Key:    objectID,
				Source: instance,
			})
		}
	}

	ctx.RespEntity(response)
}
