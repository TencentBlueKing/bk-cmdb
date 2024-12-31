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

package host

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// LockHost TODO
func (hm *hostManager) LockHost(kit *rest.Kit, input *metadata.HostLockRequest) errors.CCError {
	input.IDS = util.IntArrayUnique(input.IDS)
	condition := mapstr.MapStr{
		common.BKHostIDField: mapstr.MapStr{common.BKDBIN: input.IDS},
	}
	hostInfos := make([]metadata.HostMapStr, 0)
	limit := uint64(len(input.IDS))
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(condition).Fields(
		common.BKHostIDField).Limit(limit).All(kit.Ctx, &hostInfos)
	if err != nil {
		blog.Errorf("lock host, query host from db error, condition: %+v, err: %v, rid: %s", condition, err,
			kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommDBSelectFailed)
	}

	diffID := diffHostLockID(input.IDS, hostInfos, kit.Rid)
	if 0 != len(diffID) {
		blog.Errorf("lock host, not found, id: %+v, rid: %s", diffID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, fmt.Sprintf(" id_list %v", diffID))
	}

	user := httpheader.GetUser(kit.Header)
	var insertDataArr []interface{}
	ts := time.Now().UTC()
	for _, id := range input.IDS {
		conds := mapstr.MapStr{
			common.BKHostIDField: id,
		}
		cnt, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostLock).Find(conds).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("query host lock from db failed, err: %v, rid: %s", err, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommDBSelectFailed)
		}
		if 0 == cnt {
			insertDataArr = append(insertDataArr, metadata.HostLockData{
				User:       user,
				ID:         id,
				CreateTime: ts,
			})
		}
	}

	if 0 < len(insertDataArr) {
		err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostLock).Insert(kit.Ctx, insertDataArr)
		if err != nil {
			blog.Errorf("save host lock to db failed, err: %v, rid: %s", err, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommDBInsertFailed)
		}
	}
	return nil
}

// UnlockHost TODO
func (hm *hostManager) UnlockHost(kit *rest.Kit, input *metadata.HostLockRequest) errors.CCError {
	conds := mapstr.MapStr{
		common.BKHostIDField: mapstr.MapStr{common.BKDBIN: input.IDS},
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostLock).Delete(kit.Ctx, conds)
	if err != nil {
		blog.Errorf("delete host lock from db failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

// QueryHostLock TODO
func (hm *hostManager) QueryHostLock(kit *rest.Kit, input *metadata.QueryHostLockRequest) ([]metadata.HostLockData,
	errors.CCError) {
	hostLockInfoArr := make([]metadata.HostLockData, 0)
	conds := mapstr.MapStr{
		common.BKHostIDField: mapstr.MapStr{common.BKDBIN: input.IDS},
	}
	limit := uint64(len(input.IDS))
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostLock).Find(conds).Limit(limit).
		All(kit.Ctx, &hostLockInfoArr)
	if err != nil {
		blog.Errorf("query host lock from db failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	return hostLockInfoArr, nil
}

func diffHostLockID(ids []int64, hostInfos []metadata.HostMapStr, rid string) []int64 {
	mapInnerID := make(map[int64]bool)
	for _, hostInfo := range hostInfos {
		id, err := util.GetInt64ByInterface(hostInfo[common.BKHostIDField])
		if err != nil {
			blog.Errorf("different host lock ID not valid, hostInfo: %v, rid: %s", hostInfo, rid)
			continue
		}
		mapInnerID[id] = true
	}
	var diffIDS []int64
	for _, id := range ids {
		_, exist := mapInnerID[id]
		if !exist {
			diffIDS = append(diffIDS, id)
		}
	}
	return diffIDS
}
