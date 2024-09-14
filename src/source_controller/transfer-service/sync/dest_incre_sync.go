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
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/transfer-service/sync/util"
)

// loopPullIncrSyncData loop pull incremental sync data
func (s *Syncer) loopPullIncrSyncData(resType types.ResType) {
	ack := false

	for {
		if !s.isMaster.IsMaster() {
			blog.V(4).Infof("loop pull %s incremental sync data, but not master, skip.", resType)
			time.Sleep(time.Minute)
			ack = false
			continue
		}

		syncer := s.resSyncerMap[resType]

		hasMore, err := syncer.pullIncrSyncData(ack)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		if !hasMore {
			time.Sleep(2 * time.Second)
		}

		ack = true
	}
}

// pullIncrSyncData pull incremental sync data for one resource
func (s *resSyncer) pullIncrSyncData(ack bool) (bool, error) {
	kit := util.NewKit()
	resType := s.lgc.ResType()
	blog.Infof("start pull %s incr sync data, rid: %s", resType, kit.Rid)

	// pull incr sync data from transfer medium
	pullOpt := &types.PullSyncDataOpt{
		ResType:     resType,
		IsIncrement: true,
		Ack:         ack,
	}
	syncInfo, err := s.transMedium.PullSyncData(kit.Ctx, kit.Header, pullOpt)
	if err != nil {
		blog.Errorf("pull %s incr sync data failed, err: %v, rid: %s", resType, err, kit.Rid)
		return false, err
	}

	if len(syncInfo.Info) == 0 {
		return syncInfo.Total != 0, nil
	}

	// parse incr sync data, skip the invalid data
	syncData := new(types.IncrSyncTransData)
	err = json.Unmarshal(syncInfo.Info, &syncData)
	if err != nil {
		blog.Errorf("unmarshal %s incr sync data(%s) failed, err: %v, rid: %s", resType, syncInfo.Info, err, kit.Rid)
		return false, err
	}

	// delete data
	for subRes, rawData := range syncData.DeleteInfo {
		dataArr, err := s.lgc.ParseDataArr(syncData.Name, subRes, rawData, kit.Rid)
		if err != nil {
			blog.Errorf("parse %s-%s incr sync data(%+v) failed, err: %v, rid: %s", s.lgc.ResType(), subRes, rawData,
				err, kit.Rid)
			return syncInfo.Total != 0, nil
		}

		if err = s.lgc.DeleteData(kit, subRes, dataArr); err != nil {
			blog.Errorf("delete %s-%s data(%+v) failed, err: %v, rid: %s", resType, subRes, dataArr, err, kit.Rid)
			return false, err
		}
	}

	// cross compare data from two environments of the same interval
	for subRes, rawData := range syncData.UpsertInfo {
		dataArr, err := s.lgc.ParseDataArr(syncData.Name, subRes, rawData, kit.Rid)
		if err != nil {
			blog.Errorf("parse %s-%s incr sync data(%+v) failed, err: %v, rid: %s", s.lgc.ResType(), subRes, rawData,
				err, kit.Rid)
			return syncInfo.Total != 0, nil
		}

		insertData, updateData, err := s.lgc.ClassifyUpsertData(kit, subRes, dataArr)
		if err != nil {
			blog.Errorf("classify %s upsert data(%+v) failed, err: %v, rid: %s", resType, dataArr, err, kit.Rid)
			return false, err
		}

		if err = s.lgc.UpdateData(kit, subRes, updateData); err != nil {
			blog.Errorf("update %s-%s data(%+v) failed, err: %v, rid: %s", resType, subRes, updateData, err, kit.Rid)
			return false, err
		}

		if err = s.lgc.InsertData(kit, subRes, insertData); err != nil {
			blog.Errorf("insert %s-%s data(%+v) failed, err: %v, rid: %s", resType, subRes, insertData, err, kit.Rid)
			return false, err
		}
	}

	blog.Infof("pull %s incr sync data successfully, rid: %s", resType, kit.Rid)
	return syncInfo.Total != 0, nil
}
