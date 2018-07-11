package operation

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// BusinessOperationInterface business operation methods
type BusinessOperationInterface interface {
	CreateBusiness(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error)
	DeleteBusiness(params types.ContextParams, obj model.Object, bizID int64) error
	FindBusiness(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error)
	UpdateBusiness(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID int64) error

	SetProxy(set SetOperationInterface, module ModuleOperationInterface, inst InstOperationInterface, obj ObjectOperationInterface)
}

// NewBusinessOperation create a business instance
func NewBusinessOperation(client apimachinery.ClientSetInterface) BusinessOperationInterface {
	return &business{
		clientSet: client,
	}
}

type business struct {
	clientSet apimachinery.ClientSetInterface
	inst      InstOperationInterface
	set       SetOperationInterface
	module    ModuleOperationInterface
	obj       ObjectOperationInterface
}

func (b *business) SetProxy(set SetOperationInterface, module ModuleOperationInterface, inst InstOperationInterface, obj ObjectOperationInterface) {
	b.inst = inst
	b.set = set
	b.module = module
	b.obj = obj
}
func (b *business) CreateBusiness(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {

	data.Set(common.BKDefaultField, 0)
	data.Set(common.BKOwnerIDField, params.SupplierAccount)
	bizInst, err := b.inst.CreateInst(params, obj, data)
	if nil != err {
		blog.Errorf("[opeartion-biz] failed to create business, error info is %s", err.Error())
		return bizInst, err
	}

	bizID, err := bizInst.GetInstID()
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is %s", err.Error())
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create set
	objSet, err := b.obj.FindSingleObject(params, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	setData := mapstr.New()
	setData.Set(common.BKAppIDField, bizID)
	setData.Set(common.BKInstParentStr, bizID)
	setData.Set(common.BKSetNameField, common.DefaultResSetName)
	setData.Set(common.BKDefaultField, common.DefaultResSetFlag)
	setData.Set(common.BKOwnerIDField, params.SupplierAccount)

	setInst, err := b.set.CreateSet(params, objSet, bizID, setData)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is %s", err.Error())
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	setID, err := setInst.GetInstID()
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is %s", err.Error())
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create module
	objModule, err := b.obj.FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	moduleData := mapstr.New()
	moduleData.Set(common.BKSetIDField, setID)
	moduleData.Set(common.BKInstParentStr, setID)
	moduleData.Set(common.BKAppIDField, bizID)
	moduleData.Set(common.BKModuleNameField, common.DefaultResModuleName)
	moduleData.Set(common.BKDefaultField, common.DefaultResModuleFlag)
	moduleData.Set(common.BKOwnerIDField, params.SupplierAccount)

	_, err = b.module.CreateModule(params, objModule, bizID, setID, moduleData)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is %s", err.Error())
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create fault module
	faultModuleData := mapstr.New()
	faultModuleData.Set(common.BKSetIDField, setID)
	faultModuleData.Set(common.BKInstParentStr, setID)
	faultModuleData.Set(common.BKAppIDField, bizID)
	faultModuleData.Set(common.BKModuleNameField, common.DefaultFaultModuleName)
	faultModuleData.Set(common.BKDefaultField, common.DefaultFaultModuleFlag)
	faultModuleData.Set(common.BKOwnerIDField, params.SupplierAccount)

	_, err = b.module.CreateModule(params, objModule, bizID, setID, moduleData)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is %s", err.Error())
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	return bizInst, nil
}

func (b *business) DeleteBusiness(params types.ContextParams, obj model.Object, bizID int64) error {

	obj, err := b.obj.FindSingleObject(params, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return err
	}

	if err = b.set.DeleteSet(params, obj, bizID, nil); nil != err {
		blog.Errorf("[operation-biz] failed to delete the set, error info is %s", err.Error())
		return params.Err.New(common.CCErrTopoAppDeleteFailed, err.Error())
	}

	innerCond := condition.CreateCondition()
	innerCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	innerCond.Field(common.BKAppIDField).Eq(bizID)

	return b.inst.DeleteInst(params, obj, innerCond)
}

func (b *business) FindBusiness(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error) {
	return 0, nil, nil
}

func (b *business) UpdateBusiness(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID int64) error {

	innerCond := condition.CreateCondition()

	innerCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	innerCond.Field(common.BKAppIDField).Eq(bizID)

	return b.inst.UpdateInst(params, data, obj, innerCond)
}
