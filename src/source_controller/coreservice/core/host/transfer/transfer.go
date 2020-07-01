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
)

type genericTransfer struct {
	dbProxy   dal.DB
	eventCli  eventclient.Client
	dependent OperationDependence

	// depend parameter
	moduleIDArr []int64
	bizID       int64
	// Incr=true is added to the module
	// Incr=false delete exist module, add current module
	isIncrement bool
	// cross-business transfer module
	// From the A business to the B business module
	crossBizTransfer bool
	// cross-business transfer module, source business id
	srcBizID int64

	// delHost delete host model
	delHost bool

	// ***** cache ********
	// inner module id array
	innerModuleID []int64
	// map[bk_module_id]bk_set_id
	moduleIDSetIDMap map[int64]int64
}

// validParameter valid parameter legal
func (t *genericTransfer) ValidParameter(kit *rest.Kit) errors.CCErrorCoder {
	if len(t.innerModuleID) == 0 {
		err := t.getInnerModuleIDArr(kit)
		if err != nil {
			return err
		}
	}

	err := t.validParameterInst(kit)
	if err != nil {
		return err
	}

	err = t.validParameterModule(kit)
	if err != nil {
		return err
	}

	return nil
}

// SetCrossBusiness Set host cross-service transfer parameters
func (t *genericTransfer) SetCrossBusiness(kit *rest.Kit, bizID int64) {
	t.crossBizTransfer = true
	t.srcBizID = bizID
}

// SetCrossBusiness Set host cross-service transfer parameters
func (t *genericTransfer) SetDeleteHost(kit *rest.Kit) {
	t.delHost = true
}

func (t *genericTransfer) Transfer(kit *rest.Kit, hostID int64) errors.CCErrorCoder {
	err := t.validHost(kit, hostID)
	if err != nil {
		return err
	}

	// hostInfo
	var hostInfo mapstr.MapStr
	// transfer  host module config
	var originDatas, curDatas []mapstr.MapStr
	// must be slice ptr address, Each assignment will change the address
	defer t.generateEvent(kit, &originDatas, &curDatas, hostInfo)

	// remove service instance if necessary
	if err := t.removeHostServiceInstance(kit, hostID); err != nil {
		return err
	}

	originDatas, err = t.delHostModuleRelation(kit, hostID)
	if err != nil {
		// It is not the time to merge and base the time. When it fails,
		// it is clear that the data before the change is pushed.
		// t.origindatas = nil
		return err
	}
	// delete host.
	if t.delHost {
		hostInfo, err = t.deleteHost(kit, hostID)
		if err != nil {
			return err
		}
		return nil

	}
	// transfer host module config
	curDatas, err = t.addHostModuleRelation(kit, hostID)
	if err != nil {
		return err
	}

	// auto create service instance if necessary
	if err := t.autoCreateServiceInstance(kit, hostID); err != nil {
		return err
	}

	return nil
}

func (t *genericTransfer) deleteHost(kit *rest.Kit, hostID int64) (mapstr.MapStr, errors.CCErrorCoder) {
	hostCond := condition.CreateCondition()
	hostCond.Field(common.BKHostIDField).Eq(hostID)
	hostCondMap := util.SetQueryOwner(hostCond.ToMapStr(), kit.SupplierAccount)
	hostInfoArr := make([]metadata.HostMapStr, 0)
	err := t.dbProxy.Table(common.BKTableNameBaseHost).Find(&hostCondMap).All(kit.Ctx, &hostInfoArr)
	if err != nil {
		blog.ErrorJSON("deleteHost find data error. err:%s, cond:%s, rid:%s", err.Error(), hostCondMap, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if len(hostInfoArr) == 0 {
		blog.ErrorJSON("deleteHost not found host error. cond:%s, rid:%s", hostCond.ToMapStr(), kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotExist, hostID)
	}
	delMoudleHost := condition.CreateCondition()
	delMoudleHost.Field(common.BKHostIDField).Eq(hostID)
	delMoudleHost.Field(common.BKAppIDField).Eq(t.bizID)
	delMoudleHostMap := util.SetQueryOwner(delMoudleHost.ToMapStr(), kit.SupplierAccount)
	err = t.dbProxy.Table(common.BKTableNameModuleHostConfig).Delete(kit.Ctx, delMoudleHostMap)
	if err != nil {
		blog.ErrorJSON("deleteHost delete module hsot realtion error. err:%s, cond:%s, rid:%s", err.Error(), delMoudleHostMap, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	err = t.dbProxy.Table(common.BKTableNameBaseHost).Delete(kit.Ctx, hostCondMap)
	if err != nil {
		blog.ErrorJSON("deleteHost delete host error. err:%s, cond:%s, rid:%s", err.Error(), hostCondMap, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return mapstr.MapStr(hostInfoArr[0]), nil
}

// generateEvent handle event trigger.
// Data from before and after changes cannot be merged for historical reasons.
func (t *genericTransfer) generateEvent(kit *rest.Kit, originDatas, curDatas *[]mapstr.MapStr, hostInfo mapstr.MapStr) errors.CCErrorCoder {

	var eventArr []*metadata.EventInst
	for _, data := range *originDatas {
		event := eventclient.NewEventWithHeader(kit.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = metadata.EventObjTypeModuleTransfer
		event.Action = metadata.EventActionDelete
		event.Data = []metadata.EventData{
			{PreData: data},
		}
		eventArr = append(eventArr, event)

	}
	for _, data := range *curDatas {
		event := eventclient.NewEventWithHeader(kit.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = metadata.EventObjTypeModuleTransfer
		event.Action = metadata.EventActionCreate
		event.Data = []metadata.EventData{
			{CurData: data},
		}
		eventArr = append(eventArr, event)
	}
	if len(hostInfo) > 0 {
		if t.delHost {
			event := eventclient.NewEventWithHeader(kit.Header)
			event.EventType = metadata.EventTypeInstData
			event.ObjType = common.BKInnerObjIDHost
			event.Action = metadata.EventActionDelete
			event.Data = []metadata.EventData{
				{
					PreData: hostInfo,
				},
			}
		}

	}
	err := t.eventCli.Push(kit.Ctx, eventArr...)
	if err != nil {
		blog.Errorf("host relation event push failed, but create event error:%v, rid: %s", err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)
	}

	return nil
}

// validParameterInst  validate module, biz, srcBiz must be exist
func (t *genericTransfer) validParameterInst(kit *rest.Kit) errors.CCErrorCoder {

	appCond := condition.CreateCondition()
	appCond.Field(common.BKAppIDField).Eq(t.bizID)

	cnt, err := t.countByCond(kit, appCond.ToMapStr(), common.BKTableNameBaseApp)
	if err != nil {
		return err
	}
	if cnt == 0 {
		blog.ErrorJSON("validParameter not business host error. cond:%s, rid:%s", appCond.ToMapStr(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceBusinessNotExist, t.bizID)
	}
	// cross-business validation source business
	if t.crossBizTransfer {
		appCond := condition.CreateCondition()
		appCond.Field(common.BKAppIDField).Eq(t.srcBizID)

		cnt, err = t.countByCond(kit, appCond.ToMapStr(), common.BKTableNameBaseApp)
		if err != nil {
			return err
		}
		if cnt == 0 {
			blog.ErrorJSON("validParameter not cross-business host error. cond:%s, rid:%s", appCond.ToMapStr(), kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCoreServiceBusinessNotExist, t.srcBizID)
		}
	}
	return nil
}

// validParameterModule validate parameter module legal
// module must be exist in business
// multiple modules not default module, transfer default must be one module
func (t *genericTransfer) validParameterModule(kit *rest.Kit) errors.CCErrorCoder {
	// delete host not validation destination module
	if t.delHost {
		return nil
	}
	if len(t.moduleIDArr) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKModuleIDField)
	}
	bizID := t.bizID

	t.moduleIDArr = util.IntArrayUnique(t.moduleIDArr)
	moduleInfoArr, err := t.getModuleInfoByModuleID(kit, bizID, t.moduleIDArr, []string{common.BKModuleIDField, common.BKDefaultField, common.BKSetIDField})
	if err != nil {
		return err
	}
	//  存在不属于当前业务的模块
	if len(moduleInfoArr) != len(t.moduleIDArr) {
		blog.Errorf("validParameterModule not found module info. moduleID:%#v,bizID:%d,rid:%s", t.moduleIDArr, bizID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceHasModuleNotBelongBusiness, t.moduleIDArr, bizID)
	}

	t.moduleIDSetIDMap = make(map[int64]int64, 0)

	// When multiple modules are used, determine whether the default module .
	// has default module ,not handle transfer.
	for _, moduleInfo := range moduleInfoArr {
		// 当为多个模块的时候，不能包含默认模块。 单个模块下， 不能用附加功能。
		if len(t.moduleIDArr) != 1 || t.isIncrement {
			// 转移目标模块为多模块时，不允许包含内置模块(空闲机/故障机等)
			defaultVal, err := moduleInfo.Int64(common.BKDefaultField)
			if err != nil {
				blog.ErrorJSON("validParameter module info field default  not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfo, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule, common.BKDefaultField, "int", err.Error())
			}
			if defaultVal != 0 {
				blog.ErrorJSON("validParameter module info field  has default module.  moduleInfo:%s,rid:%s", defaultVal, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCoreServiceModuleContainDefaultModuleErr)
			}
		}
		moduleID, err := moduleInfo.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field module id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		setID, err := moduleInfo.Int64(common.BKSetIDField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field set id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule, common.BKSetIDField, "int", err.Error())
		}
		t.moduleIDSetIDMap[moduleID] = setID

	}

	return nil
}

// validParameterHostBelongbiz  legal
// check if the host belongs to the transfer business.
// check host exist
func (t *genericTransfer) validHost(kit *rest.Kit, hostID int64) errors.CCErrorCoder {
	hostCond := condition.CreateCondition()
	hostCond.Field(common.BKHostIDField).Eq(hostID)

	cnt, err := t.countByCond(kit, hostCond.ToMapStr(), common.BKTableNameBaseHost)
	if err != nil {
		return err
	}
	if cnt == 0 {
		blog.ErrorJSON("validParameter not found host error. cond:%s, rid:%s", hostCond.ToMapStr(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotExist, hostID)
	}

	bizID := t.bizID
	// transfer the host across businees,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).NotEq(bizID)
	cond.Field(common.BKHostIDField).Eq(hostID)
	condMap := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)

	cnt, dbErr := t.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(condMap).Count(kit.Ctx)
	if dbErr != nil {
		blog.ErrorJSON("validParameterHostBelongbiz find data error. err:%s,cond:%s, rid:%s", dbErr.Error(), condMap, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if cnt > 0 {
		blog.ErrorJSON("validParameterHostBelongbiz has belong to other business.cond:%s, rid:%s", condMap, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotBelongBusiness, hostID, bizID)
	}
	return nil
}

// delHostModuleRelation delete single host module relation
func (t *genericTransfer) delHostModuleRelation(kit *rest.Kit, hostID int64) ([]mapstr.MapStr, errors.CCErrorCoder) {
	bizID := t.bizID
	// transfer the host across business,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}

	if t.isIncrement {
		// delete default module
		return t.delHostModuleRelationItem(kit, bizID, hostID, true)

	} else {
		// delete all module
		return t.delHostModuleRelationItem(kit, bizID, hostID, false)
	}
}

// delHostModuleRelationItem delete single host module relation
func (t *genericTransfer) delHostModuleRelationItem(kit *rest.Kit, bizID, hostID int64, isDefault bool) ([]mapstr.MapStr, errors.CCErrorCoder) {

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	if isDefault {
		cond.Field(common.BKModuleIDField).In(t.innerModuleID)
	}
	cond.Field(common.BKHostIDField).Eq(hostID)

	delCondition := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)

	// retrieve original data
	originDatas := make([]mapstr.MapStr, 0)
	getErr := t.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(delCondition).All(kit.Ctx, &originDatas)
	if getErr != nil {
		blog.ErrorJSON("delete host relation, retrieve original data error. err:%v, cond:%s, rid:%s", getErr, delCondition, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	delCondition = util.SetModOwner(cond.ToMapStr(), kit.SupplierAccount)
	delErr := t.dbProxy.Table(common.BKTableNameModuleHostConfig).Delete(kit.Ctx, delCondition) //.DelByCondition(ModuleHostCollection, delCondition)
	if delErr != nil {
		blog.ErrorJSON("delete host relation, but del module host relation failed. err:%v, cond:%s, rid:%s", delErr, delCondition, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return originDatas, nil
}

// AddSingleHostModuleRelation add single host module relation
func (t *genericTransfer) addHostModuleRelation(kit *rest.Kit, hostID int64) ([]mapstr.MapStr, errors.CCErrorCoder) {
	bizID := t.bizID

	var moduleIDArr []int64
	// append method, filter already exist modules
	if t.isIncrement {
		cond := condition.CreateCondition()
		cond.Field(common.BKAppIDField).Eq(t.bizID)
		cond.Field(common.BKHostIDField).Eq(hostID)
		cond.Field(common.BKModuleIDField).In(t.moduleIDArr)
		condMap := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)
		relationArr := make([]metadata.ModuleHost, 0)
		err := t.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(condMap).All(kit.Ctx, &relationArr)
		if err != nil {
			blog.ErrorJSON("add host relation, retrieve original data error. err:%v, cond:%s, rid:%s", err, condMap, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		// map[moduleID]bool
		existModuleIDMap := make(map[int64]bool, 0)
		for _, item := range relationArr {
			existModuleIDMap[item.ModuleID] = true
		}
		for _, moduleID := range t.moduleIDArr {
			if _, ok := existModuleIDMap[moduleID]; !ok {
				moduleIDArr = append(moduleIDArr, moduleID)
			}
		}

	} else {
		moduleIDArr = t.moduleIDArr
	}
	if len(moduleIDArr) == 0 {
		return nil, nil
	}
	var insertDataArr []mapstr.MapStr
	for _, moduleID := range moduleIDArr {
		insertData := mapstr.MapStr{
			common.BKAppIDField: bizID, common.BKHostIDField: hostID, common.BKModuleIDField: moduleID,
			// validation parameter ensure module must be exist  t.validParameterModule
			common.BKSetIDField: t.moduleIDSetIDMap[moduleID],
		}

		insertData = util.SetModOwner(insertData, kit.SupplierAccount)
		insertDataArr = append(insertDataArr, insertData)
	}

	err := t.dbProxy.Table(common.BKTableNameModuleHostConfig).Insert(kit.Ctx, insertDataArr)
	if err != nil {
		blog.Errorf("add host module relation, add module host relation error: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	return insertDataArr, nil
}

func (t *genericTransfer) autoCreateServiceInstance(kit *rest.Kit, hostID int64) errors.CCErrorCoder {
	for _, moduleID := range t.moduleIDArr {
		if _, err := t.dependent.AutoCreateServiceInstanceModuleHost(kit, hostID, moduleID); err != nil {
			blog.Warnf("autoCreateServiceInstance failed, hostID: %d, moduleID: %d, rid: %s", hostID, moduleID, kit.Rid)
		}
	}
	return nil
}

// remove service instances bound to hosts with process instances in certain modules
func (t *genericTransfer) removeHostServiceInstance(kit *rest.Kit, hostID int64) errors.CCErrorCoder {
	// increment transfer don't need to remove service instance
	if t.isIncrement {
		return nil
	}
	// get all service instance IDs that need to be removed
	serviceInstanceFilter := map[string]interface{}{
		common.BKHostIDField: hostID,
	}
	if len(t.moduleIDArr) > 0 {
		serviceInstanceFilter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBNIN: t.moduleIDArr,
		}
	}
	instances := make([]metadata.ServiceInstance, 0)
	err := t.dbProxy.Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).Fields(common.BKFieldID).All(kit.Ctx, &instances)
	if err != nil {
		blog.ErrorJSON("removeHostServiceInstance failed, get service instance IDs failed, err: %s, filter: %s, rid: %s", err, serviceInstanceFilter, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if len(instances) == 0 {
		return nil
	}
	serviceInstanceIDs := make([]int64, 0)
	for _, instance := range instances {
		serviceInstanceIDs = append(serviceInstanceIDs, instance.ID)
	}

	// get all process IDs of the service instances to be removed
	processRelationFilter := map[string]interface{}{
		common.BKServiceInstanceIDField: map[string]interface{}{
			common.BKDBIN: serviceInstanceIDs,
		},
	}
	relations := make([]metadata.ProcessInstanceRelation, 0)
	if err := t.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(processRelationFilter).All(kit.Ctx, &relations); nil != err {
		blog.Errorf("removeHostServiceInstance failed, get process instance relation failed, err: %s, filter: %s, rid: %s", err, processRelationFilter, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	processIDs := make([]int64, 0)
	for _, relation := range relations {
		processIDs = append(processIDs, relation.ProcessID)
	}

	// delete all process relations and instances
	if len(processIDs) > 0 {
		if err := t.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Delete(kit.Ctx, processRelationFilter); nil != err {
			blog.Errorf("removeHostServiceInstance failed, delete process instance relation failed, err: %s, filter: %s, rid: %s", err, processRelationFilter, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
		}

		processFilter := map[string]interface{}{
			common.BKProcessIDField: map[string]interface{}{
				common.BKDBIN: processIDs,
			},
		}
		if err := t.dbProxy.Table(common.BKTableNameBaseProcess).Delete(kit.Ctx, processFilter); nil != err {
			blog.Errorf("removeHostServiceInstance failed, delete process instances failed, err: %s, filter: %s, rid: %s", err, processFilter, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
		}
	}

	// delete service instances
	serviceInstanceIDFilter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: serviceInstanceIDs,
		},
	}
	if err := t.dbProxy.Table(common.BKTableNameServiceInstance).Delete(kit.Ctx, serviceInstanceIDFilter); nil != err {
		blog.Errorf("removeHostServiceInstance failed, delete service instances failed, err: %s, filter: %s, rid: %s", err, serviceInstanceIDFilter, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	return nil
}

// getInnerModuleIDArr get default module
func (t *genericTransfer) getInnerModuleIDArr(kit *rest.Kit) errors.CCErrorCoder {
	bizID := t.bizID
	// transfer the host across business,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}
	moduleConds := condition.CreateCondition()
	moduleConds.Field(common.BKAppIDField).Eq(bizID)
	moduleConds.Field(common.BKDefaultField).NotEq(common.DefaultFlagDefaultValue)
	cond := util.SetQueryOwner(moduleConds.ToMapStr(), kit.SupplierAccount)

	moduleInfoArr := make([]mapstr.MapStr, 0)
	err := t.dbProxy.Table(common.BKTableNameBaseModule).Find(cond).All(kit.Ctx, &moduleInfoArr)

	if err != nil {
		blog.ErrorJSON("getInnerModuleIDArr find data error. err:%s,cond:%s, rid:%s", err.Error(), cond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if len(moduleInfoArr) == 0 {
		blog.Warnf("getInnerModuleIDArr not found default module. appID:%d, rid:%s", bizID, kit.Rid)
	}
	for _, moduleInfo := range moduleInfoArr {
		moduleID, err := moduleInfo.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("getInnerModuleIDArr module info field module id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfo, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		t.innerModuleID = append(t.innerModuleID, moduleID)
	}

	return nil
}

func (t *genericTransfer) GetInnerModuleIDArr(kit *rest.Kit) ([]int64, errors.CCError) {
	if len(t.innerModuleID) == 0 {
		err := t.getInnerModuleIDArr(kit)
		return t.innerModuleID, err
	}
	return t.innerModuleID, nil
}

func (t *genericTransfer) HasInnerModule(kit *rest.Kit) (bool, error) {
	innerModuleIDArr, err := t.GetInnerModuleIDArr(kit)
	if err != nil {
		return false, err
	}
	if len(innerModuleIDArr) == 0 {
		blog.ErrorJSON("HasInnerModule  error. module:%s, rid:%s", t.moduleIDArr, kit.Rid)
		return false, kit.CCError.CCErrorf(common.CCErrCoreServiceDefaultModuleNotExist, t.bizID)
	}
	for _, innerModuleID := range innerModuleIDArr {
		for _, moduleID := range t.moduleIDArr {
			if moduleID == innerModuleID {
				return true, nil
			}
		}

	}
	return false, nil
}

func (t *genericTransfer) getModuleInfoByModuleID(kit *rest.Kit, appID int64, moduleID []int64, fields []string) ([]mapstr.MapStr, errors.CCErrorCoder) {
	moduleConds := condition.CreateCondition()
	moduleConds.Field(common.BKAppIDField).Eq(appID)
	moduleConds.Field(common.BKModuleIDField).In(moduleID)
	cond := util.SetQueryOwner(moduleConds.ToMapStr(), kit.SupplierAccount)

	moduleInfoArr := make([]mapstr.MapStr, 0)
	err := t.dbProxy.Table(common.BKTableNameBaseModule).Find(cond).Fields(fields...).All(kit.Ctx, &moduleInfoArr)
	if err != nil {
		blog.ErrorJSON("getModuleInfoByModuleID find data CCErrorCoder. err:%s,cond:%s, rid:%s", err.Error(), cond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return moduleInfoArr, nil
}

func (t *genericTransfer) countByCond(kit *rest.Kit, conds mapstr.MapStr, tableName string) (uint64, errors.CCErrorCoder) {
	conds = util.SetQueryOwner(conds, kit.SupplierAccount)
	cnt, err := t.dbProxy.Table(tableName).Find(conds).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("countByCond find data error. err:%s, table:%s,cond:%s, rid:%s", err.Error(), tableName, conds, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return cnt, nil
}
