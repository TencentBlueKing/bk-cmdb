package api

import (
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/output"
	"time"
)

// AddInputerAndExecuteOnce execute a non-blocking inputer, only execute once
func AddInputerAndExecuteOnce(inputer input.Inputer, putter output.Puter) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: false,
		Target:  inputer,
		Kind:    input.ExecuteOnce,
		Putter:  putter,
	}), nil
}

// AddInputerAndExecuteTiming  regularly execute a non-blocking inputer
func AddInputerAndExecuteTiming(inputer input.Inputer, duration time.Duration, putter output.Puter) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: false,
		Target:  inputer,
		Kind:    input.ExecuteTiming,
		Putter:  putter,
	}), nil
}

// AddInputerAndExecuteLoop block to execute a  inputer
func AddInputerAndExecuteLoop(inputer input.Inputer, putter output.Puter) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: true,
		Target:  inputer,
		Kind:    input.ExecuteLoop,
		Putter:  putter,
	}), nil
}

// AddInputerAndExecuteTransaction execute a non-blocking inputer as a transaction
func AddInputerAndExecuteTransaction(inputer input.Inputer, putter output.Puter) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: false,
		Target:  inputer,
		Kind:    input.ExecuteTransaction,
		Putter:  putter,
	}), nil
}

// AddInputerAndExecuteTimingTransaction execute a non-blocking inputer as a timing transaction
func AddInputerAndExecuteTimingTransaction(inputer input.Inputer, duration time.Duration, putter output.Puter) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock: false,
		Target:  inputer,
		Kind:    input.ExecuteTimingTransaction,
		Putter:  putter,
	}), nil
}
