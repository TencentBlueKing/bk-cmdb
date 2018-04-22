package input

import (
	"configcenter/src/framework/core/output"
	"time"
)

// wrapInputer the Inputer wrapper
type wrapInputer struct {
	isTiming  bool
	frequency time.Duration
	kind      InputerType
	status    InputerStatus
	inputer   Inputer
	putter    output.Puter
	exception ExceptionFunc
}

func (cli *wrapInputer) SetStatus(status InputerStatus) {
	cli.status = status
}

func (cli *wrapInputer) GetStatus() InputerStatus {
	return cli.status
}

func (cli *wrapInputer) GetFrequency() time.Duration {
	return cli.frequency
}

func (cli *wrapInputer) Name() string {
	return cli.inputer.Name()
}

func (cli *wrapInputer) Input() interface{} {
	return cli.inputer.Input()
}

func (cli *wrapInputer) Stop() {
	//cli.inputer.Stop()
}
