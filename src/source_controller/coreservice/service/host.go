/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core/instances"
	"configcenter/src/storage/driver/mongodb"
)

// TransferHostToInnerModule TODO
func (s *coreService) TransferHostToInnerModule(ctx *rest.Contexts) {
	inputData := &metadata.TransferHostToInnerModule{}
	if err := ctx.DecodeInto(inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.HostOperation().TransferToInnerModule(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("transfer host to default module failed. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// TransferHostToNormalModule TODO
func (s *coreService) TransferHostToNormalModule(ctx *rest.Contexts) {
	inputData := &metadata.HostsModuleRelation{}
	if err := ctx.DecodeInto(inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.HostOperation().TransferToNormalModule(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("transfer host to normal module failed. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// TransferHostToAnotherBusiness TODO
func (s *coreService) TransferHostToAnotherBusiness(ctx *rest.Contexts) {
	inputData := &metadata.TransferHostsCrossBusinessRequest{}
	if err := ctx.DecodeInto(inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().TransferToAnotherBusiness(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("transfer host across business failed. err: %s, input: %s, rid: %s", err.Error(), inputData,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// RemoveFromModule TODO
func (s *coreService) RemoveFromModule(ctx *rest.Contexts) {
	inputData := &metadata.RemoveHostsFromModuleOption{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.HostOperation().RemoveFromModule(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("remove host from module failed. err: %s, input: %s, rid: %s", err.Error(), inputData,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// GetHostModuleRelation TODO
func (s *coreService) GetHostModuleRelation(ctx *rest.Contexts) {
	inputData := &metadata.HostModuleRelationRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	relationArr, err := s.core.HostOperation().GetHostModuleRelation(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("get host module relation error, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(relationArr)
}

// DeleteHostFromSystem TODO
func (s *coreService) DeleteHostFromSystem(ctx *rest.Contexts) {
	inputData := &metadata.DeleteHostRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().DeleteFromSystem(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("delete host error. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// HostIdentifier TODO
func (s *coreService) HostIdentifier(ctx *rest.Contexts) {
	inputData := &metadata.SearchHostIdentifierParam{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	hostIdentifierArr, err := s.core.HostOperation().Identifier(ctx.Kit, inputData)

	if err != nil {
		blog.Infof("Identifier host identifier handle error, err: %v, input: %+v, rid:%s", err, inputData, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.SearchHostIdentifierData{Info: hostIdentifierArr, Count: len(hostIdentifierArr)})
}

// GetHostByID TODO
func (s *coreService) GetHostByID(ctx *rest.Contexts) {
	hostID, err := strconv.Atoi(ctx.Request.PathParameter("bk_host_id"))
	if err != nil {
		blog.Errorf("GetHostByID failed, get host by id, but got invalid host id, hostID: %s, err: %+v, rid: %s",
			hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField))
		return
	}

	result := make(metadata.HostMapStr, 0)
	condition := common.KvMap{common.BKHostIDField: hostID}
	err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(condition).One(ctx.Kit.Ctx,
		&result)
	// TODO: return error for not found and deal error with all callers
	if err != nil && !mongodb.IsNotFoundError(err) {
		blog.Errorf("GetHostByID failed, get host by id[%d] failed, err: %+v, rid: %s", hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	if result[common.BKCreatedAt] == nil && result[common.CreateTimeField] != nil {
		result[common.BKCreatedAt] = result[common.CreateTimeField]
	}

	if result[common.BKUpdatedAt] == nil && result[common.LastTimeField] != nil {
		result[common.BKUpdatedAt] = result[common.LastTimeField]
	}

	ctx.RespEntity(result)
}

// GetHosts TODO
func (s *coreService) GetHosts(ctx *rest.Contexts) {
	var dat metadata.QueryInput
	if err := ctx.DecodeInto(&dat); err != nil {
		ctx.RespAutoError(err)
		return
	}

	condition := dat.Condition
	if dat.TimeCondition != nil {
		var err error
		condition, err = dat.TimeCondition.MergeTimeCondition(condition)
		if err != nil {
			blog.Errorf("merge time condition failed, error: %v, input: %+v, rid: %s", err, dat, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	fieldArr := util.SplitStrField(dat.Fields, ",")

	query := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(condition).Sort(dat.Sort).
		Start(uint64(dat.Start)).Limit(uint64(dat.Limit))
	info, err := instances.FindInst(ctx.Kit, fieldArr, query, common.BKInnerObjIDHost)
	if err != nil {
		blog.Errorf("failed to query the host, cond: %v, err: %v, rid: %s", condition, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var finalCount uint64

	if !dat.DisableCounter {
		count, err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(condition).
			Count(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("get object failed type:%s ,input: %v error: %v, rid: %s", common.BKInnerObjIDHost, dat, err,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostSelectInst))
			return
		}
		finalCount = count
	}

	ctx.RespEntity(metadata.HostInfo{
		Count: int(finalCount),
		Info:  info,
	})
}

// GetDistinctHostIDsByTopoRelation get all  host ids by topology relation condition
func (s *coreService) GetDistinctHostIDsByTopoRelation(ctx *rest.Contexts) {
	inputData := &metadata.DistinctHostIDByTopoRelationRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	hostIDArr, err := s.core.HostOperation().GetDistinctHostIDsByTopoRelation(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("GetDistinctHostIDsByTopoRelation  error. err:%s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.DistinctID{IDArr: hostIDArr})
}

// TransferHostResourceDirectory TODO
func (s *coreService) TransferHostResourceDirectory(ctx *rest.Contexts) {
	input := &metadata.TransferHostResourceDirectory{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().TransferResourceDirectory(ctx.Kit, input)
	if err != nil {
		blog.ErrorJSON("TransferHostResourceDirectory  error. err:%s, input:%s, rid:%s", err.Error(), input,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func delRedundantDefaultAreaHostByConds(kit *rest.Kit, defaultAreaHostCond mapstr.MapStr) error {

	hosts := make([]metadata.DefaultAreaHost, 0)
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameDefaultAreaHost).Find(defaultAreaHostCond).All(
		kit.Ctx, &hosts)
	if err != nil {
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	tenantHostIDMap := make(map[string][]int64)
	for _, item := range hosts {
		if _, ok := tenantHostIDMap[item.TenantID]; !ok {
			tenantHostIDMap[item.TenantID] = make([]int64, 0)
		}
		tenantHostIDMap[item.TenantID] = append(tenantHostIDMap[item.TenantID], item.HostID)
	}

	if err := deleteRedundantHosts(kit, tenantHostIDMap); err != nil {
		blog.Errorf("delete redundant default area host failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// DelRedundantDefaultAreaHost delete redundant default area host
func (s *coreService) DelRedundantDefaultAreaHost(ctx *rest.Contexts) {
	option := new(metadata.DelRedundantDefaultAreaHostsOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	defaultAreaCond := make(mapstr.MapStr)
	var err error
	switch option.OpType {
	case metadata.OperationByHostID:
		defaultAreaCond = getDefaultAreaHostCondByID(option.Hosts)

	case metadata.OperationByIP:
		defaultAreaCond = getDefaultAreaHostCondByIP(option.Hosts)

	default:
		blog.Errorf("operation type is invalid, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "operation_type"))
	}

	err = delRedundantDefaultAreaHostByConds(ctx.Kit, defaultAreaCond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// deleteRedundantHosts delete redundant host by condition
func deleteRedundantHosts(kit *rest.Kit, tenantHostMap map[string][]int64) error {
	removeHostIDs := make([]int64, 0)
	for tenantID, hosts := range tenantHostMap {
		findCond := mapstr.MapStr{
			common.BKHostIDField: map[string]interface{}{condition.BKDBIN: hosts},
		}
		newKit := kit.NewKit().WithTenant(tenantID)
		existHostCount, err := mongodb.Shard(newKit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(findCond).Count(
			newKit.Ctx)
		if err != nil {
			blog.Errorf("get exist host count failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		if existHostCount > 0 {
			blog.Errorf("exist host is conflict with insert host for default area host, rid: %s", kit.Rid)
			return kit.CCError.CCError(common.CCErrDefaultAreaHostIPExist)
		}

		removeHostIDs = append(removeHostIDs, hosts...)
	}

	if len(removeHostIDs) == 0 {
		return kit.CCError.CCError(common.CCErrDefaultAreaHostIPExist)
	}

	cond := mapstr.MapStr{
		common.BKHostIDField: mapstr.MapStr{
			condition.BKDBIN: removeHostIDs,
		},
	}
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameDefaultAreaHost).Delete(kit.Ctx, cond)
	if err != nil {
		blog.Errorf("delete default area host failed, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

func getDefaultAreaHostCondByID(option []metadata.DefaultAreaHost) mapstr.MapStr {
	hostIDs := make([]int64, 0)
	for _, item := range option {
		hostIDs = append(hostIDs, item.HostID)
	}

	return mapstr.MapStr{common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs}}
}

func getDefaultAreaHostCondByIP(hosts []metadata.DefaultAreaHost) mapstr.MapStr {
	ipArr := make([]string, 0)
	ipV6Arr := make([]string, 0)
	for _, host := range hosts {
		if host.InnerIP != nil {
			ipArr = append(ipArr, host.InnerIP...)
		}

		if host.InnerIPv6 != nil {
			ipV6Arr = append(ipV6Arr, host.InnerIPv6...)
		}
	}

	hostConds := make([]mapstr.MapStr, 0)
	if len(ipArr) > 0 {
		hostConds = append(hostConds, mapstr.MapStr{
			common.BKHostInnerIPField: mapstr.MapStr{
				common.BKDBIN: ipArr,
			},
		})
	}

	if len(ipV6Arr) > 0 {
		hostConds = append(hostConds, mapstr.MapStr{
			common.BKHostInnerIPv6Field: mapstr.MapStr{
				common.BKDBIN: ipV6Arr,
			},
		})
	}

	return mapstr.MapStr{condition.BKDBOR: hostConds}
}
