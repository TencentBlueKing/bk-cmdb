/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package event

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream"
	"gopkg.in/redis.v5"
)

func NewEvent(db dal.DB, rds *redis.Client, watch stream.Interface) error {
	e := Event{
		rds:   rds,
		watch: watch,
		db:    db,
	}

	if err := e.runHost(context.Background()); err != nil {
		blog.Errorf("run host event flow failed, err: %v", err)
		return err
	}

	if err := e.runModuleHostRelation(context.Background()); err != nil {
		blog.Errorf("run module host config event flow failed, err: %v", err)
		return err
	}

	return nil
}

type Event struct {
	rds   *redis.Client
	watch stream.Interface
	db    dal.DB
}

func (e *Event) runHost(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameBaseHost,
		key:        HostKey,
		rds:        e.rds,
		watch:      e.watch,
		db:         e.db,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runModuleHostRelation(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameModuleHostConfig,
		key:        ModuleHostRelationKey,
		rds:        e.rds,
		watch:      e.watch,
		db:         e.db,
	}

	return newFlow(ctx, opts)
}
