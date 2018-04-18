package api

import (
	"configcenter/src/framework/core/input"
	"time"
)

// AddInputerAndExecuteOnce execute a non-blocking inputer, only execute once
func AddInputerAndExecuteOnce(inputer input.Inputer) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(inputer), nil
}

// AddInputerAndExecuteTiming  regularly execute a non-blocking inputer
func AddInputerAndExecuteTiming(inputer input.Inputer, duration time.Duration) error {

	return nil
}

// AddInputerAndExecuteLoop block to execute a  inputer
func AddInputerAndExecuteLoop(inputer input.Inputer) error {
	return nil
}

// AddInputerAndExecuteTransaction execute a non-blocking inputer as a transaction
func AddInputerAndExecuteTransaction(inputer input.Inputer) error {

	return nil
}

// AddInputerAndExecuteTimingTransaction execute a non-blocking inputer as a timing transaction
func AddInputerAndExecuteTimingTransaction(inputer input.Inputer, duration time.Duration) error {

	return nil
}
