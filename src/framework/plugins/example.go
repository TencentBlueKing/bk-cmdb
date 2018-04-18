package example

import (
	"configcenter/src/framework/api"
	"configcenter/src/framework/core/output"
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
	return false
}

// Description the Inputer description.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *myInputer) Name() string {
	return "name_myinputer"
}

// Input the input main loop. This should block until singnalled to stop by invocation of the Stop() method.
func (cli *myInputer) Input() error {
	fmt.Println("my_inputer")

	_, sender, _ := api.CreateCustomOutputer("example_output", func(data output.MapStr) error {
		fmt.Println("outputer:", data)
		return nil
	})
	fmt.Println("put:", sender)
	sender.Put(output.MapStr{
		"test": "outputer",
		"hoid": "",
	})
	return nil
}

// Stop is the invoked to signal that the Run() method should its execution.
// It will be invoked at most once.
func (cli *myInputer) Stop() error {
	return nil
}
