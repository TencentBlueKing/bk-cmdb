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

package inst

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// BusinessOperationInterface business operation methods
type BusinessOperationInterface interface {
	// CreateBusiness create business
	CreateBusiness(kit *rest.Kit, data mapstr.MapStr) (mapstr.MapStr, error)
	DeleteBusiness(kit *rest.Kit, bizIDs []int64) error
	// FindBiz find biz
	FindBiz(kit *rest.Kit, cond *metadata.QueryCondition) (count int,
		results []mapstr.MapStr, err error)
	// HasHosts check if this business still has hosts.
	HasHosts(kit *rest.Kit, bizID int64) (bool, error)
	// GenerateAchieveBusinessName 生成归档后的业务名称
	// - 业务归档的时候，自动重命名为"foo-archived"
	// - 归档的时候，如果发现已经存在同名的"foo-archived", 自动在其后+1, 比如 "foo-archived-1", "foo-archived-2"
	GenerateAchieveBusinessName(kit *rest.Kit, bizName string) (achieveName string, err error)
	// GetBriefTopologyNodeRelation is used to get directly related business topology node information.
	// As is, you can find modules belongs to a set; or you can find the set a module belongs to.
	// It has rules as follows:
	// 1. if src object is biz, then the destination object can be any mainline object except biz.
	// 2. destination object can be biz. otherwise, src and destination object should be the neighbour.
	// this api only return business topology relations.
	GetBriefTopologyNodeRelation(kit *rest.Kit, opts *metadata.GetBriefBizRelationOptions) ([]*metadata.
		BriefBizRelations, error)
	// GetResourcePoolBusinessID search resource pool biz id
	GetResourcePoolBusinessID(kit *rest.Kit) (int64, error)
	// SetProxy SetProxy
	SetProxy(inst InstOperationInterface, module ModuleOperationInterface, set SetOperationInterface)

	// UpdateBusinessIdleSetOrModule 此函数用于更新全局的空闲机池和及下面的模块，属于管理员操作，只允许将此接口提供给前端使用
	UpdateBusinessIdleSetOrModule(kit *rest.Kit, option *metadata.ConfigUpdateSettingOption) error

	// DeleteBusinessGlobalUserModule 删除用户自定义的空闲类模块，此接口只允许提供给前端，用于管理员使用
	DeleteBusinessGlobalUserModule(kit *rest.Kit, option *metadata.BuiltInModuleDeleteOption) error
}

// NewBusinessOperation create a business instance
func NewBusinessOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) BusinessOperationInterface {
	return &business{
		clientSet:   client,
		authManager: authManager,
	}
}

type business struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	inst        InstOperationInterface
	module      ModuleOperationInterface
	set         SetOperationInterface
}

// SetProxy SetProxy
func (b *business) SetProxy(inst InstOperationInterface, module ModuleOperationInterface, set SetOperationInterface) {
	b.inst = inst
	b.module = module
	b.set = set
}

// CreateBusiness create business
func (b *business) CreateBusiness(kit *rest.Kit, data mapstr.MapStr) (mapstr.MapStr, error) {

	bizInst, err := b.inst.CreateInst(kit, common.BKInnerObjIDApp, data)
	if err != nil {
		blog.Errorf("create business failed, err: %v, data: %#v, rid: %s", err, data, kit.Rid)
		return nil, err
	}
	bizID, err := bizInst.Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("create business failed, err: %v, data: %#v, rid: %s", err, data, kit.Rid)
		return nil, err
	}

	res, err := b.clientSet.CoreService().System().SearchPlatformSetting(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("search platform setting failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	conf := res.Data
	// create set
	setData := mapstr.MapStr{
		common.BKAppIDField:    bizID,
		common.BKInstParentStr: bizID,
		common.BKSetNameField:  conf.BuiltInSetName,
		common.BKDefaultField:  common.DefaultResSetFlag,
	}
	setInst, err := b.set.CreateSet(kit, bizID, setData)
	if err != nil {
		blog.Errorf("create set failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	setID, err := setInst.Int64(common.BKSetIDField)
	if err != nil {
		blog.Errorf("create set failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	// create module
	defaultCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("failed to search default category, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	moduleData := mapstr.MapStr{
		common.BKSetIDField:             setID,
		common.BKInstParentStr:          setID,
		common.BKAppIDField:             bizID,
		common.BKModuleNameField:        conf.BuiltInModuleConfig.IdleName,
		common.BKDefaultField:           common.DefaultResModuleFlag,
		common.BKServiceTemplateIDField: common.ServiceTemplateIDNotSet,
		common.BKSetTemplateIDField:     common.SetTemplateIDNotSet,
		common.BKServiceCategoryIDField: defaultCategory.ID,
	}

	if _, err = b.module.CreateModule(kit, bizID, setID, moduleData); err != nil {
		blog.Errorf("create module failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}

	// create fault module
	moduleData.Set(common.BKModuleNameField, conf.BuiltInModuleConfig.FaultName)
	moduleData.Set(common.BKDefaultField, common.DefaultFaultModuleFlag)
	if _, err = b.module.CreateModule(kit, bizID, setID, moduleData); err != nil {
		blog.Errorf("create fault module failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}

	// create recycle module
	moduleData.Set(common.BKModuleNameField, conf.BuiltInModuleConfig.RecycleName)
	moduleData.Set(common.BKDefaultField, common.DefaultRecycleModuleFlag)
	if _, err = b.module.CreateModule(kit, bizID, setID, moduleData); err != nil {
		blog.Errorf("create recycle module failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}

	err = b.createUserDefinedModules(kit, conf, bizID, setID, defaultCategory.ID)
	if err != nil {
		blog.Errorf("create business failed, create user module failed, err: %v, rid: %s", err, kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}
	return bizInst, nil
}

func (b *business) createUserDefinedModules(kit *rest.Kit, conf metadata.PlatformSettingConfig, bizID, setID,
	defaultCategoryID int64) error {
	for _, module := range conf.BuiltInModuleConfig.UserModules {
		// create user module
		userModuleData := mapstr.MapStr{
			common.BKSetIDField:             setID,
			common.BKInstParentStr:          setID,
			common.BKAppIDField:             bizID,
			common.BKModuleNameField:        module.Value,
			common.BKDefaultField:           common.DefaultUserResModuleFlag,
			common.BKServiceTemplateIDField: common.ServiceTemplateIDNotSet,
			common.BKSetTemplateIDField:     common.SetTemplateIDNotSet,
			common.BKServiceCategoryIDField: defaultCategoryID,
		}
		_, err := b.module.CreateModule(kit, bizID, setID, userModuleData)
		if err != nil {
			return err
		}
	}
	return nil
}

// FindBiz FindBiz
func (b *business) FindBiz(kit *rest.Kit, cond *metadata.QueryCondition) (count int,
	results []mapstr.MapStr, err error) {
	if !cond.Condition.Exists(common.BKDefaultField) {
		cond.Condition[common.BKDefaultField] = 0
	}

	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("find business by query failed, condition: %s, err: %v, rid: %s", cond, err, kit.Rid)
		return 0, nil, err
	}

	return result.Count, result.Info, err
}

// HasHosts check if this business still has hosts.
func (b *business) HasHosts(kit *rest.Kit, bizID int64) (bool, error) {

	option := []map[string]interface{}{{
		common.BKAppIDField: bizID,
	}}

	rsp, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, option)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, kit.Rid)
		return false, err
	}

	return rsp[0] != 0, nil
}

var (
	numRegex = regexp.MustCompile(`^\d+$`)
)

// GenerateAchieveBusinessName 生成归档后的业务名称
// - 业务归档的时候，自动重命名为"foo-archived"
// - 归档的时候，如果发现已经存在同名的"foo-archived", 自动在其后+1, 比如 "foo-archived-1", "foo-archived-2"
func (b *business) GenerateAchieveBusinessName(kit *rest.Kit, bizName string) (achieveName string, err error) {

	queryBusinessRequest := &metadata.QueryCondition{
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
	count, data, err := b.FindBiz(kit, queryBusinessRequest)
	if err != nil {
		return "", err
	}
	if count == 0 {
		return fmt.Sprintf("%s-archived", bizName), nil
	}
	existNums := make([]int64, 0)
	for _, item := range data {
		parts := strings.Split(util.GetStrByInterface(item[common.BKAppNameField]), fmt.Sprintf("%s-archived-",
			bizName))
		if len(parts) != 2 {
			continue
		}
		numPart := parts[1]
		if !numRegex.MatchString(numPart) {
			continue
		}
		num, err := strconv.ParseInt(numPart, 10, 64)
		if err != nil {
			blog.Errorf("GetInt64ByInterface failed, numPart: %s, err: %v, rid: %s", numPart, err, kit.Rid)
			return "", err
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

// GetBriefTopologyNodeRelation is used to get directly related business topology node information.
// As is, you can find modules belongs to a set; or you can find the set a module belongs to.
// It has rules as follows:
// 1. if src object is biz, then the destination object can be any mainline object except biz.
// 2. destination object can be biz. otherwise, src and destination object should be the neighbour.
// this api only return business topology relations.
func (b *business) GetBriefTopologyNodeRelation(kit *rest.Kit, opts *metadata.GetBriefBizRelationOptions) (
	[]*metadata.BriefBizRelations, error) {

	// validate the source and destination model is mainline model or not.
	srcDestPriority, err := b.validateMainlineObjectRule(kit, opts.SrcBizObj, opts.DestBizObj)
	if err != nil {
		blog.Errorf("check object is mainline object failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	filter := make(mapstr.MapStr)
	switch opts.SrcBizObj {
	case common.BKInnerObjIDApp:
		filter[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
		return b.genBriefTopologyNodeRelation(kit, filter, opts.DestBizObj, common.BKAppIDField,
			common.GetInstIDField(opts.DestBizObj), &opts.Page)

	case common.BKInnerObjIDSet:
		switch opts.DestBizObj {
		case common.BKInnerObjIDApp:
			filter[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDSet, common.BKSetIDField,
				common.BKAppIDField, &opts.Page)

		case common.BKInnerObjIDModule:
			filter[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDModule, common.BKSetIDField,
				common.BKModuleIDField, &opts.Page)

		default:
			// search custom level model instance with set ids. which is set's parent id list
			filter[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDSet, common.BKSetIDField,
				common.BKParentIDField, &opts.Page)
		}

	case common.BKInnerObjIDModule:
		switch opts.DestBizObj {
		case common.BKInnerObjIDApp:
			filter[common.BKModuleIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDModule, common.BKModuleIDField,
				common.BKAppIDField, &opts.Page)

		case common.BKInnerObjIDSet:
			filter[common.BKModuleIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDModule, common.BKModuleIDField,
				common.BKSetIDField, &opts.Page)

		default:
			blog.Errorf("it's not allow to find destination object %s with source module model. rid: %s",
				opts.DestBizObj, kit.Rid)
			return nil, errors.New(common.CCErrCommParamsInvalid, "dest_biz_obj")
		}

	default:
		switch opts.DestBizObj {
		case common.BKInnerObjIDApp:
			filter[common.BKInstIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, opts.SrcBizObj, common.BKInstIDField,
				common.BKAppIDField, &opts.Page)

		case common.BKInnerObjIDSet:
			filter[common.BKParentIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDSet, common.BKParentIDField,
				common.BKSetIDField, &opts.Page)

		default:
			if srcDestPriority < 0 {
				// src object is the parent of destination object
				filter[common.BKParentIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
				// src inst id bk_parent_id, dest inst id bk_inst_id, this scene is different with others.
				return b.genBriefTopologyNodeRelation(kit, filter, opts.DestBizObj, common.BKParentIDField,
					common.BKInstIDField, &opts.Page)
			}

			// destination object is the parent of src object
			filter[common.BKInstIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, opts.SrcBizObj, common.BKInstIDField,
				common.BKParentIDField, &opts.Page)
		}
	}
}

func (b *business) validateMainlineObjectRule(kit *rest.Kit, src, dest string) (int, error) {

	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	asst, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return 0, err
	}

	if len(asst.Info) <= 0 {
		return 0, fmt.Errorf("invalid biz mainline object topology")
	}

	next, idx := common.BKInnerObjIDApp, 0
	// save the mainline object with it's index with map
	rankMap := make(map[string]int)
	rankMap[next] = idx
	for _, relation := range asst.Info {
		if relation.AsstObjID == next {
			next = relation.ObjectID
			idx += 1
			rankMap[next] = idx
			continue
		}

		for _, rel := range asst.Info {
			if rel.AsstObjID == next {
				next = rel.ObjectID
				idx += 1
				rankMap[next] = idx
				break
			}
		}

	}

	// src, dest object should all be mainline object.
	srcIdx, destIdx := 0, 0
	srcIdx, exist := rankMap[src]
	if !exist {
		return 0, kit.CCError.New(common.CCErrorTopoMainlineObjectAssociationNotExist, src)
	}

	destIdx, exist = rankMap[dest]
	if !exist {
		return 0, kit.CCError.New(common.CCErrorTopoMainlineObjectAssociationNotExist, dest)
	}

	srcDestPriority := srcIdx - destIdx

	if src == common.BKInnerObjIDApp {
		// if src object is biz, then do not care about if the destination object is neighbour or not.
		return srcDestPriority, nil
	}

	// if dest object is not biz, then the src and dest object should be the neighbour.
	// if dest object is biz, we do not check the src or dest is neighbour or not.
	if (dest != common.BKInnerObjIDApp) && (math.Abs(float64(srcDestPriority)) > 1) {
		return 0, kit.CCError.New(common.CCErrCommTopoModuleNotFoundError, src)
	}

	return srcDestPriority, nil
}

func (b *business) genBriefTopologyNodeRelation(kit *rest.Kit, filter mapstr.MapStr, destObj, srcInstField,
	destInstField string, page *metadata.BasePage) ([]*metadata.BriefBizRelations, error) {

	// set sort field with destination object instance field as default.
	page.Sort = common.GetInstIDField(destObj)

	input := &metadata.QueryCondition{
		// set all the possible fields.
		Fields: []string{common.BKAppIDField, common.BKParentIDField, common.BKSetIDField, common.BKModuleIDField,
			common.BKInstIDField},
		Page:           *page,
		Condition:      filter,
		DisableCounter: true,
	}
	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, destObj, input)
	if err != nil {
		blog.ErrorJSON("get biz mainline object %s instance with filter: %s failed, err: %v, rid: %s",
			destObj, input, err, kit.Rid)
		return nil, err
	}

	relations := make([]*metadata.BriefBizRelations, 0)
	for _, one := range result.Info {
		relations = append(relations, &metadata.BriefBizRelations{
			Business: one[common.BKAppIDField],
			// source object's instance id, field different with object's type
			SrcInstID: one[srcInstField],
			// destination object's instance id, field different with object's type
			DestInstID: one[destInstField],
		})
	}

	return relations, nil
}

// GetResourcePoolBusinessID search resource pool biz id
func (b *business) GetResourcePoolBusinessID(kit *rest.Kit) (int64, error) {

	cond := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField, common.BkSupplierAccount},
		Condition: map[string]interface{}{common.BKDefaultField: common.DefaultAppFlag},
	}

	rsp, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("search resource pool id failed, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	for idx, biz := range rsp.Info {
		bizSupplierAccount, err := biz.String(common.BkSupplierAccount)
		if err != nil {
			blog.Errorf("get business supplier account failed, err: %v, rid: %s", err, kit.Rid)
			return 0, err
		}

		if kit.SupplierAccount == bizSupplierAccount {
			if !rsp.Info[idx].Exists(common.BKAppIDField) {
				blog.Errorf("bk_biz_id is non-exist, rid: %s", kit.Rid)
				// this can not be happen normally.
				return 0, kit.CCError.CCError(common.CCErrTopoAppSearchFailed)
			}
			bizID, err := rsp.Info[idx].Int64(common.BKAppIDField)
			if err != nil {
				blog.Errorf("get business id failed, err: %v, rid: %s", err, kit.Rid)
				return 0, err
			}

			return bizID, nil
		}
	}

	return 0, kit.CCError.CCError(common.CCErrTopoAppSearchFailed)
}

// updateIdleModuleConfig update admin module config.
func (b *business) updateIdleModuleConfig(kit *rest.Kit, option metadata.ModuleOption,
	config metadata.PlatformSettingConfig) error {

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

	_, err := b.clientSet.CoreService().System().UpdatePlatformSetting(kit.Ctx, kit.Header, &conf)
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

	res, err := b.clientSet.CoreService().System().SearchPlatformSetting(kit.Ctx, kit.Header)
	if err != nil {
		return err
	}

	conf, oldConf := res.Data, res.Data
	conf.BuiltInSetName = metadata.ObjectString(setName)

	_, err = b.clientSet.CoreService().System().UpdatePlatformSetting(kit.Ctx, kit.Header, &conf)
	if err != nil {
		return err
	}
	updateFields := mapstr.New()
	updateFields.Set(common.SystemSetName, setName)

	err = b.savePlatformLog(kit, metadata.AuditUpdate, &oldConf, updateFields)
	if err != nil {
		blog.Errorf("generate audit log failed, config: %v, updateFields: %v, err: %v, rid: %s", oldConf,
			updateFields, err, kit.Rid)
	}
	return nil
}

// deleteUserModuleConfig 更新删除用户自定义场景下的配置
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
		blog.Errorf("generate audit log failed config :%v, updateFields: %v err: %v, rid: %s", oldConf, updateFields,
			err, kit.Rid)
		return err
	}

	return nil
}

// checkModuleNameValid check whether the module has duplicate names in business, bizID: resource pool's business id.
func (b *business) checkModuleNameValid(kit *rest.Kit, input metadata.ModuleOption, bizID int64) error {
	queryCond := []map[string]interface{}{{
		common.BKModuleNameField: input.Name,
		common.BKDefaultField:    map[string]interface{}{common.BKDBGT: common.NormalModuleFlag},
		common.BKAppIDField:      map[string]interface{}{common.BKDBNE: bizID},
	}}

	rst, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameBaseModule,
		queryCond)
	if err != nil {
		blog.Errorf("get module count failed, filter: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return err
	}

	if rst[0] > 0 {
		blog.Errorf("module %s is duplicate, rid: %s", input.Name, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, input.Name)
	}
	return nil
}

// validateIdleModuleConfigName TODO
// Validate 判断参数是否合法，注意此处合法的标准如下:
// 1、key如果存在这个key只要名字不重复即合法
// 2、如果不存在这个key那么合法即是新增场景
// 3、true是更新场景，false是新增场景
// 4、目前系统出厂的idle 、fault、recycle只支持改名字不支持删除
func (b *business) validateIdleModuleConfigName(ctx *rest.Kit, input metadata.ModuleOption) (bool, string,
	metadata.PlatformSettingConfig, error) {

	res, err := b.clientSet.CoreService().System().SearchPlatformSetting(ctx.Ctx, ctx.Header)
	if err != nil {
		return false, "", metadata.PlatformSettingConfig{}, err
	}

	conf := res.Data
	flag := false
	var oldName string
	switch input.Key {
	case common.SystemIdleModuleKey:
		if input.Name == conf.BuiltInModuleConfig.IdleName {
			return false, "", metadata.PlatformSettingConfig{}, fmt.Errorf("idle name cannot be the same")
		}
		oldName = conf.BuiltInModuleConfig.IdleName
		flag = true

	case common.SystemFaultModuleKey:
		if input.Name == conf.BuiltInModuleConfig.FaultName {
			return false, "", metadata.PlatformSettingConfig{}, fmt.Errorf("fault name cannot be the same")
		}
		oldName = conf.BuiltInModuleConfig.FaultName
		flag = true

	case common.SystemRecycleModuleKey:
		if input.Name == conf.BuiltInModuleConfig.RecycleName {
			return false, "", metadata.PlatformSettingConfig{}, fmt.Errorf("recycle name cannot be the same")
		}
		oldName = conf.BuiltInModuleConfig.RecycleName
		flag = true

	default:
		for _, m := range conf.BuiltInModuleConfig.UserModules {
			if m.Key == input.Key {
				if m.Value == input.Name {
					return false, "", metadata.PlatformSettingConfig{},
						fmt.Errorf("user defined module name cannot be the same")
				} else {
					oldName = m.Value
					flag = true
					break
				}
			}
		}
	}
	return flag, oldName, conf, nil
}

func (b *business) validateDeleteModuleName(kit *rest.Kit, option *metadata.BuiltInModuleDeleteOption) (
	metadata.PlatformSettingConfig, error) {

	header := util.BuildHeader(common.CCSystemOperatorUserName, common.BKDefaultOwnerID)

	res, err := b.clientSet.CoreService().System().SearchPlatformSetting(kit.Ctx, header)
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
func (b *business) addUserDefinedModule(kit *rest.Kit, results []mapstr.MapStr, data metadata.ModuleOption) error {

	// create module
	ds := make([]mapstr.MapStr, 0)

	defaultCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("failed to search default category, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	for _, result := range results {
		setId, err := util.GetInt64ByInterface(result[common.BKSetIDField])
		if err != nil {
			blog.Errorf("parse set id failed, data: %v, err: %s, rid: %s", result, err, kit.Rid)
			return err
		}

		bizId, err := util.GetInt64ByInterface(result[common.BKAppIDField])
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

	_, err = b.clientSet.CoreService().Instance().CreateManyInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, d)
	if nil != err {
		blog.Errorf("failed to create object instance, error info is %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	return nil
}

func (b *business) updateModuleName(kit *rest.Kit, data metadata.ModuleOption, name string, bizID int64) error {
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

	// 在更新模块时注意只更新业务下"空闲机池"的模块，需要将资源池下面的模块排除.
	inputParams := metadata.UpdateOption{
		Data: d,
		Condition: mapstr.MapStr{
			common.BKDefaultField:    defaultFlag,
			common.BKModuleNameField: name,
			common.BKAppIDField:      mapstr.MapStr{common.BKDBNE: bizID},
		},
	}

	_, err := b.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		&inputParams)
	if err != nil {
		blog.Errorf("update inst failed to request object controller, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return nil
}

func (b *business) deleteModuleName(kit *rest.Kit, op *metadata.BuiltInModuleDeleteOption) error {
	delInstIDs := make([]int64, 0)

	queryCond := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKDefaultField:    common.DefaultUserResModuleFlag,
			common.BKModuleNameField: op.ModuleName,
		},
		Fields: []string{common.BKModuleIDField},
	}

	res, err := b.inst.FindInst(kit, common.BKInnerObjIDModule, queryCond)
	if err != nil {
		blog.Errorf("failed to find the module, query %s, error is %v, rid: %s", queryCond, err, kit.Rid)
		return err
	}

	if res.Count == 0 {
		blog.Errorf("no module founded, query %s, rid: %s", queryCond, kit.Rid)
		return fmt.Errorf("no module founded")
	}

	for _, instItem := range res.Info {
		moduleId, err := util.GetInt64ByInterface(instItem[common.BKModuleIDField])
		if err != nil {
			blog.Errorf(" decode module id failed, instItems: %s, err: %s, rid: %s", res.Info, err, kit.Rid)
			return err
		}
		delInstIDs = append(delInstIDs, moduleId)
	}

	delCond := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{common.BKDBIN: delInstIDs},
	}
	countRes, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, []map[string]interface{}{delCond})
	if err != nil {
		blog.Errorf("count host object relation failed, err: %v, cond: %v, rid: %s", err, delCond, kit.Rid)
		return err
	}

	if countRes[0] > 0 {
		return fmt.Errorf("there is a host in the module")
	}

	dc := &metadata.DeleteOption{Condition: delCond}
	_, err = b.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, dc)
	if nil != err {
		blog.Errorf("delete inst failed, err: %v, cond: %s rid: %s", err, delCond, kit.Rid)
		return err
	}

	return nil
}

// updateBusinessSet rename business idle set, except resource pool's set.
func (b *business) updateBusinessSet(kit *rest.Kit, setOption metadata.SetOption, bizID int64) error {
	// determine whether there is a set with the same name in a non-resource pool.
	if err := b.checkSetNameValid(kit, setOption, bizID); err != nil {
		return err
	}

	// update idle set name
	inputParams := metadata.UpdateOption{
		Data: map[string]interface{}{
			common.BKSetNameField: setOption.Name,
		},
		Condition: map[string]interface{}{
			common.BKDefaultField: common.DefaultResSetFlag,
			common.BKAppIDField:   map[string]interface{}{common.BKDBNE: bizID},
		},
	}

	_, err := b.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet,
		&inputParams)
	if err != nil {
		blog.Errorf("update set name failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// update platform setting.
	err = b.updateBuiltInSetConfig(kit, setOption.Name)
	if err != nil {
		blog.Errorf("update config failed, set name: %s, rid: %s", setOption.Name, kit.Rid)
		return err
	}

	return nil
}

// UpdateBusinessIdleSetOrModule 此函数用于更新全局非资源池下的空闲机池和及下面的模块，属于管理员操作，只允许将此接口提供给前端使用。
func (b *business) UpdateBusinessIdleSetOrModule(kit *rest.Kit, option *metadata.ConfigUpdateSettingOption) error {

	bizID, err := b.getResourceBizID(kit)
	if err != nil {
		return err
	}

	switch option.Type {
	case metadata.ConfigUpdateTypeSet:
		err := b.updateBusinessSet(kit, option.Set, bizID)
		if err != nil {
			return err
		}
	case metadata.ConfigUpdateTypeModule:
		err := b.updateBusinessModule(kit, option.Module, bizID)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("type error")
	}

	return nil
}

// getResourceBizID get resource pool's biz ID.
func (b *business) getResourceBizID(kit *rest.Kit) (int64, error) {

	// get resource pool's biz ID
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKDefaultField: common.DefaultAppFlag},
		Fields:    []string{common.BKAppIDField},
	}
	count, bizItems, err := b.FindBiz(kit, query)
	if err != nil {
		blog.Errorf("get resource pool's biz failed, query: %+v, err: %v, rid: %s", query, err, kit.Rid)
		return 0, err
	}
	if count > 1 || count == 0 {
		blog.Errorf("get resource pool's biz num incorrect, query: %+v, err: %v, rid: %s", query, err, kit.Rid)
		return 0, err
	}

	bizID, err := bizItems[0].Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("bizID convert to Int64 failed, err: %v, rid: %v", err, kit.Rid)
		return 0, err
	}
	return bizID, nil
}

// updateBusinessModule : 对特殊空闲机池下模块(flag:1,2,3,5)做新增或改名操作
func (b *business) updateBusinessModule(kit *rest.Kit, module metadata.ModuleOption, bizID int64) error {

	// check param is legal or not
	if err := b.checkModuleNameValid(kit, module, bizID); err != nil {
		blog.Errorf("params is illegal, bizID: %d, module: %+v, err: %v, rid: %s", bizID, module, err, kit.Rid)
		return err
	}

	flag, oldname, conf, err := b.validateIdleModuleConfigName(kit, module)
	if err != nil {
		blog.Errorf("params is illegal, bizID: %d, module: %+v, err: %v, rid: %s", bizID, module, err, kit.Rid)
		return err
	}

	if flag {
		// add normal module or rename normal module.
		err := b.updateModuleName(kit, module, oldname, bizID)
		if err != nil {
			return err
		}
	} else {
		// find idle set list
		query := &metadata.QueryCondition{
			Condition: map[string]interface{}{
				common.BKDefaultField: common.DefaultResSetFlag,
				common.BKAppIDField:   map[string]interface{}{common.BKDBNE: bizID},
			},
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
			Fields: []string{common.BKSetIDField, common.BKAppIDField},
		}
		results, err := b.getIdleSetList(kit, query)
		if err != nil {
			return err
		}
		err = b.addUserDefinedModule(kit, results, module)
		if err != nil {
			return err
		}
	}

	// update platform config.
	err = b.updateIdleModuleConfig(kit, module, conf)
	if err != nil {
		blog.Errorf("update module config failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func (b *business) checkSetNameValid(kit *rest.Kit, setOption metadata.SetOption, bizID int64) error {
	queryCond := []map[string]interface{}{
		{
			common.BKSetNameField: setOption.Name,
			common.BKAppIDField:   map[string]interface{}{common.BKDBNE: bizID},
		},
	}

	rst, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameBaseSet,
		queryCond)
	if err != nil {
		blog.Errorf("get duplicate set count failed, filter: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return err
	}

	if rst[0] > 0 {
		blog.Errorf("set %s is duplicate, rid: %s", setOption.Name, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, setOption.Name)
	}
	return nil
}

func (b *business) getIdleSetList(kit *rest.Kit, querySet *metadata.QueryCondition) ([]mapstr.MapStr, error) {

	res, err := b.inst.FindInst(kit, common.BKInnerObjIDSet, querySet)
	if err != nil {
		blog.Errorf("find set failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if res.Count == 0 {
		blog.Errorf("no set founded, rid: %s", kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}
	return res.Info, nil
}

// DeleteBusinessGlobalUserModule 删除用户自定义的空闲类模块，此接口只允许提供给前端，用于管理员使用
func (b *business) DeleteBusinessGlobalUserModule(kit *rest.Kit, option *metadata.BuiltInModuleDeleteOption) error {

	// step 1: check param is legal or not
	conf, err := b.validateDeleteModuleName(kit, option)
	if err != nil {
		blog.Errorf("delete global user module fail, params is illegal config: %v, err: %v, rid: %s", option, err,
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
		blog.Errorf("fail to update admin config err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}
