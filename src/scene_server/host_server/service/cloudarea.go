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
	"strings"

	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// FindManyCloudArea  find cloud area list
func (s *Service) FindManyCloudArea(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	input := new(metadata.CloudAreaParameter)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); nil != err {
		blog.Errorf("FindManyCloudArea , but decode body failed, err: %s,rid:%s", err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// set default limit
	if input.Page.Limit == 0 {
		input.Page.Limit = common.BKDefaultLimit
	}

	if input.Page.IsIllegal() {
		blog.Errorf("FindManyCloudArea failed, parse plat page illegal, input:%#v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommPageLimitIsExceeded)})
	}

	query := &metadata.QueryCondition{
		Condition: nil,
		Limit: metadata.SearchLimit{
			Limit:  int64(input.Page.Limit),
			Offset: int64(input.Page.Start),
		},
		Fields: strings.Split(input.Page.Sort, ","),
	}

	res, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, query)
	if nil != err {
		blog.Errorf("FindManyCloudArea htt do error: %v query:%#v,rid:%s", err, query, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if false == res.Result {
		blog.Errorf("FindManyCloudArea http reply error.  query:%#v, err code:%d, err msg:%s,rid:%s", query, res.Code, res.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(res.Code, res.ErrMsg)})

	}
	platIDArr := make([]int64, 0)
	for _, item := range res.Data.Info {
		platID, err := util.GetInt64ByInterface(item[common.BKCloudIDField])
		if err != nil {
			blog.Errorf("FindManyCloudArea failed, parse plat id field failed, input:%+v,rid:%s", err.Error(), res.Data, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
			return
		}
		platIDArr = append(platIDArr, platID)
	}
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByPlatIDs(srvData.ctx, srvData.header, authmeta.Find, platIDArr...); err != nil {
		blog.Errorf("check plats authorization failed, plats: %+v, err: %v", platIDArr, err)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	retData := map[string]interface{}{
		"info":  res.Data.Info,
		"count": res.Data.Count,
		"page":  input.Page,
	}

	resp.WriteEntity(metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     retData,
	})
}
