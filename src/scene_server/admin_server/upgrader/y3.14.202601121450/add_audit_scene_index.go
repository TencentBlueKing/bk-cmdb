/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package y3_14_202601121450

import (
	"configcenter/src/common/metadata"
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addAuditLogSceneIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idxArr, err := db.Table(common.BKTableNameAuditLog).Indexes(ctx)
	if err != nil {
		blog.Errorf("get table %s index error. err:%s", common.BKTableNameAuditLog, err.Error())
		return err
	}
	err = db.Table(common.BKTableNameAuditLog).AddColumn(ctx, "audit_context", metadata.AuditSceneHeader{})
	if err != nil {
		return fmt.Errorf("cc_AuditLog add column [audit_context] err:%w", err)
	}
	createIdxArr := []types.Index{
		{
			Keys: bson.D{
				{common.BKOperationTimeField, 1},
				{common.BKAuditSceneContextSceneTraceId, 1},
			},
			Name:               "index_audit_context_scene_trace_id",
			Unique:             true,
			Background:         true,
			ExpireAfterSeconds: 0,
			PartialFilterExpression: map[string]interface{}{
				common.BKAuditSceneContextSceneTraceId: bson.D{{common.BKDBExists, true}},
			},
		},
		{
			Name: "index_audit_context_op", Keys: bson.D{
				{common.BKOperationTimeField, 1},
				{common.BKAuditSceneContextScene, 1},
				{common.BKAuditSceneContextOpUser, 1},
				{common.BKAuditAppCodeField, 1},
			}, Background: true, Unique: false,
			PartialFilterExpression: map[string]interface{}{
				common.BKOperationTimeField:            bson.D{{common.BKDBExists, true}},
				common.BKAuditSceneContextSceneTraceId: bson.D{{common.BKDBExists, true}},
			}},
	}

	for _, idx := range createIdxArr {
		exist := false
		for _, existIdx := range idxArr {
			if existIdx.Name == idx.Name {
				exist = true
				break
			}
		}
		if exist {
			if err := db.Table(common.BKTableNameAuditLog).DropIndex(ctx, idx.Name); err != nil {
				blog.Errorf("add audit log index error. err:%s", err.Error())
				return err
			}
		}
		if err := db.Table(common.BKTableNameAuditLog).CreateIndex(ctx, idx); err != nil && !db.IsDuplicatedError(err) {
			blog.ErrorJSON("create index to BKTableNameAuditLog error, err:%s, current index:%s, "+
				"all create index:%s", err.Error(), idx, createIdxArr)
			return err
		}
	}

	return nil
}
