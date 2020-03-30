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

package y3_9_202002131522

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addCloudHostAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	objID := common.BKInnerObjIDHost
	dataRows := []*Attribute{
		{ObjectID: objID, PropertyID: "bk_cloud_inst_id", PropertyName: "云主机实例ID", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		{ObjectID: objID, PropertyID: "bk_cloud_host_status", PropertyName: "云主机状态", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}

	t := new(time.Time)
	for _, r := range dataRows {
		r.OwnerID = conf.OwnerID
		r.IsPre = true
		if false != r.IsEditable {
			r.IsEditable = true
		}
		r.IsReadOnly = false
		r.CreateTime = t
		r.LastTime = t
		r.Creator = common.CCSystemOperatorUserName
		r.LastEditor = common.CCSystemOperatorUserName
		r.Description = ""

		id, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
		if err != nil {
			blog.ErrorJSON("NextSequence failed, host attrName: %s, err: %v", r.PropertyName, err)
			return err
		}
		r.ID = int64(id)

		if err := db.Table(common.BKTableNameObjAttDes).Insert(ctx, r); err != nil {
			blog.ErrorJSON("insert failed, host attrName: %s, err: %s", r.PropertyName, err)
			return err
		}
	}

	return nil
}
