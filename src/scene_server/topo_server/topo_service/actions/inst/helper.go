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

package inst

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"encoding/json"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
	"reflect"
)

func hasHost(req *restful.Request, hostCtr string, condition map[string][]int) (bool, error) {

	url := fmt.Sprintf("%s/host/v1/meta/hosts/module/config/search", hostCtr)
	inputJSON, _ := json.Marshal(condition)
	rst, rstErr := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJSON))
	if nil != rstErr {
		return false, rstErr
	}
	blog.Debug("get host(%s) config:%s input:%s", url, rst, inputJSON)
	js, err := simplejson.NewJson([]byte(rst))
	if nil != err {
		return false, err
	}
	rstData, _ := js.Map()
	if subData, ok := rstData["data"]; ok {
		if info, infoOk := subData.([]interface{}); infoOk {
			if len(info) > 0 {
				return false, nil
			}
		} else if nil == subData {
			return true, nil
		} else {
			return false, fmt.Errorf("the data is not array, the kind is %s", reflect.TypeOf(subData))
		}
	} else {
		return false, fmt.Errorf("not found the data in result")
	}
	return true, nil
}

func getInnerNameByObjectID(objID string) string {
	switch objID {
	case common.BKInnerObjIDModule:
		return common.BKInnerObjIDModule
	case common.BKInnerObjIDApp:
		return common.BKInnerObjIDApp
	case common.BKInnerObjIDSet:
		return common.BKInnerObjIDSet
	case common.BKInnerObjIDPlat:
		return common.BKInnerObjIDPlat
	case common.BKInnerObjIDProc:
		return common.BKInnerObjIDProc
	default:
		return common.BKINnerObjIDObject
	}
}
