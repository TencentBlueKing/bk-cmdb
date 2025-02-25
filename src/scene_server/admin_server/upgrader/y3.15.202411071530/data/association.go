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

var asstTypes = []associationKind{
	{
		AssociationKindID:       "belong",
		AssociationKindName:     "属于",
		SourceToDestinationNote: "属于",
		DestinationToSourceNote: "包含",
	},
	{
		AssociationKindID:       "group",
		AssociationKindName:     "组成",
		SourceToDestinationNote: "组成",
		DestinationToSourceNote: "组成于",
	},
	{
		AssociationKindID:       "bk_mainline",
		AssociationKindName:     "拓扑组成",
		SourceToDestinationNote: "组成",
		DestinationToSourceNote: "组成于",
	},
	{
		AssociationKindID:       "run",
		AssociationKindName:     "运行",
		SourceToDestinationNote: "运行于",
		DestinationToSourceNote: "运行",
	},
	{
		AssociationKindID:       "connect",
		AssociationKindName:     "上联",
		SourceToDestinationNote: "上联",
		DestinationToSourceNote: "下联",
	},
	{
		AssociationKindID:       "default",
		AssociationKindName:     "默认关联",
		SourceToDestinationNote: "关联",
		DestinationToSourceNote: "被关联",
	},
}

func addAssociationData(kit *rest.Kit, db local.DB) error {
	dataInterface := make([]mapstr.MapStr, 0)
	for _, asstType := range asstTypes {
		asstType.IsPre = &trueVar
		asstType.Direction = metadata.DestinationToSource
		item, err := util.ConvStructToMap(asstType)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v", err)
			return err
		}
		dataInterface = append(dataInterface, item)
	}

	needFields := &utils.InsertOptions{
		UniqueFields: []string{common.AssociationKindIDField},
		IgnoreKeys:   []string{"id"},
		IDField:      []string{metadata.AttributeFieldID},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelAssociationRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_asst_name",
		},
	}

	_, err := utils.InsertData(kit, db, common.BKTableNameAsstDes, dataInterface, needFields)
	if err != nil {
		blog.Errorf("insert association data for table %s failed, err: %v", common.BKTableNameAsstDes, err)
		return err
	}

	idOption := &tools.IDOptions{IDField: "id", RemoveKeys: []string{"id"}}

	err = tools.InsertTemplateData(kit, db, dataInterface, "association", []string{"data.bk_asst_id"}, idOption)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}

	return nil
}

type associationKind struct {
	ID                      int64                         `bson:"id"`
	AssociationKindID       string                        `bson:"bk_asst_id"`
	AssociationKindName     string                        `bson:"bk_asst_name"`
	SourceToDestinationNote string                        `bson:"src_des"`
	DestinationToSourceNote string                        `bson:"dest_des"`
	Direction               metadata.AssociationDirection `bson:"direction"`
	IsPre                   *bool                         `bson:"ispre"`
}
