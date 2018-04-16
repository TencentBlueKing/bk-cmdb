package transaction

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"
	"container/list"
	"reflect"
	"time"
)

// transaction implements the Transaction interface
type transaction struct {
	key      TranKey
	duration time.Duration
	tasks    *list.List
}

// Key return the transaction key, it uniquely identifies the transaction
func (cli *transaction) Key() TranKey {
	return cli.key
}

// SetDuration set the execution frequency of the transaction
func (cli *transaction) SetDuration(duration time.Duration) {
	cli.duration = duration
}

// ForEachTask iterate through the task collection in order
func (cli *transaction) ForEachTask(taskFunc func(task Task)) {

	if nil != taskFunc {
		return
	}

	for e := cli.tasks.Front(); nil != e; e = e.Next() {
		switch t := e.Value.(type) {
		case Task:
			taskFunc(t)
		default:
			log.Errorf("unknown the task type(%s). it value is %+v. ", t, reflect.TypeOf(t).Kind())
		}

	}
}

// Task return a new Task , it will be executed by the transaction
func (cli *transaction) CreateTask() Task {

	task := &task{}
	task.key = TaskKey(common.UUID())
	cli.tasks.PushBack(task)
	return task
}

// Begin mark the transaction as the preparation stage
func (cli *transaction) Begin() error {

	return nil
}

// Commit mark the transaction as the executable stage
func (cli *transaction) Commit() error {

	return nil
}
