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
	"configcenter/src/thirdparty/apigw/iam"
)

func (a *Authorize) calculatePolicy(ctx context.Context, resources []iam.Resource, p *operator.Plan) (bool, error) {
	rid := ctx.Value(common.ContextRequestIDField)
	if blog.V(5) {
		blog.InfoJSON("calculate policy, resource: %s, policy: %s, rid: %s", resources, p, rid)
	}

	if p == nil {
		return false, nil
	}

	switch p.Kind {
	case operator.AlwaysAllowedKind:
		return true, nil
	case operator.AlwaysDeniedKind:
		return false, nil
	case operator.ConditionalKind:
		if p.Condition == nil {
			return false, errors.New("conditional plan has empty condition")
		}
	default:
		return false, fmt.Errorf("unsupported plan kind: %s", p.Kind)
	}

	if len(resources) != 1 {
		return false, fmt.Errorf("auth options should have exactly one resource, but got %d", len(resources))
	}

	return a.calculateCondition(ctx, p.Condition, &resources[0])
}

func (a *Authorize) calculateCondition(ctx context.Context, cond *operator.AuthCondition, rsc *iam.Resource) (bool,
	error) {

	if cond.Operator.IsLogical() {
		return a.authContent(ctx, cond, rsc)
	}
	return a.authFieldValue(ctx, cond, rsc)
}

// calculateAnyPolicy returns true when having policy of any resource of the action
func (a *Authorize) calculateAnyPolicy(_ context.Context, _ []iam.Resource, p *operator.Plan) (bool, error) {
	if p == nil || p.Kind == operator.AlwaysDeniedKind || p.Kind == "" {
		return false, nil
	}
	return true, nil
}

// authFieldValue is to calculate the authorize status for a compare node.
func (a *Authorize) authFieldValue(ctx context.Context, cond *operator.AuthCondition, rsc *iam.Resource) (
	bool, error) {

	authorized, isAttr, err := a.authAtomPolicy(ctx, rsc, cond)
	if err != nil {
		return false, err
	}
	if !isAttr {
		return authorized, nil
	}

	return a.authResourceAttribute(ctx, cond.Operator, []*operator.AuthCondition{cond}, rsc)
}

// NOCC:golint/fnsize(设计如此)
func (a *Authorize) authContent(ctx context.Context, cond *operator.AuthCondition, rsc *iam.Resource) (bool, error) {
	content, canContent := cond.Element.(*operator.Content)
	if !canContent {
		return false, fmt.Errorf("invalid policy with unknown element type: %v", reflect.TypeOf(cond.Element))
	}

	if !cond.Operator.IsLogical() {
		return false, fmt.Errorf("invalid policy content with operator: %s ", cond.Operator)
	}

	if cond.Operator == operator.Not {
		if content == nil || len(content.Content) != 1 {
			return false, fmt.Errorf("operator %s content must have exactly one element", operator.Not)
		}
		authorized, err := a.calculateCondition(ctx, content.Content[0], rsc)
		if err != nil {
			return false, err
		}
		return !authorized, nil
	}

	// prepare for attribute match calculate
	allAttrPolicies := make([]*operator.AuthCondition, 0)

	for _, policy := range content.Content {
		var authorized bool
		var err error

		if policy.Operator.IsLogical() {
			authorized, err = a.authContent(ctx, policy, rsc)
			if err != nil {
				return false, err
			}
		} else {
			var isAttrPolicy bool
			authorized, isAttrPolicy, err = a.authAtomPolicy(ctx, rsc, policy)
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
		switch cond.Operator {
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
		yes, err := a.authResourceAttribute(ctx, cond.Operator, allAttrPolicies, rsc)
		if err != nil {
			return false, err
		}

		return yes, nil
	}

	switch cond.Operator {
	case operator.And:
		// all the content is true
		return true, nil

	case operator.Or:
		// all the content is false
		return false, nil

	default:
		return false, fmt.Errorf("invalid policy content with operator: %s ", cond.Operator)
	}
}

// authAtomPolicy evaluates a compare node against the resource.
// The second return value is true when the field is a self attribute, which should be calculated later by querying
// instances with attributes. Otherwise, the first return value is the authorized result of id or ancestor match.
func (a *Authorize) authAtomPolicy(ctx context.Context, rsc *iam.Resource, policy *operator.AuthCondition) (
	bool, bool, error) {

	rid := ctx.Value(common.ContextRequestIDField)

	// must be a FieldValue type
	fv, can := policy.Element.(*operator.FieldValue)
	if !can {
		return false, false, fmt.Errorf("invalid type %v, should be FieldValue type", reflect.TypeOf(policy.Element))
	}

	switch {
	case fv.Field.IsID():
		authorized, err := policy.Operator.Operator().Match(rsc.ID, fv.Value)
		if err != nil {
			return false, false, fmt.Errorf("do %s match calculate failed, err: %v", policy.Operator, err)
		}

		blog.Infof(">> calculate op %s, val: %v, rsc: '%s', auth: %v, rid: %v", policy.Operator, fv.Value,
			rsc.Type, authorized, rid)

		return authorized, false, nil

	case fv.Field.IsAncestor():
		authorized, err := matchAncestor(policy.Operator, rsc.Attribute, fv.Field, fv.Value)
		if err != nil {
			return false, false, err
		}
		return authorized, false, nil

	case fv.Field.IsSelfAttr():
		if fv.Field.SelfAttr() == "" {
			return false, false, fmt.Errorf("invalid self attribute field %s", fv.Field)
		}

		return false, true, nil

	default:
		return false, false, fmt.Errorf("invalid condition field %s", fv.Field)
	}
}

// matchAncestor matches the ancestor instance id in resource attributes.
// The attribute value is a string, or a string list when the resource has multiple ancestors of the same type.
func matchAncestor(op operator.OperType, attrs iam.ResourceAttributes, field operator.Field, with interface{}) (
	bool, error) {

	val, exist := attrs[field.String()]
	if !exist {
		return false, fmt.Errorf("can not find ancestor %s in resource attributes", field)
	}

	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		for i := 0; i < rv.Len(); i++ {
			authorized, err := op.Operator().Match(rv.Index(i).Interface(), with)
			if err != nil {
				return false, fmt.Errorf("do %s match calculate failed, err: %v", op, err)
			}
			if authorized {
				return true, nil
			}
		}
		return false, nil
	}

	authorized, err := op.Operator().Match(val, with)
	if err != nil {
		return false, fmt.Errorf("do %s match calculate failed, err: %v", op, err)
	}
	return authorized, nil
}

// authResourceAttribute if a user have an attribute based auth policy, then we need to use the filter constructed by
// the policy to filter out the resources. Then check the resource id is in or not in it. if yes, user is authorized.
func (a *Authorize) authResourceAttribute(ctx context.Context, op operator.OperType,
	attrPolicies []*operator.AuthCondition, rsc *iam.Resource) (bool, error) {

	listOpts := &types.ListWithAttributes{
		Operator:     op,
		AttrPolicies: attrPolicies,
		Type:         rsc.Type,
	}

	// in some cases, the resource id can be empty
	// eg: when a user has a policy on host's attribute, the action and resources is like following:
	// {"action":{"id":"edit_biz_host"},
	// "resources":[{"system":"bk_cmdb","type":"host","id":"","attribute":{"biz":"2"]}}]}
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
