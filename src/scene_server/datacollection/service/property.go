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

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateProperty create net property
func (s *Service) CreateProperty(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	netPropertyInfo := meta.NetcollectProperty{}
	if err := json.NewDecoder(req.Request.Body).Decode(&netPropertyInfo); nil != err {
		blog.Errorf("[NetProperty] add property failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result, err := s.Logics.AddProperty(pheader, netPropertyInfo)
	if nil != err {
		if err.Error() == defErr.Error(common.CCErrCollectNetPropertyCreateFail).Error() {
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}

		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}

// UpdateProperty update net property
func (s *Service) UpdateProperty(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	netPropertyID, err := checkNetPropertyIDPathParam(defErr, req.PathParameter("netcollect_property_id"))
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	netPropertyInfo := meta.NetcollectProperty{}
	if err = json.NewDecoder(req.Request.Body).Decode(&netPropertyInfo); nil != err {
		blog.Errorf("[NetProperty] update property failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if err = s.Logics.UpdateProperty(pheader, netPropertyID, netPropertyInfo); nil != err {
		if err.Error() == defErr.Error(common.CCErrCollectNetPropertyUpdateFail).Error() {
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}

		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

// BatchCreateProperty batch create net property
func (s *Service) BatchCreateProperty(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	batchAddNetProperty := new(meta.BatchAddNetProperty)
	if err := json.NewDecoder(req.Request.Body).Decode(&batchAddNetProperty); err != nil {
		blog.Errorf("[NetProperty] batch add property failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	propertyList := batchAddNetProperty.Data
	resultList, hasError := s.Logics.BatchCreateProperty(pheader, propertyList)
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

// SearchProperty search net property
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

// DeleteProperty delete net propertys
func (s *Service) DeleteProperty(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	deleteNetPropertyBatchOpt := new(meta.DeleteNetPropertyBatchOpt)
	if err := json.NewDecoder(req.Request.Body).Decode(deleteNetPropertyBatchOpt); nil != err {
		blog.Errorf("[NetProperty] delete net property batch, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	for _, netPropertyID := range deleteNetPropertyBatchOpt.NetcollectPropertyIDs {
		if err := s.Logics.DeleteProperty(pheader, netPropertyID); nil != err {
			blog.Errorf("[NetProperty] delete net property failed, with netcollect_property_id [%d], err: %v", netPropertyID, err)

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

func checkNetPropertyIDPathParam(defErr errors.DefaultCCErrorIf, ID string) (uint64, error) {
	netPropertyID, err := strconv.ParseUint(ID, 10, 64)
	if nil != err {
		blog.Errorf("[NetProperty] update net property with id[%s] to parse the net property id, error: %v", ID, err)
		return 0, defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKNetcollectPropertyIDField)
	}
	if 0 == netPropertyID {
		blog.Errorf("[NetProperty] update net property with id[%d] should not be 0", netPropertyID)
		return 0, defErr.Error(common.CCErrCommHTTPInputInvalid)
	}

	return netPropertyID, nil
}
