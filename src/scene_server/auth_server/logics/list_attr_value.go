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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/types"
)

// list enumeration attribute options of instance type resource
func (lgc *Logics) ListAttrValue(kit *rest.Kit, req types.PullResourceReq) (*types.ListAttrValueResult, error) {
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

	// get attributes' enumeration options from cache
	objID := GetInstanceResourceObjID(req.Type)
	if objID == "" {
		return &types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}}, nil
	}
	attribute, err := lgc.cache.HGetAll(common.BKCacheKeyV3Prefix + "attribute:object" + objID + "id:" + filter.Attr).Result()
	if attribute[common.BKPropertyTypeField] != common.FieldTypeEnum {
		return &types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}}, nil
	}
	options := metadata.AttributesOptions{}
	err = json.Unmarshal([]byte(attribute[common.BKOptionField]), &options)
	if err != nil {
		blog.Errorf("attribute option %s is invalid, rid: %s", attribute[common.BKOptionField], kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "option")
	}
	values := make([]types.AttrValueResource, 0)
	count := int64(len(options))
	start := req.Page.Offset
	if start >= count {
		return &types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}}, nil
	}
	end := req.Page.Offset + req.Page.Limit
	if end > count {
		end = count
	}
	for _, option := range options[start:end] {
		values = append(values, types.AttrValueResource{
			ID:          option.ID,
			DisplayName: option.Name,
		})
	}

	return &types.ListAttrValueResult{Count: count, Results: values}, nil
}
