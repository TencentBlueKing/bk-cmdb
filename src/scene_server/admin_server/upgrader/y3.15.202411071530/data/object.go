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

	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"

	"github.com/rs/xid"
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
		ObjCls:        "bk_uncategorized",
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
		ObjCls:        "bk_uncategorized",
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
	ID            int64     `bson:"id"`
	ObjCls        string    `bson:"bk_classification_id"`
	ObjIcon       string    `bson:"bk_obj_icon"`
	ObjectID      string    `bson:"bk_obj_id"`
	ObjectName    string    `bson:"bk_obj_name"`
	IsHidden      bool      `bson:"bk_ishidden"`
	IsPre         bool      `bson:"ispre"`
	IsPaused      bool      `bson:"bk_ispaused"`
	Position      string    `bson:"position"`
	Description   string    `bson:"description"`
	Creator       string    `bson:"creator"`
	Modifier      string    `bson:"modifier"`
	CreateTime    time.Time `bson:"create_time"`
	LastTime      time.Time `bson:"last_time"`
	ObjSortNumber int64     `bson:"obj_sort_number"`
	UUID          string    `bson:"uuid"`
}

// AddObjectData add object data
func AddObjectData(kit *rest.Kit, db local.DB) (map[string]string, error) {
	existObjectData := make([]Object, 0)
	err := db.Table(common.BKTableNameObjDes).Find(mapstr.MapStr{}).All(kit.Ctx, &existObjectData)
	if err != nil {
		blog.Errorf("find object data failed, err: %v", err)
		return nil, err
	}
	existObject := make(map[string]Object)
	for _, obj := range existObjectData {
		existObject[obj.ObjectID] = obj
	}

	insertData := make([]map[string]interface{}, 0)
	insertTemplateData := make([]interface{}, 0)
	objUUIDMap := map[string]string{}
	ignoreKeys := []string{common.BKFieldID, common.ObjSortNumberField, metadata.ModelFieldObjUUID}
	for _, obj := range objectData {
		insertTemplateData = append(insertTemplateData, obj)
		objMapData, err := tools.ConvStructToMap(obj)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v", err)
			return nil, err
		}
		if _, ok := existObject[obj.ObjectID]; ok {
			existMapData, err := tools.ConvStructToMap(existObject[obj.ObjectID])
			if err != nil {
				blog.Errorf("convert struct to map failed, err: %v", err)
				return nil, err
			}
			if err = tools.CmpData(objMapData, existMapData, ignoreKeys); err != nil {
				blog.Errorf("compare data failed, err: %v", err)
				return nil, err
			}
			objUUIDMap[obj.ObjectID] = existObject[obj.ObjectID].UUID
			continue
		}
		id, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequence(kit.Ctx, common.BKTableNameObjDes)
		if err != nil {
			blog.Errorf("get next sequence failed, err: %v", err)
			return nil, err
		}
		objMapData[common.BKFieldID] = id
		objMapData[common.CreateTimeField] = time.Now()
		objMapData[common.LastTimeField] = time.Now()

		objUUID := xid.New().String()
		objMapData[metadata.ModelFieldObjUUID] = objUUID
		objUUIDMap[obj.ObjectID] = objUUID
		insertData = append(insertData, objMapData)
	}

	if err = insertObjData(kit, db, insertTemplateData, insertData); err != nil {
		blog.Errorf("insert object related data failed, err: %v", err)
		return nil, err
	}

	return objUUIDMap, nil
}

func insertObjData(kit *rest.Kit, db local.DB, insertTemplateData []interface{},
	insertData []map[string]interface{}) error {

	needField := &tools.InsertOptions{
		UniqueFields: []string{common.BKObjIDField},
		IgnoreKeys:   []string{common.BKFieldID, common.ObjSortNumberField, metadata.ModelFieldObjUUID},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &tools.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModuleRes,
		},
		AuditDataField: &tools.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_name",
		},
	}

	idOptions := &tools.IDOptions{ResNameField: "bk_obj_name", RemoveKeys: []string{"id"}}
	err := tools.InsertTemplateData(kit, db, insertTemplateData, needField, idOptions, tenanttmp.TemplateTypeObject)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}

	if len(insertData) == 0 {
		return nil
	}

	err = db.Table(common.BKTableNameObjDes).Insert(kit.Ctx, insertData)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjDes, err)
		return err
	}

	auditField := &tools.AuditStruct{
		AuditDataField: &tools.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_name",
		},
		AuditTypeData: &tools.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelRes,
		},
	}

	if err = tools.AddCreateAuditLog(kit, db, insertData, auditField); err != nil {
		blog.Errorf("add audit log failed, err: %v", err)
		return err
	}

	return nil
}
