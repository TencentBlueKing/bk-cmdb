package transaction

import (
	"time"
)

// TranKey uniquely identifies a transaction
type TranKey string

// RollbackFunc the rollback function declare
type RollbackFunc func() error

// Task the transaction task interface definition
type Task interface {
	// Run the task of actual execution logic code.
	// This method cannot be blocked.
	Run() error
}

// Transaction definition
type Transaction interface {

	// Key return the transaction key, it uniquely identifies the transaction
	Key() TranKey

	// SetDuration set the execution frequency of the transaction
	SetDuration(duration time.Duration)

	// ForEachTask iterate through the task collection in order
	ForEachTask(taskFunc func(task Task))

	// Task return a new Task , it will be executed by the transaction
	CreateTask() Task

	// Begin mark the transaction as the preparation stage
	Begin() error

	// Commit mark the transaction as the executable stage
	Commit() error
}
