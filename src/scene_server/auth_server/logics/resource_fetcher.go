/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
    10| * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package logics

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/ac/iam"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	sdktypes "configcenter/src/scene_server/auth_server/sdk/types"
)

// ListInstancesWithAttributes list resource instances that user is privileged to access by policy
func (lgc *Logics) ListInstancesWithAttributes(ctx context.Context, opts *sdktypes.ListWithAttributes) ([]string,
	error) {

	resourceType := iamtypes.TypeID(opts.Type)
	rid := util.ExtractRequestIDFromContext(ctx)
	tenantID := util.ExtractTenantIDFromContext(ctx)
	if tenantID == "" {
		blog.Errorf("request tenant id is empty, rid: %s", rid)
		return nil, fmt.Errorf("request tenant id is empty")
	}

	header := headerutil.NewHeaderFromContext(ctx)
	collection := ""
	if iam.IsIAMSysInstance(resourceType) {
		obj, err := lgc.GetObjFromResourceType(ctx, header, resourceType)
		if err != nil {
			blog.Errorf("get object id from resource type(%s) failed, err: %v, rid: %s", resourceType, err, rid)
			return nil, err
		}
		collection = common.GetObjInstTableName(obj.UUID)
	} else {
		collection = getResourceTableName(resourceType)
	}
	idField := GetResourceIDField(resourceType)
	if collection == "" || idField == "" {
		return nil, fmt.Errorf("request type %s is invalid", opts.Type)
	}

	// get aggregated policy for all attribute policies
	policy, err := lgc.getAggrPolicy(opts.Operator, opts.AttrPolicies)
	if err != nil {
		blog.ErrorJSON("get aggregated policy failed, error: %s, operator:%s, attrPolicies: %s, rid: %s",
			err, opts.Operator, opts.AttrPolicies, rid)
		return nil, err
	}

	cond, err := lgc.parseFilterToMongo(ctx, header, policy, resourceType)
	if err != nil {
		blog.ErrorJSON("parse request filter expression %s failed, error: %s, rid: %s", policy, err.Error(), rid)
		return nil, err
	}
	if cond == nil {
		return make([]string, 0), nil
	}
	if len(opts.IDList) > 0 {
		idCond := make(map[string]interface{})
		if isResourceIDStringType(resourceType) {
			idCond[idField] = map[string]interface{}{common.BKDBIN: opts.IDList}
		} else {
			ids := make([]int64, len(opts.IDList))
			for idx, idStr := range opts.IDList {
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					blog.Errorf("parse id %s to int failed, error: %v, rid: %s", idStr, err, rid)
					return nil, err
				}
				ids[idx] = id
			}
			idCond[idField] = map[string]interface{}{common.BKDBIN: ids}
		}
		cond = map[string]interface{}{common.BKDBAND: []map[string]interface{}{idCond, cond}}
	}

	param := metadata.PullResourceParam{
		Collection: collection,
		Condition:  cond,
		Fields:     []string{idField},
		Limit:      common.BKNoLimit,
	}
	res, err := lgc.CoreAPI.CoreService().Auth().SearchAuthResource(ctx, header, param)
	if err != nil {
		blog.Errorf("search auth resource failed, err: %v, param: %#v, rid: %s", err, param, rid)
		return nil, err
	}
	if err := res.CCError(); err != nil {
		blog.Errorf("search auth resource failed, err: %v, param: %#v, rid: %s", err, param, rid)
		return nil, err
	}

	idList := make([]string, 0)
	for _, instance := range res.Data.Info {
		id := util.GetStrByInterface(instance[idField])
		idList = append(idList, id)
	}
	return idList, nil
}

// getAggrPolicy get aggregated policy for all attribute policies
func (lgc *Logics) getAggrPolicy(opType operator.OperType, attrPolicies []*operator.AuthCondition) (
	*operator.AuthCondition, error) {

	if len(attrPolicies) == 0 {
		return nil, fmt.Errorf("attribute policies can't be empty")
	}

	var policy *operator.AuthCondition
	// if operator is a logical one, the policy's element is used as type Content
	// else, the policy's element is used as type FieldValue
	if opType.IsLogical() {
		policy = &operator.AuthCondition{
			Operator: opType,
			Element: &operator.Content{
				Content: attrPolicies,
			},
		}
	} else {
		// in this condition, the length of attribute policies must be 1
		if len(attrPolicies) != 1 {
			return nil, fmt.Errorf("the length of attribute policies is not 1, policies:%#v", attrPolicies)
		}
		policy = attrPolicies[0]
	}

	return policy, nil
}
