package example

import (
	"configcenter/src/framework/api"
	"configcenter/src/framework/core/types"
	"fmt"
)

func init() {

	_, sender, _ := api.CreateCustomOutputer("example_output", func(data types.MapStr) error {
		fmt.Println("outputer:", data)
		return nil
	})

	api.RegisterInputer(target, sender, nil)
}

var target = &myInputer{}

type myInputer struct {
}

// Description the Inputer description.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *myInputer) Name() string {
	return "name_myinputer"
}

// Input the input should not be blocked
func (cli *myInputer) Input() interface{} {
	fmt.Println("my_inputer")

	return types.MapStr{
		"test": "outputer",
		"hoid": "",
	}
}
