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

package x19_09_03_03

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

var ProcMgrGroupID = "proc_mgr"

func AddProcAttrGroup(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	doc := map[string]interface{}{
		common.BKIsCollapseField:         true,
		common.BKPropertyGroupNameField:  "进程管理信息",
		common.BKPropertyGroupIndexField: 3,
		"ispre":                          true,
		"bk_isdefault":                   true,
		"metadata": map[string]interface{}{
			"label": make(map[string]interface{}),
		},
		common.BKPropertyGroupIDField: ProcMgrGroupID,
		common.BKObjIDField:           common.BKInnerObjIDProc,
		common.BkSupplierAccount:      conf.OwnerID,
	}
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyGroupIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNamePropertyGroup, doc, "id", uniqueFields, []string{})
	if err != nil {
		if db.IsNotFoundError(err) == false {
			return fmt.Errorf("upgrade x19_09_03_03, AddProcAttrGroup failed, err: %v", err)
		}
	}
	return nil
}

func ChangeProcFieldGroup(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKPropertyIDField: map[string]interface{}{
			common.BKDBIN: []string{"bk_func_id", "work_path", "user", "proc_num", "priority", "timeout",
				"start_cmd", "stop_cmd", "restart_cmd", "face_stop_cmd", "reload_cmd", "pid_file",
				"auto_start", "auto_time_gap"},
		},
	}
	doc := map[string]interface{}{
		"bk_property_group": ProcMgrGroupID,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("ChangeProcFieldGroup failed, err: %+v", err)
		return fmt.Errorf("ChangeProcFieldGroup failed, err: %v", err)
	}
	return nil
}
