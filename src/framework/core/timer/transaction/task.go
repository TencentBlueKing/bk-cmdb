package transaction

// TaskKey the task identifier
type TaskKey string

// task implements the Task interface
type task struct {
	key TaskKey
}

// Run the task of actual execution logic code.
func (cli *task) Run() error {
	return nil
}
