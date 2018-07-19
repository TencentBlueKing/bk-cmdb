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

package privilege

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

var system = &systemAction{}

//system Action
type systemAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/system/{flag}/{bk_supplier_account}", Params: nil, Handler: system.GetSystemFlag})

	// set cc api interface
	system.CreateAction()
}

//GetSystemFlag get the system define flag
func (cli *systemAction) GetSystemFlag(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()
		var result interface{}
		pathParams := req.PathParameters()
		ownerID := pathParams[common.BKOwnerIDField]
		flag := pathParams["flag"]
		cond := make(map[string]interface{})

		h := md5.New()
		h.Write([]byte(flag))
		cipherStr := h.Sum(nil)
		cond[flag] = hex.EncodeToString(cipherStr) + ownerID

		err := cli.CC.InstCli.GetOneByCondition(common.BKTableNameSystem, []string{}, cond, &result)
		if nil != err {
			blog.Error("get system config error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectInstFailed)
		}
		return http.StatusOK, result, nil
	}, resp)
}
