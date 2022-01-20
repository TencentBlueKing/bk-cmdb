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
	"time"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/scene_server/event_server/sync/hostidentifier"

	"github.com/tidwall/gjson"
)

// SyncHostIdentifier sync host identifier, add hostInfo message to redis fail host list
func (s *Service) SyncHostIdentifier(ctx *rest.Contexts) {
	if s.SyncData == nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommInternalServerError,
			"no start up sync hostIdentifier ability"))
		return
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

	if auth.EnableAuthorize() {
		if err := s.haveAuthority(ctx, hostIDArray.HostIDs); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	hosts, err := s.getHostInfo(ctx, hostIDArray.HostIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	task, err := s.SyncData.BatchSyncHostIdentifier(hosts.Info)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	var resultMap map[string]string
	endTime := time.Now().Add(30 * time.Second).Unix()

	for time.Now().Unix() < endTime {
		resultMap, err = s.SyncData.GetTaskExecutionResultMap(task)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		// 该任务包含的主机拿到全部的结果
		if len(task.HostInfos) == len(resultMap) {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if resultMap == nil || time.Now().Unix() >= endTime {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrEventGetTaskStatusTimeout))
		return
	}

	failHostIDs := make([]int64, 0)
	for _, hostInfo := range task.HostInfos {
		key := hostidentifier.HostKey(strconv.FormatInt(hostInfo.CloudID, 10), hostInfo.HostInnerIP)
		// 把推送失败且没超过最大重试次数的主机放到失败主机队列中
		if gjson.Get(resultMap[key], "error_code").Int() != common.CCSuccess {
			failHostIDs = append(failHostIDs, hostInfo.HostID)
		}
	}
	if len(failHostIDs) != 0 {
		ctx.RespEntityWithError(failHostIDs, ctx.Kit.CCError.CCErrorf(common.CCErrEventPushHostIdentifierFailed))
		return
	}

	ctx.RespEntity(nil)
}

func (s *Service) haveAuthority(ctx *rest.Contexts, hostIDs []int64) error {
	businessIDInfo, err := s.getHostBusinessIDInfo(ctx, hostIDs)
	if err != nil {
		return err
	}

	businessIDs := make([]int64, len(businessIDInfo.Info))
	for _, info := range businessIDInfo.Info {
		businessIDs = append(businessIDs, info.AppID)
	}

	err = s.AuthManager.AuthorizeByBusinessID(ctx.Kit.Ctx, ctx.Kit.Header, meta.ViewBusinessResource, businessIDs...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) getHostBusinessIDInfo(ctx *rest.Contexts, hostIDs []int64) (*metadata.HostConfigData, error) {
	cond := &metadata.HostModuleRelationRequest{
		HostIDArr: hostIDs,
		Fields:    []string{common.BKAppIDField},
	}
	result, err := s.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, cond)
	if err != nil {
		blog.Errorf("http do error, err: %v, input: %v, rid: %s", err, cond, ctx.Kit.Rid)
		return nil, err
	}
	if result.Count == 0 {
		blog.Errorf("get host module relation success, but result count is 0, hostIDs: %v, rid: %s",
			hostIDs, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommGetBusinessIDByHostIDFailed)
	}
	return result, nil
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
