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
	"regexp"
	"strings"
	"unicode/utf8"

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

const (
	fieldTemplateNameMaxLen = 15
	fieldTemplateDesMaxLen  = 100
)

// Validate validate FieldTemplate
func (f *FieldTemplate) Validate() ccErr.RawErrorInfo {
	if len(f.Name) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldName}}
	}

	if utf8.RuneCountInString(f.Name) > fieldTemplateNameMaxLen {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommValExceedMaxFailed,
			Args: []interface{}{common.BKFieldName, fieldTemplateNameMaxLen}}
	}

	if utf8.RuneCountInString(f.Description) > fieldTemplateDesMaxLen {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommValExceedMaxFailed,
			Args: []interface{}{common.BKDescriptionField, fieldTemplateDesMaxLen}}
	}

	return ccErr.RawErrorInfo{}
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

// Validate validate FieldTemplateAttr
func (f *FieldTemplateAttr) Validate() ccErr.RawErrorInfo {
	if f.TemplateID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKTemplateID}}
	}

	if err := f.validatePropertyID(); err.ErrCode != 0 {
		return err
	}

	if err := f.validateType(); err.ErrCode != 0 {
		return err
	}

	if err := f.validateName(); err.ErrCode != 0 {
		return err
	}

	f.Unit = strings.TrimSpace(f.Unit)
	if f.Unit != "" && common.AttributeUnitMaxLength < utf8.RuneCountInString(f.Unit) {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommValExceedMaxFailed,
			Args: []interface{}{AttributeFieldUnit, common.AttributeUnitMaxLength}}
	}

	if err := f.validatePlaceholder(); err.ErrCode != 0 {
		return err
	}

	// because there will be a package import cycle problem,
	// validate option is in src/source_controller/coreservice/core/model/field_template.go file,
	// call valid.ValidPropertyOption func

	return ccErr.RawErrorInfo{}
}

func (f *FieldTemplateAttr) validatePropertyID() ccErr.RawErrorInfo {
	f.PropertyID = strings.TrimSpace(f.PropertyID)

	if f.PropertyID == "" {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet,
			Args: []interface{}{common.BKPropertyIDField}}
	}

	if utf8.RuneCountInString(f.PropertyID) > common.AttributeIDMaxLength {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommValExceedMaxFailed,
			Args: []interface{}{common.BKPropertyIDField, common.AttributeIDMaxLength}}
	}

	match, err := regexp.MatchString(common.FieldTypeStrictCharRegexp, f.PropertyID)
	if err != nil || !match {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsIsInvalid,
			Args: []interface{}{common.BKPropertyIDField}}
	}

	if strings.HasPrefix(f.PropertyID, "bk_") || strings.HasPrefix(f.PropertyID, "_bk") {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsIsInvalid,
			Args: []interface{}{common.BKPropertyIDField}}
	}

	return ccErr.RawErrorInfo{}
}

func (f *FieldTemplateAttr) validateType() ccErr.RawErrorInfo {
	f.PropertyType = strings.TrimSpace(f.PropertyType)

	if f.PropertyType == "" {
		return ccErr.RawErrorInfo{}
	}

	switch f.PropertyType {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar, common.FieldTypeInt, common.FieldTypeFloat,
		common.FieldTypeEnumMulti, common.FieldTypeDate, common.FieldTypeTime, common.FieldTypeUser,
		common.FieldTypeOrganization, common.FieldTypeTimeZone, common.FieldTypeBool, common.FieldTypeList:

	default:
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsIsInvalid,
			Args: []interface{}{AttributeFieldPropertyType}}
	}

	return ccErr.RawErrorInfo{}
}

func (f *FieldTemplateAttr) validateName() ccErr.RawErrorInfo {
	f.PropertyName = strings.TrimSpace(f.PropertyName)

	if f.PropertyName == "" {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet,
			Args: []interface{}{common.BKPropertyNameField}}
	}

	if common.AttributeNameMaxLength < utf8.RuneCountInString(f.PropertyName) {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommValExceedMaxFailed,
			Args: []interface{}{common.BKPropertyNameField, common.AttributeNameMaxLength}}
	}

	return ccErr.RawErrorInfo{}
}

func (f *FieldTemplateAttr) validatePlaceholder() ccErr.RawErrorInfo {
	f.Placeholder.Value = strings.TrimSpace(f.Placeholder.Value)

	if f.Placeholder.Value != "" {
		if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(f.Placeholder.Value) {
			return ccErr.RawErrorInfo{ErrCode: common.CCErrCommValExceedMaxFailed,
				Args: []interface{}{AttributeFieldPlaceHolder, common.AttributePlaceHolderMaxLength}}
		}

		match, err := regexp.MatchString(common.FieldTypeLongCharRegexp, f.Placeholder.Value)
		if err != nil || !match {
			return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsIsInvalid,
				Args: []interface{}{AttributeFieldPlaceHolder}}
		}
	}

	return ccErr.RawErrorInfo{}
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

// FieldTmplUniqueCommonField field template unique common field definition
type FieldTmplUniqueCommonField struct {
	ID         int64  `json:"id" bson:"id"`
	TemplateID int64  `json:"bk_template_id" bson:"bk_template_id"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Creator    string `json:"creator" bson:"creator"`
	Modifier   string `json:"modifier" bson:"modifier"`
	CreateTime *Time  `json:"create_time" bson:"create_time"`
	LastTime   *Time  `json:"last_time" bson:"last_time"`
}

// FieldTemplateUnique field template unique definition
type FieldTemplateUnique struct {
	FieldTmplUniqueCommonField `json:",inline" bson:",inline"`
	Keys                       []int64 `json:"keys" bson:"keys"`
}

// Validate validate FieldTemplateUnique
func (f *FieldTemplateUnique) Validate() ccErr.RawErrorInfo {
	if f.TemplateID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKTemplateID}}
	}

	if len(f.Keys) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet,
			Args: []interface{}{common.BKObjectUniqueKeys}}
	}

	for _, key := range f.Keys {
		if key == 0 {
			return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid,
				Args: []interface{}{common.BKObjectUniqueKeys}}
		}
	}

	return ccErr.RawErrorInfo{}
}

// Convert convert FieldTemplateUnique to FieldTmplUniqueOption struct
func (c *FieldTemplateUnique) Convert(idToPropertyIDMap map[int64]string) (*FieldTmplUniqueOption,
	ccErr.RawErrorInfo) {

	propertyIDs := make([]string, len(c.Keys))
	for idx, key := range c.Keys {
		propertyID, exist := idToPropertyIDMap[key]
		if !exist {
			return nil, ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid,
				Args: []interface{}{common.BKObjectUniqueKeys}}
		}

		propertyIDs[idx] = propertyID
	}

	unique := new(FieldTmplUniqueOption)
	unique.FieldTmplUniqueCommonField = c.FieldTmplUniqueCommonField
	unique.Keys = propertyIDs

	return unique, ccErr.RawErrorInfo{}
}

// ObjFieldTemplateRelation the relationship between model and field template definition
type ObjFieldTemplateRelation struct {
	ObjectID   int64  `json:"object_id" bson:"object_id"`
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
	ID        int64   `json:"bk_template_id"`
	ObjectIDs []int64 `json:"object_ids"`
}

// Validate field template binding model request parameter validation function
func (option *FieldTemplateBindObjOpt) Validate() ccErr.RawErrorInfo {
	if option.ID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}
	if len(option.ObjectIDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"object_ids"},
		}
	}
	if len(option.ObjectIDs) > common.BKMaxLimitSize {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"object_ids", common.BKMaxLimitSize},
		}
	}
	return ccErr.RawErrorInfo{}
}

// FieldTemplateUnbindObjOpt field template unbinding model option
type FieldTemplateUnbindObjOpt struct {
	ID       int64 `json:"bk_template_id"`
	ObjectID int64 `json:"object_id"`
}

// Validate field template unbinding model request parameter validation function
func (option *FieldTemplateUnbindObjOpt) Validate() ccErr.RawErrorInfo {
	if option.ID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}
	if option.ObjectID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"object_id"},
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

// CreateFieldTmplOption create field template option
type CreateFieldTmplOption struct {
	FieldTemplate `json:",inline"`
	Attributes    []FieldTemplateAttr     `json:"attributes"`
	Uniques       []FieldTmplUniqueOption `json:"uniques"`
}

const (
	FieldTemplateAttrMaxCount   = 20
	FieldTemplateUniqueMaxCount = 5
)

// Validate validate create field template option
func (c *CreateFieldTmplOption) Validate() ccErr.RawErrorInfo {
	if err := c.FieldTemplate.Validate(); err.ErrCode != 0 {
		return err
	}

	if len(c.Attributes) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{"attributes"}}
	}

	if len(c.Attributes) > FieldTemplateAttrMaxCount {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit,
			Args: []interface{}{"attributes", FieldTemplateAttrMaxCount}}
	}

	if len(c.Uniques) > FieldTemplateUniqueMaxCount {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit,
			Args: []interface{}{"uniques", FieldTemplateAttrMaxCount}}
	}

	return ccErr.RawErrorInfo{}
}

// FieldTmplUniqueOption create field template unique option
type FieldTmplUniqueOption struct {
	FieldTmplUniqueCommonField `json:",inline" bson:",inline"`
	Keys                       []string `json:"keys" bson:"keys"`
}

// Convert convert FieldTmplUniqueOption to FieldTemplateUnique struct
func (c *FieldTmplUniqueOption) Convert(propertyIDToIDMap map[string]int64) (*FieldTemplateUnique,
	ccErr.RawErrorInfo) {

	ids := make([]int64, len(c.Keys))
	for idx, key := range c.Keys {
		id, exist := propertyIDToIDMap[key]
		if !exist {
			return nil, ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid,
				Args: []interface{}{common.BKObjectUniqueKeys}}
		}

		ids[idx] = id
	}

	unique := new(FieldTemplateUnique)
	unique.FieldTmplUniqueCommonField = c.FieldTmplUniqueCommonField
	unique.Keys = ids

	return unique, ccErr.RawErrorInfo{}
}

// ListFieldTmplUniqueOption list field template unique option
type ListFieldTmplUniqueOption struct {
	TemplateID        int64 `json:"bk_template_id"`
	CommonQueryOption `json:",inline"`
}

// Validate list field template unique option
func (l *ListFieldTmplUniqueOption) Validate() ccErr.RawErrorInfo {
	if l.TemplateID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKTemplateID}}
	}

	// set limit to unlimited if not set, compatible for searching all uniques, uniques amount won't be large
	if !l.Page.EnableCount && l.Page.Limit == 0 {
		l.Page.Limit = common.BKNoLimit
	}

	if rawErr := l.CommonQueryOption.Validate(); rawErr.ErrCode != 0 {
		return rawErr
	}

	return ccErr.RawErrorInfo{}
}

// FieldTemplateUniqueInfo field template unique info for list apis
type FieldTemplateUniqueInfo struct {
	Count uint64                `json:"count"`
	Info  []FieldTemplateUnique `json:"info"`
}

// ListFieldTmplUniqueResp list field template unique response
type ListFieldTmplUniqueResp struct {
	BaseResp `json:",inline"`
	Data     FieldTemplateUniqueInfo `json:"data"`
}