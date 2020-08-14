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

package y3_8_202008111026

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func removeBkSupplierIDField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKPropertyIDField: "bk_supplier_id",
		common.BKObjIDField:      common.BKInnerObjIDApp,
	}

	count, err := db.Table(common.BKTableNameObjAttDes).Find(filter).Count(ctx)
	if err != nil {
		blog.Errorf("count bk_supplier_id attribute failed, err: %s", err.Error())
		return err
	}

	if count > 0 {
		if err := db.Table(common.BKTableNameObjAttDes).Delete(ctx, filter); err != nil {
			blog.Errorf("delete bk_supplier_id attribute failed, err: %s", err.Error())
			return err
		}
	}

	if err := db.Table(common.BKTableNameBaseApp).DropColumn(ctx, "bk_supplier_id"); err != nil {
		blog.Errorf("remove biz bk_supplier_id field failed, err: %s", err.Error())
		return err
	}
	return nil
}
