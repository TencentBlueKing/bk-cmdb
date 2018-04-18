package input

import (
	"configcenter/src/framework/core/log"
)

// executeInputer start the Inputer
func (cli *manager) executeInputer(Inputer *wrapInputer) {

	log.Infof("the Inputer(%s) will to run", Inputer.Name())

	if err := Inputer.Run(); nil != err {
		log.Errorf("the Inputer(%s) exit from business cycle, the error is %s", Inputer.Name(), err.Error())
		Inputer.SetStatus(ExceptionExitStatus)
		return
	}

	log.Infof("the Inputer(%s) normal exit", Inputer.Name())
	Inputer.SetStatus(StoppedStatus)
}
