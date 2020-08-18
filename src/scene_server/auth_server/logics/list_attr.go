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
	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/types"
)

// list enumeration attributes of instance type resource
func (lgc *Logics) ListAttr(kit *rest.Kit, resourceType iam.TypeID) ([]types.AttrResource, error) {
	attrs := make([]types.AttrResource, 0)
	objID := getInstanceResourceObjID(resourceType)
	if objID == "" && resourceType != iam.SysInstance {
		return attrs, nil
	}

	param := metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKPropertyTypeField: common.FieldTypeEnum,
		},
		Fields: []string{common.BKPropertyIDField, common.BKPropertyNameField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
	}
	var res *metadata.ReadModelAttrResult
	var err error

	// read all non-inner model attributes for SysInstance resource, add object id to distinguish
	if resourceType == iam.SysInstance {
		param.Fields = append(param.Fields, common.BKObjIDField)
		param.Condition[common.BKObjIDField] = map[string]interface{}{
			common.BKDBNIN: []string{
				common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule,
				common.BKInnerObjIDHost, common.BKInnerObjIDProc, common.BKInnerObjIDPlat,
			},
		}
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
		return nil, res.CCError()
	}
	if len(res.Data.Info) == 0 {
		return attrs, nil
	}

	// get object id name map for common instances
	objIDNameMap := make(map[string]string)
	if resourceType == iam.SysInstance {
		objIDs := make([]string, 0)
		for _, attr := range res.Data.Info {
			objIDs = append(objIDs, attr.ObjectID)
		}
		modelRes, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{
			Fields: []string{common.BKObjIDField, common.BKObjNameField},
			Page:   metadata.BasePage{Limit: common.BKNoLimit},
			Condition: map[string]interface{}{
				common.BKObjIDField: map[string]interface{}{
					common.BKDBIN: objIDs,
				},
			},
		})
		if err != nil {
			blog.Errorf("get model failed, error: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		for _, obj := range modelRes.Data.Info {
			objIDNameMap[obj.Spec.ObjectID] = obj.Spec.ObjectName
		}
	}

	for _, attr := range res.Data.Info {
		displayName := attr.PropertyName
		if resourceType == iam.SysInstance {
			displayName = objIDNameMap[attr.ObjectID] + "-" + attr.PropertyName
		}
		attrs = append(attrs, types.AttrResource{
			ID:          attr.PropertyID,
			DisplayName: displayName,
		})
	}
	return attrs, nil
}
