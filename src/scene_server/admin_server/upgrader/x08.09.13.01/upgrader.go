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
	"gopkg.in/mgo.v2"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func addOperationLogIndex(db storage.DI, conf *upgrader.Config) (err error) {
	session := db.GetSession().(*mgo.Session)

	col := session.DB(db.GetDBName()).C(common.BKTableNameOperationLog)

	indexs, err := col.Indexes()
	if err != nil {
		return err
	}

	for _, index := range indexs {
		blog.V(3).Infof("droping index %s", index.Name)
		if index.Name == "_id_" {
			continue
		}
		if err = col.DropIndexName(index.Name); err != nil {
			return err
		}
	}

	idxs := []mgo.Index{
		{Name: "op_target_1_inst_id_1_op_time_-1", Key: []string{"op_target", "inst_id", "-op_time"}, Background: true},
		{Name: "bk_supplier_account_1_op_time_-1", Key: []string{"bk_supplier_account", "-op_time"}, Background: true},
		{Name: "bk_biz_id_1_bk_supplier_account_1_op_time_-1", Key: []string{"bk_biz_id", "bk_supplier_account", "-op_time"}, Background: true},
		{Name: "ext_key_1_bk_supplier_account_1_op_time_-1", Key: []string{"ext_key", "bk_supplier_account", "-op_time"}, Background: true},
	}
	for _, idx := range idxs {
		if err = col.EnsureIndex(idx); err != nil {
			return err
		}
	}
	return nil
}

func reconcileOperationLog(db storage.DI, conf *upgrader.Config) (err error) {
	all := []mapstr.MapStr{}

	if err = db.GetMutilByCondition(common.BKTableNameUserGroupPrivilege, nil, nil, &all, "", 0, 0); err != nil {
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
		if _, err = db.Insert(common.BKTableNameUserGroupPrivilege, privilege); err != nil {
			return err
		}
	}

	if err = db.DelByCondition(common.BKTableNameUserGroupPrivilege, map[string]interface{}{
		flag: map[string]interface{}{
			"$ne": true,
		},
	}); err != nil {
		return err
	}

	if err = db.DropColumn(common.BKTableNameUserGroupPrivilege, flag); err != nil {
		return err
	}

	return nil

}
