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

	synctypes "configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

const tokenTable = "cc_SrcSyncDataToken"

var _ types.TokenHandler = new(tokenHandler)

// tokenHandler is cmdb data syncer event token handler
type tokenHandler struct {
	resource synctypes.ResType
}

// newTokenHandler create a new cmdb data syncer event token handler
func newTokenHandler(resource synctypes.ResType) *tokenHandler {
	return &tokenHandler{resource: resource}
}

// tokenInfo is cmdb data syncer event token info
type tokenInfo struct {
	Resource    synctypes.ResType           `bson:"resource"`
	Token       string                      `bson:"token"`
	Cursor      map[watch.CursorType]string `bson:"cursor"`
	StartAtTime *metadata.Time              `bson:"start_at_time"`
}

// SetLastWatchToken set last event watch token
func (t *tokenHandler) SetLastWatchToken(ctx context.Context, token string) error {
	tokenData := mapstr.MapStr{
		common.BKTokenField: token,
	}
	return t.setWatchTokenInfo(ctx, tokenData)
}

// GetStartWatchToken get event start watch token
func (t *tokenHandler) GetStartWatchToken(ctx context.Context) (string, error) {
	info, err := t.getWatchTokenInfo(ctx, common.BKTokenField)
	if err != nil {
		return "", err
	}

	return info.Token, nil
}

// resetWatchToken reset watch token and start watch time
func (t *tokenHandler) resetWatchToken(startAtTime types.TimeStamp) error {
	filter := map[string]interface{}{
		"resource": t.resource,
	}
	data := mapstr.MapStr{
		common.BKCursorField:      make(map[watch.CursorType]string),
		common.BKTokenField:       "",
		common.BKStartAtTimeField: startAtTime,
	}

	if err := mongodb.Client("watch").Table(tokenTable).Upsert(context.Background(), filter, data); err != nil {
		blog.Errorf("reset %s watch token failed, data: %+v, err: %v", t.resource, data, err)
		return err
	}
	return nil
}

// getWatchTokenInfo get event watch token info
func (t *tokenHandler) getWatchTokenInfo(ctx context.Context, fields ...string) (*tokenInfo, error) {
	filter := map[string]interface{}{
		"resource": t.resource,
	}

	info := new(tokenInfo)
	if err := mongodb.Client("watch").Table(tokenTable).Find(filter).Fields(fields...).One(ctx, &info); err != nil {
		if mongodb.Client("watch").IsNotFoundError(err) {
			return new(tokenInfo), nil
		}
		blog.Errorf("get %s event watch token info failed, err: %v", t.resource, err)
		return nil, err
	}

	if info.Cursor == nil {
		info.Cursor = make(map[watch.CursorType]string)
	}

	return info, nil
}

// getWatchTokenInfo get event watch token info
func (t *tokenHandler) setWatchTokenInfo(ctx context.Context, data mapstr.MapStr) error {
	filter := map[string]interface{}{
		"resource": t.resource,
	}

	if err := mongodb.Client("watch").Table(tokenTable).Upsert(ctx, filter, data); err != nil {
		blog.Errorf("set %s watch token info failed, data: %+v, err: %v", t.resource, data, err)
		return err
	}

	return nil
}
