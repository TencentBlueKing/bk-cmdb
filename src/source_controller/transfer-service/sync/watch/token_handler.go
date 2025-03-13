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

package watch

import (
	"context"
	"time"

	synctypes "configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

const (
	tokenTable  = "SrcSyncDataToken"
	cursorTable = "SrcSyncDataCursor"
)

var _ types.TaskTokenHandler = new(tokenHandler)

// tokenHandler is cmdb data syncer event token handler
type tokenHandler struct {
	resource synctypes.ResType
}

// newTokenHandler create a new cmdb data syncer event token handler
func newTokenHandler(resource synctypes.ResType) *tokenHandler {
	return &tokenHandler{resource: resource}
}

// SetLastWatchToken set last event watch token
func (t *tokenHandler) SetLastWatchToken(ctx context.Context, uuid string, watchDB local.DB,
	token *types.TokenInfo) error {

	tokenData := mapstr.MapStr{
		common.BKTokenField:       token,
		common.BKStartAtTimeField: token.StartAtTime,
	}
	filter := map[string]interface{}{
		"resource": watch.GenDBWatchTokenID(uuid, string(t.resource)),
	}

	if err := watchDB.Table(tokenTable).Upsert(ctx, filter, tokenData); err != nil {
		blog.Errorf("set %s watch token info failed, data: %+v, err: %v", t.resource, tokenData, err)
		return err
	}
	return nil
}

// GetStartWatchToken get event start watch token
func (t *tokenHandler) GetStartWatchToken(ctx context.Context, uuid string, watchDB local.DB) (*types.TokenInfo,
	error) {

	filter := map[string]interface{}{
		"resource": watch.GenDBWatchTokenID(uuid, string(t.resource)),
	}

	info := new(types.TokenInfo)
	if err := watchDB.Table(tokenTable).Find(filter).One(ctx, &info); err != nil {
		if mongodb.IsNotFoundError(err) {
			return &types.TokenInfo{Token: "", StartAtTime: &types.TimeStamp{Sec: uint32(time.Now().Unix())}}, nil
		}
		blog.Errorf("get %s event watch token info failed, err: %v", t.resource, err)
		return nil, err
	}

	return info, nil
}

// cursorInfo is cmdb data syncer event token info
type cursorInfo struct {
	Resource    synctypes.ResType           `bson:"resource"`
	Cursor      map[watch.CursorType]string `bson:"cursor"`
	StartAtTime *metadata.Time              `bson:"start_at_time"`
}

// getWatchCursorInfo get event watch token info
func (t *tokenHandler) getWatchCursorInfo(kit *rest.Kit) (*cursorInfo, error) {
	filter := map[string]interface{}{
		"resource": t.resource,
	}

	info := new(cursorInfo)
	err := mongodb.Dal("watch").Shard(kit.ShardOpts()).Table(cursorTable).Find(filter).One(kit.Ctx, &info)
	if err != nil {
		if mongodb.IsNotFoundError(err) {
			return new(cursorInfo), nil
		}
		blog.Errorf("get %s event watch token info failed, err: %v, rid: %s", t.resource, err, kit.Rid)
		return nil, err
	}

	if info.Cursor == nil {
		info.Cursor = make(map[watch.CursorType]string)
	}

	return info, nil
}

// setWatchCursorInfo get event watch token info
func (t *tokenHandler) setWatchCursorInfo(kit *rest.Kit, data mapstr.MapStr) error {
	filter := map[string]interface{}{
		"resource": t.resource,
	}

	if err := mongodb.Dal("watch").Shard(kit.ShardOpts()).Table(cursorTable).Upsert(kit.Ctx, filter, data); err != nil {
		blog.Errorf("set %s watch token info failed, data: %+v, err: %v, rid: %s", t.resource, data, err, kit.Rid)
		return err
	}

	return nil
}
