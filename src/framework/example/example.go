package example

import (
	"configcenter/src/framework/api"
	"fmt"
)

func init() {
	api.AddInputerAndExecuteOnce(target)
}

var target = &myInputer{}

type myInputer struct {
}

// IsBlock true is block , false is non-blocking
func (cli *myInputer) IsBlock() bool {
	return true
}

// Description the Inputer description.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *myInputer) Name() string {
	return "name_myinputer"
}

// Run the input main loop. This should block until singnalled to stop by invocation of the Stop() method.
func (cli *myInputer) Run() error {
	fmt.Println("my_inputer")
	return nil
}

// Stop is the invoked to signal that the Run() method should its execution.
// It will be invoked at most once.
func (cli *myInputer) Stop() error {
	return nil
}
