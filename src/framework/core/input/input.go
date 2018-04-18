package input

import (
	"configcenter/src/framework/core/log"
	"context"
	"sync"
	"time"
)

// manager implements the Manager interface
type manager struct {
	cancel      context.CancelFunc
	inputerLock sync.RWMutex
	inputers    MapInputer
}

func (cli *manager) AddInputer(target Inputer) InputerKey {

	key := makeInputerKey()

	cli.inputerLock.Lock()
	cli.inputers[key] = &wrapInputer{inputer: target, status: NormalStatus}
	cli.inputerLock.Unlock()

	return key
}

// RemoveInputer remove the Inputer by a InputerKey
func (cli *manager) RemoveInputer(key InputerKey) {

	cli.inputerLock.Lock()
	defer cli.inputerLock.Unlock()

	deleteInputer(cli.inputers, key)
}

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

// Stop used to stop the business cycles.
func (cli *manager) Stop() error {

	// stop the all Inputers
	cli.inputerLock.Lock()
	for _, Inputer := range cli.inputers {
		Inputer.Stop()
	}
	cli.inputerLock.Unlock()

	return nil
}

// Run start the business cycle until the stop method is called.
func (cli *manager) Run(ctx context.Context, cancel context.CancelFunc) {

	// set the cancel function
	cli.cancel = cancel

	// start all existing Inputers
	for _, Inputer := range cli.inputers {
		go cli.executeInputer(Inputer)
	}

	// check the stat of the Inputer regularly, and start it if there is any new
	for {
		select {
		case <-ctx.Done():
			log.Info("will exit from inputer main business cycle")
			goto end
		case <-time.After(time.Second * 10):

			cli.inputerLock.RLock()

			// scan the all Inputers and restart the stoped Inputer
			for _, Inputer := range cli.inputers {
				switch Inputer.GetStatus() {
				case NormalStatus:
					go cli.executeInputer(Inputer)
				case WaitingToRunStatus:
					go cli.executeInputer(Inputer)
				case RunningStatus:
					// pass
				case StoppingStatus:
					// pass
				case StoppedStatus:
					// pass
				case ExceptionExitStatus:
					go cli.executeInputer(Inputer)
				default:
					log.Fatalf("unknown the Inputer status (%d)", Inputer.GetStatus())
				}
			}

			cli.inputerLock.RUnlock()
		}
	}

end:
	log.Info("finish the inputer main business cycle")
}
