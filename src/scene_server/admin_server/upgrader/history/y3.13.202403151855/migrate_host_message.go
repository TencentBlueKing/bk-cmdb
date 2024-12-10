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

package y3_13_202403151855

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func migrateHostMsg(ctx context.Context, db dal.RDB, conf *history.Config) error {
	// change host auto group name
	cond := mapstr.MapStr{
		common.BKObjIDField:           common.BKInnerObjIDHost,
		common.BKPropertyGroupIDField: "auto",
	}
	updateData := mapstr.MapStr{common.BKPropertyGroupNameField: "主机系统配置"}

	if err := db.Table(common.BKTableNamePropertyGroup).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("change host auto group name failed, err: %v, cond: %+v, data: %+v", err, cond, updateData)
		return err
	}

	// change host outer mac field name
	cond = mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: common.BKOuterMac,
	}
	updateData = mapstr.MapStr{common.BKPropertyNameField: "外网MAC地址"}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("change host outer mac field name failed, err: %v, cond: %+v, data: %+v", err, cond, updateData)
		return err
	}

	// change the placeholder of the host field
	fields := []string{common.BKHostNameField, common.BKOSTypeField, common.BKOSNameField, common.BKOSVersionField,
		common.BKOSBitField, common.BKCpuField, common.BKCpuModuleField, common.BKMemField, common.BKDiskField,
		common.BKMac, common.BKOuterMac}
	cond = mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: mapstr.MapStr{common.BKDBIN: fields},
		common.BKDBOR: []mapstr.MapStr{
			{metadata.AttributeFieldPlaceHolder: ""},
			{metadata.AttributeFieldPlaceHolder: mapstr.MapStr{common.BKDBExists: false}},
		},
	}

	attrs := make([]metadata.Attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(cond).Fields(common.BKPropertyIDField).All(ctx, &attrs)
	if err != nil {
		blog.Errorf("find host fields failed, err: %v, cond: %+v", err, cond)
		return err
	}

	if len(attrs) == 0 {
		return nil
	}

	changedFields := make([]string, 0)
	for _, attr := range attrs {
		changedFields = append(changedFields, attr.PropertyID)
	}

	cond = mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: mapstr.MapStr{common.BKDBIN: changedFields},
	}
	updateData = mapstr.MapStr{metadata.AttributeFieldPlaceHolder: "可人工设置。若需自动发现，请在主机上安装GseAgent和采集器。"}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("change host fields failed, err: %v, cond: %+v, data: %+v", err, cond, updateData)
		return err
	}

	return nil
}
