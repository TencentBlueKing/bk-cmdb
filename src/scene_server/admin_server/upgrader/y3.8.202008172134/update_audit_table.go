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

package y3_8_202008172134

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

// addAuditTableIndexes add indexes for common audit log query params
func addAuditTableIndexes(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	indexes := []types.Index{
		{Name: "index_id", Keys: map[string]int32{common.BKFieldID: 1}, Background: true},
		{Name: "index_operation_time", Keys: map[string]int32{common.BKOperationTimeField: 1}, Background: true},
		{Name: "index_bk_supplier_account", Keys: map[string]int32{common.BkSupplierAccount: 1}, Background: true},
		{Name: "index_audit_type", Keys: map[string]int32{common.BKAuditTypeField: 1}, Background: true},
		{Name: "index_action", Keys: map[string]int32{common.BKActionField: 1}, Background: true},
	}

	existIndexArr, err := db.Table(common.BKTableNameAuditLog).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for audit table failed, err: %s", err.Error())
		return err
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = true
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		if err = db.Table(common.BKTableNameAuditLog).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for audit table failed, err: %s, index: %+v", err.Error(), index)
			return err
		}
	}
	return nil
}
