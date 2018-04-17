package output

// OutputerKey the output name
type OutputerKey string

// Manager is the interface that must be implemented by every output manager.
type Manager interface {
	SetOutputer(target Outputer) OutputerKey
}

// Outputer is the interface that must be implemented by every Outputer.
type Outputer interface {
}
