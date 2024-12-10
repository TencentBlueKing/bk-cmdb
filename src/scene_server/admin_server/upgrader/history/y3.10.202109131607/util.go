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

package y3_10_202109131607

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

type parsedInstancePolicy struct {
	// instanceIDs is parsed by instance id policy element, used as an intermediate param to generate instObjIDMap
	instanceIDs []int64
	// objInstIDMap object id to instance ids mapping specified by policies with 'and' operator of instances in object
	objInstIDMap map[int64][]int64
	// objectIDs is parsed by iam path policy element
	objectIDs []int64
	// isAny specifies if the user has permissions to any of the instances
	isAny bool
}

// parseInstancePolicy parse old iam instance policy(only support id and path) to the form of object and instance ids
func parseInstancePolicy(policy *operator.Policy) (*parsedInstancePolicy, error) {
	if policy == nil || policy.Operator == "" {
		return new(parsedInstancePolicy), nil
	}

	op := policy.Operator

	if op == operator.Any {
		return &parsedInstancePolicy{isAny: true}, nil
	}

	// parse policy which is composed of multiple sub policies, combine these parsed results into one
	if op == operator.And || op == operator.Or {
		return parseCombinedInstancePolicy(policy)
	}

	return parseAutomaticInstancePolicy(policy)
}

func parseCombinedInstancePolicy(policy *operator.Policy) (*parsedInstancePolicy, error) {
	content, ok := policy.Element.(*operator.Content)
	if !ok {
		return nil, fmt.Errorf("invalid policy with unknown element type: %s", reflect.TypeOf(policy.Element))
	}
	if content == nil || len(content.Content) == 0 {
		return nil, fmt.Errorf("policy op(%s) content can't be empty", policy.Operator)
	}

	// instance id policy can only be combined with its object policy in a policy with and operator like:
	// {"op":"AND","content":[{"op":"in","field":"sys_instance.id","value":["1","2"]},
	// {"op":"starts_with","field":"sys_instance._bk_iam_path_","value":"/sys_instance_model,1/"}]}
	if policy.Operator == operator.And {
		if len(content.Content) != 2 {
			return nil, fmt.Errorf("instance policy op(%s) content length is invalid", policy.Operator)
		}

		var instanceIDs []int64
		var objectID int64

		for _, content := range content.Content {
			parsedSubPolicy, err := parseInstancePolicy(content)
			if err != nil {
				return nil, err
			}

			if len(parsedSubPolicy.instanceIDs) > 0 {
				instanceIDs = parsedSubPolicy.instanceIDs
				continue
			}

			if len(parsedSubPolicy.objectIDs) > 0 {
				if len(parsedSubPolicy.objectIDs) != 1 {
					return nil, fmt.Errorf("instance policy op(%s) has instances in multiple objects", policy.Operator)
				}
				objectID = parsedSubPolicy.objectIDs[0]
			}
		}

		return &parsedInstancePolicy{objInstIDMap: map[int64][]int64{objectID: instanceIDs}}, nil
	}

	// or operator policy contains all objects and instances in all parsed sub policies
	parsedPolicy := new(parsedInstancePolicy)
	for _, content := range content.Content {
		parsedSubPolicy, err := parseInstancePolicy(content)
		if err != nil {
			return nil, err
		}

		if parsedPolicy.isAny {
			return &parsedInstancePolicy{isAny: true}, nil
		}

		parsedPolicy.objectIDs = append(parsedPolicy.objectIDs, parsedSubPolicy.objectIDs...)
		if parsedPolicy.objInstIDMap == nil {
			parsedPolicy.objInstIDMap = make(map[int64][]int64)
		}
		for objID, instIDs := range parsedSubPolicy.objInstIDMap {
			parsedPolicy.objInstIDMap[objID] = instIDs
		}
	}
	return parsedPolicy, nil
}

func parseAutomaticInstancePolicy(policy *operator.Policy) (*parsedInstancePolicy, error) {
	fieldValue, ok := policy.Element.(*operator.FieldValue)
	if !ok {
		return nil, fmt.Errorf("invalid policy with unknown element type: %s", reflect.TypeOf(policy.Element).String())
	}
	field := fieldValue.Field

	// edit and delete instance action use sys_instance, and create instance action use sys_instance_model
	if field.Resource != "sys_instance" && field.Resource != "sys_instance_model" {
		return nil, fmt.Errorf("invalid instance policy with invalid element field resource: %s", field.Resource)
	}

	value := fieldValue.Value

	switch field.Attribute {
	case types.IamIDKey:
		var ids []int64
		// one id's policy uses equal operator, while multiple ids are aggregated into one policy with in operator
		switch policy.Operator {
		case operator.Equal:
			switch v := value.(type) {
			case string:
				id, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse policy op(%s) id(%s) failed, err: %v", policy.Operator, v, err)
				}
				ids = []int64{id}
			default:
				id, err := util.GetInt64ByInterface(v)
				if err != nil {
					return nil, fmt.Errorf("parse policy op(%s) id(%#v) failed, err: %v", policy.Operator, v, err)
				}
				ids = []int64{id}
			}
		case operator.In:
			valueArr, ok := value.([]interface{})
			if !ok || len(valueArr) == 0 {
				return nil, fmt.Errorf("policy op(%s) value(%#v) isn't array type or is empty", policy.Operator, value)
			}
			var err error
			ids = make([]int64, len(valueArr))
			for index, val := range valueArr {
				switch v := val.(type) {
				case string:
					ids[index], err = strconv.ParseInt(v, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("parse policy op(%s) id(%s) failed, err: %v", policy.Operator, v, err)
					}
				default:
					ids[index], err = util.GetInt64ByInterface(v)
					if err != nil {
						return nil, fmt.Errorf("parse policy op(%s) id(%#v) failed, err: %v", policy.Operator, v, err)
					}
				}
			}
		default:
			return nil, fmt.Errorf("policy id field operator %s not supported", policy.Operator)
		}

		switch field.Resource {
		case "sys_instance":
			return &parsedInstancePolicy{instanceIDs: ids}, nil
		case "sys_instance_model":
			return &parsedInstancePolicy{objectIDs: ids}, nil
		default:
			return nil, fmt.Errorf("invalid instance policy with invalid element field resource: %s", field.Resource)
		}
	case types.IamPathKey:
		// create instance action has no iam path, instance iam path is in the form of /sys_instance_model,${model id}/
		if field.Resource != "sys_instance" {
			return nil, fmt.Errorf("invalid iam path policy with invalid element field resource: %s", field.Resource)
		}

		switch policy.Operator {
		case operator.StartWith:
			iamPath, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("policy iam path value %#v isn't string type", value)
			}

			pathItemArr := strings.Split(strings.Trim(iamPath, "/"), "/")
			if len(pathItemArr) != 1 {
				return nil, fmt.Errorf("policy iam path value %s has weong num of items", iamPath)
			}

			typeAndID := strings.Split(pathItemArr[0], ",")
			if len(typeAndID) != 2 {
				return nil, fmt.Errorf("policy iam path item %s invalid", pathItemArr[0])
			}

			if typeAndID[0] != "sys_instance_model" {
				return nil, fmt.Errorf("instance policy iam path parent type %s invalid", typeAndID[0])
			}

			id, err := strconv.ParseInt(typeAndID[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("path iam path model id %s to int failed, error: %v", typeAndID[1], err)
			}
			return &parsedInstancePolicy{objectIDs: []int64{id}}, nil
		default:
			return nil, fmt.Errorf("policy id field operator %s not supported", policy.Operator)
		}
	default:
		return nil, fmt.Errorf("policy element field attribute %s not supported", field.Attribute)
	}
}
