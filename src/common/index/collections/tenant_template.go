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

package collections

import (
	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	registerIndexes(common.BKTableNameTenantTemplate, commTenantTemplateIndexes)
}

var commTenantTemplateIndexes = []types.Index{
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys:                    bson.D{{common.BKFieldID, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypeObjAttribute) +
			"_bkObjID_bkPropertyID",
		Keys:       bson.D{{"data.bk_obj_id", 1}, {"data.bk_property_id", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypeObjAttribute,
		},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypePropertyGroup) +
			"_bkGroupName_bkObjID",
		Keys:       bson.D{{"data.bk_group_name", 1}, {"data.bk_obj_id", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypePropertyGroup,
		},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypePropertyGroup) +
			"_bkObjID_bkGroupIndex",
		Keys:       bson.D{{"data.bk_obj_id", 1}, {"data.bk_group_index", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypePropertyGroup,
		},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypePlat) + "_bkCloudName",
		Keys:       bson.D{{"data.bk_cloud_name", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypePlat,
			"data.bk_cloud_name":             map[string]string{common.BKDBType: "string"}},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypeObjClassification) +
			"_bkClassificationID",
		Keys:       bson.D{{"data.bk_classification_id", 1}},
		Unique:     true,
		Background: false,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypeObjClassification,
			"data.bk_classification_id":      map[string]string{common.BKDBType: "string"}},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypeObjClassification) +
			"_bkClassificationName",
		Keys:       bson.D{{"data.bk_classification_name", 1}},
		Unique:     true,
		Background: false,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypeObjClassification,
			"data.bk_classification_name":    map[string]string{common.BKDBType: "string"}},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypeObject) + "_bkObjID",
		Keys:       bson.D{{"data.bk_obj_id", 1}},
		Unique:     true,
		Background: false,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypeObject,
			"data.bk_obj_id":                 map[string]string{common.BKDBType: "string"}},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypeObject) + "_bkObjName",
		Keys:       bson.D{{"data.bk_obj_name", 1}},
		Unique:     true,
		Background: false,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypeObject,
			"data.bk_obj_name":               map[string]string{common.BKDBType: "string"}},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypeBizSet) + "_bkBizSetName",
		Keys:       bson.D{{"data.bk_biz_set_name", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypeBizSet,
		},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenant.TemplateTypeServiceCategory) + "_nameParentName",
		Keys:       bson.D{{"data.name", 1}, {"data.parent_name", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKTenantTemplateTypeField: tenant.TemplateTypeServiceCategory,
		},
	},
}
