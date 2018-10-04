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

func (s *Service) CreateDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	deviceList := make([]meta.NetcollectDevice, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&deviceList); err != nil {
		blog.Errorf("[NetDevice] add device failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resultList, hasError := s.Logics.AddDevices(pheader, deviceList)
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

func (s *Service) SearchDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	body := new(meta.NetCollSearchParams)
	if err := json.NewDecoder(req.Request.Body).Decode(body); nil != err {
		blog.Errorf("[NetDevice] search net device failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	devices, err := s.Logics.SearchDevice(pheader, body)
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

func (s *Service) DeleteDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	deleteNetDeviceBatchOpt := new(meta.DeleteNetDeviceBatchOpt)
	if err := json.NewDecoder(req.Request.Body).Decode(deleteNetDeviceBatchOpt); nil != err {
		blog.Errorf("[NetDevice] delete net device batch, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	deviceIDStrArr := strings.Split(deleteNetDeviceBatchOpt.DeviceIDs, ",")
	var deviceIDArr []int64

	for _, deviceIDStr := range deviceIDStrArr {
		deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
		if nil != err {
			blog.Errorf("[NetDevice] delete net device batch, but got invalid device id, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, common.BKDeviceIDField)})
			return
		}
		deviceIDArr = append(deviceIDArr, deviceID)
	}

	for _, deviceID := range deviceIDArr {
		if err := s.Logics.DeleteDevice(pheader, deviceID); nil != err {
			blog.Errorf("[NetDevice] delete net device failed, with bk_device_id [%s], err: %v", deviceID, err)

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
