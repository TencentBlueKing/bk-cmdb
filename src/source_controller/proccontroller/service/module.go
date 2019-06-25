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
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
)

func (ps *ProctrlServer) GetProc2Module(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetLanguage(req.Request.Header)
	// get the error factory by the language
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("get process2module failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.V(5).Infof("get proc module config condition: %v ", input)
	input = util.SetModOwner(input, util.GetOwnerID(req.Request.Header))
	result := make([]interface{}, 0)

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	if err := ps.Instance.Table(common.BKTableNameProcModule).Find(input).All(ctx, &result); err != nil {
		blog.Errorf("get process2module config failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcSelectProc2Module)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}
