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
)

func generateUniqueKey(objID, propertyID string) string {
	return objID + ":" + propertyID
}

var objUniqueKeys = map[string][][]string{
	"host": {
		{"bk_cloud_id", "bk_host_innerip"},
		{"bk_cloud_inst_id", "bk_cloud_vendor"},
		{"bk_host_outerip"},
		{"bk_host_innerip_v6", "bk_cloud_id"},
		{"bk_agent_id"},
	},
	"biz":            {{"bk_biz_name"}},
	"plat":           {{"bk_cloud_name"}, {"bk_vpc_id"}},
	"set":            {{"bk_biz_id", "bk_set_name", "bk_parent_id"}},
	"module":         {{"bk_biz_id", "bk_set_id", "bk_module_name"}},
	"bk_project":     {{"bk_project_code"}, {"bk_project_id"}, {"bk_project_name"}, {"id"}},
	"bk_biz_set_obj": {{"bk_biz_set_name"}, {"bk_biz_set_id"}},
}

func getUniqueKeys(kit *rest.Kit, db local.DB) ([]objectUnique, []tenanttmp.UniqueKeyTmp, error) {
	attrArr := make([]metadata.Attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(nil).All(kit.Ctx, &attrArr)
	if err != nil {
		blog.Errorf("get host unique fields failed, err: %v", err)
		return nil, nil, err
	}

	attrIDMap := make(map[string]uint64)
	for _, attr := range attrArr {
		attrIDMap[generateUniqueKey(attr.ObjectID, attr.PropertyID)] = uint64(attr.ID)
	}
	uniqueKeys := make([]objectUnique, 0)
	var attributes []tenanttmp.UniqueKeyTmp
	for objID, value := range objUniqueKeys {
		for _, property := range value {
			keys := make([]uniqueKey, 0)
			for _, field := range property {
				keys = append(keys, uniqueKey{
					Kind: "property",
					ID:   attrIDMap[generateUniqueKey(objID, field)],
				})
			}
			attributes = append(attributes, tenanttmp.UniqueKeyTmp{
				ObjectID: objID,
				Keys:     property,
			})
			uniqueKeys = append(uniqueKeys, objectUnique{
				Keys:     keys,
				ObjID:    objID,
				IsPre:    true,
				LastTime: time.Now(),
			})
		}
	}

	return uniqueKeys, attributes, nil
}

func addObjectUniqueData(kit *rest.Kit, db local.DB) error {

	uniqueKeysArr, attributes, err := getUniqueKeys(kit, db)
	if err != nil {
		blog.Errorf("get unique keys failed, err: %v", err)
		return err
	}

	objUniqueData := make([]mapstr.MapStr, 0)
	for _, key := range uniqueKeysArr {
		item, err := tools.ConvStructToMap(key)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v", err)
			continue
		}
		objUniqueData = append(objUniqueData, item)
	}

	needField := &tools.InsertOptions{
		UniqueFields: []string{"keys"},
		IgnoreKeys:   []string{common.BKFieldID, common.BKFieldDBID},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &tools.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelUniqueRes,
		},
		AuditDataField: &tools.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_id",
		},
	}

	_, err = tools.InsertData(kit, db, common.BKTableNameObjUnique, objUniqueData, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjUnique, err)
		return err
	}
	// add tenant template data
	uniqueTmpData := make([]tenanttmp.TenantTmpData[tenanttmp.UniqueKeyTmp], 0)
	for _, data := range attributes {
		uniqueTmpData = append(uniqueTmpData, tenanttmp.TenantTmpData[tenanttmp.UniqueKeyTmp]{
			Type:  tenanttmp.TemplateTypeUniqueKeys,
			IsPre: true,
			Data:  data,
		})
	}
	err = tools.InsertUniqueKeyTmp(kit, db, uniqueTmpData)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}
	return nil
}

type objectUnique struct {
	ID       uint64      `bson:"id"`
	ObjID    string      `bson:"bk_obj_id"`
	Keys     []uniqueKey `bson:"keys"`
	IsPre    bool        `bson:"ispre"`
	LastTime time.Time   `bson:"last_time"`
}

type uniqueKey struct {
	Kind string `bson:"key_kind"`
	ID   uint64 `bson:"key_id"`
}
