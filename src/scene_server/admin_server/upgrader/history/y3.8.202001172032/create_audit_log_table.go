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

package y3_8_202001172032

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func createAuditLogTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	exists, err := db.HasTable(ctx, common.BKTableNameAuditLog)
	if err != nil {
		blog.ErrorJSON("search audit log table error, err:%s", err.Error())
		return err
	}
	if !exists {
		if err = db.CreateTable(ctx, common.BKTableNameAuditLog); err != nil && !db.IsDuplicatedError(err) {
			blog.ErrorJSON("create audit log table error, err:%s", err.Error())
			return err
		}
	}
	return nil
}

func addAuditLogTableIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idxArr, err := db.Table(common.BKTableNameAuditLog).Indexes(ctx)
	if err != nil {
		blog.Errorf("get table %s index error. err:%s", common.BKTableNameAuditLog, err.Error())
		return err
	}

	createIdxArr := []types.Index{
		{Name: "index_bk_supplier_account", Keys: map[string]int32{common.BkSupplierAccount: 1}, Background: true},
		{Name: "index_audit_type", Keys: map[string]int32{common.BKAuditTypeField: 1}, Background: true},
		{Name: "index_action", Keys: map[string]int32{common.BKActionField: 1}, Background: true},
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
			continue
		}
		if err := db.Table(common.BKTableNameAuditLog).CreateIndex(ctx, idx); err != nil {
			blog.ErrorJSON("create index to BKTableNameAuditLog error, err:%s, current index:%s, all create index:%s", err.Error(), idx, createIdxArr)
			return err
		}

	}

	return nil
}
