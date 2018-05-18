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

package util

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
)

// ValidPropertyOption valid property field option
func ValidPropertyOption(propertyType string, option interface{}) error {
	switch propertyType {
	case common.FieldTypeEnum:
		if nil == option {
			return fmt.Errorf(" option is required")
		}

		arrOption, ok := option.([]interface{})
		if false == ok {
			blog.Errorf(" option %v not enum option", option)
			return fmt.Errorf(" option %v not enum option", option)
		}
		for _, o := range arrOption {
			mapOption, ok := o.(map[string]interface{})
			if false == ok {
				blog.Errorf(" option %v not enum option, enum option item must id and name", option)
				return fmt.Errorf(" option %v not enum option, enum option item must id and name", option)
			}
			_, idOk := mapOption["id"]
			_, nameOk := mapOption["name"]
			if false == idOk || false == nameOk {
				blog.Errorf(" option %v not enum option, enum option item must id and name", option)
				return fmt.Errorf(" option %v not enum option, enum option item must id and name", option)

			}
		}
	case common.FieldTypeInt:
		if nil == option {
			return fmt.Errorf(" option is required")
		}

		_, ok := option.(map[string]interface{})
		if false == ok {
			return fmt.Errorf(" option %v not int option", option)
		}

	}
	return nil
}

// IsAssocateProperty  is Assocate property
func IsAssocateProperty(propertyType string) bool {
	if common.FieldTypeSingleAsst == propertyType || common.FieldTypeMultiAsst == propertyType {
		return true
	}

	return false
}

// IsStrProperty  is string property
func IsStrProperty(propertyType string) bool {
	if common.FieldTypeLongChar == propertyType || common.FieldTypeSingleChar == propertyType {
		return true
	}

	return false
}

// IsInnerObject is inner object model
func IsInnerObject(objID string) bool {
	switch objID {
	case common.BKInnerObjIDApp:
		return true
	case common.BKInnerObjIDHost:
		return true
	case common.BKInnerObjIDModule:
		return true
	case common.BKInnerObjIDPlat:
		return true
	case common.BKInnerObjIDProc:
		return true
	case common.BKInnerObjIDSet:
		return true
	}

	return false
}
