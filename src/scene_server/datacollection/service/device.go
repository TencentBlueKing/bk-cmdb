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

// CreateDevice create device
func (s *Service) CreateDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	deviceInfo := meta.NetcollectDevice{}
	if err := json.NewDecoder(req.Request.Body).Decode(&deviceInfo); nil != err {
		blog.Errorf("[NetDevice] add device failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result, err := s.logics.AddDevice(pheader, deviceInfo)
	if nil != err {
		if err.Error() == defErr.Error(common.CCErrCollectNetDeviceCreateFail).Error() {
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}

		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}

// UpdateDevice update device
func (s *Service) UpdateDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	netDeviceID, err := checkDeviceIDPathParam(defErr, req.PathParameter("device_id"))
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	deviceInfo := meta.NetcollectDevice{}
	if err := json.NewDecoder(req.Request.Body).Decode(&deviceInfo); nil != err {
		blog.Errorf("[NetDevice] update device failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if err = s.logics.UpdateDevice(pheader, netDeviceID, deviceInfo); nil != err {
		if err.Error() == defErr.Error(common.CCErrCollectNetDeviceUpdateFail).Error() {
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}

		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

// BatchCreateDevice batch create device
func (s *Service) BatchCreateDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	batchAddDevice := new(meta.BatchAddDevice)
	if err := json.NewDecoder(req.Request.Body).Decode(&batchAddDevice); nil != err {
		blog.Errorf("[NetDevice] batch add device failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	deviceList := batchAddDevice.Data
	resultList, hasError := s.logics.BatchCreateDevice(pheader, deviceList)
	if hasError {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{
				Result: false,
				Code:   common.CCErrCollectNetDeviceCreateFail,
				ErrMsg: defErr.Error(common.CCErrCollectNetDeviceCreateFail).Error()},
			Data: resultList,
		})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(resultList))
}

// SearchDevice search device
func (s *Service) SearchDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	body := new(meta.NetCollSearchParams)
	if err := json.NewDecoder(req.Request.Body).Decode(body); nil != err {
		blog.Errorf("[NetDevice] search net device failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	devices, err := s.logics.SearchDevice(pheader, body)
	if nil != err {
		blog.Errorf("[NetDevice] search net device failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCollectNetDeviceGetFail)})
		return
	}

	resp.WriteEntity(meta.SearchNetDeviceResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     devices,
	})
}

// DeleteDevice delete device
func (s *Service) DeleteDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	deleteNetDeviceBatchOpt := new(meta.DeleteNetDeviceBatchOpt)
	if err := json.NewDecoder(req.Request.Body).Decode(deleteNetDeviceBatchOpt); nil != err {
		blog.Errorf("[NetDevice] delete net device batch, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	for _, deviceID := range deleteNetDeviceBatchOpt.DeviceIDs {
		if err := s.logics.DeleteDevice(pheader, deviceID); nil != err {
			blog.Errorf("[NetDevice] delete net device failed, with device_id [%d], err: %v", deviceID, err)

			if defErr.Error(common.CCErrCollectNetDeviceHasPropertyDeleteFail).Error() == err.Error() {
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
				return
			}

			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCollectNetDeviceDeleteFail)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func checkDeviceIDPathParam(defErr errors.DefaultCCErrorIf, ID string) (uint64, error) {
	netDeviceID, err := strconv.ParseUint(ID, 10, 64)
	if nil != err {
		blog.Errorf("[NetDevice] update net device with id[%s] to parse the net device id, error: %v", ID, err)
		return 0, defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKDeviceIDField)
	}
	if 0 == netDeviceID {
		blog.Errorf("[NetDevice] update net device with id[%d] should not be 0", netDeviceID)
		return 0, defErr.Error(common.CCErrCommHTTPInputInvalid)
	}

	return netDeviceID, nil
}
