/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/common/valid/attribute/manager"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// AttributeFieldID TODO
	AttributeFieldID = "id"
	// AttributeFieldSupplierAccount TODO
	AttributeFieldSupplierAccount = "bk_supplier_account"
	// AttributeFieldObjectID TODO
	AttributeFieldObjectID = "bk_obj_id"
	// AttributeFieldPropertyID TODO
	AttributeFieldPropertyID = "bk_property_id"
	// AttributeFieldPropertyName TODO
	AttributeFieldPropertyName = "bk_property_name"
	// AttributeFieldPropertyGroup TODO
	AttributeFieldPropertyGroup = "bk_property_group"
	// AttributeFieldPropertyGroupName 模型属性字段分组名称
	AttributeFieldPropertyGroupName = "bk_property_group_name"
	// AttributeFieldPropertyIndex TODO
	AttributeFieldPropertyIndex = "bk_property_index"
	// AttributeFieldUnit TODO
	AttributeFieldUnit = "unit"
	// AttributeFieldPlaceHolder TODO
	AttributeFieldPlaceHolder = "placeholder"
	// AttributeFieldIsEditable TODO
	AttributeFieldIsEditable = "editable"
	// AttributeFieldIsPre TODO
	AttributeFieldIsPre = "ispre"
	// AttributeFieldIsRequired TODO
	AttributeFieldIsRequired = "isrequired"
	// AttributeFieldIsReadOnly TODO
	AttributeFieldIsReadOnly = "isreadonly"
	// AttributeFieldIsOnly TODO
	AttributeFieldIsOnly = "isonly"
	// AttributeFieldIsSystem TODO
	AttributeFieldIsSystem = "bk_issystem"
	// AttributeFieldIsAPI TODO
	AttributeFieldIsAPI = "bk_isapi"
	// AttributeFieldPropertyType TODO
	AttributeFieldPropertyType = "bk_property_type"
	// AttributeFieldOption TODO
	AttributeFieldOption = "option"
	// AttributeFieldDescription TODO
	AttributeFieldDescription = "description"
	// AttributeFieldCreator TODO
	AttributeFieldCreator = "creator"
	// AttributeFieldCreateTime TODO
	AttributeFieldCreateTime = "create_time"
	// AttributeFieldLastTime TODO
	AttributeFieldLastTime = "last_time"
	// AttributeFieldDefault attribute default value field
	AttributeFieldDefault = "default"
	// AttributeFieldIsMultiple the is multiple name field
	AttributeFieldIsMultiple = "ismultiple"
)

const (
	// TableLongCharMaxNum the maximum number of long
	// characters supported by a form field.
	TableLongCharMaxNum = 2
	// TableHeaderMaxNum the maximum length of the table header field.
	TableHeaderMaxNum = 8
	// TableDefaultMaxLines the maximum length of the table default lines.
	TableDefaultMaxLines = 10
)

// Attribute attribute metadata definition
type Attribute struct {
	BizID             int64       `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	ID                int64       `field:"id" json:"id" bson:"id" mapstructure:"id"`
	OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id" mapstructure:"bk_obj_id"`
	PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id" mapstructure:"bk_property_id"`
	PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name" mapstructure:"bk_property_name"`
	PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group" mapstructure:"bk_property_group"`
	PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-" mapstructure:"bk_property_group_name"`
	PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index" mapstructure:"bk_property_index"`
	Unit              string      `field:"unit" json:"unit" bson:"unit" mapstructure:"unit"`
	Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder" mapstructure:"placeholder"`
	IsEditable        bool        `field:"editable" json:"editable" bson:"editable" mapstructure:"editable"`
	IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre" mapstructure:"ispre"`
	IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired" mapstructure:"isrequired"`
	IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly" mapstructure:"isreadonly"`
	IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly" mapstructure:"isonly"`
	IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem" mapstructure:"bk_issystem"`
	IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi" mapstructure:"bk_isapi"`
	PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type" mapstructure:"bk_property_type"`
	Option            interface{} `field:"option" json:"option" bson:"option" mapstructure:"option"`
	Default           interface{} `field:"default" json:"default,omitempty" bson:"default" mapstructure:"default"`
	IsMultiple        *bool       `field:"ismultiple" json:"ismultiple,omitempty" bson:"ismultiple" mapstructure:"ismultiple"`
	Description       string      `field:"description" json:"description" bson:"description" mapstructure:"description"`
	TemplateID        int64       `field:"bk_template_id" json:"bk_template_id" bson:"bk_template_id" mapstructure:"bk_template_id"`
	Creator           string      `field:"creator" json:"creator" bson:"creator" mapstructure:"creator"`
	CreateTime        *Time       `json:"create_time" bson:"create_time" mapstructure:"create_time"`
	LastTime          *Time       `json:"last_time" bson:"last_time" mapstructure:"last_time"`
}

// AttributeGroup attribute metadata definition
type AttributeGroup struct {
	ID         int64  `field:"id" json:"id" bson:"id"`
	OwnerID    string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	IsDefault  bool   `field:"bk_isdefault" json:"bk_isdefault" bson:"bk_isdefault"`
	IsPre      bool   `field:"ispre" json:"ispre" bson:"ispre"`
	GroupID    string `field:"bk_group_id" json:"bk_group_id" bson:"bk_group_id"`
	GroupName  string `field:"bk_group_name" json:"bk_group_name" bson:"bk_group_name"`
	GroupIndex int64  `field:"bk_group_index" json:"bk_group_index" bson:"bk_group_index"`
}

// Parse load the data from mapstr attribute into attribute instance
func (attribute *Attribute) Parse(data mapstr.MapStr) (*Attribute, error) {

	err := mapstr.SetValueToStructByTags(attribute, data)
	if nil != err {
		return nil, err
	}

	return attribute, err
}

// ToMapStr to mapstr
func (attribute *Attribute) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(attribute)
}

// ObjAttDes 对象模型属性
type ObjAttDes struct {
	Attribute         `json:",inline" bson:",inline"`
	PropertyGroupName string `json:"bk_property_group_name"`
}

// CreateObjAttDesResp 创建对象模型属性返回结构体
type CreateObjAttDesResp struct {
	BaseResp `json:",inline"`
	Data     ObjAttDes `json:"data"`
}

// Validate Attribute
func (attribute *Attribute) Validate(ctx context.Context, data interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	var attrValidatorMap = map[string]func(context.Context, interface{}, string) errors.RawErrorInfo{
		common.FieldTypeSingleChar:   attribute.validChar,
		common.FieldTypeLongChar:     attribute.validLongChar,
		common.FieldTypeInt:          attribute.validInt,
		common.FieldTypeFloat:        attribute.validFloat,
		common.FieldTypeEnum:         attribute.validEnum,
		common.FieldTypeEnumMulti:    attribute.validEnumMulti,
		common.FieldTypeEnumQuote:    attribute.validEnumQuote,
		common.FieldTypeDate:         attribute.validDate,
		common.FieldTypeTime:         attribute.validTime,
		common.FieldTypeTimeZone:     attribute.validTimeZone,
		common.FieldTypeBool:         attribute.validBool,
		common.FieldTypeUser:         attribute.validUser,
		common.FieldTypeList:         attribute.validList,
		common.FieldObject:           attribute.validObjectCondition,
		common.FieldTypeOrganization: attribute.validOrganization,
		common.FieldTypeInnerTable:   attribute.validInnerTable,
		common.FieldTypeIDRule:       attribute.validIDRule,
	}

	rawError := errors.RawErrorInfo{}
	fieldType := attribute.PropertyType
	switch fieldType {
	case "foreignkey", "singleasst", "multiasst":
		// TODO what validation should do on these types
	case common.FieldTypeTable:
		// TODO what validation should do on these types
		rawError = attribute.validTable(ctx, data, key)
	default:
		// notice: 注意default 这里的实现逻辑， 用break 做了执行流程的终止。 pr 建议，降低圈复杂度

		validator, exists := attrValidatorMap[fieldType]
		if exists {
			rawError = validator(ctx, data, key)
			break
		}
		// 是否为扩展字段类型
		if handle, ok := manager.Get(fieldType); ok {
			if err := handle.Validate(ctx, key, fieldType, attribute.IsRequired, attribute.Option, data); err != nil {
				blog.Errorf("validate attribute fail, field type: %s, err: %v, rid: %s", fieldType, err, rid)
				return errors.RawErrorInfo{
					ErrCode: common.CCErrCommParamsInvalid,
					Args:    []interface{}{err.Error()},
				}
			}
			break
		}
		rawError = errors.RawErrorInfo{
			ErrCode: common.CCErrCommUnexpectedFieldType,
			Args:    []interface{}{fieldType},
		}

	}
	// 如果出现了问题，并且报错原内容为propertyID，则替换为propertyName。
	if rawError.ErrCode != 0 {
		if key == attribute.PropertyID || key == common.BKPropertyValueField {
			rawError.Args = []interface{}{attribute.PropertyName}
		}
	}
	return rawError
}

// validTime valid object Attribute that is time type
func (attribute *Attribute) validTime(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {

	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil || val == "" {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	if _, ok := val.(time.Time); ok {
		return errors.RawErrorInfo{}
	}

	if _, result := util.IsTime(val); !result {
		blog.Errorf("params not valid, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	return errors.RawErrorInfo{}
}

// validDate valid object Attribute that is date type
func (attribute *Attribute) validDate(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil || val == "" {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	if result := util.IsDate(val); !result {
		blog.Errorf("params is not valid, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	return errors.RawErrorInfo{}
}

// validEnum valid object attribute that is enum type
func (attribute *Attribute) validEnum(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	// validate require
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	// validate type
	valStr, ok := val.(string)
	if !ok {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	// validate within enum
	enumOption, err := ParseEnumOption(attribute.Option)
	if err != nil {
		blog.Warnf("parse enum option failed, err: %v, rid: %s", err, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	for _, k := range enumOption {
		if k.ID == valStr {
			return errors.RawErrorInfo{}
		}
	}
	blog.V(3).Infof("params %s not valid, option %#v, raw option %#v, value: %#v, rid: %s", key, enumOption,
		attribute.Option, val, rid)
	blog.Errorf("params %s not valid , enum value: %#v, rid: %s", key, val, rid)
	return errors.RawErrorInfo{
		ErrCode: common.CCErrCommParamsInvalid,
		Args:    []interface{}{key},
	}
}

// CheckInterfaceSliceType check whether propertyValue type is []interface{} or primitive.A
func CheckInterfaceSliceType(propertyValue interface{}) ([]interface{}, errors.CCErrorCoder) {
	switch t := propertyValue.(type) {
	case []interface{}:
		return t, nil
	case primitive.A:
		return t, nil
	default:
		return nil, errors.New(common.CCErrCommUnexpectedFieldType, "property value type error")
	}
}

// validEnum valid object attribute that is enum multi type
func (attribute *Attribute) validEnumMulti(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	// validate require
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}
		}
		return errors.RawErrorInfo{}
	}

	enumOption, err := ParseEnumOption(attribute.Option)
	if err != nil {
		blog.Errorf("parse enum option failed, err: %v, rid: %s", err, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	idMap := make(map[string]struct{}, 0)
	for _, option := range enumOption {
		idMap[option.ID] = struct{}{}
	}

	valIDs, ccErr := CheckInterfaceSliceType(val)
	if ccErr != nil {
		blog.Errorf("convert val to interface slice failed, val type: %T", val)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	if len(valIDs) == 0 && attribute.IsRequired {
		blog.Errorf("data must be set, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	if len(valIDs) == 0 {
		return errors.RawErrorInfo{}
	}

	if attribute.IsMultiple == nil {
		blog.Errorf("multi flag must be set, rid: %s", rid)
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{key}}
	}

	if !(*attribute.IsMultiple) && len(valIDs) != 1 {
		blog.Errorf("multiple values are not allowed, valIDs: %+v, rid: %s", valIDs, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSingleChoice,
			Args:    []interface{}{key},
		}
	}
	for _, id := range valIDs {
		idVal, ok := id.(string)
		if !ok {
			blog.Errorf("data must be string id: %v, rid: %s", id, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{key},
			}
		}
		if _, ok := idMap[idVal]; !ok {
			blog.Errorf("value entered must be in the enumerated list，id: %s, rid: %s", id, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{key},
			}
		}
	}

	return errors.RawErrorInfo{}
}

// validEnum valid object attribute that is enum quote type
func (attribute *Attribute) validEnumQuote(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	// validate require
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}
		}
		return errors.RawErrorInfo{}
	}

	switch val.(type) {
	case []interface{}:
	case bson.A:
	default:
		blog.Errorf("params should be type enum quote, but its type is %T, rid: %s", val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validBool valid object attribute that is bool type
func (attribute *Attribute) validBool(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	switch val.(type) {
	case bool:
	default:
		blog.Errorf("params should be bool, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedBool,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validTimeZone valid char valid object attribute that is timezone type
func (attribute *Attribute) validTimeZone(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	if ok := util.IsTimeZone(val); !ok {
		blog.Errorf("params should be timezone, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedTimeZone,
			Args:    []interface{}{key},
		}
	}

	return errors.RawErrorInfo{}
}

// validInt valid object attribute that is int type
func (attribute *Attribute) validInt(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	if !util.IsNumeric(val) {
		blog.Errorf("params %s:%#v not int, rid: %s", key, val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedInt,
			Args:    []interface{}{key},
		}
	}

	intOption, err := ParseIntOption(attribute.Option)
	if err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}

	value, _ := util.GetInt64ByInterface(val)
	if value > intOption.Max || value < intOption.Min {
		blog.Errorf("params %s:%#v not valid, rid: %s", key, val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validFloat valid object attribute that is float type
func (attribute *Attribute) validFloat(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	var value float64
	value, err := util.GetFloat64ByInterface(val)
	if err != nil {
		blog.Errorf("params %s:%#v not float, rid: %s", key, val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{key},
		}
	}

	floatOption, err := ParseFloatOption(attribute.Option)
	if err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}

	if value > floatOption.Max || value < floatOption.Min {
		blog.Errorf("params %s:%#v not valid, rid: %s", key, val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validLongChar valid object attribute that is long char type
func (attribute *Attribute) validLongChar(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil || val == "" {
		if attribute.IsRequired {
			blog.Errorf("params in need, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	switch value := val.(type) {
	case string:
		value = strings.TrimSpace(value)
		if len(value) > common.FieldTypeLongLenChar {
			blog.Errorf("params over length %d, rid: %s", common.FieldTypeLongLenChar, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommOverLimit,
				Args:    []interface{}{key},
			}
		}
		if len(value) == 0 {
			if attribute.IsRequired {
				blog.Errorf("params can not be empty, rid: %s", rid)
				return errors.RawErrorInfo{
					ErrCode: common.CCErrCommParamsNeedSet,
					Args:    []interface{}{key},
				}
			}
			return errors.RawErrorInfo{}
		}

		option, ok := attribute.Option.(string)
		if !ok {
			break
		}
		strReg, err := regexp.Compile(option)
		if err != nil {
			blog.Errorf(`regexp "%s" invalid, err: %v, rid: %s`, option, err, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsIsInvalid,
				Args:    []interface{}{option},
			}
		}
		if !strReg.MatchString(value) {
			blog.Errorf(`params "%s" not match regexp "%s", rid: %s`, val, option, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrFieldRegValidFailed,
				Args:    []interface{}{key},
			}
		}
	default:
		blog.Errorf("params should be string, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedString,
			Args:    []interface{}{key},
		}
	}

	return errors.RawErrorInfo{}
}

// validChar valid object attribute that is char type
func (attribute *Attribute) validChar(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params in need, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}
		}
		return errors.RawErrorInfo{}
	}

	switch value := val.(type) {
	case string:
		value = strings.TrimSpace(value)
		if len(value) > common.FieldTypeSingleLenChar {
			blog.Errorf("params over length %d, rid: %s", common.FieldTypeSingleLenChar, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommOverLimit,
				Args:    []interface{}{key},
			}
		}
		if len(value) == 0 {
			if attribute.IsRequired {
				blog.Errorf("params can not be empty, rid: %s", rid)
				return errors.RawErrorInfo{
					ErrCode: common.CCErrCommParamsNeedSet,
					Args:    []interface{}{key},
				}
			}
			return errors.RawErrorInfo{}
		}

		if key == common.BKAppNameField || key == common.BKSetNameField || key == common.BKModuleNameField {
			if strings.Contains(value, "##") {
				blog.Errorf("params %s contains TopoModuleName's split flag ##, rid: %s", value, rid)
				return errors.RawErrorInfo{
					ErrCode: common.CCErrCommParamsInvalid,
					Args:    []interface{}{value},
				}
			}
		}

		if val == "" {
			return errors.RawErrorInfo{}
		}

		option, ok := attribute.Option.(string)
		if !ok {
			break
		}
		strReg, err := regexp.Compile(option)
		if err != nil {
			blog.Errorf(`regexp "%s" invalid, err: %v, rid: %s`, option, err, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsIsInvalid,
				Args:    []interface{}{option},
			}
		}
		if !strReg.MatchString(value) {
			blog.Errorf(`params "%s" not match regexp "%s", rid: %s`, val, option, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrFieldRegValidFailed,
				Args:    []interface{}{key},
			}
		}
	default:
		blog.Errorf("params should be string, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedString,
			Args:    []interface{}{key},
		}
	}

	return errors.RawErrorInfo{}
}

// validUser valid object attribute that is user type
func (attribute *Attribute) validUser(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil || val == "" {
		if attribute.IsRequired {
			blog.Errorf("params in need, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	switch value := val.(type) {
	case string:
		value = strings.TrimSpace(value)
		if len(value) > common.FieldTypeUserLenChar {
			blog.Errorf("params over length %d, rid: %s", common.FieldTypeUserLenChar, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommOverLimit,
				Args:    []interface{}{key},
			}
		}

		if len(value) == 0 {
			if attribute.IsRequired {
				blog.Errorf("params can not be empty, rid: %s", rid)
				return errors.RawErrorInfo{
					ErrCode: common.CCErrCommParamsNeedSet,
					Args:    []interface{}{key},
				}
			}
			return errors.RawErrorInfo{}
		}

		if attribute.IsMultiple == nil {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}
		}

		if !(*attribute.IsMultiple) && len(strings.Split(value, ",")) != 1 {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{key},
			}
		}

		// regex check
		match := util.IsUser(value)
		if !match {
			blog.Errorf(`value "%s" not match regexp, rid: %s`, value, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrFieldRegValidFailed,
				Args:    []interface{}{key},
			}
		}
	default:
		blog.Errorf("params should be string, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedString,
			Args:    []interface{}{key},
		}
	}

	return errors.RawErrorInfo{}
}

// validObjectCondition valid object attribute that is user type
func (attribute *Attribute) validObjectCondition(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {

	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil || val == "" {
		if attribute.IsRequired {
			blog.Errorf("params in need, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	// 对于对象的校验只需要判断类型是否是map[string]interface和MapStr即可

	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
	case reflect.Ptr:
		switch reflect.TypeOf(val).Elem().Kind() {
		case reflect.Map:
		default:
			blog.Errorf("object type is error, must be map, type: %v, rid: %s", reflect.TypeOf(val).Elem().Kind(),
				rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{key},
			}
		}

	default:
		blog.Errorf("object type is error, must be map, type: %v, rid: %s", reflect.TypeOf(val), rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	return errors.RawErrorInfo{}
}

func (attribute *Attribute) validList(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestUserFromContext(ctx)

	if val == nil {
		if attribute.IsRequired {
			blog.Error("params can not be null, list field key: %s, rid: %s", key, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}
		}
		return errors.RawErrorInfo{}
	}

	strVal, ok := val.(string)
	if !ok {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	var listOpt []interface{}
	switch listOption := attribute.Option.(type) {
	case []interface{}:
		listOpt = listOption
	case bson.A:
		listOpt = listOption
	default:
		blog.Errorf("option %v invalid, not string type list option, but type %T, rid: %s", attribute.Option,
			attribute.Option, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	for _, inVal := range listOpt {
		inValStr, ok := inVal.(string)
		if !ok {
			blog.Errorf("inner list option convert to string failed, params %s not valid , list field value: %#v, "+
				"rid: %s", key, val, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrParseAttrOptionListFailed,
				Args:    []interface{}{key},
			}
		}
		if strVal == inValStr {
			return errors.RawErrorInfo{}
		}
	}
	blog.Errorf("params %s not valid, option %#v, raw option %#v, value: %#v, rid: %s", key, listOpt, attribute,
		val, rid)
	return errors.RawErrorInfo{
		ErrCode: common.CCErrCommParamsInvalid,
		Args:    []interface{}{key},
	}
}

// validOrganization valid object attribute that is organization type
func (attribute *Attribute) validOrganization(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{key}}
		}

		return errors.RawErrorInfo{}
	}

	switch org := val.(type) {
	case []interface{}:
		if rawErr := attribute.validOrganizationValue(org, key, rid); rawErr.ErrCode != 0 {
			return rawErr
		}
	case bson.A:
		if rawErr := attribute.validOrganizationValue(org, key, rid); rawErr.ErrCode != 0 {
			return rawErr
		}
	default:
		blog.Errorf("params should be type organization,but its type is %T, rid: %s", val, rid)
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{key}}
	}

	return errors.RawErrorInfo{}
}

func (attribute *Attribute) validOrganizationValue(org []interface{}, key string, rid string) errors.RawErrorInfo {
	if len(org) == 0 && attribute.IsRequired {
		blog.Errorf("org is required, but is null, rid: %s", rid)
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{key}}
	}

	if len(org) == 0 {
		return errors.RawErrorInfo{}
	}

	if attribute.IsMultiple == nil {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{key}}
	}

	if !(*attribute.IsMultiple) && len(org) != 1 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{key}}
	}

	for _, orgID := range org {
		if !util.IsInteger(orgID) {
			blog.Errorf("orgID params not int, type: %T, rid: %s", orgID, rid)
			return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsIsInvalid, Args: []interface{}{key}}
		}
	}
	return errors.RawErrorInfo{}
}

// validTable valid object attribute that is table type
func (attribute *Attribute) validTable(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {

	rid := util.ExtractRequestIDFromContext(ctx)
	if val == nil {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	if attribute.Option == nil {
		return errors.RawErrorInfo{}
	}

	// validate within enum
	subAttrs, err := ParseSubAttribute(ctx, attribute.Option)
	if err != nil {
		blog.Errorf("parse sub-attribute failed, err: %v, rid: %s", err, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	subAttrMap := make(map[string]SubAttribute)
	for _, subAttr := range subAttrs {
		subAttrMap[subAttr.PropertyID] = subAttr
	}

	if err := attribute.validTableValue(ctx, val, subAttrMap, rid); err != nil {
		blog.Errorf("check value type failed, err: %v, rid: %s", err, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	return errors.RawErrorInfo{}
}

// validTableValue valid object attribute that is table type value
func (attribute *Attribute) validTableValue(ctx context.Context, val interface{}, subAttrMap map[string]SubAttribute,
	rid string) error {

	valMapArr := make([]mapstr.MapStr, 0)
	switch t := val.(type) {
	case []interface{}:
		valMapArr = make([]mapstr.MapStr, len(t))
		for index, value := range t {
			var valMap mapstr.MapStr
			switch v := value.(type) {
			case mapstr.MapStr:
				valMap = v
			case map[string]interface{}:
				valMap = v
			default:
				blog.Errorf("check value type failed, valMap: %#v, rid: %s", valMap, rid)
				return fmt.Errorf("check value type failed, valMap: %v", valMap)
			}
			valMapArr[index] = valMap
		}
	case []mapstr.MapStr:
		valMapArr = t
	case []map[string]interface{}:
		valMapArr = make([]mapstr.MapStr, len(t))
		for index, value := range t {
			valMapArr[index] = value
		}
	default:
		blog.Errorf("check value type failed, val: %#v, rid: %s", val, rid)
		return fmt.Errorf("check value type failed, val: %v", val)
	}

	for _, value := range valMapArr {
		for subKey, subValue := range value {
			validator, exist := subAttrMap[subKey]
			if !exist {
				blog.Errorf("extra field, subKey: %s, subValue: %v, rid: %s", subKey, subValue, rid)
				return fmt.Errorf("extra failed, subKey: %s, subValue: %v", subKey, subValue)
			}
			if rawError := validator.Validate(ctx, subValue, subKey); rawError.ErrCode != 0 {
				blog.Errorf("validate sub-attr failed, key: %s, val: %v, rid: %s", subKey, subValue, rid)
				return fmt.Errorf("validate sub-attr failed, key: %s, val: %v", subKey, subValue)
			}
		}
	}

	return nil
}

func (attribute *Attribute) validInnerTable(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	if val == nil {
		return errors.RawErrorInfo{}
	}

	rid := util.ExtractRequestIDFromContext(ctx)

	switch t := val.(type) {
	case []interface{}:
		if len(t) == 0 {
			if attribute.IsRequired {
				blog.Errorf("required inner table attribute %s value can not be empty, rid: %s", key, rid)
				return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{key}}
			}
			return errors.RawErrorInfo{}
		}

		if len(t) > 50 {
			blog.Errorf("inner table attribute %s value length %d exceeds maximum, rid: %s", key, len(t), rid)
			return errors.RawErrorInfo{ErrCode: common.CCErrArrayLengthWrong, Args: []interface{}{key, 50}}
		}

		_, err := util.SliceInterfaceToInt64(t)
		if err != nil {
			blog.Errorf("parse inner table value to int64 array failed, err: %v, type: %T, rid: %s", err, val, rid)
			return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{key}}
		}
	default:
		blog.Errorf("params should be of inner table type, but its type is %T, rid: %s", val, rid)
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{key}}
	}
	return errors.RawErrorInfo{}
}

func (attribute *Attribute) validIDRule(ctx context.Context, val interface{}, key string) errors.RawErrorInfo {
	if val == nil {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{key}}
	}

	rid := util.ExtractRequestIDFromContext(ctx)

	switch t := val.(type) {
	case string:
		if len(t) > common.FieldTypeSingleLenChar {
			blog.Errorf("params over length %d, rid: %s", common.FieldTypeSingleLenChar, rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommOverLimit,
				Args:    []interface{}{key},
			}
		}
	default:
		blog.Errorf("params should be of string type, but its type is %T, rid: %s", val, rid)
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{key}}
	}

	return errors.RawErrorInfo{}
}

// ValidIDRuleVal validate id rule value
func ValidIDRuleVal(ctx context.Context, inst mapstr.MapStr, field Attribute, attrMap map[string]Attribute) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	val, err := inst.String(field.PropertyID)
	if err != nil {
		blog.Errorf("get property: %s failed, inst: %+v, err: %v, rid: %s", field.PropertyID, inst, err, rid)
		return err
	}

	rules, err := ParseSubIDRules(field.Option)
	if err != nil {
		blog.Errorf("parse sub id rule failed, field: %+v, err: %v, rid: %s", field, err, rid)
		return err
	}

	for _, rule := range rules {
		if rule.Kind != Attr && int64(len(val)) < rule.Len {
			blog.Errorf("val is invalid, val: %s, rule len: %d, rid: %s", val, rule.Len, rid)
			return fmt.Errorf("val is invalid, val: %s, rule len: %d", val, rule.Len)
		}

		switch rule.Kind {
		case Const:
			if prefix := val[:rule.Len]; prefix != rule.Val {
				blog.Errorf("val is invalid, val: %s, rule val: %s, rid: %s", prefix, rule.Val, rid)
				return fmt.Errorf("val is invalid, val: %s, rule val: %s", prefix, rule.Val)
			}

		case Attr:
			attr, exists := attrMap[rule.Val]
			if !exists {
				blog.Errorf("attr val %s is invalid, attribute not exists, rid: %s", rule.Val, rid)
				return fmt.Errorf("val %s related attr not exists", rule.Val)
			}
			if !IsValidAttrRuleType(attr.PropertyType) {
				blog.Errorf("attr val %s type %s is invalid, rid: %s", rule.Val, attr.PropertyType, rid)
				return fmt.Errorf("attr val %s type %s is invalid", rule.Val, attr.PropertyType)
			}

			refVal, err := inst.String(rule.Val)
			if err != nil {
				blog.Errorf("get property:%s failed, inst: %+v, err: %v, rid: %s", rule.Val, inst, err, rid)
				return err
			}
			rule.Len = int64(len(refVal))

			if int64(len(val)) < rule.Len {
				blog.Errorf("val is invalid, val: %s, rule len: %d, rid: %s", val, rule.Len, rid)
				return fmt.Errorf("val is invalid, val: %s, rule len: %d", val, rule.Len)
			}

			if prefix := val[:len(refVal)]; prefix != refVal {
				blog.Errorf("val is invalid, val: %s, rule val: %s, rid: %s", prefix, refVal, rid)
				return fmt.Errorf("val is invalid, val: %s, rule val: %s", prefix, refVal)
			}

		case GlobalID, LocalID, RandomID:
			prefix := val[:rule.Len]
			for _, c := range prefix {
				if !unicode.IsDigit(c) {
					blog.Errorf("the char value needs to be a number, val: %c, rid: %s", c, rid)
					return fmt.Errorf("the char value needs to be a number, val: %c", c)
				}
			}

		default:
			blog.Errorf("option is invalid, val: %+v, rid: %s", field.Option, rid)
			return fmt.Errorf("option is invalid, val: %+v", field.Option)
		}

		val = val[rule.Len:]
	}

	if val != "" {
		blog.Errorf("val is invalid, val: %s, rid: %s", val, rid)
		return fmt.Errorf("val is invalid")
	}

	return nil
}

// PrevIntOption previous integer option
// Deprecated: do not use anymore, use IntOption instead.
type PrevIntOption struct {
	Min string `bson:"min" json:"min"`
	Max string `bson:"max" json:"max"`
}

// IntOption integer option
type IntOption struct {
	Min int64 `bson:"min" json:"min"`
	Max int64 `bson:"max" json:"max"`
}

// ParseIntOption parse int data in option
func ParseIntOption(val interface{}) (IntOption, error) {
	if val == nil || val == "" {
		return IntOption{Max: common.MaxInt64, Min: common.MinInt64}, nil
	}

	var optMap map[string]interface{}

	switch option := val.(type) {
	case IntOption:
		return option, nil
	case string:
		return parseIntOptionMaxMin(gjson.Get(option, "max").Raw, gjson.Get(option, "min").Raw)
	case map[string]interface{}:
		optMap = option
	case bson.M:
		optMap = option
	case bson.D:
		optMap = option.Map()
	default:
		return IntOption{}, fmt.Errorf("unknow val type: %T", val)
	}

	return parseIntOptionMaxMin(optMap["max"], optMap["min"])
}

func parseIntOptionMaxMin(maxVal, minVal interface{}) (IntOption, error) {
	max, err := parseIntOptValue(maxVal, common.MaxInt64)
	if err != nil {
		return IntOption{}, fmt.Errorf("parse max int option %+v failed, err: %v", maxVal, err)
	}

	min, err := parseIntOptValue(minVal, common.MinInt64)
	if err != nil {
		return IntOption{}, fmt.Errorf("parse min int option %+v failed, err: %v", minVal, err)
	}

	return IntOption{Max: max, Min: min}, nil
}

func parseIntOptValue(value interface{}, defaultVal int64) (int64, error) {
	switch val := value.(type) {
	case string:
		if len(val) == 0 || val == `""` {
			return defaultVal, nil
		}
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, err
		}
		return intVal, nil
	default:
		intVal, err := util.GetInt64ByInterface(val)
		if err != nil {
			return 0, err
		}
		return intVal, nil
	}
}

// FloatOption float option
type FloatOption struct {
	Min float64 `bson:"min" json:"min"`
	Max float64 `bson:"max" json:"max"`
}

// ParseFloatOption parse float data in option
func ParseFloatOption(val interface{}) (FloatOption, error) {
	if val == nil || val == "" {
		return FloatOption{Max: float64(common.MaxInt64), Min: float64(common.MinInt64)}, nil
	}

	var optMap map[string]interface{}

	switch option := val.(type) {
	case FloatOption:
		return option, nil
	case string:
		return parseFloatOptionMaxMin(gjson.Get(option, "max").Raw, gjson.Get(option, "min").Raw)
	case map[string]interface{}:
		optMap = option
	case bson.M:
		optMap = option
	case bson.D:
		optMap = option.Map()
	default:
		return FloatOption{}, fmt.Errorf("unknow val type: %T", val)
	}

	return parseFloatOptionMaxMin(optMap["max"], optMap["min"])
}

func parseFloatOptionMaxMin(maxVal, minVal interface{}) (FloatOption, error) {
	max, err := parseFloatOptValue(maxVal, float64(common.MaxInt64))
	if err != nil {
		return FloatOption{}, fmt.Errorf("parse max float option %+v failed, err: %v", maxVal, err)
	}

	min, err := parseFloatOptValue(minVal, float64(common.MinInt64))
	if err != nil {
		return FloatOption{}, fmt.Errorf("parse min float option %+v failed, err: %v", minVal, err)
	}

	return FloatOption{Max: max, Min: min}, nil
}

func parseFloatOptValue(value interface{}, defaultVal float64) (float64, error) {
	switch val := value.(type) {
	case string:
		if len(val) == 0 || val == `""` {
			return defaultVal, nil
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, err
		}
		return floatVal, nil
	default:
		floatVal, err := util.GetFloat64ByInterface(val)
		if err != nil {
			return 0, err
		}
		return floatVal, nil
	}
}

func getString(val interface{}) string {
	if val == nil {
		return ""
	}

	if ret, ok := val.(string); ok {
		return ret
	}
	return ""
}

func getBool(val interface{}) bool {
	if val == nil {
		return false
	}
	if ret, ok := val.(bool); ok {
		return ret
	}
	return false
}

// EnumOption enum option
type EnumOption []EnumVal

// EnumVal enum option val
type EnumVal struct {
	ID        string `bson:"id"           json:"id"`
	Name      string `bson:"name"         json:"name"`
	Type      string `bson:"type"         json:"type"`
	IsDefault bool   `bson:"is_default"   json:"is_default"`
}

// ParseEnumOption convert val to []EnumVal
func ParseEnumOption(val interface{}) (EnumOption, error) {
	enumOptions := make([]EnumVal, 0)
	if val == nil || val == "" {
		return enumOptions, nil
	}

	var optionArr []interface{}

	switch options := val.(type) {
	case EnumOption:
		return options, nil
	case []EnumVal:
		return options, nil
	case string:
		err := json.Unmarshal([]byte(options), &enumOptions)
		if err != nil {
			return nil, err
		}
	case []interface{}:
		optionArr = options
	case bson.A:
		optionArr = options
	default:
		return nil, fmt.Errorf("unknow val type: %T for enum option", val)
	}

	for _, optionElem := range optionArr {
		enumVal, err := parseEnumVal(optionElem)
		if err != nil {
			return nil, err
		}
		enumOptions = append(enumOptions, enumVal)
	}

	return enumOptions, nil
}

// parseEnumVal parse enum options element value
func parseEnumVal(val interface{}) (EnumVal, error) {
	var valMap mapstr.MapStr

	switch optionVal := val.(type) {
	case map[string]interface{}:
		valMap = optionVal
	case bson.M:
		valMap = mapstr.MapStr(optionVal)
	case bson.D:
		valMap = mapstr.MapStr(optionVal.Map())
	default:
		return EnumVal{}, fmt.Errorf("unknow element type: %T for enum option", val)
	}

	if valMap == nil {
		return EnumVal{}, fmt.Errorf("enum option val map is nil")
	}

	enumOption := EnumVal{}
	enumOption.ID = getString(valMap["id"])
	enumOption.Name = getString(valMap["name"])
	enumOption.Type = getString(valMap["type"])
	enumOption.IsDefault = getBool(valMap["is_default"])
	if enumOption.ID == "" || enumOption.Name == "" || enumOption.Type != "text" {
		return EnumVal{}, fmt.Errorf("enum option val %#v id, name empty or not string, or type not text", val)
	}

	return enumOption, nil
}

// EnumQuoteVal enum quote option val
type EnumQuoteVal struct {
	ObjID  string `bson:"bk_obj_id" json:"bk_obj_id"`
	InstID int64  `bson:"bk_inst_id" json:"bk_inst_id"`
}

// ParseEnumQuoteOption convert val to []EnumQuoteVal
func ParseEnumQuoteOption(ctx context.Context, val interface{}) ([]EnumQuoteVal, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	enumQuoteOptions := []EnumQuoteVal{}
	if val == nil || val == "" {
		return enumQuoteOptions, nil
	}
	switch options := val.(type) {
	case []EnumQuoteVal:
		return options, nil
	case string:
		err := json.Unmarshal([]byte(options), &enumQuoteOptions)
		if nil != err {
			blog.Errorf("parse enum quote option failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
	case []interface{}:
		if err := parseEnumQuoteOption(options, &enumQuoteOptions); err != nil {
			blog.Errorf("parse enum quote option failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
	case primitive.A:
		if err := parseEnumQuoteOption(options, &enumQuoteOptions); err != nil {
			blog.Errorf("parse enum quote option failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknow val type: %#v", val)
	}
	return enumQuoteOptions, nil
}

// getEnumQuoteOptions get enum quote option value
func getEnumQuoteOptions(val map[string]interface{}, enumQuoteOptions *[]EnumQuoteVal) error {
	enumQuoteOption := EnumQuoteVal{}
	enumQuoteOption.ObjID = getString(val[common.BKObjIDField])
	if enumQuoteOption.ObjID == "" {
		return fmt.Errorf("operation %#v objID empty or not string", val)
	}
	instID, err := util.GetInt64ByInterface(val[common.BKInstIDField])
	if err != nil {
		return err
	}
	if instID == 0 {
		return fmt.Errorf("inst id cannot be 0")
	}
	enumQuoteOption.InstID = instID
	*enumQuoteOptions = append(*enumQuoteOptions, enumQuoteOption)
	return nil
}

// parseEnumQuoteOption set enum quote Options values from options
func parseEnumQuoteOption(options []interface{}, enumQuoteOptions *[]EnumQuoteVal) error {
	for _, optionVal := range options {
		switch val := optionVal.(type) {
		case map[string]interface{}:
			if err := getEnumQuoteOptions(val, enumQuoteOptions); err != nil {
				return err
			}
		case bson.M:
			if err := getEnumQuoteOptions(val, enumQuoteOptions); err != nil {
				return err
			}
		case bson.D:
			opt := val.Map()
			if err := getEnumQuoteOptions(opt, enumQuoteOptions); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknow optionVal type: %#v", optionVal)
		}
	}
	return nil
}

// ListOption list option
type ListOption []string

// ParseListOption parse 'list' type option
func ParseListOption(option interface{}) (ListOption, error) {
	if option == nil {
		return ListOption{}, fmt.Errorf("list type field option is null")
	}

	var arrOption []interface{}
	switch optionVal := option.(type) {
	case []interface{}:
		arrOption = optionVal
	case primitive.A:
		arrOption = optionVal
	case ListOption:
		return optionVal, nil
	default:
		return nil, fmt.Errorf("list option %+v type %T is invalid", option, option)
	}

	if len(arrOption) == 0 {
		return ListOption{}, fmt.Errorf("list type field option is empty")
	}

	valueList := make(ListOption, len(arrOption))
	for _, val := range arrOption {
		strVal, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("list option element %+v type %T is invalid", val, val)
		}

		valueList = append(valueList, strVal)
	}

	return valueList, nil
}

// PrettyValue TODO
func (attribute Attribute) PrettyValue(ctx context.Context, val interface{}) (string, error) {
	if val == nil {
		return "", nil
	}

	fieldType := attribute.PropertyType
	switch fieldType {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar:
		value, ok := val.(string)
		if ok == false {
			return "", fmt.Errorf("invalid value type for %s, value: %+v", fieldType, val)
		}
		return value, nil
	case common.FieldTypeInt:
		var value int64
		value, err := util.GetInt64ByInterface(val)
		if nil != err {
			return "", fmt.Errorf("invalid value type for %s, value: %+v, err: %+v", fieldType, val, err)
		}
		return strconv.FormatInt(value, 10), nil
	case common.FieldTypeFloat:
		var value float64
		value, err := util.GetFloat64ByInterface(value)
		if nil != err {
			return "", fmt.Errorf("invalid value type for %s, value: %+v, err: %+v", fieldType, value, err)
		}
		return strconv.FormatFloat(value, 'E', -1, 64), nil
	case common.FieldTypeEnum:
		valStr, ok := val.(string)
		if !ok {
			return "", fmt.Errorf("invalid value type for %s, value: %+v", fieldType, val)
		}
		// validate within enum
		enumOption, err := ParseEnumOption(attribute.Option)
		if err != nil {
			return "", fmt.Errorf("parse options for enum type failed, err: %+v", err)
		}
		for _, k := range enumOption {
			if k.ID == valStr {
				return k.Name, nil
			}
		}
		return "", fmt.Errorf("invalid value for %s, value: %s", fieldType, valStr)
	case common.FieldTypeDate:
		valStr, ok := val.(string)
		if ok == false {
			return "", fmt.Errorf("invalid data type for %s, value: %+v", fieldType, val)
		}
		return valStr, nil
	case common.FieldTypeTime:
		valStr, ok := val.(string)
		if ok == false {
			return "", fmt.Errorf("invalid value type for %s, value: %+v", fieldType, val)
		}
		return valStr, nil
	case common.FieldTypeTimeZone:
		switch value := val.(type) {
		case string:
			return value, nil
		default:
			return "", fmt.Errorf("invalid value type for %s, value: %+v", fieldType, val)
		}
	case common.FieldTypeBool:
		value, ok := val.(bool)
		if ok == false {
			return "", fmt.Errorf("invalid value type for %s, value: %+v", fieldType, val)
		}
		return strconv.FormatBool(value), nil
	case common.FieldTypeUser:
		value, ok := val.(string)
		if ok == false {
			return "", fmt.Errorf("invalid value type for %s, value: %+v", fieldType, val)
		}
		return value, nil
	case common.FieldTypeList:
		strVal, ok := val.(string)
		if !ok {
			return "", fmt.Errorf("invalid value type for %s, value: %+v", fieldType, val)
		}

		listOption, ok := attribute.Option.([]interface{})
		if false == ok {
			return "", fmt.Errorf("parse options for list type failed, option not slice type, option: %+v",
				attribute.Option)
		}
		for _, inVal := range listOption {
			inValStr, ok := inVal.(string)
			if !ok {
				return "", fmt.Errorf("parse list option failed, item not string, item: %+v", inVal)
			}
			if strVal == inValStr {
				return strVal, nil
			}
		}
		return "", fmt.Errorf("invalid value for list, value: %s, options: %+v", strVal, listOption)
	default:
		blog.V(3).Infof("unexpected property type: %s", fieldType)
		return fmt.Sprintf("%#v", val), nil
	}
}

// ValidTableDefaultAttr judging the legitimacy of the basic type in the table field.
func (attribute *Attribute) ValidTableDefaultAttr(ctx context.Context, val interface{}) errors.RawErrorInfo {
	rid := util.ExtractRequestIDFromContext(ctx)
	if attribute == nil {
		blog.Errorf("the key of the default value is illegal and not in the header list, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"default key"},
		}
	}

	switch attribute.PropertyType {
	case common.FieldTypeInt:
		return attribute.validInt(ctx, val, attribute.PropertyID)
	case common.FieldTypeFloat:
		return attribute.validFloat(ctx, val, attribute.PropertyID)
	case common.FieldTypeSingleChar:
		return attribute.validChar(ctx, val, attribute.PropertyID)
	case common.FieldTypeLongChar:
		return attribute.validLongChar(ctx, val, attribute.PropertyID)
	case common.FieldTypeEnumMulti:
		return attribute.validEnumMulti(ctx, val, attribute.PropertyID)
	case common.FieldTypeBool:
		return attribute.validBool(ctx, val, attribute.PropertyID)
	default:
		blog.Errorf("type error, type: %s, rid: %s", attribute.PropertyType, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{attribute.PropertyType},
		}
	}
}

// ParseTableAttrOption parse table attribute options.
func ParseTableAttrOption(option interface{}) (*TableAttributesOption, error) {
	marshaledOptions, err := json.Marshal(option)
	if err != nil {
		return nil, err
	}

	result := new(TableAttributesOption)
	if err := json.Unmarshal(marshaledOptions, result); err != nil {
		return nil, err
	}
	return result, nil
}

// HostApplyFieldMap TODO
var HostApplyFieldMap = map[string]bool{
	common.BKOperatorField:        true,
	common.BKBakOperatorField:     true,
	"bk_state":                    true,
	"bk_sla":                      true,
	common.BKHostInnerIPField:     false,
	common.BKHostOuterIPField:     false,
	common.BKAssetIDField:         false,
	common.BKSNField:              false,
	"bk_comment":                  false,
	"bk_service_term":             false,
	common.BKCloudIDField:         false,
	"bk_state_name":               false,
	"bk_province_name":            false,
	"bk_isp_name":                 false,
	common.BKHostNameField:        false,
	common.BKOSTypeField:          false,
	common.BKOSNameField:          false,
	"bk_os_version":               false,
	"bk_os_bit":                   false,
	"bk_cpu":                      false,
	"bk_cpu_module":               false,
	"bk_mem":                      false,
	"bk_disk":                     false,
	"bk_mac":                      false,
	"bk_outer_mac":                false,
	common.CreateTimeField:        false,
	common.LastTimeField:          false,
	common.BKImportFrom:           false,
	common.BKCloudInstIDField:     false,
	common.BKCloudHostStatusField: false,
	common.BKCloudVendor:          false,
	common.BKHostInnerIPv6Field:   false,
	common.BKHostOuterIPv6Field:   false,
	common.BKAgentIDField:         false,
	"bk_cpu_architecture":         false,
}

// CheckAllowHostApplyOnField 检查字段是否能用于主机属性自动应用
func CheckAllowHostApplyOnField(field *Attribute) bool {
	if !field.IsEditable {
		return false
	}
	// 屏蔽表格字段
	if field.PropertyType == common.FieldTypeInnerTable {
		return false
	}
	if allow, exist := HostApplyFieldMap[field.PropertyID]; exist == true {
		return allow
	}
	return true
}

// SubAttributeOption TODO
type SubAttributeOption []SubAttribute

// SubAttribute sub attribute metadata definition
type SubAttribute struct {
	PropertyID    string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
	PropertyName  string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
	Placeholder   string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
	IsEditable    bool        `field:"editable" json:"editable" bson:"editable"`
	IsRequired    bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
	IsReadOnly    bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
	IsSystem      bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
	IsAPI         bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
	PropertyType  string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
	Option        interface{} `field:"option" json:"option" bson:"option"`
	Description   string      `field:"description" json:"description" bson:"description"`
	PropertyGroup string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
}

// Validate TODO
func (sa *SubAttribute) Validate(ctx context.Context, data interface{}, key string) errors.RawErrorInfo {
	attr := Attribute{
		PropertyID:   sa.PropertyID,
		PropertyName: sa.PropertyName,
		Placeholder:  sa.Placeholder,

		IsEditable:   sa.IsEditable,
		IsRequired:   sa.IsRequired,
		IsReadOnly:   sa.IsReadOnly,
		IsSystem:     sa.IsSystem,
		IsAPI:        sa.IsAPI,
		PropertyType: sa.PropertyType,
		Option:       sa.Option,
		Description:  sa.Description,

		PropertyGroup: sa.PropertyGroup,
	}
	return attr.Validate(ctx, data, key)
}

// ParseSubAttribute convert val to []SubAttribute
func ParseSubAttribute(ctx context.Context, val interface{}) (SubAttributeOption, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	subAttrs := make([]SubAttribute, 0)
	var err error
	if val == nil || val == "" {
		return subAttrs, nil
	}
	switch options := val.(type) {
	case []SubAttribute:
		return options, nil
	case string:
		err = json.Unmarshal([]byte(options), &subAttrs)
		if nil != err {
			blog.Errorf("parse sub-attribute failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
	case []interface{}:
		subAttrs, err = parseSubAttribute(options)
		if err != nil {
			blog.Errorf("parse sub-attribute failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
	case bson.A:
		subAttrs, err = parseSubAttribute(options)
		if err != nil {
			blog.Errorf("parse sub-attribute failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknow val type: %#v", val)
	}
	return subAttrs, nil
}

func parseSubAttribute(options []interface{}) ([]SubAttribute, error) {

	subAttrs := make([]SubAttribute, 0)
	for _, optionVal := range options {
		switch option := optionVal.(type) {
		case map[string]interface{}:
			subAttrs = append(subAttrs, parseSubAttr(option))
		case bson.M:
			subAttrs = append(subAttrs, parseSubAttr(map[string]interface{}(option)))
		case bson.D:
			subAttrs = append(subAttrs, parseSubAttr(map[string]interface{}(option.Map())))
		default:
			return nil, fmt.Errorf("unknow optionVal type: %#v", optionVal)
		}
	}
	return subAttrs, nil
}

func parseSubAttr(options map[string]interface{}) SubAttribute {
	subAttr := SubAttribute{}

	subAttr.PropertyID = getString(options["bk_property_id"])
	subAttr.PropertyName = getString(options["bk_property_name"])
	subAttr.PropertyGroup = getString(options["bk_property_group"])
	subAttr.Placeholder = getString(options["placeholder"])
	subAttr.PropertyType = getString(options["bk_property_type"])
	subAttr.IsAPI = getBool(options["bk_isapi"])
	subAttr.IsEditable = getBool(options["editable"])
	subAttr.IsReadOnly = getBool(options["isreadonly"])
	subAttr.IsRequired = getBool(options["isrequired"])
	subAttr.IsSystem = getBool(options["bk_issystem"])
	subAttr.Option = options["option"]
	subAttr.Description = getString(options["description"])

	return subAttr
}

// EnumOptions TODO
type EnumOptions []AttributesOption

// AttributesOption TODO
type AttributesOption struct {
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	Type      string `json:"type" bson:"type"`
	IsDefault bool   `json:"is_default" bson:"is_default"`
}

// ListOptions TODO
type ListOptions []string

// TableAttributesOption the option of the form field, including the header and the default value.
type TableAttributesOption struct {
	Header  []Attribute              `json:"header" bson:"header" mapstructure:"header"`
	Default []map[string]interface{} `json:"default" bson:"default" mapstructure:"default"`
}

// ValidTableFieldBaseType determine the basic type supported by the form field.
func ValidTableFieldBaseType(fieldType string) bool {

	switch fieldType {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar, common.FieldTypeBool,
		common.FieldTypeEnumMulti, common.FieldTypeFloat, common.FieldTypeInt:
		return true
	default:
		return false
	}
}
