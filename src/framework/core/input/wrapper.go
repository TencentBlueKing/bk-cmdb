package input

// wrapInputer the Inputer wrapper
type wrapInputer struct {
	kind    InputerType
	status  InputerStatus
	inputer Inputer
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

func (cli *wrapInputer) Run() error {
	return cli.inputer.Input()
}

func (cli *wrapInputer) Stop() {
	cli.inputer.Stop()
}
