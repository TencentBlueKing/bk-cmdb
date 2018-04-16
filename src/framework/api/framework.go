package api

import (
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/publisher"
	"configcenter/src/framework/core/timer"
	"context"
	"sync"
	"time"
)

// Framework contains the basic framework data and the publisher client used to publis events.
type Framework struct {
	workerLock sync.RWMutex
	workers    MapWorker
	context    context.Context
	cancel     context.CancelFunc
	Timer      timer.Timer
	Publisher  publisher.Publisher
}

/* the following is private methods.*/

// executeWorker start the worker
func (cli *Framework) executeWorker(worker *wrapWorker) {

	log.Infof("the worker(%s) will to run", worker.Description())

	if err := worker.Run(cli); nil != err {
		log.Errorf("the worker(%s) exit from business cycle, the error is %s", worker.Description(), err.Error())
		worker.SetStatus(ExceptionExitStatus)
		return
	}

	log.Infof("the worker(%s) normal exit", worker.Description())
	worker.SetStatus(StoppedStatus)
}

// stop used to stop the business cycles.
func (cli *Framework) stop() error {

	// stop the main context
	cli.cancel()

	// stop the all workers
	cli.workerLock.Lock()
	for _, worker := range cli.workers {
		worker.Stop()
	}
	cli.workerLock.Unlock()

	return nil
}

/* the following is public methods.*/

// Context return the context of the framework
func (cli *Framework) Context() context.Context {
	return cli.context
}

// AddWorker add a new Worker
func (cli *Framework) AddWorker(worker Worker) WorkerKey {

	key := makeWorkerKey()

	cli.workerLock.Lock()
	cli.workers[key] = &wrapWorker{worker: worker, status: NormalStatus}
	cli.workerLock.Unlock()

	return key
}

// RemoveWorker remove the Worker by a WorkerKey
func (cli *Framework) RemoveWorker(key WorkerKey) bool {

	cli.workerLock.Lock()
	defer cli.workerLock.Unlock()

	return deleteWorker(cli.workers, key)
}

// SetLoger set a new logger, it will replace the inner one.
func (cli *Framework) SetLoger(loger log.Loger) {
	log.SetLoger(loger)
}

// Run start the business cycle until the stop method is called.
func (cli *Framework) Run() {

	// start all existing workers
	for _, worker := range cli.workers {
		go cli.executeWorker(worker)
	}

	// check the stat of the worker regularly, and start it if there is any new
	for {
		select {
		case <-cli.context.Done():
			log.Info("will exit from framework main business cycle")
			goto end
		case <-time.After(time.Second * 10):

			cli.workerLock.RLock()

			// scan the all workers and restart the stoped worker
			for _, worker := range cli.workers {
				switch worker.GetStatus() {
				case NormalStatus:
					go cli.executeWorker(worker)
				case WaitingToRunStatus:
					go cli.executeWorker(worker)
				case RunningStatus:
					// pass
				case StoppingStatus:
					// pass
				case StoppedStatus:
					// pass
				case ExceptionExitStatus:
					go cli.executeWorker(worker)
				default:
					log.Fatalf("unknown the worker status (%d)", worker.GetStatus())
				}
			}

			cli.workerLock.RUnlock()
		}
	}

end:
	log.Info("finish the framework main business cycle")
}
