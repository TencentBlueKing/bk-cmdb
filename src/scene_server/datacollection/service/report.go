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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) SearchReportSummary(req *restful.Request, resp *restful.Response) {
	var err error
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	param := metadata.ParamSearchNetcollectReport{}
	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result, err := s.Logics.SearchReportSummary(pheader, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReportSummary failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetReportSearchFail)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(result))
	return

}

func (s *Service) SearchReport(req *restful.Request, resp *restful.Response) {
	var err error
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	param := metadata.ParamSearchNetcollectReport{}
	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	count, result, err := s.Logics.SearchReport(pheader, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReportSummary failed, err: %v", err)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetReportSearchFail)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(metadata.RspNetcollectReport{
		Count: count,
		Info:  result,
	}))
	return
}

func (s *Service) ConfirmReport(req *restful.Request, resp *restful.Response) {
	var err error
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	param := metadata.ParamNetcollectComfirm{}
	if err = json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result := s.Logics.ConfirmReport(pheader, param.Reports)
	if len(result.Errors) > 0 {
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
			Msg:  defErr.Error(common.CCErrCollectNetReportConfirmFail),
			Data: result,
		})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(result))
	return

}

func (s *Service) SearchHistory(req *restful.Request, resp *restful.Response) {
	var err error
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	param := metadata.ParamSearchNetcollectReport{}
	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	count, result, err := s.Logics.SearchHistory(pheader, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReportSummary failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetHistorySearchFail)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(metadata.RspNetcollectHistory{
		Count: count,
		Info:  result,
	}))
	return
}
