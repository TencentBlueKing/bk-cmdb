package api

import (
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/manager"
	"configcenter/src/framework/core/publisher"
	"configcenter/src/framework/core/timer"
	"context"
)

//  mgr the global variable for the manager
var mgr *manager.Manager

func init() {

	ctx, cancel := context.WithCancel(context.Background())

	// create Framework
	mgr = manager.New()

	/** initialize the default configuration */

	// set inputer manager
	mgr.InputerMgr = input.New()

	// set the timer
	mgr.Timer = timer.New(ctx)

	// set publisher manager
	mgr.Publisher = publisher.New()

	/** start the main business loop */
	go mgr.Run(ctx, cancel)

}

// UnInit destory the framework
func UnInit() error {
	defer func() {
		mgr = nil
	}()
	if nil == mgr {
		return nil
	}
	return manager.Delete(mgr)
}
