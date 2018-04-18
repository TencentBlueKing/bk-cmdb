package output

type customWrapper struct {
	name    string
	runFunc func(data MapStr) error
}

// Name the Inputer description.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *customWrapper) Name() string {
	return cli.name
}

// Run the output main loop. This should block until singnalled to stop by invocation of the Stop() method.
func (cli *customWrapper) Put(data MapStr) error {
	return cli.runFunc(data)
}

// Stop is the invoked to signal that the Run() method should its execution.
// It will be invoked at most once.
func (cli *customWrapper) Stop() error {
	// only compatible with the Outputer interface
	return nil
}
