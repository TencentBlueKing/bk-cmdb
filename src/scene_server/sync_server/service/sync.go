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
	"time"

	ftypes "configcenter/pkg/types/sync/full-text-search"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
)

// SyncFullTextSearchData sync data for full-text search, NOTE: this is an async api
func (s *Service) SyncFullTextSearchData(cts *rest.Contexts) {
	opt := new(ftypes.SyncDataOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	rawErr := opt.Validate()
	if rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	blog.Infof("start sync full-text search data request, opt: %+v, rid: %s", opt, cts.Kit.Rid)

	go func() {
		err := s.lgc.FullTextSearch.SyncData(context.Background(), opt, cts.Kit.Rid)
		if err != nil {
			blog.Errorf("run sync full-text search data req failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
			return
		}
		blog.Infof("finished sync full-text search data request, opt: %+v, rid: %s", opt, cts.Kit.Rid)
	}()

	cts.RespEntity(nil)
}

// MigrateFullTextSearch migrate full-text search info, NOTE: this is an async api
func (s *Service) MigrateFullTextSearch(cts *rest.Contexts) {
	blog.Infof("start migrate full-text search request, rid: %s", cts.Kit.Rid)

	var result *ftypes.MigrateResult
	done := make(chan struct{}, 1)

	go func() {
		res, err := s.lgc.FullTextSearch.Migrate(context.Background(), cts.Kit.Rid)
		if err != nil {
			blog.Errorf("run migrate full-text search req failed, err: %v, rid: %s", err, cts.Kit.Rid)
			return
		}
		result = res
		done <- struct{}{}

		blog.Infof("finished migrate full-text search request, opt: %+v, rid: %s", cts.Kit.Rid)
	}()

	tick := time.Tick(10 * time.Second)
	select {
	case <-tick:
		cts.RespEntity(&ftypes.MigrateResult{Message: "migrate full-text search task is running"})
	case <-done:
		cts.RespEntity(result)
	}
}
