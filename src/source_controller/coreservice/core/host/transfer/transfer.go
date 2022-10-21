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

// Package transfer TODO
package transfer

import (
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/thirdparty/hooks"
)

type genericTransfer struct {
	dependent           OperationDependence
	hostApplyDependence HostApplyRuleDependence

	// depend parameter
	moduleIDArr []int64
	bizID       int64
	// Incr=true is added to the module
	// Incr=false delete exist module, add current module
	isIncrement bool
	// needAutoCreateSvcInst whether service instances should be created when target modules have template
	needAutoCreateSvcInst bool
	// cross-business transfer module
	// From the A business to the B business module
	crossBizTransfer bool
	// cross-business transfer module, source business id
	srcBizIDs []int64

	// ***** cache ********
	// inner module id array
	innerModuleID []int64
	// map[bk_module_id]bk_set_id
	moduleIDSetIDMap map[int64]int64
}

// ValidParameter valid if transfer host parameter legal
func (t *genericTransfer) ValidParameter(kit *rest.Kit) errors.CCErrorCoder {
	if err := hooks.ValidHostTransferHook(kit, mongodb.Client(), t.crossBizTransfer, t.srcBizIDs, t.bizID); err != nil {
		return err
	}

	err := t.validParameterInst(kit)
	if err != nil {
		return err
	}

	err = t.validParameterModule(kit)
	if err != nil {
		return err
	}

	archived, err := t.isAppArchived(kit)
	if err != nil {
		blog.Errorf("check if app archived failed, err: %v, rid:%s", err, kit.Rid)
		return err
	}

	if archived {
		blog.Errorf("target business or source business has been archived, dest biz ID: %d, src bizI Ds: %v, rid: %s",
			t.bizID, t.srcBizIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTransferHostToArchivedApp)
	}

	if len(t.innerModuleID) == 0 {
		err := t.getInnerModuleIDArr(kit)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetCrossBusiness Set host cross-service transfer parameters
func (t *genericTransfer) SetCrossBusiness(kit *rest.Kit, bizIDs []int64) {
	t.crossBizTransfer = true
	t.srcBizIDs = bizIDs
}

// Transfer TODO
func (t *genericTransfer) Transfer(kit *rest.Kit, hostIDs []int64, disableHostApply bool) errors.CCErrorCoder {
	err := t.validHosts(kit, hostIDs)
	if err != nil {
		return err
	}

	// remove service instance if necessary
	if err := t.removeHostServiceInstance(kit, hostIDs); err != nil {
		return err
	}

	if err := t.delHostModuleRelation(kit, hostIDs); err != nil {
		// It is not the time to merge and base the time. When it fails,
		// it is clear that the data before the change is pushed.
		return err
	}

	// transfer host module config
	if err := t.addHostModuleRelationAndHostApply(kit, hostIDs, disableHostApply); err != nil {
		return err
	}

	// auto create service instance if necessary
	if t.needAutoCreateSvcInst {
		if err := t.autoCreateServiceInstance(kit, hostIDs); err != nil {
			return err
		}
	}

	return nil
}

// DeleteHosts TODO
func (t *genericTransfer) DeleteHosts(kit *rest.Kit, hostIDs []int64) error {
	if len(hostIDs) == 0 {
		return nil
	}

	// check if biz exist
	if err := t.validParameterInst(kit); err != nil {
		return err
	}

	// check if hosts belong to biz
	if err := t.validHostsBelongBiz(kit, hostIDs); err != nil {
		return err
	}

	// remove service instances
	if err := t.removeHostServiceInstance(kit, hostIDs); err != nil {
		return err
	}

	// remove host module relations
	if err := t.delHostModuleRelation(kit, hostIDs); err != nil {
		return err
	}

	// remove hosts
	hostCond := map[string]interface{}{common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs}}
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Delete(kit.Ctx, hostCond); err != nil {
		blog.Errorf("delete host failed, err: %s, host ID: %+v, rid: %s", err.Error(), hostIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

// validParameterInst  validate module, biz, srcBiz must be exist
func (t *genericTransfer) validParameterInst(kit *rest.Kit) errors.CCErrorCoder {
	var bizIDs []int64
	if t.crossBizTransfer {
		bizIDs = t.srcBizIDs
		hasSameSrcDestBiz := false
		for _, bizID := range bizIDs {
			if bizID == t.bizID {
				hasSameSrcDestBiz = true
				break
			}
		}
		if !hasSameSrcDestBiz {
			bizIDs = append(bizIDs, t.bizID)
		}
	} else {
		bizIDs = []int64{t.bizID}
	}

	appCond := map[string]interface{}{common.BKAppIDField: mapstr.MapStr{common.BKDBIN: bizIDs}}
	appCond = util.SetQueryOwner(appCond, kit.SupplierAccount)

	cnt, err := t.countByCond(kit, appCond, common.BKTableNameBaseApp)
	if err != nil {
		return err
	}
	if int(cnt) != len(bizIDs) {
		blog.Errorf("host transfer business are not all exists. cond: %v, rid: %s", appCond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceBusinessNotExist, bizIDs)
	}
	return nil
}

// validParameterModule validate parameter module legal
// module must be exist in business
// multiple modules not default module, transfer default must be one module
func (t *genericTransfer) validParameterModule(kit *rest.Kit) errors.CCErrorCoder {
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

// validHosts TODO
// validParameterHostBelongbiz  legal
// check if the host belongs to the transfer business.
// check host exist
func (t *genericTransfer) validHosts(kit *rest.Kit, hostIDs []int64) errors.CCErrorCoder {
	if len(hostIDs) == 0 {
		return nil
	}

	// check if hosts exist and do not contain cloud host that doesn't allow cross biz transfer
	hostCond := map[string]interface{}{common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs}}
	if t.crossBizTransfer {
		hostCond[common.BKCloudHostIdentifierField] = mapstr.MapStr{common.BKDBNE: true}
	}
	hostCond = util.SetQueryOwner(hostCond, kit.SupplierAccount)

	cnt, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(&hostCond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("valid hosts, but count host failed, err: %s, hostIDs: %+v, rid: %s", err.Error(), hostIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if int(cnt) < len(hostIDs) {
		blog.Errorf("valid hosts, but some hosts not exist, hostIDs: %+v, rid: %s", hostIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
	}

	return t.validHostsBelongBiz(kit, hostIDs)
}

// validHostsBelongBiz check if hosts not belongs to other biz
func (t *genericTransfer) validHostsBelongBiz(kit *rest.Kit, hostIDs []int64) errors.CCErrorCoder {
	if len(hostIDs) == 0 {
		return nil
	}

	bizIDs := []int64{t.bizID}
	// transfer the host across business, check host belongs to the original business IDs.
	if t.crossBizTransfer {
		bizIDs = t.srcBizIDs
	}

	relationCond := map[string]interface{}{common.BKAppIDField: map[string]interface{}{common.BKDBNIN: bizIDs},
		common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs}}
	relationCond = util.SetQueryOwner(relationCond, kit.SupplierAccount)

	cnt, err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(relationCond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("valid host, but get host relation failed, err: %v, biz IDs: %v, host ID: %v, rid: %s", err, bizIDs,
			hostIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if cnt > 0 {
		blog.Errorf("delete host, but some hosts belongs to other biz, biz IDs: %v, host ID: %v, rid: %s", bizIDs,
			hostIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotBelongBusiness, hostIDs, bizIDs)
	}

	return nil
}

// delHostModuleRelation delete single host module relation
func (t *genericTransfer) delHostModuleRelation(kit *rest.Kit, hostIDs []int64) errors.CCErrorCoder {
	if len(hostIDs) == 0 {
		return nil
	}

	bizID := []int64{t.bizID}
	// transfer the host across business,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		bizID = t.srcBizIDs
	}

	if t.isIncrement {
		// delete default module
		return t.delHostModuleRelationItem(kit, bizID, hostIDs, true)

	} else {
		// delete all module
		return t.delHostModuleRelationItem(kit, bizID, hostIDs, false)
	}
}

// delHostModuleRelationItem 删除主机与模块关系，注意:此时传入的bizID 和hostIDs的对应关系可能不准确，只能保证hostIDS 在bizIDs中。
// 后续可以考虑优化将bizID和hostIDs封装成一个结构保证其对应关系。
func (t *genericTransfer) delHostModuleRelationItem(kit *rest.Kit, bizIDs []int64, hostIDs []int64,
	isDefault bool) errors.CCErrorCoder {

	relationCond := map[string]interface{}{
		common.BKAppIDField:  map[string]interface{}{common.BKDBIN: bizIDs},
		common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs},
	}
	if isDefault {
		relationCond[common.BKModuleIDField] = map[string]interface{}{common.BKDBIN: t.innerModuleID}
	} else {
		if len(t.moduleIDArr) > 0 {
			relationCond[common.BKModuleIDField] = map[string]interface{}{
				common.BKDBNIN: t.moduleIDArr,
			}
		}
	}

	err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Delete(kit.Ctx, relationCond)
	if err != nil {
		blog.Errorf("delete host, but remove host relations failed, biz IDs: %v, host ID: %v, err: %v, rid: %s", bizIDs,
			hostIDs, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

// addHostModuleRelationAndHostApply add single host module relation and implement host properties apply
func (t *genericTransfer) addHostModuleRelationAndHostApply(kit *rest.Kit, hostIDs []int64,
	disable bool) errors.CCErrorCoder {
	if len(hostIDs) == 0 || len(t.moduleIDArr) == 0 {
		return nil
	}
	existHostModuleIDMap := make(map[int64]map[int64]struct{}, 0)

	// filter already exist modules
	cond := map[string]interface{}{
		common.BKAppIDField:    t.bizID,
		common.BKHostIDField:   map[string]interface{}{common.BKDBIN: hostIDs},
		common.BKModuleIDField: map[string]interface{}{common.BKDBIN: t.moduleIDArr},
	}
	condMap := util.SetQueryOwner(cond, kit.SupplierAccount)

	relationArr := make([]metadata.ModuleHost, 0)
	err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(condMap).Fields(common.BKHostIDField,
		common.BKModuleIDField).All(kit.Ctx, &relationArr)
	if err != nil {
		blog.ErrorJSON("add host relation, retrieve original data error. err:%v, cond:%s, rid:%s", err, condMap, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	for _, item := range relationArr {
		if _, ok := existHostModuleIDMap[item.HostID]; !ok {
			existHostModuleIDMap[item.HostID] = make(map[int64]struct{}, 0)
		}
		existHostModuleIDMap[item.HostID][item.ModuleID] = struct{}{}
	}

	var insertDataArr []metadata.ModuleHost
	for _, moduleID := range t.moduleIDArr {
		for _, hostID := range hostIDs {
			if _, ok := existHostModuleIDMap[hostID][moduleID]; ok {
				continue
			}
			insertDataArr = append(insertDataArr, metadata.ModuleHost{
				SetID:    t.moduleIDSetIDMap[moduleID],
				ModuleID: moduleID,
				HostID:   hostID,
				AppID:    t.bizID,
				OwnerID:  kit.SupplierAccount,
			})
		}
	}

	if len(insertDataArr) == 0 {
		return nil
	}

	if err = mongodb.Client().Table(common.BKTableNameModuleHostConfig).Insert(kit.Ctx, insertDataArr); err != nil {
		blog.Errorf("add host module relation, add module host relation error: %v, rid: %s", err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	// 对于主机转移场景下此标记如果是true表示不进行主机属性自动应用，其余场景按照原逻辑进行主机级别的属性自动应用
	if disable {
		return nil
	}
	// run new module's host apply rule, prevent repeated execution of rules
	if _, err := t.hostApplyDependence.RunHostApplyOnHosts(kit, t.bizID, insertDataArr); err != nil {
		blog.Errorf("run host apply rule failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (t *genericTransfer) autoCreateServiceInstance(kit *rest.Kit, hostIDs []int64) errors.CCErrorCoder {
	if err := t.dependent.AutoCreateServiceInstanceModuleHost(kit, hostIDs, t.moduleIDArr); err != nil {
		blog.Warnf("autoCreateServiceInstance failed, hostIDs: %+v, moduleID: %+v, rid: %s", hostIDs, t.moduleIDArr, kit.Rid)
	}
	return nil
}

// removeHostServiceInstance TODO
// remove service instances bound to hosts with process instances in certain modules
func (t *genericTransfer) removeHostServiceInstance(kit *rest.Kit, hostIDs []int64) errors.CCErrorCoder {
	// increment transfer don't need to remove service instance
	if t.isIncrement {
		return nil
	}

	if len(hostIDs) == 0 {
		return nil
	}

	// get all service instance IDs that need to be removed
	serviceInstanceFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs},
	}
	if len(t.moduleIDArr) > 0 {
		serviceInstanceFilter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBNIN: t.moduleIDArr,
		}
	}
	instances := make([]metadata.ServiceInstance, 0)
	err := mongodb.Client().Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).All(kit.Ctx, &instances)
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

	audit := auditlog.NewSvcInstAuditGenerator()
	audit.WithServiceInstance(instances)

	// get all process IDs of the service instances to be removed
	processRelationFilter := map[string]interface{}{
		common.BKServiceInstanceIDField: map[string]interface{}{
			common.BKDBIN: serviceInstanceIDs,
		},
	}
	relations := make([]metadata.ProcessInstanceRelation, 0)
	if err := mongodb.Client().Table(common.BKTableNameProcessInstanceRelation).Find(processRelationFilter).All(kit.Ctx, &relations); nil != err {
		blog.Errorf("removeHostServiceInstance failed, get process instance relation failed, err: %s, filter: %s, rid: %s", err, processRelationFilter, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	processIDs := make([]int64, 0)
	for _, relation := range relations {
		processIDs = append(processIDs, relation.ProcessID)
	}

	// delete all process relations and instances
	if len(processIDs) > 0 {
		if err := mongodb.Client().Table(common.BKTableNameProcessInstanceRelation).Delete(kit.Ctx, processRelationFilter); nil != err {
			blog.Errorf("removeHostServiceInstance failed, delete process instance relation failed, err: %s, filter: %s, rid: %s", err, processRelationFilter, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
		}

		processFilter := map[string]interface{}{
			common.BKProcessIDField: map[string]interface{}{
				common.BKDBIN: processIDs,
			},
		}

		// generate process audit log for host transfer, use a new kit to keep it out of the txn
		auditKit := kit.NewKit()
		processes := make([]mapstr.MapStr, 0)
		err := mongodb.Client().Table(common.BKTableNameBaseProcess).Find(processFilter).All(auditKit.Ctx, &processes)
		if err != nil {
			blog.Errorf("get process instances failed, filter: %+v, err: %v, rid: %s", processFilter, err, auditKit.Rid)
			return auditKit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}

		procMap := make(map[int64]mapstr.MapStr)
		for _, proc := range processes {
			procID, err := util.GetInt64ByInterface(proc[common.BKProcessIDField])
			if err != nil {
				blog.ErrorJSON("get process(%+v) id failed, err: %v, rid: %s", proc, err, auditKit.Rid)
				return auditKit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessIDField)
			}
			procMap[procID] = proc
		}

		genAuditParam := auditlog.NewGenerateAuditCommonParameter(auditKit, metadata.AuditDelete)
		if err := audit.WithProc(genAuditParam, processes, relations); err != nil {
			return err
		}

		// delete processes
		if err := mongodb.Client().Table(common.BKTableNameBaseProcess).Delete(kit.Ctx, processFilter); nil != err {
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
	if err := mongodb.Client().Table(common.BKTableNameServiceInstance).Delete(kit.Ctx, serviceInstanceIDFilter); nil != err {
		blog.Errorf("removeHostServiceInstance failed, delete service instances failed, err: %s, filter: %s, rid: %s", err, serviceInstanceIDFilter, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	// generate and save service instance audit logs for host transfer, use a new kit to keep it out of the txn
	auditKit := kit.NewKit()
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(auditKit, metadata.AuditDelete)
	auditLogs := audit.GenerateAuditLog(generateAuditParameter)
	if err := t.dependent.CreateAuditLogDependence(auditKit, auditLogs...); err != nil {
		return kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	return nil
}

// getInnerModuleIDArr get inner module ids for source biz
func (t *genericTransfer) getInnerModuleIDArr(kit *rest.Kit) errors.CCErrorCoder {
	moduleCond := mapstr.MapStr{common.BKDefaultField: mapstr.MapStr{common.BKDBNE: common.DefaultFlagDefaultValue}}

	// transfer the host across business,
	// check host belongs to the original business ID
	if t.crossBizTransfer {
		moduleCond[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: t.srcBizIDs}
	} else {
		moduleCond[common.BKAppIDField] = t.bizID
	}
	cond := util.SetQueryOwner(moduleCond, kit.SupplierAccount)

	moduleIDArr, err := mongodb.Client().Table(common.BKTableNameBaseModule).Distinct(kit.Ctx, common.BKModuleIDField,
		cond)
	if err != nil {
		blog.Errorf("get inner module ids failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if len(moduleIDArr) == 0 {
		blog.Warnf("not found inner module. cond: %v, rid: %s", cond, kit.Rid)
		return nil
	}

	t.innerModuleID, err = util.SliceInterfaceToInt64(moduleIDArr)
	if err != nil {
		blog.Errorf("parse inner module ids failed, err: %v, raw ids: %v, rid: %s", err, moduleIDArr, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParseDBFailed)
	}

	return nil
}

// GetInnerModuleIDArr TODO
func (t *genericTransfer) GetInnerModuleIDArr(kit *rest.Kit) ([]int64, errors.CCError) {
	if len(t.innerModuleID) == 0 {
		err := t.getInnerModuleIDArr(kit)
		return t.innerModuleID, err
	}
	return t.innerModuleID, nil
}

// HasInnerModule TODO
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
	err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(cond).Fields(fields...).All(kit.Ctx, &moduleInfoArr)
	if err != nil {
		blog.ErrorJSON("getModuleInfoByModuleID find data CCErrorCoder. err:%s,cond:%s, rid:%s", err.Error(), cond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return moduleInfoArr, nil
}

func (t *genericTransfer) countByCond(kit *rest.Kit, conds mapstr.MapStr, tableName string) (uint64, errors.CCErrorCoder) {
	conds = util.SetQueryOwner(conds, kit.SupplierAccount)
	cnt, err := mongodb.Client().Table(tableName).Find(conds).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("countByCond find data error. err:%s, table:%s,cond:%s, rid:%s", err.Error(), tableName, conds, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return cnt, nil
}

func (t *genericTransfer) isAppArchived(kit *rest.Kit) (bool, errors.CCErrorCoder) {

	bizIDs := append(t.srcBizIDs, t.bizID)

	for _, bizId := range bizIDs {
		if bizId <= 0 {
			return false, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
		}
	}

	cond := mapstr.MapStr{
		common.BKAppIDField:      mapstr.MapStr{common.BKDBIN: bizIDs},
		common.BKDataStatusField: common.DataStatusDisabled,
	}

	cnt, err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count diabled biz failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return false, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if cnt != 0 {
		return true, nil
	}

	return false, nil
}
