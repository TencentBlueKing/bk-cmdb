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
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

// SyncHostIdentifier sync host identifier, add hostInfo message to redis fail host list
func (s *Service) SyncHostIdentifier(ctx *rest.Contexts) {
	if s.SyncData == nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommInternalServerError,
			"no start up sync hostIdentifier ability"))
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)

	hostIDArray := new(metadata.HostIDArray)
	if err := ctx.DecodeInto(&hostIDArray); err != nil {
		blog.Errorf("decode request body err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	rawErr := hostIDArray.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	hostModuleRelation, err := s.getHostModuleRelation(ctx, hostIDArray.HostIDs)
	if err != nil {
		ctx.RespAutoError(err)
	}

	if auth.EnableAuthorize() {
		if err := s.haveAuthority(ctx, hostModuleRelation); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	hosts, err := s.getHostInfo(ctx, hostIDArray.HostIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	s.SyncData.BatchSyncHostIdentifier(hosts.Info)
	ctx.RespEntity(nil)
}

func (s *Service) getHostModuleRelation(ctx *rest.Contexts, hostIDs []int64) (*metadata.HostConfigData, error) {
	cond := &metadata.HostModuleRelationRequest{
		HostIDArr: hostIDs,
		Fields:    []string{common.BKAppIDField},
	}
	result, err := s.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, cond)
	if err != nil {
		blog.Errorf("http do error, err: %v, input: %v, rid: %s", err, cond, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if result.Count == 0 {
		blog.Errorf("get host module relation success, but result count is 0, hostIDs: %v, rid: %s",
			hostIDs, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrHostGetFail)
	}
	return result, nil
}

func (s *Service) haveAuthority(ctx *rest.Contexts, hostConfigData *metadata.HostConfigData) error {
	authInput := meta.ListAuthorizedResourcesParam{
		UserName:     ctx.Kit.User,
		ResourceType: meta.Business,
		Action:       meta.Find,
	}
	authorizedRes, err := s.GetAuthorizer().ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header, authInput)
	if err != nil {
		blog.Errorf("list authorized resources failed, user: %s, err: %v, rid: %s", ctx.Kit.User, err, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrorTopoGetAuthorizedBusinessListFailed)
	}

	if !authorizedRes.IsAny {
		authorizedBizList := make([]int64, 0)
		for _, resourceID := range authorizedRes.Ids {
			bizID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				blog.Errorf("parse bizID failed, val: %s, err: %v, rid: %s", bizID, err, ctx.Kit.Rid)
				return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
			}
			authorizedBizList = append(authorizedBizList, bizID)
		}

		for _, info := range hostConfigData.Info {
			if !util.InArray(info.AppID, authorizedBizList) {
				return ctx.Kit.CCError.CCError(common.CCErrCommCheckAuthorizeFailed)
			}
		}
	}
	return nil
}

func (s *Service) getHostInfo(ctx *rest.Contexts, hostIDs []int64) (*metadata.ListHostResult, error) {
	options := &metadata.ListHosts{
		Fields: []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField},
		HostPropertyFilter: &querybuilder.QueryFilter{
			Rule: querybuilder.CombinedRule{
				Condition: querybuilder.ConditionAnd,
				Rules: []querybuilder.Rule{
					querybuilder.AtomRule{
						Field:    common.BKHostIDField,
						Operator: querybuilder.OperatorIn,
						Value:    hostIDs,
					},
				},
			},
		},
	}
	hosts, err := s.engine.CoreAPI.CoreService().Host().ListHosts(ctx.Kit.Ctx, ctx.Kit.Header, options)
	if err != nil {
		blog.Errorf("http do error, hostIDs: %v, err: %v, rid: %s", hostIDs, err, ctx.Kit.Rid)
		return nil, err
	}
	if hosts.Count == 0 {
		blog.Errorf("get host success, but host count is 0, hostIDs: %v, rid: %s", hostIDs, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrHostGetFail)
	}
	return hosts, nil
}
