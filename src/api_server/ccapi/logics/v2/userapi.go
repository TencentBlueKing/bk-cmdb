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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"errors"

	httpcli "configcenter/src/common/http/httpclient"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

type UserAPI struct {
}

func NewUserAPI() *UserAPI {
	return &UserAPI{}
}

//GetNameByID
func (cli *UserAPI) GetNameByID(req *restful.Request, detailURL string) (string, error, int) {
	respV3, err := httpcli.ReqHttp(req, detailURL, common.HTTPSelectGet, nil)
	//http request error
	if err != nil {
		blog.Error("getCustomerGroupList error:%v", err)
		return "", nil, common.CCErrCommHTTPDoRequestFailed
	}

	js, err := simplejson.NewJson([]byte(respV3))
	if nil != err {
		return "", err, common.CC_ERR_Comm_JSON_DECODE
	}

	resV3, err := js.Map()
	if nil != err {
		return "", err, common.CC_ERR_Comm_JSON_DECODE
	}
	result, _ := resV3["result"].(bool)

	if result {
		data, _ := resV3["data"].(map[string]interface{})
		name, _ := data["name"].(string)
		return name, nil, 0
	} else {
		return "", errors.New(resV3[common.HTTPBKAPIErrorMessage].(string)), common.CCErrCommHTTPDoRequestFailed
	}

}
