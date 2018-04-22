package output

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

// OutputerKey the output name
type OutputerKey string

// Puter send the data input
type Puter interface {
	// save into the storage
	Put(data types.MapStr) error
}

// Outputer is the interface that must be implemented by every Outputer.
type Outputer interface {

	// Name the Inputer description.
	// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
	Name() string

	// Run the output main loop. This should block until singnalled to stop by invocation of the Stop() method.
	Put(data types.MapStr) error

	// Stop is the invoked to signal that the Run() method should its execution.
	// It will be invoked at most once.
	Stop() error
}

// ModelOutputer the interface which used to maintence the model
type ModelOutputer interface {
	// CreateClassification create a new classification
	CreateClassification() model.Classification

	// FindClassificationsLikeName find a array of the classification by the name
	FindClassificationsLikeName(name string) (model.ClassificationIterator, error)

	// FindClassificationsByCondition find a array of the classification by the condition
	FindClassificationsByCondition(condition *common.Condition) (model.ClassificationIterator, error)
}

// CustomOutputer the interface which used to maintence the custom outputer
type CustomOutputer interface {
	// AddOutputer add a new outputer instance
	AddOutputer(target Outputer) OutputerKey

	// RemoveOutputer delete the outputer instace by the outputer key
	RemoveOutputer(key OutputerKey)

	// FetchOutputer find and return the puter by the outputer key
	FetchOutputer(key OutputerKey) Puter

	// CreateCustomOutputer create a new custom outputer
	CreateCustomOutputer(name string, run func(data types.MapStr) error) (OutputerKey, Puter)
}

// InstOutputer the inst outputer interface
type InstOutputer interface {
	// CreateInst create a instance for the model
	CreateInst(target model.Model) (inst.Inst, error)
	// FindInstsLikeName find all insts by the name
	FindInstsLikeName(target model.Model, instName string) (inst.Iterator, error)
	// FindInstsByCondition find all insts by the condition
	FindInstsByCondition(target model.Model, condition *common.Condition) (inst.Iterator, error)
}

// Manager is the interface that must be implemented by every output manager.
type Manager interface {
	// Model interface
	ModelOutputer

	// Custom outputer
	CustomOutputer

	// InstOutputer
	InstOutputer
}
