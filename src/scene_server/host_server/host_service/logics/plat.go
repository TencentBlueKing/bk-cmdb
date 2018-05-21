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
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

//IsExistPlat plat is exist
func IsExistPlat(req *restful.Request, objURL string, cond interface{}) (bool, error) {
	condition := make(map[string]interface{})
	condition["fields"] = common.BKAppIDField
	condition["sort"] = common.BKAppIDField
	condition["start"] = 0
	condition["limit"] = 1
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/" + common.BKInnerObjIDPlat + "/search"
	blog.Info("GetAppIDByCond url :%s", url)
	blog.Info("GetAppIDByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetAppIDByCond return :%s", string(reply))
	if err != nil {
		return false, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	cnt, err := js.Get("data").Get("count").Int()
	if nil != err {
		return false, err
	}
	if 1 == cnt {
		return true, nil
	}
	return false, nil
}
