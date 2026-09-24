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

package logics

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"configcenter/src/ac/iam"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/operator"
)

const (
	numericType = "numeric"
	booleanType = "boolean"
	stringType  = "string"
)

// parseFilterToMongo TODO
// parse filter expression to corresponding resource type's mongo query condition,
// nil means having no query condition for the resource type, and using this filter can't get any resource of this type
func (lgc *Logics) parseFilterToMongo(ctx context.Context, header http.Header, filter *operator.AuthCondition,
	resourceType iamtypes.TypeID) (map[string]interface{}, error) {
	if filter == nil || filter.Operator == "" {
		return nil, nil
	}

	op := filter.Operator

	// parse filter which is composed of multiple sub filters
	if op.IsLogical() {
		content, ok := filter.Element.(*operator.Content)
		if !ok {
			return nil, fmt.Errorf("invalid policy with unknown element type: %s",
				reflect.TypeOf(filter.Element).String())
		}
		if content == nil || len(content.Content) == 0 {
			return nil, fmt.Errorf("filter op %s content can't be empty", op)
		}
		if op == operator.Not && len(content.Content) != 1 {
			return nil, fmt.Errorf("filter op %s content must have exactly one element", op)
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
	fieldValue, ok := filter.Element.(*operator.FieldValue)
	if !ok {
		return nil, fmt.Errorf("invalid policy with unknown element type: %s", reflect.TypeOf(filter.Element).String())
	}

	attribute, err := parseConditionField(fieldValue.Field, resourceType)
	if err != nil {
		return nil, err
	}
	if attribute == "" {
		return nil, fmt.Errorf("resource %s condition field %s is invalid", resourceType, fieldValue.Field)
	}

	value := fieldValue.Value
	if fieldValue.Field.IsAncestor() && !isResourceIDStringType(iamtypes.TypeID(fieldValue.Field)) {
		value, err = convertAncestorIDsToInt(op, value)
		if err != nil {
			return nil, fmt.Errorf("convert ancestor ids(%+v) to int failed, err: %v", fieldValue.Value, err)
		}
	}

	mongoFilter, err := lgc.parseOtherFilterCond(op, value, attribute)
	if err != nil {
		return nil, err
	}

	// host ancestors are stored in module-host relation table, not host instance table.
	if fieldValue.Field.IsAncestor() && isHostResourceType(resourceType) {
		return lgc.parseHostAncestorToMongo(ctx, header, mongoFilter)
	}

	return mongoFilter, nil
}

// parseConditionField maps an authorization plan field to the corresponding mongo field.
func parseConditionField(field operator.Field, resourceType iamtypes.TypeID) (string, error) {
	if field.IsID() {
		return GetResourceIDField(resourceType), nil
	}

	if field.IsAncestor() {
		return GetResourceIDField(iamtypes.TypeID(field)), nil
	}

	if field.IsSelfAttr() {
		return field.SelfAttr(), nil
	}

	return "", fmt.Errorf("resource %s condition field %s is invalid", resourceType, field)
}

func convertAncestorIDsToInt(op operator.OperType, value interface{}) (interface{}, error) {
	switch op {
	case operator.Equal:
		return util.GetInt64ByInterface(value)
	case operator.In:
		valueArr, ok := value.([]interface{})
		if !ok || len(valueArr) == 0 {
			return nil, fmt.Errorf("filter op %s value %#v isn't array type or is empty", op, value)
		}
		ids := make([]interface{}, len(valueArr))
		for i, val := range valueArr {
			id, err := util.GetInt64ByInterface(val)
			if err != nil {
				return nil, err
			}
			ids[i] = id
		}
		return ids, nil
	default:
		return nil, fmt.Errorf("filter op %s not supported for ancestor field", op)
	}
}

func (lgc *Logics) parseOtherFilterCond(op operator.OperType, value interface{}, attribute string) (
	map[string]interface{}, error) {

	switch op {
	case operator.Equal:
		if getValueType(value) == "" {
			return nil, fmt.Errorf("filter op %s value %#v isn't string, numeric or boolean type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				operatorMap[op]: value,
			},
		}, nil
	case operator.In:
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
	case operator.StartWith:
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("filter op %s value %#v isn't string type", op, value)
		}
		return map[string]interface{}{
			attribute: map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf(operatorRegexFmtMap[op], valueStr),
			},
		}, nil
	default:
		return nil, fmt.Errorf("filter op %s not supported", op)
	}
}

// parseHostAncestorToMongo converts a host ancestor condition to a host id filter.
// Host's biz/set/module relations are stored in ModuleHostConfig, so we query host ids from that table first.
func (lgc *Logics) parseHostAncestorToMongo(ctx context.Context, header http.Header, cond map[string]interface{}) (
	map[string]interface{}, error) {

	rid := util.ExtractRequestIDFromContext(ctx)
	param := metadata.PullResourceParam{
		Collection: common.BKTableNameModuleHostConfig,
		Condition:  cond,
		Fields:     []string{common.BKHostIDField},
		Limit:      common.BKNoLimit,
	}
	res, err := lgc.CoreAPI.CoreService().Auth().SearchAuthResource(ctx, header, param)
	if err != nil {
		blog.Errorf("search host ancestor relation failed, err: %v, param: %#v, rid: %s", err, param, rid)
		return nil, err
	}
	if err := res.CCError(); err != nil {
		blog.Errorf("search host ancestor relation failed, err: %v, param: %#v, rid: %s", err, param, rid)
		return nil, err
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

var (
	operatorMap = map[operator.OperType]string{
		operator.And:   common.BKDBAND,
		operator.Or:    common.BKDBOR,
		operator.Not:   common.BKDBNOR,
		operator.Equal: common.BKDBEQ,
		operator.In:    common.BKDBIN,
	}

	operatorRegexFmtMap = map[operator.OperType]string{
		operator.StartWith: "^%s",
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
func GetResourceIDField(resourceType iamtypes.TypeID) string {
	switch resourceType {
	case iamtypes.Host, iamtypes.SysHost:
		return common.BKHostIDField
	case iamtypes.SysModelGroup, iamtypes.SysModel, iamtypes.SysInstanceModel, iamtypes.SysModelEvent,
		iamtypes.InstAsstEvent, iamtypes.MainlineModelEvent, iamtypes.SysAssociationType, iamtypes.BizCustomQuery,
		iamtypes.BizProcessServiceTemplate, iamtypes.BizProcessServiceCategory, iamtypes.BizProcessServiceInstance,
		iamtypes.BizSetTemplate, iamtypes.Project, iamtypes.FieldGroupingTemplate:
		return common.BKFieldID
	case iamtypes.SysInstance:
		return common.BKInstIDField
	case iamtypes.SysResourcePoolDirectory:
		return common.BKModuleIDField
	case iamtypes.SysCloudArea:
		return common.BKCloudIDField
	case iamtypes.Business:
		return common.BKAppIDField
	case iamtypes.BizSet:
		return common.BKBizSetIDField
	default:
		if iam.IsIAMSysInstance(resourceType) {
			return common.BKInstIDField
		}
		return ""
	}
}

// GetResourceNameField get resource display name's actual field
func GetResourceNameField(resourceType iamtypes.TypeID) string {
	switch resourceType {
	case iamtypes.Host, iamtypes.SysHost:
		return common.BKHostInnerIPField
	case iamtypes.SysModelGroup:
		return common.BKClassificationNameField
	case iamtypes.SysModel, iamtypes.SysInstanceModel, iamtypes.SysModelEvent, iamtypes.MainlineModelEvent,
		iamtypes.InstAsstEvent:
		return common.BKObjNameField
	case iamtypes.SysAssociationType:
		return common.AssociationKindNameField
	case iamtypes.SysResourcePoolDirectory:
		return common.BKModuleNameField
	case iamtypes.SysCloudArea:
		return common.BKCloudNameField
	case iamtypes.Business:
		return common.BKAppNameField
	case iamtypes.BizSet:
		return common.BKBizSetNameField
	case iamtypes.BizCustomQuery, iamtypes.BizProcessServiceTemplate, iamtypes.BizProcessServiceCategory,
		iamtypes.BizProcessServiceInstance, iamtypes.BizSetTemplate, iamtypes.FieldGroupingTemplate:
		return common.BKFieldName
	case iamtypes.Project:
		return common.BKProjectNameField
	default:
		if iam.IsIAMSysInstance(resourceType) {
			return common.BKInstNameField
		}
		return ""
	}
}
