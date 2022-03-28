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

package y3_8_202004291536

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

// createIndex create non-exist index for table
func createIndex(ctx context.Context, db dal.RDB, tableName string, createIndexArr []types.Index) error {
	existIndexArr, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		return fmt.Errorf("createIndex failed, Indexes failed, tableName: %s, err:%s", tableName, err.Error())
	}
	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = true
	}
	for _, index := range createIndexArr {
		if _, ok := existIdxMap[index.Name]; ok == true {
			continue
		}
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return fmt.Errorf("createIndex failed, tableName: %s, index: %+v, err:%s", tableName, index, err.Error())
		}
	}
	return nil
}

// CreateServiceTemplateIndex create service template table index
func CreateServiceTemplateIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	createIndexArr := []types.Index{
		{
			Keys: map[string]int32{
				common.BKFieldID: 1,
			},
			Name:       "idx_id",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKAppIDField: 1,
			},
			Name:       "idx_bkBizID",
			Background: true,
		},
	}
	return createIndex(ctx, db, common.BKTableNameServiceTemplate, createIndexArr)
}

// CreateProcessTemplateIndex create process template table index
func CreateProcessTemplateIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	createIndexArr := []types.Index{
		{
			Keys: map[string]int32{
				common.BKServiceTemplateIDField: 1,
			},
			Name:       "idx_serviceTemplateID",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKFieldID: 1,
			},
			Name:       "idx_id",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKAppIDField: 1,
			},
			Name:       "idx_bkBizID",
			Background: true,
		},
	}
	return createIndex(ctx, db, common.BKTableNameProcessTemplate, createIndexArr)
}

// CreateServiceInstanceIndex create service instance table index
func CreateServiceInstanceIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	createIndexArr := []types.Index{
		{
			Keys: map[string]int32{
				common.BKFieldID: 1,
			},
			Name:       "idx_id",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKAppIDField: 1,
			},
			Name:       "idx_bkBizID",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKServiceTemplateIDField: 1,
			},
			Name:       "idx_serviceTemplateID",
			Background: true,
		},
	}
	return createIndex(ctx, db, common.BKTableNameServiceInstance, createIndexArr)
}

// CreateProcessInstanceRelationIndex create process instance relation table index
func CreateProcessInstanceRelationIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	createIndexArr := []types.Index{
		{
			Keys: map[string]int32{
				common.BKServiceInstanceIDField: 1,
			},
			Name:       "idx_bkServiceInstanceID",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKProcessTemplateIDField: 1,
			},
			Name:       "idx_bkProcessTemplateID",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKAppIDField: 1,
			},
			Name:       "idx_bkBizID",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKProcessIDField: 1,
			},
			Name:       "idx_bkProcessID",
			Background: true,
		},
	}
	return createIndex(ctx, db, common.BKTableNameProcessInstanceRelation, createIndexArr)
}
