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
	"reflect"

	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// FindManyCloudArea  find cloud area list
func (s *Service) FindManyCloudArea(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid
	input := new(metadata.CloudAreaSearchParam)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); nil != err {
		blog.Errorf("FindManyCloudArea , but decode body failed, err: %s,rid:%s", err.Error(), rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// set default limit
	if input.Page.Limit == 0 {
		input.Page.Limit = common.BKMaxPageSize
	}
	if input.Page.IsIllegal() {
		blog.Errorf("FindManyCloudArea failed, parse plat page illegal, input:%#v,rid:%s", input, rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommPageLimitIsExceeded)})
		return
	}

	// if not exact search, change the string query to regexp
	if input.Exact != true {
		for k, v := range input.Condition {
			if reflect.TypeOf(v).Kind() == reflect.String {
				field := v.(string)
				input.Condition[k] = mapstr.MapStr{
					common.BKDBLIKE: params.SpecialCharChange(field),
					"$options":      "i",
				}
			}
		}
	}

	filter := input.Condition
	if s.AuthManager.Enabled() && !s.AuthManager.SkipReadAuthorization {
		// auth: get authorized resources
		authorizedPlatIDs, err := s.AuthManager.ListAuthorizedPlatIDs(srvData.ctx, srvData.user)
		if err != nil {
			blog.Errorf("FindManyCloudArea failed, ListAuthorizedPlatIDs failed, user: %s, err: %+v, rid: %s", srvData.user, err, rid)
			_ = resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommListAuthorizedResourceFromIAMFailed)})
			return
		}

		filter = map[string]interface{}{
			common.BKDBAND: []map[string]interface{}{
				input.Condition,
				{
					common.BKCloudIDField: map[string]interface{}{
						common.BKDBIN: authorizedPlatIDs,
					},
				},
			},
		}
	}
	query := &metadata.QueryCondition{
		Fields:    input.Fields,
		Condition: filter,
		Page:      input.Page,
	}

	res, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, query)
	if nil != err {
		blog.Errorf("FindManyCloudArea htt do error: %v query:%#v,rid:%s", err, query, rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if false == res.Result {
		blog.Errorf("FindManyCloudArea http reply error.  query:%#v, err code:%d, err msg:%s, rid:%s", query, res.Code, res.ErrMsg, rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.New(res.Code, res.ErrMsg)})
		return
	}

	// 查询云区域不需要主机数量
	if input.HostCount == false {
		retData := map[string]interface{}{
			"info":  res.Data.Info,
			"count": res.Data.Count,
		}

		_ = resp.WriteEntity(metadata.Response{
			BaseResp: metadata.SuccessBaseResp,
			Data:     retData,
		})
		return
	}

	// add host_count
	mapCloudIDInfo := make(map[int64]mapstr.MapStr, 0)
	intCloudIDArray := make([]int64, 0)
	for _, area := range res.Data.Info {
		intCloudID, err := area.Int64(common.BKCloudIDField)
		if err != nil {
			blog.ErrorJSON("FindManyCloudArea fail with cloudID convert from interface to int64 failed, err: %v, rid: %s", err, rid)
			_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDPlat, common.BKCloudIDField, "int", err.Error())})
			return
		}
		intCloudIDArray = append(intCloudIDArray, intCloudID)
		mapCloudIDInfo[intCloudID] = area
	}
	ret, err := s.findManyCloudAreaAddHostCount(srvData, intCloudIDArray, mapCloudIDInfo)
	if err != nil {
		blog.ErrorJSON("FindManyCloudArea add field host_count failed, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrHostFindManyCloudAreaAddHostCountFieldFail)})
		return
	}

	retData := map[string]interface{}{
		"info":  ret,
		"count": res.Data.Count,
	}

	_ = resp.WriteEntity(metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     retData,
	})
}

// CreatePlat create a plat instance
func (s *Service) CreatePlat(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); nil != err {
		blog.Errorf("CreatePlat , but decode body failed, err: %s,rid:%s", err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	// read supplier account from header
	input[common.BkSupplierAccount] = util.GetOwnerID(req.Request.Header)

	// bk_cloud_name is required and unique
	_, ok := input[common.BKCloudNameField]
	if !ok {
		blog.Errorf("CreatePlat failed, bk_cloud_name field is required, rid: %s", srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostPlatCloudNameIsrequired)})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeResourceCreate(srvData.ctx, srvData.header, 0, authmeta.Model); err != nil {
		blog.Errorf("check create plat authorization failed, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	instInfo := &meta.CreateModelInstance{
		Data: mapstr.NewFromMap(input),
	}

	createRes, err := s.CoreAPI.CoreService().Instance().CreateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, instInfo)
	if nil != err {
		blog.Errorf("CreatePlat error: %s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoInstCreateFailed)})
		return
	}

	if false == createRes.Result {
		blog.Errorf("CreatePlat error.err code:%d,err msg:%s,input:%+v,rid:%s", createRes.Code, createRes.ErrMsg, input, srvData.rid)
		_ = resp.WriteHeaderAndJson(http.StatusInternalServerError, createRes, "application/json")
		return
	}

	// register plat to iam
	platID := int64(createRes.Data.Created.ID)
	if err := s.AuthManager.RegisterPlatByID(srvData.ctx, srvData.header, platID); err != nil {
		blog.Errorf("CreatePlat failed, RegisterPlatByID failed, err: %s, rid:%s", err.Error(), srvData.rid)
		ccErr := srvData.ccErr.CCError(common.CCErrCommRegistResourceToIAMFailed)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: ccErr})
		return
	}

	_ = resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     createRes.Data,
	})

}

func (s *Service) DeletePlat(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	platID, convErr := util.GetInt64ByInterface(req.PathParameter(common.BKCloudIDField))
	if nil != convErr {
		blog.Errorf("the platID is invalid, error info is %s, input:%s.rid:%s", convErr.Error(), platID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, convErr.Error())})
		return
	}
	if 0 == platID {
		blog.Errorf("DelPlat search host, input:%+v,rid:%s", platID, srvData.rid)
		// try delete default area. tip: has host
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoHasHostCheckFailed)})
		return
	}

	params := new(meta.QueryInput)
	params.Fields = common.BKHostIDField
	params.Condition = map[string]interface{}{
		common.BKCloudIDField: platID,
	}

	hostRes, err := s.CoreAPI.CoreService().Host().GetHosts(srvData.ctx, srvData.header, params)
	if nil != err {
		blog.Errorf("DelPlat search host error: %s, input:%+v,rid:%s", err.Error(), platID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	if !hostRes.Result {
		blog.Errorf("DelPlat search host http response error.err code:%d,err msg:%s, input:%+v,rid:%s", hostRes.Code, hostRes.ErrMsg, platID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetFail)})
		return
	}

	// only empty plat could be delete
	if 0 < hostRes.Data.Count {
		blog.Errorf("DelPlat plat [%d] has host data, can not delete,rid:%s", platID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoHasHostCheckFailed)})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByPlatIDs(srvData.ctx, srvData.header, authmeta.Delete, platID); err != nil {
		blog.Errorf("check delete plat authorization failed, plat: %d, err: %v, rid: %s", platID, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	iamResource, err := s.AuthManager.MakeResourcesByPlatID(srvData.header, authmeta.Delete, platID)
	if err != nil {
		blog.Errorf("DelPlat failed, MakeResourcesByPlatID failed, err: %v, input:%d, rid:%s", err, platID, srvData.rid)
		result := &meta.RespError{
			Msg: srvData.ccErr.Errorf(common.CCErrTopoInstDeleteFailed),
		}
		_ = resp.WriteError(http.StatusInternalServerError, result)
		return
	}
	delCond := &meta.DeleteOption{
		Condition: mapstr.MapStr{common.BKCloudIDField: platID},
	}

	res, err := s.CoreAPI.CoreService().Instance().DeleteInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, delCond)
	if nil != err {
		blog.Errorf("DelPlat do error: %v, input:%d,rid:%s", err, platID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoInstDeleteFailed)})
		return
	}

	if false == res.Result {
		blog.Errorf("DelPlat http response error. err code:%d,err msg:%s,input:%s,rid:%s", res.Code, res.ErrMsg, platID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(res.Code, res.ErrMsg)})
		return

	}

	// deregister plat
	if err := s.AuthManager.Authorize.DeregisterResource(srvData.ctx, iamResource...); err != nil {
		blog.Errorf("DelPlat success, but DeregisterResource from iam failed, platID: %d, err: %+v,rid:%s", platID, err, srvData.rid)
		ccErr := srvData.ccErr.CCError(common.CCErrCommUnRegistResourceToIAMFailed)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: ccErr})
		return
	}

	_ = resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     "",
	})
}

func (s *Service) UpdatePlat(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	// parse platID from url
	platIDStr := req.PathParameter(common.BKCloudIDField)
	platID, err := util.GetInt64ByInterface(platIDStr)
	if nil != err {
		blog.Infof("UpdatePlat failed, parse platID failed, platID: %s, err: %s, rid:%s", platIDStr, err.Error(), srvData.rid)
		ccErr := srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: ccErr})
		return
	}
	if 0 == platID {
		blog.Infof("UpdatePlat failed, update built in cloud area forbidden, platID:%+v, rid:%s", platID, srvData.rid)
		ccErr := srvData.ccErr.Error(common.CCErrTopoUpdateBuiltInCloudForbidden)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: ccErr})
		return
	}

	// decode request body
	input := mapstr.MapStr{}
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("UpdatePlat failed, err:%+v, rid:%s", err, srvData.rid)
		ccErr := srvData.ccErr.Errorf(common.CCErrCommJSONUnmarshalFailed)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: ccErr})
		return
	}

	// update plat
	updateOption := &meta.UpdateOption{
		Data: input,
		Condition: map[string]interface{}{
			common.BKCloudIDField: platID,
		},
	}
	res, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, updateOption)
	if nil != err {
		blog.ErrorJSON("UpdatePlat failed, UpdateInstance failed, input:%s, err:%s, rid:%s", updateOption, err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoInstDeleteFailed)})
		return
	}
	if res.Result == false || res.Code != 0 {
		blog.ErrorJSON("UpdatePlat failed, UpdateInstance failed, input:%s, response:%s, rid:%s", updateOption, res, srvData.rid)
		ccErr := &meta.RespError{Msg: errors.New(res.Code, res.ErrMsg)}
		_ = resp.WriteError(http.StatusInternalServerError, ccErr)
		return
	}

	// auth: sync resource info to iam
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKCloudIDField: platID,
		},
	}

	platInfo, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, query)
	if nil != err {
		blog.Errorf("UpdatePlat ReadInstance htt do error: %v query:%#v,rid:%s", err, query, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if false == platInfo.Result {
		blog.Errorf("UpdatePlat ReadInstance http reply error.  query:%#v, err code:%d, err msg:%s, rid:%s", query, platInfo.Code, platInfo.ErrMsg, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.New(platInfo.Code, platInfo.ErrMsg)})
		return
	}
	iamPlat := extensions.PlatSimplify{
		BKCloudIDField:   platID,
		BKCloudNameField: platInfo.Data.Info[0][common.BKCloudNameField].(string),
	}
	if err := s.AuthManager.UpdateRegisteredPlat(srvData.ctx, srvData.header, iamPlat); err != nil {
		blog.Errorf("UpdatePlat success, but UpdateRegisteredPlat failed, plat: %d, err: %v, rid: %s", platID, err, srvData.rid)
		ccErr := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)}
		_ = resp.WriteError(http.StatusInternalServerError, ccErr)
		return
	}

	// response success
	_ = resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     "",
	})
}

func (s *Service) UpdateHostCloudAreaField(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	// decode request body
	input := metadata.UpdateHostCloudAreaFieldOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("UpdateHostCloudAreaField failed, decode request body failed, err:%+v, rid:%s", err, rid)
		ccErr := srvData.ccErr.Errorf(common.CCErrCommPostInputParseError)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: ccErr})
		return
	}
	if len(input.HostIDs) > common.BKMaxRecordsAtOnce {
		ccErr := srvData.ccErr.CCErrorf(common.CCErrExceedMaxOperationRecordsAtOnce, common.BKMaxRecordsAtOnce)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: ccErr})
		return
	}

	ccErr := s.CoreAPI.CoreService().Host().UpdateHostCloudAreaField(srvData.ctx, srvData.header, input)
	if ccErr != nil {
		blog.ErrorJSON("UpdateHostCloudAreaField failed, core service UpdateHostCloudAreaField failed, input: %s, err: %s, rid: %s", input, ccErr.Error(), rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: ccErr})
		return
	}

	// response success
	_ = resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     "",
	})
}

func (s *Service) findManyCloudAreaAddHostCount(srvData *srvComm, intCloudIDArray []int64, mapCloudIDInfo map[int64]mapstr.MapStr) ([]mapstr.MapStr, error) {
	condition := mapstr.MapStr{
		common.BKCloudIDField: mapstr.MapStr{common.BKDBIN: intCloudIDArray},
	}
	cond := &metadata.QueryCondition{
		Fields:    []string{common.BKCloudIDField},
		Condition: condition,
	}
	rsp, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, cond)
	if nil != err {
		blog.Errorf("findManyCloudAreaAddHostCount htt do error: %v cond:%#v,rid:%s", err, cond, srvData.rid)
		return nil, err
	}
	if false == rsp.Result {
		blog.Errorf("findManyCloudAreaAddHostCount http reply error.  cond:%#v, err code:%d, err msg:%s, rid:%s", cond, rsp.Code, rsp.ErrMsg, srvData.rid)
		return nil, srvData.ccErr.New(rsp.Code, rsp.ErrMsg)
	}

	result := make([]mapstr.MapStr, 0)
	cloudHost := make(map[int64]int64, 0)
	for _, info := range rsp.Data.Info {
		intID, err := info.Int64(common.BKCloudIDField)
		if err != nil {
			blog.ErrorJSON("findManyCloudAreaAddHostCount fail, cloudID convert from interface to int64 failed, err: %v, rid: %s", err, srvData.rid)
			return nil, err
		}
		if _, ok := cloudHost[intID]; !ok {
			cloudHost[intID] = 0
		}
		cloudHost[intID] += 1
	}
	for cloudID, cloudInfo := range mapCloudIDInfo {
		cloudInfo["host_count"] = 0
		for id, count := range cloudHost {
			if cloudID == id {
				cloudInfo["host_count"] = count
				break
			}
		}
		result = append(result, cloudInfo)
	}

	return result, nil
}
