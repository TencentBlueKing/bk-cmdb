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
	inputerChan chan *wrapInputer
}

func (cli *manager) AddInputer(params InputerParams) InputerKey {

	key := makeInputerKey()

	target := &wrapInputer{
		frequency: params.Frequency,
		isTiming:  params.IsTiming,
		inputer:   params.Target,
		status:    NormalStatus,
		kind:      params.Kind,
		putter:    params.Putter,
		exception: params.Exception,
	}

	cli.inputerLock.Lock()
	cli.inputers[key] = target
	cli.inputerLock.Unlock()

	select {
	case cli.inputerChan <- target:
	default:
		log.Fatal("failed to puth the inputer")
	}

	return key
}

// RemoveInputer remove the Inputer by a InputerKey
func (cli *manager) RemoveInputer(key InputerKey) {

	cli.inputerLock.Lock()
	defer cli.inputerLock.Unlock()

	deleteInputer(cli.inputers, key)
}

// Stop used to stop the business cycles.
func (cli *manager) Stop() error {

	// stop the main loop
	if nil != cli.cancel {
		cli.cancel()
	}
	// stop the all Inputers
	cli.inputerLock.Lock()
	for _, inputer := range cli.inputers {
		inputer.Stop()
	}
	cli.inputerLock.Unlock()

	return nil
}

// Run start the business cycle until the stop method is called.
func (cli *manager) Run(ctx context.Context, cancel context.CancelFunc) {

	// set the cancel function
	cli.cancel = cancel

	// start all existing Inputers
	for _, inputer := range cli.inputers {
		go cli.executeInputer(ctx, inputer)
	}

	// check the stat of the Inputer regularly, and start it if there is any new
	for {
		select {
		case <-ctx.Done():
			log.Info("will exit from inputer main business cycle")
			goto end

		case target := <-cli.inputerChan:
			go cli.executeInputer(ctx, target)

		case <-time.After(time.Second * 10):

			cli.inputerLock.RLock()

			// scan the all Inputers and restart the stoped Inputer
			for _, inputer := range cli.inputers {
				switch inputer.GetStatus() {
				case NormalStatus:
					go cli.executeInputer(ctx, inputer)
				case WaitingToRunStatus:
					go cli.executeInputer(ctx, inputer)
				case RunningStatus:
					// pass
				case StoppingStatus:
					// pass
				case StoppedStatus:
					// pass
				case ExceptionExitStatus:
					go cli.executeInputer(ctx, inputer)
				default:
					log.Fatalf("unknown the Inputer status (%d)", inputer.GetStatus())
				}
			}

			cli.inputerLock.RUnlock()
		}
	}

end:
	log.Info("finish the inputer main business cycle")
}
