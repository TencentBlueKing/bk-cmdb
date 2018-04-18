package output

import ()

// MapStr the common event data definition
type MapStr map[string]interface{}

// OutputerKey the output name
type OutputerKey string

// Manager is the interface that must be implemented by every output manager.
type Manager interface {

	// AddOutputer add a new outputer instance
	AddOutputer(target Outputer) OutputerKey

	// RemoveOutputer delete the outputer instace by the outputer key
	RemoveOutputer(key OutputerKey)

	// FetchOutputer find and return the puter by the outputer key
	FetchOutputer(key OutputerKey) Puter

	// CreateCustomOutputer create a new custom outputer
	CreateCustomOutputer(name string, run func(data MapStr) error) (OutputerKey, Puter)
}

// Puter send the data input
type Puter interface {
	// save into the storage
	Put(data MapStr) error
}

// Outputer is the interface that must be implemented by every Outputer.
type Outputer interface {

	// Name the Inputer description.
	// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
	Name() string

	// Run the output main loop. This should block until singnalled to stop by invocation of the Stop() method.
	Put(data MapStr) error

	// Stop is the invoked to signal that the Run() method should its execution.
	// It will be invoked at most once.
	Stop() error
}
