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
	"time"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/transfer-service/sync/medium"
	syncmeta "configcenter/src/source_controller/transfer-service/sync/metadata"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream"
)

// Watcher is cmdb data syncer event watcher
type Watcher struct {
	name          string
	loopW         stream.LoopInterface
	isMaster      discovery.ServiceManageInterface
	metadata      *syncmeta.Metadata
	cacheCli      cacheservice.CacheServiceClientInterface
	transMedium   medium.ClientI
	tokenHandlers map[types.ResType]*tokenHandler
}

// New new cmdb data syncer event watcher
func New(name string, loopW stream.LoopInterface, isMaster discovery.ServiceManageInterface, meta *syncmeta.Metadata,
	cacheCli cacheservice.CacheServiceClientInterface, transMedium medium.ClientI) (*Watcher, error) {

	// create cmdb data syncer event watch token table
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	exists, err := mongodb.Client("watch").HasTable(ctx, tokenTable)
	if err != nil {
		blog.Errorf("check if %s table exists failed, err: %v", tokenTable, err)
		return nil, err
	}

	if !exists {
		err = mongodb.Client("watch").CreateTable(ctx, tokenTable)
		if err != nil && !mongodb.Client("watch").IsDuplicatedError(err) {
			blog.Errorf("create %s table failed, err: %v", tokenTable, err)
			return nil, err
		}

		for _, resType := range types.ListAllResTypeForIncrSync() {
			token := &tokenInfo{
				Resource:    resType,
				StartAtTime: &metadata.Time{Time: time.Now()},
			}

			err = mongodb.Client("watch").Table(tokenTable).Insert(ctx, token)
			if err != nil && !mongodb.Client("watch").IsDuplicatedError(err) {
				blog.Errorf("init %s watch token failed, data: %+v, err: %v", resType, token, err)
				return nil, err
			}
		}
	}

	// generate cmdb data syncer event watcher
	watcher := &Watcher{
		name:          name,
		loopW:         loopW,
		isMaster:      isMaster,
		metadata:      meta,
		cacheCli:      cacheCli,
		transMedium:   transMedium,
		tokenHandlers: make(map[types.ResType]*tokenHandler),
	}

	for _, resType := range types.ListAllResTypeForIncrSync() {
		watcher.tokenHandlers[resType] = newTokenHandler(resType)
	}

	return watcher, nil
}

// Watch cmdb data syncer events and push the events to transfer medium
func (w *Watcher) Watch() error {
	for _, resType := range types.ListAllResTypeForIncrSync() {
		cursorTypes, exists := resTypeCursorMap[resType]
		if exists {
			for _, cursorType := range cursorTypes {
				go w.watchAPI(resType, cursorType)
			}
			continue
		}

		_, exists = resTypeWatchOptMap[resType]
		if !exists {
			continue
		}

		if err := w.watchDB(resType); err != nil {
			blog.Errorf("watch %s events from db failed, err: %v", resType, err)
			return err
		}
	}

	return nil
}
