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

package auth

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

func (a *Authorize) calculatePolicy(
	ctx context.Context,
	resources []types.Resource,
	p *operator.Policy) (bool, error) {

	rid := ctx.Value(common.ContextRequestIDField)
	if blog.V(5) {
		blog.InfoJSON("calculate policy, resource: %s, policy: %s, rid: %s", resources, p, rid)
	}

	if p == nil || p.Operator == "" {
		return false, nil
	}

	if p.Operator == operator.Any {
		return true, nil
	}

	// at least have one resources
	if len(resources) == 0 {
		return false, errors.New("auth options at least have one resource")
	}

	rscMap := make(map[string]*types.Resource)
	for idx, r := range resources {
		rscMap[string(r.Type)] = &resources[idx]
	}

	switch p.Operator {
	case operator.And, operator.Or:
		return a.authContent(ctx, p, rscMap)
	default:
		return a.authFieldValue(ctx, p, rscMap)
	}
}

// calculateAnyPolicy TODO
// returns true when having policy of any resource of the action
func (a *Authorize) calculateAnyPolicy(
	ctx context.Context,
	resources []types.Resource,
	p *operator.Policy) (bool, error) {

	if p == nil || p.Operator == "" {
		return false, nil
	}
	return true, nil
}

// authFieldValue is to calculate the authorize status for attribute.
func (a *Authorize) authFieldValue(ctx context.Context, p *operator.Policy, rscMap map[string]*types.Resource) (bool,
	error) {
	// must be a FieldValue type
	fv, can := p.Element.(*operator.FieldValue)
	if !can {
		return false, fmt.Errorf("invalid type %v, should be FieldValue type", reflect.TypeOf(p.Element))
	}

	authRsc, exist := rscMap[fv.Field.Resource]
	if !exist {
		return false, fmt.Errorf("can not find resource %s which is in iam policy", fv.Field.Resource)
	}

	// check the special resource id at first
	switch fv.Field.Attribute {
	case operator.IamIDKey:
		authorized, err := p.Operator.Operator().Match(authRsc.ID, fv.Value)
		if err != nil {
			return false, fmt.Errorf("do %s match calculate failed, err: %v", p.Operator, err)
		}
		return authorized, nil

	case operator.IamPathKey:
		authPath, err := getIamPath(authRsc.Attribute)
		if err != nil {
			return false, err
		}

		// compatible for cases when resources to be authorized hasn't put its paths in attributes
		if len(authPath) == 0 {
			// compatible for cases when resources to be authorized hasn't put all of its paths in attributes
			return a.authResourceAttribute(ctx, p.Operator, []*operator.Policy{p}, authRsc)
		}

		return a.authWithPath(p, fv, authPath)

	default:
		return a.authResourceAttribute(ctx, p.Operator, []*operator.Policy{p}, authRsc)
	}
}

// NOCC:golint/fnsize(设计如此)
func (a *Authorize) authContent(ctx context.Context, p *operator.Policy, rscMap map[string]*types.Resource) (
	bool, error) {

	content, canContent := p.Element.(*operator.Content)

	if !canContent {
		// not content and field value type at the same time.
		return false, fmt.Errorf("invalid policy with unknown element type: %v", reflect.TypeOf(p.Element))
	}

	if (p.Operator != operator.And) && (p.Operator != operator.Or) {
		return false, fmt.Errorf("invalid policy content with operator: %s ", p.Operator)
	}

	// prepare for attribute match calculate
	allAttrPolicies := make([]*operator.Policy, 0)
	var resource string

	for _, policy := range content.Content {
		var authorized bool
		var err error

		switch policy.Operator {
		case operator.And, operator.Or:
			authorized, err = a.authContent(ctx, policy, rscMap)
			if err != nil {
				return false, err
			}

		case operator.Any:
			authorized, err = policy.Operator.Operator().Match("", policy.Element)
			if err != nil {
				return false, fmt.Errorf("match any operator failed, err: %v", err)
			}

		default:
			var isAttrPolicy bool
			authorized, isAttrPolicy, resource, err = a.authAtomPolicy(ctx, rscMap, policy, resource)
			if err != nil {
				return false, err
			}

			if isAttrPolicy {
				// record these attribute for later calculate.
				allAttrPolicies = append(allAttrPolicies, policy)

				// we try to handle next attribute if it has.
				continue
			}
		}

		// do this check, so that we can return quickly.
		switch p.Operator {
		case operator.And:
			if !authorized {
				return false, nil
			}

		case operator.Or:
			if authorized {
				return true, nil
			}
		}
	}

	if len(allAttrPolicies) != 0 {
		// we have an authorized with attribute policy.
		// get the instance with these attribute
		yes, err := a.authResourceAttribute(ctx, p.Operator, allAttrPolicies, rscMap[resource])
		if err != nil {
			return false, err
		}

		return yes, nil
	}

	switch p.Operator {
	case operator.And:
		// all the content is true
		return true, nil

	case operator.Or:
		// all the content is false
		return false, nil

	default:
		return false, fmt.Errorf("invalid policy content with operator: %s ", p.Operator)
	}
}

func (a *Authorize) authAtomPolicy(ctx context.Context, rscMap map[string]*types.Resource, policy *operator.Policy,
	resource string) (bool, bool, string, error) {

	rid := ctx.Value(common.ContextRequestIDField)

	// must be a FieldValue type
	fv, can := policy.Element.(*operator.FieldValue)
	if !can {
		return false, false, "", fmt.Errorf("invalid type %v, should be FieldValue type",
			reflect.TypeOf(policy.Element))
	}

	authRsc, exist := rscMap[fv.Field.Resource]
	if !exist {
		return false, false, "", fmt.Errorf("can not find resource %s which is in iam policy", fv.Field.Resource)
	}

	// check the special resource id at first
	switch fv.Field.Attribute {
	case operator.IamIDKey:

		authorized, err := policy.Operator.Operator().Match(authRsc.ID, fv.Value)
		if err != nil {
			return false, false, "", fmt.Errorf("do %s match calculate failed, err: %v", policy.Operator, err)
		}

		blog.Infof(">> calculate op %s, val: %v, rsc: '%s', auth: %v, rid: %v", policy.Operator, fv.Value,
			fv.Field.Resource, authorized, rid)

		return authorized, false, "", nil

	case operator.IamPathKey:

		authPath, err := getIamPath(authRsc.Attribute)
		if err != nil {
			return false, false, "", err
		}

		// compatible for cases when resources to be authorized hasn't put its paths in attributes
		var authorized bool
		if len(authPath) == 0 {
			authorized, err = a.authResourceAttribute(ctx, policy.Operator, []*operator.Policy{policy}, authRsc)
		} else {
			authorized, err = a.calculateAuthPath(policy, fv, authPath)
		}

		if err != nil {
			return false, false, "", err
		}

		return authorized, false, "", nil
	default:
		// other attributes only support operator: 'eq', 'in'
		if policy.Operator != operator.Equal && policy.Operator != operator.In {
			return false, false, "", fmt.Errorf("unsupported operator %s with attribute auth", policy.Operator)
		}

		// initialize and validate the resource, can not be empty and should be all the same.
		if len(resource) == 0 {
			resource = fv.Field.Resource
		} else {
			if resource != fv.Field.Resource {
				return false, false, "", fmt.Errorf("a content have different resource %s / %s, should be same",
					authRsc, fv.Field.Resource)
			}
		}

		return false, true, resource, nil
	}
}

// authWithPath if a user has a path based auth policy, then we need to check if the user's path is matched with
// policy's path or not, if one of use's path is matched, then user is authorized.
func (a *Authorize) authWithPath(p *operator.Policy, fv *operator.FieldValue, authPath []string) (bool, error) {
	if !reflect.ValueOf(fv.Value).IsValid() && len(authPath) == 0 {
		// if policy have the path, then user's auth path must can not be empty.
		// we consider this to be unauthorized.
		return false, nil
	}

	for _, path := range authPath {
		matched, err := p.Operator.Operator().Match(path, fv.Value)
		if err != nil {
			return false, fmt.Errorf("do %s match calculate failed, err: %v", p.Operator, err)
		}
		// if one of the path is matched, the we consider it's authorized
		if matched {
			return true, nil
		}
	}

	// no path is matched, not authorized
	return false, nil
}

// authResourceAttribute TODO
// if a user have a attribute based auth policy, then we need to use the filter constructed by the policy to filter
// out the resources. Then check the resource id is in or not in it. if yes, user is authorized.
func (a *Authorize) authResourceAttribute(ctx context.Context, op operator.OperType, attrPolicies []*operator.Policy,
	rsc *types.Resource) (bool, error) {

	listOpts := &types.ListWithAttributes{
		Operator:     op,
		AttrPolicies: attrPolicies,
		Type:         rsc.Type,
	}

	// in some cases, the resource id can be empty
	// eg: when a user has a policy on host's attribute, the action and resources is like following:
	// {"action":{"id":"edit_biz_host"},
	// "resources":[{"system":"bk_cmdb","type":"host","id":"","attribute":{"_bk_iam_path_":["/biz,2/"]}}]}
	if rsc.ID != "" {
		listOpts.IDList = []string{rsc.ID}
	}

	idList, err := a.fetcher.ListInstancesWithAttributes(ctx, listOpts)
	if err != nil {
		js, _ := json.Marshal(listOpts)
		return false, fmt.Errorf("fetch instance %s with filter: %s failed, err: %s", rsc.ID, string(js), err)
	}

	if len(idList) == 0 {
		// not authorized
		return false, nil
	}

	for _, id := range idList {
		if id == rsc.ID {
			return true, nil
		}
	}

	// no id matched
	return false, nil
}

// calculateAuthPath TODO
// if a user has a path based auth policy, then we need to check if the user's path is matched with policy's path or
// not, if one of use's path is matched, then user is authorized.
func (a *Authorize) calculateAuthPath(p *operator.Policy, fv *operator.FieldValue, authPath []string) (bool, error) {
	if !reflect.ValueOf(fv.Value).IsValid() && len(authPath) == 0 {
		// if policy have the path, then user's auth path must can not be empty.
		// we consider this to be unauthorized.
		return false, nil
	}

	for _, path := range authPath {
		matched, err := p.Operator.Operator().Match(path, fv.Value)
		if err != nil {
			return false, fmt.Errorf("do %s match calculate failed, err: %v", p.Operator, err)
		}
		// if one of the path is matched, the we consider it's authorized
		if matched {
			return true, nil
		}
	}

	// no path is matched, not authorized
	return false, nil
}

func getIamPath(attr types.ResourceAttributes) ([]string, error) {
	path, exist := attr[types.IamPathKey]
	if exist {
		if path == nil {
			return nil, errors.New("have iam path key, but it's value is nil")
		}

		// iam path must be a string array
		wKind := reflect.TypeOf(path).Kind()
		if !(wKind == reflect.Slice || wKind == reflect.Array) {
			return nil, errors.New("iam path value is not array or slice type")
		}

		pathVal := reflect.ValueOf(path)
		pathLen := pathVal.Len()
		iamPathArr := make([]string, pathLen)

		for i := 0; i < pathLen; i++ {
			p, ok := pathVal.Index(i).Interface().(string)
			if !ok {
				return nil, errors.New("iam path value is not an array string type")
			}
			iamPathArr[i] = p
		}
		return iamPathArr, nil
	}
	// iam path is not exist.
	return make([]string, 0), nil
}
