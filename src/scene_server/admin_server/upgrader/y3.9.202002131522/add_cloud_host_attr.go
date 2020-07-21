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
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addCloudHostAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	objID := common.BKInnerObjIDHost
	dataRows := []*Attribute{
		{ObjectID: objID, PropertyID: "bk_cloud_inst_id", PropertyName: "云主机实例ID", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		{ObjectID: objID, PropertyID: "bk_cloud_host_status", PropertyName: "云主机状态", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: cloudInstStatusEnum},
		{ObjectID: objID, PropertyID: "bk_cloud_vendor", PropertyName: "云厂商", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: cloudVendorEnum},
	}

	now := time.Now()
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	for _, r := range dataRows {
		r.OwnerID = conf.OwnerID
		r.IsPre = true
		r.IsReadOnly = false
		r.CreateTime = &now
		r.LastTime = &now
		r.Creator = common.CCSystemOperatorUserName
		r.LastEditor = common.CCSystemOperatorUserName
		r.Description = ""

		if _, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, r, "id", uniqueFields, []string{}); err != nil {
			blog.ErrorJSON("addCloudHostAttr failed, Upsert err: %s, attribute: %#v, ", err, r)
			return err
		}
	}

	return nil
}

var cloudInstStatusEnum = []metadata.EnumVal{
	{ID: "1", Name: "未知", Type: "text"},
	{ID: "2", Name: "启动中", Type: "text"},
	{ID: "3", Name: "运行中", Type: "text"},
	{ID: "4", Name: "停止中", Type: "text"},
	{ID: "5", Name: "已停止", Type: "text"},
	{ID: "6", Name: "已销毁", Type: "text"},
}