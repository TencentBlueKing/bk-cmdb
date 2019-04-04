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

	// handle data

	// default module id array
	defaultModuleID []int64
	// transfer before host module config
	originDatas []mapstr.MapStr
	// transfer host module cofnig
	curDatas []mapstr.MapStr

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
	}
}

// validParameter valid parametere legal
func (t *transferHostModule) ValidParameter(ctx core.ContextParams) errors.CCErrorCoder {
	if len(t.defaultModuleID) == 0 {
		err := t.getDefaultModuleIDArr(ctx)
		return err
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

func (t *transferHostModule) Transfer(ctx core.ContextParams, hostID int64) errors.CCErrorCoder {
	err := t.validHost(ctx, hostID)
	if err != nil {
		return err
	}
	err = t.delHostModuleRelation(ctx, hostID)
	if err != nil {
		// It is not the time to merge and base the time. When it fails,
		// it is clear that the data before the change is pushed.
		//t.origindatas = nil
		return err
	}

	err = t.AddHostModuleRelation(ctx, hostID)
	if err != nil {
		return err
	}

	return nil
}

// generateEvent handle event trigger.
// Data from before and after changes cannot be merged for historical reasons.
func (t *transferHostModule) generateEvent(ctx core.ContextParams) errors.CCErrorCoder {
	var eventArr []*metadata.EventInst
	for _, data := range t.originDatas {
		event := eventclient.NewEventWithHeader(ctx.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = "moduletransfer"
		event.Action = metadata.EventActionDelete
		event.Data = []metadata.EventData{
			{PreData: data},
		}
		eventArr = append(eventArr, event)

	}
	for _, data := range t.curDatas {
		event := eventclient.NewEventWithHeader(ctx.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = "moduletransfer"
		event.Action = metadata.EventActionCreate
		event.Data = []metadata.EventData{
			{CurData: data},
		}
		eventArr = append(eventArr, event)
	}
	err := t.mh.EventC.Push(ctx, eventArr...)
	if err != nil {
		blog.Errorf("host relation event push failed, but create event error:%v", err)
		return err
	}
	t.originDatas = nil
	t.curDatas = nil
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
		return ctx.Error.Errorf(common.CCErrCoreServiceBusinessNotExist, t.bizID)
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
			return ctx.Error.Errorf(common.CCErrCoreServiceBusinessNotExist, t.srcBizID)
		}
	}
	return nil
}

// validParameterModule validate parameter module legal
// module must be exist in business
// multiple modules not default module, transfer default must be one module
func (t *transferHostModule) validParameterModule(ctx core.ContextParams) errors.CCErrorCoder {
	bizID := t.bizID
	// transfer the host across businees,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}

	moduleInfoArr, err := t.mh.getModuleInfoByModuleID(ctx, bizID, t.moduleIDArr, []string{common.BKModuleIDField, common.BKDefaultField})
	if err != nil {
		return err
	}
	//  存在不属于当前业务的模块
	if len(moduleInfoArr) != len(t.moduleIDArr) {
		blog.Errorf("validParameterModule not found module info. moduleID:%#v,bizID,rid:%s", t.moduleIDArr, bizID, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCoreServiceHasModuleNotBelongBusiness, t.moduleIDArr, bizID)
	}

	t.moduleIDSetIDmap = make(map[int64]int64, 0)
	// 只有一个模块，不许做其他的判断
	if len(t.moduleIDArr) == 1 {
		moduleID, err := moduleInfoArr[0].Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field module id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		setID, err := moduleInfoArr[0].Int64(common.BKSetIDField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field set id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKSetIDField, "int", err.Error())
		}
		t.moduleIDSetIDmap[moduleID] = setID
		return nil
	}

	// When multiple modules are used, determine whether the default module .
	// has default module ,not handle transfer.
	for _, moduleInfo := range moduleInfoArr {
		defaultVal, err := moduleInfo.Int64(common.BKDefaultField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field default  not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfo, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKDefaultField, "int", err.Error())
		}
		if defaultVal != 0 {
			blog.ErrorJSON("validParameter module info field  has default module. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCoreServiceModuleContainDefaultModuleErr)
		}
		moduleID, err := moduleInfoArr[0].Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field module id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		setID, err := moduleInfoArr[0].Int64(common.BKSetIDField)
		if err != nil {
			blog.ErrorJSON("validParameter module info field set id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfoArr, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKSetIDField, "int", err.Error())
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
		return ctx.Error.Errorf(common.CCErrCoreServiceHostNotExist, hostID)
	}

	bizID := t.bizID
	// transfer the host across businees,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKHostIDField).Eq(hostID)
	condMap := util.SetQueryOwner(cond, ctx.SupplierAccount)

	cnt, err = t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(cond).Count(ctx)
	if err != nil {
		blog.ErrorJSON("validParameterHostBelongbiz find data error. err:%s,cond:%s, rid:%s", err.Error(), condMap, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}
	if cnt == 0 {
		blog.ErrorJSON("validParameterHostBelongbiz not found data.cond:%s, rid:%s", condMap, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCoreServiceHostNotBelongBusiness, hostID, bizID)
	}
	return nil
}

// delHostModuleRelation delete single host module relation
func (t *transferHostModule) delHostModuleRelation(ctx core.ContextParams, hostID int64) errors.CCErrorCoder {
	bizID := t.bizID
	// transfer the host across businees,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}
	//
	err := t.delHostModuleRelationItem(ctx, bizID, hostID, true)
	if err != nil {
		return err
	}
	if t.isIncr {
		return nil
	}
	return t.delHostModuleRelationItem(ctx, bizID, hostID, false)

}

// delHostModuleRelationItem delete single host module relation
func (t *transferHostModule) delHostModuleRelationItem(ctx core.ContextParams, bizID, hostID int64, isDefault bool) errors.CCErrorCoder {

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	if isDefault {
		// 当前新加关系中存在于默认模块一致，不删除当前模块的关系。
		// 在做参数验证的时候，保证转移到默认模块只有一个模块
		for _, moduleID := range t.defaultModuleID {
			if moduleID == t.moduleIDArr[0] {
				return nil
			}
		}
		cond.Field(common.BKModuleIDField).In(t.defaultModuleID)
	}
	cond.Field(common.BKHostIDField).Eq(hostID)

	delCondition := util.SetModOwner(cond.ToMapStr(), ctx.SupplierAccount)
	num, numError := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(delCondition).Count(ctx)
	if numError != nil {
		blog.Errorf("delete host relation, but get module host relation failed, err: %v", numError)
		return numError
	}

	if num == 0 {
		return nil
	}

	// retrieve original datas
	originDatas := make([]mapstr.MapStr, 0)
	getErr := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(delCondition).All(ctx, &originDatas)
	if getErr != nil {
		blog.ErrorJSON("delete host relation, retrieve original data error. err:%v, cond:%s, rid:%s", getErr, delCondition, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}

	delErr := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Delete(ctx, delCondition) //.DelByCondition(ModuleHostCollection, delCondition)
	if delErr != nil {
		blog.ErrorJSON("delete host relation, but del module host relation failed. err:%v, cond:%s, rid:%s", delErr, delCondition, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDBDeleteFailed)
	}
	t.originDatas = append(t.originDatas, originDatas...)

	return nil
}

//AddSingleHostModuleRelation add single host module relation
func (t *transferHostModule) AddHostModuleRelation(ctx core.ContextParams, hostID int64) errors.CCErrorCoder {
	bizID := t.bizID
	// transfer the host across businees,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}

	var moduleIDArr []int64
	// append method, filter already exist modules
	if t.isIncr {
		cond := condition.CreateCondition()
		cond.Field(common.BKAppIDField).Eq(t.bizID)
		cond.Field(common.BKHostIDField).Eq(hostID)
		cond.Field(common.BKModuleIDField).In(t.moduleIDArr)
		condMap := util.SetModOwner(cond.ToMapStr(), ctx.SupplierAccount)
		relationArr := make([]metadata.ModuleHost, 0)
		err := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(condMap).All(ctx, &relationArr)
		if err != nil {
			blog.ErrorJSON("add  host relation, retrieve original data error. err:%v, cond:%s, rid:%s", err, condMap, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
		}
		//  map[moduleID]bool
		existModuleIDMap := make(map[int64]bool, 0)
		for _, item := range relationArr {
			existModuleIDMap[item.ModuleID] = true
		}
		for _, moduleID := range t.moduleIDArr {
			if _, ok := existModuleIDMap[moduleID]; !ok {
				moduleIDArr = append(moduleIDArr, moduleID)
			}
		}

	}
	if len(moduleIDArr) == 0 {
		return nil
	}
	var insertDataArr []mapstr.MapStr
	for _, moduleID := range moduleIDArr {
		insertData := mapstr.MapStr{
			common.BKAppIDField: bizID, common.BKHostIDField: hostID, common.BKModuleIDField: moduleID,
		}
		insertDataArr = append(insertDataArr, insertData)
	}

	err := t.mh.dbProxy.Table(common.BKTableNameModuleHostConfig).Insert(ctx, insertDataArr)
	if err != nil {
		blog.Errorf("add host module relation, add module host relation error: %v", err)
		return err
	}
	t.curDatas = insertDataArr
	return nil
}

// getDefaultModuleIDArr get default module
func (t *transferHostModule) getDefaultModuleIDArr(ctx core.ContextParams) errors.CCErrorCoder {
	bizID := t.bizID
	// transfer the host across businees,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizID
	}
	moduleConds := condition.CreateCondition()
	moduleConds.Field(common.BKAppIDField).Eq(bizID)
	moduleConds.Field(common.BKDefaultField).NotEq(0)
	cond := util.SetModOwner(moduleConds.ToMapStr(), ctx.SupplierAccount)

	moduleInfoArr := make([]mapstr.MapStr, 0)
	err := t.mh.dbProxy.Table(common.BKTableNameBaseModule).Find(cond).All(ctx, &moduleInfoArr)
	if err != nil {
		blog.ErrorJSON("getDefaultModuleIDArr find data error. err:%s,cond:%s, rid:%s", err.Error(), cond, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}
	if len(moduleInfoArr) == 0 {
		blog.Warnf("getDefaultModuleIDArr not found default module. appID:%d, rid:%s", bizID, ctx.ReqID)
	}
	for _, moduleInfo := range moduleInfoArr {
		moduleID, err := moduleInfo.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("getDefaultModuleIDArr module info field module id not integer. err:%s, moduleInfo:%s,rid:%s", err.Error(), moduleInfo, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		t.defaultModuleID = append(t.defaultModuleID, moduleID)
	}

	return nil
}

func (t *transferHostModule) GetDefaultModuleIDArr(ctx core.ContextParams) ([]int64, errors.CCError) {
	if len(t.defaultModuleID) == 0 {
		err := t.getDefaultModuleIDArr(ctx)
		return nil, err
	}
	return t.defaultModuleID, nil
}
