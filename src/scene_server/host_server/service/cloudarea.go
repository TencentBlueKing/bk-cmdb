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

	authmeta "configcenter/src/ac/meta"
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

	// set default sort
	if input.Page.Sort == "" {
		input.Page.Sort = "-" + common.CreateTimeField
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
		authorizedPlatIDs, err := s.AuthManager.ListAuthorizedPlatIDs(srvData.ctx, srvData.header, srvData.user)
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
		blog.Errorf("FindManyCloudArea http do error: %v query:%#v,rid:%s", err, query, rid)
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

	retData, err := s.searchPlatAddHostCount(srvData, res.Data)
	if err != nil {
		blog.ErrorJSON("FindManyCloudArea add field host_count failed, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrHostFindManyCloudAreaAddHostCountFieldFail)})
		return
	}

	_ = resp.WriteEntity(metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     retData,
	})
}

// CreatePlatBatch create plat instance in batch
func (s *Service) CreatePlatBatch(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	input := struct {
		Data []mapstr.MapStr `json:"data"`
	}{}
	if err := json.NewDecoder(req.Request.Body).Decode(&input); nil != err {
		blog.Errorf("CreatePlat , but decode body failed, err: %s,rid:%s", err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(input.Data) == 0 {
		blog.Errorf("CreatePlat , input is empty, rid:%s", srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPBodyEmpty)})
		return
	}

	instInfo := &meta.CreateManyModelInstance{
		Datas: input.Data,
	}

	result := make([]metadata.CreateManyCloudAreaElem, len(input.Data))
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		var err error
		res, err := s.CoreAPI.CoreService().Instance().CreateManyInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, instInfo)
		if nil != err {
			blog.Errorf("CreatePlatBatch failed, CreateManyInstance error: %s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
			return srvData.ccErr.Error(common.CCErrTopoInstCreateFailed)
		}

		if false == res.Result {
			blog.Errorf("CreatePlatBatch failed, CreateManyInstance error.err code:%d,err msg:%s,input:%+v,rid:%s", res.Code, res.ErrMsg, input, srvData.rid)
			return errors.New(res.Code, res.ErrMsg)
		}

		if len(res.Data.Exceptions) > 0 {
			blog.Errorf("CreatePlatBatch failed, err:#v,input:%+v,rid:%s", res.Data.Exceptions, input, srvData.rid)
			return srvData.ccErr.New(int(res.Data.Exceptions[0].Code), res.Data.Exceptions[0].Message)
		}

		if len(res.Data.Created) == 0 {
			blog.Errorf("CreatePlatBatch failed, no plat was found,input:%+v,rid:%s", input, srvData.rid)
			return srvData.ccErr.Error(common.CCErrTopoCloudNotFound)
		}

		platIDs := make([]int64, len(res.Data.Created))
		for i, created := range res.Data.Created {
			platIDs[i] = int64(created.ID)
			result[i] = metadata.CreateManyCloudAreaElem{
				CloudID: int64(created.ID),
			}
		}

		// add auditLog
		auditLog := srvData.lgc.NewCloudAreaLog(srvData.ctx, srvData.ownerID)
		if err := auditLog.WithCurrent(srvData.ctx, platIDs...); err != nil {
			blog.ErrorJSON("CreatePlatBatch failed,  WithCurrent err: %v, rid: %s", err, srvData.rid)
			return err
		}

		if err := auditLog.SaveAuditLog(srvData.ctx, metadata.AuditCreate); err != nil {
			blog.ErrorJSON("CreatePlatBatch failed,  SaveAuditLog err: %v, rid: %s", err, srvData.rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
		return
	}

	_ = resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
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
	if err := s.AuthManager.AuthorizeResourceCreate(srvData.ctx, srvData.header, 0, authmeta.Model); err != nil {
		blog.Errorf("check create plat authorization failed, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	instInfo := &meta.CreateModelInstance{
		Data: mapstr.NewFromMap(input),
	}

	var res *metadata.CreatedOneOptionResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		var err error
		res, err = s.CoreAPI.CoreService().Instance().CreateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, instInfo)
		if nil != err {
			blog.Errorf("CreatePlat error: %s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
			return srvData.ccErr.Error(common.CCErrTopoInstCreateFailed)
		}

		if false == res.Result {
			blog.Errorf("GetPlat error.err code:%d,err msg:%s,input:%+v,rid:%s", res.Code, res.ErrMsg, input, srvData.rid)
			return errors.New(res.Code, res.ErrMsg)
		}

		platID := int64(res.Data.Created.ID)

		// add auditLog
		auditLog := srvData.lgc.NewCloudAreaLog(srvData.ctx, srvData.ownerID)
		if err := auditLog.WithCurrent(srvData.ctx, platID); err != nil {
			blog.ErrorJSON("createPlat success., but add auditLog fail, err: %v, rid: %s", err, srvData.rid)
			return err
		}

		if err := auditLog.SaveAuditLog(srvData.ctx, metadata.AuditCreate); err != nil {
			blog.ErrorJSON("createPlat success., but add auditLog fail, err: %v, rid: %s", err, srvData.rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
		return
	}

	_ = resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res.Data,
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
		blog.Errorf("DelPlat failed, can't delete default cloud area, input:%+v,rid:%s", platID, srvData.rid)
		// can't delete default cloud area
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrDeleteDefaultCloudAreaFail)})
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

	// add auditLog preData
	auditLog := srvData.lgc.NewCloudAreaLog(srvData.ctx, srvData.ownerID)
	if err := auditLog.WithPrevious(srvData.ctx, platID); err != nil {
		blog.ErrorJSON("DelPlat success., but add auditLog fail, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	delCond := &meta.DeleteOption{
		Condition: mapstr.MapStr{common.BKCloudIDField: platID},
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		res, err := s.CoreAPI.CoreService().Instance().DeleteInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, delCond)
		if nil != err {
			blog.Errorf("DelPlat do error: %v, input:%d,rid:%s", err, platID, srvData.rid)
			return srvData.ccErr.Errorf(common.CCErrTopoInstDeleteFailed)
		}

		if false == res.Result {
			blog.Errorf("DelPlat http response error. err code:%d,err msg:%s,input:%s,rid:%s", res.Code, res.ErrMsg, platID, srvData.rid)
			return srvData.ccErr.New(res.Code, res.ErrMsg)
		}

		if err := auditLog.SaveAuditLog(srvData.ctx, metadata.AuditDelete); err != nil {
			blog.ErrorJSON("DelPlat success., but add auditLog fail, err: %v, rid: %s", err, srvData.rid)
			return srvData.ccErr.New(res.Code, res.ErrMsg)
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
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
		CloudName   string `json:"bk_cloud_name"`
		CloudVendor string `json:"bk_cloud_vendor"`
		Region      string `json:"bk_region"`
	}{}

	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("UpdatePlat failed, err:%+v, rid:%s", err, srvData.rid)
		ccErr := srvData.ccErr.Errorf(common.CCErrCommJSONUnmarshalFailed)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: ccErr})
		return
	}

	// auditLog preData
	auditLog := srvData.lgc.NewCloudAreaLog(srvData.ctx, srvData.ownerID)
	if err := auditLog.WithPrevious(srvData.ctx, platID); err != nil {
		blog.ErrorJSON("DelPlat success., but add auditLog fail, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	// update plat
	user := util.GetUser(req.Request.Header)

	toUpdate := mapstr.MapStr{
		common.BKLastEditor: user,
	}

	if len(input.CloudVendor) != 0 {
		toUpdate[common.BKCloudVendor] = input.CloudVendor
	}

	if len(input.Region) != 0 {
		toUpdate[common.BKRegion] = input.Region
	}

	if len(input.CloudName) != 0 {
		toUpdate[common.BKCloudNameField] = input.CloudName
	}

	updateOption := &meta.UpdateOption{
		Data: toUpdate,
		Condition: map[string]interface{}{
			common.BKCloudIDField: platID,
		},
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		res, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, updateOption)
		if nil != err {
			blog.ErrorJSON("UpdatePlat failed, UpdateInstance failed, input:%s, err:%s, rid:%s", updateOption, err.Error(), srvData.rid)
			return srvData.ccErr.Errorf(common.CCErrTopoInstDeleteFailed)
		}
		if res.Result == false || res.Code != 0 {
			blog.ErrorJSON("UpdatePlat failed, UpdateInstance failed, input:%s, response:%s, rid:%s", updateOption, res, srvData.rid)
			return errors.New(res.Code, res.ErrMsg)
		}

		// update auditLog
		if err := auditLog.WithCurrent(srvData.ctx, platID); err != nil {
			blog.ErrorJSON("UpdatePlat success., but add auditLog fail, err: %v, rid: %s", err, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if err := auditLog.SaveAuditLog(srvData.ctx, metadata.AuditUpdate); err != nil {
			blog.ErrorJSON("UpdatePlat success., but add auditLog fail, err: %v, rid: %s", err, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
		return
	}

	// response success
	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
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

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		ccErr := s.CoreAPI.CoreService().Host().UpdateHostCloudAreaField(srvData.ctx, srvData.header, input)
		if ccErr != nil {
			blog.ErrorJSON("UpdateHostCloudAreaField failed, core service UpdateHostCloudAreaField failed, input: %s, err: %s, rid: %s", input, ccErr.Error(), rid)
			return ccErr
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
		return
	}

	// response success
	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) searchPlatAddHostCount(srvData *srvComm, dataInfo metadata.InstDataInfo) (map[string]interface{}, error) {
	// add host_count
	mapCloudIDInfo := make(map[int64]mapstr.MapStr, 0)
	intCloudIDArray := make([]int64, 0)
	for _, area := range dataInfo.Info {
		intCloudID, err := area.Int64(common.BKCloudIDField)
		if err != nil {
			blog.ErrorJSON("FindManyCloudArea fail with cloudID convert from interface to int64 failed, err: %v, rid: %s", err, srvData.rid)
			return nil, err
		}
		intCloudIDArray = append(intCloudIDArray, intCloudID)
		mapCloudIDInfo[intCloudID] = area
	}

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
		if count, ok := cloudHost[cloudID]; ok {
			cloudInfo["host_count"] = count
		}
		result = append(result, cloudInfo)
	}

	retData := map[string]interface{}{
		"info":  result,
		"count": dataInfo.Count,
	}

	return retData, nil
}
