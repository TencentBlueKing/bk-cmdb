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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
)

var (
	setData = map[string]interface{}{
		common.BKSetNameField:       common.DefaultResSetName,
		common.BKDefaultField:       common.DefaultResSetFlag,
		common.BKSetEnvField:        "3",
		common.BKSetStatusField:     "1",
		common.BKSetTemplateIDField: 0,
	}
)

func addSetBaseData(kit *rest.Kit, db local.DB, bizID int64) (map[string]interface{}, error) {
	setData[common.BKAppIDField] = bizID
	setData[common.BKInstParentStr] = bizID
	setData[common.CreateTimeField] = time.Now()
	setData[common.LastTimeField] = time.Now()
	setData[common.BKSetDescField] = ""
	setData[common.BKDescriptionField] = ""

	needField := &tools.InsertOptions{
		UniqueFields: []string{common.BKAppIDField, common.BKSetNameField, common.BKInstParentStr},
		IgnoreKeys:   []string{common.BKSetIDField},
		IDField:      []string{common.BKSetIDField},
		AuditDataField: &tools.AuditDataField{
			BizIDField:   common.BKAppIDField,
			ResIDField:   common.BKSetIDField,
			ResNameField: common.BKSetNameField,
		},
		AuditTypeField: &tools.AuditResType{
			AuditType:    common.BKInnerObjIDSet,
			ResourceType: metadata.SetRes,
		},
	}

	ids, err := tools.InsertData(kit, db, common.BKTableNameBaseSet, []mapstr.MapStr{setData},
		needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameBaseApp, err)
		return nil, err
	}
	return ids, nil
}
