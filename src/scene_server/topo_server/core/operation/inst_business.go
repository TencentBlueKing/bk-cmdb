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
	DeleteBusiness(params types.ContextParams, obj model.Object, cond condition.Condition) error
	FindBusiness(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error)
	UpdateBusiness(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition) error
}

// NewBusinessOperation create a business instance
func NewBusinessOperation(set SetOperationInterface, module ModuleOperationInterface, client apimachinery.ClientSetInterface, inst InstOperationInterface) BusinessOperationInterface {
	return &business{
		clientSet: client,
		set:       set,
		module:    module,
		inst:      inst,
	}
}

type business struct {
	clientSet apimachinery.ClientSetInterface
	inst      InstOperationInterface
	set       SetOperationInterface
	module    ModuleOperationInterface
}

func (b *business) CreateBusiness(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {

	data.Set(common.BKDefaultField, 0)
	data.Set(common.BKOwnerIDField, params.SupplierAccount)
	bizInst, err := b.inst.CreateInst(params, obj, data)
	if nil != err {
		blog.Errorf("[opeartion-biz] failed to create business, error info is %s", err.Error())
		return bizInst, err
	}

	// create set

	// create module

	return nil, nil
}

func (b *business) DeleteBusiness(params types.ContextParams, obj model.Object, cond condition.Condition) error {
	return nil
}

func (b *business) FindBusiness(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error) {
	return 0, nil, nil
}

func (b *business) UpdateBusiness(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition) error {
	return nil
}
