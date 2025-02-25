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
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/service/utils"
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
	"plat":           {{"bk_cloud_name"}},
	"set":            {{"bk_biz_id", "bk_set_name", "bk_parent_id"}},
	"module":         {{"bk_biz_id", "bk_set_id", "bk_module_name"}},
	"bk_project":     {{"bk_project_code"}, {"bk_project_id"}, {"bk_project_name"}, {"id"}},
	"bk_biz_set_obj": {{"bk_biz_set_name"}, {"bk_biz_set_id"}},
}

func getUniqueKeys(kit *rest.Kit, db local.DB) ([]objectUnique, [][]string, error) {
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
	var attributes [][]string
	for objID, value := range objUniqueKeys {
		tempValue := objID
		for _, property := range value {
			keys := make([]uniqueKey, 0)
			for _, field := range property {
				keys = append(keys, uniqueKey{
					Kind: "property",
					ID:   attrIDMap[generateUniqueKey(objID, field)],
				})
				tempValue += "-" + field
			}
			attributes = append(attributes, property)
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
		item, err := util.ConvStructToMap(key)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v", err)
			continue
		}
		objUniqueData = append(objUniqueData, item)
	}

	needField := &utils.InsertOptions{
		UniqueFields: []string{"keys"},
		IgnoreKeys:   []string{common.BKFieldID},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelUniqueRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_id",
		},
	}

	_, err = utils.InsertData(kit, db, common.BKTableNameObjUnique, objUniqueData, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjUnique, err)
		return err
	}
	// add tenant template data
	err = tools.InsertUniqueKeyTmp(kit, db, objUniqueData, attributes)
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
