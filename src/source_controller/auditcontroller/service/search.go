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
	"encoding/json"
	"net/http"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *Service) Get(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	ownerID := util.GetOwnerID(req.Request.Header)

	dat := new(metadata.ObjQueryInput)
	if err := json.NewDecoder(req.Request.Body).Decode(dat); err != nil {
		blog.Error("Get json unmarshal failed,  error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	dat.Condition = util.SetModOwner(dat.Condition, ownerID)
	rows, cnt, err := s.Logics.Search(dat)
	if nil != err {
		blog.Error("get data from data  error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	data := new(metadata.AuditQueryResult) //common.KvMap{"info": rows, "count": cnt}
	data.Result = true
	data.Data.Info = rows
	data.Data.Count = cnt
	resp.WriteEntity(data)
}
