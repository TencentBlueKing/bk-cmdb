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
	"strconv"
	"strings"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *Service) CreateProperty(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	propertyList := make([]meta.NetcollectProperty, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&propertyList); err != nil {
		blog.Errorf("[NetProperty] add property failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resultList, hasError := s.Logics.AddProperty(pheader, propertyList)
	if hasError {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{
				Result: false,
				Code:   common.CCErrCollectNetPropertyCreateFail,
				ErrMsg: defErr.Error(common.CCErrCollectNetPropertyCreateFail).Error()},
			Data: resultList,
		})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(resultList))
}

func (s *Service) SearchProperty(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	body := new(meta.NetCollSearchParams)
	if err := json.NewDecoder(req.Request.Body).Decode(body); nil != err {
		blog.Errorf("[NetProperty] search net property failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	propertys, err := s.Logics.SearchProperty(pheader, body)
	if nil != err {
		blog.Errorf("[NetProperty] search net property failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCollectNetPropertyGetFail)})
		return
	}

	resp.WriteEntity(meta.SearchNetPropertyResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *propertys,
	})
}

func (s *Service) DeleteProperty(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	deleteNetPropertyBatchOpt := new(meta.DeleteNetPropertyBatchOpt)
	if err := json.NewDecoder(req.Request.Body).Decode(deleteNetPropertyBatchOpt); nil != err {
		blog.Errorf("[NetProperty] delete net property batch , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	netPropertyIDStrArr := strings.Split(deleteNetPropertyBatchOpt.NetcollectPropertyID, ",")
	var netPropertyIDArr []int64

	for _, netPropertyIDStr := range netPropertyIDStrArr {
		netPropertyID, err := strconv.ParseInt(netPropertyIDStr, 10, 64)
		if nil != err {
			blog.Errorf("[NetProperty] delete net property batch, but got invalid net property id, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, common.BKNetcollectPropertyIDField)})
			return
		}
		netPropertyIDArr = append(netPropertyIDArr, netPropertyID)
	}

	for _, netPropertyID := range netPropertyIDArr {
		if err := s.Logics.DeleteProperty(pheader, netPropertyID); nil != err {
			blog.Errorf("[NetProperty] delete net property failed, with netcollect_property_id [%s], err: %v", netPropertyID, err)

			if defErr.Error(common.CCErrCollectNetDeviceObjPropertyNotExist).Error() == err.Error() {
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
				return
			}

			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}
