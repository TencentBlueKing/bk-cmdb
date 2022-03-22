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
	"bytes"
	"context"
	"fmt"
	"strconv"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	mixevent "configcenter/src/source_controller/cacheservice/event/mix-event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream"
	"configcenter/src/storage/stream/types"
)

// newBizSetRelation init and run biz set relation event watch with sub event key
func newBizSetRelation(ctx context.Context, opts mixevent.MixEventFlowOptions) error {
	relation := bizSetRelation{
		watch:             opts.Watch,
		watchDB:           opts.WatchDB,
		ccDB:              opts.CcDB,
		mixKey:            opts.MixKey,
		key:               opts.Key,
		needCareBizFields: new(needCareBizFields),
		allBizIDStr:       new(allBizIDStr),
		metrics:           event.InitialMetrics(opts.Key.Collection(), "biz_set_relation"),
	}

	// get need care biz fields for biz event conversion, then sync it in goroutine
	fields, err := relation.getNeedCareBizFields(context.Background())
	if err != nil {
		blog.Errorf("run biz set relation watch, but get need care biz fields failed, err: %v", err)
		return err
	}
	relation.needCareBizFields.Set(fields)

	go relation.syncNeedCareBizFields()

	flow, err := mixevent.NewMixEventFlow(opts, relation.rearrangeEvents, relation.parseEvent)
	if err != nil {
		return err
	}

	return flow.RunFlow(ctx)
}

// bizSetRelation biz set relation event watch logic struct
type bizSetRelation struct {
	watch             stream.LoopInterface
	watchDB           *local.Mongo
	mixKey            event.Key
	key               event.Key
	ccDB              dal.DB
	needCareBizFields *needCareBizFields
	allBizIDStr       *allBizIDStr

	metrics *event.EventMetrics
}

// rearrangeEvents rearrange biz set and biz events into biz set events whose biz set changes
func (b *bizSetRelation) rearrangeEvents(rid string, es []*types.Event) ([]*types.Event, error) {
	switch b.key.Collection() {
	case event.BizSetKey.Collection():
		return b.rearrangeBizSetEvents(es, rid)
	case event.BizKey.Collection():
		return b.rearrangeBizEvents(es, rid)
	default:
		blog.Errorf("received unsupported biz set relation event, skip, es: %+v, rid: %s", es, rid)
		return es[:0], nil
	}
}

// parseEvent parse event into chain node and detail, detail is biz set id and its related biz ids
func (b *bizSetRelation) parseEvent(e *types.Event, id uint64, rid string) (*watch.ChainNode, []byte, bool, error) {
	switch e.OperationType {
	case types.Insert, types.Update, types.Replace, types.Delete:
	case types.Invalidate:
		blog.Errorf("biz set relation event, received invalid event operation type, doc: %s, rid: %s", e.DocBytes, rid)
		return nil, nil, false, nil
	default:
		blog.Errorf("biz set relation event, received unsupported event operation type: %s, doc: %s, rid: %s",
			e.OperationType, e.DocBytes, rid)
		return nil, nil, false, nil
	}

	name := b.key.Name(e.DocBytes)
	cursor, err := genBizSetRelationCursor(b.key.Collection(), e, rid)
	if err != nil {
		blog.Errorf("get %s event cursor failed, name: %s, err: %v, oid: %s, rid: %s", b.key.Collection(), name,
			err, e.ID(), rid)
		return nil, nil, false, err
	}

	chainNode := &watch.ChainNode{
		ID:          id,
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		// redirect all the event type to update.
		EventType: watch.ConvertOperateType(types.Update),
		Token:     e.Token.Data,
		Cursor:    cursor,
	}

	if instanceID := b.mixKey.InstanceID(e.DocBytes); instanceID > 0 {
		chainNode.InstanceID = instanceID
	}

	relationDetail, err := b.getBizSetRelationDetail(e, rid)
	if err != nil {
		return nil, nil, true, err
	}

	detail := types.EventDetail{
		Detail: types.JsonString(relationDetail),
	}
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		blog.Errorf("run %s flow, %s, marshal detail failed, name: %s, detail: %+v, err: %v, oid: %s, rid: %s",
			b.mixKey.Collection(), b.key.Collection(), name, detail, err, e.ID(), rid)
		return nil, nil, false, err
	}

	return chainNode, detailBytes, false, nil
}

// allBizIDStr struct to cache all biz ids in string form to generate detail, refreshed when rearranging events
type allBizIDStr struct {
	data string
	lock sync.RWMutex
}

// Get get all biz ids in string form
func (a *allBizIDStr) Get() string {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.data
}

// Set set all biz ids in string form
func (a *allBizIDStr) Set(data string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.data = data
}

// refreshAllBizIDStr refresh all biz ids in string form
func (b *bizSetRelation) refreshAllBizIDStr(rid string) error {
	// do not include resource pool and disabled biz in biz set
	allBizIDCond := map[string]interface{}{
		common.BKDefaultField:    mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag},
		common.BKDataStatusField: map[string]interface{}{common.BKDBNE: common.DataStatusDisabled},
	}

	allBizID, err := b.getBizIDArrStrByCond(allBizIDCond, rid)
	if err != nil {
		return err
	}

	b.allBizIDStr.Set(allBizID)
	return nil
}

// getBizSetRelationDetail get biz set relation detail by biz set event
func (b *bizSetRelation) getBizSetRelationDetail(e *types.Event, rid string) (string, error) {

	// get biz set relation detail by the scope of biz set event doc bytes
	bizSet := new(metadata.BizSetInst)
	if err := json.Unmarshal(e.DocBytes, bizSet); err != nil {
		blog.Errorf("unmarshal biz set(%s) failed, err: %v", string(e.DocBytes), err)
		return "", err
	}

	// deleted biz set is treated as updated to having no relations
	if e.OperationType == types.Delete {
		return event.GenBizSetRelationDetail(bizSet.BizSetID, ""), nil
	}

	allBizID := b.allBizIDStr.Get()

	// biz set that matches all biz uses the same all biz ids from cache
	if bizSet.Scope.MatchAll {
		if len(allBizID) == 0 {
			// do not include resource pool and disabled biz in biz set
			allBizIDCond := map[string]interface{}{
				common.BKDefaultField:    mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag},
				common.BKDataStatusField: map[string]interface{}{common.BKDBNE: common.DataStatusDisabled},
			}

			var err error
			allBizID, err = b.getBizIDArrStrByCond(allBizIDCond, rid)
			if err != nil {
				return "", err
			}
			b.allBizIDStr.Set(allBizID)
		}
		return event.GenBizSetRelationDetail(bizSet.BizSetID, allBizID), nil
	}

	// biz set scope with an empty filter is treated as having no relations
	if bizSet.Scope.Filter == nil {
		return event.GenBizSetRelationDetail(bizSet.BizSetID, ""), nil
	}

	// parse biz condition from biz set scope filter, get biz ids using it to gen relation detail
	bizSetBizCond, _, rawErr := bizSet.Scope.Filter.ToMgo()
	if rawErr != nil {
		blog.Errorf("parse biz set scope(%#v) failed, err: %v, rid: %s", bizSet.Scope, rawErr, rid)
		return "", rawErr
	}

	// do not include resource pool and disabled biz in biz set
	bizSetBizCond[common.BKDefaultField] = mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag}
	bizSetBizCond[common.BKDataStatusField] = map[string]interface{}{common.BKDBNE: common.DataStatusDisabled}

	bizIDStr, err := b.getBizIDArrStrByCond(bizSetBizCond, rid)
	if err != nil {
		return "", err
	}

	return event.GenBizSetRelationDetail(bizSet.BizSetID, bizIDStr), nil
}

func (b *bizSetRelation) getBizIDArrStrByCond(cond map[string]interface{}, rid string) (string, error) {
	const step = 500

	bizIDJson := bytes.Buffer{}

	for start := uint64(0); ; start += step {
		oneStep := make([]metadata.BizInst, 0)

		err := b.ccDB.Table(common.BKTableNameBaseApp).Find(cond).Fields(common.BKAppIDField).Start(start).
			Limit(step).Sort(common.BKAppIDField).All(context.Background(), &oneStep)
		if err != nil {
			blog.Errorf("get biz by cond(%+v) failed, err: %v, rid: %s", cond, err, rid)
			return "", err
		}

		for _, biz := range oneStep {
			bizIDJson.WriteString(strconv.FormatInt(biz.BizID, 10))
			bizIDJson.WriteByte(',')
		}

		if len(oneStep) < step {
			break
		}
	}

	// returns biz ids string form joined by comma, trim the extra trilling comma
	if bizIDJson.Len() == 0 {
		return "", nil
	}
	return bizIDJson.String()[:bizIDJson.Len()-1], nil
}

func genBizSetRelationCursor(coll string, e *types.Event, rid string) (string, error) {
	curType := watch.UnknownType
	switch coll {
	case common.BKTableNameBaseBizSet:
		curType = watch.BizSet
	case common.BKTableNameBaseApp:
		curType = watch.Biz
	default:
		blog.Errorf("unsupported biz set relation cursor type collection: %s, event: %+v, oid: %s", coll, e, rid)
		return "", fmt.Errorf("unsupported biz set relation cursor type collection: %s", coll)
	}

	cursor := &watch.Cursor{
		Type:        curType,
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		Oper:        e.OperationType,
	}

	cursorEncode, err := cursor.Encode()
	if err != nil {
		blog.Errorf("encode biz set relation cursor failed, cursor: %+v, err: %v, rid: %s", cursor, err, rid)
		return "", err
	}

	return cursorEncode, nil
}
