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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (s *Service) SaveUserCustom(ctx *rest.Contexts) {

	params := make(map[string]interface{})
	if err := ctx.DecodeInto(&params); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.CoreAPI.CoreService().Host().GetUserCustomByUser(ctx.Kit.Ctx, ctx.Kit.User, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("SaveUserCustom GetUserCustomByUser http do error,err:%s,input:%s, rid:%s", err.Error(), params, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("SaveUserCustom GetUserCustomByUser http response error,err code:%d,err msg:%s,input:%s, rid:%s", result.Code, result.ErrMsg, params, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}

	var res *metadata.BaseResp
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		if len(result.Data) == 0 {
			res, err = s.CoreAPI.CoreService().Host().AddUserCustom(ctx.Kit.Ctx, ctx.Kit.User, ctx.Kit.Header, params)
			if err != nil {
				blog.Errorf("SaveUserCustom AddUserCustom http do error,err:%s,input:%s, rid:%s", err.Error(), params, ctx.Kit.Rid)
				return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			if !res.Result {
				blog.Errorf("SaveUserCustom AddUserCustom http response error,err code:%d,err msg:%s,input:%s, rid:%s", res.Code, res.ErrMsg, params, ctx.Kit.Rid)
				return res.CCError()
			}
			return nil

		}
		id := result.Data["id"].(string)
		res, err = s.CoreAPI.CoreService().Host().UpdateUserCustomByID(ctx.Kit.Ctx, ctx.Kit.User, id, ctx.Kit.Header, params)
		if err != nil {
			blog.Errorf("SaveUserCustom UpdateUserCustomByID http do error,err:%s,input:%s, rid:%s", err.Error(), params, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !res.Result {
			blog.Errorf("SaveUserCustom UpdateUserCustomByID http response error,err code:%d,err msg:%s,input:%s, rid:%s", res.Code, res.ErrMsg, params, ctx.Kit.Rid)
			return res.CCError()
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (s *Service) GetUserCustom(ctx *rest.Contexts) {

	result, err := s.CoreAPI.CoreService().Host().GetUserCustomByUser(ctx.Kit.Ctx, ctx.Kit.User, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("GetUserCustom http do error,err:%s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("GetUserCustom http response error,err code:%d,err msg:%s, rid:%s", result.Code, result.ErrMsg, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}

	ctx.RespEntity(result.Data)
}

// GetModelDefaultCustom 获取模型在列表页面展示字段
func (s *Service) GetModelDefaultCustom(ctx *rest.Contexts) {

	result, err := s.CoreAPI.CoreService().Host().GetDefaultUserCustom(ctx.Kit.Ctx, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("GetDefaultCustom http do error,err:%s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("GetDefaultCustom http response error,err code:%d,err msg:%s, rid:%s", result.Code, result.ErrMsg, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	// ensure return {} by json decode
	if result.Data == nil {
		result.Data = make(map[string]interface{}, 0)
	}

	ctx.RespEntity(result.Data)
}

// SaveModelDefaultCustom 设置模型在列表页面展示字段
func (s *Service) SaveModelDefaultCustom(ctx *rest.Contexts) {

	objID := ctx.Request.PathParameter("obj_id")

	input := make(map[string]interface{})
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if len(input) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPBodyEmpty))
		return
	}

	userCustomInput := make(map[string]interface{}, 0)
	// add prefix all key
	for key, val := range input {
		userCustomInput[fmt.Sprintf("%s_%s", objID, key)] = val
	}

	var result *metadata.BaseResp
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		result, err = s.CoreAPI.CoreService().Host().UpdateDefaultUserCustom(ctx.Kit.Ctx, ctx.Kit.Header, userCustomInput)
		if err != nil {
			blog.ErrorJSON("SaveUserCustom GetUserCustomByUser http do error,err:%s,input:%s, rid:%s", err.Error(), input, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if err := result.CCError(); err != nil {
			blog.ErrorJSON("SaveUserCustom GetUserCustomByUser http reply error. result: %s, input: %s, rid: %s", result, input, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}
