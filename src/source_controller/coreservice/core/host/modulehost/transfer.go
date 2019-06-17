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

package modulehost

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

type transferHostModule struct {
	dependent OperationDependence

	// depend parametere
	mh          *ModuleHost
	moduleIDArr []int64
	bizID       int64
	// Incr=true is added to the module
	// Incr=false delete existing module, add current module
	isIncr bool
	// cross-business transfer module
	// From the A business to the B business module
	crossBizTransfer bool
	//   cross-business transfer module, source business id
	srcBizID int64

	// delHost delete host model
	delHost bool

	// handle data

	// inner module id array
	innerModuleID []int64

	// map[module id]set id
	moduleIDSetIDmap map[int64]int64
}

// NewHostModuleTransfer business normal module transfer
func (mh *ModuleHost) NewHostModuleTransfer(ctx core.ContextParams, bizID int64, moduleIDArr []int64, isIncr bool) *transferHostModule {
	return &transferHostModule{
		mh:          mh,
		moduleIDArr: moduleIDArr,
		bizID:       bizID,
		isIncr:      isIncr,
		dependent:   mh.dependence,
	}
}

// validParameter valid parametere legal
func (t *transferHostModule) ValidParameter(ctx core.ContextParams) errors.CCErrorCoder {
	if len(t.innerModuleID) == 0 {
		err := t.getInnerModuleIDArr(ctx)
		if err != nil {
			return err
		}
	}

	err := t.validParameterInst(ctx)
	if err != nil {
		return err
	}

	err = t.validParameterModule(ctx)
	if err != nil {
		return err
	}

	return nil
}

// SetCrossBusiness Set host cross-service transfer parameters
func (t *transferHostModule) SetCrossBusiness(ctx core.ContextParams, bizID int64) {
	t.crossBizTransfer = true
	t.srcBizID = bizID
}

// SetCrossBusiness Set host cross-service transfer parameters
func (t *transferHostModule) SetDeleteHost(ctx core.ContextParams) {
	t.delHost = true
}

func (t *transferHostModule) Transfer(ctx core.ContextParams, hostID int64) errors.CCErrorCoder {
	err := t.validHost(ctx, hostID)
	if err != nil {
		return err
	}

	// hostInfo
	var hostInfo mapstr.MapStr
	// transfer  host module config
	var originDatas, curDatas []mapstr.MapStr
	// must be silce ptr address, Each assignment will change the address
	defer t.generateEvent(ctx, &originDatas, &curDatas, hostInfo)

	originDatas, err = t.delHostModuleRelation(ctx, hostID)
	if err != nil {
		// It is not the time to merge and base the time. When it fails,
		// it is clear that the data before the change is pushed.
		// t.origindatas = nil
		return err
	}
	// delete host.
	if t.delHost {
		hostInfo, err = t.deleteHost(ctx, hostID)
		if err != nil {
			return err
		}
		return nil

	}
	// transfer host module config
	curDatas, err = t.AddHostModuleRelation(ctx, hostID)
	if err != nil {
		return err
	}

	// auto create service instance if necessary
	if err := t.autoCreateServiceInstance(ctx, hostID); err != nil {
		return err
	}

	return nil
}

func (t *transferHostModule) deleteHost(ctx core.ContextParams, hostID int64) (mapstr.MapStr, errors.CCErrorCoder) {
	hostCond := condition.CreateCondition()
	hostCond.Field(common.BKHostIDField).Eq(hostID)
	hostCondMap := util.SetQueryOwner(hostCond.ToMapStr(), ctx.SupplierAccount)
	hostInfoArr := make([]mapstr.MapStr, 0)
	err := t.mh.dbProxy.Table(common.BKTableNameBaseHost).Find(&hostCondMap).All(ctx, &hostInfoArr)
	if err != nil {
		blog.ErrorJSON("deleteHost find data error. err:%s, cond:%s, rid:%s", err.Error(), hostCondMap, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if len(hostInfoArr) == 0 {
		blog.ErrorJSON("deleteHost not found host error. cond:%s, rid:%s", hostCond.ToMapStr(), ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCoreServiceHostNotExist, hostID)
	}
	delMoudleHost := condition.CreateCondition()
	delMoudleHost.Field(common.BKHostIDField).Eq(hostID)
	delMoudleHost.Field(common.BKAppIDField).Eq(t.bizID)
	delMoudleHostMap := util.SetQueryOwner(delMoudleHost.ToMapStr(), ctx.SupplierAccount)
	err = t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Delete(ctx, delMoudleHostMap)
	if err != nil {
		blog.ErrorJSON("deleteHost delete module hsot realtion error. err:%s, cond:%s, rid:%s", err.Error(), delMoudleHostMap, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	err = t.mh.dbProxy.Table(common.BKTableNameBaseHost).Delete(ctx, hostCondMap)
	if err != nil {
		blog.ErrorJSON("deleteHost delete host error. err:%s, cond:%s, rid:%s", err.Error(), hostCondMap, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return hostInfoArr[0], nil
}

// generateEvent handle event trigger.
// Data from before and after changes cannot be merged for historical reasons.
func (t *transferHostModule) generateEvent(ctx core.ContextParams, originDatas, curDatas *[]mapstr.MapStr, hostInfo mapstr.MapStr) errors.CCErrorCoder {

	var eventArr []*metadata.EventInst
	for _, data := range *originDatas {
		event := eventclient.NewEventWithHeader(ctx.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = "moduletransfer"
		event.Action = metadata.EventActionDelete
		event.Data = []metadata.EventData{
			{PreData: data},
		}
		eventArr = append(eventArr, event)

	}
	for _, data := range *curDatas {
		event := eventclient.NewEventWithHeader(ctx.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = "moduletransfer"
		event.Action = metadata.EventActionCreate
		event.Data = []metadata.EventData{
			{CurData: data},
		}
		eventArr = append(eventArr, event)
	}
	if len(hostInfo) > 0 {
		if t.delHost {
			event := eventclient.NewEventWithHeader(ctx.Header)
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
	err := t.mh.eventC.Push(ctx, eventArr...)
	if err != nil {
		blog.Errorf("host relation event push failed, but create event error:%v", err)
		return ctx.Error.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)
	}

	return nil
}

// validParameterInst  validate module, biz, srcBiz must be exist
func (t *transferHostModule) validParameterInst(ctx core.ContextParams) errors.CCErrorCoder {

	appCond := condition.CreateCondition()
	appCond.Field(common.BKAppIDField).Eq(t.bizID)

	cnt, err := t.mh.countByCond(ctx, appCond.ToMapStr(), common.BKTableNameBaseApp)
	if err != nil {
		return err
	}
	if cnt == 0 {
		blog.ErrorJSON("validParameter not business host error. cond:%s, rid:%s", appCond.ToMapStr(), ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCoreServiceBusinessNotExist, t.bizID)
	}
	// cross-business validation source business
	if t.crossBizTransfer {
		appCond := condition.CreateCondition()
		appCond.Field(common.BKAppIDField).Eq(t.srcBizID)

		cnt, err = t.mh.countByCond(ctx, appCond.ToMapStr(), common.BKTableNameBaseApp)
		if err != nil {
			return err
		}
		if cnt == 0 {
			blog.ErrorJSON("validParameter not cross-business host error. cond:%s, rid:%s", appCond.ToMapStr(), ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCoreServiceBusinessNotExist, t.srcBizID)
		}
	}
	return nil
}

// validParameterModule validate parameter module legal
// module must be exist in business
// multiple modules not default module, transfer default must be one module
func (t *transferHostModule) validParameterModule(ctx core.ContextParams) errors.CCErrorCoder {
	// delete host not validation destination module
	if t.delHost {
		return nil
	}
	if len(t.moduleIDArr) == 0 {
		return ctx.Error.CCErrorf(common.CCErrCommParamsNeedSet, common.BKModuleIDField)
	}
	bizID := t.bizID

	moduleInfoArr, err := t.mh.getModuleInfoByModuleID(ctx, bizID, t.moduleIDArr, []string{common.BKModuleIDField, common.BKDefaultField, common.BKSetIDField})
	if err != nil {
		return err
	}
	//  存在不属于当前业务的模块
	if len(moduleInfoArr) != len(t.moduleIDArr) {
		blog.Errorf("validParameterModule not found module info. moduleID:%#v,bizID:%d,rid:%s", t.moduleIDArr, bizID, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCoreServiceHasModuleNotBelongBusiness, t.moduleIDArr, bizID)
	}

	t.moduleIDSetIDmap = make(map[int64]int64, 0)

	// When multiple modules are used, determine whether the default module .
	// has default module ,not handle transfer.
	for _, moduleInfo := range moduleInfoArr {
		if len(t.moduleIDArr) != 1 {

			// 转移目标模块为多模块时，不允许包含内置模块(空闲机/故障机等)
			defaultVal, err := moduleInfo.Int64(common.BKDefaultField)
			if err != nil {
				blog.ErrorJSON("validParameter module info field default  not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfo, ctx.ReqID)
				return ctx.Error.CCErrorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKDefaultField, "int", err.Error())
			}
			if defaultVal != 0 {
				blog.ErrorJSON("validParameter module info field  has default module.  moduleInfo:%s,rid:%s", defaultVal, ctx.ReqID)
				return ctx.Error.CCErrorf(common.CCErrCoreServiceModuleContainDefaultModuleErr)
			}
		}
		moduleID, err := moduleInfo.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field module id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		setID, err := moduleInfo.Int64(common.BKSetIDField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field set id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKSetIDField, "int", err.Error())
		}
		t.moduleIDSetIDmap[moduleID] = setID

	}

	return nil
}

// validParameterHostBelongbiz  legal
// check if the host belongs to the transfer business.
// check host exist
func (t *transferHostModule) validHost(ctx core.ContextParams, hostID int64) errors.CCErrorCoder {
	hostCond := condition.CreateCondition()
	hostCond.Field(common.BKHostIDField).Eq(hostID)

	cnt, err := t.mh.countByCond(ctx, hostCond.ToMapStr(), common.BKTableNameBaseHost)
	if err != nil {
		return err
	}
	if cnt == 0 {
		blog.ErrorJSON("validParameter not found host error. cond:%s, rid:%s", hostCond.ToMapStr(), ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCoreServiceHostNotExist, hostID)
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
	condMap := util.SetQueryOwner(cond.ToMapStr(), ctx.SupplierAccount)

	cnt, dbErr := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(condMap).Count(ctx)
	if dbErr != nil {
		blog.ErrorJSON("validParameterHostBelongbiz find data error. err:%s,cond:%s, rid:%s", dbErr.Error(), condMap, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if cnt > 0 {
		blog.ErrorJSON("validParameterHostBelongbiz has belong to other business.cond:%s, rid:%s", condMap, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCoreServiceHostNotBelongBusiness, hostID, bizID)
	}
	return nil
}

// delHostModuleRelation delete single host module relation
func (t *transferHostModule) delHostModuleRelation(ctx core.ContextParams, hostID int64) ([]mapstr.MapStr, errors.CCErrorCoder) {
	bizID := t.bizID
	// transfer the host across business,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}
	//

	if t.isIncr {
		// delete default module
		return t.delHostModuleRelationItem(ctx, bizID, hostID, true)

	} else {
		// delete all module
		return t.delHostModuleRelationItem(ctx, bizID, hostID, false)
	}
}

// delHostModuleRelationItem delete single host module relation
func (t *transferHostModule) delHostModuleRelationItem(ctx core.ContextParams, bizID, hostID int64, isDefault bool) ([]mapstr.MapStr, errors.CCErrorCoder) {

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	if isDefault {
		cond.Field(common.BKModuleIDField).In(t.innerModuleID)
	}
	cond.Field(common.BKHostIDField).Eq(hostID)

	delCondition := util.SetQueryOwner(cond.ToMapStr(), ctx.SupplierAccount)
	num, numError := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(delCondition).Count(ctx)
	if numError != nil {
		blog.Errorf("delete host relation, but get module host relation failed, err: %v", numError)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if num == 0 {
		return nil, nil
	}

	// retrieve original data
	originDatas := make([]mapstr.MapStr, 0)
	getErr := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(delCondition).All(ctx, &originDatas)
	if getErr != nil {
		blog.ErrorJSON("delete host relation, retrieve original data error. err:%v, cond:%s, rid:%s", getErr, delCondition, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	delCondition = util.SetModOwner(cond.ToMapStr(), ctx.SupplierAccount)
	delErr := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Delete(ctx, delCondition) //.DelByCondition(ModuleHostCollection, delCondition)
	if delErr != nil {
		blog.ErrorJSON("delete host relation, but del module host relation failed. err:%v, cond:%s, rid:%s", delErr, delCondition, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return originDatas, nil
}

// AddSingleHostModuleRelation add single host module relation
func (t *transferHostModule) AddHostModuleRelation(ctx core.ContextParams, hostID int64) ([]mapstr.MapStr, errors.CCErrorCoder) {
	bizID := t.bizID

	var moduleIDArr []int64
	// append method, filter already exist modules
	if t.isIncr {
		cond := condition.CreateCondition()
		cond.Field(common.BKAppIDField).Eq(t.bizID)
		cond.Field(common.BKHostIDField).Eq(hostID)
		cond.Field(common.BKModuleIDField).In(t.moduleIDArr)
		condMap := util.SetQueryOwner(cond.ToMapStr(), ctx.SupplierAccount)
		relationArr := make([]metadata.ModuleHost, 0)
		err := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(condMap).All(ctx, &relationArr)
		if err != nil {
			blog.ErrorJSON("add  host relation, retrieve original data error. err:%v, cond:%s, rid:%s", err, condMap, ctx.ReqID)
			return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
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
			common.BKSetIDField: t.moduleIDSetIDmap[moduleID],
		}

		insertData = util.SetModOwner(insertData, ctx.SupplierAccount)
		insertDataArr = append(insertDataArr, insertData)
	}

	err := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Insert(ctx, insertDataArr)
	if err != nil {
		blog.Errorf("add host module relation, add module host relation error: %v", err)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	return insertDataArr, nil
}

func (t *transferHostModule) autoCreateServiceInstance(ctx core.ContextParams, hostID int64) errors.CCErrorCoder {
	for _, moduleID := range t.moduleIDArr {
		if _, err := t.dependent.AutoCreateServiceInstanceModuleHost(ctx, hostID, moduleID); err != nil {
			blog.Warnf("autoCreateServiceInstance failed, hostID: %d, moduleID: %d, rid: %s", hostID, moduleID, ctx.ReqID)
		}
	}
	return nil
}

// getInnerModuleIDArr get default module
func (t *transferHostModule) getInnerModuleIDArr(ctx core.ContextParams) errors.CCErrorCoder {
	bizID := t.bizID
	// transfer the host across business,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}
	moduleConds := condition.CreateCondition()
	moduleConds.Field(common.BKAppIDField).Eq(bizID)
	moduleConds.Field(common.BKDefaultField).NotEq(0)
	cond := util.SetQueryOwner(moduleConds.ToMapStr(), ctx.SupplierAccount)

	moduleInfoArr := make([]mapstr.MapStr, 0)
	err := t.mh.dbProxy.Table(common.BKTableNameBaseModule).Find(cond).All(ctx, &moduleInfoArr)

	if err != nil {
		blog.ErrorJSON("getInnerModuleIDArr find data error. err:%s,cond:%s, rid:%s", err.Error(), cond, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if len(moduleInfoArr) == 0 {
		blog.Warnf("getInnerModuleIDArr not found default module. appID:%d, rid:%s", bizID, ctx.ReqID)
	}
	for _, moduleInfo := range moduleInfoArr {
		moduleID, err := moduleInfo.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("getInnerModuleIDArr module info field module id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfo, ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		t.innerModuleID = append(t.innerModuleID, moduleID)
	}

	return nil
}

func (t *transferHostModule) GetInnerModuleIDArr(ctx core.ContextParams) ([]int64, errors.CCError) {
	if len(t.innerModuleID) == 0 {
		err := t.getInnerModuleIDArr(ctx)
		return t.innerModuleID, err
	}
	return t.innerModuleID, nil
}

func (t *transferHostModule) HasInnerModule(ctx core.ContextParams) (bool, error) {
	innerModuleIDArr, err := t.GetInnerModuleIDArr(ctx)
	if err != nil {
		return false, err
	}
	if len(innerModuleIDArr) == 0 {
		blog.ErrorJSON("HasInnerModule  error. module:%s, rid:%s", t.moduleIDArr, ctx.ReqID)
		return false, ctx.Error.CCErrorf(common.CCErrCoreServiceDefaultModuleNotExist, t.bizID)
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

// DoTransferToInnerCheck check whether could be transfer to inner module
func (t *transferHostModule) DoTransferToInnerCheck(ctx core.ContextParams, hostIDs []int64) error {
	if len(hostIDs) == 0 {
		return nil
	}

	// check: 不能有服务实例/进程实例绑定主机实例
	filter := map[string]interface{}{
		common.BKHostIDField: map[string][]int64{common.BKDBIN: hostIDs},
	}
	var count uint64
	count, err := t.mh.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Count(ctx.Context)
	if err != nil {
		blog.Errorf("DoTransferToInnerCheck failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		return ctx.Error.CCErrorf(common.CCErrCoreServiceForbiddenReleaseHostReferencedByServiceInstance)
	}

	return nil
}
