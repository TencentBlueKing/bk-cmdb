package inst

import (
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

// FieldName the field name
type FieldName string

// Topo the inst topo structure
type Topo interface {
	Pre() Inst
	Next() Inst
}

// Iterator the iterator interface for the Inst
type Iterator interface {
	Next() (Inst, error)
}

// Inst the inst interface
type Inst interface {
	types.Saver

	GetModel() model.Model

	IsMainLine() bool

	GetAssociationModels() ([]model.Model, error)

	GetInstID() int
	GetInstName() string

	SetValue(key string, value interface{}) error
	GetValues() (types.MapStr, error)

	GetAssociationsByModleID(modleID string) ([]Inst, error)
	GetAllAssociations() (map[model.Model][]Inst, error)

	SetParent(parentInstID int) error
	GetParent() ([]Topo, error)
	GetChildren() ([]Topo, error)
}
