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

package transfer

import (
	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type TransferManager struct {
	dbProxy    dal.RDB
	eventCli   eventclient.Client
	cache      *redis.Client
	dependence OperationDependence
}

type OperationDependence interface {
	AutoCreateServiceInstanceModuleHost(ctx core.ContextParams, hostID int64, moduleID int64) (*metadata.ServiceInstance, errors.CCErrorCoder)
}

func New(db dal.RDB, cache *redis.Client, ec eventclient.Client, dependence OperationDependence) *TransferManager {
	return &TransferManager{
		dbProxy:    db,
		cache:      cache,
		eventCli:   ec,
		dependence: dependence,
	}
}

// NewHostModuleTransfer business normal module transfer
func (manager *TransferManager) NewHostModuleTransfer(ctx core.ContextParams, bizID int64, moduleIDArr []int64, isIncr bool) *genericTransfer {
	return &genericTransfer{
		dbProxy:     manager.dbProxy,
		eventCli:    manager.eventCli,
		dependent:   manager.dependence,
		moduleIDArr: moduleIDArr,
		bizID:       bizID,
		isIncrement: isIncr,
	}
}

// TransferHostToInnerModule transfer host to inner module, default module contain(idle module, fault module)
func (manager *TransferManager) TransferToInnerModule(ctx core.ContextParams, input *metadata.TransferHostToInnerModule) ([]metadata.ExceptionResult, error) {

	transfer := manager.NewHostModuleTransfer(ctx, input.ApplicationID, []int64{input.ModuleID}, false)

	exit, err := transfer.HasInnerModule(ctx)
	if err != nil {
		blog.ErrorJSON("TransferHostToInnerModule failed, HasInnerModule failed, input:%s, err:%s, rid:%s", input, err.Error(), ctx.ReqID)
		return nil, err
	}
	if !exit {
		blog.ErrorJSON("TransferHostToInnerModule validate module failed, module ID is not default module. input:%s, rid:%s", input, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCoreServiceModuleNotDefaultModuleErr, input.ModuleID, input.ApplicationID)
	}
	if err := transfer.DoTransferToInnerCheck(ctx, input.HostID); err != nil {
		blog.ErrorJSON("TransferHostToInnerModule failed, DoTransferToInnerCheck failed, err: %s, rid:%s", err.Error(), ctx.ReqID)
		return nil, err
	}
	err = transfer.ValidParameter(ctx)
	if err != nil {
		blog.ErrorJSON("TransferHostToInnerModule failed, ValidParameter failed, input:%s, err:%s, rid:%s", input, err.Error(), ctx.ReqID)
		return nil, err
	}

	var exceptionArr []metadata.ExceptionResult
	for _, hostID := range input.HostID {
		err := transfer.Transfer(ctx, hostID)
		if err != nil {
			blog.ErrorJSON("TransferHostToInnerModule failed, Transfer module host relation failed, input:%s, hostID:%s, err:%s, rid:%s", input, hostID, err.Error(), ctx.ReqID)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.GetCode()),
				OriginIndex: hostID,
			})
		}
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, ctx.Error.CCError(common.CCErrCoreServiceTransferHostModuleErr)
	}

	return nil, nil
}

// TransferHostModule transfer host to use add module
// 目标模块不能为空闲机模块
func (manager *TransferManager) TransferToNormalModule(ctx core.ContextParams, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error) {
	// 确保目标模块不能为空闲机模块
	defaultModuleFilter := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultFlagDefaultValue,
		},
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: input.ModuleID,
		},
	}
	defaultModuleCount, err := manager.dbProxy.Table(common.BKTableNameBaseModule).Find(defaultModuleFilter).Count(ctx.Context)
	if err != nil {
		blog.ErrorJSON("TransferToNormalModule failed, filter default module failed, filter:%s, err:%s, rid:%s", defaultModuleFilter, common.BKTableNameBaseModule, err.Error(), ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if defaultModuleCount > 0 {
		blog.ErrorJSON("TransferToNormalModule failed, target module shouldn't be default module, input:%s, defaultModuleCount:%s, rid:%s", input, defaultModuleCount, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCoreServiceTransferToDefaultModuleUseWrongMethod)
	}

	var exceptionArr []metadata.ExceptionResult

	// 检查主机从哪个模块移除，并且确认主机可以从该模块移除
	if input.IsIncrement == false {
		hostConfigFilter := map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{
				common.BKDBIN: input.HostID,
			},
		}
		hostModuleConfigs := make([]metadata.ModuleHost, 0)
		if err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(hostConfigFilter).All(ctx.Context, &hostModuleConfigs); err != nil {
			blog.ErrorJSON("TransferToNormalModule failed, find default module failed, filter:%s, hostID:%s, err:%s, rid:%s", defaultModuleFilter, common.BKTableNameBaseModule, err.Error(), ctx.ReqID)
			return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
		}
		hostModuleMap := make(map[int64][]int64)
		for _, hostConfig := range hostModuleConfigs {
			if _, exist := hostModuleMap[hostConfig.HostID]; exist == false {
				hostModuleMap[hostConfig.HostID] = make([]int64, 0)
			}
			hostModuleMap[hostConfig.HostID] = append(hostModuleMap[hostConfig.HostID], hostConfig.ModuleID)
		}

		for hostID, originalModuleIDs := range hostModuleMap {
			removedModuleIDs := make([]int64, 0)
			for _, moduleID := range originalModuleIDs {
				if util.InArray(moduleID, input.ModuleID) == false {
					removedModuleIDs = append(removedModuleIDs, moduleID)
				}
			}
			if len(removedModuleIDs) == 0 {
				continue
			}
			serviceInstanceFilter := map[string]interface{}{
				common.BKHostIDField: hostID,
				common.BKModuleIDField: map[string]interface{}{
					common.BKDBIN: removedModuleIDs,
				},
			}
			instanceCount, err := manager.dbProxy.Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).Count(ctx.Context)
			if err != nil {
				blog.ErrorJSON("TransferToNormalModule failed, find service instance failed, filter:%s, hostID:%s, err:%s, rid:%s", serviceInstanceFilter, common.BKTableNameServiceInstance, err.Error(), ctx.ReqID)
				return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
			}
			if instanceCount > 0 {
				err := ctx.Error.CCError(common.CCErrCoreServiceForbiddenReleaseHostReferencedByServiceInstance)
				exceptionArr = append(exceptionArr, metadata.ExceptionResult{
					Message:     err.Error(),
					Code:        int64(err.GetCode()),
					OriginIndex: hostID,
				})
			}
		}
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, ctx.Error.CCError(common.CCErrCoreServiceForbiddenReleaseHostReferencedByServiceInstance)
	}

	transfer := manager.NewHostModuleTransfer(ctx, input.ApplicationID, input.ModuleID, input.IsIncrement)

	err = transfer.ValidParameter(ctx)
	if err != nil {
		blog.ErrorJSON("TransferToNormalModule failed, ValidParameter failed, input:%s, err:%s, rid:%s", input, err, ctx.ReqID)
		return nil, err
	}
	for _, hostID := range input.HostID {
		err := transfer.Transfer(ctx, hostID)
		if err != nil {
			blog.ErrorJSON("TransferToNormalModule failed, Transfer module host relation failed. input:%s, hostID:%s, err:%s, rid:%s", input, hostID, err, ctx.ReqID)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.GetCode()),
				OriginIndex: hostID,
			})
		}
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, ctx.Error.CCError(common.CCErrCoreServiceTransferHostModuleErr)
	}

	return nil, nil
}

// RemoveHostFromModule 将主机从模块中移出
// 如果主机属于n+1个模块（n>0），操作之后，主机属于n个模块
// 如果主机属于1个模块, 且非空闲机模块，操作之后，主机属于空闲机模块
// 如果主机属于空闲机模块，操作失败
// 如果主机属于故障机模块，操作失败
// 如果主机不在参数指定的模块中，操作失败
func (manager *TransferManager) RemoveFromModule(ctx core.ContextParams, input *metadata.RemoveHostsFromModuleOption) ([]metadata.ExceptionResult, error) {
	hostConfigFilter := map[string]interface{}{
		common.BKHostIDField: input.HostID,
		common.BKAppIDField:  input.ApplicationID,
	}
	hostConfigs := make([]metadata.ModuleHost, 0)
	if err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(hostConfigFilter).All(ctx.Context, &hostConfigs); err != nil {
		blog.ErrorJSON("RemoveFromModule failed, find host module config failed, filter:%s, hostID:%s, err:%s, rid:%s", hostConfigFilter, common.BKTableNameModuleHostConfig, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrHostModuleConfigFailed, err.Error())
	}

	// 如果主机不在参数指定的模块中，操作失败
	if len(hostConfigs) == 0 {
		blog.ErrorJSON("RemoveFromModule failed, host invalid, host module config not found, input:%s, rid:%s", input, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrHostModuleNotExist)
	}

	moduleIDs := make([]int64, 0)
	for _, hostConfig := range hostConfigs {
		moduleIDs = append(moduleIDs, hostConfig.ModuleID)
	}

	// 检查 moduleIDs 是否有空闲机或故障机模块
	// 如果主机属于空闲机模块，操作失败
	// 如果主机属于故障机模块，操作失败
	defaultModuleFilter := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIDs,
		},
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultFlagDefaultValue,
		},
	}
	defaultModuleCount, err := manager.dbProxy.Table(common.BKTableNameBaseModule).Find(defaultModuleFilter).Count(ctx.Context)
	if err != nil {
		blog.ErrorJSON("RemoveFromModule failed, filter default module failed, filter:%s, hostID:%s, err:%s, rid:%s", defaultModuleFilter, common.BKTableNameBaseModule, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrHostGetModuleFail, err.Error())
	}
	if defaultModuleCount > 0 {
		blog.ErrorJSON("RemoveFromModule failed, default module shouldn't in target modules, input:%s, rid:%s", input, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrHostRemoveFromDefaultModuleFailed)
	}

	targetModuleIDs := make([]int64, 0)
	for _, moduleID := range moduleIDs {
		if moduleID != input.ModuleID {
			targetModuleIDs = append(targetModuleIDs, moduleID)
		}
	}
	if len(targetModuleIDs) > 0 {
		option := metadata.HostsModuleRelation{
			ApplicationID: input.ApplicationID,
			HostID:        []int64{input.HostID},
			ModuleID:      targetModuleIDs,
			IsIncrement:   false,
		}
		result, err := manager.TransferToNormalModule(ctx, &option)
		if err != nil {
			blog.ErrorJSON("RemoveFromModule failed, TransferToNormalModule failed, input:%s, option:%s, err:%s, rid:%s", input, option, err.Error(), ctx.ReqID)
			return nil, err
		}
		return result, nil
	}

	// transfer host to idle module
	idleModuleFilter := map[string]interface{}{
		common.BKAppIDField:   input.ApplicationID,
		common.BKDefaultField: common.DefaultResModuleFlag,
	}
	idleModule := metadata.ModuleHost{}
	if err := manager.dbProxy.Table(common.BKTableNameBaseModule).Find(idleModuleFilter).One(ctx.Context, &idleModule); err != nil {
		return nil, ctx.Error.CCErrorf(common.CCErrHostGetModuleFail, err.Error())
	}
	innerModuleOption := metadata.TransferHostToInnerModule{
		ApplicationID: input.ApplicationID,
		ModuleID:      idleModule.ModuleID,
		HostID:        []int64{input.HostID},
	}
	result, err := manager.TransferToInnerModule(ctx, &innerModuleOption)
	if err != nil {
		blog.ErrorJSON("RemoveFromModule failed, TransferToInnerModule failed, filter:%s, option:%s, err:%s, rid:%s", input, innerModuleOption, err.Error(), ctx.ReqID)
		return nil, err
	}
	return result, nil
}

// TransferHostCrossBusiness Host cross-business transfer
func (manager *TransferManager) TransferToAnotherBusiness(ctx core.ContextParams, input *metadata.TransferHostsCrossBusinessRequest) ([]metadata.ExceptionResult, error) {
	// check whether there is service instance bound to hosts
	serviceInstanceFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: input.HostIDArr,
		},
	}
	serviceInstanceCount, err := manager.dbProxy.Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).Count(ctx.Context)
	if err != nil {
		blog.Errorf("TransferToAnotherBusiness failed, query host related service instances failed, filter: %s, err: %s, rid: %s", serviceInstanceFilter, err.Error(), ctx.ReqID)
		return nil, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	if serviceInstanceCount > 0 {
		blog.InfoJSON("RemoveFromModule forbidden, %d service instances bound to hosts, input:%s, rid:%s", serviceInstanceCount, input, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCoreServiceForbiddenReleaseHostReferencedByServiceInstance)
	}

	transfer := manager.NewHostModuleTransfer(ctx, input.DstApplicationID, input.DstModuleIDArr, false)
	transfer.SetCrossBusiness(ctx, input.SrcApplicationID)

	err = transfer.ValidParameter(ctx)
	if err != nil {
		blog.ErrorJSON("TransferToAnotherBusiness failed, ValidParameter failed, err:%s, input:%s, rid:%s", err.Error(), input, ctx.ReqID)
		return nil, err
	}
	var exceptionArr []metadata.ExceptionResult
	for _, hostID := range input.HostIDArr {
		err := transfer.Transfer(ctx, hostID)
		if err != nil {
			blog.ErrorJSON("TransferToAnotherBusiness failed, Transfer module host relation error. err:%s, input:%s, hostID:%s, rid:%s", err.Error(), input, hostID, ctx.ReqID)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.GetCode()),
				OriginIndex: hostID,
			})
		}
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, ctx.Error.CCError(common.CCErrCoreServiceTransferHostModuleErr)
	}

	return nil, nil
}

// GetHostModuleRelation get host module relation
func (manager *TransferManager) GetHostModuleRelation(ctx core.ContextParams, input *metadata.HostModuleRelationRequest) (*metadata.HostConfigData, error) {
	if input.Empty() {
		blog.Errorf("GetHostModuleRelation input empty. input:%#v, rid:%s", input, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}
	moduleHostCond := condition.CreateCondition()
	if input.ApplicationID > 0 {
		moduleHostCond.Field(common.BKAppIDField).Eq(input.ApplicationID)
	}
	if len(input.HostIDArr) > 0 {
		moduleHostCond.Field(common.BKHostIDField).In(input.HostIDArr)
	}
	if len(input.ModuleIDArr) > 0 {
		moduleHostCond.Field(common.BKModuleIDField).In(input.ModuleIDArr)
	}
	if len(input.SetIDArr) > 0 {
		moduleHostCond.Field(common.BKSetIDField).In(input.SetIDArr)
	}
	cond := moduleHostCond.ToMapStr()
	if len(cond) == 0 {
		return nil, nil
	}
	cond = util.SetQueryOwner(moduleHostCond.ToMapStr(), ctx.SupplierAccount)

	cnt, err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(cond).Count(ctx.Context)
	if err != nil {
		blog.Errorf("get module host config count failed, err: %v, cond:%#v, rid: %s", err, cond, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	hostModuleArr := make([]metadata.ModuleHost, 0)
	db := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).
		Find(cond).
		Start(uint64(input.Page.Start)).
		Sort(input.Page.Sort)

	if input.Page.Limit > 0 {
		db = db.Limit(uint64(input.Page.Limit))
	}

	err = db.All(ctx.Context, &hostModuleArr)
	if err != nil {
		blog.Errorf("get module host config failed, err: %v, cond:%#v, rid: %s", err, cond, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return &metadata.HostConfigData{
		Count: int64(cnt),
		Info:  hostModuleArr,
		Page:  input.Page,
	}, nil
}

// DeleteHost delete host module relation and host info
func (manager *TransferManager) DeleteFromSystem(ctx core.ContextParams, input *metadata.DeleteHostRequest) ([]metadata.ExceptionResult, error) {

	transfer := manager.NewHostModuleTransfer(ctx, input.ApplicationID, nil, false)
	transfer.SetDeleteHost(ctx)

	err := transfer.ValidParameter(ctx)
	if err != nil {
		blog.ErrorJSON("DeleteFromSystem failed, ValidParameter failed, err:%s, input:%s, rid:%s", err.Error(), input, ctx.ReqID)
		return nil, err
	}

	var exceptionArr []metadata.ExceptionResult
	for _, hostID := range input.HostIDArr {
		err := transfer.Transfer(ctx, hostID)
		if err != nil {
			blog.ErrorJSON("DeleteFromSystem failed, Transfer module host relation failed. err:%s, input:%s, hostID:%s, rid:%s", err.Error(), input, hostID, ctx.ReqID)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.GetCode()),
				OriginIndex: hostID,
			})
		}
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, ctx.Error.CCError(common.CCErrCoreServiceTransferHostModuleErr)
	}

	return nil, nil
}

func (manager *TransferManager) getHostIDModuleMapByHostID(ctx core.ContextParams, appID int64, hostIDArr []int64) (map[int64][]metadata.ModuleHost, errors.CCErrorCoder) {
	moduleHostCond := condition.CreateCondition()
	moduleHostCond.Field(common.BKAppIDField).Eq(appID)
	moduleHostCond.Field(common.BKHostIDField).In(hostIDArr)
	cond := util.SetQueryOwner(moduleHostCond.ToMapStr(), ctx.SupplierAccount)

	var dataArr []metadata.ModuleHost
	err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(cond).All(ctx, &dataArr)
	if err != nil {
		blog.ErrorJSON("getHostIDModuleMapByHostID query db error. err:%s, cond:%s,rid:%s", err.Error(), cond, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	result := make(map[int64][]metadata.ModuleHost, 0)
	for _, item := range dataArr {
		result[item.HostID] = append(result[item.HostID], item)
	}
	return result, nil
}
