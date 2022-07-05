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
	"errors"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/sync/hostidentifier"

	"github.com/tidwall/gjson"
)

// SyncHostIdentifier sync host identifier, add hostInfo message to redis fail host list
func (s *Service) SyncHostIdentifier(ctx *rest.Contexts) {
	if s.SyncData == nil {
		blog.Errorf("sync host identifier disabled, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrEventSyncHostIdentifierDisabled))
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
		blog.Errorf("exceed max limit number, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	if auth.EnableAuthorize() {
		if err := s.haveAuthority(ctx.Kit, hostIDArray.HostIDs); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	hosts, err := s.getHostInfo(ctx.Kit, hostIDArray.HostIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	task, err := s.SyncData.BatchSyncHostIdentifier(hosts.Info, ctx.Kit.Header, ctx.Kit.Rid)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	retry := true
	endTime := time.Now().Add(30 * time.Second).Unix()
	result := &metadata.SyncHostIdentifierResult{
		TaskID: task.TaskID,
	}
	for retry && time.Now().Unix() < endTime {
		retry = false
		resultMap, err := s.SyncData.GetTaskExecutionResultMap(task)
		if err != nil {
			ctx.RespEntityWithError(result, err)
			return
		}

		// 该任务包含的主机没有拿到全部的结果
		if len(task.HostInfos) != len(resultMap) {
			continue
		}

		failIDs := make([]int64, 0)
		successIDs := make([]int64, 0)
		for _, hostInfo := range task.HostInfos {
			key := hostidentifier.HostKey(strconv.FormatInt(hostInfo.CloudID, 10), hostInfo.HostInnerIP)
			code := gjson.Get(resultMap[key], "error_code").Int()
			if code == common.CCSuccess {
				successIDs = append(successIDs, hostInfo.HostID)
				continue
			}

			if code == hostidentifier.Handling {
				retry = true
			}

			failIDs = append(failIDs, hostInfo.HostID)
		}

		result.FailedList = failIDs
		result.SuccessList = successIDs
	}

	if len(result.FailedList) != 0 {
		ctx.RespEntityWithError(result, ctx.Kit.CCError.CCErrorf(common.CCErrEventPushHostIdentifierFailed))
		return
	}

	ctx.RespEntity(result)
}

func (s *Service) haveAuthority(kit *rest.Kit, hostIDs []int64) error {
	businessIDs, err := s.getHostBusinessIDs(kit, hostIDs)
	if err != nil {
		blog.Errorf("get host businessIDs error, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	err = s.AuthManager.AuthorizeByBusinessID(kit.Ctx, kit.Header, meta.ViewBusinessResource, businessIDs...)
	if err != nil {
		blog.Errorf("authorize businesses failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (s *Service) getHostBusinessIDs(kit *rest.Kit, hostIDs []int64) ([]int64, error) {
	cond := &metadata.HostModuleRelationRequest{
		HostIDArr: hostIDs,
		Fields:    []string{common.BKAppIDField},
	}
	result, err := s.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("http do error, input: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}
	if result.Count == 0 {
		blog.Errorf("get host biz count is 0, hostIDs: %v, rid: %s", hostIDs, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommGetBusinessIDByHostIDFailed)
	}

	businessIDs := make([]int64, result.Count)
	for _, info := range result.Info {
		businessIDs = append(businessIDs, info.AppID)
	}

	return businessIDs, nil
}

func (s *Service) getHostInfo(kit *rest.Kit, hostIDs []int64) (*metadata.ListHostResult, error) {
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

	hosts, err := s.engine.CoreAPI.CoreService().Host().ListHosts(kit.Ctx, kit.Header, options)
	if err != nil {
		blog.Errorf("http do error, hostIDs: %v, err: %v, rid: %s", hostIDs, err, kit.Rid)
		return nil, err
	}

	if hosts.Count == 0 {
		blog.Errorf("host count is 0, hostIDs: %v, rid: %s", hostIDs, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrHostGetFail)
	}
	return hosts, nil
}

// PushHostIdentifier push host identifier message to host, returns the gse taskID that can go to gse to query the
// result of the task
func (s *Service) PushHostIdentifier(ctx *rest.Contexts) {
	if s.SyncData == nil {
		blog.Errorf("push host identifier disabled, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrEventSyncHostIdentifierDisabled))
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)

	hostIDArray := new(metadata.HostIDArray)
	if err := ctx.DecodeInto(hostIDArray); err != nil {
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
		if err := s.authByHostIDs(ctx.Kit, hostIDArray.HostIDs); err != nil {
			blog.Errorf("auth by host ids failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthNotHavePermission))
			return
		}
	}

	hosts, err := s.getHostInfo(ctx.Kit, hostIDArray.HostIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	task, err := s.SyncData.BatchSyncHostIdentifier(hosts.Info, ctx.Kit.Header, ctx.Kit.Rid)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	hostInfos := make([]metadata.HostBriefInfo, 0)
	for _, host := range task.HostInfos {
		key := hostidentifier.HostKey(strconv.FormatInt(host.CloudID, 10), host.HostInnerIP)
		info := metadata.HostBriefInfo{
			HostID:         host.HostID,
			Identification: key,
		}
		hostInfos = append(hostInfos, info)
	}

	ctx.RespEntity(&metadata.SyncIdentifierResult{
		TaskID:    task.TaskID,
		HostInfos: hostInfos,
	})
}

func (s *Service) authByHostIDs(kit *rest.Kit, hostIDs []int64) error {
	cond := &metadata.HostModuleRelationRequest{
		HostIDArr: hostIDs,
		Fields:    []string{common.BKAppIDField, common.BKHostIDField},
	}
	result, err := s.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("http do error, input: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	resourcePoolBusinessID, err := s.getResourcePoolBusinessID(kit)
	if err != nil {
		blog.Errorf("get resource pool business id failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	businessIDs := make([]int64, 0)
	resourcePoolHostIDs := make([]int64, 0)
	for _, host := range result.Info {
		if host.AppID == resourcePoolBusinessID {
			resourcePoolHostIDs = append(resourcePoolHostIDs, host.HostID)
			continue
		}

		businessIDs = append(businessIDs, host.AppID)
	}

	if len(businessIDs) != 0 {
		err = s.AuthManager.AuthorizeByBusinessID(kit.Ctx, kit.Header, meta.ViewBusinessResource, businessIDs...)
		if err != nil {
			blog.Errorf("authorize businesses failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}

	if len(resourcePoolHostIDs) != 0 {
		err = s.AuthManager.AuthorizeByHostsIDs(kit.Ctx, kit.Header, meta.Update, resourcePoolHostIDs...)
		if err != nil {
			blog.Errorf("authorize host ids failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}

	return nil
}

func (s *Service) getResourcePoolBusinessID(kit *rest.Kit) (int64, error) {
	query := &metadata.QueryCondition{
		Fields: []string{common.BKAppIDField, common.BkSupplierAccount},
		Condition: map[string]interface{}{
			common.BKDefaultField: common.DefaultAppFlag,
		},
	}

	result, err := s.engine.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp,
		query)
	if err != nil {
		blog.Errorf("get biz by query failed, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	for _, biz := range result.Info {
		if kit.SupplierAccount != util.GetStrByInterface(biz[common.BkSupplierAccount]) {
			continue
		}

		if !biz.Exists(common.BKAppIDField) {
			// this can not be happen normally.
			return 0, errors.New("can not find resource pool business id")
		}

		bizID, err := biz.Int64(common.BKAppIDField)
		if err != nil {
			return 0, fmt.Errorf("get resource pool biz id failed, err: %v", err)
		}
		return bizID, nil
	}

	return 0, errors.New("can not find resource pool business id")
}
