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
	"configcenter/src/ac"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (s *Service) LockHost(ctx *rest.Contexts) {

	input := &metadata.HostLockRequest{}

	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if 0 == len(input.IDS) {
		blog.Errorf("lock host, id_list is empty,input:%+v, rid:%s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "id_list"))
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, input.IDS...); err != nil {
		if err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v", input.IDS, err)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, input.IDS)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logic.LockHost(ctx.Kit, input)
		if nil != err {
			blog.Errorf("lock host, handle host lock error, error:%s, input:%+v,rid:%s", err.Error(), input, ctx.Kit.Rid)
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

func (s *Service) UnlockHost(ctx *rest.Contexts) {

	input := &metadata.HostLockRequest{}
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}
	if 0 == len(input.IDS) {
		blog.Errorf("unlock host, id_list is empty, input:%+v,rid:%s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "id_list"))
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, input.IDS...); err != nil {
		if err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", input.IDS, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, input.IDS)
		if err != nil {
			blog.Errorf("gen no permission response failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logic.UnlockHost(ctx.Kit, input)
		if nil != err {
			blog.Errorf("unlock host, handle host unlock error, error:%s, input:%+v,rid:%s", err.Error(), input, ctx.Kit.Rid)
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

func (s *Service) QueryHostLock(ctx *rest.Contexts) {

	input := &metadata.QueryHostLockRequest{}
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if 0 == len(input.IDS) {
		blog.Errorf("query lock host, id_list is empty, input:%+v,rid:%s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "id_list"))
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, input.IDS...); err != nil {
		if err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", input.IDS, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, input.IDS)
		if err != nil {
			blog.Errorf("gen no permission response failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	hostLockInfos, err := s.Logic.QueryHostLock(ctx.Kit, input)
	if nil != err {
		blog.Errorf("query lock host, handle query host lock error, error:%s, input:%+v,rid:%s", err.Error(), input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(hostLockInfos)
}
