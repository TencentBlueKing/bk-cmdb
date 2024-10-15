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
	"context"
	"errors"
	"time"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/lock"
	commonutil "configcenter/src/common/util"
	"configcenter/src/source_controller/transfer-service/sync/util"
	"configcenter/src/storage/driver/redis"
)

// SyncCmdbData sync cmdb data
func (s *Syncer) SyncCmdbData(kit *rest.Kit, opt *types.SyncCmdbDataOption) error {
	if !s.enableSync {
		return errors.New("sync is disabled")
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		return rawErr.ToCCError(kit.CCError)
	}

	syncer, exists := s.resSyncerMap[opt.ResType]
	if !exists {
		blog.Errorf("res type %s is invalid, rid: %s", opt.ResType, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "resource_type")
	}

	locker := lock.NewLocker(redis.Client())
	locked, err := locker.Lock(types.FullSyncLockKey, time.Hour)
	if err != nil {
		blog.Errorf("sync cmdb data but get lock failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if !locked {
		blog.Infof("sync cmdb data but there's another task running, rid: %s", kit.Rid)
		return errors.New("there is another sync task running")
	}

	blog.Infof("start sync cmdb data, opt: %+v, rid: %s", *opt, kit.Rid)

	kt := util.ConvertKit(kit)
	kt.Ctx = commonutil.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	go func() {
		defer locker.Unlock()

		isAll := false
		start, end := make(map[string]int64), make(map[string]int64)
		if !opt.IsAll {
			start = opt.Start
			end = opt.End
		}

		var err error

		for !isAll {
			isAll, start, err = syncer.doOnePushFullSyncDataStep(kt, opt.SubRes, start, end)
			if err != nil {
				blog.Errorf("sync %s-%s full sync step failed, err: %v, start: %+v, end: %+v, rid: %s", opt.ResType,
					opt.SubRes, err, start, end, kit.Rid)
				return
			}
		}

		blog.Infof("sync cmdb data successfully, opt: %+v, rid: %s", *opt, kit.Rid)
	}()

	return nil
}
