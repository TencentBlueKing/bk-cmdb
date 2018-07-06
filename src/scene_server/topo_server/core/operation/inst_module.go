package operation

import (
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ModuleOperationInterface module operation methods
type ModuleOperationInterface interface {
	CreateModule(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error)
	DeleteModule(params types.ContextParams, obj model.Object, cond condition.Condition) error
	FindModule(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error)
	UpdateModule(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition) error
}

// NewModule create a new module
func NewModule() ModuleOperationInterface {
	return &module{}
}

type module struct {
}

func (m *module) CreateModule(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {
	return nil, nil
}

func (m *module) DeleteModule(params types.ContextParams, obj model.Object, cond condition.Condition) error {
	return nil
}

func (m *module) FindModule(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error) {
	return 0, nil, nil
}

func (m *module) UpdateModule(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition) error {
	return nil
}
