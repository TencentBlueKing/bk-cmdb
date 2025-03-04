/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package tenantset

import (
	"encoding/json"
	"strconv"

	"configcenter/pkg/filter"
	tenantset "configcenter/pkg/types/tenant-set"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	ccjson "configcenter/src/common/json"
	"configcenter/src/common/metadata"
)

// ListTenantSet list tenant set
func (s *service) ListTenantSet(ctx *rest.Contexts) {
	// only system tenant can access this api
	if ctx.Kit.TenantID != common.BKDefaultTenantID {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.TenantID))
		return
	}

	opt := new(metadata.CommonQueryOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	if authResp, authorized := s.AuthManager.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.TenantSet, Action: meta.Find}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	// check if filter matches default tenant set, right now we only support default tenant set
	if opt.Filter != nil {
		matched, err := opt.Filter.Match(filter.JsonString(tenantset.DefaultTenantSetStr))
		if err != nil {
			blog.Errorf("check if filter(%s) matches failed, err: %v, rid: %s", opt.Filter, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if !matched {
			ctx.RespEntityWithCount(0, make([]interface{}, 0))
			return
		}
	}

	if opt.Page.EnableCount {
		ctx.RespEntityWithCount(1, make([]interface{}, 0))
		return
	}

	if opt.Page.Start > 0 {
		ctx.RespEntityWithCount(0, make([]interface{}, 0))
		return
	}

	tenantSet := tenantset.DefaultTenantSetStr
	if len(opt.Fields) > 0 {
		tenantSet = *ccjson.CutJsonDataWithFields(&tenantSet, opt.Fields)
	}

	ctx.RespEntityWithCount(0, []json.RawMessage{json.RawMessage(tenantSet)})
}

// ListTenantInTenantSet list tenants in tenant set
func (s *service) ListTenantInTenantSet(ctx *rest.Contexts) {
	// only system tenant can access this api
	if ctx.Kit.TenantID != common.BKDefaultTenantID {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.TenantID))
		return
	}

	tenantSetIDStr := ctx.Request.PathParameter("tenant_set_id")
	tenantSetID, err := strconv.ParseInt(tenantSetIDStr, 10, 64)
	if err != nil {
		blog.Errorf("parse tenant set id: %s failed, err: %v, rid: %s", tenantSetIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// check if tenant set is default tenant set, right now we only support default tenant set
	if tenantSetID != tenantset.DefaultTenantSetID {
		ctx.RespEntity(tenantset.ListTenantResult{Data: make([]tenantset.Tenant, 0)})
		return
	}

	// authorize
	if authResp, authorized := s.AuthManager.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.TenantSet, Action: meta.AccessTenantSet, InstanceID: tenantSetID}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	allTenants, err := s.ClientSet.CoreService().Tenant().GetAllTenants(ctx.Kit.Ctx, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("get all tenants failed, err: %v, rid: %s", tenantSetIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	tenants := make([]tenantset.Tenant, 0)
	for _, tenant := range allTenants {
		tenants = append(tenants, tenantset.Tenant{
			ID:     tenant.TenantID,
			Status: string(tenant.Status),
		})
	}

	ctx.RespEntity(tenantset.ListTenantResult{Data: tenants})
}
