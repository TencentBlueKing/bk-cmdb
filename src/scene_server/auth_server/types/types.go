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

package types

import (
	"encoding/json"
	"fmt"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/auth_server/sdk/operator"
)

const (
	SuccessCode                  = 0
	UnauthorizedErrorCode        = 401
	NotFoundErrorCode            = 404
	InternalServerErrorCode      = 500
	UnprocessableEntityErrorCode = 422

	ListAttrMethod             Method = "list_attr"
	ListAttrValueMethod        Method = "list_attr_value"
	ListInstanceMethod         Method = "list_instance"
	FetchInstanceInfoMethod    Method = "fetch_instance_info"
	ListInstanceByPolicyMethod Method = "list_instance_by_policy"
	SearchInstanceMethod       Method = "search_instance"

	IDField   = "id"
	NameField = "display_name"
)

type Method string
type PullResourceReq struct {
	Type   iam.TypeID  `json:"type"`
	Method Method      `json:"method"`
	Filter interface{} `json:"filter,omitempty"`
	Page   Page        `json:"page,omitempty"`
}

func (req *PullResourceReq) UnmarshalJSON(raw []byte) error {
	data := struct {
		Type   iam.TypeID      `json:"type"`
		Method Method          `json:"method"`
		Filter json.RawMessage `json:"filter,omitempty"`
		Page   Page            `json:"page,omitempty"`
	}{}
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return err
	}
	req.Type = data.Type
	req.Method = data.Method
	req.Page = data.Page
	if data.Filter == nil || len(data.Filter) == 0 {
		return nil
	}
	switch data.Method {
	case ListAttrValueMethod:
		filter := ListAttrValueFilter{}
		err := json.Unmarshal(data.Filter, &filter)
		if err != nil {
			return err
		}
		req.Filter = filter
	case ListInstanceMethod, SearchInstanceMethod:
		filter := ListInstanceFilter{}
		err := json.Unmarshal(data.Filter, &filter)
		if err != nil {
			return err
		}
		req.Filter = filter
	case FetchInstanceInfoMethod:
		filter := FetchInstanceInfoFilter{}
		err := json.Unmarshal(data.Filter, &filter)
		if err != nil {
			return err
		}
		req.Filter = filter
	case ListInstanceByPolicyMethod:
		filter := ListInstanceByPolicyFilter{}
		err := json.Unmarshal(data.Filter, &filter)
		if err != nil {
			return err
		}
		req.Filter = filter
	default:
		return fmt.Errorf("method %s is not supported", data.Method)
	}
	return nil
}

type Page struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (page *Page) IsIllegal() bool {
	if page.Limit == 0 {
		return false
	}
	if page.Limit > common.BKMaxPageSize && page.Limit != common.BKNoLimit {
		return true
	}
	return false
}

type ListAttrValueFilter struct {
	Attr    string `json:"attr"`
	Keyword string `json:"keyword,omitempty"`
	// id type is string, int or bool
	IDs []interface{} `json:"ids,omitempty"`
}

type ListInstanceFilter struct {
	Parent  *ParentFilter `json:"parent,omitempty"`
	Keyword string        `json:"keyword,omitempty"`
}

type ParentFilter struct {
	Type iam.TypeID `json:"type"`
	ID   string     `json:"id"`
}

type ResourceTypeChainFilter struct {
	SystemID string     `json:"system_id"`
	ID       iam.TypeID `json:"id"`
}

type FetchInstanceInfoFilter struct {
	IDs   []string `json:"ids"`
	Attrs []string `json:"attrs,omitempty"`
}

type ListInstanceByPolicyFilter struct {
	Expression *operator.Policy `json:"expression"`
}

type AttrResource struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type ListAttrValueResult struct {
	Count   int64               `json:"count"`
	Results []AttrValueResource `json:"results"`
}

type AttrValueResource struct {
	// id type is string, int or bool
	ID          interface{} `json:"id"`
	DisplayName string      `json:"display_name"`
}

type ListInstanceResult struct {
	Count   int64              `json:"count"`
	Results []InstanceResource `json:"results"`
}

type InstanceResource struct {
	ID          string         `json:"id"`
	DisplayName string         `json:"display_name"`
	Path        []InstancePath `json:"path,omitempty"`
}

type InstancePath struct {
	Type        iam.TypeID `json:"type"`
	ID          string     `json:"id"`
	DisplayName string     `json:"display_name"`
}

// ResourcePullMethod iam resource pull callback methods
type ResourcePullMethod struct {
	ListAttr      func(kit *rest.Kit, resourceType iam.TypeID) ([]AttrResource, error)
	ListAttrValue func(kit *rest.Kit, resourceType iam.TypeID, filter *ListAttrValueFilter, page Page) (
		*ListAttrValueResult, error)

	ListInstance func(kit *rest.Kit, resourceType iam.TypeID, filter *ListInstanceFilter, page Page) (
		*ListInstanceResult, error)

	FetchInstanceInfo func(kit *rest.Kit, resourceType iam.TypeID, filter *FetchInstanceInfoFilter) (
		[]map[string]interface{}, error)

	ListInstanceByPolicy func(kit *rest.Kit, resourceType iam.TypeID, filter *ListInstanceByPolicyFilter, page Page) (
		*ListInstanceResult, error)
}
