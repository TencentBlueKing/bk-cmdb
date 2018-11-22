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
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func reconcileGroupPrivilege(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
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
