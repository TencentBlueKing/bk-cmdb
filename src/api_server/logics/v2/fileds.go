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
	"encoding/json"
	"errors"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

// getObjFieldIDs get the values of properyID and properyName
func (lgc *Logics) GetObjFieldIDs(objID, user string, header http.Header) (common.KvMap, error) {
	conds := mapstr.MapStr{common.BKObjIDField: objID, common.BKOwnerIDField: user, "page": common.KvMap{"skip": 0, "limit": common.BKNoLimit}}
	result, err := lgc.CoreAPI.TopoServer().Object().SelectObjectAttWithParams(context.Background(), header, conds)
	if nil != err {
		blog.Errorf("get %s fields error:%s", objID, err.Error())
		return nil, err
	}

	blog.Info("get %s fields return:%v", objID, result)
	fields, _ := result.Data.([]interface{})
	ret := common.KvMap{}

	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})

		fieldName, _ := mapField[common.BKPropertyNameField].(string)

		blog.Debug("fieldName:%v", fieldName)
		fieldId, _ := mapField[common.BKPropertyIDField].(string)
		propertyType, _ := mapField[common.BKPropertyTypeField].(string)

		blog.Debug("fieldId:%v", fieldId)
		ret[fieldId] = common.KvMap{"name": fieldName, "type": propertyType, "require": mapField[common.BKIsRequiredField]}
	}

	return ret, nil
}

// AutoInputV3Field fields required to automatically populate the current object v3
func (lgc *Logics) AutoInputV3Field(params mapstr.MapStr, objId, user string, header http.Header) (mapstr.MapStr, error) {
	appFields, err := lgc.GetObjFieldIDs(objId, user, header)
	if nil != err {

		blog.Error("CreateApp error:%s", err.Error())
		return nil, errors.New("CC_Err_Comm_APP_Create_FAIL_STR")
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

// httpRequest http request
func httpRequest(url string, body interface{}, header http.Header) (string, error) {
	params, _ := json.Marshal(body)
	blog.Info("input:%s", string(params))
	httpClient := httpclient.NewHttpClient()
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetHeader("Accept", "application/json")
	reply, err := httpClient.POST(url, header, params)
	return string(reply), err
}
