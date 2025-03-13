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
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

type dbTokenHandler struct {
	uuid    string
	watchDB local.DB

	lastToken     *types.TokenInfo
	taskMap       map[string]*dbWatchTask
	lastTokenInfo map[string]string
	mu            sync.RWMutex
}

// newDBTokenHandler new token handler for db watch task
func newDBTokenHandler(uuid string, watchDB local.DB, taskMap map[string]*dbWatchTask) (*dbTokenHandler, error) {
	handler := &dbTokenHandler{
		uuid:          uuid,
		watchDB:       watchDB,
		taskMap:       taskMap,
		lastTokenInfo: make(map[string]string),
	}

	lastToken, err := handler.GetStartWatchToken(context.Background())
	if err != nil {
		return nil, err
	}
	handler.lastToken = lastToken

	tokenChan := make(chan struct{})

	for taskID, task := range taskMap {
		if task.lastToken != nil {
			handler.lastTokenInfo[taskID] = task.lastToken.Token
		}
		task.tokenChan = tokenChan
	}

	go func() {
		for _ = range tokenChan {
			handler.setLastWatchToken()
		}
	}()
	return handler, nil
}

func (d *dbTokenHandler) setTaskLastTokenInfo(taskLastTokenMap map[string]string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for taskID, token := range taskLastTokenMap {
		d.lastTokenInfo[taskID] = token
	}
}

func (d *dbTokenHandler) setLastWatchToken() {
	// update last token for db to the earliest last token of all db watch tasks
	// this token specifies the last event that all db watch tasks has handled
	var lastToken *types.TokenInfo
	allFinished := false

	for taskID, task := range d.taskMap {
		token := task.lastToken

		// if token is nil, skip it
		if token == nil {
			continue
		}

		isFinished := true
		d.mu.RLock()
		if token.Token < d.lastTokenInfo[taskID] {
			isFinished = false
		}
		d.mu.RUnlock()

		if lastToken == nil {
			lastToken = token
			allFinished = isFinished
			continue
		}

		if allFinished {
			// if all other tasks are finished but this task is not finished, use the last token of the unfinished task
			if !isFinished {
				allFinished = false
				lastToken = token
				continue
			}

			// if all tasks are finished, use the last token of the latest finished task
			if lastToken.Token < token.Token {
				lastToken = token
			}
			continue
		}

		// if not all tasks are finished, skip the finished tasks
		if isFinished {
			continue
		}
		if lastToken.Token > token.Token {
			// use the last token of the earliest unfinished task
			lastToken = token
		}
	}

	// if no events are handled, do not update the last token
	if lastToken == nil || lastToken.Token == "" || lastToken.Token <= d.lastToken.Token {
		return
	}

	filter := map[string]interface{}{
		"_id": d.uuid,
	}

	data := map[string]interface{}{
		common.BKTokenField:       lastToken.Token,
		common.BKStartAtTimeField: lastToken.StartAtTime,
	}

	if err := d.watchDB.Table(common.BKTableNameWatchToken).Update(context.Background(), filter, data); err != nil {
		blog.Errorf("set db %s last watch token failed, err: %v, data: %+v", d.uuid, err, data)
		return
	}
	d.lastToken = lastToken
}

// SetLastWatchToken set last watch token for db watch task
func (d *dbTokenHandler) SetLastWatchToken(ctx context.Context, token *types.TokenInfo) error {
	return nil
}

// GetStartWatchToken get start watch token of db watch task
func (d *dbTokenHandler) GetStartWatchToken(ctx context.Context) (*types.TokenInfo, error) {
	filter := map[string]interface{}{
		"_id": d.uuid,
	}

	data := new(types.TokenInfo)
	err := d.watchDB.Table(common.BKTableNameWatchToken).Find(filter).One(ctx, data)
	if err != nil {
		if !mongodb.IsNotFoundError(err) {
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
