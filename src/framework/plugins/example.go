package example

import (
	"configcenter/src/framework/api"
	"configcenter/src/framework/core/types"
	"fmt"
	"time"
)

func init() {

	_, sender, _ := api.CreateCustomOutputer("example_output", func(data types.MapStr) error {
		fmt.Println("outputer:", data)
		return nil
	})

	// api.RegisterInputer(target, sender, nil)
	api.RegisterTimingInputer(target, sender, time.Second*5, nil)
}

var target = &myInputer{}

type myInputer struct {
}

// Description the Inputer description.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *myInputer) Name() string {
	return "name_myinputer"
}

// Run the input should not be blocked
func (cli *myInputer) Run() interface{} {

	return types.MapStr{
		"test": "outputer",
		"hoid": "",
	}

}

func (cli *myInputer) Stop() error {
	return nil
}
