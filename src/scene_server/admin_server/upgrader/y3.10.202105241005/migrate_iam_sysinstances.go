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

package y3_10_202105241005

import (
	"context"

	"configcenter/src/ac"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)
// migrateIAMSysInstances migrate iam system instances
func migrateIAMSysInstances(ctx context.Context, db dal.RDB, iam ac.AuthInterface, conf *upgrader.Config) error {
	// get all custom objects
	objects := []metadata.Object{}
	condition := map[string]interface{}{
		common.BKIsPre: false,
	}
	err := db.Table(common.BKTableNameObjDes).Find(condition).All(ctx, &objects)
	if err != nil {
		blog.ErrorJSON("get all custom objects failed, err:%s", err)
		return err
	}

	return iam.SyncIAMSysInstances(ctx, objects)
}
