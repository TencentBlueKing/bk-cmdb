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

package tenant_template

// TenantTemplateType tenant template type
type TenantTemplateType string

// TenantTmpData tenant template data struct
type TenantTmpData[T any] struct {
	Type  TenantTemplateType `bson:"type"`
	IsPre bool               `bson:"is_pre"`
	ID    int64              `bson:"id"`
	Data  T                  `bson:"data"`
}

// SvrCategoryTmp tenant template for service category
type SvrCategoryTmp struct {
	Name       string `bson:"name"`
	ParentName string `bson:"parent_name"`
}

// UniqueKeyTmp tenant template for unique keys
type UniqueKeyTmp struct {
	ObjectID string   `bson:"bk_obj_id"`
	Keys     []string `bson:"keys"`
}

const (
	// BKTenantTemplateTypeField is type field for tenant template
	BKTenantTemplateTypeField = "type"

	// BKTenantTemplateIDField is id field for tenant template
	BKTenantTemplateIDField = "id"

	// BKTenantIsPreField is isPre field for tenant template
	BKTenantIsPreField = "is_pre"
)

// tenant template type
const (
	TemplateTypeAssociation       TenantTemplateType = "association"
	TemplateTypeObject            TenantTemplateType = "object"
	TemplateTypeObjAttribute      TenantTemplateType = "obj_attribute"
	TemplateTypeObjAssociation    TenantTemplateType = "obj_association"
	TemplateTypeObjClassification TenantTemplateType = "obj_classification"
	TemplateTypePlat              TenantTemplateType = "plat"
	TemplateTypePropertyGroup     TenantTemplateType = "property_group"
	TemplateTypeServiceCategory   TenantTemplateType = "service_category"
	TemplateTypeUniqueKeys        TenantTemplateType = "unique_keys"
	TemplateTypeBizSet            TenantTemplateType = "biz_set"
)

// AllTemplateTypes variables for all template types
var AllTemplateTypes = []TenantTemplateType{
	TemplateTypeAssociation,
	TemplateTypeObject,
	TemplateTypeObjAttribute,
	TemplateTypeObjAssociation,
	TemplateTypeObjClassification,
	TemplateTypePlat,
	TemplateTypePropertyGroup,
	TemplateTypeBizSet,
	TemplateTypeServiceCategory,
	TemplateTypeUniqueKeys,
}

// GetTemplateTypes get all tenant template types
func GetTemplateTypes() []TenantTemplateType {
	return AllTemplateTypes
}
