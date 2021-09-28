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
	"fmt"
	"regexp"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
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
	FindBiz(kit *rest.Kit, cond *metadata.QueryBusinessRequest) (count int, results []mapstr.MapStr, err error)
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

	setData := mapstr.New()
	setData.Set(common.BKAppIDField, bizID)
	setData.Set(common.BKInstParentStr, bizID)
	setData.Set(common.BKSetNameField, common.DefaultResSetName)
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
	idleModuleData.Set(common.BKModuleNameField, common.DefaultResModuleName)
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
	faultModuleData.Set(common.BKModuleNameField, common.DefaultFaultModuleName)
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
	recycleModuleData.Set(common.BKModuleNameField, common.DefaultRecycleModuleName)
	recycleModuleData.Set(common.BKDefaultField, common.DefaultRecycleModuleFlag)
	recycleModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	recycleModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	recycleModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(kit, objModule, bizID, setID, recycleModuleData)
	if nil != err {
		blog.Errorf("create business failed, create recycle module failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
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
func (b *business) FindBiz(kit *rest.Kit, cond *metadata.QueryBusinessRequest) (count int, results []mapstr.MapStr, err error) {
	if !cond.Condition.Exists(common.BKDefaultField) {
		cond.Condition[common.BKDefaultField] = 0
	}
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

// DeleteBusiness delete business instances by condition
func (b *business) DeleteBusiness(kit *rest.Kit, bizIDs []int64) error {
	for _, bizID := range bizIDs {
		if err := b.cleanBizAndRelatedResources(kit, bizID); err != nil {
			return err
		}
	}
	return nil
}

func (b *business) cleanBizAndRelatedResources(kit *rest.Kit, bizID int64) error {
	// 1. clean host
	// archived business has no host, need not clean host
	// 2. clean module/set template
	if err := b.cleanTemplate(kit, bizID); err != nil {
		return err
	}
	// 3. clean process
	if err := b.cleanProcess(kit, bizID); err != nil {
		return err
	}
	// 4. clean service instance
	if err := b.cleanServiceInstance(kit, bizID); err != nil {
		return err
	}
	// 5. clean module
	if err := b.cleanModule(kit, bizID); err != nil {
		return err
	}
	// 6. clean set
	if err := b.cleanSet(kit, bizID); err != nil {
		return err
	}
	// 7. clean biz
	if err := b.cleanBiz(kit, bizID); err != nil {
		return err
	}

	return nil
}

func (b *business) cleanTemplate(kit *rest.Kit, bizID int64) error {
	// 1. clean process template
	if err := b.cleanProcessTemplate(kit, bizID); err != nil {
		return err
	}
	// 2. clean service template
	if err := b.cleanServiceTemplate(kit, bizID); err != nil {
		return err
	}
	// 3. clean set template and set service template relation
	if err := b.cleanSetTemplate(kit, bizID); err != nil {
		return err
	}
	// 4. clean service category
	if err := b.cleanServiceCategory(kit, bizID); err != nil {
		return err
	}

	return nil
}

func (b *business) cleanProcessTemplate(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameProcessTemplate,
		Field:     common.BKFieldID,
		Filter:    mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}
	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get process template ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("process template ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	if err := b.clientSet.CoreService().Process().DeleteProcessTemplateBatch(kit.Ctx, kit.Header, ids); err != nil {
		blog.Errorf("batch delete process template err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (b *business) cleanServiceTemplate(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameServiceTemplate,
		Field:     common.BKFieldID,
		Filter:    mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}
	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get service template ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("service template ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	for _, id := range ids {
		if err := b.clientSet.CoreService().Process().DeleteServiceTemplate(kit.Ctx, kit.Header, id); err != nil {
			blog.Errorf("failed to delete service template, id: %v, err: %v, rid: %s", id, err, kit.Rid)
			return err
		}
	}
	return nil
}

func (b *business) cleanSetTemplate(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameSetTemplate,
		Field:     common.BKFieldID,
		Filter:    mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}
	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get set template ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("set template ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	opt := metadata.DeleteSetTemplateOption{
		SetTemplateIDs: ids,
	}
	if err := b.clientSet.CoreService().SetTemplate().DeleteSetTemplate(kit.Ctx, kit.Header, bizID, opt); err != nil {
		blog.Errorf("batch delete set template err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (b *business) cleanServiceCategory(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameServiceCategory,
		Field:     common.BKFieldID,
		Filter:    mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}
	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get service category ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("service category ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	for _, id := range ids {
		if err := b.clientSet.CoreService().Process().DeleteServiceCategory(kit.Ctx, kit.Header, id); err != nil{
			blog.Errorf("failed to delete service category, id: %v, err: %v, rid: %s", id, err, kit.Rid)
			return err
		}
	}

	return nil
}

func (b *business) cleanProcess(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameProcessInstanceRelation,
		Field:     common.BKProcessIDField,
		Filter:    mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}
	rstDist, errDist := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDist != nil {
		blog.Errorf("get process ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt, errDist,
			kit.Rid)
		return errDist
	}

	ids, err := util.SliceInterfaceToInt64(rstDist)
	if err != nil {
		blog.Errorf("process ids to int failed, ids: %v, err: %v, rid: %s", rstDist, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	// clean process instance association
	if err := b.cleanInstAsst(kit, common.BKInnerObjIDProc, ids); err != nil {
		return err
	}

	// clean process process instance
	optDelProc := metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.BKProcessIDField: mapstr.MapStr{
				common.BKDBIN: ids,
			},
		},
	}
	rstDel, err := b.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc,
		&optDelProc)
	if err != nil {
		blog.Errorf("failed to delete process instance, ids: %v, err: %v, rid: %s", ids, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rstDel.Result {
		blog.Errorf("failed to delete process instance, ids: %v, err: %v, rid: %s,", ids, rstDel.ErrMsg, kit.Rid)
		return errors.New(rstDel.Code, rstDel.ErrMsg)
	}

	// clean process instance relation
	optDelProcInstRel := metadata.DeleteProcessInstanceRelationOption{
		BusinessID: &bizID,
		ProcessIDs: ids,
	}
	if err := b.clientSet.CoreService().Process().DeleteProcessInstanceRelation(kit.Ctx, kit.Header, optDelProcInstRel);
	err != nil {
		return err
	}

	return nil
}

func (b *business) cleanServiceInstance(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameServiceInstance,
		Field:     common.BKFieldID,
		Filter:    mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}
	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get service instance ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("service instance ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	optDel := &metadata.CoreDeleteServiceInstanceOption{
		BizID: bizID,
		ServiceInstanceIDs: ids,
	}
	if err := b.clientSet.CoreService().Process().DeleteServiceInstance(kit.Ctx, kit.Header, optDel); err != nil {
		blog.Errorf("failed to delete service instance, option: %+v, err: %v, rid: %s", optDel, err, kit.Rid)
		return err
	}

	return nil
}

func (b *business) cleanModule(kit *rest.Kit, bizID int64) error {
	return b.module.DeleteModule(kit, bizID, nil, nil)
}

func (b *business) cleanSet(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameBaseSet,
		Field:     common.BKSetIDField,
		Filter:    mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}
	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get set ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt, errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("set ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	return b.set.DeleteSet(kit, bizID, ids)
}

func (b *business) cleanBiz(kit *rest.Kit, bizID int64) error {
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}
	return b.inst.DeleteInst(kit, common.BKInnerObjIDApp, cond, false)
}

func (b *business) cleanInstAsst(kit *rest.Kit, objID string, instIDs []int64) error {
	if len(instIDs) == 0 {
		return nil
	}

	cond := &metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.BKDBOR: []mapstr.MapStr{
				{
					common.BKObjIDField: objID,
					common.BKInstIDField: mapstr.MapStr{
						common.BKDBIN: instIDs,
					},
				},
				{
					common.BKObjIDField: objID,
					common.BKAsstInstIDField: mapstr.MapStr{
						common.BKDBIN: instIDs,
					},
				},
			},
		},
	}
	rsp, err := b.clientSet.CoreService().Association().DeleteInstAssociation(kit.Ctx, kit.Header, cond)
	if nil != err {
		blog.Errorf("failed to request delete inst association, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("failed to delete inst association, err: %s, rid: %s", rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}
