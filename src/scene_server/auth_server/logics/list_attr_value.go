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
	"strings"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/types"
)

// list enumeration attribute options of instance type resource
func (lgc *Logics) ListAttrValue(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListAttrValueFilter,
	page types.Page) (*types.ListAttrValueResult, error) {

	// get attributes' enumeration options from cache
	objID := getInstanceResourceObjID(resourceType)
	if objID == "" && resourceType != iam.SysInstance {
		return &types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}}, nil
	}
	var attrType string

	param := metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKPropertyIDField:   filter.Attr,
			common.BKPropertyTypeField: common.FieldTypeEnum,
		},
		Fields: []string{common.BKPropertyTypeField, common.BKOptionField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
	}
	var res *metadata.ReadModelAttrResult
	var err error

	// read all non-inner model attributes for SysInstance resource, add object id to distinguish
	if resourceType == iam.SysInstance {
		res, err = lgc.CoreAPI.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, &param)
	} else {
		res, err = lgc.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, &param)
	}
	if err != nil {
		blog.ErrorJSON("read model attribute failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}
	if !res.Result {
		blog.ErrorJSON("read model attribute failed, error code: %s, error message: %s, param: %s, rid: %s", res.Code, res.ErrMsg, param, kit.Rid)
		return nil, res.Error()
	}
	if len(res.Data.Info) == 0 {
		return &types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}}, nil
	}

	attr := res.Data.Info[0]
	attrType = attr.PropertyType
	if attrType != common.FieldTypeEnum {
		return &types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}}, nil
	}

	marshaledOptions, err := json.Marshal(attr.Option)
	if err != nil {
		blog.ErrorJSON("marshal model attribute option failed, error: %s, option: %v, rid: %s", err.Error(), attr.Option, kit.Rid)
		return nil, err
	}
	options := metadata.AttributesOptions{}
	err = json.Unmarshal(marshaledOptions, &options)
	if err != nil {
		blog.Errorf("attribute option %s is invalid, rid: %s", marshaledOptions, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "option")
	}

	// filter options by keyword and ids and pagination
	values := make([]types.AttrValueResource, 0)
	start := page.Offset
	if start >= int64(len(options)) {
		return &types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}}, nil
	}
	var count int64 = 0
	var idMap map[interface{}]bool
	if idLen := len(filter.IDs); idLen > 0 {
		idMap = make(map[interface{}]bool, idLen)
		for _, id := range filter.IDs {
			idMap[id] = true
		}
	}
	for _, option := range options[start:] {
		if count == page.Limit {
			break
		}
		if idMap != nil && !idMap[option.ID] {
			continue
		}
		if filter.Keyword != "" {
			if !strings.Contains(option.ID, filter.Keyword) && !strings.Contains(option.Name, filter.Keyword) {
				continue
			}
		}
		values = append(values, types.AttrValueResource{
			ID:          option.ID,
			DisplayName: option.Name,
		})
		count++
	}
	return &types.ListAttrValueResult{Count: int64(len(options)), Results: values}, nil
}

func (lgc *Logics) ValidateListAttrValueRequest(kit *rest.Kit, req *types.PullResourceReq) (*types.ListAttrValueFilter, error) {
	filter, ok := req.Filter.(types.ListAttrValueFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for list_attr_value method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}

	if filter.Attr == "" {
		blog.ErrorJSON("request filter %s attr not set for list_attr_value method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "filter.attr")
	}

	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}
	return &filter, nil
}
