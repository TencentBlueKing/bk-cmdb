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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/util"
)

const (
	numericType = "numeric"
	booleanType = "boolean"
	stringType  = "string"
)

// parse filter expression to corresponding resource type's mongo query condition,
// nil means having no query condition for the resource type, and using this filter can't get any resource of this type
func ParseFilterToMongo(filter *FilterExpression, resourceType string) (map[string]interface{}, error) {
	operator := filter.Operator

	// parse filter which is composed of multiple sub filters
	if operator == OperatorAnd || operator == OperatorOr {
		if filter.Content == nil || len(filter.Content) == 0 {
			return nil, fmt.Errorf("filter operator %s content can't be empty", operator)
		}
		mongoFilters := make([]map[string]interface{}, 0)
		for _, content := range filter.Content {
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
			operatorMap[operator]: mongoFilters,
		}, nil
	}

	// parse single attribute filter field to [ resourceType, attribute ]
	field := filter.Field
	if field == "" {
		return nil, fmt.Errorf("filter operator %s field can't be empty", operator)
	}
	fieldArr := strings.Split(field, ".")
	if len(fieldArr) != 2 {
		return nil, fmt.Errorf("filter operator %s field %s not in the form of 'resourceType.attribute'", operator, field)
	}
	// if field is another resource's attribute, then the filter isn't for this resource, ignore it
	if fieldArr[0] != resourceType {
		return nil, nil
	}
	value := filter.Value

	switch operator {
	case OperatorEqual, OperatorNotEqual:
		if getValueType(value) == "" {
			return nil, fmt.Errorf("filter operator %s value %#v isn't string, numeric or boolean type", operator, value)
		}
		return map[string]interface{}{
			fieldArr[1]: map[string]interface{}{
				operatorMap[operator]: value,
			},
		}, nil
	case OperatorIn, OperatorNotIn:
		valueArr, ok := value.([]interface{})
		if !ok || len(valueArr) == 0 {
			return nil, fmt.Errorf("filter operator %s value %#v isn't array type or is empty", operator, value)
		}
		valueType := getValueType(valueArr[0])
		if valueType == "" {
			return nil, fmt.Errorf("filter operator %s value %#v isn't string, numeric or boolean array type", operator, value)
		}
		for _, val := range valueArr {
			if getValueType(val) != valueType {
				return nil, fmt.Errorf("filter operator %s value %#v contains values with different types", operator, valueArr)
			}
		}
		return map[string]interface{}{
			fieldArr[1]: map[string]interface{}{
				operatorMap[operator]: valueArr,
			},
		}, nil
	case OperatorLessThan, OperatorLessThanOrEqual, OperatorGreaterThan, OperatorGreaterThanOrEqual:
		if !util.IsNumeric(value) {
			return nil, fmt.Errorf("filter operator %s value %#v isn't numeric type", operator, value)
		}
		return map[string]interface{}{
			fieldArr[1]: map[string]interface{}{
				operatorMap[operator]: value,
			},
		}, nil
	case OperatorContains, OperatorStartsWith, OperatorEndsWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter operator %s value %#v isn't string type", operator, value)
		}
		return map[string]interface{}{
			fieldArr[1]: map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf(operatorRegexFmtMap[operator], valueStr),
			},
		}, nil
	case OperatorNotContains, OperatorNotStartsWith, OperatorNotEndsWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter operator %s value %#v isn't string type", operator, value)
		}
		return map[string]interface{}{
			fieldArr[1]: map[string]interface{}{
				common.BKDBNot: fmt.Sprintf(operatorRegexFmtMap[operator], valueStr),
			},
		}, nil
	case OperatorAny:
		// operator any means having all permissions of this resource
		return make(map[string]interface{}), nil
	default:
		return nil, fmt.Errorf("filter operator %s not supported", operator)
	}
}

var (
	operatorMap = map[Operator]string{
		OperatorAnd:                common.BKDBAND,
		OperatorOr:                 common.BKDBOR,
		OperatorEqual:              common.BKDBEQ,
		OperatorNotEqual:           common.BKDBNE,
		OperatorIn:                 common.BKDBIN,
		OperatorNotIn:              common.BKDBNIN,
		OperatorLessThan:           common.BKDBLT,
		OperatorLessThanOrEqual:    common.BKDBLTE,
		OperatorGreaterThan:        common.BKDBGT,
		OperatorGreaterThanOrEqual: common.BKDBGTE,
	}

	operatorRegexFmtMap = map[Operator]string{
		OperatorContains:      "%s",
		OperatorNotContains:   "%s",
		OperatorStartsWith:    "^%s",
		OperatorNotStartsWith: "^%s",
		OperatorEndsWith:      "%s$",
		OperatorNotEndsWith:   "%s$",
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
