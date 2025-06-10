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

// Package watch defines the cmdb data syncer watch logics
package watch

import (
	"context"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/transfer-service/sync/medium"
	syncmeta "configcenter/src/source_controller/transfer-service/sync/metadata"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/task"
)

// Watcher is cmdb data syncer event watcher
type Watcher struct {
	name          string
	isMaster      discovery.ServiceManageInterface
	metadata      *syncmeta.Metadata
	cacheCli      cacheservice.CacheServiceClientInterface
	transMedium   medium.ClientI
	tokenHandlers map[types.ResType]*tokenHandler
	tenantMap     map[string]string
}

// New new cmdb data syncer event watcher
func New(name string, tenantMap map[string]string, isMaster discovery.ServiceManageInterface, meta *syncmeta.Metadata,
	cacheCli cacheservice.CacheServiceClientInterface, transMedium medium.ClientI) (*Watcher, error) {

	// create cmdb data syncer event watch token and cursor table
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	err := mongodb.Dal("watch").ExecForAllDB(func(db local.DB) error {
		for _, table := range []string{tokenTable, cursorTable} {
			exists, err := db.HasTable(ctx, table)
			if err != nil {
				blog.Errorf("check if %s table exists failed, err: %v", table, err)
				return err
			}

			if !exists {
				err = db.CreateTable(ctx, table)
				if err != nil && !mongodb.IsDuplicatedError(err) {
					blog.Errorf("create %s table failed, err: %v", table, err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// generate cmdb data syncer event watcher
	watcher := &Watcher{
		name:          name,
		isMaster:      isMaster,
		metadata:      meta,
		cacheCli:      cacheCli,
		transMedium:   transMedium,
		tokenHandlers: make(map[types.ResType]*tokenHandler),
		tenantMap:     tenantMap,
	}

	for _, resType := range types.ListAllResTypeForIncrSync() {
		watcher.tokenHandlers[resType] = newTokenHandler(resType)
	}

	return watcher, nil
}

// Watch cmdb data syncer events and push the events to transfer medium
func (w *Watcher) Watch() ([]*task.Task, error) {
	tasks := make([]*task.Task, 0)
	for _, resType := range types.ListAllResTypeForIncrSync() {
		cursorTypes, exists := resTypeCursorMap[resType]
		if exists {
			for _, cursorType := range cursorTypes {
				for tenantID := range w.tenantMap {
					kit := rest.NewKit().WithTenant(tenantID)
					go w.watchAPI(kit, resType, cursorType)
				}
			}
			continue
		}

		watchTask, err := w.watchDB(resType)
		if err != nil {
			blog.Errorf("new watch %s events task failed, err: %v", resType, err)
			return nil, err
		}
		tasks = append(tasks, watchTask)
	}

	return tasks, nil
}
