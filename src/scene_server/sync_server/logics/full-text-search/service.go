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

package fulltextsearch

import (
	"context"
	"errors"
	"sync"

	ftypes "configcenter/pkg/types/sync/full-text-search"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/upgrader"
)

// SyncData sync full-text search data
func (f *fullTextSearch) SyncData(kit *rest.Kit, opt *ftypes.SyncDataOption) error {
	if !f.enableSync {
		return errors.New("full text search sync is disabled")
	}

	// sync all data
	if opt.IsAll {
		var firstErr error
		pipeline := make(chan struct{}, 5)
		wg := sync.WaitGroup{}

		for _, index := range types.AllIndexNames {
			if firstErr != nil {
				break
			}

			pipeline <- struct{}{}
			wg.Add(1)

			go func(kit *rest.Kit, index string) {
				defer func() {
					<-pipeline
					wg.Done()
				}()

				err := f.syncDataByIndex(kit, index)
				if err != nil {
					firstErr = err
				}
			}(kit, index)
		}

		wg.Wait()
		return firstErr
	}

	if len(opt.Index) > 0 {
		return f.syncDataByIndex(kit, opt.Index)
	}

	// sync specific collection data
	index, err := getIndexByColl(opt.Collection)
	if err != nil {
		return err
	}

	_, err = f.syncCollection(kit, index, opt.Collection, opt.Oids)
	if err != nil {
		return err
	}
	return nil
}

// Migrate full-text search index info with its related data
func (f *fullTextSearch) Migrate(ctx context.Context, rid string) (*ftypes.MigrateResult, error) {
	if !f.enableSync {
		return nil, errors.New("full text search sync is disabled")
	}

	// upgrade index info
	migrateResult, err := upgrader.Upgrade(ctx, rid)
	if err != nil {
		blog.Errorf("migrate failed, err: %v, res: %v, rid: %s", err, migrateResult, rid)
		return nil, err
	}

	return migrateResult, nil
}
