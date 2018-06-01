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
 
package base

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/http/httpclient"
	"encoding/json"
)

// BaseLogic logic 实现的基类
type BaseLogic struct {
	HttpCli *httpclient.HttpClient
}

// CreateHttpClient create the default http client and set the header
func (cli *BaseLogic) CreateHttpClient() error {

	cli.HttpCli = httpclient.NewHttpClient()

	cli.HttpCli.SetHeader("Content-Type", "application/json")
	cli.HttpCli.SetHeader("Accept", "application/json")

	return nil
}

// IsSuccess check the response
func (cli *BaseLogic) IsSuccess(rst []byte) (*api.BKAPIRsp, bool) {

	var rstRes api.BKAPIRsp
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error: %s", jserr.Error())
		return &rstRes, false
	}

	if rstRes.Code != common.CCSuccess {
		return &rstRes, false
	}

	return &rstRes, true

}
