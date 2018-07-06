package operation

import (
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

// NewSet create a set instance
func NewSet() SetOperationInterface {
	return &set{}
}

type set struct {
}

func (s *set) CreateSet(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {
	return nil, nil
}

func (s *set) DeleteSet(params types.ContextParams, obj model.Object, cond condition.Condition) error {
	return nil
}

func (s *set) FindSet(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error) {
	return 0, nil, nil
}

func (s *set) UpdateSet(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition) error {
	return nil
}
