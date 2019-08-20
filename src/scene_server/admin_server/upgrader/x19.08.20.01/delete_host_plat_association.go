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
package x19_08_20_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func upgradeDeleteHostPlatAssociation(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {

	cond := map[string]interface{}{
		"bk_obj_id":           common.BKInnerObjIDPlat,
		"bk_supplier_account": common.BKDefaultOwnerID,
		"bk_asst_obj_id":      common.BKInnerObjIDHost,
		"bk_asst_id":          "default",
		"bk_obj_asst_id":      "plat_default_host",
		"bk_obj_asst_name":    "云区域",
		"ispre":               false,
		"mapping":             "1:n",
	}
	if err := db.Table(common.BKTableNameObjAsst).Delete(ctx, cond); err != nil {
		blog.ErrorJSON("upgrade delete host and plat association error. cond:%s, err:%s", cond, err.Error())
		return err
	}

	return nil
}
