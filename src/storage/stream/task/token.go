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

package task

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream/types"
)

type dbTokenHandler struct {
	uuid         string
	watchDB      local.DB
	dbWatchTasks []*dbWatchTask
}

// newDBTokenHandler new token handler for db watch task
func newDBTokenHandler(uuid string, watchDB local.DB, taskMap map[string]*dbWatchTask) types.TokenHandler {
	dbWatchTasks := make([]*dbWatchTask, 0, len(taskMap))
	for _, task := range taskMap {
		dbWatchTasks = append(dbWatchTasks, task)
	}
	return &dbTokenHandler{
		uuid:         uuid,
		watchDB:      watchDB,
		dbWatchTasks: dbWatchTasks,
	}
}

// SetLastWatchToken set last watch token for db watch task
func (d *dbTokenHandler) SetLastWatchToken(ctx context.Context, token *types.TokenInfo) error {
	// update last token for db to the earliest last token of all db watch tasks
	// this token specifies the last event that all db watch tasks has handled
	lastToken := d.dbWatchTasks[0].lastToken
	for _, task := range d.dbWatchTasks {
		if lastToken.Token > task.lastToken.Token {
			lastToken = task.lastToken
			continue
		}
	}

	filter := map[string]interface{}{
		"_id": d.uuid,
	}

	data := map[string]interface{}{
		common.BKTokenField:       lastToken.Token,
		common.BKStartAtTimeField: lastToken.StartAtTime,
	}

	if err := d.watchDB.Table(common.BKTableNameWatchToken).Update(ctx, filter, data); err != nil {
		blog.Errorf("set db %s last watch token failed, err: %v, data: %+v", d.uuid, err, data)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token of db watch task
func (d *dbTokenHandler) GetStartWatchToken(ctx context.Context) (*types.TokenInfo, error) {
	filter := map[string]interface{}{
		"_id": d.uuid,
	}

	data := new(types.TokenInfo)
	err := d.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(common.BKStartAtTimeField).One(ctx, data)
	if err != nil {
		if !d.watchDB.IsNotFoundError(err) {
			blog.Errorf("get db %s last watch token failed, err: %v", d.uuid, err)
			return nil, err
		}
		return new(types.TokenInfo), nil
	}
	return data, nil
}

// ResetWatchToken set watch token to empty and set the start watch time to the given one for next watch
func (d *dbTokenHandler) ResetWatchToken(startAtTime types.TimeStamp) error {
	filter := map[string]interface{}{
		"_id": d.uuid,
	}

	data := map[string]interface{}{
		common.BKTokenField:       "",
		common.BKStartAtTimeField: startAtTime,
	}

	if err := d.watchDB.Table(common.BKTableNameWatchToken).Update(context.Background(), filter, data); err != nil {
		blog.Errorf("reset db %s watch token failed, err: %v, data: %+v", d.uuid, err, data)
		return err
	}
	return nil
}
