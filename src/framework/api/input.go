package api

import (
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/output"
	"time"
)

// RegisterInputerAndExecuteOnce execute a non-blocking inputer, only execute once
func RegisterInputerAndExecuteOnce(inputer input.Inputer, putter output.Puter, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock:   false,
		Target:    inputer,
		Kind:      input.ExecuteOnce,
		Putter:    putter,
		Exception: exceptionFunc,
	}), nil
}

// RegisterInputerAndExecuteTiming  regularly execute a non-blocking inputer
func RegisterInputerAndExecuteTiming(inputer input.Inputer, duration time.Duration, putter output.Puter, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock:   false,
		Target:    inputer,
		Kind:      input.ExecuteTiming,
		Putter:    putter,
		Exception: exceptionFunc,
	}), nil
}

// RegisterInputerAndExecuteTransaction execute a non-blocking inputer as a transaction
func RegisterInputerAndExecuteTransaction(inputer input.Inputer, putter output.Puter, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock:   false,
		Target:    inputer,
		Kind:      input.ExecuteTransaction,
		Putter:    putter,
		Exception: exceptionFunc,
	}), nil
}

// RegisterInputerAndExecuteTimingTransaction execute a non-blocking inputer as a timing transaction
func RegisterInputerAndExecuteTimingTransaction(inputer input.Inputer, duration time.Duration, putter output.Puter, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsBlock:   false,
		Target:    inputer,
		Kind:      input.ExecuteTimingTransaction,
		Putter:    putter,
		Exception: exceptionFunc,
	}), nil
}
