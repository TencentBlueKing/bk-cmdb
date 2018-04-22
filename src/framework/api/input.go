package api

import (
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/output"
	"configcenter/src/framework/core/types"
	"time"
)

// RegisterInputer execute a non-blocking inputer, only execute once
func RegisterInputer(inputer input.Inputer, putter output.Puter, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		Target:    inputer,
		Kind:      input.ExecuteOnce,
		Putter:    putter,
		Exception: exceptionFunc,
	}), nil
}

// RegisterTimingInputer execute a non-blocking timing inputer, only execute once
func RegisterTimingInputer(inputer input.Inputer, putter output.Puter, frequency time.Duration, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) {

	return mgr.InputerMgr.AddInputer(input.InputerParams{
		IsTiming:  true,
		Frequency: frequency,
		Target:    inputer,
		Kind:      input.ExecuteOnce,
		Putter:    putter,
		Exception: exceptionFunc,
	}), nil
}

// CreateTransaction create a common transaction
func CreateTransaction() input.Transaction {
	return mgr.InputerMgr.CreateTransaction()
}

// CreateTimingTransaction create a timing transaction
func CreateTimingTransaction(duration time.Duration) input.Transaction {
	return mgr.InputerMgr.CreateTimingTransaction(duration)
}

// CreateCommonEvent create a common event
func CreateCommonEvent(saver types.Saver) interface{} {
	return saver
}
