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

package y3_10_202106101505

import (
	"context"

	"configcenter/src/ac"
	iamtype "configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// migrateIAMSysInstances migrate iam system instances
func migrateIAMSysInstances(ctx context.Context, db dal.RDB, iam ac.AuthInterface, conf *upgrader.Config) error {
	// get all custom objects(without mainline objects)
	objects := []metadata.Object{}
	condition := map[string]interface{}{
		common.BKIsPre: false,
		common.BKClassificationIDField: map[string]interface{}{
			common.BKDBNE: "bk_biz_topo",
		},
	}
	err := db.Table(common.BKTableNameObjDes).Find(condition).All(ctx, &objects)
	if err != nil {
		blog.ErrorJSON("get all custom objects failed, err:%s", err)
		return err
	}

	param := &iamtype.DeleteCMDBResourceParam{
		ActionIDs: []iamtype.ActionID{
			"create_sys_instance",
			"edit_sys_instance",
			"delete_sys_instance",
			"create_event_subscription",
			"edit_event_subscription",
			"delete_event_subscription",
			"find_event_subscription",
			"watch_set_template_event",
		},
		InstanceSelectionIDs: []iamtype.InstanceSelectionID{"sys_instance", "sys_event_pushing"},
		TypeIDs:              []iamtype.TypeID{"sys_instance", "sys_event_pushing"},
	}

	// add new system instances
	if err := iam.SyncIAMSysInstances(ctx, objects); err != nil {
		blog.ErrorJSON("sync iam system instances failed, err:%s", err)
		return err
	}

	// delete the old cmdb resource
	return iam.DeleteCMDBResource(ctx, param, objects)
}
