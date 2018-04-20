package input

import (
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/types"
)

// executeInputer start the Inputer
func (cli *manager) executeInputer(inputer *wrapInputer) {

	log.Infof("the Inputer(%s) will to run", inputer.Name())
	inputObj, err := inputer.Input()
	if nil != err {
		log.Errorf("the Inputer(%s) exit from business cycle, the error is %s", inputer.Name(), err.Error())
		inputer.SetStatus(ExceptionExitStatus)
		return
	}

	// inputer 分：事物、定时、常规实现
	switch t := inputObj.(type) {
	case types.MapStr:
		inputer.putter.Put(t)
	default:
		log.Infof("todo:need to deal:", t)
	}

	log.Infof("the Inputer(%s) normal exit", inputer.Name())
	inputer.SetStatus(StoppedStatus)
}
