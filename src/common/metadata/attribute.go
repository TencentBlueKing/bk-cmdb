/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	AttributeFieldID              = "id"
	AttributeFieldSupplierAccount = "bk_supplier_account"
	AttributeFieldObjectID        = "bk_obj_id"
	AttributeFieldPropertyID      = "bk_property_id"
	AttributeFieldPropertyName    = "bk_property_name"
	AttributeFieldPropertyGroup   = "bk_property_group"
	AttributeFieldPropertyIndex   = "bk_property_index"
	AttributeFieldUnit            = "unit"
	AttributeFieldPlaceHolder     = "placeholder"
	AttributeFieldIsEditable      = "editable"
	AttributeFieldIsPre           = "ispre"
	AttributeFieldIsRequired      = "isrequired"
	AttributeFieldIsReadOnly      = "isreadonly"
	AttributeFieldIsOnly          = "isonly"
	AttributeFieldIsSystem        = "bk_issystem"
	AttributeFieldIsAPI           = "bk_isapi"
	AttributeFieldPropertyType    = "bk_property_type"
	AttributeFieldOption          = "option"
	AttributeFieldDescription     = "description"
	AttributeFieldCreator         = "creator"
	AttributeFieldCreateTime      = "create_time"
	AttributeFieldLastTime        = "last_time"
)

// Attribute attribute metadata definition
type Attribute struct {
	BizID             int64       `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID                int64       `field:"id" json:"id" bson:"id"`
	OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
	Unit              string      `field:"unit" json:"unit" bson:"unit"`
	Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
	IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
	IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
	IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
	IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{} `field:"option" json:"option" bson:"option"`
	Description       string      `field:"description" json:"description" bson:"description"`
	Creator           string      `field:"creator" json:"creator" bson:"creator"`
	CreateTime        *Time       `json:"create_time" bson:"create_time"`
	LastTime          *Time       `json:"last_time" bson:"last_time"`
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

type HostObjAttDes struct {
	ObjAttDes        `json:",inline" bson:",inline"`
	HostApplyEnabled bool `json:"host_apply_enabled"`
}

func (attribute *Attribute) Validate(ctx context.Context, data interface{}, key string) (rawError errors.RawErrorInfo) {
	fieldType := attribute.PropertyType
	switch fieldType {
	case common.FieldTypeSingleChar:
		rawError = attribute.validChar(ctx, data, key)
	case common.FieldTypeLongChar:
		rawError = attribute.validLongChar(ctx, data, key)
	case common.FieldTypeInt:
		rawError = attribute.validInt(ctx, data, key)
	case common.FieldTypeFloat:
		rawError = attribute.validFloat(ctx, data, key)
	case common.FieldTypeEnum:
		rawError = attribute.validEnum(ctx, data, key)
	case common.FieldTypeDate:
		rawError = attribute.validDate(ctx, data, key)
	case common.FieldTypeTime:
		rawError = attribute.validTime(ctx, data, key)
	case common.FieldTypeTimeZone:
		rawError = attribute.validTimeZone(ctx, data, key)
	case common.FieldTypeBool:
		rawError = attribute.validBool(ctx, data, key)
	case common.FieldTypeUser:
		rawError = attribute.validUser(ctx, data, key)
	case common.FieldTypeList:
		rawError = attribute.validList(ctx, data, key)
	case common.FieldTypeOrganization:
		rawError = attribute.validOrganization(ctx, data, key)
	case "foreignkey", "singleasst", "multiasst":
		// TODO what validation should do on these types
	case common.FieldTypeTable:
		// TODO what validation should do on these types
		rawError = attribute.validTable(ctx, data, key)
	default:
		rawError = errors.RawErrorInfo{
			ErrCode: common.CCErrCommUnexpectedFieldType,
			Args:    []interface{}{fieldType},
		}
	}
	return rawError
}

// validTime valid object Attribute that is time type
func (attribute *Attribute) validTime(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {

	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val || "" == val {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	valStr, ok := val.(string)
	if false == ok {
		blog.Errorf("date can should be string, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsShouldBeString,
			Args:    []interface{}{key},
		}
	}

	result := util.IsTime(valStr)
	if !result {
		blog.Errorf("params not valid, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validDate valid object Attribute that is date type
func (attribute *Attribute) validDate(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val || "" == val {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}
	valStr, ok := val.(string)
	if false == ok {
		blog.Errorf("date can should be string, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsShouldBeString,
			Args:    []interface{}{key},
		}

	}
	result := util.IsDate(valStr)
	if !result {
		blog.Errorf("params is not valid, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validEnum valid object attribute that is enum type
func (attribute *Attribute) validEnum(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	// validate require
	if nil == val {
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
	enumOption, err := ParseEnumOption(ctx, attribute.Option)
	if err != nil {
		blog.Warnf("ParseEnumOption failed: %v, rid: %s", err, rid)
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
	blog.V(3).Infof("params %s not valid, option %#v, raw option %#v, value: %#v, rid: %s", key, enumOption, attribute.Option, val, rid)
	blog.Errorf("params %s not valid , enum value: %#v, rid: %s", key, val, rid)
	return errors.RawErrorInfo{
		ErrCode: common.CCErrCommParamsInvalid,
		Args:    []interface{}{key},
	}
}

// validBool valid object attribute that is bool type
func (attribute *Attribute) validBool(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
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

// valid char valid object attribute that is timezone type
func (attribute *Attribute) validTimeZone(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
		if attribute.IsRequired {
			blog.Errorf("params can not be null, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{key},
			}

		}
		return errors.RawErrorInfo{}
	}

	switch value := val.(type) {
	case string:
		isMatch := util.IsTimeZone(value)
		if false == isMatch {
			blog.Errorf("params should be timezone, rid: %s", rid)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedTimeZone,
				Args:    []interface{}{key},
			}
		}
	default:
		blog.Errorf("params should be timezone, rid: %s", rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedTimeZone,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validInt valid object attribute that is int type
func (attribute *Attribute) validInt(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
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

	value, _ := util.GetInt64ByInterface(val)

	intObjOption := ParseIntOption(ctx, attribute.Option)
	if 0 == len(intObjOption.Min) || 0 == len(intObjOption.Max) {
		return errors.RawErrorInfo{}
	}

	maxValue, err := strconv.ParseInt(intObjOption.Max, 10, 64)
	if nil != err {
		maxValue = common.MaxInt64
	}
	minValue, err := strconv.ParseInt(intObjOption.Min, 10, 64)
	if nil != err {
		minValue = common.MinInt64
	}
	if value > maxValue || value < minValue {
		blog.Errorf("params %s:%#v not valid, rid: %s", key, val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validFloat valid object attribute that is float type
func (attribute *Attribute) validFloat(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
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
	if nil != err {
		blog.Errorf("params %s:%#v not float, rid: %s", key, val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{key},
		}
	}

	intObjOption := parseFloatOption(ctx, attribute.Option)
	if 0 == len(intObjOption.Min) || 0 == len(intObjOption.Max) {
		return errors.RawErrorInfo{}
	}

	maxValue, err := strconv.ParseFloat(intObjOption.Max, 64)
	if nil != err {
		maxValue = float64(common.MaxInt64)
	}
	minValue, err := strconv.ParseFloat(intObjOption.Min, 64)
	if nil != err {
		minValue = float64(common.MinInt64)
	}
	if value > maxValue || value < minValue {
		blog.Errorf("params %s:%#v not valid, rid: %s", key, val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validLongChar valid object attribute that is long char type
func (attribute *Attribute) validLongChar(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val || "" == val {
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
		if 0 == len(value) {
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
		if nil != err {
			blog.Errorf(`regexp "%s" invalid, err: %s, rid:  %s`, option, err.Error(), rid)
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
func (attribute *Attribute) validChar(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
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
		if 0 == len(value) {
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
				blog.ErrorJSON("params %s contains TopoModuleName's split flag ##, rid: %s", value, rid)
				return errors.RawErrorInfo{
					ErrCode: common.CCErrCommParamsInvalid,
					Args:    []interface{}{value},
				}
			}
		}

		if "" == val {
			return errors.RawErrorInfo{}
		}

		option, ok := attribute.Option.(string)
		if !ok {
			break
		}
		strReg, err := regexp.Compile(option)
		if nil != err {
			blog.Errorf(`regexp "%s" invalid, err: %s, rid:  %s`, option, err.Error(), rid)
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
func (attribute *Attribute) validUser(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val || "" == val {
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

		if 0 == len(value) {
			if attribute.IsRequired {
				blog.Errorf("params can not be empty, rid: %s", rid)
				return errors.RawErrorInfo{
					ErrCode: common.CCErrCommParamsNeedSet,
					Args:    []interface{}{key},
				}
			}
			return errors.RawErrorInfo{}
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

func (attribute *Attribute) validList(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestUserFromContext(ctx)

	if nil == val {
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
		blog.Errorf("option %v invalid, not string type list option, but type %T", attribute.Option, attribute.Option)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}

	for _, inVal := range listOpt {
		inValStr, ok := inVal.(string)
		if !ok {
			blog.Errorf("inner list option convert to string  failed, params %s not valid , list field value: %#v", key, val)
			return errors.RawErrorInfo{
				ErrCode: common.CCErrParseAttrOptionListFailed,
				Args:    []interface{}{key},
			}
		}
		if strVal == inValStr {
			return errors.RawErrorInfo{}
		}
	}
	blog.Errorf("params %s not valid, option %#v, raw option %#v, value: %#v", key, listOpt, attribute, val)
	return errors.RawErrorInfo{
		ErrCode: common.CCErrCommParamsInvalid,
		Args:    []interface{}{key},
	}
}

// validBool valid object attribute that is bool type
func (attribute *Attribute) validOrganization(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
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
		blog.Errorf("params should be type organization,but its type is %T, rid: %s", val, rid)
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{key},
		}
	}
	return errors.RawErrorInfo{}
}

// validTable valid object attribute that is bool type
func (attribute *Attribute) validTable(ctx context.Context, val interface{}, key string) (rawError errors.RawErrorInfo) {
	// rid := util.ExtractRequestIDFromContext(ctx)
	// TODO 暂时不需要实现，目前只有进程和进程模板使用
	return errors.RawErrorInfo{}
}

// parseFloatOption  parse float data in option
func parseFloatOption(ctx context.Context, val interface{}) FloatOption {
	rid := util.ExtractRequestIDFromContext(ctx)
	floatOption := FloatOption{}
	if nil == val || "" == val {
		return floatOption
	}
	switch option := val.(type) {
	case string:
		floatOption.Min = gjson.Get(option, "min").Raw
		floatOption.Max = gjson.Get(option, "max").Raw
	case map[string]interface{}:
		floatOption.Min = getString(option["min"])
		floatOption.Max = getString(option["max"])
	case bson.M:
		floatOption.Min = getString(option["min"])
		floatOption.Max = getString(option["max"])
	case bson.D:
		opt := option.Map()
		floatOption.Min = getString(opt["min"])
		floatOption.Max = getString(opt["max"])
	default:
		blog.Warnf("unknow val type: %#v, rid: %s", val, rid)
	}
	return floatOption
}

// parseIntOption  parse int data in option
func ParseIntOption(ctx context.Context, val interface{}) IntOption {
	rid := util.ExtractRequestIDFromContext(ctx)
	intOption := IntOption{}
	if nil == val || "" == val {
		return intOption
	}
	switch option := val.(type) {
	case string:
		intOption.Min = gjson.Get(option, "min").Raw
		intOption.Max = gjson.Get(option, "max").Raw
	case map[string]interface{}:
		intOption.Min = getString(option["min"])
		intOption.Max = getString(option["max"])
	case bson.M:
		intOption.Min = getString(option["min"])
		intOption.Max = getString(option["max"])
	case bson.D:
		opt := option.Map()
		intOption.Min = getString(opt["min"])
		intOption.Max = getString(opt["max"])
	default:
		blog.Warnf("unknow val type: %#v, rid: %s", val, rid)
	}
	return intOption
}

// EnumOption enum option
type EnumOption []EnumVal

// IntOption integer option
type IntOption struct {
	Min string `bson:"min" json:"min"`
	Max string `bson:"max" json:"max"`
}

// FloatOption float option
type FloatOption struct {
	Min string `bson:"min" json:"min"`
	Max string `bson:"max" json:"max"`
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

// GetDefault returns EnumOption's default value
func (opt EnumOption) GetDefault() *EnumVal {
	for index := range opt {
		if opt[index].IsDefault {
			return &opt[index]
		}
	}
	return nil
}

// EnumVal enum option val
type EnumVal struct {
	ID        string `bson:"id"           json:"id"`
	Name      string `bson:"name"         json:"name"`
	Type      string `bson:"type"         json:"type"`
	IsDefault bool   `bson:"is_default"   json:"is_default"`
}

// ParseEnumOption convert val to []EnumVal
func ParseEnumOption(ctx context.Context, val interface{}) (EnumOption, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	enumOptions := []EnumVal{}
	if nil == val || "" == val {
		return enumOptions, nil
	}
	switch options := val.(type) {
	case []EnumVal:
		return options, nil
	case string:
		err := json.Unmarshal([]byte(options), &enumOptions)
		if nil != err {
			blog.Errorf("ParseEnumOption error : %s, rid: %s", err.Error(), rid)
			return nil, err
		}
	case []interface{}:
		if err := parseEnumOption(options, &enumOptions); err != nil {
			blog.Errorf("parseEnumOption error : %s, rid: %s", err.Error(), rid)
			return nil, err
		}
	case bson.A:
		if err := parseEnumOption(options, &enumOptions); err != nil {
			blog.Errorf("parseEnumOption error : %s, rid: %s", err.Error(), rid)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknow val type: %#v", val)
	}
	return enumOptions, nil
}

// parseEnumOption set enumOptions values from options
func parseEnumOption(options []interface{}, enumOptions *[]EnumVal) error {
	for _, optionVal := range options {
		if option, ok := optionVal.(map[string]interface{}); ok {
			enumOption := EnumVal{}
			enumOption.ID = getString(option["id"])
			enumOption.Name = getString(option["name"])
			enumOption.Type = getString(option["type"])
			enumOption.IsDefault = getBool(option["is_default"])
			if enumOption.ID == "" || enumOption.Name == "" || enumOption.Type != "text" {
				return fmt.Errorf("operation %#v id, name empty or not string, or type not text", option)
			}
			*enumOptions = append(*enumOptions, enumOption)
		} else if option, ok := optionVal.(bson.M); ok {
			enumOption := EnumVal{}
			enumOption.ID = getString(option["id"])
			enumOption.Name = getString(option["name"])
			enumOption.Type = getString(option["type"])
			enumOption.IsDefault = getBool(option["is_default"])
			if enumOption.ID == "" || enumOption.Name == "" || enumOption.Type != "text" {
				return fmt.Errorf("operation %#v id, name empty or not string, or type not text", option)
			}
			*enumOptions = append(*enumOptions, enumOption)
		} else if option, ok := optionVal.(bson.D); ok {
			opt := option.Map()
			enumOption := EnumVal{}
			enumOption.ID = getString(opt["id"])
			enumOption.Name = getString(opt["name"])
			enumOption.Type = getString(opt["type"])
			enumOption.IsDefault = getBool(opt["is_default"])
			if enumOption.ID == "" || enumOption.Name == "" || enumOption.Type != "text" {
				return fmt.Errorf("operation %#v id, name empty or not string, or type not text", option)
			}
			*enumOptions = append(*enumOptions, enumOption)
		} else {
			return fmt.Errorf("unknow optionVal type: %#v", optionVal)
		}
	}
	return nil
}

// parseFloatOption  parse float data in option
func ParseFloatOption(ctx context.Context, val interface{}) FloatOption {
	rid := util.ExtractRequestIDFromContext(ctx)
	floatOption := FloatOption{}
	if nil == val || "" == val {
		return floatOption
	}
	switch option := val.(type) {
	case string:
		floatOption.Min = gjson.Get(option, "min").Raw
		floatOption.Max = gjson.Get(option, "max").Raw
	case map[string]interface{}:
		floatOption.Min = getString(option["min"])
		floatOption.Max = getString(option["max"])
	case bson.M:
		floatOption.Min = getString(option["min"])
		floatOption.Max = getString(option["max"])
	case bson.D:
		opt := option.Map()
		floatOption.Min = getString(opt["min"])
		floatOption.Max = getString(opt["max"])
	default:
		blog.Warnf("unknow val type: %#v, rid: %s", val, rid)
	}
	return floatOption
}

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
		enumOption, err := ParseEnumOption(ctx, attribute.Option)
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
			return "", fmt.Errorf("parse options for list type failed, option not slice type, option: %+v", attribute.Option)
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
	return "", nil
}

var HostApplyFieldMap = map[string]bool{
	common.BKOperatorField:              true,
	common.BKBakOperatorField:           true,
	"bk_state":                          true,
	"bk_sla":                            true,
	common.BKHostInnerIPField:           false,
	common.BKHostOuterIPField:           false,
	common.BKAssetIDField:               false,
	common.BKSNField:                    false,
	"bk_comment":                        false,
	"bk_service_term":                   false,
	common.BKCloudIDField:               false,
	"bk_state_name":                     false,
	"bk_province_name":                  false,
	"bk_isp_name":                       false,
	common.BKHostNameField:              false,
	common.BKOSTypeField:                false,
	common.BKOSNameField:                false,
	"bk_os_version":                     false,
	"bk_os_bit":                         false,
	"bk_cpu":                            false,
	"bk_cpu_mhz":                        false,
	"bk_cpu_module":                     false,
	"bk_mem":                            false,
	"bk_disk":                           false,
	"bk_mac":                            false,
	"bk_outer_mac":                      false,
	common.HostFieldDockerClientVersion: false,
	common.HostFieldDockerServerVersion: false,
	common.CreateTimeField:              false,
	common.LastTimeField:                false,
	common.BKImportFrom:                 false,
	common.BKCloudInstIDField:           false,
	common.BKCloudHostStatusField:       false,
	common.BKCloudVendor:                false,
}

// CheckAllowHostApplyOnField 检查字段是否能用于主机属性自动应用
func CheckAllowHostApplyOnField(field string) bool {
	if allow, exist := HostApplyFieldMap[field]; exist == true {
		return allow
	}
	return true
}

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

func (sa *SubAttribute) Validate(ctx context.Context, data interface{}, key string) (rawError errors.RawErrorInfo) {
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
