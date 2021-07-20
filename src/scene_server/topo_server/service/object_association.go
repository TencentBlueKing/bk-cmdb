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
	"configcenter/src/common/mapstr"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// CreateObjectAssociation create a new object association
func (s *Service) CreateObjectAssociation(ctx *rest.Contexts) {
	assoc := &metadata.Association{}
	if err := ctx.DecodeInto(assoc); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	var association *metadata.Association
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		association, err = s.Core.AssociationOperation().CreateCommonAssociation(ctx.Kit, assoc)
		if nil != err {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(association)
}

// SearchObjectAssociation search  object association by object id
func (s *Service) SearchObjectAssociation(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if data.Exists("condition") {
		// ATTENTION:
		// compatible with new query structures
		// the new condition format:
		// { "condition":{}}

		cond, err := data.MapStr("condition")
		if nil != err {
			blog.Errorf("search object association, failed to get the condition, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsIsInvalid, err.Error()))
			return
		}

		if len(cond) == 0 {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
			return
		}

		resp, err := s.Core.AssociationOperation().SearchObject(ctx.Kit, &metadata.SearchAssociationObjectRequest{Condition: cond})
		if err != nil {
			blog.Errorf("search object association with cond[%v] failed, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
			return
		}

		if !resp.Result {
			blog.Errorf("search object association with cond[%v] failed, err: %s, rid: %s", cond, resp.ErrMsg, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.New(resp.Code, resp.ErrMsg))
			return
		}

		ctx.RespEntity(resp.Data)
		return
	}

	objID, err := data.String(metadata.AssociationFieldObjectID)
	if err != nil {
		blog.Errorf("search object association, but get object id failed from: %v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	if len(objID) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	resp, err := s.Core.AssociationOperation().SearchObjectAssociation(ctx.Kit, objID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// DeleteObjectAssociation delete object association
func (s *Service) DeleteObjectAssociation(ctx *rest.Contexts) {

	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("delete object association failed, got a invalid object association id[%v], err: %v, rid: %s", ctx.Request.PathParameter("id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoInvalidObjectAssociationID))
		return
	}

	if id <= 0 {
		blog.Errorf("delete object association failed, got a invalid objAsst id[%d], rid: %s", id, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoInvalidObjectAssociationID))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Core.AssociationOperation().DeleteAssociationWithPreCheck(ctx.Kit, id)
		if err != nil {
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

// UpdateObjectAssociation update object association
func (s *Service) UpdateObjectAssociation(ctx *rest.Contexts) {
	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("update object association, but got invalid id[%v], err: %v, rid: %s", ctx.Request.PathParameter("id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommParamsIsInvalid))
		return
	}

	data := new(mapstr.MapStr)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Core.AssociationOperation().UpdateAssociation(ctx.Kit, *data, id)
		if err != nil {
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

// ImportInstanceAssociation import instance  association
func (s *Service) ImportInstanceAssociation(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	request := new(metadata.RequestImportAssociation)
	if err := ctx.DecodeInto(request); err != nil {
		blog.Errorf("ImportInstanceAssociation, json unmarshal error, objID:%S, err: %v, rid:%s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	var ret metadata.ResponeImportAssociationData
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Core.AssociationOperation().ImportInstAssociation(ctx.Kit.Ctx, ctx.Kit, objID, request.AssociationInfoMap, s.Language)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ret)
}
