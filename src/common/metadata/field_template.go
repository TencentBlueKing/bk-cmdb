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

package metadata

import (
	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
)

// FieldTemplate field template definition
type FieldTemplate struct {
	ID          int64  `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	OwnerID     string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Creator     string `json:"creator" bson:"creator"`
	Modifier    string `json:"modifier" bson:"modifier"`
	CreateTime  *Time  `json:"create_time" bson:"create_time"`
	LastTime    *Time  `json:"last_time" bson:"last_time"`
}

// FieldTemplateAttr field template attribute definition
type FieldTemplateAttr struct {
	ID           int64           `json:"id" bson:"id"`
	TemplateID   int64           `json:"bk_template_id" bson:"bk_template_id"`
	PropertyID   string          `json:"bk_property_id" bson:"bk_property_id"`
	PropertyType string          `json:"bk_property_type" bson:"bk_property_type"`
	PropertyName string          `json:"bk_property_name" bson:"bk_property_name"`
	Unit         string          `json:"unit" bson:"unit"`
	Placeholder  AttrPlaceholder `json:"placeholder" bson:"placeholder"`
	Editable     AttrEditable    `json:"editable" bson:"editable"`
	Required     AttrRequired    `json:"isrequired" bson:"isrequired"`
	Option       interface{}     `json:"option" bson:"option"`
	Default      interface{}     `json:"default" bson:"default"`
	IsMultiple   bool            `json:"ismultiple" bson:"ismultiple"`
	OwnerID      string          `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Creator      string          `json:"creator" bson:"creator"`
	Modifier     string          `json:"modifier" bson:"modifier"`
	CreateTime   *Time           `json:"create_time" bson:"create_time"`
	LastTime     *Time           `json:"last_time" bson:"last_time"`
}

// AttrEditable field template attribute editable definition
type AttrEditable struct {
	Lock  bool `json:"lock" bson:"lock"`
	Value bool `json:"value" bson:"value"`
}

// AttrRequired field template attribute required definition
type AttrRequired struct {
	Lock  bool `json:"lock" bson:"lock"`
	Value bool `json:"value" bson:"value"`
}

// AttrPlaceholder field template attribute placeholder definition
type AttrPlaceholder struct {
	Lock  bool   `json:"lock" bson:"lock"`
	Value string `json:"value" bson:"value"`
}

// FieldTemplateUnique field template unique definition
type FieldTemplateUnique struct {
	ID         int64   `json:"id" bson:"id"`
	TemplateID int64   `json:"bk_template_id" bson:"bk_template_id"`
	Keys       []int64 `json:"keys" bson:"keys"`
	OwnerID    string  `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Creator    string  `json:"creator" bson:"creator"`
	Modifier   string  `json:"modifier" bson:"modifier"`
	CreateTime *Time   `json:"create_time" bson:"create_time"`
	LastTime   *Time   `json:"last_time" bson:"last_time"`
}

// ObjFieldTemplateRelation the relationship between model and field template definition
type ObjFieldTemplateRelation struct {
	ObjectID   string `json:"bk_obj_id" bson:"bk_obj_id"`
	TemplateID int64  `json:"bk_template_id" bson:"bk_template_id"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// FieldTemplateInfo field template info for list apis
type FieldTemplateInfo struct {
	Count uint64          `json:"count"`
	Info  []FieldTemplate `json:"info"`
}

// ListFieldTemplateResp list field template response
type ListFieldTemplateResp struct {
	BaseResp `json:",inline"`
	Data     FieldTemplateInfo `json:"data"`
}

// FieldTemplateBindObjOpt field template binding model option
type FieldTemplateBindObjOpt struct {
	ID        int64    `json:"id"`
	ObjectIDs []string `json:"bk_obj_ids"`
}

// Validate field template binding model request parameter validation function
func (option *FieldTemplateBindObjOpt) Validate() ccErr.RawErrorInfo {
	if option.ID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"id"},
		}
	}
	if len(option.ObjectIDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_obj_ids"},
		}
	}
	if len(option.ObjectIDs) > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"bk_obj_ids", common.BKMaxLimitSize},
		}
	}
	return ccErr.RawErrorInfo{}
}

// FieldTemplateUnbindObjOpt field template unbinding model option
type FieldTemplateUnbindObjOpt struct {
	ID       int64  `json:"id"`
	ObjectID string `json:"bk_obj_id"`
}

// Validate field template unbinding model request parameter validation function
func (option *FieldTemplateUnbindObjOpt) Validate() ccErr.RawErrorInfo {
	if option.ID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"id"},
		}
	}
	if option.ObjectID == "" {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKObjIDField},
		}
	}
	return ccErr.RawErrorInfo{}
}

// ListFieldTmplAttrOption list field template attribute option
type ListFieldTmplAttrOption struct {
	TemplateID        int64 `json:"bk_template_id"`
	CommonQueryOption `json:",inline"`
}

// Validate list field template attribute option
func (l *ListFieldTmplAttrOption) Validate() ccErr.RawErrorInfo {
	if l.TemplateID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKTemplateID}}
	}

	// set limit to unlimited if not set, compatible for searching all attributes, attributes amount won't be large
	if !l.Page.EnableCount && l.Page.Limit == 0 {
		l.Page.Limit = common.BKNoLimit
	}

	if rawErr := l.CommonQueryOption.Validate(); rawErr.ErrCode != 0 {
		return rawErr
	}

	return ccErr.RawErrorInfo{}
}

// FieldTemplateAttrInfo field template attribute info for list apis
type FieldTemplateAttrInfo struct {
	Count uint64              `json:"count"`
	Info  []FieldTemplateAttr `json:"info"`
}

// ListFieldTemplateAttrResp list field template attribute response
type ListFieldTemplateAttrResp struct {
	BaseResp `json:",inline"`
	Data     FieldTemplateAttrInfo `json:"data"`
}
