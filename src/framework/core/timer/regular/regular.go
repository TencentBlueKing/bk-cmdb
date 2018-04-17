package regular

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/timer/types"
	"context"
	"sync"
	"time"
)

// MapRegular regular timing task
type MapRegular map[types.TimerKey]*timerFunction

type timerFunction struct {
	isNew      bool
	cancelFunc context.CancelFunc
	key        types.TimerKey
	function   types.TimerFunc
	duration   time.Duration
}

type regular struct {
	taskLock sync.RWMutex
	tasks    MapRegular
}

/** the following are private methods  */

func (cli *regular) executeFunction(ctx context.Context, function *timerFunction) {
	ticker := time.NewTicker(function.duration)
	defer ticker.Stop()

	innerCtx, innerCancel := context.WithCancel(ctx)
	function.cancelFunc = innerCancel
	function.isNew = false
	for {
		select {
		case <-innerCtx.Done():
			log.Infof("finish the task(%s)", function.key)
			return
		case <-ticker.C:
			if err := function.function(); nil != err {
				log.Infof("failed to execute the task(%s), error info is %s", function.key, err.Error())
			}
		}
	}
}

func (cli *regular) run(ctx context.Context) {

	log.Info("the regular timer started")

	// start all timer tasks at the inital stage
	for key, function := range cli.tasks {
		log.Infof("the task (%s) starts on a regular basis. the duration is (%d) ", key, function.duration)
		go cli.executeFunction(ctx, function)
	}

	// check the newly added tasks on a regular basis and start it if there is a new addition.
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			log.Infof("finish the regular timer main loop")
			return
		case <-ticker.C:

			// start the new timer task
			cli.taskLock.RLock()
			for key, function := range cli.tasks {
				if function.isNew {
					log.Infof("the task (%s) starts on a regular basis. the duration is (%d) ", key, function.duration)
					go cli.executeFunction(ctx, function)
				}
			}
			cli.taskLock.RUnlock()

		}
	}
}

/** the following are public methods  */

// AddTimerFunc register the timing callback function
func (cli *regular) AddTimerFunc(timeDuration time.Duration, timerFunc types.TimerFunc) types.TimerKey {

	cli.taskLock.Lock()

	key := types.TimerKey(common.UUID())

	cli.tasks[key] = &timerFunction{
		isNew:    true,
		key:      key,
		function: timerFunc,
		duration: timeDuration,
	}

	cli.taskLock.Unlock()

	return key
}

// RemoveTimerFunc remove the timing callback function by the timer key
func (cli *regular) RemoveTimerFunc(timerKey types.TimerKey) {
	cli.taskLock.Lock()
	if item, ok := cli.tasks[timerKey]; ok {
		if nil != item.cancelFunc {
			item.cancelFunc()
		}
		delete(cli.tasks, timerKey)
	}
	cli.taskLock.Unlock()
}
