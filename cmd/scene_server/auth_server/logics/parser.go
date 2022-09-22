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

package logics

import (
	operator2 "configcenter/cmd/scene_server/auth_server/sdk/operator"
	sdktypes "configcenter/cmd/scene_server/auth_server/sdk/types"
	"configcenter/cmd/scene_server/auth_server/types"
	iamtype "configcenter/pkg/ac/iam"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/metadata"
	"configcenter/pkg/util"
)

const (
	numericType = "numeric"
	booleanType = "boolean"
	stringType  = "string"
)

// parseFilterToMongo TODO
// parse filter expression to corresponding resource type's mongo query condition,
// nil means having no query condition for the resource type, and using this filter can't get any resource of this type
func (lgc *Logics) parseFilterToMongo(ctx context.Context, header http.Header, filter *operator2.Policy, resourceType iamtype.TypeID) (map[string]interface{}, error) {
	if filter == nil || filter.Operator == "" {
		return nil, nil
	}

	op := filter.Operator

	if op == operator2.Any {
		// op any means having all permissions of this resource
		return make(map[string]interface{}), nil
	}

	// parse filter which is composed of multiple sub filters
	if op == operator2.And || op == operator2.Or {
		content, ok := filter.Element.(*operator2.Content)
		if !ok {
			return nil, fmt.Errorf("invalid policy with unknown element type: %s", reflect.TypeOf(filter.Element).String())
		}
		if content == nil || len(content.Content) == 0 {
			return nil, fmt.Errorf("filter op %s content can't be empty", op)
		}
		mongoFilters := make([]map[string]interface{}, 0)
		for _, content := range content.Content {
			mongoFilter, err := lgc.parseFilterToMongo(ctx, header, content, resourceType)
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
	fieldValue, ok := filter.Element.(*operator2.FieldValue)
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
	if attribute == types.IDField {
		attribute = GetResourceIDField(resourceType)
	}
	if attribute == "display_name" {
		attribute = GetResourceNameField(resourceType)
	}
	if attribute == sdktypes.IamPathKey {
		return lgc.parseIamPathToMongo(ctx, header, resourceType, op, value)
	}

	switch op {
	case operator2.Equal, operator2.NEqual:
		if getValueType(value) == "" {
			return nil, fmt.Errorf("filter op %s value %#v isn't string, numeric or boolean type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				operatorMap[op]: value,
			},
		}, nil
	case operator2.In, operator2.Nin:
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
	case operator2.LessThan, operator2.LessThanEqual, operator2.GreaterThan, operator2.GreaterThanEqual:
		if !util.IsNumeric(value) {
			return nil, fmt.Errorf("filter op %s value %#v isn't numeric type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				operatorMap[op]: value,
			},
		}, nil
	case operator2.Contains, operator2.StartWith, operator2.EndWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf(operatorRegexFmtMap[op], valueStr),
			},
		}, nil
	case operator2.NContains, operator2.NStartWith, operator2.NEndWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				common.BKDBNot: map[string]interface{}{common.BKDBLIKE: fmt.Sprintf(operatorRegexFmtMap[op], valueStr)},
			},
		}, nil
	default:
		return nil, fmt.Errorf("filter op %s not supported", op)
	}
}

// parseIamPathToMongo TODO
// parse iam path filter expression to corresponding resource type's mongo query condition
func (lgc *Logics) parseIamPathToMongo(ctx context.Context, header http.Header, resourceType iamtype.TypeID, op operator2.OperType, value interface{}) (map[string]interface{}, error) {
	// generate path condition
	cond := make(map[string]interface{}, 0)
	var err error
	switch op {
	case operator2.Equal, operator2.Contains, operator2.StartWith, operator2.EndWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
		}
		cond, err = parseIamPathToMongo(valueStr, common.BKDBEQ)
		if err != nil {
			return nil, err
		}
	case operator2.In:
		pathArr, ok := value.([]interface{})
		if !ok || len(pathArr) == 0 {
			return nil, fmt.Errorf("filter op %s value %#v isn't array type or is empty", op, value)
		}
		condArr := make([]map[string]interface{}, len(pathArr))
		for index, path := range pathArr {
			pathStr, ok := path.(string)
			if !ok {
				return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
			}
			subCond, err := parseIamPathToMongo(pathStr, common.BKDBEQ)
			if err != nil {
				return nil, err
			}
			condArr[index] = subCond
		}
		cond[common.BKDBOR] = condArr
	case operator2.Nin:
		pathArr, ok := value.([]interface{})
		if !ok || len(pathArr) == 0 {
			return nil, fmt.Errorf("filter op %s value %#v isn't array type or is empty", op, value)
		}
		condArr := make([]map[string]interface{}, len(pathArr))
		for index, path := range pathArr {
			pathStr, ok := path.(string)
			if !ok {
				return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
			}
			subCond, err := parseIamPathToMongo(pathStr, common.BKDBNE)
			if err != nil {
				return nil, err
			}
			condArr[index] = subCond
		}
		cond[common.BKDBAND] = condArr
	case operator2.NEqual, operator2.NContains, operator2.NStartWith, operator2.NEndWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
		}
		cond, err = parseIamPathToMongo(valueStr, common.BKDBNE)
		if err != nil {
			return nil, err
		}
	case operator2.Any:
		// op any means having all permissions of this resource
		return make(map[string]interface{}), nil
	default:
		return nil, fmt.Errorf("filter op %s not supported", op)
	}

	// resources except for host has their parent id stored in their instance table(currently all resources only have one layer TODO support multiple layers if needed)
	if resourceType != iamtype.Host {
		return cond, nil
	}

	// get host ids by path condition from host module config table
	param := metadata.PullResourceParam{
		Collection: common.BKTableNameModuleHostConfig,
		Condition:  cond,
		Fields:     []string{common.BKHostIDField},
		Limit:      common.BKNoLimit,
	}
	res, err := lgc.CoreAPI.CoreService().Auth().SearchAuthResource(ctx, header, param)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s", err.Error(), param)
		return nil, err
	}
	if !res.Result {
		blog.ErrorJSON("search auth resource failed, error code: %s, error message: %s, param: %s", res.Code, res.ErrMsg, param)
		return nil, res.Error()
	}
	if len(res.Data.Info) == 0 {
		return nil, nil
	}
	hostIDs := make([]int64, len(res.Data.Info))
	for index, data := range res.Data.Info {
		hostID, err := util.GetInt64ByInterface(data[common.BKHostIDField])
		if err != nil {
			return nil, err
		}
		hostIDs[index] = hostID
	}
	return map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}, nil
}

// parseIamPathToMongo TODO
// parse string format iam path to mongo condition
func parseIamPathToMongo(iamPath string, op string) (map[string]interface{}, error) {
	pathItemArr := strings.Split(strings.Trim(iamPath, "/"), "/")
	cond := make(map[string]interface{}, 0)

	for _, pathItem := range pathItemArr {
		typeAndID := strings.Split(pathItem, ",")
		if len(typeAndID) != 2 {
			return nil, fmt.Errorf("pathItem %s invalid", pathItem)
		}
		idStr := typeAndID[1]
		if idStr == "*" {
			continue
		}
		resourceType := iamtype.TypeID(typeAndID[0])
		idField := GetResourceIDField(resourceType)
		if isResourceIDStringType(resourceType) {
			cond[idField] = map[string]interface{}{
				op: idStr,
			}
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("id %s parse int failed, error: %s", idStr, err.Error())
		}
		cond[idField] = map[string]interface{}{
			op: id,
		}
	}
	return cond, nil
}

var (
	operatorMap = map[operator2.OperType]string{
		operator2.And:              common.BKDBAND,
		operator2.Or:               common.BKDBOR,
		operator2.Equal:            common.BKDBEQ,
		operator2.NEqual:           common.BKDBNE,
		operator2.In:               common.BKDBIN,
		operator2.Nin:              common.BKDBNIN,
		operator2.LessThan:         common.BKDBLT,
		operator2.LessThanEqual:    common.BKDBLTE,
		operator2.GreaterThan:      common.BKDBGT,
		operator2.GreaterThanEqual: common.BKDBGTE,
	}

	operatorRegexFmtMap = map[operator2.OperType]string{
		operator2.Contains:   "%s",
		operator2.NContains:  "%s",
		operator2.StartWith:  "^%s",
		operator2.NStartWith: "^%s",
		operator2.EndWith:    "%s$",
		operator2.NEndWith:   "%s$",
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

// GetResourceIDField get resource id's actual field
func GetResourceIDField(resourceType iamtype.TypeID) string {
	switch resourceType {
	case iamtype.Host:
		return common.BKHostIDField
	case iamtype.SysModelGroup:
		return common.BKFieldID
	case iamtype.SysModel:
		return common.BKFieldID
	case iamtype.SysInstanceModel:
		return common.BKFieldID
	case iamtype.SysModelEvent, iamtype.InstAsstEvent:
		return common.BKFieldID
	case iamtype.MainlineModelEvent:
		return common.BKFieldID
	case iamtype.SysInstance:
		return common.BKInstIDField
	case iamtype.SysAssociationType:
		return common.BKFieldID
	case iamtype.SysResourcePoolDirectory, iamtype.SysHostRscPoolDirectory:
		return common.BKModuleIDField
	case iamtype.SysCloudArea:
		return common.BKCloudIDField
	case iamtype.SysCloudAccount:
		return common.BKCloudAccountID
	case iamtype.SysCloudResourceTask:
		return common.BKCloudTaskID
	case iamtype.Business, iamtype.BusinessForHostTrans:
		return common.BKAppIDField
	case iamtype.BizSet:
		return common.BKBizSetIDField
	case iamtype.BizCustomQuery, iamtype.BizProcessServiceTemplate, iamtype.BizProcessServiceCategory,
		iamtype.BizProcessServiceInstance, iamtype.BizSetTemplate:
		return common.BKFieldID
	// case iam.Set:
	//	return common.BKSetIDField
	// case iam.Module:
	//	return common.BKModuleIDField
	default:
		if iamtype.IsIAMSysInstance(resourceType) {
			return common.BKInstIDField
		}
		return ""
	}
}

// GetResourceNameField get resource display name's actual field
func GetResourceNameField(resourceType iamtype.TypeID) string {
	switch resourceType {
	case iamtype.Host:
		return common.BKHostInnerIPField
	case iamtype.SysModelGroup:
		return common.BKClassificationNameField
	case iamtype.SysModel, iamtype.SysInstanceModel, iamtype.SysModelEvent, iamtype.MainlineModelEvent, iamtype.InstAsstEvent:
		return common.BKObjNameField
	case iamtype.SysAssociationType:
		return common.AssociationKindNameField
	case iamtype.SysResourcePoolDirectory, iamtype.SysHostRscPoolDirectory:
		return common.BKModuleNameField
	case iamtype.SysCloudArea:
		return common.BKCloudNameField
	case iamtype.SysCloudAccount:
		return common.BKCloudAccountName
	case iamtype.SysCloudResourceTask:
		return common.BKCloudSyncTaskName
	case iamtype.Business, iamtype.BusinessForHostTrans:
		return common.BKAppNameField
	case iamtype.BizSet:
		return common.BKBizSetNameField
	case iamtype.BizCustomQuery, iamtype.BizProcessServiceTemplate, iamtype.BizProcessServiceCategory,
		iamtype.BizProcessServiceInstance, iamtype.BizSetTemplate:
		return common.BKFieldName
	// case iam.Set:
	//	return common.BKSetNameField
	// case iam.Module:
	//	return common.BKModuleNameField
	default:
		if iamtype.IsIAMSysInstance(resourceType) {
			return common.BKInstNameField
		}
		return ""
	}
}
