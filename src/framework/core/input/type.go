package input

import (
	"context"
)

// InputerStatus the inputer status type definition.
type InputerStatus int

// InputerKey the inputer name
type InputerKey string

// Inputer status type
const (

	// NormalStatus Inputer just created
	NormalStatus InputerStatus = iota
	// WaitingToRunStatus Inputer is waiting to run
	WaitingToRunStatus
	// RunningStatus Inputer is running
	RunningStatus
	// StoppingStatus Inputer normal exiting
	StoppingStatus
	// StoppedStatus Inputer normal exit
	StoppedStatus
	// ExceptionExitStatus Inputer abnormal exit
	ExceptionExitStatus
)

// MapInputer inputer object
type MapInputer map[InputerKey]*wrapInputer

// Manager is the interface that must be implemented by every input manager.
type Manager interface {

	// AddInputer add a new inputer
	AddInputer(target Inputer) InputerKey

	// RemoveInputer remove the Inputer by a WorkerKey
	RemoveInputer(key InputerKey)

	// Run start the business cycle until the stop method is called.
	Run(ctx context.Context, cancel context.CancelFunc)

	// Stop
	Stop() error
}

// Inputer is the interface that must be implemented by every Inputer.
type Inputer interface {

	// IsBlock true is block , false is non-blocking
	IsBlock() bool

	// Description the Inputer description.
	// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
	Name() string

	// Run the input main loop. This should block until singnalled to stop by invocation of the Stop() method.
	Run() error

	// Stop is the invoked to signal that the Run() method should its execution.
	// It will be invoked at most once.
	Stop() error
}

// wrapInputer the Inputer wrapper
type wrapInputer struct {
	status  InputerStatus
	inputer Inputer
}

func (cli *wrapInputer) SetStatus(status InputerStatus) {
	cli.status = status
}

func (cli *wrapInputer) GetStatus() InputerStatus {
	return cli.status
}

func (cli *wrapInputer) Name() string {
	return cli.inputer.Name()
}

func (cli *wrapInputer) Run() error {
	return cli.inputer.Run()
}

func (cli *wrapInputer) Stop() {
	cli.inputer.Stop()
}
