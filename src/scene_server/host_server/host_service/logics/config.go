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
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"encoding/json"

	"github.com/emicklei/go-restful"
)

//GetSystemConfig get system config by condition
func GetSystemConfig(req *restful.Request, objURL string) bool {
	url := objURL + "/object/v1/system/" + common.HostCrossBizField + "/" + common.BKDefaultOwnerID
	blog.Info("GetSystemConfig url :%s", url)
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)
	if err != nil {
		blog.Error("GetSystemConfig error : %v", err)
		return false
	}
	blog.Info("GetSystemConfig return :%s", string(reply))

	var result api.APIRsp
	err = json.Unmarshal([]byte(reply), &result)
	if nil != err {
		blog.Error("GetSystemConfig error : %v", err)
		return false
	}
	if result.Result {
		return true
	}
	blog.Error("GetSystemConfig error : %v", err)
	return false
}
