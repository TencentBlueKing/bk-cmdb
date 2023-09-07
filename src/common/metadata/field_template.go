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
	"configcenter/src/common/mapstr"
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
	// fieldTemplateSyncMaxNum compare the difference status between the
	// template and the model The maximum number of models processed at one time
	fieldTemplateSyncMaxNum = 5
	fieldTemplateNameMaxLen = 128
	fieldTemplateDesMaxLen  = 2000
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
	ID           int64  `json:"id" bson:"id"`
	TemplateID   int64  `json:"bk_template_id" bson:"bk_template_id"`
	PropertyID   string `json:"bk_property_id" bson:"bk_property_id"`
	PropertyType string `json:"bk_property_type" bson:"bk_property_type"`
	PropertyName string `json:"bk_property_name" bson:"bk_property_name"`
	// PropertyIndex It is used to display field template attribute in order.
	// When a template attribute array is created or updated,
	// its value is set according to the order passed in by the front end.
	PropertyIndex int64           `json:"bk_property_index" bson:"bk_property_index"`
	Unit          string          `json:"unit" bson:"unit"`
	Placeholder   AttrPlaceholder `json:"placeholder" bson:"placeholder"`
	Editable      AttrEditable    `json:"editable" bson:"editable"`
	Required      AttrRequired    `json:"isrequired" bson:"isrequired"`
	Option        interface{}     `json:"option" bson:"option"`
	Default       interface{}     `json:"default" bson:"default"`
	IsMultiple    bool            `json:"ismultiple" bson:"ismultiple"`
	OwnerID       string          `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Creator       string          `json:"creator" bson:"creator"`
	Modifier      string          `json:"modifier" bson:"modifier"`
	CreateTime    *Time           `json:"create_time" bson:"create_time"`
	LastTime      *Time           `json:"last_time" bson:"last_time"`
}

// ValidateBase validate field template attribute
func (f *FieldTemplateAttr) ValidateBase() ccErr.RawErrorInfo {

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
	// call attrvalid.ValidPropertyOption func

	return ccErr.RawErrorInfo{}
}

// Validate validate FieldTemplateAttr
func (f *FieldTemplateAttr) Validate() ccErr.RawErrorInfo {
	if f.TemplateID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKTemplateID}}
	}

	if err := f.ValidateBase(); err.ErrCode != 0 {
		return err
	}

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
	// FieldTemplateAttrMaxCount filed template attribute max count
	FieldTemplateAttrMaxCount = 20

	// FieldTemplateUniqueMaxCount filed template unique max count
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

// ListObjFieldTmplRelOption list field template and object relation option
type ListObjFieldTmplRelOption struct {
	TemplateIDs []int64 `json:"bk_template_ids"`
	ObjectIDs   []int64 `json:"object_ids"`
}

// Validate list field template and object relation option
func (l *ListObjFieldTmplRelOption) Validate() ccErr.RawErrorInfo {
	if len(l.TemplateIDs) == 0 && len(l.ObjectIDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_template_ids and object_ids"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ObjFieldTmplRelInfo field template and object relation info for list apis
type ObjFieldTmplRelInfo struct {
	Count uint64                     `json:"count"`
	Info  []ObjFieldTemplateRelation `json:"info"`
}

// ListObjFieldTmplRelResp list field template and object relation response
type ListObjFieldTmplRelResp struct {
	BaseResp `json:",inline"`
	Data     ObjFieldTmplRelInfo `json:"data"`
}

// ListFieldTmplByObjOption list field template by related object option
type ListFieldTmplByObjOption struct {
	ObjectID int64 `json:"object_id"`
}

// Validate list field template by related object option
func (l *ListFieldTmplByObjOption) Validate() ccErr.RawErrorInfo {
	if l.ObjectID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.ObjectIDField},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ListObjByFieldTmplOption list object by related field template option
type ListObjByFieldTmplOption struct {
	TemplateID        int64 `json:"bk_template_id"`
	CommonQueryOption `json:",inline"`
}

// Validate list object by related field template option
func (l *ListObjByFieldTmplOption) Validate() ccErr.RawErrorInfo {
	if l.TemplateID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}

	if rawErr := l.CommonQueryOption.Validate(); rawErr.ErrCode != 0 {
		return rawErr
	}

	return ccErr.RawErrorInfo{}
}

// FieldTemplateSyncOption synchronization of field combination templates to model requests
type FieldTemplateSyncOption struct {
	TemplateID int64   `json:"bk_template_id"`
	ObjectIDs  []int64 `json:"object_ids"`
}

// Validate list object by related field template option
func (op *FieldTemplateSyncOption) Validate() ccErr.RawErrorInfo {
	if op.TemplateID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}
	if len(op.ObjectIDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"object_ids"},
		}
	}

	for _, objID := range op.ObjectIDs {
		if objID == 0 {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsIsInvalid,
				Args:    []interface{}{common.ObjectIDField},
			}
		}
	}
	return ccErr.RawErrorInfo{}
}

// CompareFieldTmplAttrOption compare field template attribute with object option
type CompareFieldTmplAttrOption struct {
	TemplateID int64               `json:"bk_template_id"`
	ObjectID   int64               `json:"object_id"`
	Attrs      []FieldTemplateAttr `json:"attributes"`
	// IsPartial partial comparison of template and model attre flag.
	IsPartial bool `json:"is_partial"`
}

// Validate compare field template attribute with object option
func (l *CompareFieldTmplAttrOption) Validate() ccErr.RawErrorInfo {
	if l.TemplateID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}

	if l.ObjectID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	if len(l.Attrs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"attributes"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// SyncObjectTask synchronize field combination template information to model request
type SyncObjectTask struct {
	TemplateID int64 `json:"bk_template_id"`
	ObjectID   int64 `json:"object_id"`
}

// Validate check of SyncObjectTask.
func (option *SyncObjectTask) Validate() ccErr.RawErrorInfo {
	if option.TemplateID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}
	if option.ObjectID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{common.ObjectIDField},
		}
	}
	return ccErr.RawErrorInfo{}
}

// CompareFieldTmplAttrsRes compare field template attributes with object result
type CompareFieldTmplAttrsRes struct {
	Create    []CompareOneFieldTmplAttrRes `json:"create"`
	Update    []CompareOneFieldTmplAttrRes `json:"update"`
	Conflict  []CompareOneFieldTmplAttrRes `json:"conflict"`
	Unchanged []Attribute                  `json:"unchanged"`
}

// CompareOneFieldTmplAttrRes compare one field template attribute with object result
type CompareOneFieldTmplAttrRes struct {
	// Index field template's original index in input attribute array
	Index int `json:"index"`
	// PropertyID field template attribute property id
	PropertyID string `json:"bk_property_id"`
	// Message conflict message
	Message string `json:"message,omitempty"`
	// Data original data of object for update/conflict attribute
	Data *Attribute `json:"data,omitempty"`
	// UpdateData to be updated attribute data
	UpdateData mapstr.MapStr `json:"update_data,omitempty"`
}

// CompareFieldTmplUniqueOption compare field template unique with object option
type CompareFieldTmplUniqueOption struct {
	TemplateID int64                      `json:"bk_template_id"`
	ObjectID   int64                      `json:"object_id"`
	Uniques    []FieldTmplUniqueForUpdate `json:"uniques"`
	// IsPartial partial comparison of template and model unique flag.
	IsPartial bool `json:"is_partial"`
}

// FieldTmplUniqueForUpdate field template unique for field template update scenario.
type FieldTmplUniqueForUpdate struct {
	ID int64 `json:"id"`
	// some related attributes may not be created yet, so we can only use property id for keys
	// field template attr with the same property id can not be recreated, so property id can uniques identify an attr
	Keys []string `json:"keys"`
}

// Validate compare field template unique with object option
func (l *CompareFieldTmplUniqueOption) Validate() ccErr.RawErrorInfo {
	if l.TemplateID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}

	if l.ObjectID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.ObjectIDField},
		}
	}

	return ccErr.RawErrorInfo{}
}

// CompareFieldTmplUniquesRes compare field template uniques with object result
type CompareFieldTmplUniquesRes struct {
	Create    []CompareOneFieldTmplUniqueRes `json:"create"`
	Update    []CompareOneFieldTmplUniqueRes `json:"update"`
	Conflict  []CompareOneFieldTmplUniqueRes `json:"conflict"`
	Unchanged []ObjectUnique                 `json:"unchanged"`
}

// CompareOneFieldTmplUniqueRes compare one field template unique with object result
type CompareOneFieldTmplUniqueRes struct {
	// Index field template's original index in input unique array
	Index int `json:"index"`
	// Message conflict message
	Message string `json:"message,omitempty"`
	// Data original data of object for update/conflict unique
	Data *ObjectUnique `json:"data,omitempty"`
}

// DeleteFieldTmplOption delete field template option
type DeleteFieldTmplOption struct {
	ID int64 `json:"id"`
}

// Validate delete field template option
func (d *DeleteFieldTmplOption) Validate() ccErr.RawErrorInfo {
	if d.ID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldID}}
	}

	return ccErr.RawErrorInfo{}
}

// CloneFieldTmplOption clone field template option
type CloneFieldTmplOption struct {
	ID            int64 `json:"id"`
	FieldTemplate `json:",inline"`
}

// Validate clone field template option
func (c *CloneFieldTmplOption) Validate() ccErr.RawErrorInfo {
	if c.ID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldID}}
	}

	if err := c.FieldTemplate.Validate(); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// UpdateFieldTmplOption update field template option
type UpdateFieldTmplOption struct {
	FieldTemplate `json:",inline"`
	Attributes    []FieldTemplateAttr     `json:"attributes"`
	Uniques       []FieldTmplUniqueOption `json:"uniques"`
}

// Validate update field template option
func (c *UpdateFieldTmplOption) Validate() ccErr.RawErrorInfo {
	if c.ID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldID}}
	}

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
			Args: []interface{}{"uniques", FieldTemplateUniqueMaxCount}}
	}

	return ccErr.RawErrorInfo{}
}

// ListTmplSimpleByAttrOption query the brief information of the field
// template according to the template ID of the model object attr for UI.
type ListTmplSimpleByAttrOption struct {
	TemplateID int64 `json:"bk_template_id"`
	AttrID     int64 `json:"bk_attribute_id"`
}

// Validate verify the legitimacy of the request.
func (c *ListTmplSimpleByAttrOption) Validate() ccErr.RawErrorInfo {
	if c.TemplateID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}

	if c.AttrID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAttributeIDField},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ListTmplSimpleByUniqueOption query the brief information of the field
// template according to the template ID of the model unique for UI.
type ListTmplSimpleByUniqueOption struct {
	TemplateID int64 `json:"bk_template_id"`
	UniqueID   int64 `json:"bk_unique_id"`
}

// Validate verify the legitimacy of the request.
func (c *ListTmplSimpleByUniqueOption) Validate() ccErr.RawErrorInfo {
	if c.TemplateID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKTemplateID},
		}
	}
	if c.UniqueID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_unique_id"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// ListFieldTemplateSimpleResp returns brief information about the template
type ListFieldTemplateSimpleResp struct {
	BaseResp `json:",inline"`
	Data     ListTmplSimpleResult `json:"data"`
}

// ListTmplSimpleResult returns brief information about the template
type ListTmplSimpleResult struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListFieldTmpltSyncStatusOption used to compare templates and model
// attributes or uniquely check whether there is a difference request.
type ListFieldTmpltSyncStatusOption struct {
	ID        int64   `json:"bk_template_id"`
	ObjectIDs []int64 `json:"object_ids"`
}

// Validate judging the legality of parameters
func (option *ListFieldTmpltSyncStatusOption) Validate() ccErr.RawErrorInfo {
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
	if len(option.ObjectIDs) > fieldTemplateSyncMaxNum {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"object_ids", fieldTemplateSyncMaxNum},
		}
	}
	for _, id := range option.ObjectIDs {
		if id == 0 {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{common.ObjectIDField},
			}
		}
	}
	return ccErr.RawErrorInfo{}
}

// ListFieldTmpltSyncStatusResult it is used to compare the attributes or unique
// verification status comparison results between the template and the model
type ListFieldTmpltSyncStatusResult struct {
	ObjectID int64 `json:"object_id"`
	NeedSync bool  `json:"need_sync"`
}

// ListFieldTmplModelStatusOption list field template model status option
type ListFieldTmplModelStatusOption struct {
	ID        int64   `json:"bk_template_id"`
	ObjectIDs []int64 `json:"object_ids"`
}

// Validate list field template model status option
func (option *ListFieldTmplModelStatusOption) Validate() ccErr.RawErrorInfo {
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
	if len(option.ObjectIDs) > fieldTemplateSyncMaxNum {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"object_ids", fieldTemplateSyncMaxNum},
		}
	}
	for _, id := range option.ObjectIDs {
		if id == 0 {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{common.ObjectIDField},
			}
		}
	}
	return ccErr.RawErrorInfo{}
}
