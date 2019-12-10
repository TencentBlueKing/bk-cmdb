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

	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// FindManyCloudArea  find cloud area list
func (s *Service) FindManyCloudArea(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid
	input := new(metadata.CloudAreaParameter)
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
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommPageLimitIsExceeded)})
		return
	}

	sortArr := make([]metadata.SearchSort, 0)
	if len(input.Page.Sort) != 0 {
		for _, field := range strings.Split(input.Page.Sort, ",") {
			sortArr = append(sortArr, metadata.SearchSort{
				IsDsc: true,
				Field: field,
			})
		}
	}

	// auth: get authorized resources
	authorizedPlatIDs, err := s.AuthManager.ListAuthorizedPlatIDs(srvData.ctx, srvData.user)
	if err != nil {
		blog.Errorf("FindManyCloudArea failed, ListAuthorizedPlatIDs failed, err: %+v, rid: %s", srvData.user, rid)
		_ = resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommListAuthorizedResourceFromIAMFailed)})
		return
	}

	filter := map[string]interface{}{
		common.BKDBAND: []map[string]interface{}{
			input.Condition,
			{
				common.BKCloudIDField: map[string]interface{}{
					common.BKDBIN: authorizedPlatIDs,
				},
			},
		},
	}
	query := &metadata.QueryCondition{
		Condition: filter,
		Limit: metadata.SearchLimit{
			Limit:  int64(input.Page.Limit),
			Offset: int64(input.Page.Start),
		},
		SortArr: sortArr,
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

	retData := map[string]interface{}{
		"info":  res.Data.Info,
		"count": res.Data.Count,
	}

	_ = resp.WriteEntity(metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     retData,
	})
}

// CreatePlat create a plat instance
// available fields for body are last_time, bk_cloud_name, bk_supplier_account, bk_cloud_id, create_time
// {"bk_cloud_name": "云区域", "bk_supplier_account": 0}
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

	// auth: check authorization
	if err := s.AuthManager.AuthorizeResourceCreate(srvData.ctx, srvData.header, 0, authmeta.Plat); err != nil {
		blog.Errorf("check create plat authorization failed, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	instInfo := &meta.CreateModelInstance{
		Data: mapstr.NewFromMap(input),
	}

	res, err := s.CoreAPI.CoreService().Instance().CreateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, instInfo)
	if nil != err {
		blog.Errorf("CreatePlat error: %s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoInstCreateFailed)})
		return
	}

	if false == res.Result {
		blog.Errorf("GetPlat error.err code:%d,err msg:%s,input:%+v,rid:%s", res.Code, res.ErrMsg, input, srvData.rid)
		_ = resp.WriteHeaderAndJson(http.StatusInternalServerError, res, "application/json")
		return
	}

	// register plat to iam
	platID := int64(res.Data.Created.ID)
	if err := s.AuthManager.RegisterPlatByID(srvData.ctx, srvData.header, platID); err != nil {
		blog.Errorf("CreatePlat failed, RegisterPlatByID failed, err: %s, rid:%s", err.Error(), srvData.rid)
		ccErr := srvData.ccErr.CCError(common.CCErrCommRegistResourceToIAMFailed)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: ccErr})
		return
	}

	_ = resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res.Data,
	})

}

func (s *Service) DelPlat(req *restful.Request, resp *restful.Response) {
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
	input := struct {
		CloudName string `json:"bk_cloud_name"`
	}{}
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("UpdatePlat failed, err:%+v, rid:%s", err, srvData.rid)
		ccErr := srvData.ccErr.Errorf(common.CCErrCommJSONUnmarshalFailed)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: ccErr})
		return
	}

	// update plat
	updateOption := &meta.UpdateOption{
		Data: map[string]interface{}{
			common.BKCloudNameField: input.CloudName,
		},
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
		blog.ErrorJSON("UpdatePlat failed, UpdateInstance failed, input:%s, response:%s, err:%s, rid:%s", updateOption, res, err.Error(), srvData.rid)
		ccErr := &meta.RespError{Msg: errors.New(res.Code, res.ErrMsg)}
		_ = resp.WriteError(http.StatusInternalServerError, ccErr)
		return
	}

	// auth: sync resource info to iam
	iamPlat := extensions.PlatSimplify{
		BKCloudIDField:   platID,
		BKCloudNameField: input.CloudName,
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
