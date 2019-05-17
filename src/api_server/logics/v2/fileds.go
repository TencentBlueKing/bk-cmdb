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

package logics

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

// getObjFieldIDs get the values of properyID and properyName
func (lgc *Logics) GetObjFieldIDs(ctx context.Context, objID string) (common.KvMap, error) {
	conds := mapstr.MapStr{common.BKObjIDField: objID, common.BKOwnerIDField: lgc.ownerID, "page": common.KvMap{"skip": 0, "limit": common.BKNoLimit}}
	result, err := lgc.CoreAPI.TopoServer().Object().SelectObjectAttWithParams(ctx, lgc.header, conds)
	if nil != err {
		blog.Errorf("SelectObjectAttWithParams http do error.get %s fields error:%s,rid:%s", objID, err.Error(), lgc.rid)
		return nil, err
	}
	if !result.Result {
		blog.Errorf("SelectObjectAttWithParams http do error.get %s fields err code:%d, err msg:%s,rid:%s", objID, result.Code, result.ErrMsg, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	blog.V(5).Infof("get %s fields return:%#v,rid:%s", objID, result, lgc.rid)
	fields, _ := result.Data.([]interface{})
	ret := common.KvMap{}

	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})

		fieldName, _ := mapField[common.BKPropertyNameField].(string)

		fieldId, _ := mapField[common.BKPropertyIDField].(string)
		propertyType, _ := mapField[common.BKPropertyTypeField].(string)

		ret[fieldId] = common.KvMap{"name": fieldName, "type": propertyType, "require": mapField[common.BKIsRequiredField]}
	}

	return ret, nil
}

// AutoInputV3Field fields required to automatically populate the current object v3
func (lgc *Logics) AutoInputV3Field(ctx context.Context, params mapstr.MapStr, objID string) (mapstr.MapStr, error) {
	appFields, err := lgc.GetObjFieldIDs(ctx, objID)
	if nil != err {
		return nil, err
	}
	for fieldId, item := range appFields {
		mapItem, _ := item.(common.KvMap)
		_, ok := params[fieldId]

		if !ok {
			strType, _ := mapItem["type"].(string)
			if util.IsStrProperty(strType) {
				params[fieldId] = ""
			} else {
				params[fieldId] = nil

			}
		}
	}

	return params, nil
}
