package input

import (
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/types"
	"fmt"
)

// executeInputer start the Inputer
func (cli *manager) executeInputer(inputer *wrapInputer) {

	log.Infof("the Inputer(%s) will to run", inputer.Name())
	inputObj := inputer.Input()

	// inputer 分：事物、定时、常规实现
	switch t := inputObj.(type) {
	case nil:
		log.Info("return the data is nil")
	case types.MapStr:
		inputer.putter.Put(t)
	default:
		log.Infof("todo:need to deal:", t)
		if nil != inputer.exception {
			inputer.exception(t, fmt.Errorf("unkown the input data type"))
		}
	}

	log.Infof("the Inputer(%s) normal exit", inputer.Name())
	inputer.SetStatus(StoppedStatus)
}
