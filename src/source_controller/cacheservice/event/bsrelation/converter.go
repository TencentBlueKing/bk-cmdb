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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/stream/types"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// needCareBizFields struct to get id->type map of need cared biz fields that will result in biz set relation changes
type needCareBizFields struct {
	fieldMap map[string]string
	lock     sync.RWMutex
}

// Get get need cared biz fields
func (a *needCareBizFields) Get() map[string]string {
	a.lock.RLock()
	defer a.lock.RUnlock()
	fieldMap := make(map[string]string, len(a.fieldMap))
	for key, value := range a.fieldMap {
		fieldMap[key] = value
	}
	return fieldMap
}

// Set set need cared biz fields
func (a *needCareBizFields) Set(fieldMap map[string]string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.fieldMap = fieldMap
}

// syncNeedCareBizFields refresh need cared biz fields every minutes
func (b *bizSetRelation) syncNeedCareBizFields() {
	for {
		time.Sleep(time.Minute)

		fields, err := b.getNeedCareBizFields(context.Background())
		if err != nil {
			blog.Errorf("run biz set relation watch, but get need care biz fields failed, err: %v", err)
			continue
		}
		b.needCareBizFields.Set(fields)
		blog.V(5).Infof("run biz set relation watch, sync need care biz fields done, fields: %+v", fields)
	}
}

// getNeedCareBizFields get need cared biz fields, including biz id and enum/organization type fields
func (b *bizSetRelation) getNeedCareBizFields(ctx context.Context) (map[string]string, error) {
	filter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDApp,
		common.BKPropertyTypeField: map[string]interface{}{
			common.BKDBIN: []string{common.FieldTypeEnum, common.FieldTypeOrganization},
		},
	}

	attributes := make([]metadata.Attribute, 0)
	err := b.ccDB.Table(common.BKTableNameObjAttDes).Find(filter).Fields(common.BKPropertyIDField,
		common.BKPropertyTypeField).All(ctx, &attributes)
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

// rearrangeBizSetEvents TODO
// biz set events rearrange policy:
// 1. If update event's updated fields do not contain "bk_scope" field, we will drop this event.
// 2. Aggregate multiple same biz set's events to one event, so that we can decrease the amount of biz set relation
//    events. Because we only care about which biz set's relation is changed, one event is enough for us.
func (b *bizSetRelation) rearrangeBizSetEvents(es []*types.Event, rid string) ([]*types.Event, error) {

	// get last biz set event type, in order to rearrange biz set events later, policy:
	// 2. create + update: event is changed to create event with the update event detail
	// 3. update + delete: event is changed to delete event
	lastEventMap := make(map[string]*types.Event)
	deletedOids := make([]string, 0)

	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		if _, exists := lastEventMap[e.Oid]; !exists {
			lastEventMap[e.Oid] = e

			if e.OperationType == types.Delete {
				deletedOids = append(deletedOids, e.Oid)
			}
		}
	}

	// get deleted biz set detail from del archive before adding to the events.
	oidDetailMap, err := b.getDeleteEventDetails(deletedOids, rid)
	if err != nil {
		return nil, err
	}

	hitEvents := make([]*types.Event, 0)
	// remind if a biz set events has already been hit, if yes, then skip this event.
	reminder := make(map[string]struct{})

	for _, one := range es {
		if _, yes := reminder[one.Oid]; yes {
			// this biz set event is already hit, then we aggregate the event with the former ones to only one event.
			// this is useful to decrease biz set relation events.
			blog.Infof("biz set event: %s is aggregated, rid: %s", one.ID(), rid)
			continue
		}

		// get the hit event(nil means no event is hit) and if other events of the biz set needs to be ignored
		hitEvent, needIgnored := b.parseBizSetEvent(one, oidDetailMap, lastEventMap, rid)

		if hitEvent != nil {
			hitEvents = append(hitEvents, hitEvent)
		}

		if needIgnored {
			reminder[one.Oid] = struct{}{}
		}
	}

	// refresh all biz ids cache if the events contains match all biz set, the cache is used to generate detail later
	for _, e := range hitEvents {
		if e.OperationType != types.Delete && gjson.Get(string(e.DocBytes), "bk_scope.match_all").Bool() {
			err := b.refreshAllBizIDStr(rid)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	return hitEvents, nil
}

// parseBizSetEvent parse biz set event, returns the hit event and if same oid events needs to be skipped
func (b *bizSetRelation) parseBizSetEvent(one *types.Event, oidDetailMap map[string][]byte,
	lastEventMap map[string]*types.Event, rid string) (*types.Event, bool) {

	switch one.OperationType {
	case types.Insert:
		// if events are in the order of create + update + delete, all events of this biz set will be ignored
		lastEvent := lastEventMap[one.Oid]
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
		// reset delete event detail to the value before deletion from del archive, then add to hit events.
		doc, exist := oidDetailMap[one.Oid]
		if !exist {
			blog.Errorf("%s event delete detail[oid: %s] not exists, rid: %s", b.key.Collection(), one.Oid, rid)
			return nil, false
		}
		one.DocBytes = doc
		return one, true
	case types.Update, types.Replace:
		// if events are in the order of update + delete, it is changed to delete event, update event is ignored
		if lastEventMap[one.Oid].OperationType == types.Delete {
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

// delArchive delete event archived detail struct, with oid and detail
type delArchive struct {
	oid    string                 `bson:"oid"`
	Detail map[string]interface{} `bson:"detail"`
}

// rearrangeBizEvents TODO
// biz events rearrange policy:
// 1. Biz event is redirected to its related biz sets' events by traversing all biz sets and checking if the "bk_scope"
//    field matches the biz's attribute.
// 2. Create and delete event's related biz set is judged by whether its scope contains the biz.
// 3. Update event's related biz set is judged by whether its scope contains the updated fields of the event. Since
//    we can't get the previous value of the updated fields, we can't get the exact biz sets it was in before.
// 4. Aggregate multiple biz events with the same biz set to one event.
func (b *bizSetRelation) rearrangeBizEvents(es []*types.Event, rid string) ([]*types.Event, error) {

	// get delete event oids from delete events, and get deleted biz detail by oids to find matching biz sets.
	deletedOids := make([]string, 0)
	for _, one := range es {
		if one.OperationType == types.Delete {
			deletedOids = append(deletedOids, one.Oid)
		}
	}

	deletedDetailMap := make(map[string]map[string]interface{})
	if len(deletedOids) > 0 {
		filter := map[string]interface{}{
			"oid":  map[string]interface{}{common.BKDBIN: deletedOids},
			"coll": b.key.Collection(),
		}

		docs := make([]delArchive, 0)
		err := b.ccDB.Table(common.BKTableNameDelArchive).Find(filter).Fields("detail").All(context.Background(), &docs)
		if err != nil {
			b.metrics.CollectMongoError()
			blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v, rid: %s",
				b.key.Collection(), deletedOids, err, rid)
			return nil, err
		}

		for _, doc := range docs {
			deletedDetailMap[doc.oid] = doc.Detail
		}
	}

	// parse biz events to convert biz event parameter
	params := b.parseBizEvents(es, deletedDetailMap, rid)

	// convert biz event to related biz set event whose relation is affected by biz event
	bizSetEvents, err := b.convertBizEvent(params)
	if err != nil {
		return nil, err
	}

	return bizSetEvents, nil
}

// parseBizEvents parse biz events to insert/delete event details and updated fields, used to generate biz set events
func (b *bizSetRelation) parseBizEvents(es []*types.Event, deletedDetailMap map[string]map[string]interface{},
	rid string) convertBizEventParams {

	// generate convert biz event to biz set event parameter
	params := convertBizEventParams{
		es:                            es,
		insertedAndDeletedBiz:         make([]map[string]interface{}, 0),
		insertedAndDeletedBizIndexMap: make(map[int]int, 0),
		updatedFieldsIndexMap:         make(map[string]int, 0),
		needCareFieldsMap:             b.needCareBizFields.Get(),
		rid:                           rid,
	}

	// parse biz event into inserted/deleted details, and updated fields to match biz sets whose relation changed.
	for index, one := range es {
		switch one.OperationType {
		case types.Insert:
			// use document in insert event to find matching biz sets.
			params.insertedAndDeletedBiz = append(params.insertedAndDeletedBiz, *one.Document.(*map[string]interface{}))
			params.insertedAndDeletedBizIndexMap[len(params.insertedAndDeletedBiz)-1] = index

		case types.Delete:
			// use document in delete event detail map to find matching biz sets.
			detail, exists := deletedDetailMap[one.Oid]
			if !exists {
				blog.Errorf("%s event delete detail[oid: %s] not exists, rid: %s", b.key.Collection(), one.Oid, rid)
				continue
			}
			params.insertedAndDeletedBiz = append(params.insertedAndDeletedBiz, detail)
			params.insertedAndDeletedBizIndexMap[len(params.insertedAndDeletedBiz)-1] = index

		case types.Update, types.Replace:
			// replace event's change description is empty, treat as updating all fields.
			if len(one.ChangeDesc.UpdatedFields) == 0 && len(one.ChangeDesc.RemovedFields) == 0 {
				for field := range params.needCareFieldsMap {
					if _, exists := params.updatedFieldsIndexMap[field]; !exists {
						params.updatedFieldsIndexMap[field] = index
					}
				}
				continue
			}

			// get the updated/removed fields that need to be cared, ignores the rest.
			isIgnored := true
			if len(one.ChangeDesc.UpdatedFields) > 0 {
				// for biz set relation,archive biz is treated as delete event while recover is treated as insert event
				if one.ChangeDesc.UpdatedFields[common.BKDataStatusField] == string(common.DataStatusDisabled) ||
					one.ChangeDesc.UpdatedFields[common.BKDataStatusField] == string(common.DataStatusEnable) {

					detail := *one.Document.(*map[string]interface{})
					params.insertedAndDeletedBiz = append(params.insertedAndDeletedBiz, detail)
					params.insertedAndDeletedBizIndexMap[len(params.insertedAndDeletedBiz)-1] = index
					continue
				}

				for field := range one.ChangeDesc.UpdatedFields {
					if _, exists := params.needCareFieldsMap[field]; exists {
						if _, exists := params.updatedFieldsIndexMap[field]; !exists {
							params.updatedFieldsIndexMap[field] = index
						}
						isIgnored = false
					}
				}
			}

			if len(one.ChangeDesc.RemovedFields) > 0 {
				for _, field := range one.ChangeDesc.RemovedFields {
					if _, exists := params.needCareFieldsMap[field]; exists {
						if _, exists := params.updatedFieldsIndexMap[field]; !exists {
							params.updatedFieldsIndexMap[field] = index
						}
						isIgnored = false
					}
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
	insertedAndDeletedBiz []map[string]interface{}
	// insertedAndDeletedBizIndexMap mapping of insertedAndDeletedBiz index to es index, used to locate origin event
	insertedAndDeletedBizIndexMap map[int]int
	// updatedFieldsIndexMap mapping of updated fields to the first update event index, used to locate origin event
	updatedFieldsIndexMap map[string]int
	// needCareFieldsMap need care biz fields that can be used in biz set scope
	needCareFieldsMap map[string]string
	// rid request id
	rid string
}

// convertBizEvent convert biz event to related biz set event whose relation is affected by biz event
func (b *bizSetRelation) convertBizEvent(params convertBizEventParams) ([]*types.Event, error) {

	bizSetEvents := make([]*types.Event, 0)

	if len(params.insertedAndDeletedBiz) == 0 && len(params.updatedFieldsIndexMap) == 0 {
		return bizSetEvents, nil
	}

	// get all biz sets to check if their scope matches the biz events, because filtering them in db is too complicated
	bizSets, err := b.getAllBizSets(context.Background(), params.rid)
	if err != nil {
		return nil, err
	}

	// get biz events' related biz sets
	relatedBizSets, containsMatchAllBizSet := b.getRelatedBizSets(params, bizSets)

	// refresh all biz ids cache if the events affected match all biz set, the cache is used to generate detail later
	if containsMatchAllBizSet {
		if err := b.refreshAllBizIDStr(params.rid); err != nil {
			return nil, err
		}
	}

	// parse biz sets into biz set events
	for bizEventIndex, bizSetArr := range relatedBizSets {
		bizEvent := params.es[bizEventIndex]

		for _, bizSet := range bizSetArr {
			doc, err := json.Marshal(bizSet.BizSetInst)
			if err != nil {
				blog.Errorf("marshal biz set(%+v) failed, err: %v, rid: %s", bizSet.BizSetInst, err, params.rid)
				return nil, err
			}

			bizSetEvents = append(bizSetEvents, &types.Event{
				Oid:           bizSet.Oid.Hex(),
				Document:      bizSet.BizSetInst,
				DocBytes:      doc,
				OperationType: types.Update,
				Collection:    common.BKTableNameBaseBizSet,
				ClusterTime:   bizEvent.ClusterTime,
				Token:         bizEvent.Token,
			})
		}
	}

	return bizSetEvents, nil
}

// getRelatedBizSets get biz events index to related biz sets map, which is used to generate biz set events
func (b *bizSetRelation) getRelatedBizSets(params convertBizEventParams, bizSets []bizSetWithOid) (
	map[int][]bizSetWithOid, bool) {

	containsMatchAllBizSet := false
	relatedBizSets := make(map[int][]bizSetWithOid, 0)
	for _, bizSet := range bizSets {
		// for biz set that matches all biz, only insert and delete event will affect their relations
		if bizSet.Scope.MatchAll {
			if len(params.insertedAndDeletedBiz) > 0 {
				eventIndex := params.insertedAndDeletedBizIndexMap[0]
				relatedBizSets[eventIndex] = append(relatedBizSets[eventIndex], bizSet)
				containsMatchAllBizSet = true
			}
			continue
		}

		if bizSet.Scope.Filter == nil {
			blog.Errorf("biz set(%+v) scope filter is empty, skip, rid: %s", bizSet, params.rid)
			continue
		}

		var firstEventIndex int

		// update biz event matches all biz sets whose scope contains the updated fields, get all matching fields.
		matched := bizSet.Scope.Filter.MatchAny(func(r querybuilder.AtomRule) bool {
			if index, exists := params.updatedFieldsIndexMap[r.Field]; exists {
				if firstEventIndex == 0 || index < firstEventIndex {
					firstEventIndex = index
				}
				return true
			}
			return false
		})

		// check if biz set scope filter matches the inserted/removed biz
		for index, biz := range params.insertedAndDeletedBiz {
			// if the event index already exceeds the event index of matched update fields, stop checking
			eventIndex := params.insertedAndDeletedBizIndexMap[index]
			if firstEventIndex != 0 && eventIndex >= firstEventIndex {
				break
			}

			bizMatched := bizSet.Scope.Filter.Match(func(r querybuilder.AtomRule) bool {
				// ignores the biz set filter rule that do not contain need care fields
				propertyType, exists := params.needCareFieldsMap[r.Field]
				if !exists {
					blog.Errorf("biz set(%+v) filter rule contains ignored field, rid: %s", bizSet, params.rid)
					return false
				}

				// ignores the biz that do not contain the field in filter rule
				bizVal, exists := biz[r.Field]
				if !exists {
					blog.Infof("biz(%+v) do not contain rule field %s, rid: %s", biz, r.Field, params.rid)
					return false
				}

				switch r.Operator {
				case querybuilder.OperatorEqual:
					return matchEqualOper(r.Value, bizVal, propertyType, params.rid)
				case querybuilder.OperatorIn:
					return matchInOper(r.Value, bizVal, propertyType, params.rid)
				default:
					blog.Errorf("biz set(%+v) filter rule contains invalid operator, rid: %s", bizSet, params.rid)
					return false
				}
			})

			if bizMatched {
				firstEventIndex = eventIndex
				matched = bizMatched
				break
			}
		}

		if matched {
			relatedBizSets[firstEventIndex] = append(relatedBizSets[firstEventIndex], bizSet)
		}
	}

	return relatedBizSets, containsMatchAllBizSet
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

func (b *bizSetRelation) getAllBizSets(ctx context.Context, rid string) ([]bizSetWithOid, error) {
	const step = 500

	bizSets := make([]bizSetWithOid, 0)

	cond := map[string]interface{}{}
	findOpts := dbtypes.NewFindOpts().SetWithObjectID(true)

	for {
		oneStep := make([]bizSetWithOid, 0)
		err := b.ccDB.Table(common.BKTableNameBaseBizSet).Find(cond, findOpts).Fields(common.BKBizSetIDField,
			common.BKBizSetScopeField).Limit(step).Sort(common.BKBizSetIDField).All(ctx, &oneStep)
		if err != nil {
			blog.Errorf("get biz set failed, err: %v, rid: %s", err, rid)
			return nil, err
		}

		bizSets = append(bizSets, oneStep...)

		if len(oneStep) < step {
			break
		}

		cond = map[string]interface{}{
			common.BKBizSetIDField: map[string]interface{}{common.BKDBGT: oneStep[len(oneStep)-1].BizSetID},
		}
	}

	return bizSets, nil
}

// getDeleteEventDetails get delete events' oid to related detail map from db
func (b *bizSetRelation) getDeleteEventDetails(oids []string, rid string) (map[string][]byte, error) {
	oidDetailMap := make(map[string][]byte)

	if len(oids) == 0 {
		return oidDetailMap, nil
	}

	filter := map[string]interface{}{
		"oid":  map[string]interface{}{common.BKDBIN: oids},
		"coll": b.key.Collection(),
	}

	docs := make([]map[string]interface{}, 0)
	err := b.ccDB.Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
	if err != nil {
		b.metrics.CollectMongoError()
		blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v, rid: %s",
			b.key.Collection(), oids, err, rid)
		return nil, err
	}

	for _, doc := range docs {
		oid := util.GetStrByInterface(doc["oid"])
		byt, err := json.Marshal(doc["detail"])
		if err != nil {
			blog.Errorf("get archive deleted doc for collection %s, but marshal detail to bytes failed, oid: %s, "+
				"err: %v, rid: %s", b.key.Collection(), oid, err, rid)
			return nil, err
		}
		oidDetailMap[oid] = byt
	}

	return oidDetailMap, nil
}
