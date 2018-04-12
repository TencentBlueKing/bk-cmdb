package api

// WorkerStatus the Work status type definition.
type WorkerStatus int

// Worker status type
const (

	// NormalStatus Worker just created
	NormalStatus WorkerStatus = iota
	// WaitingToRunStatus Worker is waiting to run
	WaitingToRunStatus
	// RunningStatus Worker is running
	RunningStatus
	// StoppingStatus Worker normal exiting
	StoppingStatus
	// StoppedStatus Worker normal exit
	StoppedStatus
	// ExceptionExitStatus Worker abnormal exit
	ExceptionExitStatus
)

// WorkerKey worker key definition
type WorkerKey string

// MapWorker worker object
type MapWorker map[WorkerKey]*wrapWorker

// Worker is the interface that must be implemented by every worker.
type Worker interface {

	// Description the worker description.
	// This information will be printed when the Worker is abnormal, which is convenient for debugging.
	Description() string

	// Run the input main loop. This should block until singnalled to stop by invocation of the Stop() method.
	Run(fr *Framework) error

	// Stop is the invoked to signal that the Run() method should its execution.
	// It will be invoked at most once.
	Stop()
}

// wrapWorker the Worker wrapper
type wrapWorker struct {
	status WorkerStatus
	worker Worker
}

func (cli *wrapWorker) SetStatus(status WorkerStatus) {
	cli.status = status
}

func (cli *wrapWorker) GetStatus() WorkerStatus {
	return cli.status
}

func (cli *wrapWorker) Description() string {
	return cli.worker.Description()
}

func (cli *wrapWorker) Run(fr *Framework) error {
	return cli.worker.Run(fr)
}

func (cli *wrapWorker) Stop() {
	cli.worker.Stop()
}
