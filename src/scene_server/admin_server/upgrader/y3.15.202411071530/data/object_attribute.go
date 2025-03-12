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
	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
)

var objAttrMap = map[string][]*attribute{
	"biz":            bizObjAttrs,
	"host":           hostObjAttrs,
	"set":            setObjAttrs,
	"process":        processObjAttrs,
	"module":         moduleObjAttrs,
	"plat":           platObjAttrs,
	"bk_biz_set_obj": bizSetObjAttrData,
	"bk_project":     projectObjAttrs,
}

var objPropertyMap = map[string]string{
	"set":            "default",
	"module":         "default",
	"plat":           "default",
	"bk_biz_set_obj": "default",
	"bk_project":     "default",
}

func getAttrData() []*attribute {
	for key, value := range objAttrMap {
		for _, attr := range value {
			attr.ObjectID = key
			attr.Time = tools.NewTime()
			attr.IsPre = true
			attr.Creator = "cc_system"
			if propertyGroup, ok := objPropertyMap[key]; ok {
				attr.PropertyGroup = propertyGroup
			}
		}
		objAttrData = append(objAttrData, value...)
	}
	return objAttrData
}

func addObjAttrData(kit *rest.Kit, db local.DB) error {
	if len(objAttrData) == 0 {
		getAttrData()
	}

	indexMap := make(map[string]int64)
	attributeData := make([]mapstr.MapStr, 0)
	for _, attr := range objAttrData {
		if _, ok := indexMap[attr.ObjectID+attr.PropertyGroup]; !ok {
			indexMap[attr.ObjectID+attr.PropertyGroup] = 1
		} else {
			indexMap[attr.ObjectID+attr.PropertyGroup] += 1
		}
		attr.PropertyIndex = indexMap[attr.ObjectID+attr.PropertyGroup]
		item, err := tools.ConvStructToMap(attr)
		if err != nil {
			blog.Errorf("convert attribute to mapstr failed, err: %v", err)
			return err
		}
		attributeData = append(attributeData, item)
	}

	needField := &tools.InsertOptions{
		UniqueFields: []string{common.BKObjIDField, common.BKPropertyIDField, common.BKAppIDField},
		IgnoreKeys:   []string{"id", "bk_property_index"},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &tools.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelAttributeRes,
		},
		AuditDataField: &tools.AuditDataField{
			BizIDField:   "bk_biz_id",
			ResIDField:   common.BKFieldID,
			ResNameField: "bk_property_name",
		},
	}

	_, err := tools.InsertData(kit, db, common.BKTableNameObjAttDes, attributeData, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjAttDes, err)
		return err
	}

	idOptions := &tools.IDOptions{ResNameField: "bk_property_name", RemoveKeys: []string{"id"}}
	err = tools.InsertTemplateData(kit, db, attributeData, needField, idOptions, tenanttmp.TemplateTypeObjAttribute)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}

	return nil
}
