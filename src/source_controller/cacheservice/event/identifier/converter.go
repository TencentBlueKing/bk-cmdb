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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/stream/types"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var needCareHostFields = []string{
	common.BKHostInnerIPField,
	common.BKOSTypeField,
	common.BKCloudIDField,
	common.BKHostIDField,
	common.BkSupplierAccount,
}

// rearrangeHostEvents TODO
// host events arrange policy:
// 1. do not care delete events, cause if the host is already dropped, then it's identifier event is meaningless.
// 2. if event's ChangeDesc contains fields as follows, then we care about it, otherwise we can drop this events
//    - bk_host_id
//    - bk_os_type
//    - bk_cloud_id
//    - bk_host_innerip
//    if ChangeDesc is empty, then we assume this events is need to be care.
// 3. aggregate multiple same host's event to one event, so that we can decrease the amount of host identity. because
//    we only care about which host is changed, one event is enough for us.
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

		if _, yes := reminder[one.Oid]; yes {
			// this host event has already hit, then we aggregate these events with former one to only one
			// this is useful to decrease host identify events.
			blog.Infof("host identify event, host event: %s is aggregated, rid: %s", one.ID(), rid)
			continue
		}

		// add/replace event's change description is empty.
		if len(one.ChangeDesc.UpdatedFields) == 0 && len(one.ChangeDesc.RemovedFields) == 0 {
			// we do not know what's info is changed, so we add this event directly.
			hitEvents = append(hitEvents, one)
			reminder[one.Oid] = struct{}{}
			continue
		}

		// check updated fields
		if len(one.ChangeDesc.UpdatedFields) != 0 {
			for _, care := range needCareHostFields {
				if _, yes := one.ChangeDesc.UpdatedFields[care]; yes {
					hitEvents = append(hitEvents, one)
					reminder[one.Oid] = struct{}{}
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
					reminder[one.Oid] = struct{}{}
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

// rearrangeHostRelationEvents TODO
// host relation events arrange policy:
// 1. redirect relation event to host change event.
// 2. care about all kinds of event types.
// 3. do not care the event's order, cause we all convert to host events type.
func (f *hostIdentity) rearrangeHostRelationEvents(es []*types.Event, rid string) ([]*types.Event, error) {
	deleteEventsMap := make(map[string]*types.Event, 0)
	deleteOids := make([]string, 0)
	hitEvents := make([]*types.Event, 0)
	// remind if related host's events has already been hit, if yes, then skip this event.
	reminder := make(map[int64]struct{})
	for idx := range es {
		one := es[idx]
		if one.OperationType == types.Delete {
			deleteEventsMap[one.Oid] = one
			deleteOids = append(deleteOids, one.Oid)
			continue
		}

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

	if len(deleteEventsMap) == 0 {
		// no delete type events, then return directly
		return hitEvents, nil
	}

	filter := map[string]interface{}{
		"oid":  map[string]interface{}{common.BKDBIN: deleteOids},
		"coll": common.BKTableNameModuleHostConfig,
	}

	docs := make([]bsonx.Doc, 0)
	err := f.ccDB.Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
	if err != nil {
		f.metrics.CollectMongoError()
		blog.Errorf("host identify event, get archive host relation from mongodb failed, oid: %+v, err: %v, rid: %v",
			deleteOids, err, rid)
		return nil, err
	}

	for _, doc := range docs {
		hostID := doc.Lookup("detail", common.BKHostIDField).Int64()
		if hostID <= 0 {
			blog.Errorf("host identify event, get host id from relation: %s failed, skip, rid: %s",
				doc.Lookup("detail").String(), rid)
			continue
		}

		if _, exist := reminder[hostID]; exist {
			// this host has already been hit, skip now.
			blog.Infof("host identify event, relation deleted host id: %d is aggregated, oid: %s, rid: %s",
				hostID, doc.Lookup("oid").String(), rid)
			continue
		}
		reminder[hostID] = struct{}{}

		event, exist := deleteEventsMap[doc.Lookup("oid").String()]
		if !exist {
			blog.Errorf("host identify event, get archived event's instance with oid :%s failed, skip, rid: %s",
				doc.Lookup("oid").String(), rid)
			continue
		}
		event.DocBytes = []byte(fmt.Sprintf(hostIDJson, hostID))
		event.Document = nil
		hitEvents = append(hitEvents, event)
	}

	return hitEvents, nil
}

// rearrangeProcessEvents TODO
// process events arrange policy:
// 1. redirect process event to host change event with process host relation info.
// 2. care about all kinds of event types.
// 3. do not care the event's order, cause we all convert to host events type.
func (f *hostIdentity) rearrangeProcessEvents(es []*types.Event, rid string) ([]*types.Event, error) {
	if len(es) == 0 {
		return es, nil
	}

	processIDs := make([]int64, 0)
	deleteOids := make([]string, 0)
	idToOid := make(map[int64]string)
	oidToEvent := make(map[string]*types.Event)
	reminder := make(map[string]struct{})
	for idx := range es {
		one := es[idx]

		if _, exist := reminder[one.Oid]; exist {
			// skip event's with the same oid, which means it's the same process event.
			// cause we convert a process id to host id finally.
			blog.Infof("host identify event, process: %s is aggregated, rid: %s", one.ID(), rid)
			continue
		}

		oidToEvent[one.Oid] = one
		if one.OperationType == types.Delete {
			deleteOids = append(deleteOids, one.Oid)
			reminder[one.Oid] = struct{}{}
			continue
		}

		processID := gjson.GetBytes(one.DocBytes, common.BKProcessIDField).Int()
		if processID <= 0 {
			blog.Errorf("host identify event, get process id from process: %s failed, skip, rid: %s", one.DocBytes, rid)
			continue
		}

		processIDs = append(processIDs, processID)
		idToOid[processID] = one.Oid
		reminder[one.Oid] = struct{}{}
	}

	// got 0 valid event
	if len(processIDs) == 0 && len(deleteOids) == 0 {
		return es[:0], nil
	}

	// now we need to convert these process ids and delete oids to host ids.
	// convert process ids to host ids.
	notHitProcess, hostList, err := f.convertProcessToHost(processIDs, rid)
	if err != nil {
		return nil, err
	}

	// get these process's host from cc_DelArchive
	if len(notHitProcess) != 0 {
		start := int64(es[0].ClusterTime.Sec)
		hostIDs, err := f.getHostWithProcessRelationFromDelArchive(start, notHitProcess, rid)
		if err != nil {
			return nil, err
		}
		hostList = append(hostList, hostIDs...)
	}

	if len(deleteOids) != 0 {
		start := int64(es[0].ClusterTime.Sec)
		hostIDs, err := f.getDeletedProcessHosts(start, deleteOids, rid)
		if err != nil {
			return nil, err
		}
		hostList = append(hostList, hostIDs...)
	}

	// now we get all the host's ids list
	// it should be much more less than the process's count
	hostList = util.IntArrayUnique(hostList)

	cnt := len(hostList)
	if cnt > len(es) {
		// host count is always less or equal than the count of events.
		// when this happens, somethings must be wrong.
		blog.ErrorJSON("got more host count than it's process, use host count instead, es: %s, host: %s, rid: %s",
			es, hostList, rid)
		// continue handle this, but redirect count to event's count
		cnt = len(es)
	}

	// reset the event's document info to host id field.
	for i := 0; i < cnt; i++ {
		es[i].DocBytes = []byte(fmt.Sprintf(hostIDJson, hostList[i]))
		es[i].Document = nil
	}

	// remove the unused events
	es = es[:cnt]
	return es, nil
}

type processRelation struct {
	ProcessID int64 `bson:"bk_process_id"`
	HostID    int64 `bson:"bk_host_id"`
}

// convertProcessToHost TODO
// convert process ids to host ids.
// we may can not find process's relations info, cause it may already been deleted. so we need
// to find it in cc_DelArchive collection.
func (f *hostIdentity) convertProcessToHost(pIDs []int64, rid string) ([]int64, []int64, error) {
	if len(pIDs) == 0 {
		return make([]int64, 0), make([]int64, 0), nil
	}

	filter := mapstr.MapStr{
		common.BKProcessIDField: mapstr.MapStr{common.BKDBIN: pIDs},
	}

	relations := make([]*processRelation, 0)
	err := f.ccDB.Table(common.BKTableNameProcessInstanceRelation).Find(filter).Fields(common.BKHostIDField,
		common.BKProcessIDField).All(context.Background(), &relations)
	if err != nil {
		blog.Errorf("host identify event, get process instance relation failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	idMap := make(map[int64]struct{})
	hostIDMap := make(map[int64]struct{})
	for idx := range relations {
		idMap[relations[idx].ProcessID] = struct{}{}
		hostIDMap[relations[idx].HostID] = struct{}{}
	}

	notHitProcess := make([]int64, 0)
	for _, id := range pIDs {
		if _, exist := idMap[id]; !exist {
			// this process's relations has already been deleted, so we can not find it.
			// it will be try to search in cc_DelArchive later
			notHitProcess = append(notHitProcess, id)
		}
	}

	hostList := make([]int64, 0)
	for id := range hostIDMap {
		hostList = append(hostList, id)
	}

	return notHitProcess, hostList, nil

}

// getHostWithProcessRelationFromDelArchive TODO
// get host ids from cc_DelArchive with process's ids
// a process has only one relation, so we can use process ids find it's unique relations.
func (f *hostIdentity) getHostWithProcessRelationFromDelArchive(startUnix int64, pIDs []int64, rid string) (
	[]int64, error) {
	filter := mapstr.MapStr{
		"coll": common.BKTableNameProcessInstanceRelation,
		// this archive doc's created time must be greater than start unix time.
		"_id": mapstr.MapStr{
			common.BKDBGTE: primitive.NewObjectIDFromTimestamp(time.Unix(startUnix-60, 0)),
		},
		"detail.bk_process_id": mapstr.MapStr{common.BKDBIN: pIDs},
	}

	relations := make([]map[string]*processRelation, 0)
	err := f.ccDB.Table(common.BKTableNameDelArchive).Find(filter).Fields("detail").
		All(context.Background(), &relations)
	if err != nil {
		f.metrics.CollectMongoError()
		blog.Errorf("host identify event, get archive deleted instance process relations failed, "+
			"process ids: %v, err: %v, rid: %v", f.key.Collection(), pIDs, err, rid)
		return nil, err
	}

	if len(pIDs) != len(relations) {
		blog.ErrorJSON("host identify event, can not find all process ids relations, ids: %s, relations: %s, rid: %s",
			pIDs, relations)
	}

	hostIDs := make([]int64, 0)
	for _, doc := range relations {
		relation := doc["detail"]
		hostIDs = append(hostIDs, relation.HostID)
	}
	return hostIDs, nil
}

func (f *hostIdentity) getDeletedProcessHosts(startUnix int64, oids []string, rid string) ([]int64, error) {
	filter := map[string]interface{}{
		"oid":  map[string]interface{}{common.BKDBIN: oids},
		"coll": common.BKTableNameBaseProcess,
	}

	docs := make([]bsonx.Doc, 0)
	err := f.ccDB.Table(common.BKTableNameDelArchive).Find(filter).Fields("detail").All(context.Background(), &docs)
	if err != nil {
		f.metrics.CollectMongoError()
		blog.Errorf("host identify event, get archive deleted process instances, oids: %+v, err: %v, rid: %v",
			oids, err, rid)
		return nil, err
	}

	pList := make([]int64, 0)
	for _, doc := range docs {
		pID := doc.Lookup("detail", common.BKProcessIDField).Int64()
		if pID <= 0 {
			blog.Errorf("host identify event, get process id from instance: %s failed, skip, rid: %s",
				doc.Lookup("detail").String(), rid)
			continue
		}
		pList = append(pList, pID)
	}

	if len(pList) == 0 {
		blog.Warnf("got 0 valid process from archived collection with oids: %v, rid: %s", oids, rid)
		return pList, nil
	}

	// then get hosts list with these process ids.
	return f.getHostWithProcessRelationFromDelArchive(startUnix, pList, rid)
}
