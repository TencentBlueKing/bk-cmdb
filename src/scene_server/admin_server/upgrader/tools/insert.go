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

package tools

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/service/utils"
	"configcenter/src/storage/dal/mongo/local"
)

func InsertTemplateData(kit *rest.Kit, db local.DB, data []mapstr.MapStr, dataType string, uniqueField []string,
	idOption *IDOptions) error {

	dataMap := make([]mapstr.MapStr, 0)
	for _, item := range data {
		for _, key := range idOption.RemoveKeys {
			delete(item, key)
		}

		tmp := metadata.TemplateData{
			Type:  dataType,
			IsPre: true,
			Data:  item,
		}

		tmp.Data = item
		result, err := util.ConvStructToMap(tmp)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		dataMap = append(dataMap, result)
	}

	uniqueField = append(uniqueField, "type")
	needFields := &utils.InsertOptions{
		UniqueFields: uniqueField,
		IgnoreKeys:   []string{idOption.IDField},
		IDField:      []string{idOption.IDField},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.TenantTemplate,
			ResourceType: metadata.TenantTemplateRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   idOption.IDField,
			ResNameField: "type",
		},
	}

	_, err := utils.InsertData(kit, db, common.BKTableNameTenantTemplate, dataMap, needFields)
	if err != nil {
		blog.Errorf("insert %s data for table %s failed, err: %v", dataType, common.BKTableNameAsstDes, err)
		return err
	}

	return nil
}

func InsertSvrTmp(kit *rest.Kit, db local.DB, data []mapstr.MapStr, isParent bool, parentName []string) error {

	dataMap := make([]mapstr.MapStr, 0)
	uniqueKeys := []string{"data.name", "data.bk_biz_id", "is_parent"}
	RemoveKeys := []string{"id", "bk_root_id", "bk_parent_id"}
	for i, item := range data {
		for _, key := range RemoveKeys {
			delete(item, key)
		}
		tmp := metadata.SvrCategoryTmp{
			IsParent:   isParent,
			Name:       item["name"].(string),
			ParentName: parentName[i],
			Type:       "service_category",
			IsPre:      true,
			Data:       item,
		}

		tmp.Data = item
		result, err := util.ConvStructToMap(tmp)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		dataMap = append(dataMap, result)
	}

	uniqueField := append(uniqueKeys, "type")
	needFields := &utils.InsertOptions{
		UniqueFields: uniqueField,
		IgnoreKeys:   []string{"id"},
		IDField:      []string{"id"},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.TenantTemplate,
			ResourceType: metadata.TenantTemplateRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "type",
		},
	}

	_, err := utils.InsertData(kit, db, common.BKTableNameTenantTemplate, dataMap, needFields)
	if err != nil {
		blog.Errorf("insert service_category template data for table %s failed, err: %v", common.BKTableNameAsstDes,
			err)
		return err
	}

	return nil
}

func InsertUniqueKeyTmp(kit *rest.Kit, db local.DB, data []mapstr.MapStr, attributes [][]string) error {

	dataMap := make([]mapstr.MapStr, 0)
	uniqueKeys := []string{"bk_obj_id", "attributes"}
	removeKeys := []string{"id"}
	for i, item := range data {
		for _, key := range removeKeys {
			delete(item, key)
		}
		tmp := metadata.UniqueKeyTmp{
			Type:       metadata.TemplateTypeUniqueKeys,
			Attributes: attributes[i],
			ObjectID:   item["bk_obj_id"].(string),
			IsPre:      true,
			Data:       item,
		}

		tmp.Data = item
		result, err := util.ConvStructToMap(tmp)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		dataMap = append(dataMap, result)
	}

	uniqueField := append(uniqueKeys, "type")
	needFields := &utils.InsertOptions{
		UniqueFields: uniqueField,
		IgnoreKeys:   []string{"id"},
		IDField:      []string{"id"},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.TenantTemplate,
			ResourceType: metadata.TenantTemplateRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "type",
		},
	}

	_, err := utils.InsertData(kit, db, common.BKTableNameTenantTemplate, dataMap, needFields)
	if err != nil {
		blog.Errorf("insert service_category template data for table %s failed, err: %v", common.BKTableNameAsstDes,
			err)
		return err
	}

	return nil
}

// IDOptions the options of data template id
type IDOptions struct {
	IDField    string
	RemoveKeys []string
}
