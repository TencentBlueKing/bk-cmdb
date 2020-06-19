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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/types"
)

// fetch resource instances' specified attributes info using instance ids
func (lgc *Logics) FetchInstanceInfo(kit *rest.Kit, req types.PullResourceReq) ([]map[string]interface{}, error) {
	filter, ok := req.Filter.(types.FetchInstanceInfoFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for fetch_instance_info method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}
	if len(filter.IDs) == 0 {
		blog.ErrorJSON("request filter %s ids not set for fetch_instance_info method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "filter.ids")
	}
	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}

	idField := types.GetResourceIDField(req.Type)
	nameField := types.GetResourceNameField(req.Type)
	if idField == "" || nameField == "" {
		blog.Errorf("request type %s is invalid, rid: %s", req.Type, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}

	// if attribute filter is set, add id attribute and convert display_name to the real name field TODO confirm how to provide path attribute
	var attrs []string
	if len(filter.Attrs) > 0 {
		attrs = append(filter.Attrs, idField)
		for index, attr := range attrs {
			if attr == types.NameField {
				attrs[index] = nameField
				break
			}
		}
	}

	cond := make(map[string]interface{})
	// if id type is string, use id array in filter, or else convert to int64
	if isResourceIDStringType(req.Type) {
		cond[idField] = map[string]interface{}{
			common.BKDBIN: filter.IDs,
		}
	} else {
		ids := make([]int64, len(filter.IDs))
		for _, idStr := range filter.IDs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				blog.Errorf("id %s parse int failed, error: %s, rid: %s, skip it", idStr, kit.Rid)
				continue
			}
			ids = append(ids, id)
		}
		cond[idField] = map[string]interface{}{
			common.BKDBIN: ids,
		}
	}

	param := metadata.PullResourceParam{
		Condition: cond,
		Fields:    attrs,
		Limit:     req.Page.Limit,
		Offset:    req.Page.Offset,
	}
	instances, err := lgc.searchAuthResource(kit, param, req.Type)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}

	// covert id and display_name field
	for _, instance := range instances.Info {
		instance[types.IDField] = util.GetStrByInterface(instance[idField])
		delete(instance, idField)
		if instance[nameField] != nil {
			instance[types.NameField] = util.GetStrByInterface(instance[nameField])
			delete(instance, nameField)
		}
	}
	return instances.Info, nil
}
