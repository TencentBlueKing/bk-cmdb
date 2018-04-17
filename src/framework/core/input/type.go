package input

import (
	"configcenter/src/framework/api"
)

// InputerKey the inputer name
type InputerKey string

// Manager is the interface that must be implemented by every input manager.
type Manager interface {
	SetInputer(target api.Worker) InputerKey
}
