/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

import (
	"fmt"
	"reflect"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/operator"
)

const (
	numericType = "numeric"
	booleanType = "boolean"
	stringType  = "string"
)

// parse filter expression to corresponding resource type's mongo query condition,
// nil means having no query condition for the resource type, and using this filter can't get any resource of this type
// TODO confirm how to filter path attribute
func ParseFilterToMongo(filter *operator.Policy, resourceType iam.ResourceTypeID) (map[string]interface{}, error) {
	op := filter.Operator

	// parse filter which is composed of multiple sub filters
	if op == operator.And || op == operator.Or {
		content, ok := filter.Element.(*operator.Content)
		if !ok {
			return nil, fmt.Errorf("invalid policy with unknown element type: %s", reflect.TypeOf(filter.Element).String())
		}
		if content == nil || len(content.Content) == 0 {
			return nil, fmt.Errorf("filter op %s content can't be empty", op)
		}
		mongoFilters := make([]map[string]interface{}, 0)
		for _, content := range content.Content {
			mongoFilter, err := ParseFilterToMongo(content, resourceType)
			if err != nil {
				return nil, err
			}
			// ignore other resource filter
			if mongoFilter != nil {
				mongoFilters = append(mongoFilters, mongoFilter)
			}
		}
		if len(mongoFilters) == 0 {
			return nil, nil
		}
		return map[string]interface{}{
			operatorMap[op]: mongoFilters,
		}, nil
	}

	// parse single attribute filter field to [ resourceType, attribute ]
	fieldValue, ok := filter.Element.(*operator.FieldValue)
	if !ok {
		return nil, fmt.Errorf("invalid policy with unknown element type: %s", reflect.TypeOf(filter.Element).String())
	}
	field := fieldValue.Field
	// if field is another resource's attribute, then the filter isn't for this resource, ignore it
	if field.Resource != string(resourceType) {
		return nil, nil
	}
	attribute := field.Attribute
	value := fieldValue.Value
	if value == IDField {
		value = GetResourceIDField(resourceType)
	}
	if value == "display_name" {
		value = GetResourceNameField(resourceType)
	}

	switch op {
	case operator.Equal, operator.NEqual:
		if getValueType(value) == "" {
			return nil, fmt.Errorf("filter op %s value %#v isn't string, numeric or boolean type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				operatorMap[op]: value,
			},
		}, nil
	case operator.In, operator.Nin:
		valueArr, ok := value.([]interface{})
		if !ok || len(valueArr) == 0 {
			return nil, fmt.Errorf("filter op %s value %#v isn't array type or is empty", op, value)
		}
		valueType := getValueType(valueArr[0])
		if valueType == "" {
			return nil, fmt.Errorf("filter op %s value %#v isn't string, numeric or boolean array type", op, value)
		}
		for _, val := range valueArr {
			if getValueType(val) != valueType {
				return nil, fmt.Errorf("filter op %s value %#v contains values with different types", op, valueArr)
			}
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				operatorMap[op]: valueArr,
			},
		}, nil
	case operator.LessThan, operator.LessThanEqual, operator.GreaterThan, operator.GreaterThanEqual:
		if !util.IsNumeric(value) {
			return nil, fmt.Errorf("filter op %s value %#v isn't numeric type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				operatorMap[op]: value,
			},
		}, nil
	case operator.Contains, operator.StartWith, operator.EndWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf(operatorRegexFmtMap[op], valueStr),
			},
		}, nil
	case operator.NContains, operator.NStartWith, operator.NEndWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				common.BKDBNot: fmt.Sprintf(operatorRegexFmtMap[op], valueStr),
			},
		}, nil
	case operator.Any:
		// op any means having all permissions of this resource
		return make(map[string]interface{}), nil
	default:
		return nil, fmt.Errorf("filter op %s not supported", op)
	}
}

var (
	operatorMap = map[operator.OperType]string{
		operator.And:              common.BKDBAND,
		operator.Or:               common.BKDBOR,
		operator.Equal:            common.BKDBEQ,
		operator.NEqual:           common.BKDBNE,
		operator.In:               common.BKDBIN,
		operator.Nin:              common.BKDBNIN,
		operator.LessThan:         common.BKDBLT,
		operator.LessThanEqual:    common.BKDBLTE,
		operator.GreaterThan:      common.BKDBGT,
		operator.GreaterThanEqual: common.BKDBGTE,
	}

	operatorRegexFmtMap = map[operator.OperType]string{
		operator.Contains:   "%s",
		operator.NContains:  "%s",
		operator.StartWith:  "^%s",
		operator.NStartWith: "^%s",
		operator.EndWith:    "%s$",
		operator.NEndWith:   "%s$",
	}
)

func getValueType(value interface{}) string {
	if util.IsNumeric(value) {
		return numericType
	}
	switch value.(type) {
	case string:
		return stringType
	case bool:
		return booleanType
	}
	return ""
}

// get resource id's actual field
func GetResourceIDField(resourceType iam.ResourceTypeID) string {
	switch resourceType {
	case iam.Host:
		return common.BKHostIDField
	case iam.SysEventPushing:
		return common.BKSubscriptionIDField
	case iam.SysModelGroup:
		return common.BKClassificationIDField
	case iam.SysModel, iam.SysInstanceModel:
		return common.BKObjIDField
	case iam.SysInstance:
		return common.BKInstIDField
	case iam.SysAssociationType:
		return common.AssociationKindIDField
	case iam.SysResourcePoolDirectory:
		return common.BKModuleIDField
	case iam.SysCloudArea:
		return common.BKCloudIDField
	case iam.SysCloudAccount:
		return common.BKCloudAccountID
	case iam.SysCloudResourceTask:
		return common.BKCloudTaskID
	case iam.Business:
		return common.BKAppIDField
	case iam.BizCustomQuery, iam.BizProcessServiceTemplate, iam.BizProcessServiceCategory, iam.BizProcessServiceInstance, iam.BizSetTemplate:
		return common.BKFieldID
	//case iam.Set:
	//	return common.BKSetIDField
	//case iam.Module:
	//	return common.BKModuleIDField
	default:
		return ""
	}
}

// get resource display name's actual field
func GetResourceNameField(resourceType iam.ResourceTypeID) string {
	switch resourceType {
	case iam.Host:
		return common.BKHostInnerIPField
	case iam.SysEventPushing:
		return common.BKSubscriptionNameField
	case iam.SysModelGroup:
		return common.BKClassificationNameField
	case iam.SysModel, iam.SysInstanceModel:
		return common.BKObjNameField
	case iam.SysInstance:
		return common.BKInstNameField
	case iam.SysAssociationType:
		return common.AssociationKindNameField
	case iam.SysResourcePoolDirectory:
		return common.BKModuleNameField
	case iam.SysCloudArea:
		return common.BKCloudNameField
	case iam.SysCloudAccount:
		return common.BKCloudAccountName
	case iam.SysCloudResourceTask:
		return common.BKCloudSyncTaskName
	case iam.Business:
		return common.BKAppNameField
	case iam.BizCustomQuery, iam.BizProcessServiceTemplate, iam.BizProcessServiceCategory, iam.BizProcessServiceInstance, iam.BizSetTemplate:
		return common.BKFieldName
	//case iam.Set:
	//	return common.BKSetNameField
	//case iam.Module:
	//	return common.BKModuleNameField
	default:
		return ""
	}
}
