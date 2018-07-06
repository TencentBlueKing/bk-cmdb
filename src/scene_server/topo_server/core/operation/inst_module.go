package operation

import (
	"configcenter/src/apimachinery"
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

// NewModuleOperation create a new module
func NewModuleOperation(client apimachinery.ClientSetInterface, inst InstOperationInterface) ModuleOperationInterface {
	return &module{
		clientSet: client,
		inst:      inst,
	}
}

type module struct {
	clientSet apimachinery.ClientSetInterface
	inst      InstOperationInterface
}

func (m *module) CreateModule(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {
	return m.inst.CreateInst(params, obj, data)
}

func (m *module) DeleteModule(params types.ContextParams, obj model.Object, cond condition.Condition) error {
	return m.inst.DeleteInst(params, obj, cond)
}

func (m *module) FindModule(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error) {
	return m.inst.FindInst(params, obj, cond)
}

func (m *module) UpdateModule(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition) error {
	return m.UpdateModule(params, data, obj, cond)
}
