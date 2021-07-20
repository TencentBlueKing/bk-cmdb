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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

type infoParam struct {
	ExactSearch bool     `json:"exact_search"`
	InnerIP     bool     `json:"bk_host_innerip"`
	OuterIP     bool     `json:"bk_host_outerip"`
	IPList      []string `json:"ip_list"`
}

type queryParams []queryParam
type queryParam struct {
	ObjID    string      `json:"bk_obj_id"`
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	OuterIP  bool        `json:"bk_host_outerip"`
	IPList   []string    `json:"ip_list"`
}

func (s *Service) ListHostFavourites(ctx *rest.Contexts) {
	
	query := new(metadata.QueryInput)
	if err := ctx.DecodeInto(&query); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.CoreAPI.CoreService().Host().ListHostFavourites(ctx.Kit.Ctx, ctx.Kit.User, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("GetHostFavourites http do error,err:%s,input:%+v,rid:%s", err.Error(), query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("GetHostFavourites http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	
	ctx.RespEntity(result.Data)
}

func (s *Service) AddHostFavourite(ctx *rest.Contexts) {

	param := new(metadata.FavouriteParms)
	if err := ctx.DecodeInto(&param); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if param.Name == "" {
		blog.Errorf("add host favorite, but got empty favorite name, param: %+v,rid:%s", param, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostEmptyFavName))
		return
	}

	if param.Info != "" {
		// check if the info string matches the required structure
		err := json.Unmarshal([]byte(param.Info), &infoParam{})
		if err != nil {
			blog.Errorf("AddHostFavourite info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), param.Info, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "info"))
			return
		}
	}
	if param.QueryParams != "" {
		err := json.Unmarshal([]byte(param.QueryParams), &queryParams{})
		if err != nil {
			blog.Errorf("AddHostFavourite info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), param.QueryParams, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "query params"))
			return
		}
	}

	var result *metadata.IDResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		result, err = s.CoreAPI.CoreService().Host().AddHostFavourite(ctx.Kit.Ctx, ctx.Kit.User, ctx.Kit.Header, param)
		if err != nil {
			blog.Errorf("AddHostFavourite http do error,err:%s,input:%+v,rid:%s", err.Error(), param, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("AddHostFavourite http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, param, ctx.Kit.Rid)
			return result.CCError()
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(result.Data)
}

func (s *Service) UpdateHostFavouriteByID(ctx *rest.Contexts) {

	ID := ctx.Request.PathParameter("id")

	if "" == ID || "0" == ID {
		blog.Errorf("update host favourite failed, with id  %d,rid:%s", ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPInputInvalid))
		return
	}

	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if data["name"] == nil || data["name"].(string) == "" {
		blog.Errorf("update host favorite, but got empty name, data: %+v, rid:%s", data, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavUpdateFail))
		return
	}

	if info, exists := data["info"]; exists {
		info := info.(string)
		if info != "" {
			// check if the info string matches the required structure
			err := json.Unmarshal([]byte(info), &infoParam{})
			if err != nil {
				blog.Errorf("AddHostFavourite info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), info, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "info"))
				return
			}
		}
	}
	if queryParam, exists := data["query_params"]; exists {
		queryParam := queryParam.(string)
		if queryParam != "" {
			// check if the info string matches the required structure
			err := json.Unmarshal([]byte(queryParam), &queryParams{})
			if err != nil {
				blog.Errorf("AddHostFavourite info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), queryParam, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "query params"))
				return
			}
		}
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		result, err := s.CoreAPI.CoreService().Host().UpdateHostFavouriteByID(ctx.Kit.Ctx, ctx.Kit.User, ID, ctx.Kit.Header, data)
		if err != nil {
			blog.Errorf("UpdateHostFavouriteByID http do error,err:%s,input:%+v,rid:%s", err.Error(), data, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("UpdateHostFavouriteByID http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, data, ctx.Kit.Rid)
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

func (s *Service) DeleteHostFavouriteByID(ctx *rest.Contexts) {
	
	ID := ctx.Request.PathParameter("id")

	if "" == ID || "0" == ID {
		blog.Errorf("delete host favourite failed, with id  %d,rid:%s", ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPInputInvalid))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		result, err := s.CoreAPI.CoreService().Host().DeleteHostFavouriteByID(ctx.Kit.Ctx, ctx.Kit.User, ID, ctx.Kit.Header)
		if err != nil {
			blog.Errorf("DeleteHostFavouriteByID http do error,err:%s,input:%+v,rid:%s", err.Error(), ID, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("DeleteHostFavouriteByID http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, ID, ctx.Kit.Rid)
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

func (s *Service) IncrHostFavouritesCount(ctx *rest.Contexts) {
	
	ID := ctx.Request.PathParameter("id")
	if "" == ID || "0" == ID {
		blog.Errorf("delete host favourite failed, with id  %s, rid:%s", ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPInputInvalid))
		return
	}

	result, err := s.CoreAPI.CoreService().Host().GetHostFavouriteByID(ctx.Kit.Ctx, ctx.Kit.User, ID, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("IncrHostFavouritesCount GetHostFavouriteByID http do error,err:%s,input:%+v,rid:%s", err.Error(), ID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("IncrHostFavouritesCount GetHostFavouriteByID http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, ID, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}

	count := result.Data.Count + 1
	data := map[string]interface{}{"count": count}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		uResult, err := s.CoreAPI.CoreService().Host().UpdateHostFavouriteByID(ctx.Kit.Ctx, ctx.Kit.User, ID, ctx.Kit.Header, data)
		if err != nil {
			blog.Errorf("IncrHostFavouritesCount UpdateHostFavouriteByID http do error,err:%s,input:%+v,rid:%s", err.Error(), data, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !uResult.Result {
			blog.Errorf("IncrHostFavouritesCount UpdateHostFavouriteByID http response error,err code:%d,err msg:%s,input:%+v,rid:%s", uResult.Code, uResult.ErrMsg, data, ctx.Kit.Rid)
			return uResult.CCError()
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	info := make(map[string]interface{})
	info["id"] = ID
	info["count"] = count
	ctx.RespEntity(info)

}