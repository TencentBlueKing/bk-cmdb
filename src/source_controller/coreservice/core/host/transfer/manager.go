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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
	"gopkg.in/redis.v5"
)

type TransferManager struct {
	dbProxy             dal.RDB
	eventCli            eventclient.Client
	cache               *redis.Client
	dependence          OperationDependence
	hostApplyDependence HostApplyRuleDependence
}

type OperationDependence interface {
	AutoCreateServiceInstanceModuleHost(kit *rest.Kit, hostID int64, moduleID int64) (*metadata.ServiceInstance, errors.CCErrorCoder)
	SelectObjectAttWithParams(kit *rest.Kit, objID string, bizID int64) (attribute []metadata.Attribute, err error)
	UpdateModelInstance(kit *rest.Kit, objID string, param metadata.UpdateOption) (*metadata.UpdatedCount, error)
}

type HostApplyRuleDependence interface {
	RunHostApplyOnHosts(kit *rest.Kit, bizID int64, option metadata.UpdateHostByHostApplyRuleOption) (metadata.MultipleHostApplyResult, errors.CCErrorCoder)
}

func New(db dal.RDB, cache *redis.Client, ec eventclient.Client, dependence OperationDependence, hostApplyDependence HostApplyRuleDependence) *TransferManager {
	return &TransferManager{
		dbProxy:             db,
		cache:               cache,
		eventCli:            ec,
		dependence:          dependence,
		hostApplyDependence: hostApplyDependence,
	}
}

// NewHostModuleTransfer business normal module transfer
func (manager *TransferManager) NewHostModuleTransfer(kit *rest.Kit, bizID int64, moduleIDArr []int64, isIncr bool) *genericTransfer {
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
func (manager *TransferManager) TransferToInnerModule(kit *rest.Kit, input *metadata.TransferHostToInnerModule) ([]metadata.ExceptionResult, error) {

	transfer := manager.NewHostModuleTransfer(kit, input.ApplicationID, []int64{input.ModuleID}, false)

	exit, err := transfer.HasInnerModule(kit)
	if err != nil {
		blog.ErrorJSON("TransferHostToInnerModule failed, HasInnerModule failed, input:%s, err:%s, rid:%s", input, err.Error(), kit.Rid)
		return nil, err
	}
	if !exit {
		blog.ErrorJSON("TransferHostToInnerModule validate module failed, module ID is not default module. input:%s, rid:%s", input, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCoreServiceModuleNotDefaultModuleErr, input.ModuleID, input.ApplicationID)
	}
	if err := transfer.DoTransferToInnerCheck(kit, input.HostID); err != nil {
		blog.ErrorJSON("TransferHostToInnerModule failed, DoTransferToInnerCheck failed, err: %s, rid:%s", err.Error(), kit.Rid)
		return nil, err
	}
	err = transfer.ValidParameter(kit)
	if err != nil {
		blog.ErrorJSON("TransferHostToInnerModule failed, ValidParameter failed, input:%s, err:%s, rid:%s", input, err.Error(), kit.Rid)
		return nil, err
	}

	var exceptionArr []metadata.ExceptionResult
	for _, hostID := range input.HostID {
		err := transfer.Transfer(kit, hostID)
		if err != nil {
			blog.ErrorJSON("TransferHostToInnerModule failed, Transfer module host relation failed, input:%s, hostID:%s, err:%s, rid:%s", input, hostID, err.Error(), kit.Rid)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.GetCode()),
				OriginIndex: hostID,
			})
		}
	}
	updateHostOption := metadata.UpdateHostByHostApplyRuleOption{
		HostIDs: input.HostID,
	}
	if _, err := manager.hostApplyDependence.RunHostApplyOnHosts(kit, input.ApplicationID, updateHostOption); err != nil {
		blog.Warnf("TransferHostToInnerModule success, but RunHostApplyOnHosts failed, bizID: %d, option: %+v, err: %+v, rid: %s", input.ApplicationID, updateHostOption, err, kit.Rid)
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, kit.CCError.CCError(common.CCErrCoreServiceTransferHostModuleErr)
	}

	return nil, nil
}

// TransferHostModule transfer host to use add module
// 目标模块不能为空闲机模块
func (manager *TransferManager) TransferToNormalModule(kit *rest.Kit, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error) {
	// 确保目标模块不能为空闲机模块
	defaultModuleFilter := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultFlagDefaultValue,
		},
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: input.ModuleID,
		},
	}
	defaultModuleCount, err := manager.dbProxy.Table(common.BKTableNameBaseModule).Find(defaultModuleFilter).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("TransferToNormalModule failed, filter default module failed, filter:%s, err:%s, rid:%s", defaultModuleFilter, common.BKTableNameBaseModule, err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if defaultModuleCount > 0 {
		blog.ErrorJSON("TransferToNormalModule failed, target module shouldn't be default module, input:%s, defaultModuleCount:%s, rid:%s", input, defaultModuleCount, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCoreServiceTransferToDefaultModuleUseWrongMethod)
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
		if err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(hostConfigFilter).All(kit.Ctx, &hostModuleConfigs); err != nil {
			blog.ErrorJSON("TransferToNormalModule failed, find default module failed, filter:%s, hostID:%s, err:%s, rid:%s", defaultModuleFilter, common.BKTableNameBaseModule, err.Error(), kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
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
			instanceCount, err := manager.dbProxy.Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).Count(kit.Ctx)
			if err != nil {
				blog.ErrorJSON("TransferToNormalModule failed, find service instance failed, filter:%s, hostID:%s, err:%s, rid:%s", serviceInstanceFilter, common.BKTableNameServiceInstance, err.Error(), kit.Rid)
				return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			}
			if instanceCount > 0 {
				err := kit.CCError.CCError(common.CCErrCoreServiceForbiddenReleaseHostReferencedByServiceInstance)
				exceptionArr = append(exceptionArr, metadata.ExceptionResult{
					Message:     err.Error(),
					Code:        int64(err.GetCode()),
					OriginIndex: hostID,
				})
			}
		}
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, kit.CCError.CCError(common.CCErrCoreServiceForbiddenReleaseHostReferencedByServiceInstance)
	}

	transfer := manager.NewHostModuleTransfer(kit, input.ApplicationID, input.ModuleID, input.IsIncrement)

	err = transfer.ValidParameter(kit)
	if err != nil {
		blog.ErrorJSON("TransferToNormalModule failed, ValidParameter failed, input:%s, err:%s, rid:%s", input, err, kit.Rid)
		return nil, err
	}
	for _, hostID := range input.HostID {
		err := transfer.Transfer(kit, hostID)
		if err != nil {
			blog.ErrorJSON("TransferToNormalModule failed, Transfer module host relation failed. input:%s, hostID:%s, err:%s, rid:%s", input, hostID, err, kit.Rid)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.GetCode()),
				OriginIndex: hostID,
			})
		}
	}
	updateHostOption := metadata.UpdateHostByHostApplyRuleOption{
		HostIDs: input.HostID,
	}
	if _, err := manager.hostApplyDependence.RunHostApplyOnHosts(kit, input.ApplicationID, updateHostOption); err != nil {
		blog.Warnf("TransferToNormalModule success, but RunHostApplyOnHosts failed, bizID: %d, option: %+v, err: %+v, rid: %s", input.ApplicationID, updateHostOption, err, kit.Rid)
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, kit.CCError.CCError(common.CCErrCoreServiceTransferHostModuleErr)
	}

	return nil, nil
}

// RemoveHostFromModule 将主机从模块中移出
// 如果主机属于n+1个模块（n>0），操作之后，主机属于n个模块
// 如果主机属于1个模块, 且非空闲机模块，操作之后，主机属于空闲机模块
// 如果主机属于空闲机模块，操作失败
// 如果主机属于故障机模块，操作失败
// 如果主机不在参数指定的模块中，操作失败
func (manager *TransferManager) RemoveFromModule(kit *rest.Kit, input *metadata.RemoveHostsFromModuleOption) ([]metadata.ExceptionResult, error) {
	hostConfigFilter := map[string]interface{}{
		common.BKHostIDField: input.HostID,
		common.BKAppIDField:  input.ApplicationID,
	}
	hostConfigs := make([]metadata.ModuleHost, 0)
	if err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(hostConfigFilter).All(kit.Ctx, &hostConfigs); err != nil {
		blog.ErrorJSON("RemoveFromModule failed, find host module config failed, filter:%s, hostID:%s, err:%s, rid:%s", hostConfigFilter, common.BKTableNameModuleHostConfig, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrHostModuleConfigFailed, err.Error())
	}

	// 如果主机不在参数指定的模块中，操作失败
	if len(hostConfigs) == 0 {
		blog.ErrorJSON("RemoveFromModule failed, host invalid, host module config not found, input:%s, rid:%s", input, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrHostModuleNotExist)
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
	defaultModuleCount, err := manager.dbProxy.Table(common.BKTableNameBaseModule).Find(defaultModuleFilter).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("RemoveFromModule failed, filter default module failed, filter:%s, hostID:%s, err:%s, rid:%s", defaultModuleFilter, common.BKTableNameBaseModule, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrHostGetModuleFail, err.Error())
	}
	if defaultModuleCount > 0 {
		blog.ErrorJSON("RemoveFromModule failed, default module shouldn't in target modules, input:%s, rid:%s", input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrHostRemoveFromDefaultModuleFailed)
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
		result, err := manager.TransferToNormalModule(kit, &option)
		if err != nil {
			blog.ErrorJSON("RemoveFromModule failed, TransferToNormalModule failed, input:%s, option:%s, err:%s, rid:%s", input, option, err.Error(), kit.Rid)
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
	if err := manager.dbProxy.Table(common.BKTableNameBaseModule).Find(idleModuleFilter).One(kit.Ctx, &idleModule); err != nil {
		return nil, kit.CCError.CCErrorf(common.CCErrHostGetModuleFail, err.Error())
	}
	innerModuleOption := metadata.TransferHostToInnerModule{
		ApplicationID: input.ApplicationID,
		ModuleID:      idleModule.ModuleID,
		HostID:        []int64{input.HostID},
	}
	result, err := manager.TransferToInnerModule(kit, &innerModuleOption)
	if err != nil {
		blog.ErrorJSON("RemoveFromModule failed, TransferToInnerModule failed, filter:%s, option:%s, err:%s, rid:%s", input, innerModuleOption, err.Error(), kit.Rid)
		return nil, err
	}
	return result, nil
}

// TransferHostCrossBusiness Host cross-business transfer
func (manager *TransferManager) TransferToAnotherBusiness(kit *rest.Kit, input *metadata.TransferHostsCrossBusinessRequest) ([]metadata.ExceptionResult, error) {
	// check whether there is service instance bound to hosts
	serviceInstanceFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: input.HostIDArr,
		},
	}
	serviceInstanceCount, err := manager.dbProxy.Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("TransferToAnotherBusiness failed, query host related service instances failed, filter: %s, err: %s, rid: %s", serviceInstanceFilter, err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	if serviceInstanceCount > 0 {
		blog.InfoJSON("RemoveFromModule forbidden, %d service instances bound to hosts, input:%s, rid:%s", serviceInstanceCount, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCoreServiceForbiddenReleaseHostReferencedByServiceInstance)
	}

	transfer := manager.NewHostModuleTransfer(kit, input.DstApplicationID, input.DstModuleIDArr, false)
	transfer.SetCrossBusiness(kit, input.SrcApplicationID)

	err = transfer.ValidParameter(kit)
	if err != nil {
		blog.ErrorJSON("TransferToAnotherBusiness failed, ValidParameter failed, err:%s, input:%s, rid:%s", err.Error(), input, kit.Rid)
		return nil, err
	}

	// attributes in legacy business
	legacyAttributes, err := transfer.dependent.SelectObjectAttWithParams(kit, common.BKInnerObjIDHost, input.SrcApplicationID)
	if err != nil {
		blog.ErrorJSON("TransferToAnotherBusiness failed, SelectObjectAttWithParams failed, bizID: %s, err:%s, rid:%s", input.SrcApplicationID, err.Error(), kit.Rid)
		return nil, err
	}

	// attributes in new business
	newAttributes, err := transfer.dependent.SelectObjectAttWithParams(kit, common.BKInnerObjIDHost, input.DstApplicationID)
	if err != nil {
		blog.ErrorJSON("TransferToAnotherBusiness failed, SelectObjectAttWithParams failed, bizID: %s, err:%s, rid:%s", input.DstApplicationID, err.Error(), kit.Rid)
		return nil, err
	}

	var exceptionArr []metadata.ExceptionResult
	successHostIDs := make([]int64, 0)
	for _, hostID := range input.HostIDArr {
		err := transfer.Transfer(kit, hostID)
		if err != nil {
			blog.ErrorJSON("TransferToAnotherBusiness failed, Transfer module host relation error. err:%s, input:%s, hostID:%s, rid:%s", err.Error(), input, hostID, kit.Rid)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.GetCode()),
				OriginIndex: hostID,
			})
			continue
		}
		successHostIDs = append(successHostIDs, hostID)
	}

	if len(successHostIDs) > 0 {
		// reset private field in legacy business
		if err := manager.clearLegacyPrivateField(kit, legacyAttributes, successHostIDs...); err != nil {
			blog.ErrorJSON("TransferToAnotherBusiness failed, clearLegacyPrivateField failed, hostID:%s, attributes:%s, err:%s, rid:%s", successHostIDs, legacyAttributes, err.Error(), kit.Rid)
			// we should go on setting default value for new private field
		}

		// set default value for private field in new business
		if err := manager.setDefaultPrivateField(kit, newAttributes, successHostIDs...); err != nil {
			blog.ErrorJSON("TransferToAnotherBusiness failed, setDefaultPrivateField failed, hostID:%s, attributes:%s, err:%s, rid:%s", successHostIDs, newAttributes, err.Error(), kit.Rid)
			for _, hostID := range successHostIDs {
				exceptionArr = append(exceptionArr, metadata.ExceptionResult{
					Message:     err.Error(),
					Code:        int64(err.GetCode()),
					OriginIndex: hostID,
				})
			}
		}
	}

	updateHostOption := metadata.UpdateHostByHostApplyRuleOption{
		HostIDs: input.HostIDArr,
	}
	if hostApplyResult, err := manager.hostApplyDependence.RunHostApplyOnHosts(kit, input.DstApplicationID, updateHostOption); err != nil {
		blog.Warnf("TransferToAnotherBusiness success, but RunHostApplyOnHosts failed, bizID: %d, option: %+v, hostApplyResult: %+v, err: %+v, rid: %s", input.DstApplicationID, updateHostOption, hostApplyResult, err, kit.Rid)
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, kit.CCError.CCError(common.CCErrCoreServiceTransferHostModuleErr)
	}

	return nil, nil
}

func (manager *TransferManager) clearLegacyPrivateField(kit *rest.Kit, attributes []metadata.Attribute, hostIDs ...int64) errors.CCErrorCoder {
	doc := make(map[string]interface{}, 0)
	for _, attribute := range attributes {
		bizID, err := attribute.Metadata.ParseBizID()
		if err != nil {
			blog.Warnf("clearLegacyPrivateField, parse bizID from attribute failed, attribute: %+v, err: %s, rid: %s", attribute, err.Error(), kit.Rid)
			continue
		}
		if bizID == 0 {
			continue
		}
		doc[attribute.PropertyID] = nil
	}
	if len(doc) == 0 {
		return nil
	}
	reset := types.ModeUpdate{
		Op:  "unset",
		Doc: doc,
	}
	filter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	if err := manager.dbProxy.Table(common.BKTableNameBaseHost).UpdateMultiModel(kit.Ctx, filter, reset); err != nil {
		blog.ErrorJSON("clearLegacyPrivateField failed. table: %s, filter: %s, doc: %s, err: %s, rid:%s", common.BKTableNameBaseHost, filter, doc, err.Error(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func (manager *TransferManager) setDefaultPrivateField(kit *rest.Kit, attributes []metadata.Attribute, hostID ...int64) errors.CCErrorCoder {
	doc := make(map[string]interface{})
	for _, attribute := range attributes {
		bizID, err := attribute.Metadata.ParseBizID()
		if err != nil {
			blog.Warnf("clearLegacyPrivateField, parse bizID from attribute failed, attribute: %+v, err: %s, rid: %s", attribute, err.Error(), kit.Rid)
			continue
		}
		if bizID == 0 {
			continue
		}
		doc[attribute.PropertyID] = nil
	}
	if len(doc) == 0 {
		return nil
	}
	updateOption := metadata.UpdateOption{
		Data: doc,
		Condition: map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{
				common.BKDBIN: hostID,
			},
		},
	}
	_, err := manager.dependence.UpdateModelInstance(kit, common.BKInnerObjIDHost, updateOption)
	if err != nil {
		blog.ErrorJSON("setDefaultPrivateField failed. UpdateModelInstance failed, option: %s, err: %s, rid:%s", common.BKTableNameBaseHost, updateOption, err.Error(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

// GetHostModuleRelation get host module relation
func (manager *TransferManager) GetHostModuleRelation(kit *rest.Kit, input *metadata.HostModuleRelationRequest) (*metadata.HostConfigData, error) {
	if input.Empty() {
		blog.Errorf("GetHostModuleRelation input empty. input:%#v, rid:%s", input, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
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
	cond = util.SetQueryOwner(moduleHostCond.ToMapStr(), kit.SupplierAccount)

	cnt, err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get module host config count failed, err: %v, cond:%#v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	hostModuleArr := make([]metadata.ModuleHost, 0)
	db := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).
		Find(cond).
		Start(uint64(input.Page.Start)).
		Sort(input.Page.Sort)

	if input.Page.Limit > 0 {
		db = db.Limit(uint64(input.Page.Limit))
	}

	err = db.All(kit.Ctx, &hostModuleArr)
	if err != nil {
		blog.Errorf("get module host config failed, err: %v, cond:%#v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &metadata.HostConfigData{
		Count: int64(cnt),
		Info:  hostModuleArr,
		Page:  input.Page,
	}, nil
}

// DeleteHost delete host module relation and host info
func (manager *TransferManager) DeleteFromSystem(kit *rest.Kit, input *metadata.DeleteHostRequest) ([]metadata.ExceptionResult, error) {

	transfer := manager.NewHostModuleTransfer(kit, input.ApplicationID, nil, false)
	transfer.SetDeleteHost(kit)

	err := transfer.ValidParameter(kit)
	if err != nil {
		blog.ErrorJSON("DeleteFromSystem failed, ValidParameter failed, err:%s, input:%s, rid:%s", err.Error(), input, kit.Rid)
		return nil, err
	}

	var exceptionArr []metadata.ExceptionResult
	for _, hostID := range input.HostIDArr {
		err := transfer.Transfer(kit, hostID)
		if err != nil {
			blog.ErrorJSON("DeleteFromSystem failed, Transfer module host relation failed. err:%s, input:%s, hostID:%s, rid:%s", err.Error(), input, hostID, kit.Rid)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.GetCode()),
				OriginIndex: hostID,
			})
		}
	}
	if len(exceptionArr) > 0 {
		return exceptionArr, kit.CCError.CCError(common.CCErrCoreServiceTransferHostModuleErr)
	}

	return nil, nil
}

func (manager *TransferManager) getHostIDModuleMapByHostID(kit *rest.Kit, appID int64, hostIDArr []int64) (map[int64][]metadata.ModuleHost, errors.CCErrorCoder) {
	moduleHostCond := condition.CreateCondition()
	moduleHostCond.Field(common.BKAppIDField).Eq(appID)
	moduleHostCond.Field(common.BKHostIDField).In(hostIDArr)
	cond := util.SetQueryOwner(moduleHostCond.ToMapStr(), kit.SupplierAccount)

	var dataArr []metadata.ModuleHost
	err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(cond).All(kit.Ctx, &dataArr)
	if err != nil {
		blog.ErrorJSON("getHostIDModuleMapByHostID query db error. err:%s, cond:%s,rid:%s", err.Error(), cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	result := make(map[int64][]metadata.ModuleHost, 0)
	for _, item := range dataArr {
		result[item.HostID] = append(result[item.HostID], item)
	}
	return result, nil
}

func (manager *TransferManager) TransferResourceDirectory(kit *rest.Kit, input *metadata.TransferHostResourceDirectory) errors.CCErrorCoder {
	// validate input bk_module_id
	existHostIDs, err := manager.validTransferResourceDirParams(kit, input)
	if err != nil {
		blog.ErrorJSON("TransferResourceDirectory fail with validTransferResourceDirParams failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	wrongHostIDs := make([]int64, 0)
	if len(existHostIDs) < len(input.HostID) {
		for _, id := range input.HostID {
			if !util.ContainsInt64(existHostIDs, id) {
				wrongHostIDs = append(wrongHostIDs, id)
			}
		}
	}

	cond := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: existHostIDs,
		},
	}
	data := map[string]interface{}{
		common.BKModuleIDField: input.ModuleID,
	}
	updateErr := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Update(kit.Ctx, cond, data)
	if updateErr != nil {
		blog.ErrorJSON("TransferResourceDirectory fail with update table error: %v, cond: %v, data: %v, rid: %v", updateErr, cond, data, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if len(wrongHostIDs) > 0 {
		return kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotUnderAnyResourceDirectory, wrongHostIDs)
	}

	return nil
}

func (manager *TransferManager) validTransferResourceDirParams(kit *rest.Kit, input *metadata.TransferHostResourceDirectory) ([]int64, errors.CCErrorCoder) {
	// valid bk_module_id,资源池目录default=4,空闲机default=1
	cond := mapstr.MapStr{}
	cond[common.BKModuleIDField] = input.ModuleID
	cond[common.BKDBOR] = []mapstr.MapStr{{common.BKDefaultField: 1}, {common.BKDefaultField: 4}}
	count, err := manager.dbProxy.Table(common.BKTableNameBaseModule).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("validTransferResourceDirParams find data error, err: %s, cond:%s, rid: %s", err.Error(), cond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count <= 0 {
		blog.ErrorJSON("validTransferResourceDirParams bk_module_id bind resource directory not exist")
		return nil, kit.CCError.CCError(common.CCErrCoreServiceResourceDirectoryNotExistErr)
	}

	// 确保主机在资源池目录下(default=1的业务)
	bizInfo := &metadata.BizInst{}
	filter := mapstr.MapStr{common.BKDefaultField: 1}
	err = manager.dbProxy.Table(common.BKTableNameBaseApp).Find(filter).One(kit.Ctx, bizInfo)
	if err != nil {
		blog.ErrorJSON("validTransferResourceDirParams find data error, err: %s, cond:%s, rid: %s", err.Error(), cond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	opt := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: input.HostID}, common.BKAppIDField: bizInfo.BizID}
	moduleHost := make([]metadata.ModuleHost, 0)
	if err := manager.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(opt).All(kit.Ctx, &moduleHost); err != nil {
		blog.ErrorJSON("validTransferResourceDirParams failed, find %s failed, filter: %s, err: %s, rid: %s", common.BKTableNameModuleHostConfig, opt, err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	existHostIDs := make([]int64, 0)
	for _, host := range moduleHost {
		existHostIDs = append(existHostIDs, host.HostID)
	}

	return existHostIDs, nil
}
