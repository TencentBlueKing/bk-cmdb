package input

import (
	"configcenter/src/framework/core/output"
)

// wrapInputer the Inputer wrapper
type wrapInputer struct {
	kind    InputerType
	status  InputerStatus
	inputer Inputer
	putter  output.Puter
}

func (cli *wrapInputer) SetStatus(status InputerStatus) {
	cli.status = status
}

func (cli *wrapInputer) GetStatus() InputerStatus {
	return cli.status
}

func (cli *wrapInputer) Name() string {
	return cli.inputer.Name()
}

func (cli *wrapInputer) Input() (interface{}, error) {
	return cli.inputer.Input()
}

func (cli *wrapInputer) Stop() {
	//cli.inputer.Stop()
}
