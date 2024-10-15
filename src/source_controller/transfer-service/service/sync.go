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

package service

import (
	"context"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
)

// SyncCmdbData sync cmdb data
func (s *Service) SyncCmdbData(cts *rest.Contexts) {
	opt := new(types.SyncCmdbDataOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	blog.Infof("start sync cmdb data, opt: %+v, rid: %s", opt, cts.Kit.Rid)

	cts.Kit.Ctx = context.Background()
	go func() {
		err := s.syncer.SyncCmdbData(cts.Kit, opt)
		if err != nil {
			blog.Errorf("sync cmdb data failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
			return
		}
		blog.Infof("finished sync cmdb data, opt: %+v, rid: %s", opt, cts.Kit.Rid)
	}()

	cts.RespEntity(nil)
}
