package api

import (
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/manager"
	"configcenter/src/framework/core/output"
	"context"
)

//  mgr the global variable for the manager
var mgr *manager.Manager

// Init init the framework
func Init() {

	ctx, cancel := context.WithCancel(context.Background())

	// create Framework
	mgr = manager.New()

	/** initialize the default configuration */

	// set outputer manager
	mgr.OutputerMgr = output.New()

	// set inputer manager
	mgr.InputerMgr = input.New()

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
