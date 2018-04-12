package worker

import (
	"configcenter/src/framework/api"
)

type WorkerOne struct {
}

// Describtion the worker description.
// This information will be printed when the Worker is abnormal, which is convenient for debugging.
func (cli *WorkerOne) Description() string {
	return "worker_one_description"
}

// Run the input main loop. This should block until singnalled to stop by invocation of the Stop() method.
func (cli *WorkerOne) Run(fr *api.Framework) error {
	return nil
}

// Stop is the invoked to signal that the Run() method should its execution.
// It will be invoked at most once.
func (cli *WorkerOne) Stop() {

}

type WorkerTwo struct {
}

// Description the worker description.
// This information will be printed when the Worker is abnormal, which is convenient for debugging.
func (cli *WorkerTwo) Description() string {
	return "worker_two_description"
}

// Run the input main loop. This should block until singnalled to stop by invocation of the Stop() method.
func (cli *WorkerTwo) Run(fr *api.Framework) error {
	return nil
}

// Stop is the invoked to signal that the Run() method should its execution.
// It will be invoked at most once.
func (cli *WorkerTwo) Stop() {

}
