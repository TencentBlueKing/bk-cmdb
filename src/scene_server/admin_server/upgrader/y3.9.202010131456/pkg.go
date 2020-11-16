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

package y3_9_202010131456

import (
	"context"
	"time"

	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// UserConfig is old dynamic groupping metadata struct, now only used in migrate.
type UserConfig struct {
	Info       string    `json:"info" bson:"info"`
	Name       string    `json:"name" bson:"name"`
	ID         string    `json:"id" bson:"id"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	UpdateTime time.Time `json:"last_time" bson:"last_time"`
	AppID      int64     `json:"bk_biz_id" bson:"bk_biz_id"`
	CreateUser string    `json:"create_user" bson:"create_user"`
	ModifyUser string    `json:"modify_user" bson:"modify_user"`
}

func init() {
	upgrader.RegistUpgrader("y3.9.202010131456", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	err = createTable(ctx, db, conf)
	if err != nil {
		return err
	}

	err = migrateHistory(ctx, db, conf)
	if err != nil {
		return err
	}
	return
}
