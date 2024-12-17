/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_13_202404221100

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func updateHostField(ctx context.Context, db dal.RDB, conf *history.Config) error {
	cond := mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: common.BKCloudIDField,
		common.BKDBOR: []mapstr.MapStr{
			{metadata.AttributeFieldPlaceHolder: ""},
			{metadata.AttributeFieldPlaceHolder: mapstr.MapStr{common.BKDBExists: false}},
		},
	}
	msg := "管控区域为“未分配”的主机可在节点管理中重新指定管控区域并安装Agent"
	updateData := mapstr.MapStr{metadata.AttributeFieldPlaceHolder: msg}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("change host cloud area place holder failed, err: %v, cond: %+v, data: %+v", err, cond, updateData)
		return err
	}

	cond = mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: common.BKHostInnerIPField,
	}
	updateData = mapstr.MapStr{common.BKPropertyNameField: "内网IPv4"}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("change host inner ipv4 field name failed, err: %v, cond: %+v, data: %+v", err, cond, updateData)
		return err
	}

	cond = mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: common.BKHostOuterIPField,
	}
	updateData = mapstr.MapStr{common.BKPropertyNameField: "外网IPv4"}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("change host outer ipv4 field name failed, err: %v, cond: %+v, data: %+v", err, cond, updateData)
		return err
	}

	return nil
}
