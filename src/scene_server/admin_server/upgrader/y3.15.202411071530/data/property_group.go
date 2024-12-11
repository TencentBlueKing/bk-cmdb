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

package data

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal"
)

func addPropertyGroupData(kit *rest.Kit, db dal.Dal) error {
	var propertyGroupArr []interface{}
	var propertyGroupAudit []tools.AuditType
	for _, group := range propertyGroupData {
		group.IsDefault = true
		propertyGroupArr = append(propertyGroupArr, group)
		propertyGroupAudit = append(propertyGroupAudit, tools.AuditType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelGroupRes,
		})
	}

	cmpField := &tools.CmpFiled{
		UniqueFields: []string{common.BKObjIDField, common.BKAppIDField, common.BKPropertyGroupIndexField},
		IgnoreKeys:   []string{common.BKFieldID, common.BKPropertyGroupIndexField},
		IDField:      common.BKFieldID,
	}
	auditDataField := &tools.AuditDataField{
		BusinessID:   "bk_biz_id",
		ResourceID:   common.BKFieldID,
		ResourceName: "bk_group_name",
	}
	_, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNamePropertyGroup,
		propertyGroupArr, cmpField, propertyGroupAudit, auditDataField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameBaseBizSet, err)
		return err
	}

	return nil
}

var propertyGroupData = []metadata.Group{
	{
		ObjectID:   common.BKInnerObjIDApp,
		GroupID:    mCommon.BaseInfo,
		GroupName:  mCommon.BaseInfoName,
		GroupIndex: 1,
	},
	{
		ObjectID:   common.BKInnerObjIDApp,
		GroupID:    mCommon.AppRole,
		GroupName:  mCommon.AppRoleName,
		GroupIndex: 2,
	},
	{
		ObjectID:   common.BKInnerObjIDSet,
		GroupID:    mCommon.BaseInfo,
		GroupName:  mCommon.BaseInfoName,
		GroupIndex: 1,
	},
	{
		ObjectID:   common.BKInnerObjIDModule,
		GroupID:    mCommon.BaseInfo,
		GroupName:  mCommon.BaseInfoName,
		GroupIndex: 1,
	},
	{
		ObjectID:   common.BKInnerObjIDHost,
		GroupID:    mCommon.BaseInfo,
		GroupName:  mCommon.BaseInfoName,
		GroupIndex: 1,
	},
	{
		ObjectID:   common.BKInnerObjIDHost,
		GroupID:    mCommon.HostAutoFields,
		GroupName:  "主机系统配置",
		GroupIndex: 2,
	},
	{
		ObjectID:   common.BKInnerObjIDProc,
		GroupID:    mCommon.BaseInfo,
		GroupName:  mCommon.BaseInfoName,
		GroupIndex: 1,
	},
	{
		ObjectID:   common.BKInnerObjIDProc,
		GroupID:    mCommon.ProcPort,
		GroupName:  mCommon.ProcPortName,
		GroupIndex: 2,
	},
	{
		ObjectID:   common.BKInnerObjIDProc,
		GroupID:    mCommon.ProcGsekitBaseInfo,
		GroupName:  mCommon.ProcGsekitBaseInfoName,
		GroupIndex: 3,
	},
	{
		ObjectID:   common.BKInnerObjIDProc,
		GroupID:    mCommon.ProcGsekitManageInfo,
		GroupName:  mCommon.ProcGsekitManageInfoName,
		GroupIndex: 4,
	},
	{
		ObjectID:   common.BKInnerObjIDPlat,
		GroupID:    mCommon.BaseInfo,
		GroupName:  mCommon.BaseInfoName,
		GroupIndex: 1,
	},
	{
		ObjectID:   common.BKInnerObjIDProc,
		GroupID:    "network_proxy",
		GroupName:  "外网代理信息",
		IsPre:      true,
		IsCollapse: true,
		GroupIndex: 5,
	},
	{
		ObjectID:   common.BKInnerObjIDProc,
		GroupID:    "proc_mgr",
		GroupName:  "进程管理信息",
		IsCollapse: true,
		GroupIndex: 6,
	},
	{
		ObjectID:   common.BKInnerObjIDBizSet,
		GroupID:    "default",
		GroupName:  "基础信息",
		GroupIndex: 1,
	},
	{
		ObjectID:   common.BKInnerObjIDBizSet,
		GroupID:    "default",
		GroupName:  "角色",
		GroupIndex: 2,
	},
	{
		ObjectID:   common.BKInnerObjIDProject,
		GroupID:    "default",
		GroupName:  "基础信息",
		GroupIndex: 1,
	},
}
