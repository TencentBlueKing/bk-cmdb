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

package x08_09_13_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addOperationLogIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	indexs, err := db.Table(common.BKTableNameOperationLog).Indexes(ctx)
	if err != nil {
		return err
	}

	for _, index := range indexs {
		blog.V(3).Infof("droping index %s", index.Name)
		if index.Name == "_id_" {
			continue
		}
		if err = db.Table(common.BKTableNameOperationLog).DropIndex(ctx, index.Name); err != nil {
			return err
		}
	}

	idxs := []dal.Index{
		{Name: "op_target_1_inst_id_1_op_time_-1", Keys: map[string]interface{}{"op_target": 1, "inst_id": 1, "op_time": -1}, Background: true},
		{Name: "bk_supplier_account_1_op_time_-1", Keys: map[string]interface{}{"bk_supplier_account": 1, "op_time": -1}, Background: true},
		{Name: "bk_biz_id_1_bk_supplier_account_1_op_time_-1", Keys: map[string]interface{}{"bk_biz_id": 1, "bk_supplier_account": 1, "op_time": -1}, Background: true},
		{Name: "ext_key_1_bk_supplier_account_1_op_time_-1", Keys: map[string]interface{}{"ext_key": 1, "bk_supplier_account": 1, "op_time": -1}, Background: true},
	}
	for _, idx := range idxs {
		if err = db.Table(common.BKTableNameOperationLog).CreateIndex(ctx, idx); err != nil {
			return err
		}
	}
	return nil
}

func reconcileOperationLog(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	all := []mapstr.MapStr{}

	if err = db.Table(common.BKTableNameUserGroupPrivilege).Find(nil).All(ctx, &all); nil != err {
		return err
	}
	flag := "updateflag"
	expectM := map[string]mapstr.MapStr{}
	for _, privilege := range all {
		groupID, err := privilege.String("group_id")
		if err != nil {
			return err
		}
		privilege.Set(flag, true)
		expectM[groupID] = privilege
	}

	for _, privilege := range expectM {
		if err = db.Table(common.BKTableNameUserGroupPrivilege).Insert(ctx, privilege); nil != err {
			return err
		}
	}

	if err = db.Table(common.BKTableNameUserGroupPrivilege).Delete(ctx, map[string]interface{}{
		flag: map[string]interface{}{
			"$ne": true,
		},
	}); err != nil {
		return err
	}

	if err = db.Table(common.BKTableNameUserGroupPrivilege).DropColumn(ctx, flag); err != nil {
		return err
	}

	return nil

}
