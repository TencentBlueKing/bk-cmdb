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

package identifier

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"

	"github.com/tidwall/gjson"
)

var needCareHostFields = []string{
	common.BKHostInnerIPField,
	common.BKOSTypeField,
	common.BKCloudIDField,
	common.BKHostIDField,
	common.BKAgentIDField,
	common.BKAddressingField,
}

// rearrangeHostEvents TODO
// host events arrange policy:
// 1. do not care delete events, cause if the host is already dropped, then it's identifier event is meaningless.
// 2. if event's ChangeDesc contains fields as follows, then we care about it, otherwise we can drop this events
//   - bk_host_id
//   - bk_os_type
//   - bk_cloud_id
//   - bk_host_innerip
//     if ChangeDesc is empty, then we assume this events is need to be care.
//     3. aggregate multiple same host's event to one event, so that we can decrease the amount of host identity.
//     because we only care about which host is changed, one event is enough for us.
func (f *hostIdentity) rearrangeHostEvents(es []*types.Event, rid string) []*types.Event {

	hitEvents := make([]*types.Event, 0)
	// remind if a host events has already been hit, if yes, then skip this event.
	reminder := make(map[string]struct{})
loop:
	for idx := range es {
		one := es[idx]
		if one.OperationType == types.Delete {
			blog.Infof("host identify event, received delete host event, skip, oid: %s, rid: %s", one.ID(), rid)
			// skip delete event
			continue
		}

		if _, yes := reminder[genUniqueKey(one)]; yes {
			// this host event has already hit, then we aggregate these events with former one to only one
			// this is useful to decrease host identify events.
			blog.Infof("host identify event, host event: %s is aggregated, rid: %s", one.ID(), rid)
			continue
		}

		// add/replace event's change description is empty.
		if len(one.ChangeDesc.UpdatedFields) == 0 && len(one.ChangeDesc.RemovedFields) == 0 {
			// we do not know what's info is changed, so we add this event directly.
			hitEvents = append(hitEvents, one)
			reminder[genUniqueKey(one)] = struct{}{}
			continue
		}

		// check updated fields
		if len(one.ChangeDesc.UpdatedFields) != 0 {
			for _, care := range needCareHostFields {
				if _, yes := one.ChangeDesc.UpdatedFields[care]; yes {
					hitEvents = append(hitEvents, one)
					reminder[genUniqueKey(one)] = struct{}{}
					continue loop
				}
			}
		}

		// check removed fields
		if len(one.ChangeDesc.RemovedFields) != 0 {
			check := make(map[string]struct{})
			for _, key := range one.ChangeDesc.RemovedFields {
				check[key] = struct{}{}
			}

			for _, care := range needCareHostFields {
				if _, yes := check[care]; yes {
					// one of the cared fields is removed, we do need to care.
					hitEvents = append(hitEvents, one)
					reminder[genUniqueKey(one)] = struct{}{}
					continue loop
				}
			}
		}

		// this is event is need to be ignored when comes to here.
		blog.Infof("host identify event, host changed detail do not care, skip, oid: %s, detail: %+v, rid: %s",
			one.ID(), one.ChangeDesc, rid)
	}

	return hitEvents
}

var hostIDJson = `{"bk_host_id":%d}`

// rearrangeHostRelationEvents rearrange host relation events
// host relation events arrange policy:
// 1. redirect relation event to host change event.
// 2. care about all kinds of event types.
// 3. do not care the event's order, cause we all convert to host events type.
func (f *hostIdentity) rearrangeHostRelationEvents(es []*types.Event, rid string) ([]*types.Event, error) {
	hitEvents := make([]*types.Event, 0)
	// remind if related host's events has already been hit, if yes, then skip this event.
	reminder := make(map[int64]struct{})
	for idx := range es {
		one := es[idx]

		hostID := gjson.GetBytes(one.DocBytes, common.BKHostIDField).Int()
		if hostID <= 0 {
			blog.Errorf("host identify event, get host id from relation: %s failed, skip, rid: %s", one.DocBytes, rid)
			continue
		}

		if _, exist := reminder[hostID]; exist {
			// this host has already been hit, skip now.
			blog.Infof("host identify event, relation host id: %d/%s is aggregated, rid: %s", hostID, one.ID(), rid)
			continue
		}

		// convert to host id event detail
		one.DocBytes = []byte(fmt.Sprintf(hostIDJson, hostID))
		one.Document = nil
		hitEvents = append(hitEvents, one)
		// update reminder
		reminder[hostID] = struct{}{}
	}

	return hitEvents, nil
}

// rearrangeProcessEvents process events arrange policy:
// 1. redirect process event to host change event with process host relation info.
// 2. care about all kinds of event types.
// 3. do not care the event's order, cause we all convert to host events type.
func (f *hostIdentity) rearrangeProcessEvents(es []*types.Event, rid string) ([]*types.Event, error) {
	if len(es) == 0 {
		return es, nil
	}

	processIDsMap := make(map[string][]int64)
	reminder := make(map[string]struct{})
	for idx := range es {
		one := es[idx]

		if _, exist := reminder[genUniqueKey(one)]; exist {
			// skip event's with the same oid, which means it's the same process event.
			// cause we convert a process id to host id finally.
			blog.Infof("host identify event, process: %s is aggregated, rid: %s", one.ID(), rid)
			continue
		}

		// process delete and insert event are handled by related process relation event.
		if one.OperationType == types.Delete || one.OperationType == types.Insert {
			reminder[genUniqueKey(one)] = struct{}{}
			continue
		}

		processID := gjson.GetBytes(one.DocBytes, common.BKProcessIDField).Int()
		if processID <= 0 {
			blog.Errorf("host identify event, get process id from process: %s failed, skip, rid: %s", one.DocBytes, rid)
			continue
		}

		processIDsMap[one.TenantID] = append(processIDsMap[one.TenantID], processID)
		reminder[genUniqueKey(one)] = struct{}{}
	}

	// got 0 valid event
	if len(processIDsMap) == 0 {
		return es[:0], nil
	}

	// now we need to convert these process ids to host ids.
	hostListMap := make(map[string][]int64)
	for tenant, processIDs := range processIDsMap {
		hostList, err := f.convertProcessToHost(tenant, processIDs, rid)
		if err != nil {
			return nil, err
		}

		hostListMap[tenant] = hostList
	}

	for tenantID, hostList := range hostListMap {
		// now we get all the host's ids list
		// it should be much more less than the process's count
		hostListMap[tenantID] = util.IntArrayUnique(hostList)
	}

	events := make([]*types.Event, 0)
	for _, e := range es {
		hostList := hostListMap[e.TenantID]
		if len(hostList) == 0 {
			continue
		}
		// reset the event's document info to host id field.
		e.DocBytes = []byte(fmt.Sprintf(hostIDJson, hostList[0]))
		hostListMap[e.TenantID] = hostList[1:]
		e.Document = nil
		events = append(events, e)
	}
	return es, nil
}

type processRelation struct {
	ProcessID int64 `bson:"bk_process_id"`
	HostID    int64 `bson:"bk_host_id"`
}

// convertProcessToHost convert process ids to host ids.
func (f *hostIdentity) convertProcessToHost(tenantID string, pIDs []int64, rid string) ([]int64, error) {
	if len(pIDs) == 0 {
		return make([]int64, 0), nil
	}

	filter := mapstr.MapStr{
		common.BKProcessIDField: mapstr.MapStr{common.BKDBIN: pIDs},
	}

	relations := make([]*processRelation, 0)
	err := mongodb.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(common.BKTableNameProcessInstanceRelation).
		Find(filter).Fields(common.BKHostIDField, common.BKProcessIDField).All(context.Background(), &relations)
	if err != nil {
		blog.Errorf("host identify event, get process instance relation failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	hostIDMap := make(map[int64]struct{})
	for idx := range relations {
		hostIDMap[relations[idx].HostID] = struct{}{}
	}

	hostList := make([]int64, 0)
	for id := range hostIDMap {
		hostList = append(hostList, id)
	}

	return hostList, nil
}

func genUniqueKey(e *types.Event) string {
	return e.Collection + "-" + e.Oid
}
