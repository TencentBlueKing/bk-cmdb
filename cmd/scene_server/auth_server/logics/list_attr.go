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
	"configcenter/cmd/scene_server/auth_server/types"
	iamtype "configcenter/pkg/ac/iam"
	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/http/rest"
	"configcenter/pkg/metadata"
)

// ListAttr list enumeration and list type attributes of instance type resource
func (lgc *Logics) ListAttr(kit *rest.Kit, resourceType iamtype.TypeID) ([]types.AttrResource, error) {
	attrs := make([]types.AttrResource, 0)
	objID := getInstanceResourceObjID(resourceType)
	if objID == "" && !iamtype.IsIAMSysInstance(resourceType) {
		return attrs, nil
	}

	param := metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKPropertyTypeField: map[string]interface{}{
				common.BKDBIN: []interface{}{
					common.FieldTypeEnum,
					common.FieldTypeList,
				},
			},
		},
		Fields: []string{common.BKPropertyIDField, common.BKPropertyNameField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
	}

	var err error
	if iamtype.IsIAMSysInstance(resourceType) {
		objID, err = lgc.GetObjIDFromResourceType(kit.Ctx, kit.Header, resourceType)
		if err != nil {
			blog.ErrorJSON("get object id from resource type failed, error: %s, resource type: %s, rid: %s",
				err, resourceType, kit.Rid)
			return nil, err
		}
	}

	res, err := lgc.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, &param)
	if err != nil {
		blog.ErrorJSON("read model attribute failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}

	if len(res.Info) == 0 {
		return attrs, nil
	}

	for _, attr := range res.Info {
		displayName := attr.PropertyName
		attrs = append(attrs, types.AttrResource{
			ID:          attr.PropertyID,
			DisplayName: displayName,
		})
	}
	return attrs, nil
}
