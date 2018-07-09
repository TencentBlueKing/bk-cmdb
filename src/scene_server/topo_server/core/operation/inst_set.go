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

// SetOperationInterface set operation methods
type SetOperationInterface interface {
	CreateSet(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error)
	DeleteSet(params types.ContextParams, obj model.Object, cond condition.Condition) error
	FindSet(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error)
	UpdateSet(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition) error
}

// NewSetOperation create a set instance
func NewSetOperation(client apimachinery.ClientSetInterface, inst InstOperationInterface) SetOperationInterface {
	return &set{
		clientSet: client,
		inst:      inst,
	}
}

type set struct {
	clientSet apimachinery.ClientSetInterface
	inst      InstOperationInterface
}

func (s *set) CreateSet(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {

	return s.inst.CreateInst(params, obj, data)
}

func (s *set) DeleteSet(params types.ContextParams, obj model.Object, cond condition.Condition) error {
	return s.inst.DeleteInst(params, obj, cond)
}

func (s *set) FindSet(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error) {
	return s.inst.FindInst(params, obj, cond)
}

func (s *set) UpdateSet(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition) error {
	return s.inst.UpdateInst(params, data, obj, cond)
}
