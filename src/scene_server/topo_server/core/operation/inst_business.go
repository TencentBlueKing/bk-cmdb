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

package operation

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// BusinessOperationInterface business operation methods
type BusinessOperationInterface interface {
	CreateBusiness(kit *rest.Kit, obj model.Object, data mapstr.MapStr) (inst.Inst, error)

	// UpdateBusinessIdleSetOrModule 此函数用于更新全局的空闲机池和及下面的模块，属于管理员操作，只允许将此接口提供给前端使用
	UpdateBusinessIdleSetOrModule(kit *rest.Kit, option *metadata.ConfigUpdateSettingOption) error

	// DeleteBusinessGlobalUserModule 删除用户自定义的空闲类模块，此接口只允许提供给前端，用于管理员使用
	DeleteBusinessGlobalUserModule(kit *rest.Kit, obj model.Object, option *metadata.BuiltInModuleDeleteOption) error
	FindBiz(kit *rest.Kit, cond *metadata.QueryCondition) (count int, results []mapstr.MapStr, err error)
	GetInternalModule(kit *rest.Kit, bizID int64) (count int, result *metadata.InnterAppTopo, err errors.CCErrorCoder)
	UpdateBusiness(kit *rest.Kit, data mapstr.MapStr, obj model.Object, bizID int64) error
	UpdateBusinessByCond(kit *rest.Kit, data mapstr.MapStr, obj model.Object, cond mapstr.MapStr) error
	DeleteBusiness(kit *rest.Kit, bizIDs []int64) error
	HasHosts(kit *rest.Kit, bizID int64) (bool, error)
	SetProxy(set SetOperationInterface, module ModuleOperationInterface, inst InstOperationInterface, obj ObjectOperationInterface)
	GenerateAchieveBusinessName(kit *rest.Kit, bizName string) (achieveName string, err error)
}

// NewBusinessOperation create a business instance
func NewBusinessOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) BusinessOperationInterface {
	return &business{
		clientSet:   client,
		authManager: authManager,
	}
}

type business struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	inst        InstOperationInterface
	set         SetOperationInterface
	module      ModuleOperationInterface
	obj         ObjectOperationInterface
}

func (b *business) SetProxy(set SetOperationInterface, module ModuleOperationInterface, inst InstOperationInterface, obj ObjectOperationInterface) {
	b.inst = inst
	b.set = set
	b.module = module
	b.obj = obj
}

func (b *business) HasHosts(kit *rest.Kit, bizID int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		Fields:        []string{common.BKHostIDField},
		Page:          metadata.BasePage{Limit: 1},
	}
	rsp, err := b.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
	if nil != err {
		blog.Errorf("[operation-set] failed to request the object controller, error info is %s", err.Error())
		return false, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-set]  failed to search the host set configures, error info is %s", rsp.ErrMsg)
		return false, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}

// updateIdleModuleConfig update admin module config.
func (b *business) updateIdleModuleConfig(kit *rest.Kit, option metadata.ModuleOption,
	config metadata.PlatformSettingConfig) error {

	// 获取db中platform的配置
	//res, err := b.clientSet.CoreService().System().SearchPlatformSetting(kit.Ctx, kit.Header)
	//if err != nil {
	//	return err
	//}

	conf, oldConf := config, config
	flag := false
	updateFields := mapstr.New()

	switch option.Key {
	case common.SystemIdleModuleKey:
		conf.BuiltInModuleConfig.IdleName = option.Name
		updateFields.Set(common.SystemIdleModuleKey, option.Name)
		flag = true

	case common.SystemFaultModuleKey:
		conf.BuiltInModuleConfig.FaultName = option.Name
		updateFields.Set(common.SystemFaultModuleKey, option.Name)
		flag = true

	case common.SystemRecycleModuleKey:
		conf.BuiltInModuleConfig.RecycleName = option.Name
		updateFields.Set(common.SystemRecycleModuleKey, option.Name)
		flag = true

	default:
		for index, module := range conf.BuiltInModuleConfig.UserModules {
			if module.Key == option.Key {
				conf.BuiltInModuleConfig.UserModules[index].Value = option.Name
				updateFields.Set("module_name", option.Name)
				updateFields.Set("module_key", option.Key)
				flag = true
				break
			}
		}
	}
	// flag: false 用户新增模块场景
	if !flag {
		conf.BuiltInModuleConfig.UserModules = append(conf.BuiltInModuleConfig.UserModules, metadata.UserModuleList{
			Key:   option.Key,
			Value: option.Name,
		})
		updateFields.Set(common.UserDefinedModules, metadata.UserModuleList{
			Key:   option.Key,
			Value: option.Name,
		})
	}

	_, err := b.clientSet.CoreService().System().UpdatePlatformSetting(context.Background(), kit.Header, &conf)
	if err != nil {
		return err
	}

	err = b.savePlatformLog(kit, metadata.AuditUpdate, &oldConf, updateFields)
	if err != nil {
		blog.Errorf("generate audit log failed config: %v, updateFields: %v, err: %v, rid: %s", oldConf, updateFields,
			err, kit.Rid)
		return err
	}
	return nil
}

// savePlatformLog 平台管理的审计日志，无论是新建模块，修改模块或者集群名字还是删除用户自定义场景都是对空闲机池的修改，统一走update
func (b *business) savePlatformLog(kit *rest.Kit, auditLog metadata.ActionType, oldConf *metadata.PlatformSettingConfig,
	updateFields mapstr.MapStr) error {

	audit := auditlog.NewPlatFormSettingAuditLog(b.clientSet.CoreService())
	auditParam := auditlog.NewGenerateAuditCommonParameter(kit, auditLog).WithUpdateFields(updateFields)

	auditLogs, err := audit.GenerateAuditLog(auditParam, oldConf)

	if err != nil {
		return err
	}
	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		return err
	}
	return nil
}

func (b *business) updateBuiltInSetConfig(kit *rest.Kit, setName string) error {

	//header := util.BuildHeader(common.CCSystemOperatorUserName, common.BKDefaultOwnerID)

	res, err := b.clientSet.CoreService().System().SearchPlatformSetting(context.Background(), kit.Header)
	if err != nil {
		return err
	}

	conf, oldConf := res.Data, res.Data
	conf.BuiltInSetName = metadata.ObjectString(setName)

	_, err = b.clientSet.CoreService().System().UpdatePlatformSetting(context.Background(), kit.Header, &conf)
	if err != nil {
		return err
	}
	updateFields := mapstr.New()
	updateFields.Set(common.SystemSetName, setName)

	err = b.savePlatformLog(kit, metadata.AuditUpdate, &oldConf, updateFields)
	if err != nil {
		blog.Errorf("generate audit log failed, config: %v,updateFields: %v err: %v, rid: %s", oldConf,
			updateFields, err, kit.Rid)
	}
	return nil
}

// deleteIdlePoolConfig 更新删除用户自定义场景下的配置
func (b *business) deleteUserModuleConfig(kit *rest.Kit, option *metadata.BuiltInModuleDeleteOption,
	config metadata.PlatformSettingConfig) error {

	conf, oldConf := config, config
	updateFields := mapstr.New()

	for index, module := range conf.BuiltInModuleConfig.UserModules {
		if module.Key == option.ModuleKey {
			conf.BuiltInModuleConfig.UserModules = append(conf.BuiltInModuleConfig.UserModules[:index],
				conf.BuiltInModuleConfig.UserModules[index+1:]...)
			updateFields.Set(common.UserDefinedModules, metadata.UserModuleList{
				Key:   option.ModuleKey,
				Value: option.ModuleName,
			})
			break
		}
	}

	_, err := b.clientSet.CoreService().System().UpdatePlatformSetting(kit.Ctx, kit.Header, &conf)
	if err != nil {
		return err
	}

	err = b.savePlatformLog(kit, metadata.AuditUpdate, &oldConf, updateFields)
	if err != nil {
		blog.Errorf("generate audit log failed config :v,updateFields: %v err: %v, rid: %s", oldConf, updateFields,
			err, kit.Rid)
		return err
	}

	return nil
}

// checkModuleNameValid check whether the module has duplicate names.
func (b *business) checkModuleNameValid(ctx *rest.Kit, input metadata.ModuleOption) error {
	obj, err := b.obj.FindSingleObject(ctx, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), ctx.Rid)
		return err
	}

	queryCond := &metadata.QueryInput{
		Condition: map[string]interface{}{
			common.BKModuleNameField: input.Name,
		},
	}

	cnt, _, err := b.module.FindModule(ctx, obj, queryCond)
	if err != nil {
		blog.Errorf("find module fail, queryCond: %+v,error %v, rid: %s", queryCond, err, ctx.Rid)
		return err
	}
	if cnt > 0 {
		return fmt.Errorf("update module name fail, duplicate module name")
	}
	return nil
}

// Validate 判断参数是否合法，注意此处合法的标准如下:
// 1、key如果存在这个key只要名字不重复即合法
// 2、如果不存在这个key那么合法即是新增场景
// 3、true是更新场景，false是新增场景
// 4、目前系统出厂的idle 、fault、recycle只支持改名字不支持删除
func (b *business) validateIdleModuleConfigName(ctx *rest.Kit, input metadata.ModuleOption) (error, bool, string,
	metadata.PlatformSettingConfig) {

	res, err := b.clientSet.CoreService().System().SearchPlatformSetting(ctx.Ctx, ctx.Header)
	if err != nil {
		return err, false, "", metadata.PlatformSettingConfig{}
	}

	conf := res.Data
	flag := false
	var oldName string
	switch input.Key {
	case common.SystemIdleModuleKey:
		if input.Name == conf.BuiltInModuleConfig.IdleName {
			return fmt.Errorf("idle name cannot be the same"), false, "", metadata.PlatformSettingConfig{}
		}
		oldName = conf.BuiltInModuleConfig.IdleName
		flag = true

	case common.SystemFaultModuleKey:
		if input.Name == conf.BuiltInModuleConfig.FaultName {
			return fmt.Errorf("fault name cannot be the same"), false, "", metadata.PlatformSettingConfig{}
		}
		oldName = conf.BuiltInModuleConfig.FaultName
		flag = true

	case common.SystemRecycleModuleKey:
		if input.Name == conf.BuiltInModuleConfig.RecycleName {
			return fmt.Errorf("recycle name cannot be the same"), false, "", metadata.PlatformSettingConfig{}
		}
		oldName = conf.BuiltInModuleConfig.RecycleName
		flag = true

	default:
		for _, m := range conf.BuiltInModuleConfig.UserModules {
			if m.Key == input.Key {
				if m.Value == input.Name {
					return fmt.Errorf("user defined module name cannot be the same"), false, "",
						metadata.PlatformSettingConfig{}
				} else {
					oldName = m.Value
					flag = true
					break
				}
			}
		}
	}
	return nil, flag, oldName, conf
}

func (b *business) validateDeleteModuleName(ctx context.Context, option *metadata.BuiltInModuleDeleteOption) (metadata.PlatformSettingConfig, error) {

	header := util.BuildHeader(common.CCSystemOperatorUserName, common.BKDefaultOwnerID)

	res, err := b.clientSet.CoreService().System().SearchPlatformSetting(context.Background(), header)
	if err != nil {
		return metadata.PlatformSettingConfig{}, err
	}
	conf := res.Data

	for _, userModule := range conf.BuiltInModuleConfig.UserModules {
		if userModule.Key == option.ModuleKey {
			return conf, nil
		}
	}
	return metadata.PlatformSettingConfig{}, fmt.Errorf("no key founded")
}

// addUserDefinedModule 增加用户的自定义模块操作
func (b *business) addUserDefinedModule(kit *rest.Kit, results []inst.Inst, data metadata.ModuleOption) error {

	// create module
	objModule, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("add user define module fail, failed to search the set, %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	ds := make([]mapstr.MapStr, 0)

	defaultCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("failed to search default category, err: %+v, rid: %s", err, kit.Rid)
		return err
	}

	for _, result := range results {

		d := result.GetValues()

		setId, err := util.GetInt64ByInterface(d[common.BKSetIDField])
		if err != nil {
			blog.Errorf("add user define module fail, decode set failed, data: %v, err: %s, rid: %s", d, err,
				kit.Rid)
			return err
		}

		bizId, err := result.GetBizID()
		if nil != err {
			blog.Errorf("add user define module fail, error is %v, rid: %s", err, kit.Rid)
			return err
		}

		userModuleData := mapstr.New()
		userModuleData.Set(common.BKSetIDField, setId)
		userModuleData.Set(common.BKInstParentStr, setId)
		userModuleData.Set(common.BKAppIDField, bizId)

		userModuleData.Set(common.BKModuleNameField, data.Name)

		// 设置用户新增的模块标识
		userModuleData.Set(common.BKDefaultField, common.DefaultUserResModuleFlag)
		userModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
		userModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
		userModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)
		ds = append(ds, userModuleData)
	}

	d := &metadata.CreateManyModelInstance{Datas: ds}

	objID := objModule.GetObjectID()
	rsp, err := b.clientSet.CoreService().Instance().CreateManyInstance(kit.Ctx, kit.Header, objID, d)
	if nil != err {
		blog.Errorf("failed to create object instance, error info is %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	if !rsp.Result {
		blog.Errorf("failed to create object instance ,error info is %v, rid: %s", rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (b *business) updateModuleName(kit *rest.Kit, data metadata.ModuleOption, oldName string) error {

	obj, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	var defaultFlag int

	// 根据不同类型的模块获取不同的模块标记,最终获取到原来的模块名称
	switch data.Key {
	case common.SystemIdleModuleKey:
		defaultFlag = common.DefaultResModuleFlag
	case common.SystemFaultModuleKey:
		defaultFlag = common.DefaultFaultModuleFlag
	case common.SystemRecycleModuleKey:
		defaultFlag = common.DefaultRecycleModuleFlag
	default:
		defaultFlag = common.DefaultUserResModuleFlag
	}

	d := mapstr.New()
	d.Set(common.BKModuleNameField, data.Name)

	inputParams := metadata.UpdateOption{
		Data: d,
		Condition: mapstr.MapStr{
			common.BKDefaultField:    defaultFlag,
			common.BKModuleNameField: oldName},
	}

	rsp, err := b.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, obj.GetObjectID(),
		&inputParams)
	if nil != err {
		blog.Errorf("update inst failed to request object controller, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("update inst failed to set the object(%s) inst by the condition(%#v), err: %s, rid: %s",
			obj.Object().ObjectID, inputParams, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (b *business) deleteModuleName(kit *rest.Kit, op *metadata.BuiltInModuleDeleteOption) error {

	obj, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), kit.Rid)
		return err
	}
	delInstIDs := make([]int64, 0)

	queryCond := &metadata.QueryInput{
		Condition: map[string]interface{}{
			common.BKDefaultField:    common.DefaultUserResModuleFlag,
			common.BKModuleNameField: op.ModuleName,
		},
		Fields: common.BKModuleIDField,
	}

	cnt, instItems, err := b.module.FindModule(kit, obj, queryCond)
	if err != nil {
		blog.Errorf("failed to find the module, query %s, error is %v, rid: %s", queryCond, err, kit.Rid)
		return err
	}

	if cnt == 0 {
		blog.Errorf("no module founded, query %s, rid: %s", queryCond, kit.Rid)
		return fmt.Errorf("no module founded")
	}
	for _, instItem := range instItems {
		moduleId, err := util.GetInt64ByInterface(instItem[common.BKModuleIDField])
		if err != nil {
			blog.Errorf(" decode module id failed, instItems: %s, err: %s, rid: %s", instItems, err, kit.Rid)
			return err
		}
		delInstIDs = append(delInstIDs, moduleId)
	}

	delCond := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{common.BKDBIN: delInstIDs},
	}
	flag, err := b.module.HasHostInModules(kit, delInstIDs)
	if err != nil {
		blog.Errorf("check whether there is a host in the module failed, query %s, rid: %s", queryCond, kit.Rid)
		return err
	}

	if flag {
		return fmt.Errorf("there is a host in the module")
	}
	dc := &metadata.DeleteOption{Condition: delCond}
	rsp, err := b.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, obj.GetObjectID(), dc)
	if nil != err {
		blog.Errorf("delete inst failed, err: %v, cond: %s rid: %s", err, delCond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err := rsp.CCError(); err != nil {
		blog.Errorf("delete inst failed, err: %v, cond: %s rid: %s", err, delCond, kit.Rid)
		return err
	}

	return nil
}

// updateBusinessSet rename business idle set.
func (b *business) updateBusinessSet(kit *rest.Kit, setOptin metadata.SetOption) error {

	obj, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("get set object failed,err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// verify whether the cluster name is duplicate.
	if err := b.checkSetNameValid(kit, obj, setOptin); err != nil {
		return err
	}

	querySet := &metadata.QueryInput{
		Condition: map[string]interface{}{
			common.BKDefaultField: common.DefaultResSetFlag,
		},
		Limit: common.BKNoLimit,
	}
	s := mapstr.New()
	s.Set(common.BKSetNameField, setOptin.Name)
	// update idle set name
	err = b.set.UpdateSetForPlatform(kit, s, obj, querySet)
	if err != nil {
		blog.Errorf("update set failed, rid: %s", kit.Rid)
		return err
	}

	// update platform setting.
	err = b.updateBuiltInSetConfig(kit, setOptin.Name)
	if err != nil {
		blog.Errorf("update set config failed, rid: %s", kit.Rid)
		return err
	}

	return nil
}

// UpdateBusinessIdleSetOrModule 此函数用于更新全局的空闲机池和及下面的模块，属于管理员操作，只允许将此接口提供给前端使用
func (b *business) UpdateBusinessIdleSetOrModule(kit *rest.Kit, option *metadata.ConfigUpdateSettingOption) error {

	switch option.Type {
	case metadata.ConfigUpdateTypeSet:
		err := b.updateBusinessSet(kit, option.Set)
		if err != nil {
			return err
		}
	case metadata.ConfigUpdateTypeModule:
		err := b.updateBusinessModule(kit, option.Module)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("type error")
	}

	return nil
}

//updateBusinessModule: 对特殊空闲机池下模块(flag:1,2,3,5)做新增或改名操作
func (b *business) updateBusinessModule(kit *rest.Kit, module metadata.ModuleOption) error {

	// check param is legal or not
	if err := b.checkModuleNameValid(kit, module); err != nil {
		blog.Errorf("params is illegal err: %v, config: %v,rid: %s", err, module, kit.Rid)
		return err
	}

	err, flag, oldname, conf := b.validateIdleModuleConfigName(kit, module)
	if err != nil {
		blog.Errorf("params is illegal err: %v, config: %v,rid: %s", err, module, kit.Rid)
		return err
	}

	if flag {
		// add user module or rename
		err := b.updateModuleName(kit, module, oldname)
		if err != nil {
			return err
		}
	} else {
		// find idle set list
		results, err := b.getIdleSetList(kit)
		if err != nil {
			return err
		}
		err = b.addUserDefinedModule(kit, results, module)
		if err != nil {
			return err
		}
	}

	// update platform config
	err = b.updateIdleModuleConfig(kit, module, conf)
	if err != nil {
		blog.Errorf("update module config failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func (b *business) checkSetNameValid(kit *rest.Kit, obj model.Object, setOptin metadata.SetOption) error {

	querySet := &metadata.QueryInput{
		Condition: map[string]interface{}{
			common.BKSetNameField: setOptin.Name,
		},
	}

	count, _, err := b.set.FindSet(kit, obj, querySet)
	if err != nil {
		blog.Errorf("find set failed err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if count > 0 {
		blog.Errorf("update set name fail, duplicate cluster name exists,set name: %v, rid: %s",
			setOptin.Name, kit.Rid)
		return fmt.Errorf("update set name fail, duplicate set name")
	}
	return nil
}

func (b *business) getIdleSetList(kit *rest.Kit) (results []inst.Inst, err error) {

	obj, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}
	querySet := &metadata.QueryInput{
		Condition: map[string]interface{}{
			common.BKDefaultField: common.DefaultResSetFlag,
		},
		Limit:  common.BKNoLimit,
		Fields: fmt.Sprintf("%s,%s", common.BKSetIDField, common.BKAppIDField),
	}

	count, results, err := b.set.FindSet(kit, obj, querySet)
	if err != nil {
		blog.Errorf("find set failed err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if count == 0 {
		blog.Errorf("no set founded, rid: %s", kit.Rid)
		return nil, fmt.Errorf("no set founded")
	}
	return results, nil
}

// DeleteBusinessGlobalUserModule 删除用户自定义的空闲类模块，此接口只允许提供给前端，用于管理员使用
func (b *business) DeleteBusinessGlobalUserModule(kit *rest.Kit, obj model.Object,
	option *metadata.BuiltInModuleDeleteOption) error {

	// step 1: check param is legal or not
	conf, err := b.validateDeleteModuleName(kit.Ctx, option)
	if err != nil {
		blog.Errorf("delete global user module fail, params is illegal config: %v, err: %v, rid:%s", option, err,
			kit.Rid)
		return err
	}

	// step 2: delete user module.
	err = b.deleteModuleName(kit, option)
	if err != nil {
		return err
	}

	// step 3: update platform config.
	err = b.deleteUserModuleConfig(kit, option, conf)
	if err != nil {
		blog.Errorf("fail to update admin config err: %v,rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (b *business) CreateBusiness(kit *rest.Kit, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {

	defaultFieldVal, err := data.Int64(common.BKDefaultField)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is did not set the default field, %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}
	if defaultFieldVal == int64(common.DefaultAppFlag) && kit.SupplierAccount != common.BKDefaultOwnerID {
		// this is a new supplier owner and prepare to create a new business.
		asstQuery := map[string]interface{}{
			common.BKOwnerIDField: common.BKDefaultOwnerID,
		}
		defaultOwnerHeader := util.CloneHeader(kit.Header)
		defaultOwnerHeader.Set(common.BKHTTPOwnerID, common.BKDefaultOwnerID)

		asstRsp, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, defaultOwnerHeader, &metadata.QueryCondition{Condition: asstQuery})
		if nil != err {
			blog.Errorf("create business failed to get default assoc, error info is %s, rid: %s", err.Error(), kit.Rid)
			return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
		}
		if !asstRsp.Result {
			return nil, kit.CCError.Error(asstRsp.Code)
		}
		expectAssts := asstRsp.Data.Info
		blog.Infof("copy asst for %s, %+v, rid: %s", kit.SupplierAccount, expectAssts, kit.Rid)

		existAsstRsp, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: asstQuery})
		if nil != err {
			blog.Errorf("create business failed to get default assoc, error info is %s, rid: %s", err.Error(), kit.Rid)
			return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
		}
		if !existAsstRsp.Result {
			return nil, kit.CCError.Error(existAsstRsp.Code)
		}
		existAssts := existAsstRsp.Data.Info

	expectLoop:
		for _, asst := range expectAssts {
			asst.OwnerID = kit.SupplierAccount
			for _, existAsst := range existAssts {
				if existAsst.ObjectID == asst.ObjectID &&
					existAsst.AsstObjID == asst.AsstObjID &&
					existAsst.AsstKindID == asst.AsstKindID {
					continue expectLoop
				}
			}

			var createAsstRsp *metadata.CreatedOneOptionResult
			var err error
			if asst.AsstKindID == common.AssociationKindMainline {
				// bk_mainline is a inner association type that can only create in special case,
				// so we separate bk_mainline association type creation with a independent method,
				createAsstRsp, err = b.clientSet.CoreService().Association().CreateMainlineModelAssociation(kit.Ctx, kit.Header, &metadata.CreateModelAssociation{Spec: asst})
			} else {
				createAsstRsp, err = b.clientSet.CoreService().Association().CreateModelAssociation(kit.Ctx, kit.Header, &metadata.CreateModelAssociation{Spec: asst})
			}
			if nil != err {
				blog.Errorf("create business failed to copy default assoc, error info is %s, rid: %s", err.Error(), kit.Rid)
				return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
			}
			if !createAsstRsp.Result {
				return nil, kit.CCError.Error(createAsstRsp.Code)
			}

		}
	}

	bizInst, err := b.inst.CreateInst(kit, obj, data)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, err
	}

	bizID, err := bizInst.GetInstID()
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create set
	objSet, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}
	header := util.BuildHeader(common.CCSystemOperatorUserName, common.BKDefaultOwnerID)

	res, err := b.clientSet.CoreService().System().SearchPlatformSetting(context.Background(), header)
	if err != nil {
		return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	conf := res.Data

	setData := mapstr.New()
	setData.Set(common.BKAppIDField, bizID)
	setData.Set(common.BKInstParentStr, bizID)
	setData.Set(common.BKSetNameField, conf.BuiltInSetName)

	setData.Set(common.BKDefaultField, common.DefaultResSetFlag)

	setInst, err := b.set.CreateSet(kit, objSet, bizID, setData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	setID, err := setInst.GetInstID()
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create module
	objModule, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	defaultCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("failed to search default category, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrProcGetDefaultServiceCategoryFailed, err.Error())
	}

	idleModuleData := mapstr.New()
	idleModuleData.Set(common.BKSetIDField, setID)
	idleModuleData.Set(common.BKInstParentStr, setID)
	idleModuleData.Set(common.BKAppIDField, bizID)
	idleModuleData.Set(common.BKModuleNameField, conf.BuiltInModuleConfig.IdleName)

	idleModuleData.Set(common.BKDefaultField, common.DefaultResModuleFlag)
	idleModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	idleModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	idleModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(kit, objModule, bizID, setID, idleModuleData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create fault module
	faultModuleData := mapstr.New()
	faultModuleData.Set(common.BKSetIDField, setID)
	faultModuleData.Set(common.BKInstParentStr, setID)
	faultModuleData.Set(common.BKAppIDField, bizID)
	faultModuleData.Set(common.BKModuleNameField, conf.BuiltInModuleConfig.FaultName)

	faultModuleData.Set(common.BKDefaultField, common.DefaultFaultModuleFlag)
	faultModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	faultModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	faultModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(kit, objModule, bizID, setID, faultModuleData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create recycle module
	recycleModuleData := mapstr.New()
	recycleModuleData.Set(common.BKSetIDField, setID)
	recycleModuleData.Set(common.BKInstParentStr, setID)
	recycleModuleData.Set(common.BKAppIDField, bizID)
	recycleModuleData.Set(common.BKModuleNameField, conf.BuiltInModuleConfig.RecycleName)

	recycleModuleData.Set(common.BKDefaultField, common.DefaultRecycleModuleFlag)
	recycleModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	recycleModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	recycleModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(kit, objModule, bizID, setID, recycleModuleData)
	if nil != err {
		blog.Errorf("create business failed, create recycle module failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	for _, module := range conf.BuiltInModuleConfig.UserModules {

		// create user module
		userModuleData := mapstr.New()
		userModuleData.Set(common.BKSetIDField, setID)
		userModuleData.Set(common.BKInstParentStr, setID)
		userModuleData.Set(common.BKAppIDField, bizID)
		userModuleData.Set(common.BKModuleNameField, module.Value)

		userModuleData.Set(common.BKDefaultField, common.DefaultUserResModuleFlag)
		userModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
		userModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
		userModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

		_, err = b.module.CreateModule(kit, objModule, bizID, setID, userModuleData)
		if nil != err {
			blog.Errorf("create business failed, create user module failed, err: %v, rid: %s", err, kit.Rid)
			return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
		}

	}
	return bizInst, nil
}
func (b *business) FindBusiness(kit *rest.Kit, cond *metadata.QueryBusinessRequest) (count int, results []mapstr.MapStr, err error) {

	cond.Condition[common.BKDefaultField] = 0
	query := &metadata.QueryCondition{
		Fields:    cond.Fields,
		Condition: cond.Condition,
		Page:      cond.Page,
	}

	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.ErrorJSON("failed to find business by query condition: %s, err: %s, rid: %s", query, err.Error(), kit.Rid)
		return 0, nil, err
	}

	if !result.Result {
		return 0, nil, kit.CCError.Errorf(result.Code, result.ErrMsg)
	}

	return result.Data.Count, result.Data.Info, err
}

// obtain business information according to condition "cond".
func (b *business) FindBiz(kit *rest.Kit, cond *metadata.QueryCondition) (count int, results []mapstr.MapStr,
	err error) {
	if !cond.Condition.Exists(common.BKDefaultField) {
		cond.Condition[common.BKDefaultField] = 0
	}

	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.ErrorJSON("failed to find business by query condition: %s, err: %s, rid: %s", cond, err.Error(), kit.Rid)
		return 0, nil, err
	}

	if !result.Result {
		return 0, nil, kit.CCError.Errorf(result.Code, result.ErrMsg)
	}

	return result.Data.Count, result.Data.Info, err
}

var (
	NumRegex = regexp.MustCompile(`^\d+$`)
)

/*
GenerateAchieveBusinessName 生成归档后的业务名称
	- 业务归档的时候，自动重命名为"foo-archived"
	- 归档的时候，如果发现已经存在同名的"foo-archived", 自动在其后+1, 比如 "foo-archived-1", "foo-archived-2"
*/
func (b *business) GenerateAchieveBusinessName(kit *rest.Kit, bizName string) (achieveName string, err error) {
	queryBusinessRequest := &metadata.QueryBusinessRequest{
		Fields: []string{common.BKAppNameField},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKAppNameField: map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf(`^%s-archived`, regexp.QuoteMeta(bizName)),
			},
		},
	}
	count, data, err := b.FindBusiness(kit, queryBusinessRequest)
	if err != nil {
		return "", err
	}
	if count == 0 {
		return fmt.Sprintf("%s-archived", bizName), nil
	}
	existNums := make([]int64, 0)
	for _, item := range data {
		biz := metadata.BizBasicInfo{}
		if err := mapstruct.Decode2Struct(item, &biz); err != nil {
			blog.Errorf("GenerateBusinessAchieveName failed, Decode2Struct failed, biz: %+v, err: %+v, rid: %s", item, err, kit.Rid)
			return "", kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
		}
		parts := strings.Split(biz.BizName, fmt.Sprintf("%s-archived-", bizName))
		if len(parts) != 2 {
			continue
		}
		numPart := parts[1]
		if !NumRegex.MatchString(numPart) {
			continue
		}
		num, err := util.GetInt64ByInterface(numPart)
		if err != nil {
			blog.Errorf("GenerateBusinessAchieveName failed, GetInt64ByInterface failed, numPart: %s, err: %+v, rid: %s", numPart, err, kit.Rid)
			return "", kit.CCError.CCError(common.CCErrCommParseDataFailed)
		}
		existNums = append(existNums, num)
	}
	// 空数组时默认填充
	existNums = append(existNums, 0)
	maxNum := existNums[0]
	for _, num := range existNums {
		if num > maxNum {
			maxNum = num
		}
	}

	return fmt.Sprintf("%s-archived-%d", bizName, maxNum+1), nil
}

func (b *business) GetInternalModule(kit *rest.Kit,
	bizID int64) (count int, result *metadata.InnterAppTopo, err errors.CCErrorCoder) {
	// get set model
	querySet := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKDefaultField: common.DefaultResSetFlag,
		},
		Fields: []string{common.BKSetIDField, common.BKSetNameField},
	}
	querySet.Page.Limit = 1

	/* setRsp, err := b.inst.FindOriginInst(kit, common.BKInnerObjIDSet, querySet)
	if nil != err {
		return 0, nil, kit.CCError.New(common.CCErrTopoAppSearchFailed, err.Error())
	} */

	setRsp := &metadata.ResponseSetInstance{}
	// 返回数据不包含自定义字段
	if err = b.clientSet.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDSet, querySet, setRsp); err != nil {
		return 0, nil, err
	}
	if err := setRsp.CCError(); err != nil {
		blog.ErrorJSON("query set error. filter: %s, result: %s, rid: %s", querySet, setRsp, kit.Rid)
		return 0, nil, err
	}

	// search modules
	queryModule := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: 0,
			},
		},
		Fields: []string{common.BKModuleIDField, common.BKModuleNameField, common.BKDefaultField, common.HostApplyEnabledField},
	}
	queryModule.Page.Limit = common.BKNoLimit

	/*
		moduleRsp, err := b.inst.FindOriginInst(kit, common.BKInnerObjIDModule, queryModule)
		if nil != err {
			return 0, nil, kit.CCError.New(common.CCErrTopoAppSearchFailed, err.Error())
		}
	*/

	moduleResp := &metadata.ResponseModuleInstance{}
	// 返回数据不包含自定义字段
	if err = b.clientSet.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, queryModule, moduleResp); err != nil {
		return 0, nil, err
	}
	if err := moduleResp.CCError(); err != nil {
		blog.ErrorJSON("query module error. filter: %s, result: %s, rid: %s", queryModule, moduleResp, kit.Rid)
		return 0, nil, err
	}

	// construct result
	result = &metadata.InnterAppTopo{}
	for _, set := range setRsp.Data.Info {
		result.SetID = set.SetID
		result.SetName = set.SetName
		break // should be only one set
	}

	for _, module := range moduleResp.Data.Info {

		result.Module = append(result.Module, metadata.InnerModule{
			ModuleID:         module.ModuleID,
			ModuleName:       module.ModuleName,
			Default:          module.Default,
			HostApplyEnabled: module.HostApplyEnabled,
		})
	}

	return 0, result, nil
}

// UpdateBusiness update business instances by bizID
func (b *business) UpdateBusiness(kit *rest.Kit, data mapstr.MapStr, obj model.Object, bizID int64) error {
	cond := mapstr.MapStr{
		common.BKAppIDField: mapstr.MapStr{
			common.BKDBEQ: bizID,
		},
	}
	return b.inst.UpdateInst(kit, data, obj, cond)
}

// UpdateBusinessByCond update business instances by condition
func (b *business) UpdateBusinessByCond(kit *rest.Kit, data mapstr.MapStr, obj model.Object, cond mapstr.MapStr) error {
	return b.inst.UpdateInst(kit, data, obj, cond)
}
