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

	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

func (a *Authorize) countPolicy(ctx context.Context, p *operator.Policy, resourceType types.ResourceType) ([]string, error) {

	if hasIamPath(p) {
		return nil, errors.New("policy content has _bk_iam_path_, not support for now")
	}

	switch p.Operator {
	case operator.And, operator.Or:
		content, can := p.Element.(*operator.Content)
		if !can {
			return nil, errors.New("policy with invalid content field")
		}

		list, err := a.countContent(ctx, p.Operator, content, resourceType)
		if err != nil {
			return nil, err
		}

		return list, nil

	case operator.Any:
		ids, err := a.countAny(ctx, p, resourceType)
		if err != nil {
			return nil, err
		}

		return ids, nil

	default:
		fv, can := p.Element.(*operator.FieldValue)
		if !can {
			return nil, errors.New("policy with invalid FieldValue field")
		}

		if fv.Field.Attribute == operator.IamIDKey {

			ids, err := a.countIamIDKey(p.Operator, fv)
			if err != nil {
				return nil, err
			}

			return ids, nil

		} else {
			// TODO: cause we do not support _bk_iam_path_ field for now
			// So we only need to check resource's attribute field.

			opts := &types.ListWithAttributes{
				Operator:   p.Operator,
				Attributes: []*operator.FieldValue{fv},
				Type:       resourceType,
			}

			ids, err := a.fetcher.ListInstancesWithAttributes(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("list instance with %s attribute failed, err: %v", p.Operator, err)
			}

			return ids, nil
		}

	}

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

// count all the resource ids according to the operator and content, eg policies.
func (a *Authorize) countContent(ctx context.Context, op operator.OperType, content *operator.Content, resourceType types.ResourceType) (
	[]string, error) {

	allAttrFieldValue := make([]*operator.FieldValue, 0)
	allList := make([][]string, 0)
	for _, policy := range content.Content {

		switch policy.Operator {
		case operator.And, operator.Or:
			content, can := policy.Element.(*operator.Content)
			if !can {
				return nil, errors.New("policy with invalid content field")
			}

			list, err := a.countContent(ctx, policy.Operator, content, resourceType)
			if err != nil {
				return nil, err
			}
			allList = append(allList, list)

		case operator.Any:
			ids, err := a.countAny(ctx, policy, resourceType)
			if err != nil {
				return nil, err
			}

			allList = append(allList, ids)

		default:
			fv, can := policy.Element.(*operator.FieldValue)
			if !can {
				return nil, errors.New("policy with invalid FieldValue field")
			}

			if fv.Field.Attribute == operator.IamIDKey {
				list, err := a.countIamIDKey(policy.Operator, fv)
				if err != nil {
					return nil, err
				}

				allList = append(allList, list)

			} else {
				// TODO: cause we do not support _bk_iam_path_ field for now
				// So we only need to check resource's attribute field.
				allAttrFieldValue = append(allAttrFieldValue, fv)
			}
		}
	}

	if len(allAttrFieldValue) != 0 {
		opts := &types.ListWithAttributes{
			Operator:   op,
			Attributes: allAttrFieldValue,
			Type:       resourceType,
		}

		ids, err := a.fetcher.ListInstancesWithAttributes(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("list instance with any attribute failed, err: %v", err)
		}

		allList = append(allList, ids)
	}

	return calculateSet(op, allList)

}

func (a *Authorize) countAny(ctx context.Context, policy *operator.Policy, resourceType types.ResourceType) ([]string, error) {
	fv, can := policy.Element.(*operator.FieldValue)
	if !can {
		return nil, errors.New("policy operator any with invalid FieldValue field")
	}

	opts := &types.ListWithAttributes{
		Operator:   policy.Operator,
		Attributes: []*operator.FieldValue{fv},
		Type:       resourceType,
	}

	ids, err := a.fetcher.ListInstancesWithAttributes(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("list instance with any attribute failed, err: %v", err)
	}
	return ids, nil
}

// Note: op must be one of And or Or
func calculateSet(op operator.OperType, sets [][]string) ([]string, error) {
	if sets == nil {
		return nil, errors.New("sets can not be nil")
	}

	cnt := len(sets)
	if cnt == 0 {
		return make([]string, 0), nil
	}

	switch op {
	case operator.Or:
		if cnt == 1 {
			return sets[0], nil
		}
		// now, at least we have two set

		all := make(map[string]struct{})

		for _, set := range sets {
			for _, ele := range set {
				all[ele] = struct{}{}
			}
		}

		set := make([]string, 0)
		for ele := range all {
			set = append(set, ele)
		}

		return set, nil

	case operator.And:
		if cnt == 1 {
			return make([]string, 0), nil
		}
		// now, at least we have two set

		set := make([]string, 0)
		// now, we have at least two set to compare.
		// we use the first set as the base, and compare base element with each
		// element of the rest sets. if base element is hit at the reset of each
		// set, then this element is hit.
		for _, base := range sets[0] {
			hitOuter := true
			for _, set := range sets[1:] {
				hit := false
				for _, ele := range set {
					if ele == base {
						// hit in this set
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

		return set, nil

	default:
		return nil, fmt.Errorf("operator %s is not support to calculate set", op)
	}
}

// check user's policy has _bk_iam_path_ or not.
func hasIamPath(p *operator.Policy) bool {
	switch p.Operator {
	case operator.And, operator.Or:
		content, can := p.Element.(*operator.Content)
		if !can {
			// a policy with invalid content
			return false
		}

		for _, c := range content.Content {
			if hasIamPath(c) {
				return true
			}
		}
		return false
	default:
		fv, can := p.Element.(*operator.FieldValue)
		if !can {
			// a policy with invalid FieldValue type
			return false
		}

		if fv.Field.Attribute == types.IamPathKey {
			return true
		}

		return false
	}
}
