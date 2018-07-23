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

package service

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

//GetSystemFlag get the system define flag
func (cli *Service) GetSystemFlag(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	var result interface{}
	pathParams := req.PathParameters()
	ownerID := pathParams[common.BKOwnerIDField]
	flag := pathParams["flag"]
	cond := make(map[string]interface{})

	h := md5.New()
	h.Write([]byte(flag))
	cipherStr := h.Sum(nil)
	cond[flag] = hex.EncodeToString(cipherStr) + ownerID

	err := cli.Instance.GetOneByCondition(common.BKTableNameSystem, []string{}, cond, &result)
	if nil != err {
		blog.Errorf("get system config error :%v, cond:%#v", err, cond)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectSelectInstFailed, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})
}
