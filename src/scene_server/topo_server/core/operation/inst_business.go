package operation

import (
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

// NewBusiness create a business instance
func NewBusiness() BusinessOperationInterface {
	return &business{}
}

type business struct {
}

func (b *business) CreateBusiness(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {
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
