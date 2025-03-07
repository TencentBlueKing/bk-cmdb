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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/service/utils"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
)

var objClassificationDataData = []Classification{
	{
		ClassificationID:   "bk_host_manage",
		ClassificationName: "主机管理",
		ClassificationIcon: "icon-cc-host",
	},
	{
		ClassificationID:   "bk_biz_topo",
		ClassificationName: "业务拓扑",
		ClassificationIcon: "icon-cc-business",
	},
	{
		ClassificationID:   "bk_organization",
		ClassificationName: "组织架构",
		ClassificationIcon: "icon-cc-organization",
	},
	{
		ClassificationIcon: "icon-cc-default",
		ClassificationType: "inner",
		ClassificationID:   "bk_uncategorized",
		ClassificationName: "未分类",
	},
	{
		ClassificationType: "hidden",
		ClassificationID:   "bk_table_classification",
		ClassificationName: "表格分类",
	},
}

func addObjClassificationData(kit *rest.Kit, db local.DB) error {
	objClassification := make([]mapstr.MapStr, 0)
	for _, asst := range objClassificationDataData {
		item, err := util.ConvStructToMap(asst)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v", err)
			return err
		}
		objClassification = append(objClassification, item)
	}

	needField := &utils.InsertOptions{
		UniqueFields: []string{"bk_classification_name"},
		IgnoreKeys:   []string{"id"},
		IDField:      []string{metadata.ClassificationFieldID},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_classification_name",
		},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelClassificationRes,
		},
	}

	_, err := utils.InsertData(kit, db, common.BKTableNameObjClassification, objClassification,
		needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjClassification, err)
		return err
	}

	idOptions := &tools.IDOptions{IDField: "id", RemoveKeys: []string{"id"}}
	err = tools.InsertTemplateData(kit, db, objClassification, "obj_classification",
		[]string{"data.bk_classification_name"}, idOptions)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}

	return nil
}

// Classification the classification metadata definition
type Classification struct {
	ID                 int64  `bson:"id"`
	ClassificationID   string `bson:"bk_classification_id"`
	ClassificationName string `bson:"bk_classification_name"`
	ClassificationType string `bson:"bk_classification_type"`
	ClassificationIcon string `bson:"bk_classification_icon"`
}
