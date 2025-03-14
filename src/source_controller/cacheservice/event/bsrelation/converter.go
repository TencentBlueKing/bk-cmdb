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

package bsrelation

import (
	"context"
	"reflect"
	"sync"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo/sharding"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// needCareBizFields struct to get id->type map of need cared biz fields that will result in biz set relation changes
type needCareBizFields struct {
	fieldMap map[string]map[string]string
	lock     sync.RWMutex
}

// Get get need cared biz fields
func (a *needCareBizFields) Get() map[string]map[string]string {
	a.lock.RLock()
	defer a.lock.RUnlock()
	fieldMap := make(map[string]map[string]string, len(a.fieldMap))
	for key, value := range a.fieldMap {
		fieldMap[key] = value
	}
	return fieldMap
}

// Set set need cared biz fields
func (a *needCareBizFields) Set(tenantID string, fieldMap map[string]string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.fieldMap[tenantID] = fieldMap
}

// syncNeedCareBizFields refresh need cared biz fields every minutes
func (b *bizSetRelation) syncNeedCareBizFields() {
	for {
		time.Sleep(time.Minute)

		err := tenant.ExecForAllTenants(func(tenantID string) error {
			fields, err := b.getNeedCareBizFields(context.Background(), tenantID)
			if err != nil {
				blog.Errorf("run biz set relation watch, but get need care biz fields failed, err: %v", err)
				return err
			}
			b.needCareBizFields.Set(tenantID, fields)
			blog.V(5).Infof("run biz set relation watch, sync tenant %s need care biz fields done, fields: %+v",
				tenantID, fields)
			return nil
		})
		if err != nil {
			continue
		}
	}
}

// getNeedCareBizFields get need cared biz fields, including biz id and enum/organization type fields
func (b *bizSetRelation) getNeedCareBizFields(ctx context.Context, tenantID string) (map[string]string, error) {
	filter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDApp,
		common.BKPropertyTypeField: map[string]interface{}{
			common.BKDBIN: []string{common.FieldTypeEnum, common.FieldTypeOrganization},
		},
	}

	attributes := make([]metadata.Attribute, 0)
	err := mongodb.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(common.BKTableNameObjAttDes).Find(filter).
		Fields(common.BKPropertyIDField, common.BKPropertyTypeField).All(ctx, &attributes)
	if err != nil {
		blog.Errorf("get need care biz attributes failed, filter: %+v, err: %v", filter, err)
		return nil, err
	}

	fieldMap := map[string]string{common.BKAppIDField: common.FieldTypeInt}
	for _, attr := range attributes {
		fieldMap[attr.PropertyID] = attr.PropertyType
	}
	return fieldMap, nil
}

// rearrangeBizSetEvents biz set events rearrange policy:
//  1. If update event's updated fields do not contain "bk_scope" field, we will drop this event.
//  2. Aggregate multiple same biz set's events to one event, so that we can decrease the amount of biz set relation
//     events. Because we only care about which biz set's relation is changed, one event is enough for us.
func (b *bizSetRelation) rearrangeBizSetEvents(es []*types.Event, rid string) ([]*types.Event, error) {

	// get last biz set event type, in order to rearrange biz set events later, policy:
	// 2. create + update: event is changed to create event with the update event detail
	// 3. update + delete: event is changed to delete event
	lastEventMap := make(map[string]*types.Event)

	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		if _, exists := lastEventMap[genUniqueKey(e)]; !exists {
			lastEventMap[genUniqueKey(e)] = e
		}
	}

	hitEvents := make([]*types.Event, 0)
	// remind if a biz set events has already been hit, if yes, then skip this event.
	reminder := make(map[string]struct{})

	for _, one := range es {
		if _, yes := reminder[genUniqueKey(one)]; yes {
			// this biz set event is already hit, then we aggregate the event with the former ones to only one event.
			// this is useful to decrease biz set relation events.
			blog.Infof("biz set event: %s is aggregated, rid: %s", one.ID(), rid)
			continue
		}

		// get the hit event(nil means no event is hit) and if other events of the biz set needs to be ignored
		hitEvent, needIgnored := b.parseBizSetEvent(one, lastEventMap, rid)

		if hitEvent != nil {
			hitEvents = append(hitEvents, hitEvent)
		}

		if needIgnored {
			reminder[genUniqueKey(one)] = struct{}{}
		}
	}

	// refresh all biz ids cache if the events contains match all biz set, the cache is used to generate detail later
	matchAllTenantIDs := make([]string, 0)
	for _, e := range hitEvents {
		if e.OperationType != types.Delete && gjson.Get(string(e.DocBytes), "bk_scope.match_all").Bool() {
			matchAllTenantIDs = append(matchAllTenantIDs, e.TenantID)
			break
		}

		matchAllTenantIDs = util.StrArrayUnique(matchAllTenantIDs)
		for _, tenantID := range matchAllTenantIDs {
			if err := b.refreshAllBizIDStr(tenantID, rid); err != nil {
				return nil, err
			}
		}
	}

	return hitEvents, nil
}

// parseBizSetEvent parse biz set event, returns the hit event and if same oid events needs to be skipped
func (b *bizSetRelation) parseBizSetEvent(one *types.Event, lastEventMap map[string]*types.Event, rid string) (
	*types.Event, bool) {

	switch one.OperationType {
	case types.Insert:
		// if events are in the order of create + update + delete, all events of this biz set will be ignored
		lastEvent := lastEventMap[genUniqueKey(one)]
		if lastEvent.OperationType == types.Delete {
			return nil, true
		}

		// if events are in the order of create + update, it is changed to create event with the update event detail
		if lastEvent.OperationType == types.Update {
			one.DocBytes = lastEvent.DocBytes
			one.Document = lastEvent.Document
		}

		// insert event is added directly.
		return one, true
	case types.Delete:
		// delete event is added directly.
		return one, true
	case types.Update, types.Replace:
		// if events are in the order of update + delete, it is changed to delete event, update event is ignored
		if lastEventMap[genUniqueKey(one)].OperationType == types.Delete {
			return nil, false
		}

		// replace event's change description is empty, we do not know what is changed, so we add it directly.
		if len(one.ChangeDesc.UpdatedFields) == 0 && len(one.ChangeDesc.RemovedFields) == 0 {
			return one, true
		}

		// check if updated/removed fields contains "bk_scope" field, only add the events that changed scope.
		if len(one.ChangeDesc.UpdatedFields) != 0 {
			if _, exists := one.ChangeDesc.UpdatedFields[common.BKBizSetScopeField]; exists {
				return one, true
			}
		}

		if len(one.ChangeDesc.RemovedFields) != 0 {
			for _, field := range one.ChangeDesc.RemovedFields {
				if field == common.BKBizSetScopeField {
					return one, true
				}
			}
		}

		// event that needs to be ignored comes here.
		blog.Infof("biz set relation event, biz set change detail do not care, skip, oid: %s, detail: %+v, rid: %s",
			one.ID(), one.ChangeDesc, rid)
	default:
		blog.Errorf("biz set event operation type(%s) is invalid, event: %+v, rid: %s", one.OperationType, one, rid)
		return nil, false
	}
	return nil, false
}

// rearrangeBizEvents biz events rearrange policy:
//  1. Biz event is redirected to its related biz sets' events by traversing all biz sets and checking if the "bk_scope"
//     field matches the biz's attribute.
//  2. Create and delete event's related biz set is judged by whether its scope contains the biz.
//  3. Update event's related biz set is judged by whether its scope contains the updated fields of the event. Since
//     we can't get the previous value of the updated fields, we can't get the exact biz sets it was in before.
//  4. Aggregate multiple biz events with the same biz set to one event.
func (b *bizSetRelation) rearrangeBizEvents(es []*types.Event, rid string) ([]*types.Event, error) {
	// parse biz events to convert biz event parameter
	params := b.parseBizEvents(es, rid)

	// convert biz event to related biz set event whose relation is affected by biz event
	bizSetEvents, err := b.convertBizEvent(params, rid)
	if err != nil {
		return nil, err
	}

	return bizSetEvents, nil
}

// parseBizEvents parse biz events to insert/delete event details and updated fields, used to generate biz set events
func (b *bizSetRelation) parseBizEvents(es []*types.Event, rid string) convertBizEventParams {

	// generate convert biz event to biz set event parameter
	params := convertBizEventParams{
		es:                            es,
		insertedAndDeletedBiz:         make(map[string][]map[string]interface{}),
		insertedAndDeletedBizIndexMap: make(map[string]map[int]int),
		updatedFieldsIndexMap:         make(map[string]map[string]int),
		needCareFieldsMap:             b.needCareBizFields.Get(),
	}

	// parse biz event into inserted/deleted details, and updated fields to match biz sets whose relation changed.
	for index, one := range es {
		tenantID := one.TenantID
		switch one.OperationType {
		case types.Insert, types.Delete:
			// use document in insert event to find matching biz sets.
			params.convertInsertedAndDeletedBizEvent(tenantID, one.Document, index)

		case types.Update, types.Replace:
			// replace event's change description is empty, treat as updating all fields.
			if len(one.ChangeDesc.UpdatedFields) == 0 && len(one.ChangeDesc.RemovedFields) == 0 {
				for field := range params.needCareFieldsMap[tenantID] {
					params.convertUpdateBizEvent(tenantID, field, index)
				}
				continue
			}

			// get the updated/removed fields that need to be cared, ignores the rest.
			isIgnored := true
			updateFields := make([]string, 0)
			if len(one.ChangeDesc.UpdatedFields) > 0 {
				// for biz set relation,archive biz is treated as delete event while recover is treated as insert event
				if one.ChangeDesc.UpdatedFields[common.BKDataStatusField] == string(common.DataStatusDisabled) ||
					one.ChangeDesc.UpdatedFields[common.BKDataStatusField] == string(common.DataStatusEnable) {
					params.convertInsertedAndDeletedBizEvent(tenantID, one.Document, index)
					continue
				}

				for field := range one.ChangeDesc.UpdatedFields {
					updateFields = append(updateFields, field)
				}
			}

			for _, field := range append(updateFields, one.ChangeDesc.RemovedFields...) {
				if _, exists := params.needCareFieldsMap[tenantID][field]; exists {
					params.convertUpdateBizEvent(tenantID, field, index)
					isIgnored = false
				}
			}

			if isIgnored {
				blog.Infof("biz set relation event, biz change detail do not care, skip, oid: %s, detail: %+v, rid: %s",
					one.ID(), one.ChangeDesc, rid)
			}

		default:
			blog.Errorf("biz event operation type(%s) is invalid, event: %+v, rid: %s", one.OperationType, one, rid)
		}
	}

	return params
}

// convertBizEventParams parameter for convertBizEvent function
type convertBizEventParams struct {
	// es biz events
	es []*types.Event
	// insertedAndDeletedBiz inserted and deleted biz event details, handled in the same way
	insertedAndDeletedBiz map[string][]map[string]interface{}
	// insertedAndDeletedBizIndexMap mapping of insertedAndDeletedBiz index to es index, used to locate origin event
	insertedAndDeletedBizIndexMap map[string]map[int]int
	// updatedFieldsIndexMap mapping of updated fields to the first update event index, used to locate origin event
	updatedFieldsIndexMap map[string]map[string]int
	// needCareFieldsMap need care biz fields that can be used in biz set scope
	needCareFieldsMap map[string]map[string]string
}

func (params *convertBizEventParams) convertInsertedAndDeletedBizEvent(tenantID string, doc interface{}, index int) {
	detail, ok := doc.(*map[string]interface{})
	if !ok || detail == nil {
		return
	}

	params.insertedAndDeletedBiz[tenantID] = append(params.insertedAndDeletedBiz[tenantID], *detail)
	if _, exists := params.insertedAndDeletedBizIndexMap[tenantID]; !exists {
		params.insertedAndDeletedBizIndexMap[tenantID] = make(map[int]int)
	}
	params.insertedAndDeletedBizIndexMap[tenantID][len(params.insertedAndDeletedBiz)-1] = index
}

func (params *convertBizEventParams) convertUpdateBizEvent(tenantID string, field string, index int) {
	if _, exists := params.updatedFieldsIndexMap[tenantID]; !exists {
		params.updatedFieldsIndexMap[tenantID] = make(map[string]int)
	}
	if _, exists := params.updatedFieldsIndexMap[tenantID][field]; !exists {
		params.updatedFieldsIndexMap[tenantID][field] = index
	}
}

// convertBizEvent convert biz event to related biz set event whose relation is affected by biz event
func (b *bizSetRelation) convertBizEvent(params convertBizEventParams, rid string) ([]*types.Event, error) {
	bizSetEvents := make([]*types.Event, 0)

	if len(params.insertedAndDeletedBiz) == 0 && len(params.updatedFieldsIndexMap) == 0 {
		return bizSetEvents, nil
	}

	// get all biz sets to check if their scope matches the biz events, because filtering them in db is too complicated
	tenantIDs := make([]string, 0)
	for tenantID := range params.insertedAndDeletedBiz {
		tenantIDs = append(tenantIDs, tenantID)
	}
	for tenantID := range params.updatedFieldsIndexMap {
		tenantIDs = append(tenantIDs, tenantID)
	}
	bizSetsMap, err := b.getAllBizSets(tenantIDs, rid)
	if err != nil {
		return nil, err
	}

	// get biz events' related biz sets
	relatedBizSets, containsMatchAllBizSetTenants := b.getRelatedBizSets(params, bizSetsMap, rid)

	// refresh all biz ids cache if the events affected match all biz set, the cache is used to generate detail later
	for _, tenantID := range containsMatchAllBizSetTenants {
		if err := b.refreshAllBizIDStr(tenantID, rid); err != nil {
			return nil, err
		}
	}

	// parse biz sets into biz set events
	for bizEventIndex, bizSetArr := range relatedBizSets {
		bizEvent := params.es[bizEventIndex]

		for _, bizSet := range bizSetArr {
			doc, err := json.Marshal(bizSet.BizSetInst)
			if err != nil {
				blog.Errorf("marshal biz set(%+v) failed, err: %v, rid: %s", bizSet.BizSetInst, err, rid)
				return nil, err
			}

			bizSetEvents = append(bizSetEvents, &types.Event{
				Oid:           bizSet.Oid.Hex(),
				Document:      bizSet.BizSetInst,
				DocBytes:      doc,
				OperationType: types.Update,
				CollectionInfo: types.CollectionInfo{
					Collection: common.GenTenantTableName(bizEvent.TenantID, common.BKTableNameBaseBizSet),
					TenantID:   bizEvent.TenantID,
					ParsedColl: common.BKTableNameBaseBizSet,
				},
				ClusterTime: bizEvent.ClusterTime,
				Token:       bizEvent.Token,
			})
		}
	}

	return bizSetEvents, nil
}

// getRelatedBizSets get biz events index to related biz sets map, which is used to generate biz set events
func (b *bizSetRelation) getRelatedBizSets(params convertBizEventParams, bizSetsMap map[string][]bizSetWithOid,
	rid string) (map[int][]bizSetWithOid, []string) {

	containsMatchAllBizSetTenants := make([]string, 0)
	relatedBizSets := make(map[int][]bizSetWithOid, 0)
	for tenantID, bizSets := range bizSetsMap {
		for _, bizSet := range bizSets {
			// for biz set that matches all biz, only insert and delete event will affect their relations
			if bizSet.Scope.MatchAll {
				if len(params.insertedAndDeletedBiz[tenantID]) > 0 {
					eventIndex := params.insertedAndDeletedBizIndexMap[tenantID][0]
					relatedBizSets[eventIndex] = append(relatedBizSets[eventIndex], bizSet)
					containsMatchAllBizSetTenants = append(containsMatchAllBizSetTenants, tenantID)
				}
				continue
			}

			if bizSet.Scope.Filter == nil {
				blog.Errorf("biz set(%+v) scope filter is empty, skip, rid: %s", bizSet, rid)
				continue
			}

			firstEventIndex, matched := b.getFirstMatchedEvent(params, bizSet, tenantID, rid)
			if matched {
				relatedBizSets[firstEventIndex] = append(relatedBizSets[firstEventIndex], bizSet)
			}
		}
	}

	return relatedBizSets, containsMatchAllBizSetTenants
}

func (b *bizSetRelation) getFirstMatchedEvent(params convertBizEventParams, bizSet bizSetWithOid, tenantID string,
	rid string) (int, bool) {

	var firstEventIndex int

	// update biz event matches all biz sets whose scope contains the updated fields, get all matching fields.
	matched := bizSet.Scope.Filter.MatchAny(func(r querybuilder.AtomRule) bool {
		updatedFieldsIndexMap, exists := params.updatedFieldsIndexMap[tenantID]
		if !exists {
			return false
		}
		if index, exists := updatedFieldsIndexMap[r.Field]; exists {
			if firstEventIndex == 0 || index < firstEventIndex {
				firstEventIndex = index
			}
			return true
		}
		return false
	})

	// check if biz set scope filter matches the inserted/removed biz
	for index, biz := range params.insertedAndDeletedBiz[tenantID] {
		// if the event index already exceeds the event index of matched update fields, stop checking
		eventIndex := params.insertedAndDeletedBizIndexMap[tenantID][index]
		if firstEventIndex != 0 && eventIndex >= firstEventIndex {
			break
		}

		bizMatched := bizSet.Scope.Filter.Match(func(r querybuilder.AtomRule) bool {
			// ignores the biz set filter rule that do not contain need care fields
			propertyType, exists := params.needCareFieldsMap[tenantID][r.Field]
			if !exists {
				blog.Errorf("biz set(%+v) filter rule contains ignored field, rid: %s", bizSet, rid)
				return false
			}

			// ignores the biz that do not contain the field in filter rule
			bizVal, exists := biz[r.Field]
			if !exists {
				blog.Infof("biz(%+v) do not contain rule field %s, rid: %s", biz, r.Field, rid)
				return false
			}

			switch r.Operator {
			case querybuilder.OperatorEqual:
				return matchEqualOper(r.Value, bizVal, propertyType, rid)
			case querybuilder.OperatorIn:
				return matchInOper(r.Value, bizVal, propertyType, rid)
			default:
				blog.Errorf("biz set(%+v) filter rule contains invalid operator, rid: %s", bizSet, rid)
				return false
			}
		})

		if bizMatched {
			firstEventIndex = eventIndex
			matched = bizMatched
			break
		}
	}
	return firstEventIndex, matched
}

// matchEqualOper check if biz set scope filter rule with equal operator matches biz value
func matchEqualOper(ruleVal, bizVal interface{}, propertyType string, rid string) bool {
	switch propertyType {
	case common.FieldTypeEnum:
		ruleValStr, ok := ruleVal.(string)
		if !ok {
			blog.Errorf("enum type field rule value(%+v) is not string type, rid: %s", ruleVal, rid)
			return false
		}

		bizValStr, ok := bizVal.(string)
		if !ok {
			blog.Errorf("enum type field biz value(%+v) is not string type, rid: %s", bizVal, rid)
			return false
		}

		if ruleValStr == bizValStr {
			return true
		}
		return false
	case common.FieldTypeInt, common.FieldTypeOrganization:
		ruleValInt, err := util.GetIntByInterface(ruleVal)
		if err != nil {
			blog.Errorf("parse rule value(%+v) to int failed, rid: %s", ruleVal, err, rid)
			return false
		}

		bizValInt, err := util.GetIntByInterface(bizVal)
		if err != nil {
			blog.Errorf("parse biz value(%+v) to int failed, rid: %s", bizVal, err, rid)
			return false
		}

		if ruleValInt == bizValInt {
			return true
		}
		return false
	default:
		blog.Errorf("rule filed type(%s) is invalid, rid: %s", propertyType, rid)
		return false
	}
}

// matchInOper check if biz set scope filter rule with in operator matches biz value
func matchInOper(ruleVal, bizVal interface{}, propertyType string, rid string) bool {
	switch reflect.TypeOf(ruleVal).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		blog.Errorf("rule value(%+v) type is invalid, rid: %s", ruleVal, rid)
		return false
	}

	ruleValArr := reflect.ValueOf(ruleVal)
	ruleValLen := ruleValArr.Len()
	for i := 0; i < ruleValLen; i++ {
		// check if any of the rule value matches biz value
		matched := matchEqualOper(ruleValArr.Index(i).Interface(), bizVal, propertyType, rid)
		if matched {
			return true
		}
	}
	return false
}

// bizSetWithOid biz set struct with oid
type bizSetWithOid struct {
	Oid                 primitive.ObjectID `bson:"_id"`
	metadata.BizSetInst `bson:",inline"`
}

func (b *bizSetRelation) getAllBizSets(tenantIDs []string, rid string) (map[string][]bizSetWithOid, error) {
	const step = 500

	bizSetsMap := make(map[string][]bizSetWithOid, 0)

	findOpts := dbtypes.NewFindOpts().SetWithObjectID(true)

	for _, tenantID := range tenantIDs {
		kit := rest.NewKit().WithTenant(tenantID).WithRid(rid)
		cond := map[string]interface{}{}

		for {
			oneStep := make([]bizSetWithOid, 0)
			err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseBizSet).Find(cond, findOpts).
				Fields(common.BKBizSetIDField, common.BKBizSetScopeField).Limit(step).
				Sort(common.BKBizSetIDField).All(kit.Ctx, &oneStep)
			if err != nil {
				blog.Errorf("get biz set failed, err: %v, rid: %s", err, rid)
				return nil, err
			}

			bizSetsMap[tenantID] = append(bizSetsMap[tenantID], oneStep...)

			if len(oneStep) < step {
				break
			}

			cond = map[string]interface{}{
				common.BKBizSetIDField: map[string]interface{}{common.BKDBGT: oneStep[len(oneStep)-1].BizSetID},
			}
		}
	}

	return bizSetsMap, nil
}

func genUniqueKey(e *types.Event) string {
	return e.Collection + "-" + e.Oid
}
