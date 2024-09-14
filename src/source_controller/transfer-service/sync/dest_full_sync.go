/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package sync

import (
	"encoding/json"
	"time"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/lock"
	"configcenter/src/source_controller/transfer-service/sync/util"
	"configcenter/src/storage/driver/redis"
)

// loopPullFullSyncData loop pull full sync data
func (s *Syncer) loopPullFullSyncData() {
	ack := false

	for {
		if !s.isMaster.IsMaster() {
			blog.V(4).Infof("loop pull full sync data, but not master, skip")
			time.Sleep(5 * time.Minute)
			ack = false
			continue
		}

		locker := lock.NewLocker(redis.Client())
		locked, err := locker.Lock(types.FullSyncLockKey, time.Hour)
		if err != nil || !locked {
			blog.Errorf("do not get %s lock, err: %v, locked: %v", types.FullSyncLockKey, err, locked)
			time.Sleep(5 * time.Minute)
			continue
		}

		// get object ids for object instance resource sync
		var objIDs, quotedObjIDs []string
		util.RetryWrapper(3, func() (bool, error) {
			objIDs, quotedObjIDs, err = s.metadata.GetCommonObjIDs()
			if err != nil {
				blog.Errorf("get object ids failed, err: %v", err)
				return true, err
			}
			return false, nil
		})

		for _, resType := range types.ListAllResType() {
			syncer := s.resSyncerMap[resType]

			switch resType {
			case types.ObjectInstance:
				for _, objID := range objIDs {
					syncer.pullFullSyncData(objID, ack)
				}
			case types.InstAsst:
				for _, objID := range append(objIDs, common.BKInnerObjIDHost) {
					syncer.pullFullSyncData(objID, ack)
				}
			case types.QuotedInstance:
				for _, objID := range quotedObjIDs {
					syncer.pullFullSyncData(objID, ack)
				}
			default:
				syncer.pullFullSyncData("", ack)
			}
		}
		ack = true

		locker.Unlock()

		time.Sleep(5 * time.Minute)
	}
}

// pullFullSyncData pull full sync data for one resource
func (s *resSyncer) pullFullSyncData(subRes string, ack bool) {
	kit := util.NewKit()
	startTime := time.Now()
	blog.Infof("start pull %s-%s full sync data, start time: %s, rid: %s", s.lgc.ResType(), subRes, startTime, kit.Rid)

	hasMore := true
	var err error

	for hasMore && err == nil {
		util.RetryWrapper(3, func() (bool, error) {
			hasMore, err = s.doOnePullFullSyncDataStep(kit, subRes, ack)
			if err != nil {
				blog.Errorf("try %s-%s full sync step failed, err: %v, rid: %s", s.lgc.ResType(), subRes, err, kit.Rid)
				return true, err
			}
			return false, nil
		})
		ack = true
	}

	blog.Infof("pull %s-%s full sync data successfully, start time: %s, cost: %s, rid: %s", s.lgc.ResType(), subRes,
		time.Now(), time.Since(startTime), kit.Rid)
}

// doOnePullFullSyncDataStep do one pull full sync data step
func (s *resSyncer) doOnePullFullSyncDataStep(kit *util.Kit, subRes string, ack bool) (bool, error) {
	// pull full sync data from transfer medium
	pullOpt := &types.PullSyncDataOpt{
		ResType:     s.lgc.ResType(),
		SubRes:      subRes,
		IsIncrement: false,
		Ack:         ack,
	}
	syncInfo, err := s.transMedium.PullSyncData(kit.Ctx, kit.Header, pullOpt)
	if err != nil {
		blog.Errorf("pull %s-%s full sync data failed, err: %v, rid: %s", s.lgc.ResType(), subRes, err, kit.Rid)
		return false, err
	}

	if len(syncInfo.Info) == 0 {
		return syncInfo.Total != 0, nil
	}

	// parse full sync data, skip the invalid data
	syncData := new(types.FullSyncTransData)
	rawDataArr := make([]json.RawMessage, 0)
	syncData.Data = &rawDataArr
	err = json.Unmarshal(syncInfo.Info, syncData)
	if err != nil {
		blog.Errorf("unmarshal %s-%s full sync data(%s) failed, err: %v, skip these data, rid: %s", s.lgc.ResType(),
			subRes, syncInfo.Info, err, kit.Rid)
		return syncInfo.Total != 0, nil
	}

	dataArr, err := s.lgc.ParseDataArr(syncData.Name, subRes, rawDataArr, kit.Rid)
	if err != nil {
		blog.Errorf("parse %s-%s full sync data(%+v) failed, err: %v, rid: %s", s.lgc.ResType(), subRes, rawDataArr,
			err, kit.Rid)
		return syncInfo.Total != 0, nil
	}
	syncData.Data = dataArr

	// loop handle full sync data
	isAll := false
	for !isAll {
		var nextStart map[string]int64
		var remainingData any

		util.RetryWrapper(3, func() (bool, error) {
			isAll, nextStart, remainingData, err = s.handleFullSyncData(kit, subRes, syncData)
			if err != nil {
				blog.Errorf("handle %s-%s full sync data failed, err: %v, data: %+v, rid: %s", s.lgc.ResType(), subRes,
					err, *syncData, kit.Rid)
				return true, err
			}

			return false, nil
		})

		syncData.Start = nextStart
		syncData.Data = remainingData
	}

	return syncInfo.Total != 0, nil
}

// handleFullSyncData handle full sync data
func (s *resSyncer) handleFullSyncData(kit *util.Kit, subRes string, syncData *types.FullSyncTransData) (
	bool, map[string]int64, any, error) {

	resType := s.lgc.ResType()

	// list data of the corresponding interval
	listOpt := &types.ListDataOpt{
		SubRes: subRes,
		Start:  syncData.Start,
		End:    syncData.End,
	}
	listRes, err := s.lgc.ListData(kit, listOpt)
	if err != nil {
		blog.Errorf("list %s data failed, err: %v, opt: %+v, rid: %s", resType, err, *listOpt, kit.Rid)
		// start from the next interval
		nextStart := make(map[string]int64)
		for field, id := range syncData.Start {
			nextStart[field] = id + 1
		}
		return false, nextStart, syncData.Data, err
	}

	// cross compare data from two environments of the same interval
	compRes, err := s.lgc.CompareData(kit, subRes, syncData, listRes)
	if err != nil {
		blog.Errorf("compare %s data failed, err: %v, src: %+v, dest: %+v, rid: %s", resType, err, *syncData,
			listRes.Data, kit.Rid)
		return listRes.IsAll, listRes.NextStart, syncData.Data, err
	}

	// insert/update/delete data by compare result
	err = s.lgc.DeleteData(kit, subRes, compRes.Delete)
	if err != nil {
		blog.Errorf("delete %s-%s data(%+v) failed, err: %v, rid: %s", resType, subRes, compRes.Delete, err, kit.Rid)
		return listRes.IsAll, listRes.NextStart, compRes.RemainingSrc, err
	}

	err = s.lgc.UpdateData(kit, subRes, compRes.Update)
	if err != nil {
		blog.Errorf("update %s-%s data(%+v) failed, err: %v, rid: %s", resType, subRes, compRes.Update, err, kit.Rid)
		return listRes.IsAll, listRes.NextStart, compRes.RemainingSrc, err
	}

	err = s.lgc.InsertData(kit, subRes, compRes.Insert)
	if err != nil {
		blog.Errorf("insert %s-%s data(%+v) failed, err: %v, rid: %s", resType, subRes, compRes.Insert, err, kit.Rid)
		return listRes.IsAll, listRes.NextStart, compRes.RemainingSrc, err
	}

	return listRes.IsAll, listRes.NextStart, compRes.RemainingSrc, nil
}
