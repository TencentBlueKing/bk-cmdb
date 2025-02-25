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

var objectData = []Object{
	{
		ObjCls:        "bk_host_manage",
		ObjectID:      common.BKInnerObjIDHost,
		ObjectName:    "主机",
		IsPre:         true,
		ObjIcon:       "icon-cc-host",
		Position:      `{"bk_host_manage":{"x":-600,"y":-650}}`,
		Creator:       common.CCSystemOperatorUserName,
		ObjSortNumber: 1,
	},
	{
		ObjCls:        "bk_biz_topo",
		ObjectID:      common.BKInnerObjIDModule,
		ObjectName:    "模块",
		IsPre:         true,
		ObjIcon:       "icon-cc-module",
		Creator:       common.CCSystemOperatorUserName,
		ObjSortNumber: 1,
	},
	{
		ObjCls:        "bk_biz_topo",
		ObjectID:      common.BKInnerObjIDSet,
		ObjectName:    "集群",
		IsPre:         true,
		ObjIcon:       "icon-cc-set",
		Creator:       common.CCSystemOperatorUserName,
		ObjSortNumber: 2,
	},
	{
		ObjCls:        "bk_organization",
		ObjectID:      common.BKInnerObjIDApp,
		ObjectName:    "业务",
		IsPre:         true,
		ObjIcon:       "icon-cc-business",
		Position:      `{"bk_organization":{"x":-100,"y":-100}}`,
		Creator:       common.CCSystemOperatorUserName,
		IsPaused:      false,
		ObjSortNumber: 1,
	},
	{
		ObjCls:        "bk_host_manage",
		ObjectID:      common.BKInnerObjIDProc,
		ObjectName:    "进程",
		IsPre:         true,
		ObjIcon:       "icon-cc-process",
		Position:      `{"bk_host_manage":{"x":-450,"y":-650}}`,
		IsHidden:      true,
		Creator:       common.CCSystemOperatorUserName,
		ObjSortNumber: 2,
	},
	{
		ObjCls:        "bk_host_manage",
		ObjectID:      common.BKInnerObjIDPlat,
		ObjectName:    "云区域",
		IsPre:         true,
		ObjIcon:       "icon-cc-subnet",
		Position:      `{"bk_host_manage":{"x":-600,"y":-500}}`,
		IsHidden:      true,
		Creator:       common.CCSystemOperatorUserName,
		ObjSortNumber: 3,
	},
	{
		ObjCls:        "bk_organization",
		ObjectID:      common.BKInnerObjIDBizSet,
		ObjectName:    "业务集",
		ObjIcon:       "icon-cc-business-set",
		IsPre:         true,
		Creator:       common.CCSystemOperatorUserName,
		ObjSortNumber: 2,
	},
	{
		ObjCls:        "bk_organization",
		ObjectID:      common.BKInnerObjIDProject,
		ObjectName:    "项目",
		ObjIcon:       "icon-cc-project",
		IsPre:         true,
		Creator:       common.CCSystemOperatorUserName,
		ObjSortNumber: 3,
	},
}

// Object object metadata definition
type Object struct {
	ID            int64       `bson:"id"`
	ObjCls        string      `bson:"bk_classification_id"`
	ObjIcon       string      `bson:"bk_obj_icon"`
	ObjectID      string      `bson:"bk_obj_id"`
	ObjectName    string      `bson:"bk_obj_name"`
	IsHidden      bool        `bson:"bk_ishidden"`
	IsPre         bool        `bson:"ispre"`
	IsPaused      bool        `bson:"bk_ispaused"`
	Position      string      `bson:"position"`
	Description   string      `bson:"description"`
	Creator       string      `bson:"creator"`
	Modifier      string      `bson:"modifier"`
	Time          *tools.Time `bson:",inline"`
	ObjSortNumber int64       `bson:"obj_sort_number"`
}

func addObjectData(kit *rest.Kit, db local.DB) error {
	objectDataArr := make([]mapstr.MapStr, 0)
	for _, obj := range objectData {
		obj.Time = tools.NewTime()
		item, err := util.ConvStructToMap(obj)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v", err)
			return err
		}
		objectDataArr = append(objectDataArr, item)
	}

	needField := &utils.InsertOptions{
		UniqueFields: []string{"bk_obj_id"},
		IgnoreKeys:   []string{"id", "obj_sort_number"},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModuleRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_name",
		},
	}

	_, err := utils.InsertData(kit, db, common.BKTableNameObjDes, objectDataArr, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjDes, err)
		return err
	}

	idOptions := &tools.IDOptions{IDField: "id", RemoveKeys: []string{"id"}}
	err = tools.InsertTemplateData(kit, db, objectDataArr, "object", []string{"data.bk_obj_id"}, idOptions)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}
	return nil
}
