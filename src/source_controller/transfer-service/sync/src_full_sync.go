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
	"time"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/lock"
	"configcenter/src/source_controller/transfer-service/sync/util"
	"configcenter/src/storage/driver/redis"
)

// loopPushFullSyncData loop push full sync data
func (s *Syncer) loopPushFullSyncData(interval time.Duration) {
	time.Sleep(15 * time.Minute)

	for {
		if !s.isMaster.IsMaster() {
			blog.V(4).Infof("loop push full sync data, but not master, skip")
			time.Sleep(5 * time.Minute)
			continue
		}

		locker := lock.NewLocker(redis.Client())
		locked, err := locker.Lock(types.FullSyncLockKey, time.Hour)
		if err != nil || !locked {
			blog.Errorf("do not get %s lock, err: %v, locked: %v", types.FullSyncLockKey, err, locked)
			time.Sleep(5 * time.Minute)
			continue
		}

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
					syncer.pushFullSyncData(objID)
				}
			case types.InstAsst:
				// TODO 目前没有同步业务相关的关联，因为现在产品形态是不支持的，后续如果需要支持的话实例关联同步都需要调整
				for _, objID := range append(objIDs, common.BKInnerObjIDHost) {
					syncer.pushFullSyncData(objID)
				}
			case types.QuotedInstance:
				for _, objID := range quotedObjIDs {
					syncer.pushFullSyncData(objID)
				}
			default:
				syncer.pushFullSyncData("")
			}
		}

		locker.Unlock()

		time.Sleep(interval)
	}
}

// pushFullSyncData push full sync data for one resource
func (s *resSyncer) pushFullSyncData(subRes string) {
	kit := util.NewKit()
	startTime := time.Now()
	blog.Infof("start push %s-%s full sync data, start time: %s, rid: %s", s.lgc.ResType(), subRes, startTime, kit.Rid)

	isAll := false
	start := make(map[string]int64)
	var err error

	for !isAll {
		var nextStart map[string]int64

		util.RetryWrapper(3, func() (bool, error) {
			isAll, nextStart, err = s.doOnePushFullSyncDataStep(kit, subRes, start, nil)
			if err != nil {
				blog.Errorf("try %s-%s full sync step failed, err: %v, start: %+v, rid: %s", s.lgc.ResType(), subRes,
					err, start, kit.Rid)
				return true, err
			}
			return false, nil
		})

		start = nextStart
	}

	blog.Infof("push %s-%s full sync data successfully, start time: %s, cost: %s, rid: %s", s.lgc.ResType(), subRes,
		time.Now(), time.Since(startTime), kit.Rid)
}

// doOnePushFullSyncDataStep do one push full sync data step
func (s *resSyncer) doOnePushFullSyncDataStep(kit *util.Kit, subRes string, start, end map[string]int64) (bool,
	map[string]int64, error) {

	// list data from the start index
	listOpt := &types.ListDataOpt{
		SubRes: subRes,
		Start:  start,
		End:    end,
	}

	info, err := s.lgc.ListData(kit, listOpt)
	if err != nil {
		blog.Errorf("list %s data failed, err: %v, opt: %+v, rid: %s", s.lgc.ResType(), err, *listOpt, kit.Rid)
		// start from the next interval
		nextStart := make(map[string]int64)
		for field, id := range start {
			nextStart[field] = id + 1
		}
		return false, nextStart, err
	}

	// all data has been listed, do not have a sync interval end
	syncEnd := info.NextStart
	if info.IsAll {
		syncEnd = make(map[string]int64)
	}

	// push full sync data to transfer medium
	pushOpt := &types.PushSyncDataOpt{
		ResType:     s.lgc.ResType(),
		SubRes:      subRes,
		IsIncrement: false,
		Data: &types.FullSyncTransData{
			Name:  s.name,
			Start: start,
			End:   syncEnd,
			Data:  info.Data,
		},
	}
	err = s.transMedium.PushSyncData(kit.Ctx, kit.Header, pushOpt)
	if err != nil {
		blog.Errorf("push %s-%s full sync data failed, err: %v, opt: %+v, rid: %s", s.lgc.ResType(), subRes, err,
			*pushOpt, kit.Rid)
		return info.IsAll, info.NextStart, err
	}

	return info.IsAll, info.NextStart, nil
}
