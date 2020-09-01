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
	"fmt"
	"strconv"
	"time"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	meta "configcenter/src/common/metadata"
	parser "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
)

func (s *Service) AddUserCustomQuery(ctx *rest.Contexts) {

	ucq := new(meta.UserConfig)
	if err := ctx.DecodeInto(&ucq); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if "" == ucq.Name {
		blog.Error("AddUserCustomQuery add user custom query parameter name is required,input:%+v,rid:%s", ucq, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "name"))
		return
	}

	if "" == ucq.Info {
		blog.Error("AddUserCustomQuery add user custom query info is required,input:%+v,rid:%s", ucq, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "info"))
		return
	}
	// check if the info string matches the required structure
	err := json.Unmarshal([]byte(ucq.Info), &meta.HostCommonSearch{})
	if err != nil {
		blog.Errorf("AddUserCustomQuery info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), ucq, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	if 0 >= ucq.AppID {
		blog.Error("AddUserCustomQuery add user custom query parameter ApplicationID is required,input:%+v,rid:%s", ucq, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField))
		return
	}

	ucq.CreateUser = ctx.Kit.User

	var result *meta.IDResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		result, err = s.CoreAPI.CoreService().Host().AddUserConfig(ctx.Kit.Ctx, ctx.Kit.Header, ucq)
		if err != nil {
			blog.Errorf("GetUserCustom http do error, err:%s, input:%+v,rid:%s", err.Error(), ucq, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("GetUserCustom http response error, err code:%d,err msg:%s, input:%+v,rid:%s", result.Code, result.ErrMsg, ucq, ctx.Kit.Rid)
			return result.CCError()
		}

		// register custom query resource creator action to iam
		if auth.EnableAuthorize() {
			res, err := s.CoreAPI.CoreService().Host().GetUserConfigDetail(ctx.Kit.Ctx, strconv.FormatInt(ucq.AppID, 10), result.Data.ID, ctx.Kit.Header)
			if err != nil {
				blog.Errorf("get created custom query failed, err: %s, biz: %d, ID: %s, rid: %s", err.Error(), ucq.AppID, result.Data.ID, ctx.Kit.Rid)
				return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			if !res.Result {
				blog.Errorf("get created custom query failed, err: %s, biz: %d, ID: %s, rid: %s", res.ErrMsg, ucq.AppID, result.Data.ID, ctx.Kit.Rid)
				return res.CCError()
			}

			iamInstance := meta.IamInstanceWithCreator{
				Type:    string(iam.BizCustomQuery),
				ID:      res.Data.ID,
				Name:    res.Data.Name,
				Creator: ctx.Kit.User,
			}
			_, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created custom query to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(result.Data)
}

func (s *Service) UpdateUserCustomQuery(ctx *rest.Contexts) {

	params := make(map[string]interface{})
	req := ctx.Request
	if err := json.NewDecoder(req.Request.Body).Decode(&params); nil != err {
		blog.Errorf("update user custom query failed with decode body err: %v,rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	params["modify_user"] = ctx.Kit.User
	params[common.LastTimeField] = time.Now().UTC()

	if info, exists := params["info"]; exists {
		info := info.(string)
		if len(info) != 0 {
			// check if the info string matches the required structure
			err := json.Unmarshal([]byte(info), &meta.HostCommonSearch{})
			if err != nil {
				blog.Errorf("UpdateUserCustomQuery info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), params, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
				return
			}
		}
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		bizID := req.PathParameter("bk_biz_id")
		result, err := s.CoreAPI.CoreService().Host().UpdateUserConfig(ctx.Kit.Ctx, bizID, req.PathParameter("id"), ctx.Kit.Header, params)
		if err != nil {
			blog.Errorf("UpdateUserCustomQuery http do error,err:%s, biz:%v,input:%+v,rid:%s", err.Error(), bizID, params, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("UpdateUserCustomQuery http response error,err code:%d,err msg:%s, bizID:%v,input:%+v,rid:%s", result.Code, result.ErrMsg, bizID, params, ctx.Kit.Rid)
			return result.CCError()
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

func (s *Service) DeleteUserCustomQuery(ctx *rest.Contexts) {

	req := ctx.Request
	dynamicID := req.PathParameter("id")
	appID := req.PathParameter("bk_biz_id")

	dyResult, err := s.CoreAPI.CoreService().Host().GetUserConfigDetail(ctx.Kit.Ctx, appID, dynamicID, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("DeleteUserCustomQuery http do error,err:%s, biz:%v, rid:%s", err.Error(), appID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	if !dyResult.Result {
		blog.Errorf("DeleteUserCustomQuery http response error,err code:%d,err msg:%s, bizID:%v,rid:%s", dyResult.Code, dyResult.ErrMsg, appID, ctx.Kit.Rid)
		ctx.RespAutoError(dyResult.CCError())
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		result, err := s.CoreAPI.CoreService().Host().DeleteUserConfig(ctx.Kit.Ctx, appID, dynamicID, ctx.Kit.Header)
		if err != nil {
			blog.Errorf("DeleteUserCustomQuery http do error,err:%s, biz:%v, rid:%s", err.Error(), appID, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("DeleteUserCustomQuery http response error,err code:%d,err msg:%s, bizID:%v,rid:%s", result.Code, result.ErrMsg, appID, ctx.Kit.Rid)
			return result.CCError()
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (s *Service) GetUserCustomQuery(ctx *rest.Contexts) {

	req := ctx.Request
	input := &meta.QueryInput{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); nil != err {
		blog.Errorf("get user custom query failed with decode body err: %v,rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	var condition map[string]interface{}
	if nil != input.Condition {
		condition = input.Condition
	} else {
		condition = make(map[string]interface{})
	}
	// if name in condition , add like search
	name, ok := condition["name"].(string)
	if ok && "" != name {
		condition["name"] = common.KvMap{common.BKDBLIKE: parser.SpecialCharChange(name)}
	}

	var err error
	condition[common.BKAppIDField], err = util.GetInt64ByInterface(req.PathParameter("bk_biz_id"))
	if nil != err {
		blog.Error("GetUserCustomQuery query user custom query parameter ApplicationID not integer in url,bizID:%s,rid:%s", req.PathParameter("bk_biz_id"), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	input.Condition = condition

	result, err := s.CoreAPI.CoreService().Host().GetUserConfig(ctx.Kit.Ctx, ctx.Kit.Header, input)
	if err != nil {
		blog.Errorf("GetUserCustomQuery http do error,err:%s, biz:%v,input:%+v,rid:%s", err.Error(), req.PathParameter("bk_biz_id"), input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("GetUserCustomQuery http response error,err code:%d,err msg:%s, bizID:%v,input:%+v,rid:%s", result.Code, result.ErrMsg, req.PathParameter("bk_biz_id"), input, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}

	ctx.RespEntity(result.Data)
}

func (s *Service) GetUserCustomQueryDetail(ctx *rest.Contexts) {

	req := ctx.Request
	appID := req.PathParameter("bk_biz_id")
	ID := req.PathParameter("id")

	result, err := s.CoreAPI.CoreService().Host().GetUserConfigDetail(ctx.Kit.Ctx, appID, ID, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("GetUserCustomQueryDetail http do error,err:%s, biz:%v,ID:%+v,rid:%s", err.Error(), appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("GetUserCustomQueryDetail http response error,err code:%d,err msg:%s, bizID:%v,ID:%+v,rid:%s", result.Code, result.ErrMsg, appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}

	ctx.RespEntity(result.Data)
}

func (s *Service) GetUserCustomQueryResult(ctx *rest.Contexts) {

	req := ctx.Request
	appID := req.PathParameter("bk_biz_id")
	ID := req.PathParameter("id")

	intAppID, err := util.GetInt64ByInterface(appID)
	if nil != err {
		blog.Errorf("UserAPIResult custom query failed,  err: %v, appid: %s, id:%s, logID:%s", err.Error(), appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID"))
		return
	}

	result, err := s.CoreAPI.CoreService().Host().GetUserConfigDetail(ctx.Kit.Ctx, appID, ID, ctx.Kit.Header)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("UserAPIResult custom query failed,  err: %v, appid:%s, id:%s, rid: %s", err.Error(), appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrGetUserCustomQueryDetailFailed, err.Error()))
		return
	}

	if "" == result.Data.Name {
		blog.Errorf("UserAPIResult custom query not found, appid:%s, id:%s, logID:%s", appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommNotFound))
		return
	}

	var input meta.HostCommonSearch
	input.AppID = intAppID

	err = json.Unmarshal([]byte(result.Data.Info), &input)
	if nil != err {
		blog.Errorf("UserAPIResult custom unmarshal failed,  err: %v, appid:%s, id:%s, logID:%s", err.Error(), appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	input.Page.Start, err = util.GetIntByInterface(req.PathParameter("start"))
	if err != nil {
		blog.Errorf("UserAPIResult start invalid, err: %v, appid:%s, id:%s, logID:%s", err.Error(), appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "start"))
		return
	}
	input.Page.Limit, err = util.GetIntByInterface(req.PathParameter("limit"))
	if err != nil {
		blog.Errorf("UserAPIResult limit invalid, err: %v, appid:%s, id:%s, logID:%s", err.Error(), appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "limit"))
		return
	}
	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	retData, err := lgc.SearchHost(ctx.Kit.Ctx, &input, false)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("UserAPIResult custom query search host failed,  err: %v, appid:%s, id:%s, rid: %s", err.Error(), appID, ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrGetUserCustomQueryDetailFailed, err.Error()))
		return
	}

	ctx.RespEntity(meta.SearchHost{
		Count: retData.Count,
		Info:  retData.Info,
	})
	return
}
