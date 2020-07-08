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

	"configcenter/src/common/json"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

func (a *Authorize) calculatePolicy(
	ctx context.Context,
	resources []types.Resource,
	p *operator.Policy) (bool, error) {

	// at least have one resources
	if len(resources) == 0 {
		return false, errors.New("auth options at least have one resource")
	}

	if len(resources) > 1 {
		// TODO: support auth functionality across different system.
		return false, errors.New("do not support auth across different system for now")
	}

	if p == nil || p.Operator == "" {
		return false, nil
	}

	if p.Operator == operator.Any {
		return true, nil
	}

	resourceID := resources[0].ID
	authPath, err := getIamPath(resources[0].Attribute)
	if err != nil {
		return false, fmt.Errorf("parse iam path failed, err: %v", err)
	}

	switch p.Operator {
	case operator.And, operator.Or:
		return a.calculateContent(ctx, p, resourceID, authPath, resources[0].Type)
	default:
		return a.calculateFieldValue(ctx, p, resourceID, authPath, resources[0].Type)
	}
}

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

// calculateFieldValue is to calculate the authorize status for attribute.
func (a *Authorize) calculateFieldValue(ctx context.Context, p *operator.Policy, rscID string, authPath []string, resourceType types.ResourceType) (bool, error) {
	// must be a FieldValue type
	fv, can := p.Element.(*operator.FieldValue)
	if !can {
		return false, fmt.Errorf("invalid type %v, should be FieldValue type", reflect.TypeOf(p.Element))
	}

	// check the special resource id at first
	switch fv.Field.Attribute {
	case operator.IamIDKey:
		authorized, err := p.Operator.Operator().Match(rscID, fv.Value)
		if err != nil {
			return false, fmt.Errorf("do %s match calculate failed, err: %v", p.Operator, err)
		}
		return authorized, nil
	case operator.IamPathKey:
		return a.calculateAuthPath(p, fv, authPath)
	default:
		return a.calculateResourceAttribute(ctx, p.Operator, rscID, []*operator.FieldValue{fv}, resourceType)
	}
}

// calculateContent is to calculate the final authorize status, authorized or not.
func (a *Authorize) calculateContent(ctx context.Context, p *operator.Policy, rscID string, authPath []string, resourceType types.ResourceType) (
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
	allAttributes := make([]*operator.FieldValue, 0)
	var resource string

	results := make([]bool, 0)
	for _, policy := range content.Content {
		var authorized bool
		var err error

		switch policy.Operator {
		case operator.And:
			authorized, err = a.calculateContent(ctx, policy, rscID, authPath, resourceType)
			if err != nil {
				return false, err
			}

		case operator.Or:
			authorized, err = a.calculateContent(ctx, policy, rscID, authPath, resourceType)
			if err != nil {
				return false, err
			}

		case operator.Any:
			authorized, err = policy.Operator.Operator().Match(rscID, policy.Element)
			if err != nil {
				return false, fmt.Errorf("match any operator failed, err: %v", err)
			}

		default:

			// must be a FieldValue type
			fv, can := policy.Element.(*operator.FieldValue)
			if !can {
				return false, fmt.Errorf("invalid type %v, should be FieldValue type", reflect.TypeOf(policy.Element))
			}

			// check the special resource id at first
			switch fv.Field.Attribute {
			case operator.IamIDKey:
				authorized, err = policy.Operator.Operator().Match(rscID, fv.Value)
				if err != nil {
					return false, fmt.Errorf("do %s match calculate failed, err: %v", p.Operator, err)
				}

			case operator.IamPathKey:
				authorized, err = a.calculateAuthPath(policy, fv, authPath)
				if err != nil {
					return false, err
				}

			default:

				if policy.Operator != operator.Equal {
					// TODO: confirm this logic with iam.
					// Normally, we need attribute policy should all be "eq" operator.
					return false, fmt.Errorf("unsupported operator %s with attribute auth", policy.Operator)
				}

				// record these attribute for later calculate.
				allAttributes = append(allAttributes, fv)

				// initialize and validate the resource, can not be empty and should be all the same.
				if len(resource) == 0 {
					resource = fv.Field.Resource
				} else {
					if resource != fv.Field.Resource {
						return false, fmt.Errorf("a content have different resource %s / %s, should be same",
							resource, fv.Field.Resource)
					}
				}

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

		// save the result.
		results = append(results, authorized)
	}

	if len(allAttributes) != 0 {
		// we have an authorized with attribute policy.
		// get the instance with these attribute
		yes, err := a.calculateResourceAttribute(ctx, p.Operator, rscID, allAttributes, resourceType)
		if err != nil {
			return false, err
		}
		results = append(results, yes)
	}

	switch p.Operator {
	case operator.And:
		for _, yes := range results {
			if !yes {
				return false, nil
			}
		}
		// all the content is true
		return true, nil

	case operator.Or:
		for _, yes := range results {
			if yes {
				return true, nil
			}
		}
		// all the content is false
		return false, nil

	default:
		return false, fmt.Errorf("invalid policy content with operator: %s ", p.Operator)
	}
}

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

// if a user have a attribute based auth policy, then we need to use the filter constructed by the policy to filter
// out the resources. Then check the resource id is in or not in it. if yes, user is authorized.
func (a *Authorize) calculateResourceAttribute(ctx context.Context, op operator.OperType, rscID string,
	fv []*operator.FieldValue, resourceType types.ResourceType) (bool, error) {

	listOpts := &types.ListWithAttributes{
		Operator:   op,
		IDList:     []string{rscID},
		Attributes: fv,
		Type:       resourceType,
	}

	idList, err := a.fetcher.ListInstancesWithAttributes(ctx, listOpts)
	if err != nil {
		js, _ := json.Marshal(fv)
		return false, fmt.Errorf("fetch instance %s with filter: %s failed, err: %s", rscID, string(js), err)
	}

	if len(idList) == 0 {
		// not authorized
		return false, nil
	}

	for _, id := range idList {
		if id == rscID {
			return true, nil
		}
	}

	// no id matched
	return false, nil
}

func getIamPath(attr types.ResourceAttributes) ([]string, error) {
	path, exist := attr[types.IamPathKey]
	if exist {
		if path == nil {
			return nil, errors.New("have iam path key, but it's value is nil")
		}

		// must be a array string
		p, ok := path.([]string)
		if ok {
			// we got a iam path.
			return p, nil
		}
		return nil, errors.New("iam path value is not an array string type")
	}
	// iam path is not exist.
	return make([]string, 0), nil
}
