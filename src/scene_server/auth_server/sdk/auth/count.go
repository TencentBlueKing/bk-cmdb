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

	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/thirdparty/apigw/iam"
)

func (a *Authorize) countPolicy(ctx context.Context, p *operator.Plan, resourceType iam.IamResourceType) (
	*iam.AuthorizeList, error) {

	if p == nil {
		return &iam.AuthorizeList{}, nil
	}

	switch p.Kind {
	case operator.AlwaysAllowedKind:
		return &iam.AuthorizeList{IsAny: true}, nil
	case operator.AlwaysDeniedKind, "":
		return &iam.AuthorizeList{IsAny: false}, nil
	case operator.ConditionalKind:
		if p.Condition == nil {
			return nil, errors.New("conditional plan has empty condition")
		}
		return a.countCondition(ctx, p.Condition, resourceType)
	default:
		return nil, fmt.Errorf("unsupported plan kind: %s", p.Kind)
	}
}

func (a *Authorize) countCondition(ctx context.Context, cond *operator.AuthCondition,
	resourceType iam.IamResourceType) (*iam.AuthorizeList, error) {

	if hasAncestor(cond) {
		return nil, errors.New("plan condition has ancestor, not support for now")
	}

	//  please refer to the issue #5579 for specific permission scenario classification.
	switch cond.Operator {
	case operator.And, operator.Or:
		content, can := cond.Element.(*operator.Content)
		if !can {
			return nil, errors.New("policy with invalid content field")
		}

		list, err := a.countContent(ctx, cond.Operator, content, resourceType)
		if err != nil {
			return nil, err
		}

		return list, nil

	case operator.Not:
		return a.countNot(ctx, cond, resourceType)

	default:
		fv, can := cond.Element.(*operator.FieldValue)
		if !can {
			return nil, errors.New("policy with invalid FieldValue field")
		}

		if fv.Field.IsID() {
			ids, err := a.countIamIDKey(cond.Operator, fv)
			if err != nil {
				return nil, err
			}

			return &iam.AuthorizeList{Ids: ids}, nil
		}

		// TODO: cause we do not support ancestor field for now
		// So we only need to get resource's other attribute policy.
		opts := &types.ListWithAttributes{
			Operator:     cond.Operator,
			AttrPolicies: []*operator.AuthCondition{cond},
			Type:         resourceType,
		}

		ids, err := a.fetcher.ListInstancesWithAttributes(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("list instance with %s attribute failed, err: %v", cond.Operator, err)
		}

		return &iam.AuthorizeList{Ids: ids}, nil
	}
}

func (a *Authorize) countNot(ctx context.Context, cond *operator.AuthCondition, resourceType iam.IamResourceType) (
	*iam.AuthorizeList, error) {

	content, can := cond.Element.(*operator.Content)
	if !can || content == nil || len(content.Content) != 1 {
		return nil, fmt.Errorf("operator %s content must have exactly one element", operator.Not)
	}

	// not cannot be inverted from a finite id set in memory, query the inverted filter instead.
	opts := &types.ListWithAttributes{
		Operator:     operator.Not,
		AttrPolicies: content.Content,
		Type:         resourceType,
	}

	ids, err := a.fetcher.ListInstancesWithAttributes(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("list instance with %s attribute failed, err: %v", cond.Operator, err)
	}

	return &iam.AuthorizeList{Ids: ids}, nil
}

func (a *Authorize) countIamIDKey(op operator.OperType, fv *operator.FieldValue) ([]string, error) {
	if op == operator.Equal {
		strValue, ok := fv.Value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid policy with operator eq value %v, should be string", fv.Value)
		}
		return []string{strValue}, nil
	}

	if op != operator.In {
		return nil, errors.New("unsupported policy with iam \"id\" key, op is not \"in\"")
	}

	arrayValue, ok := fv.Value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid policy with operator in value %v", fv.Value)
	}

	ids := make([]string, 0)
	for _, id := range arrayValue {
		strID, ok := id.(string)
		if !ok {
			return nil, fmt.Errorf("invalid policy with operator in value: %v, should be string", id)
		}

		ids = append(ids, strID)
	}
	return ids, nil
}

// countContent count all the resource ids according to the operator and content, eg policies.
func (a *Authorize) countContent(ctx context.Context, op operator.OperType, content *operator.Content,
	resourceType iam.IamResourceType) (idList *iam.AuthorizeList, err error) {

	allAttrPolicies := make([]*operator.AuthCondition, 0)
	allList := make([]iam.AuthorizeList, 0)
	idList = new(iam.AuthorizeList)

	for _, policy := range content.Content {
		if policy.Operator.IsLogical() {
			list, err := a.countCondition(ctx, policy, resourceType)
			if err != nil {
				return nil, err
			}
			allList = append(allList, *list)
			continue
		}

		fv, can := policy.Element.(*operator.FieldValue)
		if !can {
			return nil, errors.New("policy with invalid FieldValue field")
		}

		if fv.Field.IsID() {
			list, err := a.countIamIDKey(policy.Operator, fv)
			if err != nil {
				return nil, err
			}
			allList = append(allList, iam.AuthorizeList{Ids: list})
			continue
		}

		// TODO: cause we do not support ancestor field for now
		// So we only need to get resource's other attribute policy.
		allAttrPolicies = append(allAttrPolicies, policy)
	}

	if len(allAttrPolicies) != 0 {
		opts := &types.ListWithAttributes{
			Operator:     op,
			AttrPolicies: allAttrPolicies,
			Type:         resourceType,
		}

		ids, err := a.fetcher.ListInstancesWithAttributes(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("list instance with any attribute failed, err: %v", err)
		}

		allList = append(allList, iam.AuthorizeList{Ids: ids})
	}

	return calculateSet(op, allList)
}

func calculateSetForAnd(sets []iam.AuthorizeList, cnt int) (*iam.AuthorizeList, error) {

	if cnt == 1 {
		return &sets[0], nil
	}

	// now, at least we have two set
	set := make([]string, 0)
	// now, we have at least two set to compare.
	// we use the first set as the base, and compare base element with each
	// element of the rest sets. if base element is hit at the reset of each
	// set, then this element is hit.
	var (
		idBase   int
		baseFlag bool
	)
	// we find base set firstly.
	for id, setBase := range sets {
		if !setBase.IsAny {
			idBase = id
			baseFlag = true
			break
		}
	}
	// if all the sets's isAny is true ,set the isAny flag is true.
	if !baseFlag {
		return &iam.AuthorizeList{IsAny: true}, nil
	}
	for _, base := range sets[idBase].Ids {

		hitOuter := true
		for _, set := range sets[1:] {
			// if this set'isAny is true skip it.
			if set.IsAny {
				continue
			}
			hit := false
			for _, ele := range set.Ids {
				if ele == base {
					// hit in this set.
					hit = true
					break
				}
			}

			if !hit {
				// one of the sets not not hit, then all sets is not hit.
				hitOuter = false
				break
			}

		}

		if hitOuter {
			// all the sets has this element.
			set = append(set, base)
		}
	}
	return &iam.AuthorizeList{Ids: set}, nil
}

func calculateSetForOr(sets []iam.AuthorizeList, cnt int) (*iam.AuthorizeList, error) {
	if cnt == 1 {
		return &sets[0], nil
	}
	// now, at least we have two set.
	all := make(map[string]struct{})
	for _, set := range sets {
		// op is "OR" and the set's isAny is true,return flag true.
		if set.IsAny {
			return &iam.AuthorizeList{IsAny: true}, nil
		}

		for _, ele := range set.Ids {
			all[ele] = struct{}{}
		}
	}

	set := make([]string, 0)
	for ele := range all {
		set = append(set, ele)
	}

	return &iam.AuthorizeList{Ids: set}, nil

}

// calculateSet : put the authorized instance ID into the Ids, op must be one of And or Or.
func calculateSet(op operator.OperType, sets []iam.AuthorizeList) (*iam.AuthorizeList, error) {
	if sets == nil {
		return nil, errors.New("sets can not be nil")
	}

	cnt := len(sets)
	if cnt == 0 {
		return &iam.AuthorizeList{}, nil
	}

	switch op {
	case operator.Or:
		return calculateSetForOr(sets, cnt)
	case operator.And:
		return calculateSetForAnd(sets, cnt)
	default:
		return nil, fmt.Errorf("operator %s is not support to calculate set", op)
	}
}

// hasAncestor returns whether the condition contains an ancestor resource field.
func hasAncestor(cond *operator.AuthCondition) bool {
	if cond == nil {
		return false
	}

	if cond.Operator.IsLogical() {
		content, can := cond.Element.(*operator.Content)
		if !can || content == nil {
			// a plan with invalid content
			return false
		}

		for _, child := range content.Content {
			if hasAncestor(child) {
				return true
			}
		}
		return false
	}

	fv, can := cond.Element.(*operator.FieldValue)
	if !can {
		return false
	}

	return fv.Field.IsAncestor()
}
