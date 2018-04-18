package api

import (
	"configcenter/src/framework/core/input"
	"time"
)

// AddInputerAndExecuteOnce execute a non-blocking inputer, only execute once
func AddInputerAndExecuteOnce(inputer input.Inputer) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: false,
		Target:  inputer,
		Kind:    input.ExecuteOnce,
	}), nil
}

// AddInputerAndExecuteTiming  regularly execute a non-blocking inputer
func AddInputerAndExecuteTiming(inputer input.Inputer, duration time.Duration) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: false,
		Target:  inputer,
		Kind:    input.ExecuteTiming,
	}), nil
}

// AddInputerAndExecuteLoop block to execute a  inputer
func AddInputerAndExecuteLoop(inputer input.Inputer) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: true,
		Target:  inputer,
		Kind:    input.ExecuteLoop,
	}), nil
}

// AddInputerAndExecuteTransaction execute a non-blocking inputer as a transaction
func AddInputerAndExecuteTransaction(inputer input.Inputer) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: false,
		Target:  inputer,
		Kind:    input.ExecuteTransaction,
	}), nil
}

// AddInputerAndExecuteTimingTransaction execute a non-blocking inputer as a timing transaction
func AddInputerAndExecuteTimingTransaction(inputer input.Inputer, duration time.Duration) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: false,
		Target:  inputer,
		Kind:    input.ExecuteTimingTransaction,
	}), nil
}
