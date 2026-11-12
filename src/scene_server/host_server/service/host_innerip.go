/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"errors"

	"configcenter/src/ac"
	authmeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	meta "configcenter/src/common/metadata"
)

// modifyHostInnerIPConfigKey the config key of the switch of the host inner ip modification api.
//
// NOTICE: the feature is disabled by default, it has to be turned on manually in common.yaml,
// otherwise the api always returns the CCErrHostModifyInnerIPDisabled error.
const modifyHostInnerIPConfigKey = "hostServer.enableModifyHostInnerIP"

// ModifyHostInnerIP modify the inner ipv4/ipv6 of the host which is using the static addressing mode.
//
// NOTICE: this is a dangerous api, the inner ip is a part of the identity of the host. It is
// disabled by default. The caller must have the edit permission of the host, and the system
// level config management permission which is checked by the ac parser in api server, exactly
// the same as the other config admin apis.
//
// The source ip is the confirmed current value of the host, the modification is rejected if
// the current value does not match it. The destination ip replaces the whole inner ip list.
// The data validation and the audit log are handled by the coreservice api.
func (s *Service) ModifyHostInnerIP(ctx *rest.Contexts) {
	kit := ctx.Kit

	// the feature must be turned on manually. The config is read on each request so that the
	// api can be closed as soon as the config is updated, no service restart is needed.
	// NOTICE: coreservice checks the same switch as well, this check only rejects the disabled
	// request earlier, without starting a transaction.
	enabled, err := cc.Bool(modifyHostInnerIPConfigKey)
	if err != nil {
		blog.Errorf("get config %s failed, err: %v, rid: %s", modifyHostInnerIPConfigKey, err, kit.Rid)
		ctx.RespAutoError(kit.CCError.CCError(common.CCErrHostModifyInnerIPDisabled))
		return
	}
	if !enabled {
		blog.Errorf("config %s is not enabled, host inner ip modification is forbidden, rid: %s",
			modifyHostInnerIPConfigKey, kit.Rid)
		ctx.RespAutoError(kit.CCError.CCError(common.CCErrHostModifyInnerIPDisabled))
		return
	}

	opt := new(meta.ModifyHostInnerIPOpt)
	if err := ctx.DecodeInto(opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	// validate the request shape early for a friendly error; the coreservice api validates
	// again and is authoritative.
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(kit.CCError))
		return
	}

	// the caller must have the edit permission of the host, it is the same as the other host
	// modification apis.
	perm, err := s.authorizeModifyHostInnerIP(kit, opt.HostID)
	if err != nil {
		if errors.Is(err, ac.NoAuthorizeError) {
			ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
			return
		}
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		// the update and its audit log are saved in this transaction by coreservice.
		if _, err := s.CoreAPI.CoreService().Host().UpdateHostInnerIP(kit.Ctx, kit.Header, opt); err != nil {
			blog.Errorf("update the inner ip of host %d failed, err: %v, rid: %s", opt.HostID, err, kit.Rid)
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

// authorizeModifyHostInnerIP checks the edit permission of the host, the global config
// management permission is checked by the ac parser in api server, so that the behavior is
// the same as the other host modification apis, e.g. /updatemany/hosts/all/property.
func (s *Service) authorizeModifyHostInnerIP(kit *rest.Kit, hostID int64) (*meta.BaseResp, error) {
	if err := s.AuthManager.AuthorizeByHostsIDs(kit.Ctx, kit.Header, authmeta.Update, hostID); err != nil {
		if !errors.Is(err, ac.NoAuthorizeError) {
			blog.Errorf("check host authorization failed, host: %d, err: %v, rid: %s", hostID, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommAuthorizeFailed)
		}

		perm, err := s.AuthManager.GenHostBatchNoPermissionResp(kit.Ctx, kit.Header, authmeta.Update,
			[]int64{hostID})
		if err != nil && !errors.Is(err, ac.NoAuthorizeError) {
			blog.Errorf("check host authorization get permission failed, host: %d, err: %v, rid: %s", hostID, err,
				kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommAuthorizeFailed)
		}

		blog.Errorf("host %d is not authorized to be modified, rid: %s", hostID, kit.Rid)
		return perm, ac.NoAuthorizeError
	}

	return nil, nil
}
