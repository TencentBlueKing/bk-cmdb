/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"encoding/json"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful/v3"
)

// CountInstance count instance.
func (s *service) CountInstance(req *restful.Request, resp *restful.Response) {
	// only allow calling from web-server
	if !httpheader.IsReqFromWeb(req.Request.Header) {
		resp.WriteAsJson(&metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommAuthNotHavePermission,
			ErrMsg: "not allowed to call count api",
		})
		return
	}

	kit := rest.NewKitFromHeader(req.Request.Header, s.engine.CCErr)

	input := make([]map[string]interface{}, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("decode request body failed, err: %v, rid: %s", err, kit.Rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	objID := req.PathParameter("bk_obj_id")

	tableName := common.GetInstTableName(objID, common.BKDefaultOwnerID)
	count, err := s.engine.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, tableName, input)
	if err != nil {
		blog.Errorf("get %s instance count failed, err: %v, rid: %s", objID, err, kit.Rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(count))
}
